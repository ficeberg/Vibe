package controllers

import (
	"../models"
	"../utils"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

type (
	User models.User
)

var (
	conf        = models.Config{}.Init()
	TokenTTL    = conf.JWT.TokenTTL * time.Hour
	TokenLeeway = 1 * time.Minute

	tokenCache   = make(map[string]string)
	tokenCacheMu sync.Mutex
)

func (u *User) Create() error {
	sa := new(utils.SaltAuth)
	u.EncryptedPassword, u.Salt, _ = sa.Gen(u.Password)
	err := userCrud(u, "create")

	return err
}

func (u *User) Get() error {
	err := userCrud(u, "read")
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Update(username string) error {
	userPrev := User{Username: username}
	err := userCrud(&userPrev, "read")
	if err != nil {
		return err
	}

	changed := structs.Map(u)
	changedFields := structs.Names(u)
	s := reflect.ValueOf(&userPrev).Elem()

	for _, chField := range changedFields {
		// exported field
		f := s.FieldByName(chField)
		if f.IsValid() && f.CanSet() {
			if fmt.Sprintf("%s", f.Type()) == "int64" {
				f.SetInt(changed[chField].(int64))
			}
			if fmt.Sprintf("%s", f.Type()) == "string" && changed[chField].(string) != "" {
				f.SetString(changed[chField].(string))
			}
		}
	}

	//TODO: Need to limit what could be updated by role

	err = userCrud(&userPrev, "update")
	if err != nil {
		return err
	}

	u.Username = username
	err = userCrud(u, "read")
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	err := userCrud(u, "delete")
	if err != nil {
		return err
	}
	return nil
}

func (u *User) IsPass(pw string) bool {
	err := userCrud(u, "read")
	if err != nil {
		return false
	}
	//Verify password
	sa := new(utils.SaltAuth)
	return sa.Check(pw, u.Salt, u.EncryptedPassword)
}

// generateToken returns a JWT token string. Please see the URL for details:
// http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-13#section-4.1
func (u *User) GenerateToken(username, privateKey string, ttl int) (string, error) {
	tokenCacheMu.Lock()
	defer tokenCacheMu.Unlock()

	uniqKey := u.Email + username + u.Username // neglect privateKey, its always the same
	signed, ok := tokenCache[uniqKey]
	if ok {
		return signed, nil
	}

	tknID := uuid.NewV4()

	// Identifies the expiration time after which the JWT MUST NOT be accepted
	// for processing.
	if ttl < 0 {
		ttl = int(TokenTTL)
	}

	if privateKey == "" {
		privateKey = conf.JWT.SigningKey
	}

	// Implementers MAY provide for some small leeway, usually no more than
	// a few minutes, to account for clock skew.
	leeway := TokenLeeway

	tkn := jwt.New(jwt.GetSigningMethod(conf.JWT.SigningMethod))
	tkn.Claims["iss"] = u.Username                                                  // Issuer
	tkn.Claims["sub"] = username                                                    // Subject
	tkn.Claims["aud"] = u.Email                                                     // Audience
	tkn.Claims["exp"] = time.Now().UTC().Add(time.Duration(ttl)).Add(leeway).Unix() // Expiration Time
	tkn.Claims["nbf"] = time.Now().UTC().Add(-leeway).Unix()                        // Not Before
	tkn.Claims["iat"] = time.Now().UTC().Unix()                                     // Issued At
	tkn.Claims["jti"] = tknID.String()                                              // JWT ID
	tkn.Claims["role"] = u.Role

	var err error
	signed, err = tkn.SignedString([]byte(privateKey))
	if err != nil {
		return "", errors.New("Server error: Cannot generate a token")
	}

	// cache our token
	tokenCache[uniqKey] = signed

	// cache invalidation, because we cache the token in tokenCache we need to
	// invalidate it expiration time. This was handled usually within JWT, but
	// now we have to do it manually for our own cache.
	time.AfterFunc(TokenTTL-TokenLeeway, func() {
		tokenCacheMu.Lock()
		defer tokenCacheMu.Unlock()

		delete(tokenCache, uniqKey)
	})

	return signed, nil
}

func (u *User) ParseToken(ut interface{}) map[string]interface{} {
	token := ut.(*jwt.Token)

	return token.Claims
}

type mongoConnectionDetails struct {
	password string
	hostPort string
}

func mongoConnDetailsFromCfg() *mongoConnectionDetails {
	password := ""
	host := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	port := os.Getenv("MONGO_PORT_27017_TCP_PORT")

	return &mongoConnectionDetails{
		password: password,
		hostPort: host + ":" + port,
	}
}

func userCrud(user *User, action string) error {

	mdb, err := mgo.Dial(mongoConnDetailsFromCfg().hostPort)
	if err != nil {
		return err
	}
	mdb.SetMode(mgo.Monotonic, true)
	defer mdb.Close()

	_, table := getTable("user")

	col := mdb.DB(conf.DB.Name).C(table)
	err = col.EnsureIndex(mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})
	if err != nil {
		return err
	}

	colQuerier := bson.M{}
	if user.Email != "" {
		colQuerier = bson.M{"email": user.Email}
	} else if user.Username != "" {
		colQuerier = bson.M{"username": user.Username}
	} else {
		return errors.New("BAD_KEY_INDEX")
	}

	user.Password = ""
	switch action {
	case "create":
		user.IsDisabled = false
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.LastLogin = time.Now()
		err = col.Insert(&user)
		if err != nil {
			//E11000: conflict
			return errors.New(strings.Fields(err.Error())[0])
		}
		err = col.Find(colQuerier).Sort("-timestamp").One(&user)
		if err != nil {
			return errors.New("CHECK_CREATED_ACCOUNT_FAILED")
		}
	case "read":
		err = col.Find(colQuerier).Sort("-timestamp").One(&user)
		if err != nil {
			return err
		}
	case "update":
		user.UpdatedAt = time.Now()

		m := new(utils.Marshal)
		change, err := m.S2M(user)
		if err != nil {
			return err
		}

		//Need to reset non-json fields
		change["encrypted_password"] = user.EncryptedPassword
		change["salt"] = user.Salt

		err = col.Update(colQuerier, change)
		if err != nil {
			return err
		}
	case "delete":
		_, err = col.RemoveAll(colQuerier)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTable(scene string) ([]string, string) {
	key := []string{}
	table := "user"
	switch scene {
	case "user":
		key = []string{"email"}
	case "social":
		key = []string{"provider", "data.userid"}
	case "usersocial":
		key = []string{"social"}
	}

	return key, conf.DB.Table[table]
}

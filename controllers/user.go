package controllers

import (
	"../models"
	"../utils"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
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
	u.CreatedAt, u.UpdatedAt, u.LastLogin = time.Now(), time.Now(), time.Now()
	u.IsDisabled = false
	govalidator.TagMap["role"] = govalidator.Validator(func(str string) bool {
		valid_roles := []string{
			"admin",
			"member",
			"guest",
			"bot",
			"api",
		}
		for _, r := range valid_roles {
			if r == str {
				return true
			}
		}
		return false
	})

	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}
	err = userCrud(u, "create")

	return err
}

func (u *User) Get() error {
	err := userCrud(u, "read")
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Update() error {
	orgUser := User{Username: u.Username}
	if err := orgUser.Get(); err != nil {
		return err
	}

	changed, changedFields := structs.Map(u), structs.Names(u)
	s := reflect.ValueOf(&orgUser).Elem()

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
			if fmt.Sprintf("%s", f.Type()) == "bool" {
				f.SetBool(changed[chField].(bool))
			}
		}
	}

	if u.Role != "" && (u.Role != "admin" && u.Role != "member" && u.Role != "guest" && u.Role != "bot" && u.Role != "api") {
		return errors.New("Role is incorrect")
	}

	err := userCrud(&orgUser, "update")
	if err != nil {
		return err
	}

	if err := u.Get(); err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	if err := userCrud(u, "delete"); err != nil {
		return err
	}
	return nil
}

func (u *User) IsPass(pw string) bool {
	if err := u.Get(); err != nil {
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

	if err := u.Get(); err != nil { //fetch extra data by key to fullfill token fields
		return "", err
	}

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
	claims := tkn.Claims.(jwt.MapClaims)
	claims["iss"] = u.Username                                                  // Issuer
	claims["sub"] = username                                                    // Subject
	claims["aud"] = u.Email                                                     // Audience
	claims["exp"] = time.Now().UTC().Add(time.Duration(ttl)).Add(leeway).Unix() // Expiration Time
	claims["nbf"] = time.Now().UTC().Add(-leeway).Unix()                        // Not Before
	claims["iat"] = time.Now().UTC().Unix()                                     // Issued At
	claims["jti"] = tknID.String()                                              // JWT ID
	claims["role"] = u.Role

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

	return token.Claims.(jwt.MapClaims)
}

func (u *User) GenerateBase64Token(username, privateKey string, ttl int) (string, error) {
	token, err := u.GenerateToken(username, privateKey, ttl)

	return b64.StdEncoding.EncodeToString([]byte(token)), err
}

func (u *User) ParseBase64Token(encToken string) map[string]interface{} {
	ut, err := b64.StdEncoding.DecodeString(encToken)
	if err != nil {
		return map[string]interface{}{}
	}
	token := interface{}(string(ut)).(*jwt.Token)

	return token.Claims.(jwt.MapClaims)
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
		return errors.New("Ensure Error: " + err.Error())
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

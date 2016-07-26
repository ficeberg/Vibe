package models

import (
	"github.com/markbates/goth"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// User contains details about the user of a page.
type User struct {
	ID                bson.ObjectId `json:"-" bson:"_id,omitempty"` //ObjectId().getTimestamp() can get created date
	Email             string        `json:"email" valid:email,required`
	Username          string        `json:"username" valid:"alphanum,required"`
	Password          string        `json:"password,omitempty" bson:"-" valid:"alphanum,required"`
	EncryptedPassword string        `json:"-" bson:"encrypted_password"`
	Salt              string        `json:"-" bson:"salt"`
	Role              string        `json:"role" valid:"role,required"` //sysadmin,admin,member,vip,banned
	DisplayName       string        `json:"display_name" bson:"display_name"`
	GivenName         string        `json:"given_name" bson:"given_name"`
	FamilyName        string        `json:"family_name" bson:"family_name"`
	Language          string        `json:"language"`
	Avatar            string        `json:"avatar"`
	ShortBio          string        `json:"short_bio" bson:"short_bio"`
	LongBio           string        `json:"long_bio" bson:"long_bio"`
	Country           string        `json:"country"`
	Phone             string        `json:"phone"`
	Birth             time.Time     `json:"birth"`
	Age               int64         `json:"age"`
	Gender            int64         `json:"gender"` //0=undefined;1=male;2=female;3=shemale
	Social            UserSocial    `json:"social"`
	IsDisabled        bool          `json:"is_disabled" bson:"is_disabled" valid:"boolean, required"`
	CreatedAt         time.Time     `json:"created_at" bson:"created_at" valid:"required"`
	UpdatedAt         time.Time     `json:"updated_at" bson:"updated_at" valid:"required"`
	LastLogin         time.Time     `json:"last_login" bson:"last_login" valid:"required"`
}

type UserToken struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Language string `json:"language"`
	Expire   int64  `json:"expire"`
}

func (u *User) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&u.Email: binding.Field{
			Form:     "email",
			Required: true,
		},
		&u.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
		&u.Role:        "role",
		&u.DisplayName: "display_name",
		&u.GivenName:   "given_name",
		&u.FamilyName:  "family_name",
		&u.Language:    "language",
		&u.Avatar:      "avatar",
		&u.ShortBio:    "short_bio",
		&u.LongBio:     "long_bio",
		&u.Country:     "country",
		&u.Phone:       "phone",
		&u.Birth:       "birth",
		&u.Age:         "age",
		&u.Gender:      "gender",
		&u.Social:      "social",
	}
}

// UserSocial is a place to put social details per user. These are the
// standard keys that themes will expect to have available, but can be
// expanded to any others on a per site basis
// - website
// - facebook
// - googleplus
// - twitter
// - linkedin
// - github
// - bitbucket
// - pinterest
// - instagram
// - youtube
// - skype
// - weibo
// - baidu
type UserSocial map[string]string

type UserStatus struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	Enabled        bool          `json:"enabled"`
	EmailActivated bool          `json:"email_activated"`
	PhoneActivated bool          `json:"phone_activated"`
	Login          []LoginStatus `json:"login"`
}

type LoginStatus struct {
	When     time.Time `json:"when"`
	IPv4     string    `json:"ipv4"`
	Location string    `json:"location"`
}

type Billing struct {
	CardNumber  string `json:"card_number"`
	Cid         string `json:"cid"`
	Expire      string `json:"expire"` //MMYY
	Type        string `json:"type"`
	DisplayName string `json:"display_name"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Country     string `json:"country"`
	Address     string `json:"address"`
}

type Social struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Provider string        `json:"provider"`
	Data     goth.User
}

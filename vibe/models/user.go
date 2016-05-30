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
	Email             string        `json:"email"`
	Username          string        `json:"username"`
	Password          string        `json:"password,omitempty" bson:"-"`
	EncryptedPassword string        `json:"-" bson:"encrypted_password"`
	Salt              string        `json:"-" bson:"salt"`
	Role              string        `json:"role"` //sysadmin,admin,member,vip,banned
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
	IsDisabled        bool          `json:"is_disabled" bson:"is_disabled"`
	CreatedAt         time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" bson:"updated_at"`
	LastLogin         time.Time     `json:"last_login" bson:"last_login"`
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

var Country = [...]string{
	"Afghanistan",
	"Albania",
	"Algeria",
	"American Samoa",
	"Andorra",
	"Angola",
	"Anguilla",
	"Antigua and Barbuda",
	"Argentina",
	"Armenia",
	"Aruba",
	"Australia",
	"Austria",
	"Azerbaijan",
	"The Bahamas",
	"Bahrain",
	"Bangladesh",
	"Barbados",
	"Belarus",
	"Belgium",
	"Belize",
	"Benin",
	"Bermuda",
	"Bhutan",
	"Bolivia",
	"Bosnia and Herzegovina",
	"Botswana",
	"Brazil",
	"Brunei",
	"Bulgaria",
	"Burkina Faso",
	"Myanmar",
	"Burundi",
	"Cambodia",
	"Cameroon",
	"Canada",
	"Cape Verde",
	"Central African Republic",
	"Chad",
	"Chile",
	"People's Republic of China",
	"Colombia",
	"Comoros",
	"Democratic Republic of the Congo",
	"Republic of the Congo",
	"Costa Rica",
	"C???e d'Ivoire",
	"Croatia",
	"Cuba",
	"Republic of Cyprus",
	"Czech Republic",
	"Denmark",
	"Djibouti",
	"Dominica",
	"Dominican Republic",
	"East Timor",
	"Ecuador",
	"Egypt",
	"El Salvador",
	"Equatorial Guinea",
	"Eritrea",
	"Estonia",
	"Ethiopia",
	"Fiji",
	"Finland",
	"France",
	"Faroe Islands",
	"Gabon",
	"The Gambia",
	"Georgia",
	"Germany",
	"Ghana",
	"Greece",
	"Grenada",
	"Guatemala",
	"Guinea",
	"Guinea-Bissau",
	"Guyana",
	"Haiti",
	"Honduras",
	"Hungary",
	"Iceland",
	"India",
	"Indonesia",
	"Iran",
	"Iraq",
	"Ireland",
	"Israel",
	"Italy",
	"Jamaica",
	"Japan",
	"Jordan",
	"Kazakhstan",
	"Kenya",
	"Kiribati",
	"Kuwait",
	"Kyrgyzstan",
	"Laos",
	"Latvia",
	"Lebanon",
	"Lesotho",
	"Liberia",
	"Libya",
	"Liechtenstein",
	"Lithuania",
	"Luxembourg",
	"Republic of Macedonia",
	"Madagascar",
	"Malawi",
	"Malaysia",
	"Maldives",
	"Mali",
	"Malta",
	"Marshall Islands",
	"Mauritania",
	"Mauritius",
	"Mexico",
	"Federated States of Micronesia",
	"Moldova",
	"Monaco",
	"Mongolia",
	"Montenegro",
	"Morocco",
	"Mozambique",
	"Myanmar",
	"Namibia",
	"Nauru",
	"Nepal",
	"Netherlands",
	"New Zealand",
	"Nicaragua",
	"Niger",
	"Nigeria",
	"Niue",
	"North Korea",
	"Norway",
	"Oman",
	"Pakistan",
	"Palau",
	"Palestine",
	"Panama",
	"Papua New Guinea",
	"Paraguay",
	"Peru",
	"Philippines",
	"Poland",
	"Portugal",
	"Puerto Rico",
	"Qatar",
	"Romania",
	"Russia",
	"Rwanda",
	"Saint Kitts and Nevis",
	"Saint Lucia",
	"Saint Vincent and the Grenadines",
	"Samoa",
	"San Marino",
	"S???o Tom??? and Pr???ncipe",
	"Saudi Arabia",
	"Senegal",
	"Serbia",
	"Seychelles",
	"Sierra Leone",
	"Singapore",
	"Slovakia",
	"Slovenia",
	"Solomon Islands",
	"Somalia",
	"South Africa",
	"South Korea",
	"South Sudan",
	"Spain",
	"Sri Lanka",
	"Sudan",
	"Suriname",
	"Swaziland",
	"Sweden",
	"Switzerland",
	"Syria",
	"Taiwan",
	"Tajikistan",
	"Tanzania",
	"Tatarstan",
	"Thailand",
	"Tibet",
	"Togo",
	"Tonga",
	"Trinidad and Tobago",
	"Tunisia",
	"Turkey",
	"Turkmenistan",
	"Tuvalu",
	"Uganda",
	"Ukraine",
	"United Arab Emirates",
	"United Kingdom",
	"United States",
	"Uruguay",
	"Uzbekistan",
	"Vanuatu",
	"Vatican City",
	"Venezuela",
	"Vietnam",
	"Western Sahara",
	"Yemen",
	"Zambia",
	"Zimbabwe",
}

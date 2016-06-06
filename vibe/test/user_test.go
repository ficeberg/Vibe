package controllers_test

import (
	"../controllers"
	"../utils"
	"errors"
	"testing"
	"time"
)

func UserCRUD(t *testing.T) {
	tn := time.Now()
	birth := tn.AddDate(-30, 0, 0)
	u := &controllers.User{
		Email:       "test_account@vibe.me",
		Username:    "TEST_USER_001",
		Password:    "just_a_pass_123",
		Role:        "member",
		DisplayName: "FS",
		GivenName:   "RT",
		FamilyName:  "Qin",
		Language:    "en-US",
		Country:     "US",
		Phone:       "+886987654321",
		Birth:       birth,
		Gender:      1,
		IsDisabled:  false,
	}
	sa := new(utils.SaltAuth)
	saerr := errors.New("")
	u.EncryptedPassword, u.Salt, saerr = sa.Gen(u.Password)
	if saerr != nil {
		t.Error("Password encryption failed")
	}
	if err := u.Create(); err != nil {
		t.Error("User cannot be created")
	}
	u2 := &controllers.User{
		Email:    u.Email,
		Username: u.Username,
	}
	if err := u2.Get(); err != nil {
		t.Error("User cannot get user")
	}
	if u.Phone != u2.Phone {
		t.Error("User profile is inconsistent")
	}
}

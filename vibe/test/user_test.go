package controllers_test

import (
	"../controllers"
	"../utils"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)
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
	err := u.Create()
	if assert.NotNil(err) {
		if err.Error() != "E11000" {
			t.Error("User cannot be created: " + err.Error())
		}
		u.Delete()
		u.Create()
		err = nil
	}

	u2 := &controllers.User{
		Email:    u.Email,
		Username: u.Username,
	}
	if err := u2.Get(); err != nil {
		t.Error("User cannot get user: " + err.Error())
	}
	assert.Equal(u.Phone, u2.Phone, "User profile is not consistent")

	u.GivenName = "TEST1"
	if err := u.Update(u.Username); err != nil {
		t.Error("User cannot be updated: " + err.Error())
	}
	u2.Get()
	assert.Equal(u2.GivenName, "TEST1", "User update failed")

	err = u.Delete()
	assert.Nil(err)
}

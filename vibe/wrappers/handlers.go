package wrappers

import (
	"../controllers"
	"../models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

type (
	Handlers struct {
		User controllers.User
	}
)

var (
	conf = models.Config{}.Init()
)

func (h *Handlers) Init() {}

func (h *Handlers) Get(c echo.Context) error {
	u := new(controllers.User)
	ut := u.ParseToken(c.Get("user"))
	u.Username = ut["iss"].(string)
	if err := u.Get(); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, u)
}

func (h *Handlers) Update(c echo.Context) error {
	u := new(controllers.User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := u.Update(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}

func (h *Handlers) Delete(c echo.Context) error {
	u := new(controllers.User)
	if err := c.Bind(u); err != nil {
		ut := u.ParseToken(c.Get("user"))
		u.Username = ut["iss"].(string)
		u.Get()
	}
	if err := u.Delete(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (h *Handlers) JWTCheck() middleware.JWTAuthConfig {
	return middleware.JWTAuthConfig{
		SigningKey:    []byte(conf.JWT.SigningKey),
		SigningMethod: conf.JWT.SigningMethod,
	}
}

func (h *Handlers) Login(c echo.Context) error {
	u, pw, err := getLoginName(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if !u.IsPass(pw) {
		return c.NoContent(http.StatusUnauthorized)
	}
	token, err := u.GenerateToken("", "", -1)
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *Handlers) Register(c echo.Context) error {
	u := new(controllers.User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := u.Create(); err != nil {
		switch err.Error() {
		case "E11000":
			return c.NoContent(http.StatusConflict)
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	//TODO: Mask passwords as asterisk
	return c.JSON(http.StatusCreated, u)
}

func getLoginName(c echo.Context) (*controllers.User, string, error) {
	u := new(controllers.User)
	user := &struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	pw := ""
	if err := c.Bind(user); err != nil {
		pw = c.FormValue("password")
		if c.FormValue("email") != "" {
			u.Email = c.FormValue("email")
		}
		if c.FormValue("username") != "" {
			u.Username = c.FormValue("username")
		}
		if u.Email == "" && u.Username == "" {
			return u, pw, err
		}
	} else {
		u.Username = user.Username
		u.Email = user.Email
		pw = user.Password
	}

	return u, pw, nil
}

func (h *Handlers) Social(c echo.Context) error {
	return nil
}

func (h *Handlers) Check(c echo.Context) error {
	return nil
}

func (h *Handlers) Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func (h *Handlers) TokenResolve(c echo.Context) error {
	u := new(controllers.User)
	ut := u.ParseToken(c.Get("user"))

	return c.String(http.StatusOK, "Welcome "+ut["iss"].(string)+"!")
}

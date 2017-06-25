package api

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/zengfu/web/common"
	"net/http"
)

type Info struct {
	Err     int    `json:"err"`
	User    string `json:"username"`
	Message string `json:"message"`
}

func CLeanSession(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Values["auth"] = false
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}
func Signin(c echo.Context) error {
	sess, _ := session.Get("session", c)
	username := c.FormValue("username")
	password := c.FormValue("password2")
	email := c.FormValue("email")
	sess.Values["auth"] = true
	sess.Values["username"] = username
	common.Add(username, password, email)
	sess.Save(c.Request(), c.Response())
	r, _ := json.Marshal(Info{
		Err: 0,
	})
	return c.JSON(http.StatusOK, string(r))
}
func CheckUsername(c echo.Context) error {
	username := c.QueryParam("username")
	var info Info
	if common.CheckName(username) {
		info.Err = 0
		info.User = username
	} else {
		info.Err = 1
		info.User = username
		info.Message = "the username existed"
	}
	r, _ := json.Marshal(info)
	return c.JSON(http.StatusOK, string(r))
}
func CheckSession(c echo.Context) error {
	sess, _ := session.Get("session", c)
	var info Info
	var ok2 bool
	auth, ok1 := sess.Values["auth"].(bool)
	if ok1 && auth {
		info.User, ok2 = sess.Values["username"].(string)
		if ok2 {
			info.Err = 0
		} else {
			info.Err = 1
			info.Message = "unvalid session"
		}
	} else {
		info.Err = 1
		info.Message = "unvalid session"
	}
	r, _ := json.Marshal(info)
	return c.JSON(http.StatusOK, string(r))
}
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	//remem := c.FormValue("remember")
	//fmt.Println(name, email, remem)
	a, err := common.Authenticate(username, password)
	if err != nil {
		return err
	}
	var info Info
	info.User = username
	if a {
		info.Err = 0
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
		}
		sess.Values["auth"] = true
		sess.Values["username"] = username
		//sess.Values["password"] = password
		sess.Save(c.Request(), c.Response())
		// cookie := new(http.Cookie)
		// cookie.Name = whyengineer.com
		// cookie.Value = username
		// cookie.Expires = time.Now().Add(24 * time.Hour)
		// c.SetCookie(cookie)
	} else {
		info.Err = 1
		info.Message = "auth error"
	}
	r, err := json.Marshal(info)
	if err != nil {
		return nil
	}

	return c.JSON(http.StatusOK, string(r))
}

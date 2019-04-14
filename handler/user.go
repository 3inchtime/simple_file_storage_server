package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	db "simple_file_storage_server/dbops"
	"simple_file_storage_server/util"
	"time"
)

const pwdSalt = "FK996"

//用户注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()

	username := r.Form.Get("username")
	pwd := r.Form.Get("password")

	if len(pwd) < 5 {
		w.Write([]byte("Invalid password"))
		return
	}

	if len(username) < 5 {
		w.Write([]byte("Invalid username"))
		return
	}

	encPwd := util.Sha1([]byte(pwd + pwdSalt))

	suc := db.UserSignUp(username, encPwd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}

	r.ParseForm()
	username := r.Form.Get("username")
	pwd := r.Form.Get("password")

	encPwd := util.Sha1([]byte(pwd + pwdSalt))

	pwdCheck := db.UserSignIn(username, encPwd)

	if !pwdCheck {
		w.Write([]byte("FAILED"))
		return
	}

	token := GenToken(username)

	upToken := db.UpdateToken(username, token)
	if !upToken {
		w.Write([]byte("FAILED"))
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token := util.MD5([]byte(username + ts + pwdSalt))
	token = token + ts[:8]
	return token
}

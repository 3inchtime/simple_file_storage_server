package handler

import (
	"fmt"
	"net/http"
	db "simple_file_storage_server/dbops"
	"simple_file_storage_server/util"
	"time"
)

const pwdSalt = "FK996"

//用户注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
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
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
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

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	isTokenValid := IsTokenValid(token)
	if !isTokenValid {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token := util.MD5([]byte(username + ts + pwdSalt))
	token = token + ts[:8]
	return token
}

func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}

	return true
}

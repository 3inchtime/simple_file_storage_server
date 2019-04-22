package handler

import "net/http"

//验证用户登录
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h(w, r)
		})
}

func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	return true
}

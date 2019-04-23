package main

import (
	"fmt"
	"net/http"
	"simple_file_storage_server/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadFileHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/update", handler.FileUpdateHandler)
	http.HandleFunc("/file/download", handler.FileDownloadHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/fastupload", handler.FastUploadHandler)

	http.HandleFunc("/file/mpupload/init",
		handler.HttpInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",
		handler.HttpInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete",
		handler.HttpInterceptor(handler.CompleteUploadHandler))

	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HttpInterceptor(handler.UserInfoHandler))

	// 配置静态资源处理
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}

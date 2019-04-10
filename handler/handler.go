package handler

import (
	"encoding/json"
	"filestore_server/meta"
	"filestore_server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)


// 上传文件接口
func UploadFileHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		index, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internet Error")
			return
		}
		io.WriteString(w, string(index))

	}else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Printf("Upload File Error: %s\n", err.Error())
			return
		}
		defer file.Close()

		fileLocation := "/tmp/" + head.Filename
		uploadTime := time.Now().Format("2006-01-02 15:04:05")

		newFile, err := os.Create(fileLocation)
		if err != nil {
			fmt.Printf("Create File Fail: %s\n", err.Error())
			return
		}

		fileSize, err := io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Save File Fail: %s\n", err.Error())
			return
		}

		// 保存文件元信息
		newFile.Seek(0, 0)
		fileSha1 := util.FileSha1(newFile)

		fileMeta := meta.FileMeta{
			FileSha1:fileSha1,
			FileName:head.Filename,
			Location:fileLocation,
			FileSize:fileSize,
			UploadTime:uploadTime,
		}

		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "Upload File Success")
}

func GetFileMetaHnadler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	filehash:=r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		 w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
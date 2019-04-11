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
	"strconv"
	"time"
)

// 上传文件接口
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		index, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internet Error")
			return
		}
		io.WriteString(w, string(index))

	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
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
			FileSha1:   fileSha1,
			FileName:   head.Filename,
			Location:   fileLocation,
			FileSize:   fileSize,
			UploadTime: uploadTime,
		}

		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload File Success")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limit, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas := meta.GetLastFileMetas(limit)

	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func FileDownloadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")

	file := meta.GetFileMeta(fileSha1)

	f ,err := os.Open(file.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Descrption", "attachment;filename=\""+file.FileName+"\"")
	w.Write(data)
}

func FileUpdateHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)

	curFileMeta.FileName = newFileName

	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)

	data, err := json.Marshal(curFileMeta)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")

	fileMeta := meta.GetFileMeta(fileSha1)

	err := os.Remove(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}

package handler

import (
	"fmt"
	"go-common/library/cache/redis"
	"math"
	"net/http"
	"os"
	rPool "simple_file_storage_server/cache/redis"
	dblayer "simple_file_storage_server/db"
	"simple_file_storage_server/util"
	"strconv"
	"strings"
	"time"
)

type MulitpartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))

	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
	}

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	upInfo := MulitpartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())

}

func UploadPartHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	fd, err := os.Create("/data/" + uploadID + "/" + chunkIndex)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part filed", nil).JSONBytes())
		return
	}

	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	rConn.Do("HSET", "MP_" + uploadID, "chkidx_" + chunkIndex, 1)
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}

	totalCount := 0
	chunkCount := 0

	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))

		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount ++
		}
	}

	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-1, "incaild request, complete upload error", nil).JSONBytes())
		return
	}

	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")

	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())

}
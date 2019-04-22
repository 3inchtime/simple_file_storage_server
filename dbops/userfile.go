package dbops

import (
	"fmt"
	"simple_file_storage_server/dbops/mysql"
	"time"
)

type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize string
	UploadAt string
	LastUpdated string
}

func UserFileUpload(username, filehash, filename string, filesize int64) bool{
	db := mysql.DBCon()
	stmt, err := db.Prepare("insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
		"`file_size`,`upload_at`,`status`) values (?,?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to Prepare, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		fmt.Println("Failed to insert user file, err:" + err.Error())
		return false
	}
	return true
}

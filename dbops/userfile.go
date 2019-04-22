package dbops

import (
	"fmt"
	"simple_file_storage_server/dbops/mysql"
	"time"
)

type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

func UserFileUpload(username, filehash, filename string, filesize int64) bool {
	db := mysql.DBCon()
	stmt, err := db.Prepare(
		"insert into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`,`status`) values (?,?,?,?,?,1)")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		return false
	}
	return true
}

//获取用户文件
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	db := mysql.DBCon()
	stmt, err := db.Prepare(
		"select file_sha1,file_name,file_size,upload_at," +
			"last_update from tbl_user_file where user_name=? and status!=2 limit ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(username, limit)

	if err != nil {
		return nil, err
	}

	var userFiles []UserFile

	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}

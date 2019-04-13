package dbops

import (
	"fmt"
	"simple_file_storage_server/meta"
)

func FileUpload(file meta.FileMeta) bool{
	filehash := file.FileSha1
	filename := file.FileName
	filesize := file.FileSize
	fileaddr := file.Location

	db := DBCon()
	stmt, err := db.Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false

}
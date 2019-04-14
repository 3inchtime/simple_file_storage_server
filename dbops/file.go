package dbops

import (
	"database/sql"
	"fmt"
	"simple_file_storage_server/dbops/mysql"
)

func UploadFileMetaDB(filehash string, filename string, filesize int64, fileaddr string) bool {

	db := mysql.DBCon()
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

	//文件Hash为唯一约束，插入受影响行数为0，则该文件已在数据库中
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false

}

type FileResult struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMetaDB(filehash string) (*FileResult, error) {
	db := mysql.DBCon()
	stmt, err := db.Prepare("SELECT file_sha1, file_name, file_addr, file_size FROM tbl_file WHERE file_sha1=? AND status = 1")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return nil, err
	}

	defer stmt.Close()

	fresult := FileResult{}
	err = stmt.QueryRow(filehash).Scan(&fresult.FileHash, &fresult.FileName, &fresult.FileAddr, &fresult.FileSize)
	if err != nil {
		return nil, err
	}

	return &fresult, nil
}

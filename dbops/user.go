package dbops

import (
	"fmt"
	"simple_file_storage_server/dbops/mysql"
)

func UserSignUp(username string, pwd string) bool {
	db := mysql.DBCon()
	stmt, err := db.Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, pwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("User with username: %s has been signed up", username)
		}
		return true
	}
	return false
}

func UserSignIn(username string, encpwd string) bool {
	db := mysql.DBCon()
	stmt, err := db.Prepare("select user_pwd from tbl_user where user_name=?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)

	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Printf("%s: user not found", username)
		return false
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

func UpdateToken(username string, token string) bool {
	db := mysql.DBCon()
	stmt, err := db.Prepare("replace into tbl_user_token (`user_name`, `user_token`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("Failed to query token, err:" + err.Error())
		return false
	}
	return true
}

package dbops

import (
	"fmt"
	"simple_file_storage_server/dbops/mysql"
)

func UserSignUp(username string, pwd string) bool {
	db := mysql.DBCon()
	stmt, err := db.Prepare(
		"insert into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to Prepare, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, pwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}

	rowsAffected, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	if rowsAffected > 0 {
		fmt.Println("Signup succeed")
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

type User struct {
	Username   string
	Email      string
	Phtone     string
	SignupAt   string
	LastActive string
	Status     string
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	db := mysql.DBCon()
	stmt, err := db.Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

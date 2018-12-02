package db

import (
	"database/sql"
	"fmt"

	"github.com/pingguoxueyuan/gostudy/mercury/common"
	"github.com/pingguoxueyuan/gostudy/mercury/util"
)

const (
	PasswordSalt = "HBZciU2SiSDr4uPeJ1e7qlIlMbyusQ0v"
)

func Register(user *common.UserInfo) (err error) {

	var userId int64
	sqlstr := "select user_id from user where username=?"
	fmt.Printf("db:%p user:%#v\n", DB, user)
	err = DB.Get(&userId, sqlstr, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if userId > 0 {
		err = ErrUserExists
		return
	}

	passwd := user.Password + PasswordSalt
	dbPassword := util.Md5([]byte(passwd))

	sqlstr = "insert into user(username,  password, email, user_id, sex, nickname)values(?,?,?,?,?,?)"
	_, err = DB.Exec(sqlstr, user.Username, dbPassword, user.Email,
		user.UserId, user.Sex, user.Nickname)
	return
}

func Login(user *common.UserInfo) (err error) {

	originPassword := user.Password
	sqlstr := "select username,password, user_id from user where username=?"
	fmt.Printf("db:%p user:%#v\n", DB, user)
	err = DB.Get(user, sqlstr, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = ErrUserNotExists
		return
	}

	passwd := originPassword + PasswordSalt
	originPasswordSalt := util.Md5([]byte(passwd))

	if originPasswordSalt != user.Password {
		err = ErrUserPasswordWrong
		return
	}

	return
}

package mysql

import (
	"bluebell/model"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "baobaobao123"



// QueryUserByUsername 检查指定用户名的用户是否存在
func QueryUserByUsername(username string) (err error){
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}

	if count > 0 {
		return ErrUserExist
	}

	return

}

// CreateUser 向数据库中插入一条新的用户记录
func CreateUser(user *model.User) (u *model.User, err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行数据库入库
	sqlStr := "insert into user(user_id, username, password) values(?, ?, ?)"
	_, err  = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err != nil {
		return &model.User{}, err
	}
	return user, nil
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *model.User) (err error){
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}

	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrInvalidPassword
	}
	return
}

// GetUserById 根据id查询用户信息
func GetUserByID(uid int64) (user *model.User, err error){
	user = new(model.User)
	sqlStr := "select user_id, username from user where user_id = ?"
	db.Get(user, sqlStr, uid)
	return
}
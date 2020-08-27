package mysql

import "github.com/pkg/errors"

var (
	ErrUserExist = errors.New("用户名已存在")
	ErrUserNotExist = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("用户名或密码错误")
	ErrInvalidID = errors.New("无效的ID")
)

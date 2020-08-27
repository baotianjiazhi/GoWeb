package service

import (
	"bluebell/dao/mysql"
	"bluebell/model"
)

func GetCommunityList() ([]*model.Community, error){
	// 查询数据库 查找到所有的community 并返回
	return mysql.GetCommunityList()
}


func GetCommunityDetail(id int64) (*model.CommunityDetail,error) {
	return mysql.GetCommunityByID(id)
}
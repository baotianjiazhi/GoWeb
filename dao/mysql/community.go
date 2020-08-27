package mysql

import (
	"bluebell/model"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (data []*model.Community,err error){

	sqlStr := "select community_id, community_name from community"
	if err := db.Select(&data,sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}

	return
}

func GetCommunityByID(id int64) (c *model.CommunityDetail,err error) {
	sqlStr := `select community_id, community_name, introduction, create_time from community where community_id = ?`
	c = new(model.CommunityDetail)
	if err = db.Get(c, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidID
		}
		return nil, err
	}
	return c, err
}

package mysql

import (
	"bluebell/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *model.Post) (err error){
	sqlStr := `insert into posts(post_id, title, content, author_id, community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostByID 根据ID查询单个帖子的详情数据
func GetPostByID(pid int64) (post *model.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from posts where post_id = ?`
	post = new(model.Post)
	err = db.Get(post, sqlStr, pid)
	if err != nil {
		return
	}

	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*model.Post,err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from posts ORDER BY create_time DESC limit ?,?`
	posts = make([]*model.Post, 0, 2)  // 不要写成make([]*model.Post, 2)
 	err = db.Select(&posts, sqlStr, (page-1)*size, size)
 	return
}


// 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*model.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from posts where post_id in (?) order by FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}

	query = db.Rebind(query)
	db.Select(&postList, query, args...)
	return
}

package mysql

import "bluebell/model"

func CreatePost(p *model.Post) (err error){
	sqlStr := `insert into posts(post_id, title, content, author_id, community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (post *model.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from posts where post_id = ?`
	post = new(model.Post)
	err = db.Get(post, sqlStr, pid)
	if err != nil {
		return
	}

	return
}

func GetPostList(page, size int64) (posts []*model.Post,err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from posts limit ?,?`
	posts = make([]*model.Post, 0, 2)  // 不要写成make([]*model.Post, 2)
 	err = db.Select(&posts, sqlStr, (page-1)*size, size)
 	return
}
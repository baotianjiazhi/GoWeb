package service

import (
	"bluebell/dao/mysql"
	"bluebell/model"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *model.Post) (err error) {
	// 1. 生成post_id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	return
	// 3. 返回
}


func GetPostByID(pid int64) (data *model.ApiPostDetail, err error) {

	// 查询拼接所使用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error(" mysql.GetPostByID(pid)", zap.Error(err))
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID)", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}

	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID(post.CommunityID)", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}
	data = &model.ApiPostDetail{
		AuthorName: user.Username,
		Post: post,
		CommunityDetail: community,
	}

	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*model.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	data = make([]*model.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID)", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityID)", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &model.ApiPostDetail{
			AuthorName: user.Username,
			Post: post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}
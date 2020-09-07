package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/model"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// PostListService 获取帖子列表query string
type PostListService struct {
	CommunityID int64 `json:"community_id" form:"community_id"` //可以为空
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}


func CreatePost(p *model.Post) (err error) {
	// 1. 生成post_id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
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
		AuthorName:      user.Username,
		Post:            post,
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
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(servicer *PostListService) (data []*model.ApiPostDetail, err error) {

	// 1. 去redis查询id列表
	ids, err := redis.GetPostIDInOrder(servicer.Size, servicer.Page, servicer.Order)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(servier) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 3. 根据id去数据库查询帖子详情
	// 返回的数据还要按照我给定的id的顺序
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("getposts2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for idx, post := range posts {
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
			AuthorName:      user.Username,
			VoteNumber:      votedata[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList2(servicer *PostListService) (data []*model.ApiPostDetail, err error) {
	// 1. 去redis查询id列表
	ids, err := redis.GetCommunityPostIDInOrder(servicer.Order, servicer.CommunityID, servicer.Page, servicer.Size)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(servier) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 3. 根据id去数据库查询帖子详情
	// 返回的数据还要按照我给定的id的顺序
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("getposts2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for idx, post := range posts {
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
			AuthorName:      user.Username,
			VoteNumber:      votedata[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

// GetPostListNew 将两个查询帖子列表合二为一的函数
func GetPostListNew(servicer *PostListService) (data []*model.ApiPostDetail, err error){
	if servicer.CommunityID == 0 {
		// 查询所有
		data, err = GetPostList2(servicer)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList2(servicer)
	}

	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return

}
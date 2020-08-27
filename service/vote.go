package service

import (
	"bluebell/dao/redis"
	"fmt"
	"go.uber.org/zap"
)

type VoteData struct {
	// UserID从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`                // 帖子ID
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票（1）还是反对票（-1）取消投票（0）
}



// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, servicer *VoteData) (error) {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.String("postID", servicer.PostID),
		zap.Int8("direction", servicer.Direction))
	return redis.VoteForPost(fmt.Sprintf("%d", userID), servicer.PostID, float64(servicer.Direction))
}

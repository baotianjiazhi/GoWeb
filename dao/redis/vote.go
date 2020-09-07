package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math"
	"time"
)

const(
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote = 238 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested = errors.New("不允许重复投票")
)


// 投一票加248分 86400/200 -> 200张赞成票可以给帖子续一天 -> 《redis实战》
// 投票功能

/*
direction=1时，两种情况:
	1. 之前没有投过票，现在投赞成票 --> 更新分数和投票记录
	2. 之前投反对票，现在投赞成票
direction=0时,两种情况：
	1. 之前投过赞成票，现在要取消投票
	2. 之前投过反对票，现在要取消投票
direction=-1时，两种情况：
	1. 之前没有投过票，现在投反对票
	2. 之前投赞成票，现在改投反对票

投票的限制：
每个帖子自发表起一个星期之内允许用户投票，超过一个星期就不允许投票了
	1.到期之后将redis中保存的赞成票数及反对票存储到mysql表中
	2.到期之后删除KeyPostVotedZsetPf
*/
func VoteForPost(userID, postID string, value float64) (err error) {

	// 1.判断投票限制
	// 去redis里拿发帖时间
	postTime, err := rdb.ZScore(getRedisKey(KeyPostTimeZset), postID).Result()
	zap.L().Error("rdb.ZScore(getRedisKey(KeyPostTimeZset), postID).Result() err", zap.Error(err))
	if float64(time.Now().Unix()) - postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2.更新帖子的分数
	// 先查当前用户给当前帖子的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZsetPf+postID), userID).Val()
	// 更新：如果这一次投票和上一次一样，就不需要重复投票
	if value == ov {
		return ErrVoteRepested
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	pipeline := rdb.Pipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZset), op*diff*scorePerVote, postID)
	if err != nil {
		return err
	}
	// 3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZsetPf+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZsetPf+postID), redis.Z{
			Score: value,
			Member: userID,
		}).Result()
	}
	_, err = pipeline.Exec()
	return err
}

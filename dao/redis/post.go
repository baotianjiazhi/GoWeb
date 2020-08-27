package redis

import (
	"github.com/go-redis/redis"
	"time"
)

func CreatePost(postID int64) error {
	pipeline := rdb.Pipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZset), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	}).Result()

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZset), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	}).Result()
	_, err := pipeline.Exec()
	return err
}

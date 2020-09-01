package redis

import (
	"bluebell/model"
	"fmt"
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


func GetPostIDInOrder(Size, Page int64, Order string) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求中携带的order
	key := getRedisKey(KeyPostTimeZset)
	if Order == model.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}
	// 2. 确定查询的索引起始位置
	start := (Page - 1) * Size
	end := start + Size - 1
	// 3. ZREVARANGE 查询
	fmt.Println(start, end)
	return rdb.ZRevRange(key, start, end).Result()
}
package redis

import (
	"bluebell/model"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func CreatePost(postID, communityID int64) error {
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
	// 更新 把帖子id加到社区的set中
	cKey := getRedisKey(KeyCommunitySetPf+strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 2. 确定查询的索引起始位置
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZREVARANGE 查询
	fmt.Println(start, end)
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDInOrder(Size, Page int64, Order string) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求中携带的order
	key := getRedisKey(KeyPostTimeZset)
	if Order == model.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}

	return getIDsFromKey(key, Size, Page)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {

	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZsetPf+ id)
	//	// 查找key中分数是1的元素数量
	//	v := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDInOrder 按社区查找
func GetCommunityPostIDInOrder(order string, communityID, page, size int64) ([]string, error){
	orderKey := getRedisKey(KeyPostTimeZset)
	if order == model.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZset)
	}
	// 使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据M

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPf+strconv.Itoa(int(communityID)))
	// 缓存的key
	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(communityID))
	if rdb.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsFromKey(key, page, size)
}
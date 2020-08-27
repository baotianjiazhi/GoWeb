package redis

// redis key

// redis key注意用命名空间的方式，方便业务查询和拆分
const (
	Prefix = "bluebell:"
	KeyPostTimeZset  = "post:time" // zset；帖子及发帖时间
	KeyPostScoreZset = "post:score" // zset：帖子及投票的分数
	KeyPostVotedZsetPf = "post:voted" // zset；记录用户及投票类型；参数是post_id
)

// 给redis加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
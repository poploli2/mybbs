package common

const (
	KeyUserTokenPrefix = "mybbs:user:token:"
	KeyPostTimeZSet    = "mybbs:post:time"
	KeyPostScoreZSet   = "mybbs:post:score"
	KeyPostVotedPrefix = "mybbs:post:voted:" // wobbs:post:voted:post_id 记录用户及投票的类型
)

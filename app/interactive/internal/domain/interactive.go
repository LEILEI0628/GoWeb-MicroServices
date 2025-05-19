package domain

// Interactive 总体交互计数
type Interactive struct {
	Biz   string
	BizId int64

	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`
	// 此字段是文章有无点赞或收藏
	// 也可以把这两个字段分离，作为一个单独的结构体
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

// max(发送者总速率/单一分区写入速率, 发送者总速率/单一消费者速率) + buffer

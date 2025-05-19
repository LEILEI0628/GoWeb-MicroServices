package events

import "context"

type Producer interface {
	ProduceReadEvent(ctx context.Context, evt ReadEvent) error // 具体事件而不是阅读计数+1这种
}

type ReadEvent struct {
	Uid int64
	Aid int64
}

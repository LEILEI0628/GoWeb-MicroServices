package service

import (
	"context"
	loggerx "github.com/LEILEI0628/GinPro/middleware/logger"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/events/article"
	"github.com/LEILEI0628/GoWeb/internal/repository/article"
	"time"
)

type ArticleServiceInterface interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error)
}
type ArticleService struct {
	repo repository.ArticleRepository
	// V1 依靠两个不同的 repository 来解决这种跨表，或者跨库的问题
	author   repository.ArticleAuthorRepository
	reader   repository.ArticleReaderRepository
	l        loggerx.Logger
	producer events.Producer
}

func NewArticleService(repo repository.ArticleRepository) ArticleServiceInterface {
	return &ArticleService{repo: repo}
}

func (svc *ArticleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		// 修改帖子（有Id）
		err := svc.repo.Update(ctx, article)
		return article.Id, err
	}
	return svc.repo.Create(ctx, article) // 新建帖子
}

func (a *ArticleService) Withdraw(ctx context.Context, art domain.Article) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (svc *ArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	// 制作库
	//id, err := a.repo.Create(ctx, art)
	//// 线上库
	//a.repo.SyncToLiveDB(ctx, art)
	panic("implement me")
}

func (a *ArticleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败，保存到线上库失败",
			loggerx.Int64("art_id", art.Id),
			loggerx.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败，重试彻底失败",
			loggerx.Int64("art_id", art.Id),
			loggerx.Error(err))
		// 接入你的告警系统，手工处理一下
		// 走异步，我直接保存到本地文件
		// 走 Canal
		// 打 MQ
	}
	return id, err
}
func (svc *ArticleService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			er := svc.producer.ProduceReadEvent( // 在此处发送读者阅读事件
				ctx,
				events.ReadEvent{
					// 即使需要article中的数据也不要在这里传，因为在消费的时候可能已经被修改了
					Uid: uid,
					Aid: id,
				})
			if er == nil {
				svc.l.Error("发送读者阅读事件失败")
			}
		}()
	}
	return art, err
}

func (a *ArticleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *ArticleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

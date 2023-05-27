package service

import (
	"context"

	"github.com/sjxiang/xiaohongshu-demo/internal/domain"
)


// 因为我们的业务逻辑过于简单，所以 service 和 repository 的基本一样
// 复杂的业务会使用不同的 repo 和其它组件共同组成一个整体
type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (domain.Article, error)
	Publish(ctx context.Context, article domain.Article) (domain.Article, error)
	Get(ctx context.Context, id uint64) (domain.Article, error)
}


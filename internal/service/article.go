package service

import (
	"context"

	"github.com/sjxiang/xiaohongshu-demo/internal/domain"
	"github.com/sjxiang/xiaohongshu-demo/internal/repository"
)

type service struct {
	bRepo repository.ArticleRepo
	// cRepo repository.ArticleRepo
}

func NewService(bRepo repository.ArticleRepo) ArticleService {
	return &service{
		bRepo: bRepo,
		// cRepo: cRepo,
	}
}

func (s *service) Publish(ctx context.Context, article domain.Article) (domain.Article, error) {
	// 先保存 B 端
	art, err := s.Save(ctx, article)
	if err != nil {
		return domain.Article{}, err
	}
	// 同步过去 C 端，这里因为是不同的数据库，所以你不能搞本地事务
	// 只能是考虑重试 + 监控 + 告警
	// _, err = s.cRepo.CreateAndCached(ctx, art)
	return art, err
}

// Save 在 service 层面上才会有创建或者更新的概念。repository 的职责更加单纯一点
func (s *service) Save(ctx context.Context, article domain.Article) (domain.Article, error) {
	if article.ID == 0 {
		id, err := s.bRepo.Create(ctx, article)
		if err != nil {
			return domain.Article{}, err
		}
		article.ID = id
		return article, nil
	}
	err := s.bRepo.Update(ctx, article)
	return article, err
}

// Get 这是 C 端查看
func (s *service) Get(ctx context.Context, id uint64) (domain.Article, error) {
	return s.bRepo.Get(ctx, id)
}
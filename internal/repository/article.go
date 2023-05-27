package repository

import (
	"context"
	"time"


	"github.com/sjxiang/xiaohongshu-demo/internal/domain"
	"github.com/sjxiang/xiaohongshu-demo/internal/repository/dao"
)

type ArticleRepo interface {
	// Create 创建一篇文章
	Create(ctx context.Context, article domain.Article) (uint64, error)
	
	// CreateAndCached(ctx context.Context, article domain.Article) (uint64, error)
	
	// Update 更新一篇文章
	Update(ctx context.Context, article domain.Article) error
	
	// Get 方法应该负责将 Author 也一并组装起来。
	// 这里就会有一个很重要的概念，叫做延迟加载，但是 GO 是做不了的，
	// 所以只能考虑传递标记位，或者使用新方法来控制要不要把 Author 组装起来
	Get(ctx context.Context, id uint64) (domain.Article, error)
}

func NewArticleRepo(dao dao.ArticleDAO) ArticleRepo {
	return &articleRepoImpl{
		dao:   dao,
	}
}

type articleRepoImpl struct {
	dao   dao.ArticleDAO
}

// func (a *articleRepo) CreateAndCached(ctx context.Context, article domain.Article) (uint64, error) {
// 	now := time.Now().UnixMilli()
// 	entity := dao.Article{
// 		ID:      article.ID,
// 		Title:   article.Title,
// 		Content: article.Content,
// 		Author:  article.Author.ID,
// 		Ctime:   now,
// 		Utime:   now,
// 	}
// 	// 你也可以在这里 error 再次封装一遍
// 	id, err := a.dao.Insert(ctx, entity)
// 	if err != nil {
// 		return 0, err
// 	}
// 	article.ID = id
// 	// err = a.cache.Set(ctx, article)
// 	return id, err
// }


func (a *articleRepoImpl) Create(ctx context.Context, article domain.Article) (uint64, error) {
	now := time.Now().Unix()

	entity := dao.Article{
		ID:      article.ID,
		Title:   article.Title,
		Content: article.Content,
		Author:  article.Author.ID,
		Ctime:   now,
		Utime:   now,
	}
	// 你也可以在这里 error 再次封装一遍
	return a.dao.Insert(ctx, entity)
}


func (a *articleRepoImpl) Update(ctx context.Context, article domain.Article) error {
	now := time.Now().Unix()

	entity := dao.Article{
		ID:      article.ID,
		Title:   article.Title,
		Content: article.Content,
		Author:  article.Author.ID,
		Ctime:   now,
		Utime:   now,
	}
	return a.dao.Update(ctx, entity)
}

func (a *articleRepoImpl) Get(ctx context.Context, id uint64) (domain.Article, error) {
	// res, err := a.cache.Get(ctx, id)
	// if err == nil {
	// 	return res, nil
	// }
	
	entity, err := a.dao.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	
	// 按照道理来说，这里要组装好 Author 的
	art := domain.Article{
		ID:      entity.ID,
		Title:   entity.Title,
		Content: entity.Content,
		Ctime:   time.UnixMilli(entity.Ctime),
		Utime:   time.UnixMilli(entity.Utime),
		Author: domain.Author{
			ID: entity.Author,
		},
	}
	// err = a.cache.Set(ctx, art)
	// if err != nil {
	// 	// 这个 error 你可以考虑吞掉。因为缓存虽然更新失败了，但是实际上你数据库是拿到了
	// 	// 不过实际上这个很危险，因为如果 Redis 整个崩溃了，那么数据库也扛不住压力
	// 	zap.L().Error("换成数据失败", zap.Error(err))
	// }
	return art, nil
	// return nil, nil 

}
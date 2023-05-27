package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sjxiang/xiaohongshu-demo/internal/domain"
)

type ArticleCache interface {
	// Set 理论上来说，ArticleCache 也应该有自己的 Article 定义
	// 比如说你并不需要缓存全部字段
	// 但是我们这里直接缓存全部
	Set(ctx context.Context, article domain.Article) error
	Get(ctx context.Context, id uint64) (domain.Article, error)
}

type articleRedisCache struct {
	client *redis.Client
}

func (a *articleRedisCache) Set(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	res, err := a.client.Set(ctx, fmt.Sprintf("article_%d", article.ID), string(data), time.Hour).Result()
	if res != "OK" {
		return errors.New("插入失败")
	}
	return err
}

func (a *articleRedisCache) Get(ctx context.Context, id uint64) (domain.Article, error) {
	data, err := a.client.Get(ctx, fmt.Sprintf("article_%d", id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(data, &art)
	return art, err
}

func NewArticleRedisCache(client *redis.Client) ArticleCache {
	return &articleRedisCache{
		client: client,
	}
}
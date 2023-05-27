package dao

import (
	"context"
	"time"
	 
	"gorm.io/gorm"
)

// Article 直接对应到表结构
type Article struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	Title   string `form:"title"`
	Content string `form:"content"`
	// 作者 ID
	Author  uint64 `gorm:"index,not null"`
	// 创建时间，毫秒作为单位
	Ctime   int64
	// 更新时间，毫秒作为单位
	Utime   int64
}


// ArticleDAO 有些人会在这里继续操作 domain
// 但是我更加认为 DAO 没有 domain 的概念， repository 要负责转换
type ArticleDAO interface {
	// Insert 的概念更加贴近关系型数据库，所以这里就不再是用 CREATE 这种说法了
	Insert(ctx context.Context, article Article) (uint64, error)
	Update(ctx context.Context, article Article) error
	GetByID(ctx context.Context, id uint64) (Article, error)
}

type articleGORM struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &articleGORM{db: db}
}

func (impl *articleGORM) Insert(ctx context.Context, article Article) (uint64, error) {
	err := impl.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

func (a *articleGORM) Update(ctx context.Context, article Article) error {
	// 我一般都是显式指定更新条件、更新字段也尽可能指定，绝对不依赖于默认行为
	// 默认行为对后面的维护者很不好
	// 尤其是依赖于更新非零值的特性，你看代码是不知道哪些字段是零值，哪些字段是非零值
	article.Utime = time.Now().Unix()
	return a.db.WithContext(ctx).Model(&article).Updates(article).Error
}

func (a *articleGORM) GetByID(ctx context.Context, id uint64) (Article, error) {
	var art Article
	err := a.db.WithContext(ctx).Where("id=?", id).First(&art).Error
	return art, err
}
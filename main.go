package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db := InitDB()

	server := gin.Default()

	server.LoadHTMLGlob("template/*")

	article := server.Group("/article")
	
	{
		// GET /article/new 编辑页面
		article.GET("/new", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "write_article.html", nil)
		})

		// POST /article/new 发表帖子
		article.POST("/new", func(ctx *gin.Context) {
			var article Article
			if err := ctx.Bind(&article); err != nil {
				log.Println(err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "输入有误",
				})
				return
			}

			now := time.Now().Unix()
			article.Ctime = now
			article.Utime = now
			article.Author = "sjxiang"  // 理论上，应该从 Session 里读取，未实现，故写死
			
			if err := db.Create(&article).Error; err != nil {
				log.Println(err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "系统故障",
				})
				return
			}

			ctx.Redirect(http.StatusTemporaryRedirect, "/article/ok")
		})

		article.Any( "/ok", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "投稿成功")
		})
	}

	if err := server.Run(":9000"); err != nil {
		panic(err)
	}
}

type Article struct {
	ID      uint64 `gorm:"primarykey;autoincrement"`
	Title   string `form:"title"`
	Content string `form:"content"`
	Author  string `gorm:"not null"`
	// 创建时间，毫秒作为单位（避免时区）
	Ctime   int64
	// 更新时间，毫秒作为单位
	Utime   int64
}

func (Article) TableName() string {
	return "article"
}

const (
	MySQLDefaultDSN = "root:dangerous@tcp(localhost:3306)/xiaohongshu?charset=utf8&parseTime=True&loc=Local"
	RedisAddr = "localhost:6379"
)

func InitDB() *gorm.DB {
	
	db, err := gorm.Open(mysql.Open(MySQLDefaultDSN),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,  // 不允许外键
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		panic(err.Error())
	}

	if err := db.AutoMigrate(new(Article)); err != nil {
		panic(err.Error())
	} 		
	
	return db
}

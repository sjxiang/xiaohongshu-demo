package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sjxiang/xiaohongshu-demo/internal/conf"
	"github.com/sjxiang/xiaohongshu-demo/internal/repository"
	"github.com/sjxiang/xiaohongshu-demo/internal/repository/dao"
	"github.com/sjxiang/xiaohongshu-demo/internal/service"
	"github.com/sjxiang/xiaohongshu-demo/internal/web"
)

func main() {
	// 初始化 DB
	db := InitDB()

	// 开始组装各种服务，本质上是一个依赖注入
	articleRepo := repository.NewArticleRepo(dao.NewArticleDAO(db))
	articleService := service.NewService(articleRepo)
	articleCtrl := web.NewArticleController(articleService)
	

	server := gin.Default()
	server.LoadHTMLGlob(conf.PATH + "/*")
	articleCtrl.RegisterRoutes(server)
	if err := server.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}

// NewHttpEngine 创建了一个绑定了路由的Web引擎
func NewHttpEngine() (*gin.Engine, error) {
	// 设置为Release，为的是默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	// 默认启动一个Web引擎
	r := gin.New()

	// 业务绑定路由操作
	Routes(r)
	
	// 返回绑定路由后的Web引擎
	return r, nil
}

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	// 如果配置了swagger，则显示swagger的中间件
	// if configService.GetBool("app.swagger") == true {
	// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	// 用户模块
	// user.RegisterRoutes(r)
	// 问答模块
	// qa.RegisterRoutes(r)
}


func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(conf.MySQLDefaultDSN),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,  // 不允许外键
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		panic(err.Error())
	}

	if err := db.AutoMigrate(new(dao.Article)); err != nil {
		panic(err.Error())
	} 		

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// check db connection
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
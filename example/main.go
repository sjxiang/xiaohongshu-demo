package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 准备渲染模板
	// 在 IDE 直接运行的时候，要将 workspace 调整为当前目录，相对路径
	r.LoadHTMLFiles("./login.html")

	v1 := r.Group("/v1")
	{
		// 静态路由
		v1.GET("/health", healthCheck)
		
		// 静态路由（查询参数，例如：敲下 enter 键，/orders?id=123，命中该静态路由）
		v1.GET("/orders", func(ctx *gin.Context) {
			id := ctx.Query("id")
			ctx.String(http.StatusOK, "你传递过来的订单编号是：%s", id)
		})

		// 参数路由
		v1.GET("/users/:name", func(ctx *gin.Context) {
			name := ctx.Param("name")
			ctx.String(http.StatusOK, "你传递过来的用户名是：%s", name)
		})

		// 通配符路由
		v1.GET("/views/*.html", func(ctx *gin.Context) {
			path := ctx.Param(".html")
			ctx.String(http.StatusOK, "匹配上的值是：%s", path)
		})
	}

	user := r.Group("/user")
	{
		// 获得表单页面
		user.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", nil)
		})

		// 提交表单
		user.POST("/login", func(ctx *gin.Context) {
			var userLoginRequest struct{
				Name     string `form:"username" binding:"required,min=3,max=100"`  // 字段，首字母大写
				Password string `form:"password" binding:"required"`
			}

			// 要传入指针
			if err := ctx.Bind(&userLoginRequest); err != nil {
				log.Println(err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "输入错误",
				})
				// 写回响应之后，别忘了 return
				return
			}

			log.Printf("登录信息记录：%v", userLoginRequest)
			ctx.Redirect(http.StatusTemporaryRedirect, "/user/ok")
		})

		user.Any("/ok", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "您已经登录成功")
		})
	}


	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}


func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

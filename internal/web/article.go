package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sjxiang/xiaohongshu-demo/internal/domain"
	"github.com/sjxiang/xiaohongshu-demo/internal/service"
)

type ArticleController struct {
	service service.ArticleService
}

func NewArticleController(service service.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

func (impl *ArticleController) RegisterRoutes(router *gin.Engine) {
	
	article := router.Group("/article")

	{
		article.Any("/:id", impl.GetByID)
		
		article.GET("/new", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "write_article.html", nil)
		})

		article.POST("/new", impl.Save)

		article.Any("/new/success", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "发表成功")
		})
		article.Any("/new/failed", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "发表失败")
		})
	}
}



func (impl *ArticleController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 随便给一个错误码
		ctx.JSON(http.StatusBadRequest, Resp{Code: 1, Msg: "ID 错误"})
		return
	}
	art, err := impl.service.Get(ctx.Request.Context(), id)
	if err != nil {
		// 如果代码严谨的话，这边要区别是真的没有数据，还是服务器出现了异常
		ctx.JSON(http.StatusInternalServerError, Resp{Code: 2, Msg: "找不到对应的帖子"})
		return
	}
	var vo ArticleVO
	vo.init(art)
	ctx.JSON(http.StatusOK, Resp{Data: vo})
}

func (impl *ArticleController) Save(ctx *gin.Context) {
	var vo ArticleVO
	if err := ctx.Bind(&vo); err != nil {
		// 出现 error 的情况下，实际上已经返回le
		return
	}

	// 缺乏登录部分，所以直接写死
	var authorID uint64 = 123
	article, err := impl.service.Save(ctx.Request.Context(), domain.Article{
		Title: vo.Title,
		Content: vo.Content,
		Author: domain.Author{
			ID: authorID,
		},
	})

	if err != nil {
		// 这边不能把 error 写回去，暂时直接输出到控制台
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统异常，请重试",
		})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/article/%d", article.ID))
}



// 跟前端打交道
type ArticleVO struct {
	ID      uint64 `form:"id"`
	Title   string `form:"title"`
	Content string `form:"content"`
	// 一般来说，考虑各种 APP 发版本不容易，所以数字、货币、日期、国际化之类的都是后端做的，前端就是无脑展示 
	Ctime   string
	Utime   string
}


type Resp struct {
	// 业务错误码，不为 0 则是表示出错了
	Code int
	Msg  string
	Data any
}


func (a *ArticleVO) init(article domain.Article) {
	a.ID      = article.ID
	a.Utime   = article.Utime.String()
	a.Ctime   = article.Ctime.String()
	a.Content = article.Content
	a.Title   = article.Title
}
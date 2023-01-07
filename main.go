package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"log"
)

// 自定义中间件 拦截器
func handle() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("usersession", "userid-1")
		//context.Next()放行
		context.Next()
		//context.Abort()拦截
		context.Abort()
	}
}

func main() {
	//创建一个服务
	ginServer := gin.Default()
	//加载静态资源
	ginServer.LoadHTMLGlob("templates/*")
	//加载静态资源文件
	ginServer.Static("/static", "./static")
	//设置接口图标
	ginServer.Use(favicon.New("./fav/favicon.ico"))
	//访问地址,处理请求 Request Response
	ginServer.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{"msg": "GET请求成功"}) //请求成功,返回json数据
	})
	ginServer.POST("/user", func(context *gin.Context) {
		context.JSON(200, gin.H{"msg": "POST请求成功"})
	})
	//响应一个页面给前端
	ginServer.GET("/index", func(context *gin.Context) {
		context.HTML(200, "index.html", gin.H{
			"msg": "前端请求成功",
		})
	})
	//获取请求的参数
	//传统的参数 user?userid=1&name=cxk
	ginServer.GET("/user/info", func(context *gin.Context) {
		userId := context.Query("userid")
		userName := context.Query("name")
		context.JSON(200, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	//restful风格的参数 user/1/cxk
	ginServer.GET("/user/info/:userId/:userName", func(context *gin.Context) {

		userId := context.Param("userId")
		userName := context.Param("userName")
		context.JSON(200, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})
	//
	ginServer.POST("/json", func(context *gin.Context) {
		data, _ := context.GetRawData()
		var m map[string]interface {
		}
		_ = json.Unmarshal(data, &m)
		for s := range m {
			println(s)
		}
		context.JSON(200, m)
	})

	//处理表单数据
	ginServer.POST("/user/add", func(context *gin.Context) {
		username := context.PostForm("username")
		password := context.PostForm("password")
		context.JSON(200, gin.H{
			"username": username,
			"password": password,
		})
	})

	//路由,重定向
	ginServer.GET("/test", func(context *gin.Context) {
		context.Redirect(301, "https://jojo.lw123.top")
	})
	//404页面
	ginServer.NoRoute(func(context *gin.Context) {
		context.HTML(404, "./templates/404.html", nil)
	})
	//路由组,统一管理
	userGruop := ginServer.Group("/user")
	{
		userGruop.GET("/add")
		userGruop.GET("/update")
		userGruop.GET("/delete")
	}
	//拦截器
	ginServer.GET("/test/lanjie", handle(), func(context *gin.Context) {
		log.Println("拦截成功")
		mustGet := context.MustGet("usersession")
		println(mustGet)
		log.Println(mustGet)
	})

	//启动服务:设置端口为8082
	ginServer.Run(":8082")

}

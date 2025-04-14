package router

import (
	"houserent/db"
	"houserent/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func StartServer(port string) {
	var err error
	err = db.InitDb()
	if err != nil {
		return
	}

	r := gin.Default()

	// ✅ 初始化 session 存储（必须要这一步）
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("my_session", store))

	api := r.Group("/api")

	api.POST("/register", RegisterUser)
	api.POST("/login", LoginUser)

	api.Use(middleware.AuthRequired())
	// 需要登录后的api接口设置
	api.POST("/logout", LogoutUser)
	api.POST("/update", UpdateUser)
	api.POST("/delete", DeleteUser)
	api.POST("/user", GetUser)

	// 房源相关接口
	api.POST("/listings/create", CreateListing) // 创建房源
	api.POST("/listings/update", UpdateListing) // 更新房源
	api.POST("/listings/delete", DeleteListing) // 删除房源
	api.POST("/listings/get", GetListing)       // 获取单个房源
	api.POST("/listings/list", GetListings)     // 获取房源列表

	r.Run(":" + port)
}

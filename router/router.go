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

	// 交易相关接口
	api.POST("/transaction/create", CreateTransaction)            // 创建交易
	api.POST("/transaction/update", UpdateTransaction)            // 更新交易
	api.POST("/transaction/get", GetTransaction)                  // 获取单个交易
	api.POST("/transactions/landlord", GetTransactionsByLandlord) // 获取房东的交易
	api.POST("/transactions/tenant", GetTransactionsByTenant)     // 获取租客的交易
	api.POST("/transactions/listing", GetTransactionsByListing)   // 获取房源的交易
	api.POST("/transactions/pending", GetPendingTransactions)     // 获取待处理交易

	// 评价相关路由
	reviewGroup := api.Group("/reviews")
	{
		reviewGroup.POST("/create", CreateReview)       // 创建评价
		reviewGroup.POST("/listing", GetListingReviews) // 获取房源评价
		reviewGroup.POST("/update", UpdateReview)       // 更新评价
		reviewGroup.POST("/delete", DeleteReview)       // 删除评价
	}

	r.Run(":" + port)
}

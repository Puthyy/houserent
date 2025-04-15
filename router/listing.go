package router

import (
	"houserent/db"
	"houserent/model"

	"github.com/gin-gonic/gin"
)

// CreateListing 创建房源
func CreateListing(ctx *gin.Context) {
	var listing model.Listing
	if err := ctx.ShouldBindJSON(&listing); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 设置初始状态
	listing.Status = "available"

	if err := db.AddListing(&listing); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "创建成功", "listing": listing})
}

// UpdateListing 更新房源信息
func UpdateListing(ctx *gin.Context) {
	var listing model.Listing
	if err := ctx.ShouldBindJSON(&listing); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 检查房源是否存在
	existingListing, err := db.FindListingByID(listing.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "房源不存在"})
		return
	}

	// 保留原有的一些字段
	listing.CreatedAt = existingListing.CreatedAt
	listing.UpdatedAt = existingListing.UpdatedAt
	listing.DeletedAt = existingListing.DeletedAt

	if err := db.UpdateListing(&listing); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "更新成功", "listing": listing})
}

// DeleteListing 删除房源
func DeleteListing(ctx *gin.Context) {
	// 从请求体中获取ID
	var request struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 检查房源是否存在
	existingListing, err := db.FindListingByID(request.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "房源不存在"})
		return
	}

	if err := db.DeleteListing(existingListing); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "删除成功"})
}

// GetListing 获取单个房源信息
func GetListing(ctx *gin.Context) {
	// 从请求体中获取ID
	var request struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	existingListing, err := db.FindListingByID(request.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "房源不存在"})
		return
	}

	ctx.JSON(200, gin.H{"listing": existingListing})
}

// GetListings 获取所有房源
func GetListings(c *gin.Context) {
	listings, err := db.FindAllListings()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, listings)
}

// GetLandlordListings 获取房东的所有房源
func GetLandlordListings(c *gin.Context) {
	var req struct {
		LandlordID uint `json:"landlord_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证房东ID是否有效
	if req.LandlordID == 0 {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "房东ID不能为空",
		})
		return
	}

	listings, err := db.FindListingsByLandlord(req.LandlordID)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取房源列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"total":    len(listings),
			"listings": listings,
		},
	})
}

func SetupListingRoutes(r *gin.Engine) {
	listing := r.Group("/api/listings")
	{
		listing.POST("/create", CreateListing)
		listing.POST("/list", GetListings)
		listing.POST("/update", UpdateListing)
		listing.POST("/delete", DeleteListing)
		listing.POST("/landlord", GetLandlordListings)
	}
}

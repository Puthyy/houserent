package router

import (
	"houserent/db"
	"houserent/model"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTransaction 创建交易记录
func CreateTransaction(ctx *gin.Context) {
	var transaction model.Transaction
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 设置初始状态
	transaction.Status = "pending"

	// 检查房源是否存在
	listing, err := db.FindListingByID(transaction.ListingID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "房源不存在"})
		return
	}

	// 检查房源是否可租
	if listing.Status != "available" {
		ctx.JSON(400, gin.H{"error": "房源已被租用"})
		return
	}

	// 检查房东是否存在
	if _, err := db.FindUserByID(transaction.LandlordID); err != nil {
		ctx.JSON(400, gin.H{"error": "房东不存在"})
		return
	}

	// 检查租客是否存在
	if _, err := db.FindUserByID(transaction.TenantID); err != nil {
		ctx.JSON(400, gin.H{"error": "租客不存在"})
		return
	}

	// 验证日期格式
	startDate, err := time.Parse("2006-01-02", transaction.StartDate)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "开始日期格式错误，应为YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", transaction.EndDate)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "结束日期格式错误，应为YYYY-MM-DD"})
		return
	}

	// 验证日期逻辑
	if endDate.Before(startDate) {
		ctx.JSON(400, gin.H{"error": "结束日期不能早于开始日期"})
		return
	}

	// 创建交易记录
	if err := db.AddTransaction(&transaction); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 更新房源状态为已租用
	listing.Status = "rented"
	listing.TenantID = transaction.TenantID
	if err := db.UpdateListing(listing); err != nil {
		ctx.JSON(400, gin.H{"error": "更新房源状态失败"})
		return
	}

	ctx.JSON(200, gin.H{
		"message":     "创建成功",
		"transaction": transaction,
	})
}

// UpdateTransaction 更新交易记录
func UpdateTransaction(ctx *gin.Context) {
	var transaction model.Transaction
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 检查交易记录是否存在
	existingTransaction, err := db.FindTransactionByID(transaction.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "交易记录不存在"})
		return
	}

	// 保留原有的一些字段
	transaction.ListingID = existingTransaction.ListingID
	transaction.LandlordID = existingTransaction.LandlordID
	transaction.TenantID = existingTransaction.TenantID
	transaction.CreatedAt = existingTransaction.CreatedAt

	// 如果更新了日期，验证日期格式和逻辑
	if transaction.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", transaction.StartDate)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "开始日期格式错误，应为YYYY-MM-DD"})
			return
		}
		endDate, err := time.Parse("2006-01-02", transaction.EndDate)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "结束日期格式错误，应为YYYY-MM-DD"})
			return
		}
		if endDate.Before(startDate) {
			ctx.JSON(400, gin.H{"error": "结束日期不能早于开始日期"})
			return
		}
	}

	// 更新交易记录
	if err := db.UpdateTransaction(&transaction); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 如果交易取消，更新房源状态为可租
	if transaction.Status == "cancelled" {
		listing, err := db.FindListingByID(transaction.ListingID)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "房源不存在"})
			return
		}
		listing.Status = "available"
		listing.TenantID = 0
		if err := db.UpdateListing(listing); err != nil {
			ctx.JSON(400, gin.H{"error": "更新房源状态失败"})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"message":     "更新成功",
		"transaction": transaction,
	})
}

// GetTransaction 获取单个交易记录
func GetTransaction(ctx *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transaction, err := db.FindTransactionByID(req.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "交易记录不存在"})
		return
	}

	ctx.JSON(200, gin.H{
		"transaction": transaction,
	})
}

// GetTransactionsByLandlord 获取房东的所有交易记录
func GetTransactionsByLandlord(ctx *gin.Context) {
	var req struct {
		LandlordID uint `json:"landlord_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transactions, err := db.FindTransactionsByLandlord(req.LandlordID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"transactions": transactions,
	})
}

// GetTransactionsByTenant 获取租客的所有交易记录
func GetTransactionsByTenant(ctx *gin.Context) {
	var req struct {
		TenantID uint `json:"tenant_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transactions, err := db.FindTransactionsByTenant(req.TenantID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"transactions": transactions,
	})
}

// GetTransactionsByListing 获取房源的所有交易记录
func GetTransactionsByListing(ctx *gin.Context) {
	var req struct {
		ListingID uint `json:"listing_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transactions, err := db.FindTransactionsByListing(req.ListingID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"transactions": transactions,
	})
}

// GetPendingTransactions 获取所有待处理的交易记录
func GetPendingTransactions(ctx *gin.Context) {
	transactions, err := db.FindPendingTransactions()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"transactions": transactions,
	})
}

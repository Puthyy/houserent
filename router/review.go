package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"houserent/db"
	"houserent/model"
)

// CreateReview 创建评价
func CreateReview(c *gin.Context) {
	var review model.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.AddReview(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评价创建成功", "review": review})
}

// GetListingReviews 获取房源评价
func GetListingReviews(c *gin.Context) {
	var req struct {
		ListingID uint `json:"listing_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviews, err := db.GetListingReviews(req.ListingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// UpdateReview 更新评价
func UpdateReview(c *gin.Context) {
	var review model.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.UpdateReview(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评价更新成功", "review": review})
}

// DeleteReview 删除评价
func DeleteReview(c *gin.Context) {
	var req struct {
		ReviewID uint `json:"review_id"`
		TenantID uint `json:"tenant_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DeleteReview(req.ReviewID, req.TenantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评价删除成功"})
}

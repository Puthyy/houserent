package model

import (
	"gorm.io/gorm"
)

// Review 租户评价结构
type Review struct {
	gorm.Model
	ListingID uint   `json:"listing_id"`
	TenantID  uint   `json:"tenant_id"`
	Rating    int    `json:"rating"` // 1-5星评价
	Comment   string `json:"comment"`
}

// Listing 房源信息结构
type Listing struct {
	gorm.Model
	Housename   string   `json:"housename" binding:"required"`   // 房屋名
	Description string   `json:"description" binding:"required"` // 房源描述
	Price       float64  `json:"price" binding:"required"`       // 租金
	Location    string   `json:"location" binding:"required"`    // 地点
	Images      string   `json:"images"`                         // 图片URL，用逗号分隔
	LandlordID  uint     `json:"landlord_id" binding:"required"` // 外键，关联 User 表
	TenantID    uint     `json:"tenant_id"`                      // 当前租客ID
	Status      string   `json:"status" binding:"required"`      // "available", "rented", "removed"
	ChainTx     string   `json:"chain_tx"`                       // 链上交易ID
	Reviews     []Review `json:"reviews" gorm:"foreignKey:ListingID"`
}

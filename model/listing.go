package model

import "gorm.io/gorm"

type Listing struct {
	gorm.Model
	Housename   string  `json:"housename" binding:"required"`   // 房屋名
	Description string  `json:"description" binding:"required"` // 房源描述
	Price       float64 `json:"price" binding:"required"`       // 租金
	Location    string  `json:"location" binding:"required"`    // 地点
	Images      string  `json:"images"`                         // 图片URL，用逗号分隔
	LandlordID  uint    `json:"landlord_id" binding:"required"` // 外键，关联 User 表
	TenantID    uint    `json:"tenant_id"`                      // 当前租客ID
	Status      string  `json:"status" binding:"required"`      // "available", "rented", "removed"
}

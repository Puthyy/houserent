package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required" gorm:"uniqueIndex"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`     // "landlord" 或 "tenant"
	Email    string `json:"email"`    // 联系方式
	ChainTx  string `json:"chain_tx"` // 链上交易hash，初始为空
}

package model

import "gorm.io/gorm"

// Transaction 交易记录
type Transaction struct {
	gorm.Model
	ListingID  uint    `json:"listing_id" binding:"required"`  // 房源ID
	LandlordID uint    `json:"landlord_id" binding:"required"` // 房东ID
	TenantID   uint    `json:"tenant_id" binding:"required"`   // 租客ID
	Amount     float64 `json:"amount"`                         // 交易金额
	Status     string  `json:"status"`                         // 交易状态：pending, completed, cancelled
	StartDate  string  `json:"start_date"`                     // 租期开始日期
	EndDate    string  `json:"end_date"`                       // 租期结束日期
	ChainTx    string  `json:"chain_tx"`                       // 区块链交易哈希
}

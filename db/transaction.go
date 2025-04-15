package db

import (
	"houserent/model"
)

// AddTransaction 添加交易记录
func AddTransaction(transaction *model.Transaction) error {
	return Db.Create(transaction).Error
}

// UpdateTransaction 更新交易记录
func UpdateTransaction(transaction *model.Transaction) error {
	return Db.Save(transaction).Error
}

// DeleteTransaction 删除交易记录
func DeleteTransaction(transaction *model.Transaction) error {
	return Db.Delete(transaction).Error
}

// FindTransactionByID 根据ID查找交易记录
func FindTransactionByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	err := Db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// FindTransactionsByLandlord 查找房东的所有交易记录
func FindTransactionsByLandlord(landlordID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := Db.Where("landlord_id = ?", landlordID).Find(&transactions).Error
	return transactions, err
}

// FindTransactionsByTenant 查找租客的所有交易记录
func FindTransactionsByTenant(tenantID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := Db.Where("tenant_id = ?", tenantID).Find(&transactions).Error
	return transactions, err
}

// FindTransactionsByListing 查找房源的所有交易记录
func FindTransactionsByListing(listingID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := Db.Where("listing_id = ?", listingID).Find(&transactions).Error
	return transactions, err
}

// FindPendingTransactions 查找所有待处理的交易记录
func FindPendingTransactions() ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := Db.Where("status = ?", "pending").Find(&transactions).Error
	return transactions, err
}

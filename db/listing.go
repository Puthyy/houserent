package db

import (
	"houserent/model"
)

// AddListing 添加房源
func AddListing(listing *model.Listing) error {
	return Db.Create(listing).Error
}

// UpdateListing 更新房源信息
func UpdateListing(listing *model.Listing) error {
	return Db.Save(listing).Error
}

// DeleteListing 删除房源
func DeleteListing(listing *model.Listing) error {
	return Db.Delete(listing).Error
}

// FindListingByID 通过ID查找房源
func FindListingByID(id uint) (*model.Listing, error) {
	var listing model.Listing
	err := Db.Preload("Reviews").First(&listing, id).Error
	return &listing, err
}

// FindListingsByLandlordID 查找房东的所有房源
func FindListingsByLandlordID(landlordID uint) ([]model.Listing, error) {
	var listings []model.Listing
	err := Db.Preload("Reviews").Where("landlord_id = ?", landlordID).Find(&listings).Error
	return listings, err
}

// FindAvailableListings 查找所有可用房源
func FindAvailableListings() ([]model.Listing, error) {
	var listings []model.Listing
	err := Db.Preload("Reviews").Where("status = ?", "available").Find(&listings).Error
	return listings, err
}

// SearchListings 搜索房源
func SearchListings(query map[string]interface{}) ([]model.Listing, error) {
	var listings []model.Listing
	err := Db.Preload("Reviews").Where(query).Find(&listings).Error
	return listings, err
}

// FindAllListings 获取所有房源
func FindAllListings() ([]model.Listing, error) {
	var listings []model.Listing
	if err := Db.Preload("Reviews").Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

// FindListingsByLandlord 获取房东的所有房源
func FindListingsByLandlord(landlordID uint) ([]model.Listing, error) {
	var listings []model.Listing
	if err := Db.Preload("Reviews").
		Where("landlord_id = ?", landlordID).
		Order("created_at DESC").
		Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

// FindListingsByTenant 获取租客的所有房源
func FindListingsByTenant(tenantID uint) ([]model.Listing, error) {
	var listings []model.Listing
	if err := Db.Preload("Reviews").
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

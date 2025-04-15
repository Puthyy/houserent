package db

import (
	"errors"
	"houserent/model"
)

// AddReview 添加评价
func AddReview(review *model.Review) error {
	// 检查房源是否存在
	var listing model.Listing
	if err := Db.First(&listing, review.ListingID).Error; err != nil {
		return errors.New("房源不存在")
	}

	// 检查用户是否租过该房源
	var transaction model.Transaction
	err := Db.Where("listing_id = ? AND tenant_id = ? AND status = ?",
		review.ListingID, review.TenantID, "completed").First(&transaction).Error
	if err != nil {
		return errors.New("只有租过该房源的租客才能评价")
	}

	// 检查是否已经评价过
	var existingReview model.Review
	err = Db.Where("listing_id = ? AND tenant_id = ?",
		review.ListingID, review.TenantID).First(&existingReview).Error
	if err == nil {
		return errors.New("您已经评价过该房源")
	}

	// 创建评价
	return Db.Create(review).Error
}

// GetListingReviews 获取房源的所有评价
func GetListingReviews(listingID uint) ([]model.Review, error) {
	var reviews []model.Review
	err := Db.Where("listing_id = ?", listingID).Find(&reviews).Error
	return reviews, err
}

// UpdateReview 更新评价
func UpdateReview(review *model.Review) error {
	// 检查评价是否存在
	var existingReview model.Review
	if err := Db.First(&existingReview, review.ID).Error; err != nil {
		return errors.New("评价不存在")
	}

	// 检查是否是评价作者
	if existingReview.TenantID != review.TenantID {
		return errors.New("只能修改自己的评价")
	}

	// 更新评价
	return Db.Model(&existingReview).Updates(map[string]interface{}{
		"rating":  review.Rating,
		"comment": review.Comment,
	}).Error
}

// DeleteReview 删除评价
func DeleteReview(reviewID, tenantID uint) error {
	// 检查评价是否存在
	var review model.Review
	if err := Db.First(&review, reviewID).Error; err != nil {
		return errors.New("评价不存在")
	}

	// 检查是否是评价作者
	if review.TenantID != tenantID {
		return errors.New("只能删除自己的评价")
	}

	// 删除评价
	return Db.Delete(&review).Error
}

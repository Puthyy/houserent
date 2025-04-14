package router

import (
	"houserent/db"
	"houserent/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	existingUser, err := db.FindUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		ctx.JSON(400, gin.H{
			"error": "用户名已存在",
		})
		return
	}

	if err := db.AddUser(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "注册成功",
		"user":    user,
	})
}

func LoginUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	u, err := db.FindUserByUsername(user.Username)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "用户不存在",
		})
		return
	}

	if u.Password != user.Password {
		ctx.JSON(400, gin.H{
			"error": "密码错误",
		})
		return
	}

	// 设置 Session
	session := sessions.Default(ctx)
	session.Set("user_id", u.ID)
	session.Set("username", u.Username)
	err = session.Save()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Session保存失败"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "登录成功",
		"user":    u,
	})
}

func LogoutUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	err := session.Save()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "登出失败，Session保存失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "登出成功"})
}

func UpdateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查用户是否存在
	existingUser, err := db.FindUserByID(user.ID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 保留原有的一些字段
	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = existingUser.UpdatedAt
	user.DeletedAt = existingUser.DeletedAt

	if err := db.UpdateUser(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "更新成功",
		"user":    user,
	})
}

func GetUser(ctx *gin.Context) {
	// 从请求体中获取ID
	var request struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 通过ID查找用户
	u, err := db.FindUserByID(request.ID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 返回完整的用户信息
	ctx.JSON(200, gin.H{
		"id":         u.ID,
		"username":   u.Username,
		"password":   u.Password,
		"role":       u.Role,
		"email":      u.Email,
		"chain_tx":   u.ChainTx,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
		"deleted_at": u.DeletedAt,
	})
}

func DeleteUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查用户是否存在
	existingUser, err := db.FindUserByID(user.ID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "用户不存在",
		})
		return
	}

	if err := db.DeleteUser(existingUser); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "删除成功",
	})
}

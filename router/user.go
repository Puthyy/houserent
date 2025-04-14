package router

import (
	"houserent/db"
	"houserent/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {

	var user model.User
	//接口传值
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

	//添加用户到数据库
	if err := db.AddUser(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "注册成功",
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

	// 先查找用户是否存在
	existingUser, err := db.FindUserByUsername(user.Username)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 保留原有的ID
	user.ID = existingUser.ID
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
	})
}

func GetUser(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		ctx.JSON(400, gin.H{
			"error": "用户名不能为空",
		})
		return
	}

	u, err := db.FindUserByUsername(username)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "用户不存在",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"username": u.Username,
		"password": u.Password,
		"role":     u.Role,
		"email":    u.Email,
		"chain_tx": u.ChainTx,
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

	// 先查找用户是否存在
	u, err := db.FindUserByUsername(user.Username)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 删除用户
	if err := db.DeleteUser(u); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "删除成功",
	})
}

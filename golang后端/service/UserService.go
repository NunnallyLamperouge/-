package service

import (
	"example/config"
	"example/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMyInfo(c *gin.Context) {
	var user entity.User
	userId, _ := c.Get("userId")
	if err := config.MysqlDataBase.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "在获取用户信息时出错！详细信息:"+err.Error()))
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, entity.SuccessResponse(user))
}

// GetAllUsers 获取所有用户列表（仅ID、用户名、头像）
func GetAllUsers(c *gin.Context) {
	type SimpleUser struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		Avatar     string `json:"avatar"`
		Apartment  string `json:"apartment"`
		Profession string `json:"profession"`
	}

	var users []SimpleUser

	// 查询所有用户，只获取ID、用户名和头像字段
	if err := config.MysqlDataBase.Model(&entity.User{}).
		Select("id, username, avatar, apartment, profession").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取用户列表失败:"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"users": users,
	}))
}

// GetAllTeachers 获取所有用户列表（仅ID、用户名、头像）
func GetAllTeachers(c *gin.Context) {
	type SimpleUser struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		Avatar     string `json:"avatar"`
		Apartment  string `json:"apartment"`
		Profession string `json:"profession"`
	}

	var users []SimpleUser

	// 查询所有用户，只获取ID、用户名和头像字段
	if err := config.MysqlDataBase.Model(&entity.User{}).
		Select("id, username, avatar, apartment, profession").
		Where("type = ?", 1).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取用户列表失败:"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"users": users,
	}))
}

// GetAllStudents 获取所有用户列表（仅ID、用户名、头像）
func GetAllStudents(c *gin.Context) {
	type SimpleUser struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		Avatar     string `json:"avatar"`
		Apartment  string `json:"apartment"`
		Profession string `json:"profession"`
	}

	var users []SimpleUser
	teacherID, _ := c.Get("userId")

	// 查询状态为1且对应老师ID的指导关系中的学生
	query := config.MysqlDataBase.Model(&entity.User{}).
		Select("users.id, users.username, users.avatar, users.apartment, users.profession").
		Joins("JOIN guides ON users.id = guides.student_user_id").
		Where("guides.status = ?", 1)

	// 如果提供了特定教师ID，则筛选该教师指导的学生
	if teacherID != "" {
		query = query.Where("guides.teacher_user_id = ?", teacherID)
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取学生列表失败:"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"users": users,
	}))
}

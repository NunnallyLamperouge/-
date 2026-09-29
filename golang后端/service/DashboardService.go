package service

import (
	"example/config"
	"example/entity"
	"example/utility"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 学生仪表盘数据结构
type StudentDashboardData struct {
	TemplateEditHeatmap      []TemplateEditHeatmapItem      `json:"templateEditHeatmap"`      // 模板编辑次数热力图数据
	TaskCompletionRate       float64                        `json:"taskCompletionRate"`       // 任务完成率
	NotificationStatistics   []NotificationStatisticsItem   `json:"notificationStatistics"`   // 通知类型统计
	TemplateTypeDistribution []TemplateTypeDistributionItem `json:"templateTypeDistribution"` // 模板类型分布
}

// 模板编辑热力图数据项
type TemplateEditHeatmapItem struct {
	Date  string `json:"date"`  // 日期，格式：YYYY-MM-DD
	Count int    `json:"count"` // 编辑次数
}

// 通知统计数据项
type NotificationStatisticsItem struct {
	Type  string `json:"type"`  // 通知类型
	Count int    `json:"count"` // 数量
}

// 模板类型分布数据项
type TemplateTypeDistributionItem struct {
	Type  string `json:"type"`  // 模板类型
	Count int    `json:"count"` // 数量
}

// GetStudentDashboardData 获取学生仪表盘数据
func GetStudentDashboardData(c *gin.Context) {
	// 获取当前学生用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	studentId := utility.GetUserIDFromContext(userId)

	// 获取时间范围参数（可选）
	timeRange := c.DefaultQuery("timeRange", "30") // 默认30天

	// 构建仪表盘数据
	dashboardData := buildStudentDashboardData(studentId, timeRange)

	c.JSON(http.StatusOK, entity.SuccessResponse(dashboardData))
}

// buildStudentDashboardData 构建学生仪表盘数据
func buildStudentDashboardData(studentId int, timeRange string) StudentDashboardData {
	// 构建空的仪表盘数据结构
	data := StudentDashboardData{
		TemplateEditHeatmap:      []TemplateEditHeatmapItem{},
		TaskCompletionRate:       0,
		NotificationStatistics:   []NotificationStatisticsItem{},
		TemplateTypeDistribution: []TemplateTypeDistributionItem{},
	}

	// 获取时间范围
	days := 30 // 默认30天
	if timeRange == "7" {
		days = 7
	} else if timeRange == "90" {
		days = 90
	} else if timeRange == "365" {
		days = 365
	}

	// 设置时间范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 1. 获取模板编辑次数热力图数据
	data.TemplateEditHeatmap = getTemplateEditHeatmapData(studentId, startDate, endDate)

	// 2. 获取任务完成率
	data.TaskCompletionRate = getTaskCompletionRate(studentId)

	// 3. 获取通知类型统计
	data.NotificationStatistics = getNotificationStatistics(studentId)

	// 4. 获取模板类型分布
	data.TemplateTypeDistribution = getTemplateTypeDistribution(studentId)

	return data
}

// getTemplateEditHeatmapData 获取模板编辑热力图数据
func getTemplateEditHeatmapData(studentId int, startDate, endDate time.Time) []TemplateEditHeatmapItem {
	var result []TemplateEditHeatmapItem

	// 获取学生的所有模板ID
	var templateIds []uint
	config.MysqlDataBase.Model(&entity.FakeTemplate{}).
		Where("user_id = ?", studentId).
		Pluck("id", &templateIds)

	if len(templateIds) == 0 {
		return result
	}

	// 查询模板历史记录
	var histories []struct {
		Date  time.Time
		Count int
	}

	config.MysqlDataBase.Model(&entity.TemplateHistory{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("template_id IN (?) AND created_at BETWEEN ? AND ?", templateIds, startDate, endDate).
		Group("DATE(created_at)").
		Order("date").
		Scan(&histories)

	// 转换为前端需要的格式
	for _, h := range histories {
		result = append(result, TemplateEditHeatmapItem{
			Date:  h.Date.Format("2006-01-02"),
			Count: h.Count,
		})
	}

	// 处理没有编辑记录的日期
	dateMap := make(map[string]bool)
	for _, item := range result {
		dateMap[item.Date] = true
	}

	// 填充没有编辑记录的日期
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if _, exists := dateMap[dateStr]; !exists {
			result = append(result, TemplateEditHeatmapItem{
				Date:  dateStr,
				Count: 0,
			})
		}
	}

	return result
}

// getTaskCompletionRate 获取任务完成率
func getTaskCompletionRate(studentId int) float64 {
	var totalTasks, completedTasks int64

	// 计算总任务数
	config.MysqlDataBase.Model(&entity.Mission{}).
		Where("student_user_id = ?", studentId).
		Count(&totalTasks)

	if totalTasks == 0 {
		return 0
	}

	// 计算已完成任务数（状态为2表示已完成）
	config.MysqlDataBase.Model(&entity.Mission{}).
		Where("student_user_id = ? AND status = 2", studentId).
		Count(&completedTasks)

	// 计算完成率
	return float64(completedTasks) / float64(totalTasks) * 100
}

// getNotificationStatistics 获取通知类型统计
func getNotificationStatistics(studentId int) []NotificationStatisticsItem {
	var result []NotificationStatisticsItem

	// 通知类型映射
	notificationTypes := map[int]string{
		int(entity.SystemNotification):           "系统通知",
		int(entity.ThesisProgressNotification):   "论文进度",
		int(entity.ReviewFeedbackNotification):   "评审反馈",
		int(entity.MilestoneNotification):        "里程碑",
		int(entity.CoauthoringInvitation):        "合著邀请",
		int(entity.FormatTemplateUpdate):         "模板更新",
		int(entity.DefenseScheduleNotification):  "答辩安排",
		int(entity.ThesisSubmissionNotification): "论文提交",
	}

	// 查询各类型通知数量
	var notifications []struct {
		Type  int
		Count int
	}

	config.MysqlDataBase.Model(&entity.Notification{}).
		Select("type, COUNT(*) as count").
		Where("receiver_id = ?", studentId).
		Group("type").
		Scan(&notifications)

	// 转换为前端需要的格式
	for _, n := range notifications {
		typeName, exists := notificationTypes[n.Type]
		if !exists {
			typeName = "其他"
		}

		result = append(result, NotificationStatisticsItem{
			Type:  typeName,
			Count: n.Count,
		})
	}

	return result
}

// getTemplateTypeDistribution 获取模板类型分布
func getTemplateTypeDistribution(studentId int) []TemplateTypeDistributionItem {
	var result []TemplateTypeDistributionItem

	// 查询各类型模板数量
	var distributions []struct {
		Type  string
		Count int
	}

	config.MysqlDataBase.Model(&entity.FakeTemplate{}).
		Select("type, COUNT(*) as count").
		Where("user_id = ?", studentId).
		Group("type").
		Scan(&distributions)

	// 转换为前端需要的格式
	for _, d := range distributions {
		result = append(result, TemplateTypeDistributionItem{
			Type:  d.Type,
			Count: d.Count,
		})
	}

	return result
}

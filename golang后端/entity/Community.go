package entity

import (
	"time"

	"gorm.io/gorm"
)

type Guide struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey" json:"id"`
	StudentUserId uint   `gorm:"not null" json:"studentUserId"`
	TeacherUserId uint   `gorm:"not null" json:"teacherUserId"`
	Status        int    `gorm:"not null" json:"status"`
	Message       string `gorm:"type:text;comment:申请理由"`
	Student       User   `gorm:"foreignKey:StudentUserId;references:ID"`
	Teacher       User   `gorm:"foreignKey:TeacherUserId;references:ID"`
}

// MissionStatus 任务状态枚举
type MissionStatus int

const (
	MissionStatusDraft      MissionStatus = iota + 1 // 草稿
	MissionStatusInProgress                          // 进行中
	MissionStatusSubmitted                           // 已提交待审核
	MissionStatusRejected                            // 已打回
	MissionStatusCompleted                           // 已完成
)

type Mission struct {
	gorm.Model
	ID            uint          `gorm:"primaryKey" json:"id"`
	Name          string        `gorm:"not null" json:"name"`                   // 任务名称
	Description   string        `gorm:"not null" json:"description"`            // 任务描述
	StudentUserId uint          `gorm:"not null" json:"studentUserId"`          // 学生ID
	TeacherUserId uint          `gorm:"not null" json:"teacherUserId"`          // 教师ID
	Status        MissionStatus `gorm:"not null" json:"status"`                 // 任务状态
	EndTime       time.Time     `gorm:"not null" json:"endTime"`                // 截止时间
	StartTime     time.Time     `gorm:"not null" json:"startTime"`              // 任务开始时间
	Priority      int           `gorm:"type:tinyint;default:0" json:"priority"` // 任务优先级
	TemplateId    uint          `gorm:"not null" json:"templateId"`             // 关联的模板ID
	Template      FakeTemplate  `gorm:"foreignKey:TemplateId" json:"template"`  // 关联的模板
	Content       string        `gorm:"type:longtext" json:"content"`           // 学生提交的内容
	Feedback      string        `gorm:"type:text" json:"feedback"`              // 教师反馈
	Student       User          `gorm:"foreignKey:StudentUserId;references:ID"` // 关联学生
	Teacher       User          `gorm:"foreignKey:TeacherUserId;references:ID"` // 关联教师

	// 提醒相关
	LastReminderTime *time.Time `gorm:"comment:上次提醒时间" json:"lastReminderTime"` // 上次提醒时间
	ReminderSent     bool       `gorm:"default:false" json:"reminderSent"`      // 是否已发送提醒

	// 审核相关
	SubmissionTime  *time.Time `gorm:"comment:提交时间" json:"submissionTime"` // 提交时间
	ReviewTime      *time.Time `gorm:"comment:审核时间" json:"reviewTime"`     // 审核时间
	RejectionReason string     `gorm:"type:text" json:"rejectionReason"`   // 打回原因
	RevisionCount   int        `gorm:"default:0" json:"revisionCount"`     // 修改次数
}

// NotificationType 通知类型枚举
type NotificationType int

const (
	SystemNotification           NotificationType = iota + 1 // 系统通知
	ThesisProgressNotification                               // 论文进度通知
	ReviewFeedbackNotification                               // 导师评审反馈通知
	MilestoneNotification                                    // 里程碑达成通知
	CoauthoringInvitation                                    // 论文合著邀请通知
	FormatTemplateUpdate                                     // 格式模板更新通知
	DefenseScheduleNotification                              // 答辩安排通知
	ThesisSubmissionNotification                             // 论文提交通知
)

// Notification 通知模型
type Notification struct {
	gorm.Model
	Type        NotificationType `gorm:"type:int;not null;comment:通知类型"`  // 通知类型
	Title       string           `gorm:"type:varchar(100);comment:通知标题"`  // 通知标题
	Content     string           `gorm:"type:text;comment:通知内容"`          // 通知内容
	SenderID    uint             `gorm:"comment:发送者ID"`                   // 发送者ID，可为0表示系统通知
	ReceiverID  uint             `gorm:"not null;comment:接收者ID"`          // 接收者ID
	IsRead      bool             `gorm:"default:false;comment:是否已读"`      // 是否已读
	ReadTime    *time.Time       `gorm:"comment:阅读时间"`                    // 阅读时间
	RelatedID   uint             `gorm:"comment:相关内容ID"`                  // 相关内容ID，如作品ID、评论ID等
	RelatedType string           `gorm:"type:varchar(50);comment:相关内容类型"` // 相关内容类型
	ExtraData   string           `gorm:"type:text;comment:额外数据JSON"`      // 额外数据，JSON格式
}

// MessageType 消息类型枚举
type MessageType int

const (
	PrivateMessage MessageType = iota + 1 // 私聊消息
	GroupMessage                          // 群聊消息
)

// 媒体消息类型
const (
	MessageTypeText  = "text"  // 文本消息
	MessageTypeImage = "image" // 图片消息
	MessageTypeAudio = "audio" // 音频消息
	MessageTypeVideo = "video" // 视频消息
	MessageTypeFile  = "file"  // 文件消息
)

// Message 私信模型
type Message struct {
	gorm.Model
	Type       MessageType `gorm:"type:int;not null;comment:消息类型"`  // 消息类型
	Content    string      `gorm:"type:text;not null;comment:消息内容"` // 消息内容
	SenderID   uint        `gorm:"not null;comment:发送者ID"`          // 发送者ID
	ReceiverID uint        `gorm:"comment:接收者ID"`                   // 接收者ID，私聊时使用
	GroupID    uint        `gorm:"comment:群组ID"`                    // 群组ID，群聊时使用
	IsRead     bool        `gorm:"default:false;comment:是否已读"`      // 是否已读
	ReadTime   *time.Time  `gorm:"comment:阅读时间"`                    // 阅读时间
	MediaType  string      `gorm:"type:varchar(20);comment:媒体类型"`   // 媒体类型：text, image, audio, video等
	MediaURL   string      `gorm:"type:varchar(255);comment:媒体URL"` // 媒体URL
}

// ChatGroup 聊天群组模型
type ChatGroup struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;comment:群组名称"` // 群组名称
	Description string `gorm:"type:text;comment:群组描述"`                  // 群组描述
	CreatorID   uint   `gorm:"not null;comment:创建者ID"`                  // 创建者ID
	AvatarURL   string `gorm:"type:varchar(255);comment:群组头像"`          // 群组头像
	MemberCount int    `gorm:"default:1;comment:成员数量"`                  // 成员数量
}

// GroupMember 群组成员模型
type GroupMember struct {
	gorm.Model
	GroupID  uint      `gorm:"not null;comment:群组ID"`               // 群组ID
	UserID   uint      `gorm:"not null;comment:用户ID"`               // 用户ID
	Nickname string    `gorm:"type:varchar(50);comment:群内昵称"`       // 群内昵称
	JoinTime time.Time `gorm:"not null;comment:加入时间"`               // 加入时间
	IsAdmin  bool      `gorm:"default:false;comment:是否为管理员"`        // 是否为管理员
	Status   int       `gorm:"default:1;comment:状态 1:正常 2:禁言 3:退出"` // 成员状态
}

// NotificationSetting 通知设置模型
type NotificationSetting struct {
	gorm.Model
	UserID                       uint `gorm:"not null;uniqueIndex;comment:用户ID"` // 用户ID
	SystemNotification           bool `gorm:"default:true;comment:是否接收系统通知"`     // 是否接收系统通知
	ThesisProgressNotification   bool `gorm:"default:true;comment:是否接收论文进度通知"`   // 是否接收论文进度通知
	ReviewFeedbackNotification   bool `gorm:"default:true;comment:是否接收导师评审反馈通知"` // 是否接收导师评审反馈通知
	MilestoneNotification        bool `gorm:"default:true;comment:是否接收里程碑达成通知"`  // 是否接收里程碑达成通知
	CoauthoringInvitation        bool `gorm:"default:true;comment:是否接收论文合著邀请通知"` // 是否接收论文合著邀请通知
	FormatTemplateUpdate         bool `gorm:"default:true;comment:是否接收格式模板更新通知"` // 是否接收格式模板更新通知
	DefenseScheduleNotification  bool `gorm:"default:true;comment:是否接收答辩安排通知"`   // 是否接收答辩安排通知
	ThesisSubmissionNotification bool `gorm:"default:true;comment:是否接收论文提交通知"`   // 是否接收论文提交通知
	EmailNotification            bool `gorm:"default:false;comment:是否接收邮件通知"`    // 是否接收邮件通知
	PushNotification             bool `gorm:"default:true;comment:是否接收推送通知"`     // 是否接收推送通知
}

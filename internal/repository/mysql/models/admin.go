package models

import "time"

// Admin 管理员表
type Admin struct {
	Id          int32     // 主键
	Username    string    // 用户名
	Password    string    // 密码
	Nickname    string    // 昵称
	Mobile      string    // 手机号
	IsUsed      int8      // 是否启用(1:是 -1:否)
	CreatedUser string    // 创建人
	CreatedAt   time.Time `gorm:"time"` // 创建时间
	UpdatedUser string    // 更新人
	UpdatedAt   time.Time `gorm:"time"` // 更新时间
}

// 定义枚举常量
const (
	ADMIN_ISUSED_YES = 1  // 启用
	ADMIN_ISUSED_NOT = -1 // 禁用
)

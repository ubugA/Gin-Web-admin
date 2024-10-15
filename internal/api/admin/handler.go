package admin

import (
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/repository/mysql"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Register 管理员注册
	Register() core.HandlerFunc

	// Login 管理员登录
	Login() core.HandlerFunc

	// List 获取管理员列表
	List() core.HandlerFunc

	// One 获取单条管理员信息
	One() core.HandlerFunc

	// Update 更新管理员信息
	Update() core.HandlerFunc

	// Delete 删除管理员
	Delete() core.HandlerFunc
}

type handler struct {
	logger *zap.Logger
	db     mysql.Repo
}

func New(logger *zap.Logger, db mysql.Repo) Handler {
	return &handler{
		logger: logger,
		db:     db,
	}
}

func (h *handler) i() {}

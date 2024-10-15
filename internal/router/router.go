package router

import (
	"gin-api-admin/internal/api/admin"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/repository/mysql"
	"gin-api-admin/internal/router/interceptor"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewHTTPMux(logger *zap.Logger, db mysql.Repo) (core.Mux, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	if db == nil {
		return nil, errors.New("db required")
	}

	mux, err := core.New(logger,
		core.WithEnableCors(),
		core.WithEnableSwagger(),
		core.WithEnablePProf(),
	)

	if err != nil {
		panic(err)
	}

	// 初始化拦截器
	interceptorHandler := interceptor.New(logger, db)

	// 管理员后台模块
	adminHandler := admin.New(logger, db)
	adminRouter := mux.Group("/api", core.WrapAuthHandler(interceptorHandler.JWTokenAuthVerify))
	{
		// 管理员列表
		adminRouter.GET("/admins", adminHandler.List())

		// 查询单个管理员
		adminRouter.GET("/admin/:username", adminHandler.One())

		// 修改管理员信息
		adminRouter.PUT("/admin/:username", adminHandler.Update())

		// 删除个人管理员
		adminRouter.DELETE("/admin/:username", adminHandler.Delete())
	}

	// 登录&注册
	loginRouter := mux.Group("/api")
	{
		// 注册
		loginRouter.POST("/admin/register", adminHandler.Register())

		// 登录
		loginRouter.POST("/admin/login", adminHandler.Login())
	}

	return mux, nil
}

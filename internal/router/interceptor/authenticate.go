package interceptor

import (
	"net/http"

	"gin-api-admin/internal/pkg/core"
)

func (i *interceptor) Authenticate() core.HandlerFunc {
	return func(ctx core.Context) {

		// 身份信息
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				10104,
				"缺少 Authorization 信息错误"),
			)
			return
		}

		// 根据自己业务编写验证逻辑
	}
}

package admin

import (
	"fmt"
	"net/http"

	"gin-api-admin/internal/code"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/pkg/timeutil"
	"gin-api-admin/internal/pkg/utils"
	"gin-api-admin/internal/pkg/validation"
	"gin-api-admin/internal/repository/mysql/models"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
	Nickname string `json:"nickname"`                    // 昵称
	Mobile   string `json:"mobile" binding:"required"`   // 手机号
}

type registerResponse struct {
	Id          int32  `json:"id"`           // 主键ID
	Username    string `json:"username"`     // 用户名
	Nickname    string `json:"nickname"`     // 昵称
	Mobile      string `json:"mobile"`       // 手机号
	IsUsed      int8   `json:"is_used"`      // 是否启用(1:是 -1:否)
	CreatedUser string `json:"created_user"` // 创建人
	CreatedAt   string `json:"created_at"`   // 创建时间
}

// Register 管理员注册
// @Summary 管理员注册
// @Description 管理员注册
// @Tags API.register
// @Accept json
// @Produce json
// @Param RequestBody body registerRequest true "请求参数"
// @Success 200 {object} registerResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/register [post]
func (h *handler) Register() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(registerRequest)
		res := new(registerResponse)
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)),
			)
			return
		}

		registerData := new(models.Admin)
		registerData.Username = req.Username
		registerData.Nickname = req.Nickname
		registerData.Mobile = req.Mobile

		// 生成密码
		hashedPassword, err := utils.GenerateHashedPassword(req.Password)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminRegisterError,
				fmt.Sprintf("%s: 生成密码失败", code.Text(code.AdminRegisterError))),
			)
			return
		}

		registerData.Password = hashedPassword

		// 设置为启用
		registerData.IsUsed = models.ADMIN_ISUSED_YES

		// 创建人
		registerData.CreatedUser = "主动注册"

		dbResult := h.db.GetDbW().WithContext(ctx.RequestContext()).Create(registerData)

		if dbResult.Error != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminRegisterError,
				fmt.Sprintf("%s: %s", code.Text(code.AdminRegisterError), dbResult.Error.Error())),
			)
			return
		}

		res.Id = registerData.Id
		res.Username = registerData.Username
		res.Nickname = registerData.Nickname
		res.Mobile = registerData.Mobile
		res.IsUsed = registerData.IsUsed
		res.CreatedUser = registerData.CreatedUser
		res.CreatedAt = registerData.CreatedAt.Format(timeutil.CSTLayout)

		ctx.Payload(res)
	}
}

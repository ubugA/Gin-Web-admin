package admin

import (
	"fmt"
	"net/http"

	"gin-api-admin/internal/code"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/pkg/timeutil"
	"gin-api-admin/internal/pkg/validation"
	"gin-api-admin/internal/repository/mysql/models"

	"gorm.io/gorm"
)

type oneRequest struct {
	Username string `uri:"username" binding:"required"` // 用户名
}

type oneResponse struct {
	Id          int32  `json:"id"`           // 主键ID
	Username    string `json:"username"`     // 用户名
	Nickname    string `json:"nickname"`     // 昵称
	Mobile      string `json:"mobile"`       // 手机号
	IsUsed      int8   `json:"is_used"`      // 是否启用(1:是 -1:否)
	CreatedUser string `json:"created_user"` // 创建人
	CreatedAt   string `json:"created_at"`   // 创建时间
	UpdatedUser string `json:"updated_user"` // 更新人
	UpdatedAt   string `json:"updated_at"`   // 更新时间
}

// One 获取单条管理员信息
// @Summary 获取单条管理员信息
// @Description 获取单条管理员信息
// @Tags API.admin
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} oneResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/{username} [get]
// @Security LoginVerifyToken
func (h *handler) One() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(oneRequest)
		res := new(oneResponse)
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)),
			)
			return
		}

		var resultData *models.Admin
		dbResult := h.db.GetDbR().WithContext(ctx.RequestContext()).
			Where("username = ?", req.Username).
			First(&resultData)

		if dbResult.Error != nil && dbResult.Error != gorm.ErrRecordNotFound {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminOneError,
				fmt.Sprintf("%s: %s", code.Text(code.AdminOneError), dbResult.Error.Error())),
			)
			return
		}

		res.Id = resultData.Id
		res.Username = resultData.Username
		res.Nickname = resultData.Nickname
		res.Mobile = resultData.Mobile
		res.IsUsed = resultData.IsUsed
		res.CreatedUser = resultData.CreatedUser
		res.CreatedAt = resultData.CreatedAt.Format(timeutil.CSTLayout)
		res.UpdatedUser = resultData.UpdatedUser
		res.UpdatedAt = resultData.UpdatedAt.Format(timeutil.CSTLayout)

		ctx.Payload(res)
	}
}

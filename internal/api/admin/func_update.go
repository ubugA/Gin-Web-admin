package admin

import (
	"fmt"
	"net/http"
	"time"

	"gin-api-admin/internal/code"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/pkg/timeutil"
	"gin-api-admin/internal/pkg/utils"
	"gin-api-admin/internal/pkg/validation"
	"gin-api-admin/internal/repository/mysql/models"
)

type updateRequest struct {
	Password string `json:"password"` // 密码
	Nickname string `json:"nickname"` // 昵称
	Mobile   string `json:"mobile"`   // 手机号
	IsUsed   int8   `json:"is_used"`  // 是否启用(1:是 -1:否)
}

type updateResponse struct {
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

// Update 更新管理员信息
// @Summary 更新管理员信息
// @Description 更新管理员信息
// @Tags API.admin
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Param RequestBody body updateRequest true "请求参数"
// @Success 200 {object} updateResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/{username} [put]
// @Security LoginVerifyToken
func (h *handler) Update() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(updateRequest)
		res := new(updateResponse)
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)),
			)
			return
		}

		username := ctx.RequestPathParams("username")
		if username == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminUpdateError,
				fmt.Sprintf("%s：缺少 username 参数", code.Text(code.AdminUpdateError))),
			)
			return
		}

		// 查询用户名是否存在
		var searchAdminData *models.Admin
		searchAdminResult := h.db.GetDbR().WithContext(ctx.RequestContext()).
			Where("username = ?", username).
			First(&searchAdminData)

		if searchAdminResult.Error != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminUpdateError,
				fmt.Sprintf("%s：用户名(%s)不存在", code.Text(code.AdminUpdateError), username)),
			)
			return
		}

		// 执行更新操作
		updateMap := map[string]interface{}{}

		if req.Password != "" {
			// 生成密码
			hashedPassword, err := utils.GenerateHashedPassword(req.Password)
			if err != nil {
				ctx.AbortWithError(core.Error(
					http.StatusBadRequest,
					code.AdminUpdateError,
					fmt.Sprintf("%s: 生成密码失败", code.Text(code.AdminUpdateError))),
				)
				return
			}

			updateMap["password"] = hashedPassword
		}

		if req.Nickname != "" {
			updateMap["nickname"] = req.Nickname
		}

		if req.Mobile != "" {
			updateMap["mobile"] = req.Mobile
		}

		if req.IsUsed != 0 && req.IsUsed != searchAdminData.IsUsed {
			if req.IsUsed != models.ADMIN_ISUSED_YES && req.IsUsed != models.ADMIN_ISUSED_NOT {
				ctx.AbortWithError(core.Error(
					http.StatusBadRequest,
					code.AdminUpdateError,
					fmt.Sprintf("%s：is_used 参数错误，只能为 %d 或 %d", code.Text(code.AdminUpdateError), models.ADMIN_ISUSED_YES, models.ADMIN_ISUSED_NOT)),
				)
				return
			}

			updateMap["is_used"] = req.IsUsed
		}

		if len(updateMap) == 0 {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminUpdateError,
				fmt.Sprintf("%s：未提交需要修改的参数", code.Text(code.AdminUpdateError))),
			)
			return
		}

		updateMap["updated_user"] = ctx.SessionUserInfo().UserName
		updateMap["updated_at"] = time.Now().Format(timeutil.CSTLayout)

		var updateModel *models.Admin
		updateResult := h.db.GetDbW().WithContext(ctx.RequestContext()).
			Model(&models.Admin{}).
			Where("username = ?", username).
			Updates(updateMap).
			Find(&updateModel)

		if updateResult.Error != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminUpdateError,
				fmt.Sprintf("%s: %s", code.Text(code.AdminUpdateError), updateResult.Error.Error()),
			))
			return
		}

		res.Id = updateModel.Id
		res.Username = updateModel.Username
		res.Nickname = updateModel.Nickname
		res.Mobile = updateModel.Mobile
		res.IsUsed = updateModel.IsUsed
		res.CreatedUser = updateModel.CreatedUser
		res.CreatedAt = updateModel.CreatedAt.Format(timeutil.CSTLayout)
		res.UpdatedUser = updateModel.UpdatedUser
		res.UpdatedAt = updateModel.UpdatedAt.Format(timeutil.CSTLayout)

		ctx.Payload(res)
	}
}

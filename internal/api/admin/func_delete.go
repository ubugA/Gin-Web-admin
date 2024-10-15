package admin

import (
	"fmt"
	"net/http"

	"gin-api-admin/internal/code"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/pkg/validation"
	"gin-api-admin/internal/repository/mysql/models"
)

type deleteRequest struct {
	Username string `uri:"username" binding:"required"` // 用户名
}

type deleteResponse struct {
	Username string `json:"username"` // 用户名
}

// Delete 删除管理员信息
// @Summary 删除管理员信息
// @Description 删除管理员信息
// @Tags API.admin
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} deleteResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/{username} [delete]
// @Security LoginVerifyToken
func (h *handler) Delete() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(deleteRequest)
		res := new(deleteResponse)
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)),
			)
			return
		}

		// 查询用户名是否存在
		var searchAdminData *models.Admin
		searchAdminResult := h.db.GetDbR().WithContext(ctx.RequestContext()).
			Where("username = ?", req.Username).
			First(&searchAdminData)

		if searchAdminResult.Error != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminUpdateError,
				fmt.Sprintf("%s：用户名(%s)不存在", code.Text(code.AdminDeleteError), req.Username)),
			)
			return
		}

		dbResult := h.db.GetDbW().WithContext(ctx.RequestContext()).
			Where("username = ?", req.Username).
			Delete(&models.Admin{})

		if dbResult.Error != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDeleteError,
				fmt.Sprintf("%s: %s", code.Text(code.AdminDeleteError), dbResult.Error.Error())),
			)
			return
		}

		res.Username = req.Username
		ctx.Payload(res)
	}
}

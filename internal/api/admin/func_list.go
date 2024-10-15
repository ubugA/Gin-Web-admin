package admin

import (
	"fmt"
	"net/http"

	"gin-api-admin/internal/code"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/pkg/timeutil"
	"gin-api-admin/internal/pkg/validation"
	"gin-api-admin/internal/repository/mysql/models"
)

type listRequest struct {
	Page      int    `form:"page"`       // 当前页码，默认为第一页
	PageSize  int    `form:"page_size"`  // 每页返回的数据量
	SortField string `form:"sort_field"` // 排序字段的名称，默认为 id
	SortOrder string `form:"sort_order"` // 排序的顺序，可以是 asc(升序) 或 desc(降序)。
	Username  string `form:"username"`   // 用户名，支持模糊查询
	Nickname  string `form:"nickname"`   // 昵称，支持模糊查询
	Mobile    string `form:"mobile"`     // 手机号，支持模糊查询
}

type listData struct {
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

type listResponse struct {
	Page     int        `json:"page"`      // 当前页码
	PageSize int        `json:"page_size"` // 每页返回的数据量
	Total    int64      `json:"total"`     // 符合查询条件的总记录数
	Items    []listData `json:"items"`     // 当前页的数据列表
}

// List 获取管理员列表
// @Summary 获取管理员列表
// @Description 获取管理员列表
// @Tags API.admin
// @Accept json
// @Produce json
// @Param page query int true "当前页码" default(1)
// @Param page_size query int true "每页返回的数据量，最多 200 条" default(20)
// @Param sort_field query string true "排序字段的名称" default(id)
// @Param sort_order query string true "排序的顺序" default(desc)
// @Param username query string false "用户名，支持模糊查询"
// @Param nickname query string false "昵称，支持模糊查询"
// @Param mobile query string false "手机号，支持模糊查询"
// @Success 200 {object} listResponse
// @Failure 400 {object} code.Failure
// @Router /api/admins [get]
// @Security LoginVerifyToken
func (h *handler) List() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(listRequest)
		res := new(listResponse)
		if err := ctx.ShouldBindForm(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)),
			)
			return
		}

		db := h.db.GetDbR().WithContext(ctx.RequestContext())

		if req.Username != "" {
			db = db.Where("username like ?", "%"+req.Username+"%")
		}

		if req.Mobile != "" {
			db = db.Where("mobile like ?", "%"+req.Mobile+"%")
		}

		if req.Nickname != "" {
			db = db.Where("nickname like ?", "%"+req.Nickname+"%")
		}

		if req.Page == 0 {
			req.Page = 1
		}

		if req.PageSize == 0 {
			req.PageSize = 20
		}

		if req.PageSize > 200 {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminListError,
				fmt.Sprintf("%s: 一次最多只能查询 200 条", code.Text(code.AdminListError)),
			))
			return
		}

		if req.SortField == "" {
			req.SortField = "id"
		}

		if req.SortOrder == "" {
			req.SortOrder = "desc"
		}

		var resultData []*models.Admin
		dbResult := db.
			Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder)).
			Limit(req.PageSize).
			Offset((req.Page - 1) * req.PageSize).
			Find(&resultData)

		if dbResult.Error != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminListError,
				fmt.Sprintf("%s: %s", code.Text(code.AdminListError), dbResult.Error.Error()),
			))
			return
		}

		res.Page = req.Page
		res.PageSize = req.PageSize

		var total int64
		db.Model(&resultData).Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Count(&total)
		res.Total = total

		res.Items = make([]listData, len(resultData))
		for k, v := range resultData {
			res.Items[k] = listData{
				Id:          v.Id,
				Username:    v.Username,
				Nickname:    v.Nickname,
				Mobile:      v.Mobile,
				IsUsed:      v.IsUsed,
				CreatedUser: v.CreatedUser,
				CreatedAt:   v.CreatedAt.Format(timeutil.CSTLayout),
				UpdatedUser: v.UpdatedUser,
				UpdatedAt:   v.UpdatedAt.Format(timeutil.CSTLayout),
			}
		}

		ctx.Payload(res)
	}
}

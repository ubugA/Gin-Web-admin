package admin

import (
	"fmt"
	"net/http"
	"time"

	"gin-api-admin/configs"
	"gin-api-admin/internal/code"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/pkg/jwtoken"
	"gin-api-admin/internal/pkg/timeutil"
	"gin-api-admin/internal/pkg/utils"
	"gin-api-admin/internal/pkg/validation"
	"gin-api-admin/internal/proposal"
	"gin-api-admin/internal/repository/mysql/models"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type loginResponse struct {
	Token    string `json:"token"` // 登录成功后颁发的 Token
	UserInfo struct {
		Id          int32  `json:"id"`           // 主键ID
		Username    string `json:"username"`     // 用户名
		Nickname    string `json:"nickname"`     // 昵称
		Mobile      string `json:"mobile"`       // 手机号
		IsUsed      int8   `json:"is_used"`      // 是否启用(1:是 -1:否)
		CreatedUser string `json:"created_user"` // 创建人
		CreatedAt   string `json:"created_at"`   // 创建时间
		UpdatedUser string `json:"updated_user"` // 更新人
		UpdatedAt   string `json:"updated_at"`   // 更新时间
	} `json:"user_info"`
}

// Login 管理员登录
// @Summary 管理员登录
// @Description 管理员登录
// @Tags API.login
// @Accept json
// @Produce json
// @Param RequestBody body loginRequest true "请求参数"
// @Success 200 {object} loginResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/login [post]
func (h *handler) Login() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(loginRequest)
		res := new(loginResponse)
		if err := ctx.ShouldBindJSON(req); err != nil {
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
				code.AdminLoginError,
				fmt.Sprintf("%s：用户名(%s)不存在", code.Text(code.AdminLoginError), req.Username)),
			)
			return
		}

		// 验证状态
		if searchAdminData.IsUsed != models.ADMIN_ISUSED_YES {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				fmt.Sprintf("%s：用户名(%s)已被禁用", code.Text(code.AdminLoginError), req.Username)),
			)
			return
		}

		// 验证密码
		if !utils.VerifyHashedPassword(searchAdminData.Password, req.Password) {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				fmt.Sprintf("%s：用户名或密码错误", code.Text(code.AdminLoginError))),
			)
			return
		}

		// 设置 Session 信息
		sessionUserInfo := proposal.SessionUserInfo{
			Id:       searchAdminData.Id,
			UserName: searchAdminData.Username,
			NickName: searchAdminData.Nickname,
		}

		// 设置载荷数据、有效期，生成 JWT Token String
		tokenString, err := jwtoken.New(configs.Get().JWT.Secret).Sign(sessionUserInfo, 24*time.Hour)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				fmt.Sprintf("%s：token 生成失败(%s)", code.Text(code.AdminLoginError), err.Error())),
			)
			return
		}

		res.Token = tokenString

		res.UserInfo.Id = searchAdminData.Id
		res.UserInfo.Username = searchAdminData.Username
		res.UserInfo.Nickname = searchAdminData.Nickname
		res.UserInfo.Mobile = searchAdminData.Mobile
		res.UserInfo.IsUsed = searchAdminData.IsUsed
		res.UserInfo.CreatedUser = searchAdminData.CreatedUser
		res.UserInfo.CreatedAt = searchAdminData.CreatedAt.Format(timeutil.CSTLayout)
		res.UserInfo.UpdatedUser = searchAdminData.UpdatedUser
		res.UserInfo.UpdatedAt = searchAdminData.UpdatedAt.Format(timeutil.CSTLayout)

		ctx.Payload(res)
	}
}

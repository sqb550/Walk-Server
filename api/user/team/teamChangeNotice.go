package team

import (
	"reflect"
	"runtime"

	teamCache "app/dao/cache/team"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"

	"app/comm"
)

func TeamChangeNoticeHandler() gin.HandlerFunc {
	api := TeamChangeNoticeApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfTeamChangeNotice).Pointer()).Name()] = api
	return hfTeamChangeNotice
}

type TeamChangeNoticeApi struct {
	Info     struct{} `name:"获取团队变更通知" desc:"返回当前用户尚未确认的团队密码、路线或被移除通知"`
	Request  struct{}
	Response TeamChangeNoticeApiResponse
}

type TeamChangeNoticeApiResponse struct {
	PasswordChanged bool   `json:"password_changed" desc:"团队密码是否已修改"`
	RouteChanged    bool   `json:"route_changed" desc:"团队路线是否已修改"`
	RemovedFromTeam bool   `json:"removed_from_team" desc:"是否已被移出队伍"`
	RemovedTeamName string `json:"removed_team_name" desc:"被移出的原队伍名称"`
}

func (h *TeamChangeNoticeApi) Run(ctx *gin.Context) kit.Code {
	userID, err := comm.GetUserIDFromCtx(ctx)
	if err != nil || userID <= 0 {
		return comm.CodeNotLoggedIn
	}
	notice, err := teamCache.GetTeamChangeNotice(ctx, userID)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Warn("读取团队变更通知失败")
		return comm.CodeServerError
	}
	h.Response.PasswordChanged = notice.PasswordChanged
	h.Response.RouteChanged = notice.RouteChanged
	h.Response.RemovedFromTeam = notice.RemovedFromTeam
	h.Response.RemovedTeamName = notice.RemovedTeamName
	return comm.CodeOK
}

func hfTeamChangeNotice(ctx *gin.Context) {
	api := &TeamChangeNoticeApi{}
	if code := api.Run(ctx); code == comm.CodeOK {
		reply.Reply(ctx, comm.CodeOK, api.Response)
	} else {
		reply.Fail(ctx, code)
	}
}

func TeamChangeNoticeAckHandler() gin.HandlerFunc {
	api := TeamChangeNoticeAckApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfTeamChangeNoticeAck).Pointer()).Name()] = api
	return hfTeamChangeNoticeAck
}

type TeamChangeNoticeAckApi struct {
	Info     struct{} `name:"确认团队变更通知" desc:"清除当前用户明确确认过的团队变更通知"`
	Request  TeamChangeNoticeAckApiRequest
	Response struct{}
}

type TeamChangeNoticeAckApiRequest struct {
	Body struct {
		PasswordChanged bool `json:"password_changed" desc:"确认团队密码变更通知"`
		RouteChanged    bool `json:"route_changed" desc:"确认团队路线变更通知"`
		RemovedFromTeam bool `json:"removed_from_team" desc:"确认被移出队伍通知"`
	}
}

func (h *TeamChangeNoticeAckApi) Init(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&h.Request.Body)
}

func (h *TeamChangeNoticeAckApi) Run(ctx *gin.Context) kit.Code {
	if !h.Request.Body.PasswordChanged && !h.Request.Body.RouteChanged && !h.Request.Body.RemovedFromTeam {
		return comm.CodeParameterInvalid
	}
	userID, err := comm.GetUserIDFromCtx(ctx)
	if err != nil || userID <= 0 {
		return comm.CodeNotLoggedIn
	}
	if err := teamCache.AckTeamChangeNotice(
		ctx,
		userID,
		h.Request.Body.PasswordChanged,
		h.Request.Body.RouteChanged,
		h.Request.Body.RemovedFromTeam,
	); err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Warn("确认团队变更通知失败")
		return comm.CodeServerError
	}
	return comm.CodeOK
}

func hfTeamChangeNoticeAck(ctx *gin.Context) {
	api := &TeamChangeNoticeAckApi{}
	if err := api.Init(ctx); err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Warn("参数绑定校验错误")
		reply.Fail(ctx, comm.CodeParameterInvalid)
		return
	}
	if code := api.Run(ctx); code == comm.CodeOK {
		reply.Reply(ctx, comm.CodeOK, api.Response)
	} else {
		reply.Fail(ctx, code)
	}
}

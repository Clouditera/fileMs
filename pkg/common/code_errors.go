package common

import (
	"github.com/mlogclub/simple/web"
)

var (
	ErrorNotLogin = web.NewError(1, "请先登录")

	//ErrParamInvalid = web.NewError(1000, "传递参数异常")
	//CaptchaError        = web.NewError(1000, "验证码错误")
	//ForbiddenError      = web.NewError(1001, "已被禁言")
	//UserDisabled        = web.NewError(1002, "账号已禁用")
	//InObservationPeriod = web.NewError(1003, "账号尚在观察期")
	//EmailNotVerified    = web.NewError(1004, "请先验证邮箱")
)

const (
	ServiceNum = 1000000
	ModuleNum  = 10000
)

//  服务编码号
const (
	ServiceZhurong  = 1 * ServiceNum
	ServiceCodeFuzz = 2 * ServiceNum
	ServiceFS       = 3 * ServiceNum
	ServiceOther    = 9 * ServiceNum
)

// 模块编码
const (
	ModuleUser  = 01 * ModuleNum
	ModuleFS    = 02 * ModuleNum
	ModuleMinio = 03 * ModuleNum
)

// errCode
const (
	ErrParamInvalid = 1000 // 传递参数异常
	ErrPramNull     = 1001 // 参数为空
	ErrPramType     = 1002 // 参数类型错误
	ErrDBOperate    = 2000 // 操作数据库异常

	ErrBusiness    = 3000 // 业务异常
	ErrBusUpload   = 3001 // 上传失败
	ErrBusDownload = 3002 // 下载失败
	ErrBusDelete   = 3003 // 删除失败

	ErrLogin  = 4000 // 用户名密码错误
	ErrSign   = 4001 // 注册失败
	ErrLogout = 4002 // 退出失败
	ErrUpdate = 4003 // 用户信息更新失败

	ErrUnauthorized = 8000 // 未授权
	ErrPermission   = 8001 // 权限不足
	ErrNotFound     = 8002 // 不存在
	ErrForbidden    = 8003 // 禁止
	ErrExisted      = 8004 // 已存在

	ErrInternal = 9000 // 系统内部错误

	ErrUnknown = 9999 // 未知错误
)

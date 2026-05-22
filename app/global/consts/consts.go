package consts

const (

	ProcessKilled string = "收到信号，进程被结束"

	ValidatorPrefix              string = "Form_Validator_"
	ValidatorParamsCheckFailCode int    = -400300
	ValidatorParamsCheckFailMsg  string = "参数校验失败"

	ServerOccurredErrorCode int    = -500100
	ServerOccurredErrorMsg  string = "服务器内部发生代码执行错误, "
	GinSetTrustProxyError   string = "Gin 设置信任代理服务器出错"

	JwtTokenOK            int    = 200100
	JwtTokenInvalid       int    = -400100
	JwtTokenExpired       int    = -400101
	JwtTokenFormatErrCode int    = -400102
	JwtTokenFormatErrMsg  string = "提交的 token 格式错误"
	JwtTokenMustValid     string = "token为必填项,请在请求header部分提交!"

	StartTimeStamp = int64(1483228800000)
	MachineIdBits  = uint(10)
	SequenceBits   = uint(12)

	SequenceMask   = int64(-1 ^ (-1 << SequenceBits))
	MachineIdShift = SequenceBits
	TimestampShift = SequenceBits + MachineIdBits

	CurdStatusOkCode         int    = 200
	CurdStatusOkMsg          string = "Success"
	CurdCreatFailCode        int    = -400200
	CurdCreatFailMsg         string = "新增失败"
	CurdUpdateFailCode       int    = -400201
	CurdUpdateFailMsg        string = "更新失败"
	CurdDeleteFailCode       int    = -400202
	CurdDeleteFailMsg        string = "删除失败"
	CurdSelectFailCode       int    = -400203
	CurdSelectFailMsg        string = "查询无数据"
	CurdRegisterFailCode     int    = -400204
	CurdRegisterFailMsg      string = "注册失败"
	CurdLoginFailCode        int    = -400205
	CurdLoginFailMsg         string = "登录失败"
	CurdRefreshTokenFailCode int    = -400206
	CurdRefreshTokenFailMsg  string = "刷新Token失败"

	FilesUploadFailCode            int    = -400250
	FilesUploadFailMsg             string = "文件上传失败, 获取上传文件发生错误!"
	FilesUploadMoreThanMaxSizeCode int    = -400251
	FilesUploadMoreThanMaxSizeMsg  string = "长传文件超过系统设定的最大值,系统允许的最大值："
	FilesUploadMimeTypeFailCode    int    = -400252
	FilesUploadMimeTypeFailMsg     string = "文件mime类型不允许"
	FilesUploadIsEmpty             string = "不允许上传空文件"

	WsServerNotStartCode int    = -400300
	WsServerNotStartMsg  string = "websocket 服务没有开启，请在配置文件开启，相关路径：config/config.yml"
	WsOpenFailCode       int    = -400301
	WsOpenFailMsg        string = "websocket open阶段初始化基本参数失败"

	CaptchaGetParamsInvalidMsg    string = "获取验证码：提交的验证码参数无效,请检查验证码ID以及文件名后缀是否完整"
	CaptchaGetParamsInvalidCode   int    = -400350
	CaptchaCheckParamsInvalidMsg  string = "校验验证码：提交的参数无效，请检查 【验证码ID、验证码值】 提交时的键名是否与配置项一致"
	CaptchaCheckParamsInvalidCode int    = -400351
	CaptchaCheckOkMsg             string = "验证码校验通过"
	CaptchaCheckFailCode          int    = -400355
	CaptchaCheckFailMsg           string = "验证码校验失败"
)

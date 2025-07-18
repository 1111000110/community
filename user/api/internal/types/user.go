// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3

package types

type UserDeleteReq struct {
	Phone    string `json:"phone"`    // 根据手机号删除用户
	UserId   int64  `json:"user_id"`  // 要删除的用户ID
	Password string `json:"password"` // 用户密码，用于身份验证
}

type UserDeleteResp struct {
}

type UserLoginReq struct {
	Phone    string `json:"phone"`    // 手机号登录（推荐方式）
	UserId   int64  `json:"userId"`   // 用户ID登录（备选方式）
	Password string `json:"password"` // 用户密码，明文传输（建议HTTPS）
}

type UserLoginResp struct {
	Token string `json:"token"` // JWT访问令牌，有效期见配置
}

type UserQueryReq struct {
	QueryUserId int64  `json:"query_user_id"` // 要查询的用户ID
	Type        string `json:"type"`          // 查询类型："basic"-基本信息 "private"-私密信息
}

type UserQueryResp struct {
	User UserInfo `json:"user_info"` // 用户信息，内容根据type参数决定
}

type UserRegisterReq struct {
	Phone    string `json:"phone"`    // 手机号，必须唯一，用作登录凭证
	Password string `json:"password"` // 用户密码，建议8位以上包含字母数字
	Role     string `json:"role"`     // 用户角色，默认为"user"
}

type UserRegisterResp struct {
	UserId int64  `json:"user_id"` // 新创建的用户ID
	Token  string `json:"token"`   // JWT访问令牌，注册即登录
}

type UserUpdateReq struct {
	UserInfo   UserInfo `json:"user_info"`   // 包含更新字段的用户基本信息
	UpdateType string   `json:"update_type"` // 修改的类型，基础信息，私有信息
}

type UserUpdateResp struct {
}

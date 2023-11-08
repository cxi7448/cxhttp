package clFeishuBot


// 获取accessToken的结构体
type RespGetAccessToken struct {
	AppAccessToken string `json:"app_access_token"`
	Code uint32 `db:"code"`
	Expire uint32 `db:"expire"`
	Msg string `db:"msg"`
	TenantAccessToken string `db:"tenant_access_token"`
}
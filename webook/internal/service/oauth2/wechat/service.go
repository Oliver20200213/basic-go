package wechat

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// 需要先编码一下，在构造url路径的时候，当路径中包含特殊字符或非ASCII字符时，需要对其进行编码（:和//为特殊字符，中文为非ASCII字符）
var redirectURI = url.PathEscape("https://www.meoying.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string, state string) (domain.WeChatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewService(appId, appSecret string) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		// 依赖注入但是没有完全注入
		client: http.DefaultClient, //偷懒的写法
	}
}

// 不偷懒的写法，依赖注入的写法
func NewServiceV1(appId, appSecret string, client *http.Client) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
	}
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	// 构建微信扫码认证的url
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	//state := uuid.New()

	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string, state string) (domain.WeChatInfo, error) {
	// 构造url，通过携带code的url来获取access_token
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, redirectURI, code)
	// request请求写法一：
	//req,err := http.Get(target)
	// 写法二：
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	// 或者
	//req, err := http.NewRequest(http.MethodGet, target, nil)
	//if err != nil {
	//	return domain.WeChatInfo{}, err
	//}
	//req = req.WithContext(ctx) //这里会产生复制，坏处是性能极差，比如说你的URL很长，需要复制一遍性能就会差
	if err != nil {
		return domain.WeChatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WeChatInfo{}, err
	}

	// 只读一遍body
	decoder := json.NewDecoder(resp.Body) // 得到的是流式JSON数据
	var res Result
	err = decoder.Decode(&res) // 将JSON数据解码到res

	// 不优的写法，不推荐，ReadAll整个响应都读出来,因为Unmarshal的时候在读一遍，一共读两遍body
	//body, err := io.ReadAll(resp.Body)
	//err := json.Unmarshal(body,&res)

	if err != nil {
		return domain.WeChatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WeChatInfo{}, fmt.Errorf("微信返回错误响应，错误码：%d,错误信息：%s", res.ErrCode, res.ErrMsg)
	}
	return domain.WeChatInfo{
		OpenId:  res.Openid,
		UnionId: res.UnionID,
	}, nil
}

type Result struct {
	/*
		根据微信开放平台response的值来构造
		正确返回值：
		{
		"access_token":"ACCESS_TOKEN",
		"expires_in":7200,
		"refresh_token":"REFRESH_TOKEN",
		"openid":"OPENID",
		"scope":"SCOPE",
		"unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL"
		}
		错误返回值
		{"errcode":40029,"errmsg":"invalid code"}
	*/
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	Openid  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}

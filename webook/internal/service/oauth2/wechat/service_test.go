//go:build manual

package wechat

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// 手动跑的，提前验证代码
func Test_service_manual_VerifyCode(t *testing.T) {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量WECHAT_APP_ID")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_KEY")
	if !ok {
		panic("没有找到环境该变量WECHAT_APP_KEY")
	}

	svc := NewService(appId, appKey)
	res, err := svc.VerifyCode(context.Background(), "手动复制浏览器url里面的code", "state")
	require.NoError(t, err)
	t.Log(res)
}

func Test_service_VerifyCode(t *testing.T) {

}

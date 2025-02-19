package aliyunv1

import (
	"basic-go/webook/internal/service/sms"
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestService_SendV1(t *testing.T) {
	AccessKeyId, ok := os.LookupEnv("AccessKeyId")
	if !ok {
		t.Fatal("AccessKeyId not found")

	}
	AccessKeySecret, ok := os.LookupEnv("AccessKeySecret")
	if !ok {
		t.Fatal("AccessKeySecret not found")
	}
	config := &openapi.Config{
		AccessKeyId:     &AccessKeyId,
		AccessKeySecret: &AccessKeySecret,
	}
	config.Endpoint = ekit.ToPtr[string]("dysmsapi.aliyuncs.com")

	c, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(c, "go学习")
	testCase := []struct {
		name    string
		tplId   string
		params  []sms.NamedArg
		phone   string
		wantErr error
	}{
		{
			name:  "发送成功",
			tplId: "SMS_477385010",
			params: []sms.NamedArg{
				{Name: "code", Val: "123456"}, {Name: "time", Val: "10"},
			},
			phone:   "15562666678",
			wantErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err = s.SendV1(context.Background(), tc.tplId, tc.params, tc.phone)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}

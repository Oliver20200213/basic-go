package async

//
//import (
//	"basic-go/webook/internal/service/sms"
//	"context"
//)
//
//type SMSService struct {
//	svc  sms.Service
//	repo respository.SMSAsyncReqRepository
//}
//
//func NewSMSService() *SMSService {
//	return &SMSService{}
//}
//
//func (s *SMSService) StartAsync() {
//	go func() {
//		reqs := s.repo.Find没发出去的请求()
//		for _, req := range reqs {
//			// 在这里发送，并控制重试
//			s.svc.Send(ctx, req.biz, req.args, req.number...)
//		}
//	}()
//}
//
//func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
//	// 首先是正常路径
//	err := s.svc.Send(ctx, biz, args, numbers...)
//	if err != nil {
//		// 判定是不是崩溃了
//		if 崩溃了 {
//			s.repo.Store()
//		}
//	}
//}

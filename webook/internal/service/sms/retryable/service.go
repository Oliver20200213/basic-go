package retryable

// 这个要小心并发问题
//type Service struct {
//	svc sms.Service
//	// 重试次数
//	retryCnt int
//}
//
//func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) {
//	err := s.svc.Send(ctx, tpl, args, numbers...)
//	if err != nil && s.retryCnt < 10 {
//		err = s.svc.Send(ctx, tpl, args, numbers...)
//		s.retryCnt++
//	}
//}

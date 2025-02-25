package failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

// failover的第一种实现
//type FailoverSMSService struct {
//	svcs []sms.Service
//}

// failover的第二种实现
type FailoverSMSService struct {
	svcs []sms.Service

	idx uint64
}

func NewFailoverSMSService(svcs []sms.Service) *FailoverSMSService {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
		// 正常这边，输出日志
		// 要做好监控
		log.Println(err)
	}
	return errors.New("所有的服务商都发送失败了") // 走到这里表示全部轮询了一遍都没有成功，基本上可能是自己网络崩了
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers []string) error {
	// 我取下一个节点为发送节点
	//  atomic.AddUint64是一个原子操作函数，用于在并发环境下安全地对 uint64 类型的变量进行加法操作
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ { // 这个地方怎写都行，循环length也行，只要循环的次数满足就行
		// i%length 是为了防止索引越界，当i超过length-1时会导致索引越界
		// i%length 会将值限定在[0,length-1]
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			// 输出日志
		}
	}
	return errors.New("全部服务商都失败了")
}

// Atomic 原子操作是轻量级并发工具
// 面试中一种并发优化的思路，就是使用原子操作
// 原子操作在atomic包中，注意原子操作操作的都是指针
// 记住一个原则：任何变量的任何操作，在没有并发控制的情况下，都不是并发安全的
func Atomic() {
	var val int32 = 12
	// 原子读，不会读到修改了一半的数据
	val = atomic.LoadInt32(&val)
	println(val)
	// 原子写，即便别的Goroutine在别的CPU核行，也能立刻看到
	atomic.StoreInt32(&val, 12)
	// 原子自增， 返回的是自增后的结果
	newVal := atomic.AddInt32(&val, 1)
	println(newVal)
	// CAS操作
	// 如果val的值是13，就修改为15
	swapped := atomic.CompareAndSwapInt32(&val, 13, 15)
	println(swapped)
}

/*
调用实例
假设服务商列表为 [svcA, svcB, svcC]，初始 f.idx=0：
第一次调用 SendV1：
idx = 1（原子操作后）。
遍历服务商 svcB -> svcC -> svcA。
第二次调用 SendV1：
idx = 2。
遍历服务商 svcC -> svcA -> svcB。


为什么是循环的时候是 i < idx+length？
因为i是从1开始的，需要满足循环一遍也就循环三次
循环范围：
i = 1,2,3
循环的条件也就是i<1+3 即i<4

配合取余
f.svc[i%length]
取余之后索引的取值范围[0,length-1]即[0,2]
*/

package lock

import "sync"

// 优先使用RWMutex，优先加读锁
// 常用的并发优化手段，用读写锁来优化读锁
// 有经验的使用这种方式，一定要把接收器设置成指针
type LockDemo struct {
	lock sync.Mutex
}

// 这种方式不用显示的初始初始化
func NewLockDemo() *LockDemo {
	return &LockDemo{}
}

func (l *LockDemo) PanicDemo() {
	l.lock.Lock()
	// 在这中间panic了，无法释放锁
	panic("abc")
	l.lock.Unlock()
}

func (l *LockDemo) DeferDemo() {
	// 正确的用法
	l.lock.Lock()
	defer l.lock.Unlock()
}

// 不加指针会报错，会产生锁的复制，就会有两把锁
func (l LockDemo) NoPointerDemo() {
	// 正确的用法
	l.lock.Lock()
	defer l.lock.Unlock()
}

// 另一种写法,这样是没问题
// 新人推荐这种用法
type LockDemoV1 struct {
	lock *sync.Mutex //如果这里没用指针，接收器必须用指针，如果这里用了指针，接收器用不用指针无所谓
}

// 这种写法需要显示的初始化他
func NewLockDemoV1() *LockDemoV1 {
	return &LockDemoV1{
		// 如果不初始化，lock就是nil
		lock: &sync.Mutex{},
	}
}
func (l LockDemoV1) NoPointerDemo() {
	// 正确的用法
	l.lock.Lock()
	defer l.lock.Unlock()
}

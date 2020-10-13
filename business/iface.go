package business

type IModel interface {
	GetId() int
}

type IUser interface {
	GetId() int
}

type ICorp interface {
	GetId() int
	GetPlatformId() int
	IsPlatform() bool
	IsValid() bool
}

type IOrder interface {
	GetId() int
	GetBid() string
	GetDeductableMoney() int
}

type IResource interface {
	GetType() string //获得资源类型
	GetDeductionMoney(deductableMoney int) int //获得资源抵扣的金额
	GetPrice() int //获得资源价格
	GetPostage() int //获得资源的运费
	IsAllocated() bool //设置资源已申请成功
	ResetAllocation() //设置资源未申请
	SetAllocated() //设置资源未申请
	CanSplit() bool //是否可以切分给不同的supplier
	IsNeedLockWhenConsume() bool //在消费时是否需要加锁
	GetLockName() string //获得加锁时的锁名
	IsValid() error //验证资源是否合法
	ToMap() map[string]interface{} //转换为map
	SaveForOrder(order IOrder) error //存储资源自身
	GetRawResourceObject() interface{}
}

type IResourceAllocator interface {
	Allocate(resource IResource, newOrder IOrder) error //申请资源
	Release(resource IResource) //释放资源
}
package event_infra

import (
	"errors"
	"sync/atomic"

	iface "poprako-s/internal/domain/event"

	"github.com/panjf2000/ants"
	"go.uber.org/zap"
)

type handlerTable = *atomic.Pointer[map[iface.EventType][]iface.EventHandler]

type eventBusImpl struct {
	tbl        handlerTable
	asyncEvCh  chan iface.Event
	workerDone chan struct{}
}

func NewEventBus() (iface.EventBus, error) {
	// 初始化并存储空的分发映射
	tbl := new(atomic.Pointer[map[iface.EventType][]iface.EventHandler])
	initTbl := make(map[iface.EventType][]iface.EventHandler)

	tbl.Store(&initTbl)

	// 初始化工作池和相关的 channel
	const numWorkers = 64

	// 使用带缓冲的通道来避免阻塞，缓冲大小大于工作池的大小以保证可以短暂积压
	asyncEvCh := make(chan iface.Event, numWorkers*4)
	// 工作池的结束信号通道，使用带缓冲的通道以避免在发送结束信号时阻塞
	workerDone := make(chan struct{}, 1)

	pool, err := ants.NewPool(64)
	if err != nil {
		return nil, err
	}

	go asyncWorker(pool, handlerTable(tbl), asyncEvCh, workerDone)

	return &eventBusImpl{
		tbl:        handlerTable(tbl),
		asyncEvCh:  asyncEvCh,
		workerDone: workerDone,
	}, nil
}

func asyncWorker(
	pool *ants.Pool,
	mapper handlerTable,
	asyncEventCh <-chan iface.Event,
	workerDone <-chan struct{},
) {

	defer pool.Release()

	for {
		select {
		case ev := <-asyncEventCh:
			{
				// 收到任务信号，执行相应的任务
				m := (*mapper).Load()
				if m == nil {
					continue
				}

				handlers := (*m)[ev.EventType()]
				for _, handler := range handlers {
					if err := handler.Handle(ev); err != nil {
						// 出现错误，仅记录日志，不重试，继续处理下一个处理器
						zap.L().Error("[asyncWorker] 执行任务失败",
							zap.String("event_type", string(ev.EventType())),
							zap.Error(err),
						)
					}
				}
			}

		case <-workerDone:
			{
				// 收到结束信号，退出循环
				return
			}
		}
	}
}

func (b *eventBusImpl) Close() {
	// 发送结束信号，通知工作池中的工作协程退出
	b.workerDone <- struct{}{}
}

func (b *eventBusImpl) Pub(ev []iface.Event) error {
	m := (*b.tbl).Load()
	if m == nil {
		return errors.New("事件总线未正确初始化")
	}

	// 先按事件自己的发布类型拆分，同步事件先执行，异步事件后入队
	syn := make([]iface.Event, 0, 4)
	asyn := make([]iface.Event, 0, 4)

	for _, e := range ev {
		switch e.PubType() {
		case iface.PubTypeSync:
			syn = append(syn, e)
		case iface.PubTypeAsync:
			asyn = append(asyn, e)
		}
	}

	for _, e := range syn {
		handlers := (*m)[e.EventType()]
		for _, handler := range handlers {
			if err := handler.Handle(e); err != nil {
				return err
			}
		}
	}

	// 复用 PubAsync 逻辑
	if err := b.PubAsync(asyn); err != nil {
		return err
	}

	return nil
}

func (b *eventBusImpl) PubAsync(ev []iface.Event) error {
	// 无论事件自己的发布类型是什么，都强制按异步方式入队
	for _, e := range ev {
		select {
		case b.asyncEvCh <- e:
			// enqueued
		default:
			// 如果通道已满，返回错误，避免阻塞
			return errors.New("事件总线已过载，无法处理更多事件")
		}
	}

	return nil
}

func (b *eventBusImpl) Sub(h iface.EventHandler) error {
	// 先获取当前的分发映射
	m := (*b.tbl).Load()
	if m == nil {
		// 如果映射未初始化，返回错误
		return errors.New("事件总线未正确初始化")
	}

	// 创建新的分发映射，复制现有的映射内容
	newMap := make(map[iface.EventType][]iface.EventHandler)

	for k, v := range *m {
		newMap[k] = v
	}

	// 将新的处理器添加到对应事件类型的处理器列表中
	ty := h.EventType()
	newMap[ty] = append(newMap[ty], h)

	// 使用原子操作更新分发映射，确保线程安全
	b.tbl.Store(&newMap)

	return nil
}

func (b *eventBusImpl) SubUnsafe(h iface.EventHandler) error {
	// 直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	m := (*b.tbl).Load()
	if m == nil {
		return errors.New("事件总线未正确初始化")
	}

	// 直接修改当前的分发映射，适用于单线程环境或已知没有并发访问的情况
	ty := h.EventType()
	(*m)[ty] = append((*m)[ty], h)

	return nil
}

// FIXME: unsub 的 handler == h 是否足够？是否需要更严格的比较方式？

func (b *eventBusImpl) Unsub(h iface.EventHandler) error {
	// 先获取当前的分发映射
	m := (*b.tbl).Load()
	if m == nil {
		// 如果映射未初始化，返回错误
		return errors.New("事件总线未正确初始化")
	}

	// 创建新的分发映射，复制现有的映射内容
	newMap := make(map[iface.EventType][]iface.EventHandler)

	for k, v := range *m {
		newMap[k] = v
	}

	// 从对应事件类型的处理器列表中移除指定的处理器
	ty := h.EventType()

	handlers := newMap[ty]
	for i, handler := range handlers {
		if handler == h {
			// 找到要移除的处理器，使用切片操作移除它
			newMap[ty] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	// 使用原子操作更新分发映射，确保线程安全
	b.tbl.Store(&newMap)

	return nil
}

func (b *eventBusImpl) UnsubUnsafe(h iface.EventHandler) error {
	// 直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	m := (*b.tbl).Load()
	if m == nil {
		return errors.New("事件总线未正确初始化")
	}

	// 直接修改当前的分发映射，适用于单线程环境或已知没有并发访问的情况
	ty := h.EventType()

	handlers := (*m)[ty]

	for i, handler := range handlers {
		if handler == h {
			// 找到要移除的处理器，使用切片操作移除它
			(*m)[ty] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

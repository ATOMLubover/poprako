package event

import (
	"errors"
	"sync/atomic"

	"github.com/panjf2000/ants"
	"go.uber.org/zap"
)

// EventBus 是一个事件总线接口，定义了发布和订阅事件的方法
// 支持同步地消费事件（Pub）和异步地消费事件（PubAsync），以及订阅（Sub）和取消订阅（Unsub）事件处理器的方法
type EventBus interface {
	// Pub 使用同步方式发布事件，等待所有处理器完成后返回
	// 注意，它会截取调用时的一个 snapshot 来保证在事件处理过程中分发映射的一致性
	// 这意味着在事件处理过程中对订阅的修改不会影响当前事件的分发，但会影响后续事件的分发
	Pub(ev []Event) error
	// PubAsync 使用异步方式发布事件，将事件发送到工作池中处理，立即返回
	// 注意，它不会截取调用时的一个 snapshot，而是直接使用当前的分发映射来处理事件
	// 而且它不保证所有的事件都原子地入队成功，如果事件总线过载，可能会返回错误，导致部分事件无法处理
	PubAsync(ev []Event) error
	// Sub 订阅事件处理器，使用线程安全的方式修改分发映射，适用于多线程环境
	Sub(ty EventType, h EventHandler) error
	// SubUnsafe 订阅事件处理器，直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	SubUnsafe(ty EventType, h EventHandler) error
	// Unsub 取消订阅事件处理器，使用线程安全的方式修改分发映射，适用于多线程环境
	Unsub(ty EventType, h EventHandler) error
	// UnsubUnsafe 取消订阅事件处理器，直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	UnsubUnsafe(ty EventType, h EventHandler) error
}

type handlerTable = *atomic.Pointer[map[EventType][]EventHandler]

type eventBusImpl struct {
	tbl        handlerTable
	asyncEvCh  chan Event
	workerDone chan struct{}
}

func NewEventBus() (EventBus, error) {
	// 初始化并存储空的分发映射
	tbl := new(atomic.Pointer[map[EventType][]EventHandler])
	initTbl := make(map[EventType][]EventHandler)

	tbl.Store(&initTbl)

	// 初始化工作池和相关的 channel
	const numWorkers = 64

	// 使用带缓冲的通道来避免阻塞，缓冲大小大于工作池的大小以保证可以短暂积压
	asyncEvCh := make(chan Event, numWorkers*4)
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
	asyncEventCh <-chan Event,
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

func (b *eventBusImpl) Pub(ev []Event) error {
	// 按顺序执行事件处理器，等待每个处理器完成后再执行下一个
	m := (*b.tbl).Load()
	if m == nil {
		return errors.New("事件总线未正确初始化")
	}

	for _, e := range ev {
		handlers := (*m)[e.EventType()]
		for _, handler := range handlers {
			if err := handler.Handle(e); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *eventBusImpl) PubAsync(ev []Event) error {
	// 将事件发送到异步事件通道，由工作池中的协程处理
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

func (b *eventBusImpl) Sub(ty EventType, h EventHandler) error {
	// 先获取当前的分发映射
	m := (*b.tbl).Load()
	if m == nil {
		// 如果映射未初始化，返回错误
		return errors.New("事件总线未正确初始化")
	}

	// 创建新的分发映射，复制现有的映射内容
	newMap := make(map[EventType][]EventHandler)

	for k, v := range *m {
		newMap[k] = v
	}

	// 将新的处理器添加到对应事件类型的处理器列表中
	newMap[ty] = append(newMap[ty], h)

	// 使用原子操作更新分发映射，确保线程安全
	b.tbl.Store(&newMap)

	return nil
}

func (b *eventBusImpl) SubUnsafe(ty EventType, h EventHandler) error {
	// 直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	m := (*b.tbl).Load()
	if m == nil {
		return errors.New("事件总线未正确初始化")
	}

	// 直接修改当前的分发映射，适用于单线程环境或已知没有并发访问的情况
	(*m)[ty] = append((*m)[ty], h)

	return nil
}

// FIXME: unsub 的 handler == h 是否足够？是否需要更严格的比较方式？

func (b *eventBusImpl) Unsub(ty EventType, h EventHandler) error {
	// 先获取当前的分发映射
	m := (*b.tbl).Load()
	if m == nil {
		// 如果映射未初始化，返回错误
		return errors.New("事件总线未正确初始化")
	}

	// 创建新的分发映射，复制现有的映射内容
	newMap := make(map[EventType][]EventHandler)

	for k, v := range *m {
		newMap[k] = v
	}

	// 从对应事件类型的处理器列表中移除指定的处理器
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

func (b *eventBusImpl) UnsubUnsafe(ty EventType, h EventHandler) error {
	// 直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	m := (*b.tbl).Load()
	if m == nil {
		return errors.New("事件总线未正确初始化")
	}

	// 直接修改当前的分发映射，适用于单线程环境或已知没有并发访问的情况
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

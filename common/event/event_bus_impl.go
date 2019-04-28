/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package event

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/qlcchain/go-qlc/common/hashmap"
)

const DefaultQueueSize = 100

// DefaultEventBus
type DefaultEventBus struct {
	handlers  map[string][]*eventHandler
	wg        sync.WaitGroup
	queueSize int
	lock      sync.RWMutex
}

type eventOption struct {
	async bool
}

type eventHandler struct {
	callBack reflect.Value
	ctx      context.Context
	cancel   context.CancelFunc
	queue    chan []reflect.Value
	option   *eventOption
}

// New returns new DefaultEventBus with empty handlers.
func New() EventBus {
	b := &DefaultEventBus{
		handlers:  make(map[string][]*eventHandler),
		queueSize: DefaultQueueSize,
	}
	return EventBus(b)
}

func NewEventBus(queueSize int) EventBus {
	b := &DefaultEventBus{
		handlers:  make(map[string][]*eventHandler),
		queueSize: queueSize,
	}
	return EventBus(b)
}

var (
	once  sync.Once
	eb    EventBus
	cache *hashmap.HashMap
)

func init() {
	cache = hashmap.New(50)
}

func SimpleEventBus() EventBus {
	once.Do(func() {
		eb = New()
	})

	return eb
}

func GetEventBus(id string) EventBus {
	if len(id) == 0 {
		return SimpleEventBus()
	}

	if v, ok := cache.GetStringKey(id); ok {
		return v.(EventBus)
	} else {
		eb := New()
		cache.Set(id, eb)
		return eb
	}
}

// doSubscribe handles the subscription logic and is utilized by the public Subscribe functions
func (eb *DefaultEventBus) doSubscribe(topic string, fn interface{}, option *eventOption) error {
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}

	ctx, cancel := context.WithCancel(context.Background())

	handler := &eventHandler{
		callBack: reflect.ValueOf(fn), ctx: ctx, cancel: cancel,
		queue: make(chan []reflect.Value, eb.queueSize), option: option,
	}

	eb.wg.Add(1)
	go func() {
		defer eb.wg.Done()
		for {
			select {
			case args, ok := <-handler.queue:
				if ok {
					if handler.option.async {
						go handler.callBack.Call(args)
					} else {
						handler.callBack.Call(args)
					}
				}
			case <-handler.ctx.Done():
				return
			}
		}
	}()
	eb.lock.Lock()
	defer eb.lock.Unlock()

	eb.handlers[topic] = append(eb.handlers[topic], handler)
	return nil
}

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (eb *DefaultEventBus) Subscribe(topic string, fn interface{}) error {
	return eb.doSubscribe(topic, fn, &eventOption{async: true})
}

// HasCallback returns true if exists any callback subscribed to the topic.
func (eb *DefaultEventBus) HasCallback(topic string) bool {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	_, ok := eb.handlers[topic]
	if ok {
		return len(eb.handlers[topic]) > 0
	}
	return false
}

// Close unsubscribe all handlers from given topic
func (eb *DefaultEventBus) Close(topic string) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	for topicPattern, handlers := range eb.handlers {
		if len(handlers) > 0 && MatchSimple(topicPattern, topic) {
			for i := range handlers {
				eb.removeHandler(topic, i)
			}
			delete(eb.handlers, topicPattern)
		}
	}
}

// Unsubscribe removes callback defined for a topic.
// Returns error if there are no callbacks subscribed to the topic.
func (eb *DefaultEventBus) Unsubscribe(topic string, handler interface{}) error {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	counter := 0
	for topicPattern, handlers := range eb.handlers {
		if len(handlers) > 0 && MatchSimple(topicPattern, topic) {
			counter++
			eb.removeHandler(topicPattern, eb.findHandlerIdx(topicPattern, reflect.ValueOf(handler)))
		}
	}

	if counter > 0 {
		return nil
	}
	return fmt.Errorf("topic %s doesn't exist", topic)
}

// Publish executes callback defined for a topic. Any additional argument will be transferred to the callback.
func (eb *DefaultEventBus) Publish(topic string, args ...interface{}) {
	rArgs := eb.setUpPublish(topic, args...)

	eb.lock.RLock()
	defer eb.lock.RUnlock()
	for topicPattern, handlers := range eb.handlers {
		if 0 < len(handlers) && MatchSimple(topicPattern, topic) {
			for _, handler := range handlers {
				handler.queue <- rArgs
			}
		}
	}
}

func (eb *DefaultEventBus) removeHandler(topic string, idx int) {
	if _, ok := eb.handlers[topic]; !ok {
		return
	}
	l := len(eb.handlers[topic])

	if !(0 <= idx && idx < l) {
		return
	}

	h := eb.handlers[topic][idx]
	h.cancel()
	close(h.queue)

	eb.handlers[topic] = append(eb.handlers[topic][:idx], eb.handlers[topic][idx+1:]...)
}

func (eb *DefaultEventBus) findHandlerIdx(topic string, callback reflect.Value) int {
	if _, ok := eb.handlers[topic]; ok {
		for idx, handler := range eb.handlers[topic] {
			if handler.callBack == callback {
				return idx
			}
		}
	}
	return -1
}

func (eb *DefaultEventBus) setUpPublish(topic string, args ...interface{}) []reflect.Value {
	passedArguments := make([]reflect.Value, 0)
	for _, arg := range args {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}
	return passedArguments
}

// WaitAsync waits for all async callbacks to complete
func (eb *DefaultEventBus) WaitAsync() {
	eb.wg.Wait()
}

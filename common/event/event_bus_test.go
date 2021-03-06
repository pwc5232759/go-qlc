/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package event

import (
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bus := New()
	if bus == nil {
		t.Log("New EventBus not created!")
		t.Fail()
	}
}

func TestHasCallback(t *testing.T) {
	bus := New()
	err := bus.Subscribe("topic", func() {})
	if err != nil {
		t.Fatal(err)
	}
	if bus.HasCallback("topic_topic") {
		t.Fail()
	}
	if !bus.HasCallback("topic") {
		t.Fail()
	}
}

func TestSubscribe(t *testing.T) {
	bus := New()
	if bus.Subscribe("topic", func() {}) != nil {
		t.Fail()
	}
	if bus.Subscribe("topic", "String") == nil {
		t.Fail()
	}
}

func TestSubscribeOnce(t *testing.T) {
	bus := New()
	if bus.SubscribeOnce("topic", func() {}) != nil {
		t.Fail()
	}
	if bus.SubscribeOnce("topic", "String") == nil {
		t.Fail()
	}
}

func TestSubscribeOnceAndManySubscribe(t *testing.T) {
	bus := New()
	event := "topic"
	flag := 0
	fn := func(f int) {
		flag += 1
		t.Log(f, flag)
	}
	_ = bus.SubscribeOnce(event, fn)
	_ = bus.Subscribe(event, fn)
	_ = bus.Subscribe(event, fn)
	bus.Publish(event, flag)

	if flag != 3 {
		t.Fail()
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := New()
	handler := func() {}
	err := bus.Subscribe("topic", handler)
	if err != nil {
		t.Fatal(err)
	}
	if bus.Unsubscribe("topic", handler) != nil {
		t.Fail()
	}
	if bus.Unsubscribe("topic", handler) == nil {
		t.Fail()
	}
}

func TestPublish(t *testing.T) {
	bus := New()
	err := bus.Subscribe("topic", func(a int, b int) {
		if a != b {
			t.Fail()
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	bus.Publish("topic", 10, 10)
}

func TestSubscribeOnceAsync(t *testing.T) {
	results := make([]int, 0)

	bus := New()
	err := bus.SubscribeOnceAsync("topic", func(a int, out *[]int) {
		*out = append(*out, a)
	})
	if err != nil {
		t.Fatal(err)
	}

	bus.Publish("topic", 10, &results)
	bus.Publish("topic", 10, &results)

	bus.WaitAsync()

	if len(results) != 1 {
		t.Fail()
	}

	if bus.HasCallback("topic") {
		t.Fail()
	}
}

func TestSubscribeAsyncTransactional(t *testing.T) {
	results := make([]int, 0)

	bus := New()
	err := bus.SubscribeAsync("topic", func(a int, out *[]int, dur string) {
		sleep, _ := time.ParseDuration(dur)
		time.Sleep(sleep)
		*out = append(*out, a)
	}, true)

	if err != nil {
		t.Fatal(err)
	}

	bus.Publish("topic", 1, &results, "1s")
	bus.Publish("topic", 2, &results, "0s")

	bus.WaitAsync()

	if len(results) != 2 {
		t.Fail()
	}

	if results[0] != 1 || results[1] != 2 {
		t.Fail()
	}
}

func TestSubscribeAsyncTopic(t *testing.T) {
	results := make([]int, 0)

	type result struct {
		A int
	}

	bus := New()
	fn := func(a *result, out *[]int, dur string) {
		sleep, _ := time.ParseDuration(dur)
		time.Sleep(sleep)
		*out = append(*out, a.A)
		t.Log(&a, a.A)
	}
	err := bus.SubscribeAsync("topic:*", fn, true)

	if err != nil {
		t.Fatal(err)
	}
	bus.Publish("topic:", &result{3}, &results, "1s")
	bus.Publish("topic:1st", &result{1}, &results, "1s")
	bus.Publish("topic:2nd", &result{2}, &results, "0s")

	bus.WaitAsync()

	if len(results) != 3 {
		t.Fail()
	}

	if results[0] != 3 || results[1] != 1 || results[2] != 2 {
		t.Fail()
	}
}

func TestSubscribeAsync(t *testing.T) {
	results := make(chan int)

	bus := New()
	err := bus.SubscribeAsync("topic", func(a int, out chan<- int) {
		out <- a
	}, false)

	if err != nil {
		t.Fatal(err)
	}

	bus.Publish("topic", 1, results)
	bus.Publish("topic", 2, results)

	numResults := 0

	go func() {
		for range results {
			numResults++
		}
	}()

	bus.WaitAsync()

	time.Sleep(10 * time.Millisecond)

	if numResults != 2 {
		t.Fail()
	}
}

func TestSimpleEventBus(t *testing.T) {
	eb1 := SimpleEventBus()
	eb2 := SimpleEventBus()

	if eb1 != eb2 {
		t.Fatal("eb1!=eb2")
	}
}

func TestGetEventBus(t *testing.T) {
	eb0 := SimpleEventBus()
	eb1 := GetEventBus("")
	if eb0 != eb1 {
		t.Fatal("invalid default eb")
	}

	id1 := "111111"
	eb2 := GetEventBus(id1)
	eb3 := GetEventBus(id1)

	if eb2 != eb3 {
		t.Fatal("invalid eb of same id")
	}

	id2 := "222222"
	eb4 := GetEventBus(id2)
	if eb3 == eb4 {
		t.Fatal("invalid eb of diff ids")
	}
}

func TestAsyncGetEvent(t *testing.T) {
	var eb0 EventBus
	var eb1 EventBus
	id1 := "sssdfasdfasd"
	eb0 = GetEventBus(id1)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(id string) {
		defer wg.Done()
		eb1 = GetEventBus(id)
	}(id1)

	wg.Wait()

	if eb0 == nil || eb1 == nil || eb0 != eb1 {
		t.Fatalf("invalid eb0 %p, eb1 %p", eb0, eb1)
	}

}

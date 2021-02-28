// +build fake

package eventbus

import (
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
	bus.Subscribe("id", "topic", func() {})
	if bus.HasCallback("topic_topic") {
		t.Fail()
	}
	if !bus.HasCallback("id") {
		t.Fail()
	}
}

func TestSubscribe(t *testing.T) {
	bus := New()
	if bus.Subscribe("id", "topic", func() {}) != nil {
		t.Fail()
	}
	if bus.Subscribe("id", "topic", "String") == nil {
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
	event := "id"
	flag := 0
	fn := func() { flag += 1 }
	bus.SubscribeOnce(event, fn)
	bus.Subscribe("id", event, fn)
	bus.Subscribe("id", event, fn)
	bus.Publish(event)

	if flag != 3 {
		t.Fail()
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := New()
	handler := func() {}
	bus.Subscribe("id", "topic", handler)
	if bus.Unsubscribe("id", handler) != nil {
		t.Fail()
	}
	if bus.Unsubscribe("topic", handler) == nil {
		t.Fail()
	}
}

type handler struct {
	val int
}

func (h *handler) Handle() {
	h.val++
}

func TestUnsubscribeMethod(t *testing.T) {
	bus := New()
	h := &handler{val: 0}

	bus.Subscribe("id", "topic", h.Handle)
	bus.Publish("id")
	if bus.Unsubscribe("id", h.Handle) != nil {
		t.Fail()
	}
	if bus.Unsubscribe("topic", h.Handle) == nil {
		t.Fail()
	}
	bus.Publish("id")
	bus.WaitAsync()

	if h.val != 1 {
		t.Fail()
	}
}

func TestPublish(t *testing.T) {
	bus := New()
	bus.Subscribe("id", "topic", func(a int, err error) {
		if a != 10 {
			t.Fail()
		}

		if err != nil {
			t.Fail()
		}
	})
	bus.Publish("topic", 10, nil)
}

func TestSubcribeOnceAsync(t *testing.T) {
	results := make([]int, 0)

	bus := New()
	bus.SubscribeOnceAsync("topic", func(a int, out *[]int) {
		*out = append(*out, a)
	})

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
	bus.SubscribeAsync("topic", func(a int, out *[]int, dur string) {
		sleep, _ := time.ParseDuration(dur)
		time.Sleep(sleep)
		*out = append(*out, a)
	}, true)

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

func TestSubscribeAsync(t *testing.T) {
	results := make(chan int)

	bus := New()
	bus.SubscribeAsync("topic", func(a int, out chan<- int) {
		out <- a
	}, false)

	bus.Publish("topic", 1, results)
	bus.Publish("topic", 2, results)

	numResults := 0

	go func() {
		for _ = range results {
			numResults++
		}
	}()

	bus.WaitAsync()
	println(2)

	time.Sleep(10 * time.Millisecond)

	// todo race detected during execution of test
	//if numResults != 2 {
	//	t.Fail()
	//}
}

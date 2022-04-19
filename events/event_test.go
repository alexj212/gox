package events_test

import (
	"github.com/alexj212/gox/events"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestNewEvent(t *testing.T) {
	event := events.New[string]()
	if event == nil {
		t.Error("Expected event to be created")
	}
}

func TestSubscribeAndDispatch(t *testing.T) {
	s := is.New(t)
	var results []int

	event := events.New[int]()
	event.Subscribe(func(data int) {
		results = append(results, data)
	})
	event.Dispatch(1)
	s.Equal(1, len(results))
}

func TestSubscribeMany(t *testing.T) {
	s := is.New(t)
	var results []int

	event := events.New[int]()
	event.Subscribe(func(data int) {
		results = append(results, data)
	})
	event.Subscribe(func(data int) {
		results = append(results, data)
	})
	event.Dispatch(1)
	s.Equal(2, len(results))
}

func TestSubscribeOnce(t *testing.T) {
	s := is.New(t)
	var results []int

	event := events.New[int]()
	event.SubscribeOnce(func(data int) {
		results = append(results, data)
	})
	event.Dispatch(1)
	event.Dispatch(2)
	s.Equal(1, len(results))
}

func TestUnsubscribe(t *testing.T) {
	s := is.New(t)
	var results []int

	event := events.New[int]()
	unsubscribe := event.Subscribe(func(data int) {
		results = append(results, data)
	})
	unsubscribe()
	event.Dispatch(1)
	s.Equal(0, len(results))
}

func TestSubscribeAsync(t *testing.T) {
	s := is.New(t)
	var results []int

	event := events.New[int]()
	event.SubscribeAsync(func(data int) {
		results = append(results, data)
	})
	event.Dispatch(1)
	event.Wait()
	s.Equal(1, len(results))
}

func TestSubscribeAsyncTransactional(t *testing.T) {
	s := is.New(t)
	var results []int
	type Data struct {
		dur   time.Duration
		value int
	}
	event := events.New[Data]()
	event.SubscribeAsync(func(data Data) {
		time.Sleep(data.dur)
		results = append(results, data.value)
	}, true)

	event.Dispatch(Data{dur: time.Millisecond * 100, value: 1})
	event.Dispatch(Data{dur: 0, value: 2})
	event.Dispatch(Data{dur: time.Millisecond * 100, value: 3})
	event.Dispatch(Data{dur: 10, value: 4})
	event.Dispatch(Data{dur: 0, value: 5})

	event.Wait()

	s.Equal(5, len(results))
	for i, v := range results {
		s.Equal(i+1, v)
	}
}

func TestSubscribeOnceAsync(t *testing.T) {
	s := is.New(t)
	var results []int

	event := events.New[int]()
	event.SubscribeOnceAsync(func(data int) {
		results = append(results, data)
	})
	event.Dispatch(1)
	event.Dispatch(2)
	event.Wait()
	s.Equal(1, len(results))
}

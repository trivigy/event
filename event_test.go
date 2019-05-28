package event

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EventSuite struct {
	suite.Suite
}

func (r *EventSuite) TestIsSet() {
	event := New()
	assert.False(r.T(), event.IsSet())
	event.Set()
	assert.True(r.T(), event.IsSet())
	event.Set()
	assert.True(r.T(), event.IsSet())
	event.Clear()
	assert.False(r.T(), event.IsSet())
	event.Clear()
	assert.False(r.T(), event.IsSet())
}

func (r *EventSuite) TestNotify() {
	event := New()
	r.checkNotify(event)
	event.Set()
	event.Clear()
	r.checkNotify(event)
}

func (r *EventSuite) checkNotify(event *Event) {
	mutex1 := sync.Mutex{}
	results1 := make([]bool, 0)

	mutex2 := sync.Mutex{}
	results2 := make([]bool, 0)

	start := sync.WaitGroup{}
	finish := sync.WaitGroup{}
	N := 5
	start.Add(N)
	finish.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			start.Done()
			val1 := event.Wait(nil)
			mutex1.Lock()
			results1 = append(results1, val1)
			mutex1.Unlock()

			val2 := event.Wait(nil)
			mutex2.Lock()
			results2 = append(results2, val2)
			mutex2.Unlock()
			finish.Done()
		}()
	}

	start.Wait()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(r.T(), 0, len(results1))
	event.Set()
	finish.Wait()

	expected := make([]bool, N)
	for i := range expected {
		expected[i] = true
	}

	assert.Equal(r.T(), expected, results1)
	assert.Equal(r.T(), expected, results2)
}

type outcome struct {
	result bool
	lapse  time.Duration
}

func (r *EventSuite) TestTimeout() {
	event := New()

	N := 5
	var mutex1 sync.Mutex
	var results1 []bool

	var mutex2 sync.Mutex
	var results2 []outcome

	var expected []bool
	finish := sync.WaitGroup{}
	f := func() {
		ctx1, cancel1 := context.WithTimeout(context.Background(), 0*time.Millisecond)
		defer cancel1()
		val1 := event.Wait(ctx1)
		mutex1.Lock()
		results1 = append(results1, val1)
		mutex1.Unlock()

		ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel2()
		t1 := time.Now()
		r := event.Wait(ctx2)
		t2 := time.Now()
		val2 := outcome{r, t2.Sub(t1)}
		mutex2.Lock()
		results2 = append(results2, val2)
		mutex2.Unlock()

		finish.Done()
	}

	results1 = make([]bool, 0)
	results2 = make([]outcome, 0)
	finish.Add(N)
	for i := 0; i < N; i++ {
		go f()
	}
	finish.Wait()
	expected = make([]bool, N)
	for i := range expected {
		expected[i] = false
	}
	assert.Equal(r.T(), expected, results1)
	for _, o := range results2 {
		assert.False(r.T(), o.result)
		assert.True(r.T(), o.lapse >= 0.6*500*time.Millisecond)
		assert.True(r.T(), o.lapse < 1.1*500*time.Millisecond)
	}

	results1 = make([]bool, 0)
	results2 = make([]outcome, 0)
	event.Set()
	finish.Add(N)
	for i := 0; i < N; i++ {
		go f()
	}
	finish.Wait()
	expected = make([]bool, N)
	for i := range expected {
		expected[i] = true
	}
	assert.Equal(r.T(), expected, results1)
	for _, o := range results2 {
		assert.True(r.T(), o.result)
	}
}

func (r *EventSuite) TestSetAndClear() {
	event := New()
	mutex := sync.Mutex{}
	results := make([]bool, 0)
	timeout := 250 * time.Millisecond

	finish := sync.WaitGroup{}
	N := 5
	finish.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), timeout*4)
			defer cancel()

			val := event.Wait(ctx)
			mutex.Lock()
			results = append(results, val)
			mutex.Unlock()
			finish.Done()
		}()
	}

	time.Sleep(timeout)
	event.Set()
	event.Clear()
	finish.Wait()
	expected := make([]bool, N)
	for i := range expected {
		expected[i] = true
	}
	assert.Equal(r.T(), expected, results)
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventSuite))
}

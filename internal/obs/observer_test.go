package obs

import (
	"math"
	"testing"
	"time"
)

type testStruct struct {
	testing.TB
	expected int
	actual int
}

type testEvent struct {
	int
}

func (t testStruct) Handle(event Event) {
	//utils.Equals(t, t.expected, event.Value.(testEvent).int)
}

func TestWatch_Full(t *testing.T) {
	w := NewWatch()
	ts := testStruct{t, 27, 0}
	w.Register(ts)
	w.NotifyAll(Event{Value: testEvent{27}})
}

func TestWatch_Deadlock(t *testing.T) {
	w := NewWatch()
	ts := testStruct{t, 27, 0}
	s := 2000
	c1 := make(chan bool)
	c2 := make(chan bool)
	t1 := time.Now()
	go func() {
		for i := 0; i < s; i++ {
			w.Register(ts)
			time.Sleep(1*time.Millisecond)
		}
		c1 <- true
	}()
	go func() {
		for i := 0; i < s; i++ {
			w.NotifyAll(Event{Value: testEvent{27}})
			time.Sleep(1*time.Millisecond)
		}
		c2 <- true
	}()
	<-c1
	<-c2
	elapsedTime := time.Since(t1).Milliseconds()
	t.Log(elapsedTime)
	dif := math.Abs(float64(s)-float64(elapsedTime))
	if dif > float64(s/10) { // 200 milliseconds
		t.Fatalf("Elapsed time: %v - %v", elapsedTime, dif)
	}
}

func TestWatch_Deregister(t *testing.T) {

}

func TestWatch_Listen(t *testing.T) {

}

func TestWatch_NotifyAll(t *testing.T) {

}

func TestWatch_Register(t *testing.T) {

}

func TestWatch_callbackOnNotification(t *testing.T) {

}


func BenchmarkObserver_NotifyAll(b *testing.B) {

	w := NewWatch()

	for i := 0; i < b.N; i++ {
		ts := testStruct{b, i, 0}
		w.Register(ts)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.NotifyAll(Event{Value: testEvent{i}})
	}
}
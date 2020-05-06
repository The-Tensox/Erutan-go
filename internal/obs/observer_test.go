package obs

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type a struct {}
func (a a) Handle(event Event) {}
func (a a) Priority() int {return 0}

type b struct {}
func (b b) Handle(event Event) {}
func (b b) Priority() int {return 2}

type c struct {}
func (c c) Handle(event Event) {}
func (c c) Priority() int {return 56}


func TestObserver_Priority(t *testing.T) {
	w := Watch{}
	a := a{}
	b := b{}
	c := c{}
	w.Register(c)
	w.Register(a)
	w.Register(b)
	equals(t, c, w.observers[0])
	equals(t, b, w.observers[1])
	equals(t, a, w.observers[2])
}
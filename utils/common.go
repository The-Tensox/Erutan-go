package utils

import (
	"fmt"
	"github.com/The-Tensox/erutan/cfg"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
)

const timeFormat = "03:04:05 PM"

func tsToTime(ts *timestamp.Timestamp) time.Time {
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		return time.Now()
	}
	return t.In(time.Local)
}

func ClientLogf(ts time.Time, format string, args ...interface{}) {
	log.Printf("[%s] <<Client>>: "+format, append([]interface{}{ts.Format(timeFormat)}, args...)...)
}

func ServerLogf(ts time.Time, format string, args ...interface{}) {
	log.Printf("[%s] <<Server>>: "+format, append([]interface{}{ts.Format(timeFormat)}, args...)...)
}

func MessageLog(ts time.Time, name, msg string) {
	log.Printf("[%s] %s: %s", ts.Format(timeFormat), name, msg)
}

func DebugLogf(format string, args ...interface{}) {
	if !cfg.Global.DebugMode {
		return
	}
	// Add more information  about the log, such as file name, function ...
	pc := make([]uintptr, 10)  // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	if s := strings.Split(file, "/"); len(s)>0 {
		file = s[len(s)-1] // Only keep the file name, drop the path
	}
	var functionName string
	if s := strings.Split(f.Name(), "/"); len(s)>0 {
		functionName = s[len(s)-1] // Only keep the package.(class).function
	}
	file = fmt.Sprintf("%s - %s - L%d", functionName, file, line)
	formattedString := append([]interface{}{time.Now().Format(timeFormat)}, []interface{}{file}...)
	formattedString = append(formattedString, args...)
	log.Printf("[%s] - [%s] <<Debug>>: "+format, formattedString...)
}

func SignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		DebugLogf("listening for shutdown signal")
		<-sigs
		DebugLogf("shutdown signal received")
		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// RandomString generates a random string of 4 bytes
func RandomString() string {
	str := make([]byte, 4)
	rand.Read(str)
	return fmt.Sprintf("%x", str)
}

func RandFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GetProtoTime() float64 {
	return float64(ptypes.TimestampNow().Seconds)*math.Pow(10, 9) + float64(ptypes.TimestampNow().Nanos)
}

func DoEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

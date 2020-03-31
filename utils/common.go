package utils

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"
	"testing"
	"time"

	erutan "github.com/The-Tensox/erutan/protobuf"
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
	if !Config.DebugMode {
		return
	}
	log.Printf("[%s] <<Debug>>: "+format, append([]interface{}{time.Now().Format(timeFormat)}, args...)...)
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

func RandomPositionInsideCircle(center *erutan.NetVector2, radius float64) *erutan.NetVector3 {
	return &erutan.NetVector3{X: RandFloats(-radius+center.X, radius+center.X), Y: 0, Z: RandFloats(-radius+center.Y, radius+center.Y)}
}

func RandomPositionInsideSphere(center *erutan.NetVector3, radius float64) *erutan.NetVector3 {
	return &erutan.NetVector3{X: RandFloats(-radius+center.X, radius+center.X),
		Y: RandFloats(-radius+center.Y, radius+center.Y),
		Z: RandFloats(-radius+center.Y, radius+center.Y)}
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

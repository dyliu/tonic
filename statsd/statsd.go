package statsd

import (
	"fmt"
	"time"

	"github.com/lingmiaotech/tonic/configs"
	"github.com/lingmiaotech/tonic/logging"
	"gopkg.in/alexcesaro/statsd.v2"
)

type InstanceClass struct {
	AppName string
	Enabled bool
	Client  *statsd.Client
}

type Timer struct {
	start time.Time
}

var Instance InstanceClass

func Increment(bucket string) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		logging.GetDefaultLogger().Infof("[STATSD] key=%s count=1", b)
		return
	}
	Instance.Client.Increment(b)
}

// Timing takes bucket name and delta in milliseconds
func Timing(bucket string, delta int) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		logging.GetDefaultLogger().Infof("[STATSD] key=%s time_delta=%d(ms)", b, delta)
		return
	}
	Instance.Client.Timing(b, delta)
}

func Count(bucket string, n int) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		logging.GetDefaultLogger().Infof("[STATSD] key=%v count=%d", b, n)
		return
	}
	Instance.Client.Count(b, n)
}

func Gauge(bucket string, n int) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		logging.GetDefaultLogger().Infof("[STATSD] key=%v gauge=%d", b, n)
		return
	}
	Instance.Client.Gauge(b, n)
}

func getBucket(bucket string) string {
	return fmt.Sprintf("%v.%v", Instance.AppName, bucket)
}

func NewTimer() Timer {
	return Timer{start: time.Now()}
}

func NewCustomTimer(t time.Time) Timer {
	return Timer{start: t}
}

// Send sends the time elapsed since the creation of the Timing.
func (t Timer) Send(bucket string) {
	Timing(bucket, int(t.Duration()/time.Millisecond))
}

// Duration returns the time elapsed since the creation of the Timing.
func (t Timer) Duration() time.Duration {
	return time.Now().Sub(t.start)
}

func InitStatsd() error {

	Instance.AppName = configs.GetString("app_name")
	Instance.Enabled = configs.GetBool("statsd.enabled")

	if !Instance.Enabled {
		return nil
	}

	host := configs.GetString("statsd.host")
	port := configs.GetString("statsd.port")
	address := fmt.Sprintf("%v:%v", host, port)
	c, err := statsd.New(statsd.Address(address))
	if err != nil || c == nil {
		return err
	}

	Instance.Client = c
	return nil

}

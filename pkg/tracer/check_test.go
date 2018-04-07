package tracer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckDuration(t *testing.T) {
	assert := assert.New(t)

	// a valid duration
	c := &check{Duration: "1h"}
	assert.Equal(1*time.Hour, c.duration(), "should be a valid duration")

	// an invalid duration
	c = &check{Duration: "asdf"}
	assert.Equal(1*time.Hour, c.duration(), "should return the default 1h duration")
}

func TestCheckMonitor(t *testing.T) {
	assert := assert.New(t)
	c := &check{Name: "should-not-panic", Duration: "1h", Panic: "touch /tmp/should-not-panic.tracer"}
	go func(c *check) {
		<-time.After(1 * time.Second)
		c.ok()
	}(c)
	assert.Equal("OK", c.monitor(time.Hour))
}

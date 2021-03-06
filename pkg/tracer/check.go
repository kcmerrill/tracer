package tracer

import (
	"bytes"
	"os/exec"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
)

// check contains information needed to monitor a trace
type check struct {
	Name      string    `json:"name"`
	Panic     string    `json:"panic"`
	Duration  string    `json:"duration"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	completed chan string
}

func (c *check) cancel() {
	c.Status = "Cancelled"
	c.completed <- c.Status
}

func (c *check) ok() {
	c.Status = "OK"
	c.completed <- c.Status
}

func (c *check) parsePanic(panic string) string {
	template := template.Must(template.New("translated").Funcs(sprig.TxtFuncMap()).Parse(panic))
	b := new(bytes.Buffer)
	err := template.Execute(b, c)
	if err != nil {
		log.WithFields(log.Fields{
			"name":  c.Name,
			"panic": panic,
		}).Error("Unable to parse panic")
		return panic
	}
	return b.String()
}

func initCheck(name, duration, panic string) *check {
	c := &check{
		Name:      name,
		Duration:  duration,
		Panic:     panic,
		Status:    "Pending",
		Created:   time.Now(),
		completed: make(chan string, 1),
	}

	log.WithFields(log.Fields{
		"check":    c.Name,
		"duration": c.Duration,
		"created":  c.Created,
	}).Info("Creating check")

	// validate the duration
	dur, err := time.ParseDuration(c.Duration)
	if err != nil {
		// default to 3 hours
		dur = 3 * time.Hour
		log.WithFields(log.Fields{
			"duration": c.Duration,
			"default":  "3h",
		}).Error("Unable to parse duration")
	}

	go func(c *check) {
		select {
		case status := <-c.completed:
			log.WithFields(log.Fields{
				"check":    c.Name,
				"duration": c.Duration,
				"created":  c.Created,
			}).Info(status)
			// whew! good to go
			return
		case <-time.After(dur):
			parsedPanic := c.parsePanic(c.Panic)
			log.WithFields(log.Fields{
				"check":    c.Name,
				"duration": c.Duration,
				"panic":    parsedPanic,
				"created":  c.Created,
			}).Error("Panicking ...")

			cmd := exec.Command("bash", "-c", parsedPanic)
			cmd.CombinedOutput()
		}
	}(c)
	return c
}

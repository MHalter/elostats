// eloextplugin

package eloextplugin

import (
	"context"
	"math/rand"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type EloPerfData struct {
	ValueName       string          `toml:"value_name"`
	Min             int64           `toml:"min"`
	Max             int64           `toml:"max"`
	SampleFrequency config.Duration `toml:"sample_frequency"`
	ctx             context.Context
	cancel          context.CancelFunc

	Log telegraf.Logger `toml:"-"`
}

func init() {
	inputs.Add("eloextplugin", func() telegraf.Input {
		return &EloPerfData{
			ValueName:       "checkoutSord",
			Min:             0,
			Max:             100,
			SampleFrequency: config.Duration(1 * time.Second),
		}
	})
}

func (r *EloPerfData) Init() error {
	return nil
}

func (r *EloPerfData) SampleConfig() string {
	r.Log.Infof("SampleConfig called")
	return `
  ## Generates random numbers
	[inputs.eloextplugin]
	# the name of the measurement to write out to.
	# value_name = "checkoutSord"
	# min = 0
	# max = 100
	# sample_frequency = "1s"
`
}

func (r *EloPerfData) Description() string {
	return "Generates a random number"
}

func (r *EloPerfData) Gather(a telegraf.Accumulator) error {
	ticker := time.NewTicker(time.Duration(r.SampleFrequency))
	defer ticker.Stop()

	for range ticker.C {
		r.sendMetric(a)
	}

	return nil
}

// // provide the extra functions so we can also run as a service input.
// func (r *RandomNumberGenerator) Start(a telegraf.Accumulator) error {
// 	println("Started as service")
// 	r.ctx, r.cancel = context.WithCancel(context.Background())
// 	go func() {
// 		t := time.NewTicker(r.SampleFrequency)
// 		for {
// 			select {
// 			case <-r.ctx.Done():
// 				t.Stop()
// 				return
// 			case <-t.C:
// 				r.sendMetric(a)
// 			}
// 		}
// 	}()
// 	return nil
// }

func (r *EloPerfData) Stop() {
	r.cancel()
}

func (r *EloPerfData) sendMetric(a telegraf.Accumulator) {
	n := rand.Int63n(r.Max-r.Min) + r.Min

	tags := map[string]string{
		"host":   "your_host",
		"metric": r.ValueName,
	}

	fields := map[string]interface{}{
		r.ValueName: n,
	}

	now := time.Now()

	a.AddFields("measurement_name", fields, tags, now)
}

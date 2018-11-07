package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/ssokssok/wmibeat/config"
)

// Wmibeat configuration.
type Wmibeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of wmibeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Wmibeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts wmibeat.
func (bt *Wmibeat) Run(b *beat.Beat) error {
	logp.Info("wmibeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": counter,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
		counter++
	}
}

// Stop stops wmibeat.
func (bt *Wmibeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

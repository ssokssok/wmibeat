package beater

import (
	"fmt"
	"time"
	"strings"
	"bytes"  

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

    var allValues common.MapStr

    for _, class := range bt.config.Classes {

      var query bytes.Buffer

      if len(class.Fields) > 0 {
        wmiFields := class.Fields 
        query.WriteString("SELECT ")
        query.WriteString(strings.Join(wmiFields, ","))
        query.WriteString(" FROM ")
        query.WriteString(class.Class)
        if class.WhereClause != "" {
          query.WriteString(" WHERE ")
          query.WriteString(class.WhereClause)
        }
      } else {
        query.WriteString("SELECT ")
        query.WriteString(" * ")
        query.WriteString("FROM ")
        query.WriteString(class.Class)
        if class.WhereClause != "" {
          query.WriteString(" WHERE ")
          query.WriteString(class.WhereClause)
        }
      }

      logp.Info("Query: "+query.String())
      eachValues, err := WmiQuery(query.String(), class.Fields)
      if err != nil {
        logp.Warn("WmiQuery error:",err)
        break
      }

      allValues = common.MapStrUnion(allValues, common.MapStr {class.Class: eachValues })
    }


		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"wmi":     allValues,
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

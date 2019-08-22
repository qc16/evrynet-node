package core

import (
	"log"
	"reflect"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
)

// ----------------------------------------------------------------------------

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		tendermint.RequestEvent{},
		tendermint.MessageEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
}

func (c *core) handleEvents() {
	// Clear state
	defer func() {
		c.handlerWg.Done()
	}()

	c.handlerWg.Add(1)

	for {
		select {
		case event, ok := <-c.events.Chan():
			if !ok {
				return
			}
			// A real event arrived, process interesting content
			switch ev := event.Data.(type) {
			case tendermint.RequestEvent:
				//TODO: Handle block proposal and remove this log
				log.Printf("--- Type of event.Data: %+v\n", reflect.TypeOf(ev))
				log.Printf("--- Value of event.Data: %+v\n", event.Data)
			case tendermint.MessageEvent:
				log.Printf("--- Type of event.Data: %+v\n", reflect.TypeOf(ev))
				log.Printf("--- Value of event.Data: %+v\n", ev.Payload)
				//TODO: Handle ev.Payload, if got error then call c.backend.Gossip()
			default:
				log.Printf("--- Unknow event :%v", ev)
			}
		}
	}
}

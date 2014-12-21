// Package pinger is an helper library for building a pinger service for
// availability monitoring
package pinger

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Pinger is an HTTP pinger service,
// it periodically updates the list of hosts to ping,
// every 5 minutes by default
type Pinger struct {
	httpClient  *http.Client
	hosts       []Host
	getter      Getter
	alertSender AlertSender
	mu          sync.RWMutex
}

// NewPinger returns a new Pinger object
func NewPinger(h *http.Client, g Getter, a AlertSender) *Pinger {
	hosts, err := g.Hosts()
	if err != nil {
		log.Fatal(err)
	}

	p := Pinger{httpClient: h, hosts: hosts, getter: g, alertSender: a}

	// start periodic hosts list updates
	go p.update(5 * time.Minute)

	return &p
}

// ping the hosts and return a map[Host[Response
func (p *Pinger) ping() map[Host]Response {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make(map[Host]Response)

	var wg sync.WaitGroup

	for _, h := range p.hosts {
		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Printf("pinging %s...", h.Name)

			status, body, err := h.Ping(p.httpClient)
			out[h] = Response{
				Error:      err,
				StatusCode: status,
				Body:       body,
			}

			if err != nil {
				log.Printf("ERROR %s: %s", h.Name, err.Error())
			} else {
				log.Printf("%s OK", h.Name)
			}
		}()
	}

	wg.Wait()

	return out
}

// update periodically updates the hosts list
func (p *Pinger) update(d time.Duration) {
	for {
		time.Sleep(d)

		func() {
			p.mu.Lock()
			defer p.mu.Unlock()

			hosts, err := p.getter.Hosts()
			if err != nil {
				log.Printf("error updating hosts: %s", err.Error())
			} else {
				p.hosts = hosts
			}
		}()
	}
}

// Ping pings the hosts list every time.Duration and sends alerts to the
// alerts sender.
//
// it is a blocking function.
func (p Pinger) Ping(d time.Duration) {
	for {
		resp := p.ping()
		for h, r := range resp {
			if r.Error != nil {
				p.alertSender.Send(fmt.Sprintf("host %s, error %s", h.Name,
					r.Error))
			}
		}

		time.Sleep(d)
	}
}

package pinger

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

type getter struct{}

func (g getter) Hosts() ([]Host, error) {
	return testHosts, nil
}

type alertSender struct{}

func (a alertSender) Send(msg string) error {
	return nil
}

var testHosts = []Host{
	Host{
		Name: "test1",
		Url:  "url1",
	},
}

func TestNewPinger(t *testing.T) {
	var g getter
	var a alertSender
	p := NewPinger(http.DefaultClient, g, a)

	if p.httpClient != http.DefaultClient {
		t.Fail()
	}

	if p.getter != g {
		t.Fail()
	}

	if p.alertSender != a {
		t.Fail()
	}

	if !reflect.DeepEqual(p.hosts, testHosts) {
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	var g getter
	p := Pinger{
		getter: g,
		hosts:  []Host{},
	}

	go p.update(time.Second)
	time.Sleep(2 * time.Second)

	if !reflect.DeepEqual(p.hosts, testHosts) {
		t.Fail()
	}
}

func ExamplePinger() {
	// getter and alertSender implement Hosts() and Send() respectively
	p := NewPinger(http.DefaultClient, getter{}, alertSender{})

	p.Ping(30 * time.Second)
}

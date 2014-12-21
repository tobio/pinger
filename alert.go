package pinger

// AlertSender is an interface for sending alerts for failed hosts
type AlertSender interface {
	Send(string) error
}

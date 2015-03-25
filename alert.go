package pinger

// AlertSender is an interface for sending alerts for health check results
type AlertSender interface {
	NotifyFailure(string,error) error
	NotifySuccess(string) error
}

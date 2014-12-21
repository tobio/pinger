package pinger

// Getter is an interface for gettings a list of hosts
type Getter interface {
	Hosts() ([]Host, error)
}

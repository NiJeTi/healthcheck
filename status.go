package healthcheck

// Status represents a health status state in the system.
type Status int

// StatusUnknown represents an undefined or uninitialized health status.
const StatusUnknown Status = -1

const (

	// StatusHealthy represents a healthy status state in the system.
	StatusHealthy = Status(iota)

	// StatusDegraded represents a degraded status state in the system, indicating reduced functionality.
	StatusDegraded

	// StatusUnhealthy represents an unhealthy status state in the system.
	StatusUnhealthy
)

// Int converts the Status value to its corresponding integer representation.
func (s Status) Int() int {
	return int(s)
}

// String converts the Status value to its corresponding string representation.
func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusDegraded:
		return "degraded"
	case StatusUnhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

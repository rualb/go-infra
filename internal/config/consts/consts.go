// Package consts app const
package consts

const AppName = "go-infra"

const (
	// DefaultTextSize default size of text field
	DefaultTextSize = 100

	// WF_STATUS_NEW       = 0
	// WF_STATUS_PROGRESS  = 6
	// WF_STATUS_DELETE    = 7
	// WF_STATUS_ERROR     = 10
	// WF_STATUS_SUCCESS   = 15
	// WF_STATUS_VOID      = 17
	// WF_STATUS_SIGNED    = 4
	// WF_STATUS_DELIVERED = 5
	// WF_STATUS_OUTBOX    = 3
	// WF_STATUS_READONLY  = 32
	// WF_STATUS_UNPAID    = 19
	// WF_STATUS_PAID      = 21
	// WF_STATUS_INQUEUE   = 31
)

// // log levels
// const (
// 	LogLevelError = 0
// 	LogLevelWarn  = 1
// 	LogLevelInfo  = 2
// 	LogLevelDebug = 3
// )

const (
	PathSysMetricsAPI = "/sys/api/metrics"

	PathTestPingAPI   = "/test/api/ping"
	PathTestHealthAPI = "/test/api/health"
	PathMessengerAPI  = "/messenger/api"
	PathConfigsAPI    = "/configs/api"
)

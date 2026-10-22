package netbox

// restyLogger adapts resty's internal Logger interface to the package-level
// logrus logger used by the NetBox client (see SetLogger/getLogger). Resty
// invokes it for its own request/response detail records; every call is
// forwarded to logrus so output honours LOG_LEVEL/LOG_FORMAT and the existing
// custom access-log line ("netbox request") is unaffected.
//
// resty v3.0.0-rc.4 exposes no trace level in its Logger interface (only
// Errorf/Warnf/Debugf), so the finest resty records arrive as Debugf and are
// gated by the logrus level at emit time: they surface under both
// LOG_LEVEL=debug and LOG_LEVEL=trace, and are suppressed otherwise.
type restyLogger struct{}

// Errorf forwards a resty error record to logrus.
func (restyLogger) Errorf(format string, v ...any) {
	getLogger().Errorf(format, v...)
}

// Warnf forwards a resty warning record to logrus.
func (restyLogger) Warnf(format string, v ...any) {
	getLogger().Warnf(format, v...)
}

// Debugf forwards a resty debug/request-response record to logrus. Resty's
// own request and response dumps arrive here; the Authorization header is
// redacted by resty's sanitizeHeaders before the dump is formatted, so the
// NetBox token is never written to the log.
func (restyLogger) Debugf(format string, v ...any) {
	getLogger().Debugf(format, v...)
}

//go:build ios && !production

package application

// iosVerboseLogging enables the framework's internal iOS diagnostics
// (request tracing, bridge messages, lifecycle markers). Debug builds only;
// production builds compile these call sites away via the constant.
const iosVerboseLogging = true

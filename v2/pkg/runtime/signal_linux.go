//go:build linux

package runtime

/*
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

static void fix_signal(int signum)
{
    struct sigaction st;

    if (sigaction(signum, NULL, &st) < 0) {
        return;
    }
    st.sa_flags |= SA_ONSTACK;
    sigaction(signum, &st, NULL);
}

static void fix_all_signals()
{
#if defined(SIGSEGV)
    fix_signal(SIGSEGV);
#endif
#if defined(SIGBUS)
    fix_signal(SIGBUS);
#endif
#if defined(SIGFPE)
    fix_signal(SIGFPE);
#endif
#if defined(SIGABRT)
    fix_signal(SIGABRT);
#endif
}
*/
import "C"

// ResetSignalHandlers resets signal handlers to allow panic recovery.
//
// On Linux, WebKit (used for the webview) may install signal handlers without
// the SA_ONSTACK flag, which prevents Go from properly recovering from panics
// caused by nil pointer dereferences or other memory access violations.
//
// Call this function immediately before code that might panic to ensure
// the signal handlers are properly configured for Go's panic recovery mechanism.
//
// Example usage:
//
//	go func() {
//	    defer func() {
//	        if err := recover(); err != nil {
//	            log.Printf("Recovered from panic: %v", err)
//	        }
//	    }()
//	    runtime.ResetSignalHandlers()
//	    // Code that might panic...
//	}()
//
// Note: This function only has an effect on Linux. On other platforms,
// it is a no-op.
func ResetSignalHandlers() {
	C.fix_all_signals()
}

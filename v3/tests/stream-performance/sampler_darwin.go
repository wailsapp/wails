//go:build darwin

package main

/*
#include <libproc.h>
#include <sys/proc_info.h>
#include <sys/resource.h>
#include <stdlib.h>
#include <string.h>

static int hp_listpids(int *buf, int count) {
	return proc_listallpids((void *)buf, count * (int)sizeof(int));
}

static int hp_name(int pid, char *out, int outlen) {
	struct proc_bsdinfo bsd;
	int r = proc_pidinfo(pid, PROC_PIDTBSDINFO, 0, &bsd, PROC_PIDTBSDINFO_SIZE);
	if (r != (int)PROC_PIDTBSDINFO_SIZE) return -1;
	strncpy(out, bsd.pbi_name, outlen - 1);
	out[outlen - 1] = '\0';
	return 0;
}

// ri_phys_footprint is what Activity Monitor shows and what jetsam acts on.
// The harness's own RSS is not the signal; this is.
static unsigned long long hp_footprint(int pid, int *ok) {
	rusage_info_current ri;
	if (proc_pid_rusage(pid, RUSAGE_INFO_CURRENT, (rusage_info_t *)&ri) != 0) {
		*ok = 0;
		return 0ULL;
	}
	*ok = 1;
	return (unsigned long long)ri.ri_phys_footprint;
}
*/
import "C"

import "fmt"

const samplerSupported = true

const webContentProcName = "com.apple.WebKit.WebContent"

// webContentPids enumerates every WebContent process on the machine. Callers
// must diff against a pre-launch snapshot — summing all of them would fold in
// other applications' Safari tabs.
//
// proc_listallpids returns the number of PIDS written, not a byte count.
// (Verified: 907 returned against 907 pids present.) Dividing by sizeof(int)
// here would scan only a quarter of the table and could miss our own process.
func webContentPids() ([]int, error) {
	const maxPids = 32768
	buf := make([]C.int, maxPids)
	n := int(C.hp_listpids(&buf[0], C.int(maxPids)))
	if n <= 0 {
		return nil, fmt.Errorf("proc_listallpids failed (returned %d)", n)
	}
	if n > maxPids {
		n = maxPids
	}

	name := make([]C.char, 64)
	var out []int
	for i := 0; i < n; i++ {
		pid := int(buf[i])
		if pid <= 0 {
			continue
		}
		if C.hp_name(C.int(pid), &name[0], C.int(len(name))) != 0 {
			continue
		}
		if C.GoString(&name[0]) == webContentProcName {
			out = append(out, pid)
		}
	}
	return out, nil
}

// footprint returns ri_phys_footprint for pid. ok=false means the process is
// gone or unreadable — for a tracked WebContent pid that means it was replaced.
func footprint(pid int) (uint64, bool) {
	var ok C.int
	v := C.hp_footprint(C.int(pid), &ok)
	if ok == 0 {
		return 0, false
	}
	return uint64(v), true
}

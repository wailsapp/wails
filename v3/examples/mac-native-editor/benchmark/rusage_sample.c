#include <errno.h>
#include <libproc.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <unistd.h>

static int snapshot(pid_t pid, struct rusage_info_v4 *usage) {
    memset(usage, 0, sizeof(*usage));
    return proc_pid_rusage(pid, RUSAGE_INFO_V4, (rusage_info_t *)usage);
}

int main(int argc, char **argv) {
    if (argc != 3) {
        fprintf(stderr, "usage: %s pid seconds\n", argv[0]);
        return 2;
    }
    pid_t pid = (pid_t)strtol(argv[1], NULL, 10);
    unsigned seconds = (unsigned)strtoul(argv[2], NULL, 10);
    struct rusage_info_v4 before;
    struct rusage_info_v4 after;
    if (snapshot(pid, &before) != 0) {
        fprintf(stderr, "initial proc_pid_rusage: %s\n", strerror(errno));
        return 1;
    }
    sleep(seconds);
    if (snapshot(pid, &after) != 0) {
        fprintf(stderr, "final proc_pid_rusage: %s\n", strerror(errno));
        return 1;
    }
    printf(
        "{\"pid\":%d,\"seconds\":%u,"
        "\"resident_bytes_start\":%llu,\"resident_bytes_end\":%llu,"
        "\"footprint_bytes_start\":%llu,\"footprint_bytes_end\":%llu,"
        "\"lifetime_max_footprint_bytes\":%llu,"
        "\"user_time_ns_delta\":%llu,\"system_time_ns_delta\":%llu,"
        "\"idle_wakeups_delta\":%llu,\"interrupt_wakeups_delta\":%llu,"
        "\"pageins_delta\":%llu,\"instructions_delta\":%llu,"
        "\"cycles_delta\":%llu}\n",
        pid,
        seconds,
        before.ri_resident_size,
        after.ri_resident_size,
        before.ri_phys_footprint,
        after.ri_phys_footprint,
        after.ri_lifetime_max_phys_footprint,
        after.ri_user_time - before.ri_user_time,
        after.ri_system_time - before.ri_system_time,
        after.ri_pkg_idle_wkups - before.ri_pkg_idle_wkups,
        after.ri_interrupt_wkups - before.ri_interrupt_wkups,
        after.ri_pageins - before.ri_pageins,
        after.ri_instructions - before.ri_instructions,
        after.ri_cycles - before.ri_cycles
    );
    return 0;
}

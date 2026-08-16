# Dev Session invalidation and process lifecycle

Type: grilling
Status: open
Blocked by: 03

## Question

How does the long-lived Dev Session map file changes to affected Nodes,
coordinate frontend HMR with backend rebuilds, coalesce bursts, cancel stale
Plans, restart or preserve the application process, report readiness and
failures, and recover when watched files or configuration change? Define the
fast path and lifecycle contract without turning persistent processes into DAG
Nodes.

## Comments

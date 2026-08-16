# Build performance acceptance matrix and regression budget

Type: grilling
Status: resolved
Blocked by: 03, 04, 05, 08

## Question

Using the measured current baseline and the resolved cache, generated-asset,
and hook contracts, which workloads, project sizes, Targets, package formats,
and development changes form the permanent benchmark matrix? Define the
wall-clock targets, variance treatment, regression thresholds, and CI or
release gates that make speed a verifiable product requirement.

## Comments

### 2026-08-16 — Accepted implementation default

- Correct no-op desktop builds target under 100ms on the badge fixture and may
  not regress more than 20% from the checked-in benchmark median without an
  explicit performance note.
- Incremental builds must execute only the affected semantic Nodes. Ordinary
  Go edits must not rebuild the frontend when binding bytes are unchanged;
  frontend-only edits must not rerun binding analysis.
- Cold-build orchestration overhead must remain below 5% of the critical path.
  Package and multi-Target work should run independent branches concurrently,
  bounded by declared CPU, memory, and exclusive-tool claims.
- CI runs correctness and deterministic-plan tests on every change; noisy wall
  benchmarks run in a controlled release job and compare medians, not single
  samples.

### 2026-08-16 — MVP evidence

The integrated manifest path currently satisfies the primary threshold:
70–80ms warm badge builds and 70–90ms warm DEB packaging. Semantic Go
snapshots keep method implementation edits off the binding/frontend path while
retaining free-function bodies that may delegate service registration.

## Answer

The permanent acceptance contract is recorded in
[`performance-acceptance.md`](../performance-acceptance.md). Controlled release
runs use two warmups and seven samples, compare medians, reject evidence above
15% MAD, enforce the 100ms badge no-op and 20% regression budgets, verify work
counts in every sample, and cap cold orchestration overhead at 5%.

Ordinary CI owns deterministic target/format planning, semantic invalidation,
package and multi-Target overlap, Dev lifecycle, race, and vet gates. Matching
native release hosts own Linux, Windows, macOS, iOS, Android, cross-architecture,
package, and credentialed signing evidence. Signing and iOS package actions are
deliberately never accepted as reusable cache hits.

The benchmark harness now enforces the complete budget. The native verifier
accepts host Target matrices and mobile/signing modes and requires zero-work
reruns for cacheable commands. Final Linux evidence passes: 93.503ms badge
median with 8.19% MAD and 0/4 work, 3.22% cold overhead, every Linux format,
and a 0/9 warm rerun after AppImage packaging.

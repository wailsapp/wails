# Taskfile AST migration and unsupported customization policy

Type: grilling
Status: open
Blocked by: 01

## Question

What generated Taskfile patterns and user modifications can be translated into
manifest fields, profiles, hooks, or template references? What must be reported
for manual migration? Define dry-run output, machine-readable diagnostics,
write/remove behavior, behavior when `wails.toml` already exists, and the
preconditions for removing legacy files when unsupported Taskfile logic
remains. Define a machine-readable migration completion state, any explicit
cutover marker, and provenance-aware removal that cannot delete modified,
user-owned, or incompletely represented files. Define the safety rule that
inline shell blocks are never silently converted into generated scripts.

## Comments

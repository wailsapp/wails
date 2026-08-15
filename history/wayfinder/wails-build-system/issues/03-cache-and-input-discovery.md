# Automatic input discovery and cache identity

Type: grilling
Status: open
Blocked by: 02

## Question

How does the new runner discover and hash inputs for Go packages, embedded
files, frontend lockfiles and outputs, generated bindings, platform metadata,
toolchains, environment variables, and user-owned scripts? Which operations are
cacheable by default, and which must always run because they may have external
side effects?

## Comments

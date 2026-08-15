# Script hook contract and safe extension boundary

Type: grilling
Status: open
Blocked by: 02

## Question

What phases may invoke user-owned scripts, what arguments and environment are
stable, how are target/profile/output values passed, how are failures reported,
and what input/output declaration is required for a hook to participate in the
cache? The cache contract must treat the referenced script content and relevant
executable metadata as implicit inputs. Confirm that inline shell remains
unsupported in the new manifest.

## Comments

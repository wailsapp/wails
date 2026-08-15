# Issue tracker: local Markdown

Planning issues for this repository are local Markdown files under
`history/wayfinder/`. All tracker operations read and write these files; GitHub
issues and pull requests are outside this workflow.

## General conventions

- One effort per directory: `history/wayfinder/<effort>/`.
- The canonical map is `history/wayfinder/<effort>/map.md`.
- Child tickets are individual files under
  `history/wayfinder/<effort>/issues/`, numbered from `01`.
- Conversation history is appended under `## Comments`.
- Resolutions are appended under `## Answer`.

## Ticket header

Each child ticket starts with:

```markdown
# <ticket name>

Type: research|prototype|grilling|task
Status: open|claimed|resolved
Blocked by: none|NN, NN
```

The filename number is the ticket identity. Human-facing references use the
ticket name as a Markdown link rather than a bare number.

## Wayfinding operations

- **Create:** add `issues/NN-<slug>.md` with a unique next number and an
  `## Question` body.
- **Wire:** after ticket files exist, record blocking ticket numbers on the
  `Blocked by:` line.
- **Frontier:** scan numbered tickets in order. A frontier ticket has
  `Status: open` and every ticket in `Blocked by:` has `Status: resolved`.
- **Claim:** set `Status: claimed` and save before reading deeply or doing work.
- **Resolve:** append the decision under `## Answer`, set `Status: resolved`,
  then add a one-line gist and relative link under the map's
  `## Decisions so far`.
- **Comment:** append dated discussion under `## Comments`; keep the ticket's
  question unchanged.
- **Out of scope:** resolve the ticket as out of scope and link its gist under
  the map's `## Out of scope`, not `## Decisions so far`.

Open tickets are discovered from `issues/`; the map does not duplicate them.

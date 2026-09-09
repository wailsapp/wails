# Wails star history

This is the dependency-free Go generator for `website/static/img/star-history.svg`.
It calls GitHub's timestamped stargazers endpoint, groups the visible stars into
weekly points, and renders a self-contained SVG using the Digital Wales artwork
as a faded background with the horizontal Wails logo embedded in the header.

GitHub now requires the token to belong to a repository admin or collaborator for
the stargazers endpoint. The workflow passes the existing `WAILS_REPO_TOKEN`
when available, with `STAR_HISTORY_TOKEN` and `GITHUB_TOKEN` as fallbacks.

Run it locally from this directory:

```sh
GOWORK=off GITHUB_TOKEN="$(gh auth token)" go run .
```

The scheduled workflow renders the SVG weekly and retains it as a workflow
artifact. A manually dispatched run also publishes the changed SVG, keeping
repository writes behind an explicit confirmation.

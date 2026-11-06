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

The scheduled workflow renders and publishes the SVG weekly, retaining it as a
workflow artifact and committing the updated image to the repository. A manually
dispatched run performs the same refresh immediately. Publishing runs only in
wailsapp/wails and uses the existing repository token; its identity must be
allowed to push to the target branch under the repository rules.

#!/usr/bin/env bash
#
# GitHub's Ubuntu runner images preinstall Microsoft's apt repositories (azure-cli and the
# packages.microsoft.com prod repo). Wails uses neither, and when that host serves a 403 for
# them - see https://github.com/wailsapp/wails/issues/6040 - `apt-get update` fails and takes
# the whole job down with it. Dropping the two list files makes our CI independent of it.
#
# Only do this on ephemeral GitHub-hosted VMs. On a self-hosted runner, or on a developer's
# machine running these workflows locally with `act`, /etc/apt belongs to the host: deleting
# real apt sources there is destructive and not something a CI job should do.
set -euo pipefail

if [ "${ACT:-}" = "true" ] || [ "${RUNNER_ENVIRONMENT:-}" != "github-hosted" ]; then
	echo "Leaving apt sources untouched (ACT=${ACT:-unset}, RUNNER_ENVIRONMENT=${RUNNER_ENVIRONMENT:-unset})."
	exit 0
fi

sudo rm -f /etc/apt/sources.list.d/microsoft-prod.list /etc/apt/sources.list.d/azure-cli.list

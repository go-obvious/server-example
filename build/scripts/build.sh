#!/usr/bin/env bash
set -euo pipefail

project_root=$(git rev-parse --show-toplevel)
apps_dir="cmd"
dist_dir="dist"

die () {
    echo >&2 "$@"
    exit 1
}

service_name=$1
[ -z "${service_name:-}" ] && { die "\"service_name\" parameter not provided"; }

mkdir -p ${project_root:-}/${dist_dir:-}/

function get_ld_flags() {
    BUILD_TIME="$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
    TAG="current"
    REVISION="current"
    if hash git 2>/dev/null && [ -e ${project_root}/.git ]; then
        TAG="$(git describe --tags 2>/dev/null || true)"
        [[ -z "$TAG" ]] && TAG="notag"
        REVISION="$(git rev-parse HEAD)"
    fi
    echo "-s -w \
        -X github.com/go-obvious/server-example/internal/build.Time=${BUILD_TIME} \
        -X github.com/go-obvious/server-example/internal/build.Rev=${REVISION} \
        -X github.com/go-obvious/server-example/internal/build.Tag=${TAG}"
}
LD_FLAGS=$(get_ld_flags)

(
    pushd ${project_root}/${apps_dir}/${service_name}
    for GOARCH in amd64 arm64; do
            CGO_ENABLED=0 GOARCH=${GOARCH} go build -mod=mod -trimpath -ldflags="${LD_FLAGS}" -tags 'netgo osusergo' -o "${project_root}/${dist_dir}/${service_name}-${GOARCH}"
    done
)
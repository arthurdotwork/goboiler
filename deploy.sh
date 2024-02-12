#!/bin/sh

set -e

if [ "$#" -lt 1 ] || [ -z "$GITHUB_SHA" ] || [ -z "$GITHUB_REF" ]; then
    echo 'usage: env GITHUB_SHA=<commit-hash> GITHUB_REF=<git-ref> deploy.sh <image> [...docker build args]' >&2
    exit 1
fi

img="$1"

shift

branch="$(echo "$GITHUB_REF" | sed -n 's#refs/heads/##p')"

if [ "$branch" = "master" ] ; then
    tag="latest"
else
    tag="$(echo "$GITHUB_REF" | sed -n 's#refs/tags/##p')"
fi

if [ -n "$tag" ]; then
    docker build -t "${img}:${GITHUB_SHA}" "$@"
    docker tag "${img}:${GITHUB_SHA}" "${img}:${tag}"
    docker push "${img}:${tag}"
else
    echo  'no action taken'
fi

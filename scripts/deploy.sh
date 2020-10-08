#!/usr/bin/env bash
set -euox pipefail
SCRIPT_DIR=$(dirname "$0")
echo "${SCRIPT_DIR}"

project=$(gcloud secrets versions access latest --secret="project-id")
if [[ -z "${project}" ]]; then
  echo -n "need project"
  exit 1
fi
echo "${project}"

kubectl run go-subscriber-fs \
  --image gcr.io/"${project}"/go-subscriber-fs:latest \
  --env="PUB_PROJECT=${project}"

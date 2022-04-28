#!/bin/bash

UUID=$(jq -r ".last_run_uuid" manifest.json)
ARTIFACT=$(jq -r ".builds[] | \
                    select(.packer_run_uuid == \"${UUID}\") |\
                    .artifact_id" manifest.json)

# NOTE(jkoelker) extract just the timestamp as the chars after the final `-`
TS="${ARTIFACT##*-}"

# NOTE(jkoelker) extract the image name as the chars after the final `/`
IMG="${ARTIFACT##*/}"

# NOTE(jkoelker) further refine the image name by stripping off the timestamp
IMG="${IMG%-"$TS"}"

az login \
    --service-principal \
    --username "${AZURE_CLIENT_ID}" \
    --password "${AZURE_CLIENT_SECRET}" \
    --tenant "${AZURE_TENANT_ID}"

az sig image-version delete \
    --resource-group dkp \
    --gallery-name dkp \
    --gallery-image-definition "${IMG}" \
    --gallery-image-version "0.0.${TS}"

#!/bin/bash

LAST_RUN_UUID=$(jq -r ".last_run_uuid" manifest.json)
FULL_ID=$(jq -r ".builds[] | select(.packer_run_uuid == \"$LAST_RUN_UUID\") | .artifact_id" manifest.json)
REGION=$(echo "$FULL_ID" | cut -d':' -f1)
ID=$(echo "$FULL_ID" | cut -d':' -f2)

IMAGE=$(aws ec2 describe-images --region "$REGION" --image-ids "$ID")
NAME=$(echo "$IMAGE" | jq -r ".Images[0].Name")
SNAPSHOT=$(aws ec2 describe-snapshots --region "$REGION" --filters "Name=tag:ami_name,Values=$NAME")
SNAPSNOT_IDS=$(echo "$SNAPSHOT" | jq -r ".Snapshots[] | .SnapshotId")

echo "Deregistering Image:"
echo "$IMAGE"
aws ec2 deregister-image --region "$REGION" --image-id "$ID"

echo "Deleting Snapshots:"
echo "$SNAPSHOT"

for SNAPSNOT_ID in $SNAPSNOT_IDS
do
  aws ec2 delete-snapshot --region "$REGION" --snapshot-id "$SNAPSNOT_ID"
done

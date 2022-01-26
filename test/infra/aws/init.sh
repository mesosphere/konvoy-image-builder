set -o errexit
set -o pipefail
set -o nounset

OWNER="${OWNER:-$(whoami)}"
EXPIRATION="${EXPIRATION:-4h}"
CLUSTER_NAME="${CLUSTER_NAME:-$OWNER}"

cat >tmp <<EOF
{
    "tags": {
        "owner": "$OWNER",
        "expiration": "$EXPIRATION"
    }
}
EOF

# Do not write an empty, or invalid file
mv tmp terraform.tfvars.json

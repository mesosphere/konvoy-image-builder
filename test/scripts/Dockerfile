FROM debian:buster-slim

# NOTE(schung) Ignore "Pin versions in apt get install"
# hadolint ignore=DL3008
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        awscli jq \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

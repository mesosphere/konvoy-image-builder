# syntax=docker/dockerfile:1.4

ARG BASE=mesosphere/konvoy-image-builder:latest-devkit
# NOTE(jkoelker) Ignore "Always tag the version of an image explicitly"
# hadolint ignore=DL3006
FROM ${BASE} as devkit

ARG TARGETPLATFORM
FROM --platform=${TARGETPLATFORM} alpine:3.15.2

ENV ANSIBLE_PATH=/usr
ENV PYTHON_PATH=/usr

COPY requirements.txt /tmp/
# NOTE(jkoelker) Ignore "Pin versions in [pip | apk add]"
# hadolint ignore=DL3013,DL3018
RUN apk add --no-cache \
        # NOTE(jkoelker) for packer-provisioner-goss / konvoy-image
        libc6-compat \
        openssh-client \
        python3 \
        py3-cryptography \
        py3-pip \
        py3-wheel \
    && pip3 install --no-cache-dir --requirement /tmp/requirements.txt \
    && rm -rf /root/.cache

ARG BUILDARCH
COPY --from=devkit /usr/local/bin/goss-${BUILDARCH} /usr/local/bin/goss
COPY --from=devkit /usr/local/bin/packer-${BUILDARCH} /usr/local/bin/packer
COPY --from=devkit /usr/local/bin/packer-provisioner-goss-${BUILDARCH} /usr/local/bin/packer-provisioner-goss
COPY konvoy-image /usr/local/bin
COPY images /root/images
COPY ansible /root/ansible
COPY packer /root/packer
WORKDIR /root
ENTRYPOINT ["/usr/local/bin/konvoy-image"]

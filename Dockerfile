#
ARG BASE=mesosphere/konvoy-image-builder:latest-devkit
# NOTE(jkoelker) Ignore "Always tag the version of an image explicitly"
# hadolint ignore=DL3006
FROM ${BASE} as devkit

FROM alpine:3.15.4

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

COPY --from=devkit /usr/local/bin/goss /usr/local/bin/
COPY --from=devkit /usr/local/bin/packer /usr/local/bin/
COPY --from=devkit /usr/local/bin/packer-provisioner-goss /usr/local/bin/
COPY --from=devkit /root/.config/packer/plugins/ /root/.config/packer/plugins/
COPY bin/konvoy-image /usr/local/bin
COPY images /root/images
COPY ansible /root/ansible
COPY packer /root/packer

# this is needed for containerd tar
# place it in /usr/share/ansible/collections, the container will be run with a different user
RUN mkdir -p /usr/share/ansible/collections && ansible-galaxy collection install ansible.utils -p /usr/share/ansible/collections

WORKDIR /root
ENTRYPOINT ["/usr/local/bin/konvoy-image"]

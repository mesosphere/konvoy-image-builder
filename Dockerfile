ARG BUILDARCH
ARG BASE=mesosphere/konvoy-image-builder:latest-devkit-${BUILDARCH}
# NOTE(jkoelker) Ignore "Always tag the version of an image explicitly"
# hadolint ignore=DL3006
FROM ${BASE} as devkit

ARG TARGETPLATFORM
# hadolint ignore=DL3029
FROM --platform=${TARGETPLATFORM} alpine:3.15.4

ENV ANSIBLE_PATH=/usr
ENV PYTHON_PATH=/usr
ENV PACKER_PLUGIN_PATH=/opt/packer/plugins

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
        xorriso \
    && pip3 install --no-cache-dir --requirement /tmp/requirements.txt \
    && rm -rf /root/.cache

ARG BUILDARCH
# we copy this to remote hosts to execute GOSS
COPY --from=devkit /usr/local/bin/goss-amd64 /usr/local/bin/goss-amd64
COPY --from=devkit /usr/local/bin/goss-${BUILDARCH} /usr/local/bin/goss
# we copy this to remote hosts to execute mindthegap so its always amd64
COPY --from=devkit /usr/local/bin/mindthegap /usr/local/bin/
COPY --from=devkit /usr/local/bin/packer-${BUILDARCH} /usr/local/bin/packer
COPY --from=devkit /usr/local/bin/packer-provisioner-goss-${BUILDARCH} /usr/local/bin/packer-provisioner-goss
COPY --from=devkit /usr/local/bin/govc /usr/local/bin/
COPY --from=devkit /root/.config/packer/plugins/ ${PACKER_PLUGIN_PATH}
COPY --from=devkit /usr/share/ansible/collections/ansible_collections/ /usr/share/ansible/collections/ansible_collections/
COPY bin/konvoy-image-${BUILDARCH} /usr/local/bin/konvoy-image
COPY images /root/images
COPY ansible /root/ansible
COPY packer /root/packer

WORKDIR /root
ENTRYPOINT ["/usr/local/bin/konvoy-image"]

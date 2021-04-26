FROM python:3-slim-buster

ENTRYPOINT ["/entrypoint.sh"]
CMD ["sleep", "infinity"]

EXPOSE 22
# NOTE(jkoelker) Ignore "Pin versions in apt get install"
# hadolint ignore=DL3008
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        openssh-server \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd sshgroup \
    && useradd -ms /bin/bash -g sshgroup sshuser \
    && mkdir /setup

COPY entrypoint.sh /entrypoint.sh
COPY sshd_config /etc/ssh/sshd_config

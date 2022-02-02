FROM registry.access.redhat.com/ubi8/ubi:8.2

COPY *.repo /etc/yum.repos.d/

# hadolint ignore=DL3033
RUN yum -y install yum-utils epel-release modulemd-tools && yum clean all
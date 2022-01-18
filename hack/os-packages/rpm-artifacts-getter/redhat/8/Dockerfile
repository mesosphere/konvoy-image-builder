FROM registry.access.redhat.com/ubi8/ubi:8.2

COPY *.repo /etc/yum.repos.d/

# hadolint ignore=DL3033
RUN yum -y install yum-utils https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm modulemd-tools && yum clean all
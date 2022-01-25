FROM registry.access.redhat.com/ubi7/ubi:7.9

COPY *.repo /etc/yum.repos.d/

# hadolint ignore=DL3033
RUN yum -y install yum-utils https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm createrepo && yum clean all
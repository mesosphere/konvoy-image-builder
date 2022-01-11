ARG BASE_IMAGE=mesosphere/centos:7.9.2009.minimal
# hadolint ignore=DL3006
FROM $BASE_IMAGE

COPY *.repo /etc/yum.repos.d/

# hadolint ignore=DL3033
RUN yum -y install epel-release yum-utils createrepo && \
    yum clean all

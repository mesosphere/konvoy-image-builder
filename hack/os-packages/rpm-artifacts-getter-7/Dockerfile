FROM mesosphere/centos:7.9.2009.minimal

COPY *.repo /etc/yum.repos.d/

# hadolint ignore=DL3033
RUN yum install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm && \
    yum -y install yum-utils createrepo && \
    yum clean all

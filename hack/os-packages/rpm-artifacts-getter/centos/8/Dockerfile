FROM centos:8

COPY *.repo /etc/yum.repos.d/

# hadolint ignore=DL3033
RUN yum -y install epel-release yum-utils createrepo && \
    yum clean all

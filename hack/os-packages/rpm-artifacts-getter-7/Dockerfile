FROM centos:7.9.2009

COPY *.repo /etc/yum.repos.d/

RUN yum install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm && \
    yum -y install yum-utils createrepo && \
    yum clean all
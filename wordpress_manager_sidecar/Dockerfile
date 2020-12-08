FROM centos:7

LABEL Developer="Lakshya Nanjangud <lakshya.nanjangud@in.ibm.com>"

LABEL Maintainer="Openshift Team - Resiliency Orchestration"

WORKDIR /usr/src/app

COPY . .

RUN mkdir dbfile

CMD "./wpbin"

# RUN With This VOLUME: docker run -v wpvol:/dbfile/ wpr
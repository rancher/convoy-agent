FROM ubuntu:15.04
MAINTAINER Rancher Labs, Inc.

RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	curl

RUN curl --insecure -sSL -o convoy.tar.gz https://github.com/rancher/convoy/releases/download/v0.3.1-gfs-rc1/convoy.tar.gz; tar -xvzf convoy.tar.gz; mv convoy/convoy /usr/bin/convoy

COPY convoy-agent /usr/bin/
COPY convoy-gluster-launcher /
ENTRYPOINT ["/convoy-gluster-launcher"]

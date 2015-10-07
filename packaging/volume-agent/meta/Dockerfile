FROM ubuntu:15.04
MAINTAINER Rancher Labs, Inc.

RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	curl \
	nfs-common \
	uuid-runtime

RUN curl --insecure -sSL https://get.docker.com | sh
COPY meta-launcher /
ENTRYPOINT ["/meta-launcher"]

FROM ubuntu:15.04
MAINTAINER Rancher Labs, Inc.
RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	curl
RUN curl --insecure -sSL -o convoy.tar.gz https://github.com/rancher/convoy/releases/download/v0.3/convoy.tar.gz; tar -xvzf convoy.tar.gz; mv convoy/convoy /usr/bin/convoy
#RUN curl --insecure -sSL -o convoy-agent.tar.gz CONVOY_AGENT_DOWNLOAD_URL; tar -xvzf convoy-agent.tar.gz; mv convoy-agent /usr/bin/convoy-agent
COPY convoy-agent /usr/bin/
COPY launch /
ENTRYPOINT ["/launch"]

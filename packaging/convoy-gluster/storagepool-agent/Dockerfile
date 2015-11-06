FROM ubuntu:15.04
MAINTAINER Rancher Labs, Inc.
RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	curl \
	uuid-runtime
#RUN curl --insecure -sSL -o convoy-agent.tar.gz CONVOY_SP_AGENT_DOWNLOAD_URL; tar -xvzf convoy-agent.tar.gz; mv convoy-agent /usr/bin/convoy-agent
COPY convoy-agent /usr/bin/
COPY launch /
ENTRYPOINT ["/launch"]


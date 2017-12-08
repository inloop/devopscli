FROM alpine:3.5

ADD https://get.docker.com/builds/Linux/x86_64/docker-1.10.3 /usr/local/bin/docker

RUN apk update && \
    apk --update add ruby ruby-json ruby-bigdecimal ruby-io-console \
    ca-certificates libssl1.0 openssl libstdc++ && \
    chmod +x /usr/local/bin/docker

RUN apk --update add --virtual build-dependencies ruby-dev build-base openssl-dev && \
    gem install kontena-cli --no-rdoc --no-ri -v ${CLI_VERSION} && \
    apk del build-dependencies

COPY bin/devops-alpine /devops

# https://serverfault.com/questions/772227/chmod-not-working-correctly-in-docker
RUN mv /devops /usr/local/bin/devops && chmod +x /usr/local/bin/devops

ENTRYPOINT []

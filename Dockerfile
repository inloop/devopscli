FROM alpine:3.5

COPY bin/devops-alpine /devops

# https://serverfault.com/questions/772227/chmod-not-working-correctly-in-docker
RUN apk --update add docker && \
    mv /devops /usr/local/bin/devops && \
    chmod +x /usr/local/bin/devops

ENTRYPOINT []

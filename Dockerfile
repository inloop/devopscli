FROM alpine

COPY bin/devops-alpine /usr/bin/devops
RUN chmod +x /usr/bin/devops

ENTRYPOINT []

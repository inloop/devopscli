FROM alpine

COPY bin/devops-alpine /devops

# https://serverfault.com/questions/772227/chmod-not-working-correctly-in-docker
RUN mv /devops /usr/bin/devops && chmod +x /usr/bin/devops 

ENTRYPOINT []

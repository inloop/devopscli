FROM kontena/cli

COPY bin/devops-alpine /usr/local/bin/devops
RUN chmod +x /usr/local/bin/devops

ENTRYPOINT []

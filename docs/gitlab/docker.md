# Gitlab CI with docker

## Linked dependencies during build

For example:

* `Node.js` - submodules has to be installed during build process, because some dependencies contains create binaries for current OS

```
# Dockerfile
FROM node:8.9.0-alpine

COPY . /code
WORKDIR /code

RUN rm -rf node_modules && yarn install --production --ignore-engines
```

```
# .gitlab-ci.yml
stages:
  - build

build:
  image: inloopeu/devops
  stage: build
  tags:
    - docker
  services:
    - docker:dind
  script:
    - devops gitlab docker build
```

## Image tagging

Tag name is deducted from `CI_COMMIT_REF_NAME` with following rules:

* master -> latest
* develop -> unstable
* release/x -> x
* other -> other

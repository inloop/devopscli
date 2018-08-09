# Gitlab CI with docker (Create React App)

CRA apps can be deployed using docker easily. If you don't need the SSR, the image is very small and nginx footprint is very low.

Basically all You need is to create `Dockerfile` for building you app:

```
FROM inloopx/cra-docker

COPY build /app
```

Then the pipeline configuration in `.gitlab-ci.yml` should look like this to build your image and push it to gitlab registry:

```
stages:
  - compile
  - build

compile:
  image: node:8.4.0
  stage: compile
  tags:
    - docker
  artifacts:
    expire_in: 10 minutes
    paths:
      - build/
  script:
    # you can also use `npm install` and `npm run build`
    - yarn install
    - yarn build

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

_NOTE: image inloopx/cra-docker is used, you can find more info at [https://github.com/inloop/cra-docker](https://github.com/inloop/cra-docker)_

_NOTE2: images are tagged automatically, read more in [Docker image tagging](../docker-image-tagging.md)_

_NOTE3: Image is built using https://github.com/inloop/devopscli which is our tool to ease the deployment tasks_

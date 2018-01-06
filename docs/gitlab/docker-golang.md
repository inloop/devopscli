# Gitlab CI with docker (Golang)

## Create cross compiled binaries

```
# .gitlab-ci.yml
stages:
  - test
  - compile

test:
  image: inloopeu/devops:latest-golang
  tags:
    - docker
  stage: test
  script:
    - devops gitlab go test

compile:
  tags:
    - docker
  stage: compile
  image: inloopeu/devops:latest-golang
  only:
    - develop
    - master
    - tags
  script:
    - devops gitlab go compile
```

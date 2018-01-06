# Gitlab CI with docker (Golang)

## Cross compile and upload to S3

```
# .gitlab-ci.yml
stages:
  - test
  - compile
  - deploy

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
  artifacts:
    expire_in: 10 minutes
    paths:
    - bin/
  script:
    - devops gitlab go compile


deploy:
  image: xueshanf/awscli
  tags:
    - docker
  stage: deploy
  only:
    - master
    - develop
  artifacts:
    expire_in: 10 minutes
    paths:
    - bin/
  script:
    - aws s3 sync bin s3://static.inloop.eu/project-setup/$CI_COMMIT_REF_NAME
```

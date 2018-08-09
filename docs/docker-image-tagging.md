## Docker Image tagging

Image tag is deducted from `CI_COMMIT_REF_NAME` with following rules:

- `master` -> `latest`
- `develop` -> `unstable`
- `release/x` -> `x`
- `other` -> `other` (eg. feature branches, tags)

_NOTE: this naming is opinionated and if it doesn't suit your needs, you can always use `-t,tag` attribute_

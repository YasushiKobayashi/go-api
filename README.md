[![CircleCI](https://circleci.com/gh/YasushiKobayashi/go-api.svg?style=svg)](https://circleci.com/gh/YasushiKobayashi/go-api)
[![Build Status](https://travis-ci.org/YasushiKobayashi/go-api.svg?branch=master)](https://travis-ci.org/YasushiKobayashi/go-api)

### front repo
https://github.com/YasushiKobayashi/react-cms

### env
- golang version 1.8
- mysql5.6

### set up golang
```bash
# mac
brew install go
brew install glide
brew install direnv
echo 'eval "$(direnv hook bash)"' >> ~/.bash_profile
```

### setup project
```bash
echo 'export DOCUMENT_ENV=develop' >> ~/.bash_profile
cd src/document
glide install
```

### doucment
`godoc -http=:6060`

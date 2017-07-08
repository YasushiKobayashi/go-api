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
cd src/document
glide install
```


### migrate
```bash
cd src/document
goose up
```

### doucment
`godoc -http=:6060`

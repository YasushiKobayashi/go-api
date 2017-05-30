
### front repo
https://github.com/YasushiKobayashi/react-cms

### env
- golang version 1.8
- mysql5.6

### set up golang
```shell
# mac
brew install go
brew install glide
brew install direnv

# centos
yum -y install golang
```

### setup project
```shell
cd src/document
glide install
```


### migrate
```shell
cd src/document
goose up
```

FROM golang:1.11

LABEL version="1.0.0"
LABEL maintainer="Yongjin Jung (eniv)"
LABEL repository="https://github.com/ainvyu/radibrary-cli.go"
LABEL homepage="https://github.com/ainvyu/radibrary-cli.go"
LABEL "com.github.actions.name"="Go"
LABEL "com.github.actions.description"="Golang"
LABEL "com.github.actions.icon"="code"
LABEL "com.github.actions.color"="white"

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

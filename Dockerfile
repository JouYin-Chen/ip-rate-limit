# 改用 alpine 當基底，縮減 docker image 的大小
FROM golang:alpine

# 加入 git
RUN apk add --no-cache git

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

ENV APP_DIR /app
RUN mkdir -p $APP_DIR

ADD . $APP_DIR
ENTRYPOINT (cd $APP_DIR && ./ip-rate-limit)

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 go build -ldflags '-d -w -s'

EXPOSE 3000
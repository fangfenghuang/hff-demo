# 基于 Debian 镜像
由于 Debian 镜像中已经包含了tzdata，因此设置时区的方法比较简单，只需添加环境变量TZ即可。
```
FROM debian:stretch

ENV TZ=Asia/Shanghai
```

# 基于 Alpine 镜像
```
FROM alpine:3.9

ENV TZ=Asia/Shanghai

RUN apk update \
    && apk add tzdata \
    && echo "${TZ}" > /etc/timezone \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && rm /var/cache/apk/*
```

# 基于 Ubuntu 镜像
```
FROM ubuntu:bionic

ENV TZ=Asia/Shanghai

RUN echo "${TZ}" > /etc/timezone \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && apt update \
    && apt install -y tzdata \
    && rm -rf /var/lib/apt/lists/*
```
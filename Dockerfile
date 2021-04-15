FROM ubuntu:20.04 as ubuntu_base

RUN apt-get update && apt-get install -y git

FROM ubuntu_base as ubuntu_wget

RUN apt-get update && apt-get install -y wget

FROM ubuntu_base as ubuntu_git

RUN git clone --single-branch --depth 1 --branch v2.0.0 https://github.com/bazelbuild/remote-apis.git

FROM alpine:3.13 as alpine_base

RUN apk update && apk add git

FROM alpine_base as alpine_wget

RUN apk update && apk add wget

FROM alpine_base as alpine_git

RUN git clone --single-branch --depth 1 --branch v2.0.0 https://github.com/bazelbuild/remote-apis.git

FROM debian:stable as debian_base

RUN apt-get update && apt-get install -y git

FROM debian_base as debian_wget

RUN apt-get update && apt-get install -y wget

FROM debian_base as debian_git

RUN git clone --single-branch --depth 1 --branch v2.0.0 https://github.com/bazelbuild/remote-apis.git

FROM ubuntu:14.04

RUN apt-get update && apt-get -y install wget

ENV PATH /usr/local/go/bin:$PATH

ADD check /tmp/check
ADD in /tmp/in

RUN wget -qO- https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz | tar -C /usr/local -xzf - && \
      go build -o /tmp/resource/check /tmp/check/*.go && \
      go build -o /tmp/resource/in /tmp/in/*.go && \
      rm -rf /tmp/in /tmp/check && \
      rm -rf /usr/local/go

FROM ubuntu:xenial
MAINTAINER florije4ex <florije4ex@gmail.com>

RUN sed -i s/archive.ubuntu.com/mirrors.aliyun.com/g /etc/apt/sources.list \
    && sed -i s/security.ubuntu.com/mirrors.aliyun.com/g /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y build-essential wget vim\
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /root/.litecoin
ADD . /root/.litecoin
COPY . /data/
WORKDIR /data/
RUN wget https://download.litecoin.org/litecoin-0.16.0/linux/litecoin-0.16.0-x86_64-linux-gnu.tar.gz
RUN tar -zxvf litecoin-0.16.0-x86_64-linux-gnu.tar.gz
RUN cp litecoin-0.16.0/bin/* /usr/local/bin/
RUN chmod +x /usr/local/bin/litecoin*

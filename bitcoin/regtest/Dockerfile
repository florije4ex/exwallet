FROM ubuntu:xenial
MAINTAINER florije4ex <florije4ex@gmail.com>

RUN sed -i s/archive.ubuntu.com/mirrors.aliyun.com/g /etc/apt/sources.list \
    && sed -i s/security.ubuntu.com/mirrors.aliyun.com/g /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y wget \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /root/.bitcoin
ADD . /root/.bitcoin
COPY . /data/
WORKDIR /data/
RUN wget https://bitcoincore.org/bin/bitcoin-core-0.16.3/bitcoin-0.16.3-x86_64-linux-gnu.tar.gz
RUN tar -zxvf bitcoin-0.16.3-x86_64-linux-gnu.tar.gz
RUN cp bitcoin-0.16.3/bin/* /usr/local/bin/
RUN chmod +x /usr/local/bin/*

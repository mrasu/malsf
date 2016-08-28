FROM ubuntu:16.10

RUN apt update
RUN apt -y upgrade
RUN apt install -y apt-utils
RUN apt install -y wget
RUN apt install -y unzip
RUN apt install -y net-tools
RUN apt install -y git
RUN apt install -y vim
RUN apt install -y python3

WORKDIR /mine
RUN wget https://releases.hashicorp.com/consul/0.6.4/consul_0.6.4_linux_amd64.zip
RUN unzip consul_0.6.4_linux_amd64.zip
RUN wget https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz
RUN tar xzf go1.7.linux-amd64.tar.gz

ENV PATH $PATH:/mine
ENV PATH $PATH:/mine/go/bin
ENV GOPATH /mine/gopath
ENV GOBIN $GOPATH/bin
ENV GOROOT /mine/go

CMD go get github.com/mrasu/malsf/example

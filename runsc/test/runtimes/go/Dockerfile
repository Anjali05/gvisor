FROM ubuntu:bionic
ENV LANG_VER=1.12.5

RUN apt-get update && apt-get install -y \
  curl \
  gcc \
  git

WORKDIR /root

RUN curl -o go1.4.linux-amd64.tar.gz https://dl.google.com/go/go1.4.linux-amd64.tar.gz
RUN git clone https://go.googlesource.com/go
RUN mkdir bootstr
RUN tar -C bootstr -xzf go1.4.linux-amd64.tar.gz

ENV GOROOT=/root/go
ENV GOROOT_BOOTSTRAP=/root/bootstr/go
ENV LANG_DIR=${GOROOT}

WORKDIR ${GOROOT}
RUN git checkout go${LANG_VER}
WORKDIR /root/go/src
RUN ./make.bash
ENV PATH=$PATH:/root/go/bin

WORKDIR $GOROOT

COPY proctor-go.go ${LANG_DIR}

ENTRYPOINT ["/root/go/bin/go", "run", "proctor-go.go"]

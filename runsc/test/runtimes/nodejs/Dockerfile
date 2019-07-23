FROM ubuntu:bionic
ENV LANG_VER=12.4.0

RUN apt-get update && apt-get install -y \
  curl \
  dumb-init \
  g++ \
  make \
  python

WORKDIR /root
RUN curl -o go.tar.gz https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
RUN tar -zxf go.tar.gz
ENV PATH=$PATH:/root/go/bin

RUN curl -o node-v${LANG_VER}.tar.gz https://nodejs.org/dist/v${LANG_VER}/node-v${LANG_VER}.tar.gz
RUN tar -zxf node-v${LANG_VER}.tar.gz
ENV LANG_DIR=/root/node-v${LANG_VER}

WORKDIR ${LANG_DIR}
RUN ./configure
RUN make
RUN make test-build

COPY proctor-nodejs.go ${LANG_DIR}

ENTRYPOINT ["/usr/bin/dumb-init", "/root/go/bin/go", "run", "proctor-nodejs.go"]

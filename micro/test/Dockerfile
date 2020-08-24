FROM alpine:latest
RUN apk --no-cache add make git go gcc libtool musl-dev curl

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin


RUN apk add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

RUN         apk add --update ca-certificates openssl tar && \
            wget https://github.com/coreos/etcd/releases/download/v3.4.7/etcd-v3.4.7-linux-amd64.tar.gz && \
            tar xzvf etcd-v3.4.7-linux-amd64.tar.gz && \
            mv etcd-v3.4.7-linux-amd64/etcd* /bin/ && \
            apk del --purge tar openssl && \
            rm -Rf etcd-v3.4.7-linux-amd64* /var/cache/apk/*
VOLUME      /data
EXPOSE      2379 2380 4001 7001
ADD         scripts/run-etcd.sh /bin/run.sh

ENV MICRO_REGISTRY=etcd

COPY ./test/example-service ${GOPATH}/src/example-service
COPY ./test/config-example-service ${GOPATH}/src/config-example-service

# Speeding up tests by predownloading dependencies for services used.
RUN cd ${GOPATH}/src/example-service && go get && \
    cd ${GOPATH}/src/config-example-service && go get && \
    go get github.com/micro/services/location && \
    go get github.com/micro/services/helloworld

# RUN go get github.com/crufter/micro-services/logspammer

COPY ./micro /micro

ENTRYPOINT ["sh", "/bin/run.sh"]

ARG GOVERSION=1.15

FROM golang:${GOVERSION}-alpine AS build

RUN apk update && apk add make git

WORKDIR /go/src/crowdsec
COPY . .

RUN make release && \ 
    tar xzvf crowdsec-cloud-firewall-bouncer.tgz && \
    cd crowdsec-cloud-firewall-bouncer-v*/ && \
    install -v -m 755 -D ./crowdsec-cloud-firewall-bouncer "/usr/local/bin/crowdsec-cloud-firewall-bouncer"

FROM alpine:latest
COPY --from=build /usr/local/bin/crowdsec-cloud-firewall-bouncer /usr/local/bin/crowdsec-cloud-firewall-bouncer

COPY docker/entrypoint.sh /etc/crowdsec/start.sh

ENV CONFIG_PATH="/etc/crowdsec/config.d/config.yaml"

VOLUME /etc/crowdsec/config.d/

ENTRYPOINT /etc/crowdsec/start.sh
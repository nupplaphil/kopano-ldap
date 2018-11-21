FROM golang:1.10.0 as builder
WORKDIR /go/src/github.com/nupplaphil/kopano-ldap
COPY kopano-ld.go .
COPY cmd/*.go ./cmd/
COPY kopano/*.go ./kopano/
RUN go get -t -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kopano-ld .

FROM osixia/openldap
MAINTAINER Philipp Holzer <admin@philipp.info>

COPY --from=builder /go/src/github.com/nupplaphil/kopano-ldap/kopano-ld /usr/bin/

ADD docker/schema/*.schema /container/service/slapd/assets/config/bootstrap/schema
ADD docker/ldif/*.ldif /container/service/slapd/assets/config/bootstrap/ldif
ADD docker/environment /container/environment/01-custom
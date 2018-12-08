FROM golang:1.10.0 as builder
WORKDIR /go/src/github.com/nupplaphil/kopano-ldap
COPY kopano-ld.go .
COPY cmd/*.go ./cmd/
COPY kopano/*.go ./kopano/
RUN go get -t -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o kopano-ld .

FROM fbartels/openldap:1.2.3
MAINTAINER Philipp Holzer <admin@philipp.info>

COPY --from=builder /go/src/github.com/nupplaphil/kopano-ldap/kopano-ld /usr/bin/

ADD bootstrap /container/service/slapd/assets/config/bootstrap
ADD docker/environment /container/environment/01-custom
RUN rm /container/service/slapd/assets/config/bootstrap/schema/mmc/mail.schema
RUN touch /etc/ldap/slapd.conf
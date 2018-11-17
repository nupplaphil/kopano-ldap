FROM osixia/openldap
MAINTAINER Philipp Holzer <admin@philipp.info>

ADD schema/*.schema /container/service/slapd/assets/config/bootstrap/schema
ADD ldif/*.ldif /container/service/slapd/assets/config/bootstrap/ldif
ADD environment /container/environment/01-custom
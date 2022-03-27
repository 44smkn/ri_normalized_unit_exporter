ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="44smkn"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/aws_ri_exporter /bin/aws_ri_exporter

EXPOSE      9981
USER        nobody
ENTRYPOINT  [ "/bin/aws_ri_exporter" ]
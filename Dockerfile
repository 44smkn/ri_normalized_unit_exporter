ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="44smkn"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/ri_normalized_unit_exporter /bin/ri_normalized_unit_exporter

EXPOSE      9981
USER        nobody
ENTRYPOINT  [ "/bin/ri_normalized_unit_exporter" ]
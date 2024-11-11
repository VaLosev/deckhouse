ARG BASE_DISTROLESS
ARG BASE_GOLANG_23_ALPINE

FROM $BASE_GOLANG_23_ALPINE as artifact

ARG VERSION=3.15.0
ARG COMMIT_REF=b203438d231cbe743ea5bc5763ca56b2d92e4221

ARG SOURCE_REPO
ENV SOURCE_REPO=${SOURCE_REPO}

ARG GOPROXY
ENV GOPROXY=${GOPROXY}

RUN apk add --no-cache go git make
# RUN git clone --depth 1 --branch v3.17.0 ${SOURCE_REPO}/enix/x509-certificate-exporter.git /x509-certificate-exporter
# HACK: Using gh for testing, change to local repo
RUN git clone --depth 1 --branch v3.15.0 https://github.com/enix/x509-certificate-exporter.git /x509-certificate-exporter

WORKDIR /x509-certificate-exporter
RUN go build ./cmd/x509-certificate-exporter && \
    chown -R 64535:64535 /x509-certificate-exporter/ && \
    chmod 0700 /x509-certificate-exporter/x509-certificate-exporter
FROM $BASE_DISTROLESS

COPY --from=artifact /x509-certificate-exporter/x509-certificate-exporter /x509-certificate-exporter

ENTRYPOINT ["/x509-certificate-exporter"]

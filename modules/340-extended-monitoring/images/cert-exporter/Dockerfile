ARG BASE_DISTROLESS
ARG BASE_GOLANG_23_ALPINE

FROM $BASE_GOLANG_23_ALPINE as artifact

ARG VERSION=3.17.0
ARG COMMIT_REF=8f97b98c862f83d0c25c2994942b1ea90c6459da

ARG SOURCE_REPO
ENV SOURCE_REPO=${SOURCE_REPO}

ARG GOPROXY
ENV GOPROXY=${GOPROXY}

RUN apk add --no-cache go git make
# RUN git clone --depth 1 --branch v3.17.0 ${SOURCE_REPO}/enix/x509-certificate-exporter.git /x509-certificate-exporter
RUN git clone --depth 1 --branch v3.17.0 https://github.com/enix/x509-certificate-exporter.git /x509-certificate-exporter

WORKDIR /x509-certificate-exporter
RUN go build ./cmd/x509-certificate-exporter && \
    chown -R 64535:64535 /x509-certificate-exporter/ && \
    chmod 0700 /x509-certificate-exporter/x509-certificate-exporter
FROM $BASE_DISTROLESS

COPY --from=artifact /x509-certificate-exporter/x509-certificate-exporter /x509-certificate-exporter

ENTRYPOINT ["/x509-certificate-exporter"]

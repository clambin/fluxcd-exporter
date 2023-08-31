FROM alpine

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY bld/fluxcd-exporter-$TARGETOS-$TARGETARCH /app/fluxcd-exporter

EXPOSE 8080
ENTRYPOINT ["/app/fluxcd-exporter"]
CMD []

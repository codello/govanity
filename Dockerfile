ARG TARGETPLATFORM=linux/amd64
ARG TARGETOS=linux
ARG TARGETARCH=amd64
FROM --platform=$TARGETPLATFORM gcr.io/distroless/static

COPY --chmod=755 build/${TARGETOS}-${TARGETARCH}/govanity /govanity
USER nonroot:nonroot

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/govanity"]

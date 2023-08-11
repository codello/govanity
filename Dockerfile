FROM --platform=$TARGETPLATFORM gcr.io/distroless/static

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

COPY --chmod=755 build/${TARGETOS}-${TARGETARCH}/govanity /govanity
USER nonroot:nonroot

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/govanity"]

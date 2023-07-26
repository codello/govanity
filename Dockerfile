ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:alpine as builder

RUN apk add --no-cache git
WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

ARG VERSION
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} VERSION=${VERSION}

COPY . .
RUN go build -ldflags="-w -s -X main.Version=$VERSION" -o build/govanity .


FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static

COPY --from=builder /work/build/govanity /govanity

EXPOSE 8080
EXPOSE 9090
USER nonroot:nonroot
ENTRYPOINT ["/govanity"]

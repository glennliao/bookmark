FROM --platform=$BUILDPLATFORM golang:1.20 as builder
ARG TARGETARCH
WORKDIR /app
COPY main.go /app/main.go
COPY app /app/app
COPY packed /app/packed
COPY go.* /app/

RUN GOOS=linux GOARCH=$TARGETARCH  \
    gf build --gf.gcfg.file=gf.yaml

FROM alpine:latest

WORKDIR /root
COPY --from=builder /app/output/bookmark .
CMD /root/bookmark
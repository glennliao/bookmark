FROM --platform=$BUILDPLATFORM golang:1.20 as builder
ARG TARGETARCH
WORKDIR /app
COPY main.go /app/main.go
COPY app /app/app
COPY packed /app/packed
COPY config.toml.example /app/config.toml.example
COPY go.* /app

RUN wget -O gf https://github.com/gogf/gf/releases/latest/download/gf_$(go env GOOS)_$(go env GOARCH) && chmod +x gf && ./gf install -y && rm ./gf
RUN GOOS=linux GOARCH=$TARGETARCH  \
    gf build --gf.gcfg.file=gf.yaml

FROM alpine:latest

WORKDIR /root
COPY --from=builder /app/output/bookmark .
CMD /root/bookmark
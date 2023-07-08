FROM --platform=$BUILDPLATFORM golang:1.20 as builder
ARG TARGETARCH
WORKDIR /app
COPY main.go /app/main.go
COPY app /app/app
COPY packed /app/packed
COPY config.yaml.example /app/config.yaml.example
COPY go.* /app
COPY gf.yaml /app/gf.yaml

RUN wget -q -O gf https://github.com/gogf/gf/releases/latest/download/gf_$(go env GOOS)_$(go env GOARCH) && chmod +x gf && ./gf install -y && rm ./gf
RUN GOOS=linux GOARCH=$TARGETARCH  \
    gf build --gf.gcfg.file=gf.yaml

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bookmark .
RUN /app/bookmark init
CMD /app/bookmark
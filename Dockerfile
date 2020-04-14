FROM golang:1.14.2 AS builder
RUN useradd --create-home --shell /bin/bash docker
USER docker
RUN mkdir -p /go/github.com/AAkindele/k8s_pod_logger
WORKDIR /go/github.com/AAkindele/k8s_pod_logger
ADD --chown=docker . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/k8s_pod_logger github.com/AAkindele/k8s_pod_logger


FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser docker --home /home/docker --disabled-password
USER docker
WORKDIR /home/docker
COPY --from=builder /go/bin/k8s_pod_logger .
ENTRYPOINT ["./k8s_pod_logger"]
CMD []

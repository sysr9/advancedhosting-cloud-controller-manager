FROM golang:1.16.0-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /go/src/github.com/advancedhosting/advancedhosting-cloud-controller-manager

COPY . .


ARG TAG
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-w -s -X github.com/advancedhosting/advancedhosting-cloud-controller-manager/advancedhosting.version=${TAG}" -o advancedhosting-cloud-controller-manager ./cmd/cloud-controller-manager

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/advancedhosting/advancedhosting-cloud-controller-manager/advancedhosting-cloud-controller-manager .
ENTRYPOINT ["/advancedhosting-cloud-controller-manager"]

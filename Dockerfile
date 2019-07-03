# Build the manager binary
FROM golang:1.11.4 as builder

ENV GO111MODULE=on

# Copy in the go src
WORKDIR /go/src/github.com/kubeflow/xgboost-operator

COPY go.mod .
COPY go.sum .
# COPY vendor/ vendor/

RUN go mod download

COPY pkg/    pkg/
COPY cmd/    cmd/
COPY config/ config/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager github.com/kubeflow/xgboost-operator/cmd/manager

# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /
COPY --from=builder /go/src/github.com/kubeflow/xgboost-operator/manager .
ENTRYPOINT ["/manager"]

################################################################################################
# Builder image
# See https://hub.docker.com/_/golang/
################################################################################################
FROM golang:1.20 as builder

ARG OS=linux
ARG ARCH=amd64

WORKDIR /usr/src/app

# pre-copy/cache parent go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o check-argocd main.go

################################################################################################
# user-identity-mapper image to be run by the job on OpenShift
################################################################################################
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest as check-argocd

# Copy the generated binary into the $PATH so it can be invoked
COPY --from=builder /usr/src/app/check-argocd /usr/local/bin/

# Run as non-root user
USER 1001

ENTRYPOINT ["/usr/local/bin/check-argocd", "check"]
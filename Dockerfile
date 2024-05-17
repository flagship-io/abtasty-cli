FROM golang:1.18-alpine as builder
RUN apk add --update make
WORKDIR /go/src/github/flagship-io/abtasty-cli

ARG ABTASTY_CLI_VERSION
ENV ABTASTY_CLI_VERSION $ABTASTY_CLI_VERSION

# Download dependencies before building
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github/flagship-io/abtasty-cli/abtasty-cli ./
CMD ["/bin/sh"]
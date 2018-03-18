# Golang build step
FROM golang:1.10-alpine as builder
RUN apk update && apk add curl git
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep
ADD . /go/src/github.com/kamaln7/klein
WORKDIR /go/src/github.com/kamaln7/klein
RUN dep ensure -vendor-only
RUN go install

# Copy go binary and static assets
FROM alpine:3.7
# Add ca-certificates so that we can talk to DigitalOcean Spaces
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/bin/klein .
COPY --from=builder /go/src/github.com/kamaln7/klein/404.html .

CMD ["./klein"]

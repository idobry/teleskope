FROM golang:latest AS builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /go/src/teleskope
WORKDIR /go/src/teleskope
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/teleskope ./
CMD ["./teleskope run"]

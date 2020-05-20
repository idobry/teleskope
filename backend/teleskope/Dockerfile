FROM golang:latest AS builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
#RUN go get -u "github.com/gorilla/handlers"
#RUN go get -u "github.com/gorilla/mux"
#RUN go get -u "github.com/mitchellh/go-homedir"
#RUN go get -u "github.com/spf13/cobra"
#RUN go get -u "github.com/spf13/viper"
WORKDIR /go/src/teleskope
#WORKDIR /go/src/github.com/idobry/teleskope/backend/teleskope
ADD . .
RUN dep ensure
#RUN CGO_ENABLED=0 GOOS=linux go build -a
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o teleskope .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/teleskope ./
CMD ["./teleskope run"]

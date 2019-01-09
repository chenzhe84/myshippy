FROM golang as builder
WORKDIR /go/src/vessel-service
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/micro/go-micro
# RUN dep init && dep ensure
RUN go build main.go

FROM debian
WORKDIR /app
COPY --from=builder /go/src/vessel-service/main .
CMD ["./main"]
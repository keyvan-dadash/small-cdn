FROM golang AS builder

WORKDIR /go/src/small-cdn

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o small-cdn .

FROM ubuntu:latest  

WORKDIR /root/

COPY --from=builder /go/src/small-cdn .

EXPOSE 8080

CMD ["./small-cdn"]

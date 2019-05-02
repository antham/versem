FROM golang:1.12 as build

ENV GO111MODULE=on
WORKDIR /go/src/github.com/antham/versem
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:3.6

COPY --from=build /go/src/github.com/antham/versem/main /usr/sbin/
ENTRYPOINT ["/usr/sbin/main"]

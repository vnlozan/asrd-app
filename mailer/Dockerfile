FROM golang:1.24 AS builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o . ./...

FROM scratch

COPY --from=builder /build/app .

EXPOSE 8080

CMD ["/app"]
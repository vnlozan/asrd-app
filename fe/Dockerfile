FROM golang:1.24 AS builder

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o . ./...

FROM scratch

COPY --from=builder /build/app .
COPY /assets /assets

EXPOSE 8080

CMD ["/app"]
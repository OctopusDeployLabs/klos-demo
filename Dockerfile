FROM golang:1.24 AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin ./...

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/* /
CMD ["/worker"]
FROM golang:1.14 as builder

WORKDIR /usr/src/avatars.io

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /usr/src/avatars.io/avatars.io /usr/bin/

CMD ["/usr/bin/avatars.io"]

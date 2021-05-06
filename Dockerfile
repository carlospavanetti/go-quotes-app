# Build phase
FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o main .

# Final image phase
FROM  gcr.io/distroless/static:nonroot

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

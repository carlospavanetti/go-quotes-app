# Build phase
FROM golang:alpine as builder

RUN apk --no-cache add ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \  
    "${USER}"

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o main .

# Final image phase
FROM scratch

WORKDIR /root/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/main .

USER appuser:appuser

EXPOSE 8080

CMD ["./main"]

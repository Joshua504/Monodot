# ---------- Builder ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o monodot ./cmd/web

# ---------- Runtime ----------
FROM alpine:latest

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/monodot .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

RUN mkdir uploads outputs && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./monodot"]
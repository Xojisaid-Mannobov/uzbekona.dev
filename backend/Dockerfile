# ─── Build ───────────────────────────────────────────────
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Statik binary: kichik image, qo'shimcha runtime kerak emas
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ─── Runtime ─────────────────────────────────────────────
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app \
    && mkdir -p /app/storage/uploads && chown -R app:app /app

WORKDIR /app
COPY --from=build /out/api /app/api

USER app
ENV PORT=8080 UPLOAD_DIR=/app/storage/uploads
EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=3s --retries=5 \
  CMD wget -qO- http://127.0.0.1:8080/api/v1/health || exit 1

ENTRYPOINT ["/app/api"]

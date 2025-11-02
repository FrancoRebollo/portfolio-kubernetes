# =======================
# STAGE 1: Build
# =======================
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/apis

# Crear la carpeta logs y asignar permisos desde el builder
RUN mkdir -p /app/logs && chown -R 65532:65532 /app/logs
# 65532 = UID/GID de `nonroot` en distroless

# =======================
# STAGE 2: Run
# =======================
FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/logs /app/logs
COPY --from=builder /app/.env /app/.env

USER nonroot
EXPOSE 3004

ENTRYPOINT ["/app/server"]
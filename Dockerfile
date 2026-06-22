# --- Stage 1: Build Go backend ---
FROM golang:1.23-alpine AS go-builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# --- Stage 2: Build Astro frontend ---
FROM node:20-alpine AS web-builder

WORKDIR /web

COPY web/package.json web/package-lock.json ./
RUN npm ci --ignore-scripts

COPY web/ ./

RUN npm run build

# --- Stage 3: Production image ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates nodejs npm

WORKDIR /app

COPY --from=go-builder /app/server ./server
COPY --from=web-builder /web/dist ./web/dist

RUN mkdir -p /app/data

ENV PORT=8080
ENV ALLOWED_ORIGINS=*

EXPOSE 8080 4321

COPY docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/api/health || exit 1

ENTRYPOINT ["./docker-entrypoint.sh"]

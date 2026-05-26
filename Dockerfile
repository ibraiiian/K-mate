# Stage 1: Build Frontend Vue
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend Go
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/backend
# Tambahkan env agar compile menjadi single binary tanpa cgo dependency
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Build binary Go (dengan flag -s -w untuk memperkecil ukuran file)
RUN go build -ldflags="-s -w" -o kmate-api ./cmd/api

# Stage 3: Production Image (Sangat ringan)
FROM alpine:3.19
WORKDIR /app

# Install package dasar untuk HTTP dan timezone
RUN apk add --no-cache tzdata ca-certificates

# Copy backend binary
COPY --from=backend-builder /app/backend/kmate-api .

# Copy frontend hasil build (Vue dist) ke folder ./dist
COPY --from=frontend-builder /app/frontend/dist ./dist

# Di Cloud Run, Port selalu disediakan melalui Env Var PORT
# Default fallback port (biasanya 8080)
ENV PORT=8080
EXPOSE 8080

CMD ["./kmate-api"]

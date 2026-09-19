# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o go-rest-wol .

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.19
WORKDIR /app

# 拷贝二进制和前端页面
COPY --from=builder /app/go-rest-wol .
COPY --from=builder /app/pages/index.html ./pages/index.html

# 数据目录（挂载持久化 SQLite 用）
RUN mkdir -p /data

ENV WOLHTTPPORT=8080
ENV WOLDB=/data/computer.db

EXPOSE 8080

CMD ["/app/go-rest-wol"]

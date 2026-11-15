# ---------- 构建阶段 ----------
# 使用 buildx 提供的 TARGETOS / TARGETARCH 实现多架构交叉编译
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /src

# 先拷贝依赖描述文件，利用构建缓存
COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/move_mp4 .

# ---------- 运行阶段 ----------
FROM alpine:latest

# 安装 ca-certificates 和时区数据（确保时间戳正确）
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /out/move_mp4 /app/move_mp4

# 监控目录 /xunlei，目标目录 /tmp，运行时按需挂载
ENTRYPOINT ["/app/move_mp4"]

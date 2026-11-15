# mp4-mover

一个极简的 Go 小工具：定时把迅雷下载目录里的 `.mp4` 文件搬运到目标目录。

- 每 5 分钟扫描一次 `/xunlei`
- 只处理 `.mp4`（大小写不敏感），自动跳过子目录
- 目标目录已有同名文件时，自动追加时间戳重命名，不覆盖
- 同分区走 `rename`（瞬时完成），跨分区自动降级为「复制 + 删除」
- 单文件静态二进制，无任何运行时依赖；镜像仅 ~10MB

## 快速开始

### Docker（推荐）

```bash
docker run -d --name mp4-mover \
  -v /你的/迅雷下载目录:/xunlei \
  -v /你的/目标目录:/tmp \
  --restart unless-stopped \
  ghcr.io/<你的GitHub用户名>/mp4-mover:latest
```

### 二进制

从 [Releases](https://github.com/<你的GitHub用户名>/mp4-mover/releases) 下载对应架构的文件：

```bash
chmod +x move_mp4_linux_amd64
./move_mp4_linux_amd64
```

### 本地构建

```bash
go build -o move_mp4 .                       # 本机架构
docker build -t mp4-mover .                  # 镜像（多阶段构建，无需预编译）
```

## 镜像与标签

镜像发布在 GitHub Container Registry：

| 触发方式 | 产物 |
| --- | --- |
| 推送到 `main` / `master` | `ghcr.io/<owner>/mp4-mover:latest`、`:main`、`:sha-xxxxxxx` |
| 推送 `v1.2.3` 标签 | `:1.2.3`、`:1.2`、`:latest`、`:sha-xxxxxxx`，同时创建 GitHub Release 并附带二进制 |
| Pull Request | 仅构建验证，不推送 |

支持架构：`linux/amd64`、`linux/arm64`、`linux/arm/v7`。

## CI / CD 说明

| 工作流 | 文件 | 作用 |
| --- | --- | --- |
| CI | `.github/workflows/ci.yml` | `gofmt`、`go vet`、linux 交叉编译检查 |
| Release | `.github/workflows/release.yml` | 打 `v*` 标签时编译多架构二进制 + `checksums.txt` 并发布 Release |
| Docker | `.github/workflows/docker.yml` | 构建多架构镜像并推送到 ghcr.io |

### 发布新版本

```bash
git tag v1.0.0
git push origin v1.0.0
```

工作流会自动完成：GitHub Release + 多架构镜像推送。

> **首次使用需要先让镜像公开**：GHCR 包第一次推送后默认是 private（即使仓库是公开的）。
> 进入仓库右侧 **Packages → mp4-mover → Package settings → Change visibility → Public** 一次即可，
> 之后匿名 `docker pull` 就能用了。

## 配置

路径与扫描间隔目前是 `main.go` 顶部的常量，需要改动时改完重新构建即可：

```go
sourceDir     = "/xunlei"
targetDir     = "/tmp"
scanInterval  = 5 * time.Minute
fileExtension = ".mp4"
```

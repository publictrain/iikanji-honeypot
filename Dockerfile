# 基本イメージとしてgolangの公式イメージを使用
FROM golang:1.20 as builder

# ワークディレクトリの設定
WORKDIR /app

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o ssh-honeypot

# マルチステージビルド
# 新しいステージを開始して、alpineイメージをベースにする
FROM alpine:latest  

# 作業ディレクトリを/appに設定
WORKDIR /app

# ビルダーステージからビルド済みバイナリと必要なファイルをコピー
COPY --from=builder /app/ssh-honeypot .
# COPY cmds.txt .

# 実行時にアプリケーションを起動
CMD ["./ssh-honeypot"]

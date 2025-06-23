# GoアプリケーションのDockerコンテナセットアップ
# 依存関係のダウンロードを行い、'app'という名前の実行可能ファイルを生成
# './cmd'はアプリケーションのエントリポイントがあるディレクトリを指します。

FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

CMD ["./app"]
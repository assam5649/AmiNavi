# GoアプリケーションのDockerコンテナセットアップ
# 依存関係のダウンロードを行い、'app'という名前の実行可能ファイルを生成
# './cmd'はアプリケーションのエントリポイントがあるディレクトリを指します。
FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.0/migrate.linux-amd64.tar.gz -o migrate.tar.gz \
    && tar xvz -f migrate.tar.gz\
    && mv migrate /usr/bin/migrate \
    && rm migrate.tar.gz

RUN go build -o server ./cmd/api

COPY scripts/entry.sh /app/scripts/entry.sh
RUN chmod +x /app/scripts/entry.sh

CMD ["./scripts/entry.sh"]

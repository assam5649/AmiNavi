## あみナビ

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://golang.org/)
[![MySQL](https://img.shields.io/badge/Database-MySQL-blue.svg?logo=mysql)](https://www.mysql.com/)
[![Docker](https://img.shields.io/badge/Container-Docker-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

本リポジトリは、Androidアプリとセンサー付き編み棒からなるスマート編み物支援システム「あみナビ」のバックエンドです。<br>
本バックエンドではユーザー管理、編み図データの管理、モバイルアプリとの同期処理を担当します。

## API

### 認証
Firebase Authentication による認証を用い、ユーザーの識別を行う。<br>
ユーザー名やアイコンなどのアプリケーション固有のプロフィール情報は、 Firebase UIDとサーバー側のDBを紐付ける
 - **POST /v1/register** - 新しいユーザーを登録します

### 編み図管理
 - **GET /v1/works** - すべての作品を一覧を返します(ユーザー別にフィルタリングします)
 - **POST /v1/works** - 新しい作品を登録します
 - **PUT /v1/works/{id}** - 特定の作品のメタデータを更新します
 - **DELETE /v1/works/{id}** - 作品を削除します
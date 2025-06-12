## あみナビ

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://golang.org/)
[![MySQL](https://img.shields.io/badge/Database-MySQL-blue.svg?logo=mysql)](https://www.mysql.com/)
[![Docker](https://img.shields.io/badge/Container-Docker-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

本リポジトリは、Androidアプリとセンサー付き編み棒からなるスマート編み物支援システム「あみナビ」のバックエンドです。<br>
本バックエンドではユーザー管理、編み図データの管理、モバイルアプリとの同期処理を担当します。

## API

### 認証
 - **POST /register** - 新しいユーザーを登録します
 - **POST /login** - 既存のユーザーを認証します

### 編み図管理
 - **GET /works** - すべての作品を一覧を返します(ユーザー別にフィルタリングします)
 - **POST /works** - 新しい作品を登録します
 - **PUT /works/{id}** - 特定の作品のメタデータを更新します
 - **DELETE /works/{id}** - 作品を削除します
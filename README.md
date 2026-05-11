# my-go-app

Go API 専用のサンプルプロジェクトです。  
Docker の学習を兼ねて、**開発用** と **本番想定用** の 2 つの実行方式を同じリポジトリで扱える構成にしています。

## このプロジェクトの構成

- 開発用: `docker-compose.yml` + `Dockerfile` の `dev` ステージ + `.air.toml`
- 本番想定: `Dockerfile` の `prod` ステージ（マルチステージビルド + 非 root）
- API: `cmd/api/main.go`（`/healthz` と `/` の最小エンドポイント）

## 前提

- Docker Desktop がインストール済み
- `docker compose` が使えること

確認コマンド:

```bash
docker --version
docker compose version
```

## コミット時のコード整形（Git フック）

コミット直前に、**ステージ済みの `.go` ファイル**へ `gofmt` をかけ、整形結果を同じコミットに含めます。  
（いわゆる **pre-commit フック** です。ショートカットキーではなく、Git がコミット時に自動で呼び出します。）

### 初回だけ（このリポジトリで有効化）

```bash
git config core.hooksPath .githooks
chmod +x .githooks/pre-commit
```

別マシンでクローンしたら、`core.hooksPath` をもう一度設定してください（ローカル設定のためリポジトリには含まれません）。

### チェックだけしてコミットを止めたい場合

デフォルトは自動整形です。整形漏れを許さず **失敗させたい** ときは、環境変数を付けます。

```bash
PRE_COMMIT_GOFMT_CHECK=1 git commit
```

### 手動

```bash
make fmt        # 全 .go を整形
make fmt-check  # 未整形があれば終了コード 1（CI 向け）
```

## 開発環境の構築手順

`air` によるホットリロードを使います。  
Go ファイルを保存すると、コンテナ内で自動ビルド・自動再起動されます。

1. イメージをビルドして起動

```bash
docker compose up --build
```

2. 別ターミナルで API を確認

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/
```

期待値:

- `/healthz` -> `ok`
- `/` -> `{"message":"Go API is running"}`

※ 現在のエンドポイント設計は REST API の思想に合わせ、アクション名ではなくリソース指向のパス（例: `/todos`）を採用しています。将来的に `POST/PUT/PATCH/DELETE` などへ拡張する際も、HTTP メソッドごとに責務を分けます。

3. 開発を終了する

```bash
docker compose down
```

## 本番想定イメージの構築手順

本番想定では `air` を使わず、ビルド済みバイナリを直接起動します。

1. 本番想定イメージをビルド

```bash
docker build --target prod -t my-go-api:prod .
```

2. コンテナ起動

```bash
docker run --rm -p 8080:8080 my-go-api:prod
```

3. 動作確認

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/
```

## よく使う補助コマンド

- コンテナをバックグラウンド起動

```bash
docker compose up -d --build
```

- ログ確認

```bash
docker compose logs -f app
```

- 完全停止（ネットワーク/コンテナ削除）

```bash
docker compose down
```

- ボリュームも削除して初期化

```bash
docker compose down -v
```

## ファイルの役割

- `Dockerfile`  
  `dev` / `builder` / `prod` のマルチステージ構成。  
  `dev` ステージには開発用ツールとして `air` と SQLBoiler CLI を入れています。  
  `prod` ステージにはビルド済みバイナリのみを配置します。

- `docker-compose.yml`  
  開発時の起動設定。  
  `app` service では `dev` ステージを使い、コードマウントとポート公開を行います。  
  `db` service では PostgreSQL を起動します。

- `db/init.sql`  
  PostgreSQL 初期化用 SQL です。  
  初回 DB volume 作成時に実行され、`todos` テーブルなどを作成します。

- `sqlboiler.toml`  
  SQLBoiler の設定ファイルです。  
  DB 接続先や生成コードの出力先を定義します。

- `.air.toml`  
  開発専用のホットリロード設定ファイルです（本番想定では使用しません）。

- `.dockerignore`  
  Docker ビルド時に不要ファイルを除外し、ビルドを軽量化します。

- `cmd/api/main.go`  
  API サーバのエントリーポイントです。  
  ルーティング、サーバ起動、graceful shutdown などを扱います。

- `internal/domain/todo`  
  Todo のドメインモデル、Repository interface、ドメインエラーなどを定義します。

- `internal/usecase/todo`  
  Todo に関するアプリケーションロジックを扱います。  
  保存先がメモリか DB かは知りません。

- `internal/infrastructure/todo`  
  メモリ保存など、具体的な保存実装を置きます。

- `internal/infrastructure/sqlboiler/models`  
  SQLBoiler によって自動生成される DB モデルです。  
  手で編集せず、DB スキーマ変更後に再生成します。

- `internal/handler/todo`  
  HTTP リクエストとレスポンスを扱う層です。  
  JSON の decode / encode や HTTP status の変換を担当します。

- `.githooks/pre-commit`  
  コミット前に `gofmt` を実行する Git フックです。

- `Makefile`  
  `make fmt` / `make fmt-check` で手動整形・整形チェックができます。

## 依存関係メモ

SQLBoiler 関連の主な依存は以下です。

```txt
github.com/aarondl/sqlboiler/v4 v4.19.7
github.com/aarondl/null/v8 v8.1.3
github.com/lib/pq v1.10.9
```

SQLBoiler CLI は `go.mod` ではなく、`Dockerfile` の `dev` ステージでインストールします。

```dockerfile
RUN go install github.com/aarondl/sqlboiler/v4@v4.19.7
RUN go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@v4.19.7
```
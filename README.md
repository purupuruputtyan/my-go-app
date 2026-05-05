# my-go-app

Go API 専用のサンプルプロジェクトです。  
Docker の学習を兼ねて、**開発用** と **本番想定用** の 2 つの実行方式を同じリポジトリで扱える構成にしています。

## このプロジェクトの構成

- 開発用: `docker-compose.yml` + `Dockerfile` の `dev` ステージ + `.air.toml`
- 本番想定: `Dockerfile` の `prod` ステージ（マルチステージビルド + 非 root）
- API: `main.go`（`/healthz` と `/` の最小エンドポイント）

## 前提

- Docker Desktop がインストール済み
- `docker compose` が使えること

確認コマンド:

```bash
docker --version
docker compose version
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
  `dev` / `builder` / `prod` のマルチステージ構成。開発と本番想定を分離します。

- `docker-compose.yml`  
  開発時の起動設定。`dev` ステージを使い、コードマウントとポート公開を行います。

- `.air.toml`  
  開発専用のホットリロード設定ファイルです（本番想定では使用しません）。

- `.dockerignore`  
  Docker ビルド時に不要ファイルを除外し、ビルドを軽量化します。

- `main.go`  
  最小限の API サーバ実装です。

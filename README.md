# vrc-outfit-shelf

VRChat向けのアセット管理サービス。購入したアセットと、それをどのアバターにどう使ったかを記録します。

## なぜ作ったか

VRChatで活動する中で、「以前買ったアセットがどこにあるか分からない」「どのアバターに適用したのか分からない」という声を周りから何度も聞きました。困りごとは「何を持っているか」と「それをどこにどう使ったか」の二つに分かれるため、この二つを別々に記録し、後から結びつけられる構造を前提に設計しています。

## 技術構成

Go / PostgreSQL / Docker Compose（pgx, net/http）

## 現在の状態

開発中です。データベースの設計・構築と、アセットの一覧表示・登録まで実装しています。
今後は認証、所有権チェック、残りのAPI、html/template による画面、テスト、デプロイを予定しています。

## 動かし方

```bash
# 1. データベースを起動
docker compose up -d

# 2. .env を用意（.env.example をコピーして値を入れる）
cp .env.example .env

# 3. スキーマとシードデータを投入
docker compose exec -T db psql -U seolhui -d maindb < schema.sql
docker compose exec -T db psql -U seolhui -d maindb < seed.sql

# 4. 接続情報を環境変数に設定
export DATABASE_URL="postgres://seolhui:<password>@localhost:5432/maindb"

# 5. サーバーを起動
go run .
```

PowerShell の場合は 3 と 4 を以下に置き換えてください。

```powershell
Get-Content schema.sql | docker compose exec -T db psql -U seolhui -d maindb
Get-Content seed.sql | docker compose exec -T db psql -U seolhui -d maindb
$env:DATABASE_URL = "postgres://seolhui:<password>@localhost:5432/maindb"
```

起動後 http://localhost:8080/assets で確認できます。
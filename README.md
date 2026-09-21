# vrc-outfit-shelf

VRChat向けのアセット管理サービス。購入したアセットと、それをどのアバターにどう使ったかを記録します。

## なぜ作ったか

VRChatで活動する中で、「以前買ったアセットがどこにあるか分からない」「どのアバターに適用したのか分からない」という声を周りから何度も聞きました。困りごとは「何を持っているか」と「それをどこにどう使ったか」の二つに分かれるため、この二つを別々に記録し、後から結びつけられる構造を前提に設計しています。

## 技術構成

Go / PostgreSQL / Docker Compose（pgx, net/http, html/template, bcrypt）

## 現在の状態

開発中です。現在までに以下を実装しています。

- データベースの設計・構築とシードデータ
- アセットの一覧表示・登録
- html/template による画面描画
- 会員登録（パスワードは bcrypt でハッシュ化して保存）
- ログイン時のパスワード照合

今後はセッションによるログイン状態の保持、ログアウト、所有権チェック、残りのAPI、テスト、デプロイを予定しています。

## 設計上の判断

- 認証は JWT ではなくセッション方式を採用予定です。サーバー1台・サーバーサイドレンダリングの構成では JWT の利点が生きず、ログアウトを確実に処理できるためです
- 画面は html/template で描画し、利用者の入力を自動でエスケープすることで XSS を防いでいます
- ログイン失敗時は「ユーザーが存在しない」と「パスワードが違う」を区別せず、同じメッセージとステータスコードを返します

## 動かし方

```bash
# 1. .env を用意（.env.example をコピーし、パスワードを設定）
cp .env.example .env
# .env の中身: POSTGRES_PASSWORD=任意のパスワード

# 2. データベースを起動
docker compose up -d

# 3. スキーマとシードデータを投入
docker compose exec -T db psql -U seolhui -d maindb < schema.sql
docker compose exec -T db psql -U seolhui -d maindb < seed.sql

# 4. 接続情報を環境変数に設定（パスワードは .env と同じもの）
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

起動後、以下のページで確認できます。

- アセット一覧・登録: http://localhost:8080/assets
- 会員登録: http://localhost:8080/signup
- ログイン: http://localhost:8080/signin

シードデータのユーザーはパスワードが設定されていないためログインできません。動作確認の際は先に会員登録を行ってください。
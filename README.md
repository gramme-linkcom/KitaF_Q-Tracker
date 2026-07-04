# Q-Tracker
待ち時間をリアルタイムで確認できるWebアプリ

# 使用技術
本システムに使用している技術は下記のとおりです。
- Docker (コンテナ仮想化プラットフォーム)
- Docker compose (コンテナ管理ツール)
- Golang (バックエンド)
- Next.js (フロントエンド)
- Daisy UI (CSSフレームワーク)
- Bootstrap (CSSフレームワーク)
- Sqlite (データベース)

# 利用方法
本システムの推奨動作環境はDockerになります。
## Docker composeを用いたサービスの開始方法
1. 本リポジトリを任意のディレクトリにて`git clone`します
2. ポート番号などを変更する必要がある場合は、`docker-compose.yml`の内容を編集してください。
3. `docker compose up -d --build`を実行します。
4. 環境にもよりますが、おおよそ2分前後でサービスが立ち上がります。

# ライセンス
本システムのソースコードはMIT LICENSEにて提供します。
MIT LICENSEの詳細は、本リポジトリの[LICENSE](https://github.com/gramme-linkcom/KitaF_Q-Tracker/blob/main/LICENSE)よりご確認ください。

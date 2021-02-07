# slack-egosa
## 環境変数
| 環境変数 | 説明 | 値 |
| ---- | ---- | ---- |
| SLACK_TOKEN | OAuth Token | 必須項目。`xoxp-` から始まるもの。`search:read` 権限が必要。 |
| WEBHOOK_URL | Webhook URL | 非デバッグ時は必須項目。結果の通知先。`https://hooks.slack.com.services/` から始まるもの。`links:write` 権限を推奨。|
| KEYWORDS | 検索候補のキーワードの列 | 必須項目。`,` 区切りで与える。例: `hoge,fuga` |
| MAX_SEARCH_COUNT | 検索数の上限値 | デフォルト20。最大100。|
| DURATION_HOURS | 現在からの検索期間の時間。now から (now - DURATION_HOURS) までの間が検索候補となる | デフォルト24。|
| DEBUG | デバッグを有効にするかどうか。有効の場合は，Slackへの投稿は行わずに，標準出力に結果を出力する。| 非空の場合にデバッグモード有効となる。|

## 使い方
```
$ SLACK_TOKEN=xoxp-... WEBHOOK_URL=https://hooks.slack.com/services/... KEYWORDS=hoge,fuga... MAX_SEARCH_COUNT=3 DURATION_HOURS=23 go run main.go
```

### デバッグ
```
$ DEBUG=true SLACK_TOKEN=xoxp-... WEBHOOK_URL=https://hooks.slack.com/services/... KEYWORDS=hoge,fuga... MAX_SEARCH_COUNT=3 DURATION_HOURS=23 go run main.go
```

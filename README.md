# mastrain

Extracts Mastodon bookmarks and puts them into the unsorted folder in Raindrop.

Automatically eliminates bookmarks that are already copied over.

Create a .env file with credentials - you'll need app credentials for Mastodon and an app token for Raindrop.

```
MASTODON_SERVER="https://mastodon.org.uk"
MASTODON_CLIENT_ID=""
MASTODON_CLIENT_SECRET=""
MASTODON_APP_TOKEN=""
RAINDROP_SERVER="https://api.raindrop.io/rest/v1/"
RAINDROP_APP_TOKEN=""
```

Run with:

`go run .`

for now. Better deployment coming soon.

\* considering moving to dropping them into unsorted. Change my mind.


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

Run from source with:

`go run .`

Or download a binary from the releases. You still need to make a .env file.




# Bellboy

Usage:

`bellboy`

Example configuration file (~/.bellboy/bellboy.conf):

```go
{
  "tumblr": {
    "consumerKey": "SECRET_KEY",
    "consumerSecret": "SECRET_KEY",
    "oauthKey": "SECRET_KEY",
    "oauthSecret": "SECRET_KEY",
    "blog": "myblog"
  },

  "db": "~/.bellboy/bellboy.db",
  "mediaFolder": "~/.bellboy/media"
}
```

Dependencies (for tests):

- go get gopkg.in/jarcoal/httpmock.v1

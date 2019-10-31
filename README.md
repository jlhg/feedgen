# feedgen

Convert website content to RSS Feed.

## Getting Started

Start the web app by docker-compose:

```
docker-compose up -d
```

Then go to `http://localhost:8080/<site>`.

## Supported Sites and Examples

- `hackernews`: Hacker News
    - Top link: https://feedgen.org/hackernews/best
- `ptt`: 批踢踢實業坊
    - 八卦版：https://feedgen.org/ptt/Gossiping
    - 表特版 (20 推以上): https://feedgen.org/ptt/Beauty?q=recommend:30
- `gamer_forum`: 巴哈姆特-哈拉區
    - 場外休憩區：https://feedgen.org/gamer_forum/60076
    - 魔物獵人 (20 推以上)：https://feedgen.org/gamer_forum/5786?gp=20

## License

feedgen is released under the MIT license.

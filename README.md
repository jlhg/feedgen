# feedgen

A middleware to generate Atom feeds from websites. https://feedgen.org

## Supported Sites and Examples

- `hackernews`: Hacker News
    - Top link: https://feedgen.org/hackernews?category=best
- `ptt`: 批踢踢實業坊
    - Steam 版：https://feedgen.org/ptt?b=Steam
    - 電影版 (30 推以上): https://feedgen.org/ptt?b=movie&q=recommend:30
- `gamer_forum`: 巴哈姆特-哈拉區
    - 精靈寶可夢：https://feedgen.org/gamer_forum?bsn=1647
    - 魔物獵人 (20 推以上)：https://feedgen.org/gamer_forum?bsn=5786&gp=20
- `udn_game`: 遊戲角落
    - 最新文章: https://feedgen.org/udn_game?section=rank&by=newest
    - 最多瀏覽: https://feedgen.org/udn_game?section=rank&by=pv

## Getting Started

Start the web by docker-compose:

```
docker-compose up -d
```

It's on `http://localhost:8080`.

## License

feedgen is released under the MIT license.

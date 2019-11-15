# feedgen

A middleware to generate Atom feeds from websites. https://feedgen.org

## Supported Sites and Examples

- `hackernews`: Hacker News
    - Top link: https://feedgen.org/hackernews?category=best
- `ptt`: 批踢踢實業坊
    - 八卦版：https://feedgen.org/ptt?b=Gossiping
    - 表特版 (30 推以上): https://feedgen.org/ptt?b=Beauty&q=recommend:30
- `gamer_forum`: 巴哈姆特-哈拉區
    - 場外休憩區：https://feedgen.org/gamer_forum?bsn=60076
    - 魔物獵人 (20 推以上)：https://feedgen.org/gamer_forum?bsn=5786&gp=20

## Getting Started

Start the web by docker-compose:

```
docker-compose up -d
```

It's on `http://localhost:8080`.

## License

feedgen is released under the MIT license.

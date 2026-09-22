# pico-rss: an RSS server for ants

pico-rss is a very minimal zero-dependency personal RSS feed server that serves `.md` files using [gorilla/feeds](https://github.com/gorilla/feeds) with a Dockerhub release at [mcbalaam/pico-rss](https://hub.docker.com/r/mcbalaam/pico-rss)

## Running pico-rss

### From source

```sh
git clone https://github.com/mcbalaam/pico-rss && cd pico-rss
docker compose up -d --build
```

### From Docker Hub with compose

```yaml
services:
  pico-rss:
    image: mcbalaam/pico-rss:latest
    container_name: pico-rss
    ports:
      - "8066:8066"
    environment:
      RSS_FEED_TITLE: "mcbalaam notes"
      RSS_FEED_DESCRIPTION: "Notes by mcbalaam"
      RSS_AUTHOR_NAME: "mcbalaam"
      RSS_AUTHOR_EMAIL: "mcbalaam@example.com"
      RSS_ORIGIN_URL: "https://rss.mcblm.xyz"
      RSS_TARGET_DIR: /data
      RSS_TIMEOUT: 10
    volumes:
      - "/home/mcbalaam/rss-notes:/data:ro"
    restart: unless-stopped
```

```sh
docker compose up -d
```

> Replace `/home/mcbalaam/rss-notes` with your notes dir. All configuration goes in `compose.yml`, but you can use `.env` if you prefer to.

### From Docker Hub directly

```sh
docker pull mcbalaam/pico-rss:latest

docker run -d --name pico-rss --restart unless-stopped \
  -p 8066:8066 \
  -v /home/user/rss-notes:/data:ro \
  -e RSS_TARGET_DIR=/data \
  -e RSS_ORIGIN_URL=https://rss.mcblm.xyz \
  -e RSS_FEED_TITLE="mcbalaam notes" \
  -e RSS_FEED_DESCRIPTION="Notes by mcbalaam" \
  -e RSS_AUTHOR_NAME=mcbalaam \
  -e RSS_AUTHOR_EMAIL=mcbalaam@example.com \
  -e RSS_TIMEOUT=10 \
  mcbalaam/pico-rss:latest
```

## Configuring pico-rss

- `RSS_FEED_TITLE`: feed title, fallback `"<author> notes"` (`author` from `RSS_AUTHOR_NAME`);
- `RSS_FEED_DESCRIPTION`: feed description, fallback `"Notes by <author>"`;
- `RSS_AUTHOR_NAME`: author name;
- `RSS_AUTHOR_EMAIL`: author email, empty by default;
- `RSS_ORIGIN_URL` is the URL used to host your feed;
- `RSS_TARGET_DIR`: where will pico-rss look for the `.md` files (`/data` inside container, mapped via `volumes`);
- `RSS_TIMEOUT`: how long will the server retry writing/reading for.

## Using pico-rss
### Publishing

pico-rss parses the `.md` files to generate the feed. Here's how to structure yours to get the desired output:

```md
# Hello, ants!    ⟸ this is the title of the post, should be <h1> (# header)
22.09.2026    ⟸ this is the publishing date, should be dd.mm.yyyy
Short teaser shown in readers as description.    ⟸ this is the description, should be under 300 symbols

Everything below is the full post body, rendered to HTML
into `content:encoded`. **Bold** and *italic* work.
```

If those are unavailable, metadata is used as a fallback: file name as the title, modification date as the publishing date.

### Reading

Available endpoints:
- `GET /rss.xml`: application/rss+xml;
- `GET /atom.xml`: application/atom+xml;
- `GET /feed.json`: application/feed+json;
- `GET /{slug}`: full post as HTML, e.g. `GET /hello-ants` for `hello-ants.md`;
- `GET /`: index with all posts and auto-discovery.

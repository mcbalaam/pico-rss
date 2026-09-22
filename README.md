# pico-rss: an RSS server for ants

pico-rss is a very minimal personal RSS feed server that serves `.md` files using [gorilla/feeds](https://github.com/gorilla/feeds) with a Dockerhub release at [mcbalaam/pico-rss](https://hub.docker.com/r/mcbalaam/pico-rss)

## Running pico-rss

You can...
- clone the project, reconfigure [compose.yml](https://github.com/mcbalaam/pico-rss/blob/master/compose.yml) and run `docker compose up -d`;
- create and configure an `.env` file and  use the following `compose.yml` to pull the image from Dockerhub;
```
services:
  go-rss:
    image: mcbalaam/pico-rss:latest
    ports:
      - "${RSS_PORT:-8080}:8080"
    volumes:
      - "${NOTES_DIR}:/data:ro"
```
- pull it directly from Dockerhub:

```sh
docker pull mcbalaam/pico-rss:latest

docker run -d --name pico-rss --restart unless-stopped \
  -p 8066:8066 \
  -v /home/user/rss-notes:/data:ro \
  -e RSS_TARGET_DIR=/data \
  -e RSS_HOST=0.0.0.0 \
  -e RSS_PORT=8066 \
  -e RSS_ORIGIN_URL=https://rss.mcblm.xyz \
  -e RSS_AUTHOR_USERNAME=mcbalaam \
  -e RSS_TIMEOUT=10 \
  mcbalaam/pico-rss:latest
```

## Configuring pico-rss

- `RSS_AUTHOR_USERNAME` is what people will see as the author of all the posts served;
- `RSS_ORIGIN_URL` is the URL used to host your feed;
- `RSS_TARGET_DIR`: where will pico-rss look for the `.md` files;
- `RSS_HOST`: for running pico-rss outside of containers on localhost or local network;
- `RSS_PORT`: the used port;
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

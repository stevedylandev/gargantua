# Gargantua 

// COVER

> This little maneuver is going to cost us years

A simple scraper tarpit and maze written in Go


## Quick Start 

Clone the repo and install

```bash
git clone https://github.com/stevedylandev/gargantua
cd gargantua
go install .
```

Run the server

```bash
gargantua
```

## Configuration 

Use the following environment variables to configure Gargantua

- `GARGANTUA_CORPUS` - Text used for Markov babbler (default looks for `moby.txt` in the repo)
- `GARGANTUA_HOST`- Host used (default `127.0.0.1`)
- `GARGANTUA_PORT` - Port used (default `8090`)
- Gargantua uses a slow drip method for writing the html that will keep the scraper longer. You can tweak the range of delay used between writes and the amount of data written.
  - `GARGANTUA_DELAY_MIN_MS` - Minimum delay in ms (default `150`)
  - `GARGANTUA_DELAY_MAX_MS` - Maximum delay in ms (default `250`)
  - `GARGANTUA_CHUNK_MAX_BYTES` - Maximum number of bytes in a written chunk (default `45)
- `robots.txt` - Defaults to `User-agent: * Disallow: /` to prevent indexing; the goal i to direct a disobedient crawler to this deployed server and get it stuck.


## Reasoning

LLM scrapers love to eat content, so this gives them something to feed on. Heavily inspired by Nepenthes, this sever will create an endless loop of content and relative links that are dripped slowly to consume resources and time. Instead of writing a proxy to go on top of a site, the idea of Gargantua is to have your personal site `robots.txt` say something like: 

```
User-agent: *
Disallow: /gargantua
```

Then have that `/gargantua` path do a permanent redirect to a server running the program. Crawlers are supposed to obey the `robots.txt`, and by knowingly disregarding that rule, they are subject to Gargantua. 

I also wanted to write this in Go instead of Lua to make it more resource efficient. Go really is amazing for this kind of project, as it uses goroutines by default to enable multithreaded performance, and the design of the server is to keep each crawl within its own goroutine. The server can handle upwards of 1000 connections at once and not break a sweat, making it ideal to run on smaller hardware like a Pi. 

# Matrix schema server

Simple and slim server to answer `well-known` requests for [matrix.org](https://matrix.org) servers based on the request
host header, to be used behind a reverse proxy of your choice.

# Setup example

For my own domain `klonmaschine.de` I use the following docker-compose setup:

```
version: '3.7'

services:
  updater:
    image: adrianrudnik/maitrx-schema-server
    volumes:
      - ./schema:/var/schema
    ports:
      - 8080:8080
```

Now let's define the response to `/.well-known/matrix/server` if a request with an `Host: klonmaschine.de` header
is recieved by placing a file into `./schema/klonmaschine.de.server.json`:

```json
{
  "m.server": "matrix.klonmaschine:de:443"
}
```

The same for `/.well-known/matrix/client` requests, placed into `./schema/klonmaschine.de.client.json`:

```json
{
  "m.homeserver": {
    "base_url": "https://matrix.klonmaschine.de"
  }
}
```

Now boot it up with `docker-compose up -d` and the following requests are answered correctly and offer CORS support as well:

http://klonmaschine.de:8080/.well-known/matrix/server
http://klonmaschine.de:8080/.well-known/matrix/client

# Environment variables

The following configuration can be done by passing environment variables to this app:

| Key | Description | Default |
| --- | --- |
| BIND_ADDR | HTTP bind address | `:8080` |
| SCHEMA_ROOT | Root folder for schema JSON files | `/var/schema` |

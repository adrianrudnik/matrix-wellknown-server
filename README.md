# Matrix schema server

Simple and slim server to answer `well-known` requests for [matrix.org](https://matrix.org) servers based on the request
host header, to be used behind a reverse proxy of your choice.

Features:

- Have a single server instance serve as many different domains on the same reverse proxy as required.
- Simple managment by passing through your JSON files queried by the `Host`-header.
- Low footprint (~8MB RAM, 10MB binary), no heavy dependencies (single binary), fast responses (~1ms).

# Setup example

For my own domain `klonmaschine.de` I use the following docker-compose setup:

```
version: '3.7'

services:
  updater:
    image: adrianrudnik/matrix-wellknown-server:latest
    volumes:
      - ./schema:/var/schema
    ports:
      - 8080:8080
```

Now let's define the response to `/.well-known/matrix/server` if a request with an `Host: klonmaschine.de` header is
received by placing a file into `./schema/klonmaschine.de.server.json`:

```json
{
  "m.server": "matrix.klonmaschine.de:8448"
}
```

The same for `/.well-known/matrix/client` requests, placed into `./schema/klonmaschine.de.client.json`:

```json
{
  "m.homeserver": {
    "base_url": "https://matrix.klonmaschine.de:8448"
  }
}
```

Now boot it up with `docker-compose up -d` and the following requests are answered correctly and offer CORS support as
well:

https://klonmaschine.de/.well-known/matrix/server  
http://klonmaschine.de/.well-known/matrix/client

# Environment variables

The following configuration can be done by passing environment variables to this app:

| Key | Description | Default |
| --- | --- | --- |
| BIND_ADDR | HTTP bind address | `:8080` |
| SCHEMA_ROOT | Root folder for schema JSON files | `/var/schema` |

# Remarks

The element.io client does not use the server schema to resolve the homeserver. You need to specify the matrix server as
full domain, in my example `matrix.klonmaschine.de` to connect to it.

# Reserver proxy examples

## Traefik

I use this slightly redacted version to serve the app:

```yaml
http:
  routers:
    matrix-wellknown-server:
      rule: "PathPrefix(`/.well-known/matrix/`)"
      priority: 1000
      entryPoints:
        - https
      service: wellknown-server
      tls:
        certResolver: http-le
        
  services:
    wellknown-server:
      loadBalancer:
        servers:
          - url: "http://wellknown-server:8080"
```

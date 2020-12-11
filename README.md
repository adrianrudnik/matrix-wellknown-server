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

## Traefik v2.2

I use this slightly redacted version to serve the app:

```yaml
http:
  routers:
    matrix.klonmaschine.de-schema:
      rule: "Host(`klonmaschine.de`) && PathPrefix(`/.well-known/matrix/`)"
      entryPoints:
        - https
      service: matrix.klonmaschine.de-schema
      tls:
        certResolver: http-le

    matrix.klonmaschine.de-server:
      rule: "Host(`matrix.klonmaschine.de`)"
      entryPoints:
        - https
        - matrixfed
      service: matrix.klonmaschine.de-server
      tls:
        certResolver: http-le

  services:
    matrix.klonmaschine.de-schema:
      loadBalancer:
        servers:
          - url: "http://schemaserver:8080"

    matrix.klonmaschine.de-server:
      loadBalancer:
        servers:
          - url: "http://matrixserver:8008"
```

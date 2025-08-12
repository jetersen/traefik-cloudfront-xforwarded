# CloudFront To X-Forwarded Headers - Traefik Plugin

A Traefik middleware plugin that extracts AWS CloudFront headers and translates them into X-Forwarded-* headers.
It is useful for applications behind CloudFront that need to be aware of the original client IP and protocol.

This also preserves the remote port information on the X-Forwarded-For header.
Which can be important for applications that need to know the original port used by the client.

## Features

- Extracts AWS CloudFront headers and translates them into X-Forwarded-* headers.
- Preserves the remote port information on the X-Forwarded-For header.

## Static Configuration

```yaml
# Static Configuration
experimental:
  plugins:
    cloudfront-to-xforwarded:
      moduleName: github.com/jetersen/traefik-cloudfront-xforwarded
      version: v0.1.0
```

```yaml
# Dynamic Configuration
http:
  middlewares:
    cloudfront-to-xforwarded:
      plugin:
        cloudfront-to-xforwarded: {}

  routers:
    my-router:
      rule: "Host(`example.com`)"
      service: my-service
      middlewares:
        - cloudfront-to-xforwarded
      entryPoints:
        - websecure

  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://internal-service:8080"
```

### Configuration Options

N/A

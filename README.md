# Traefik Aggregator

This is a simple tool to aggregate multiple traefik instances into one.
The use case is for an ingress proxy in e.g. a different machine than the actual traefik instances.

## Usage

### Setup as Local Plugin

Clone the repository

```bash
git clone https://github.com/MadddinTribleD/traefikaggregator traefikaggregator
```

The plugin needs to be located in the `plugins-local/src/github.com/MadddinTribleD/traefikaggregator` relative to the traefik binary.  
For docker-compose the volume mount looks like this:

```yml
    volumes:
      - ./traefikaggregator:/plugins-local/src/github.com/MadddinTribleD/traefikaggregator
```

### Configuration

Add the plugin to your `traefik.yml`:

```yml
providers:
  plugin:
    traefikaggregator:
      pollInterval: 2s
      instances:
        - apiEndpoint: http://<traefikInstance_1>:8080/api
          service:
            name: traefikInstance-1-web
            loadBalancer:
              servers:
                - url: http://<traefikInstance_1>:80
          allowedEndpoints:
            - web
          router:
            entryPoints:
              - web
            middlewares:
              - redirect-to-https@file
        - apiEndpoint: http://<traefikInstance_1>:8080/api
          service:
            name: traefikInstance-1-web-secured
            loadBalancer:
              servers:
                - url: https://<traefikInstance_1>:443
          allowedEndpoints:
            - web-secured
          router:
            entryPoints:
              - web-secured
            tls:
              enabled: true
          certResolverMapping:
            <remoteTlsResolve1>: <localTlsResolve1>
            <remoteTlsResolve2>: <localTlsResolve2>
```

It is also possible to proxy the remote traefik dashboard:

```yml
        - apiEndpoint: http://<traefikInstance_1>/api
          service:
            name: traefikInstance-1-traefik
            loadBalancer:
              servers:
                - url: http://<traefikInstance_1>:8080
          allowedEndpoints:
            - traefik
          router:
            entryPoints:
              - traefik
            middlewares:
              - auth-for-dashboard@file
```

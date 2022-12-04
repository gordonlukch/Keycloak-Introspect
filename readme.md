### Configuration

```yaml
# Static configuration

experimental:
  plugins:
    KcIntrospect:
      moduleName: github.com/traefik/Keycloak-Introspect
      version: v0.1.0
```

```yaml
# Dynamic configuration

http:
  routers:
    ......
    ...
      middlewares:
        - my-plugin
    ......
    ...

  services:
    ......
    ...

  middlewares:
    my-plugin:
      plugin:
        KcIntrospect:
        	- hostname: "keyclock hostname"
	        - client_id: "xxxxxxxxxx"
	        - client_secret: "xxxxxxxxxx"
	        - realm: "xxxxxxxxxx"
```

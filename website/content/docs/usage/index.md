---
title: Usage
weight: 2
---

`kubectl-power-forward` is used with a configuration file. Per default, it searches for a file `.power-forward.yaml` in the current directory. Use can use the `--config` (short `-f`) flag to reference a different file.

## Forwarding to Kubernetes
Forwarding to Kubernetes works with targeting a Kubernetes *service*. You need the following information to describe a service configuration:

- **namespace** - The namespace the service is deployed in.
- **serviceName** - The name of the service to be forwarded to.
- **podPort** - The port on the pod (not service!) to be forwarded to.
- **localPort** - The local port used to forward to Kubernetes.

You can define multiple forwards:

```yaml
forwards:
  - namespace: default
    serviceName: echoserver
    podPort: 5678
    localPort: 8080
  - namespace: production
    serviceName: very-important-app
    podPort: 80
    localPort: 8081
```

If a connection is lost (rescheduling of a pod), `kubectl-powre-forward` will reconnect to a different pod immediately.

## Adding a SOCKS proxy
You can also add a SOCKS proxy to allow using custom domains locally. You have to explicitly enable it via the `proxy.socks.enabled` parameter. To define a domain that should be intercepted and pointed to `localhost`, add a `hostname` to your forwards:

```yaml
proxy:
  socks:
    enabled: true
    listenAddress: "0.0.0.0:1080" # default
forwards:
  - namespace: default
    serviceName: echoserver
    podPort: 5678
    localPort: 8080
    hostname: echo.example.org
```

Now you can access the forwarded server with the custom URL

```shell
$ curl --socks5-hostname 127.0.0.1:1080 http://echo.example.org:8080
```

## Adding a HTTP reverse proxy
You can add a HTTP server to use custom domains with real-world ports. As you already set the hostnames with the SOCKS proxy configuration, you only have to enable the reverse proxy:


```yaml
proxy:
  socks:
    enabled: true
  httpReverse:
    enabled: true
    listenAddress: "0.0.0.0:80" # default, requires NET_BIND_SERVICE capability
forwards:
  - namespace: default
    serviceName: echoserver
    podPort: 5678
    localPort: 8081
    hostname: echo.example.org

```

{{% alert title="Use different ports for services and the proxy" %}}
Your service will be forwarded from a local port in all cases. You cannot bind the HTTP reverse proxy to the same port as a local port of a forwarded service. When caring about custom domain and port, set the HTTP proxy's port as the relevant one and give the local forward a different one. It will be discovered and set up automatically.
{{% /alert %}}

This allows you to query you workload with the given domain and port

```shell
$ curl --socks5-hostname 127.0.0.1:1080 http://echo.example.org
```

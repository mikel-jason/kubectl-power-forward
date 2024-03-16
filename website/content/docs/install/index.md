---
title: Installation
weight: 1
---

## Introduction
This installation covers installing the `kubectl-power-forward` binary to your system. You will need some kind of client to use it, like a browser or application. The later is not part of this documentation.

## System requirements
`kubectl-power-forward` is built to run on Linux and Mac, both with amd64 and arm64 (e.g. Apple Silicon) CPUs.

## Installation / Update steps
The following procedure explains how to install `kubectl-power-forward` as binary.

{{% alert title="No stable release yet" color="warning" %}}
As the application is under initial development, there is no stable version to install yet! Follow along to install a pre-release version.
{{% /alert %}}


1. Go to [the release page](https://github.com/mikel-jason/kubectl-power-forward/releases), choose your release and download the archive matching you system
2. Unpack the archive, e.g. `tar -xf <PAHT_TO_ARCHIVE>`
3. Place the binary from the archive to a localtion in your `PATH` variable, e.g. `/usr/local/bin`


### Optional: Allow binding privileged ports
If you want to use `kubectl-power-forward` with ports < 1024 (e.g. 80 or 443 for easier usage with HTTP traffic), you can either use `sudo` mode or add the `NET_BIND_SERVICE` capability to the binary. You can do this with

```shell
$ sudo setcap 'cap_net_bind_service+ep' "$(which kubectl-power-forward)"
```


## Uninstalling
Remove the binary with

```shell
$ rm $(which kubectl-power-forward)
```


## Next steps
Read the [Usage]({{< ref "/docs/usage/index.md" >}}) to learn how to set up forwarding.

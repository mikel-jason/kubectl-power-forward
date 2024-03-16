---
title: Clients
weight: 3
---

# Browser: Firefox
Follow [this guide](https://support.mozilla.org/en-US/kb/connection-settings-firefox). Make sure to use *Proxy DNS when using SOCKS5 v5* to be able to use custom domains.

# CLI: `curl`
Use the `--socks5-hostname` option and point it to your SOCKS proxy address. Check all options with `curl --help all | grep -E 'socks|proxy'`.

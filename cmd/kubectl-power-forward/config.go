package main

import "github.com/mikel-jason/kube-power-forward/pkg/cmd/powerforward"

type ConfigFile struct {
	Proxy    ProxyConfig            `json:"proxy"`
	Forwards []powerforward.Forward `json:"forwards"`
}

type ProxyConfig struct {
	Socks       SocksProxyConfig       `json:"socks"`
	HttpReverse HttpReverseProxyConfig `json:"httpReverse"`
}

type SocksProxyConfig struct {
	Enabled    bool   `json:"enabled"`
	ListenAddr string `json:"listenAddress"`
}

type HttpReverseProxyConfig struct {
	Enabled    bool   `json:"enabled"`
	ListenAddr string `json:"listenAddress"`
}

func (c *ConfigFile) FillDefaults() {
	if c.Proxy.HttpReverse.ListenAddr == "" {
		c.Proxy.HttpReverse.ListenAddr = "0.0.0.0:80"
	}
	if c.Proxy.Socks.ListenAddr == "" {
		c.Proxy.Socks.ListenAddr = "0.0.0.0:1080"
	}
}

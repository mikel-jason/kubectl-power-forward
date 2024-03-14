package proxy

import (
	"context"
	"log"
	"net"
)

type resolver struct {
	localNames []string
}

func (r *resolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	for _, localName := range r.localNames {
		if localName == name {
			log.Printf("intercepting '%s' -> localhost\n", name)
			return ctx, net.ParseIP("127.0.0.1"), nil
		}
	}

	// fallback to usual DNS
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, addr.IP, nil
}

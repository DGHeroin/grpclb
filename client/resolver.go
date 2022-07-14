package client

import (
    "google.golang.org/grpc/resolver"
    "log"
)

type (
    Resolver struct {
        cc       resolver.ClientConn
        discover Discover
        scheme   string
    }
)

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (r *Resolver) Close() {

}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
    r.cc = cc
    go r.watch()
    return r, nil
}

func (r *Resolver) Scheme() string {
    return r.scheme
}

func (r *Resolver) watch() {
    for {
        select {
        case addrs := <-r.discover.WatchUpdate():
            var addrList []resolver.Address
            for _, addr := range addrs {
                addrList = append(addrList, resolver.Address{
                    Addr: addr,
                })
            }
            if err := r.cc.UpdateState(resolver.State{Addresses: addrList}); err != nil {
                log.Println(err)
            }
        }
    }
}

func newResolver(discover Discover, scheme string) resolver.Builder {
    return &Resolver{
        scheme:   scheme,
        discover: discover,
    }
}

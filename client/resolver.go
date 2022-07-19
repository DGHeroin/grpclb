package client

import (
    "google.golang.org/grpc/resolver"
    "log"
    "sync"
)

type (
    Resolver struct {
        clientConn resolver.ClientConn
        discover   Discover
        scheme     string
        closeCh    chan struct{}
        once       sync.Once
    }
)

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (r *Resolver) Close() {

}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
    r.clientConn = cc
    r.once.Do(func() {
        go r.watch()
    })
    r.updateAddress(r.discover.Latest())
    return r, nil
}

func (r *Resolver) Scheme() string {
    return r.scheme
}

func (r *Resolver) watch() {
    for {
        select {
        case address := <-r.discover.WatchUpdate():
            r.updateAddress(address)
        case <-r.closeCh:
            return
        }
    }
}
func (r *Resolver) updateAddress(address []string) {
    var addrList []resolver.Address
    for _, addr := range address {
        addrList = append(addrList, resolver.Address{
            Addr: addr,
        })
    }
    if err := r.clientConn.UpdateState(resolver.State{Addresses: addrList}); err != nil {
        log.Println("UpdateState:", err)
    }
}

func newResolver(discover Discover) *Resolver {
    r := &Resolver{
        discover: discover,
        scheme:   "dns",
        closeCh:  make(chan struct{}),
    }

    return r
}

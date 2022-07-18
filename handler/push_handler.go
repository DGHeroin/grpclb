package handler

import (
    "context"
    "fmt"
)

type (
    PushHandlers struct {
        handler map[string]PushHandlerFunc
    }
    PushHandlerFunc func(ctx context.Context, payload []byte) ([]byte, error)
)

func NewPushHandlers() *PushHandlers {
    return &PushHandlers{
        handler: make(map[string]PushHandlerFunc),
    }
}

func (s *PushHandlers) Register(serviceName string, fn PushHandlerFunc) bool {
    s.handler[serviceName] = fn
    return true
}

func (s *PushHandlers) HandleMessage(ctx context.Context, name string, payload []byte) ([]byte, error) {
    fn, ok := s.handler[name]
    if !ok {
        return nil, fmt.Errorf("push handler [%v] not found", name)
    }
    return fn(ctx, payload)
}

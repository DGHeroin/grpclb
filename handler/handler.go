package handler

import (
    "context"
    "fmt"
    "github.com/DGHeroin/grpclb/utils"
    "reflect"
)

type (
    Handlers struct {
        handler map[string]*handler
    }
    handler struct {
        fn reflect.Value
        r  reflect.Type
        w  reflect.Type
    }
)

func NewHandlers() *Handlers {
    return &Handlers{
        handler: make(map[string]*handler),
    }
}

func (s *Handlers) Register(serviceName string, i interface{}) bool {
    sh, ok := checkFunc(i)
    if !ok {
        return false
    }
    s.handler[serviceName] = sh
    return true
}

func (s *Handlers) HandleMessage(ctx context.Context, name string, payload []byte) ([]byte, error) {
    if ctx == nil {
        ctx = context.Background()
    }
    defer utils.Recover()
    sh, ok := s.handler[name]
    if !ok {
        return nil, fmt.Errorf("handler[%s] not found", name)
    }

    t0 := reflect.ValueOf(ctx)
    t1 := reflect.New(sh.r)
    t2 := reflect.New(sh.w)
    in := []reflect.Value{
        t0, t1, t2,
    }

    err := Unmarshal(payload, t1.Interface())
    if err != nil {
        return nil, fmt.Errorf("handler Unmarshal error:%v", err)
    }

    rs := sh.fn.Call(in)
    r1 := rs[0]
    if r1.Interface() != nil {
        return nil, fmt.Errorf("handler invoke error:%v", err)
    }
    data, err := Marshal(in[2].Interface())
    if err != nil {
        return nil, fmt.Errorf("handler Marshal error:%v", err)
    }

    return data, nil
}

func checkFunc(fn interface{}) (sh *handler, ok bool) {
    defer utils.Recover()
    // 检查传入的函数是否符合格式要求
    var (
        typeOfError = reflect.TypeOf((*error)(nil)).Elem()
    )
    f, ok := fn.(reflect.Value)
    if !ok {
        f = reflect.ValueOf(fn)
    }

    t := f.Type()
    if t.NumIn() != 3 { // context/request/response
        return
    }
    if t.NumOut() != 1 {
        return
    }
    if returnType := t.Out(0); returnType != typeOfError {
        return
    }
    r := t.In(1)
    w := t.In(2)

    ok = true
    sh = &handler{
        fn: f,
        r:  r.Elem(),
        w:  w.Elem(),
    }
    return
}

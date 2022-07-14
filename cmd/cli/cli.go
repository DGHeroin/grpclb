package main

import (
    "bufio"
    "fmt"
    "github.com/DGHeroin/grpclb/client"
    "log"
    "os"
    "strings"
    "time"
)

type discovery struct {
    ch chan []string
}

func (d *discovery) WatchUpdate() <-chan []string {
    return d.ch
}

func main() {
    dis := &discovery{
        ch: make(chan []string),
    }
    go func() {
        dis.ch <- []string{"localhost:30001", "localhost:30002"}
    }()
    go func() {
        r := bufio.NewReader(os.Stdin)
        for {
            line, _, err := r.ReadLine()
            if err == nil {
                str := strings.TrimSpace(string(line))
                if str == "" {
                    continue
                }
                infos := strings.Split(str, ",")
                log.Println(infos)
                dis.ch <- infos
            }
        }
    }()

    cli := client.NewClient(dis, client.WithPushMessage(func(name string, payload []byte) error {
        log.Println("收到推送哎")
        return nil
    }))

    request(cli, time.Second)
}
func request(cli client.Client, duration time.Duration) {
    i := 0
    for {
        i++
        if duration > 0 {
            time.Sleep(duration)
        }
        _, err := cli.Send(nil, "hello", []byte(fmt.Sprintf("hello world:%d", i)))
        if err != nil {
            continue
        }
    }
}

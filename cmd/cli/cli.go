package main

import (
    "bufio"
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
    // go func() {
    //     dis.ch <- []string{"localhost:30001", "localhost:30002"}
    // }()
    go func() {
        r := bufio.NewReader(os.Stdin)
        for {
            line, _, err := r.ReadLine()
            if err == nil {
                str := strings.TrimSpace(string(line))
                // if str == "" {
                //     continue
                // }
                infos := strings.Split(str, ",")
                log.Println(infos)
                dis.ch <- infos
            }
        }
    }()
    cli := client.NewClient(dis)
    request(cli, time.Second)
    // for i := 0; i < 1; i++ {
    //     go request(cli, time.Second)
    // }
    // select {}
}
func request(cli client.Client, duration time.Duration) {
    for {
        if duration > 0 {
            time.Sleep(duration)
        }
        _, err := cli.Send(nil, "hello", []byte("hello world"))
        if err != nil {
            continue
        }

    }
}

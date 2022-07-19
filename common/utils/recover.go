package utils

import (
    "bytes"
    "fmt"
    "log"
    "runtime"
)

func Recover() {
    if e := recover(); e != nil {
        buffer := bytes.NewBufferString(fmt.Sprint(e))
        // 打印调用栈信息
        buf := make([]byte, 2048)
        n := runtime.Stack(buf, false)
        stackInfo := fmt.Sprintf("\n%s", buf[:n])
        buffer.WriteString(fmt.Sprintf("panic stack info %s", stackInfo))
        log.Println(buffer)
    }
}

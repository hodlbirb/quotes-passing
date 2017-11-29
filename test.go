package main

import (
    "fmt"
    "reflect"
)

func main() {
    x := uint64(1)
    px := &x
    fmt.Printf("%d\n", px == do(&x))
    fmt.Printf("[x] T:%T S:%d V:%d\n", x, reflect.TypeOf(x).Size(), *do(&x))
}

func do(p *uint64) *uint64 {
    v := uint64(20)
    p = &v
    return p
}
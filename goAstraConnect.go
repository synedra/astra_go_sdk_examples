package main

import (
    "fmt"
    "github.com/datastax-ext/astra-go-sdk"
    "os"
)

func main() {
    token := os.Getenv("ASTRA_TOKEN")
    secureBundle := os.Getenv("ASTRA_SECURE_BUNDLE_LOCATION")
    keyspace := os.Getenv("ASTRA_KEYSPACE")

    fmt.Println("Building client connection")

    c, err := astra.NewStaticTokenClient(
        token,
        astra.WithSecureConnectBundle(secureBundle),
	      astra.WithDefaultKeyspace(keyspace),
    )
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("SELECTing from system.local")

    rows, err := c.Query("SELECT cluster_name FROM system.local").Exec()
    if err != nil {
        fmt.Println(err)
    }

    for _, r := range rows {
        vals := r.Values()
        strClusterName := vals[0].(string)
        fmt.Println("cluster_name:", strClusterName)
    }
}

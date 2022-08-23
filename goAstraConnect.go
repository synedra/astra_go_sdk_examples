package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "github.com/datastax-ext/astra-go-sdk"
    "io/ioutil"
    "log"
    "os"
    "time"
)

func main() {
    astraURI := os.Getenv("ASTRA_URI_w_PORT")
    token := os.Getenv("ASTRA_TOKEN")
    secureBundleDir := os.Getenv("ASTRA_SECURE_BUNDLE_DIR")

    certPath := secureBundleDir + "cert"
    keyPath := secureBundleDir + "key"
    caPath := secureBundleDir + "ca.crt"

    //keyspace := "stackoverflow"

    cert, err := tls.LoadX509KeyPair(certPath, keyPath)
    if err != nil {
     log.Fatalf("failed to load key pair: %s", err)
    }
    caCert, err := ioutil.ReadFile(caPath)
    if err != nil {
     log.Fatalf("failed to read CA cert: %s", err)
    }
    caCertPool, err := x509.SystemCertPool()
    //caCertPool, err := x509.NewCertPool()
    if err != nil {
     log.Fatalf("failed to read system cert pool: %s", err)
    }
    if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
     log.Fatal("failed to append caCert")
    }

    fmt.Println("Defining TLS config")

    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
    }

    fmt.Println("Building client connection")

    c, err := astra.NewStaticTokenClient(
        // URL of the Stargate service to use.
        // Example: "<cluster ID>-<cluster region>.apps.astra.datastax.com:443"
        astraURI,
        // Static auth token to use.
        token,
        // Optional deadline for initial connection.
        astra.WithDeadline(time.Second * 10),
        // Optional per-query timeout.
        astra.WithTimeout(time.Second * 5),
        // Optional TLS config. Assumes insecure if not specified.
        astra.WithTLSConfig(tlsConfig),
        // Optional default keyspace in which to run queries that do not specify
        // keyspace.
        //astra.WithDefaultKeyspace(keyspace),
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

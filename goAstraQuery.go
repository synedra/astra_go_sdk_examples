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
    //"github.com/pborman/uuid"
)

func main() {

    var name string

    if len(os.Args) != 2 {
        fmt.Println("This program only takes one argument.")
        os.Exit(1);
    } else {
        name = os.Args[1]
    }

    astraURI := os.Getenv("ASTRA_URI_w_PORT")
    token := os.Getenv("ASTRA_TOKEN")
    secureBundleDir := os.Getenv("ASTRA_SECURE_BUNDLE_DIR")

    certPath := secureBundleDir + "cert"
    keyPath := secureBundleDir + "key"
    caPath := secureBundleDir + "ca.crt"

    keyspace := "stackoverflow"

    // TLS config
    fmt.Println("Defining TLS config")

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

    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
    }

    fmt.Println("Building client connection")

    c, err := astra.NewStaticTokenClient(
        // URL of the Stargate service to use.
        // Example: "localhost:8090"
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
        astra.WithDefaultKeyspace(keyspace),
    )
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("SELECTing from stackoverflow.user_offers")

//    type UserOffers struct {
//    	userId    string
//    	group_id   uuid.UUID
//      offer_id   uuid.UUID
//      offer string
//    }

    var userId, offer string
    var groupId, offerId string

    rows, err := c.Query("SELECT user_id, group_id, offer_id, offer_desc FROM stackoverflow.user_offers WHERE user_id=?", name).Exec()
    if err != nil {
        fmt.Println(err)
    }

    for _, row := range rows {
        //uoffers := &UserOffers{}
        err := row.Scan(&userId, &groupId, &offerId, &offer)

        if err != nil {
      		log.Fatalf("failed to scan row: %v", err)
      	}

        fmt.Printf("%s - %s - %s - %s \n", userId, groupId, offerId, offer)
    }
}

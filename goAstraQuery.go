package main

import (
    "fmt"
    "github.com/datastax-ext/astra-go-sdk"
    "log"
    "os"
)

func main() {

    var name string

    if len(os.Args) != 2 {
        fmt.Println("This program only takes one argument.")
        os.Exit(1);
    } else {
        name = os.Args[1]
    }

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

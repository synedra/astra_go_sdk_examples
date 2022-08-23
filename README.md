# Astra Go SDK Examples

These are simple Go programs which demonstrate how to use DataStax's Astra Go SDK.

## Prerequisites

 - A DataStax [Astra DB](https://astra.datastax.com) account.
 - Go 1.17 or 1.18 installed locally.

Verify that Go is installed with `go version`.

```bash
% go version
go version go1.18.2 darwin/amd64
```

## Environment Variables

These programs each require three environment variables to be set.

```
ASTRA_URI_w_PORT = "<DATABASEID>-<REGION>.apps.astra.datastax.com:443"
ASTRA_TOKEN = "<TOKEN>"
ASTRA_SECURE_BUNDLE_DIR = "<LOCATION OF UNZIPPED SECURE BUNDLE>"
```

The SDK connects with gRPC via the StarGate API, so the port should always be `443`.  The Astra token will start with a "AstraCS" prefix, which should also be included in the variable definition.

For the secure bundle, download it from the "Connect" tab on the Astra dashboard.  By default it will be in your `~/Downloads` directory.  Feel free to move it to another location, and then unzip it.  It must be unzipped (for now) so that the program can get at the TLS certificate files (`key`, `cert`, and `ca.crt`).  Example:

```
export ASTRA_SECURE_BUNDLE_DIR = "/Users/aaronploetz/local/stackoverflow/"
```

To assist with this part of the setup, copy the `sample.env` file to `.env` and edit `.env` with your Astra information.

```bash
% cp sample.env .env
```

Then simply `source` the file to instantiate the environment variables.

```bash
% source .env
```

## goAstraConnect.go

This program uses the SDK to connect to Astra DB, and returns the value of `cluster_name` from the `system.local` table.  With DataStax Astra DB, this should always be "cndb."

To run:

```bash
% go run goAstraConnect.go
Defining TLS config
Building client connection
SELECTing from system.local
cluster_name: cndb
```

If the cluster name is returned successfully, then the environment variables were correctly supplied.

## goAstraQuery.go

To get this program to work, you'll need an additional table and keyspace set up on your Astra DB cluster.  From the Astra Dashboard, click on the database name, and then click the "Add Keyspace" button.  For the purposes of this example, the keyspace is named "stackoverflow."

With the keyspace created, click on the "CQL Console" tab and run the following CQL:

```SQL
use stackoverflow;
CREATE TABLE stackoverflow.user_offers (
    user_id text,
    group_id timeuuid,
    offer_id timeuuid,
    offer_desc text,
    PRIMARY KEY (user_id, group_id, offer_id)
) WITH CLUSTERING ORDER BY (group_id DESC, offer_id DESC);
```

Then `INSERT` some data:

```SQL
INSERT INTO user_offers (user_id,group_id,offer_id,offer_desc)
VALUES ('Aaron',now(),now(),'Free stuff!');
INSERT INTO user_offers (user_id,group_id,offer_id,offer_desc)
VALUES ('Aaron',c84f7dc0-b5ba-11ec-bdcd-bdf5ae06f7ae,now(),'More free stuff!');
INSERT INTO user_offers (user_id,group_id,offer_id,offer_desc)
VALUES ('Aaron',now(),now(),'Amz $20 gift card.');
INSERT INTO user_offers (user_id,group_id,offer_id,offer_desc)
VALUES ('Aaron',2e753db0-b5bb-11ec-bdcd-bdf5ae06f7ae,now(),'Amz $50 gift card.');
INSERT INTO user_offers (user_id,group_id,offer_id,offer_desc)
VALUES ('Aaron',2e753db0-b5bb-11ec-bdcd-bdf5ae06f7ae,now(),'Amz $100 gift card.');
INSERT INTO user_offers (user_id,group_id,offer_id,offer_desc)
VALUES ('Patrick',now(),now(),'Free stuff!');
INSERT INTO user_offers (user_id,group_id,offer_id,offer_desc)
VALUES ('Patrick',now(),now(),'Amz $20 gift card.');
```

Then, running this program with the name "Aaron" or "Patrick" should retrieve results from the table.

```
% go run goAstraQuery.go Patrick
Defining TLS config
Building client connection
SELECTing from stackoverflow.user_offers
Patrick - 65054640-2325-11ed-9388-cb18a63d6f83 - 65054641-2325-11ed-9388-cb18a63d6f83 - Amz $20 gift card.
Patrick - d3a4cc20-b5ba-11ec-bdcd-bdf5ae06f7ae - d3a4cc21-b5ba-11ec-bdcd-bdf5ae06f7ae - Free stuff!
```

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"rhize-data-collection-import/auth"
	"rhize-data-collection-import/importSheet"
	"rhize-data-collection-import/types"
	"time"
)

const versionString = "rhize-data-collection-import v1.2.1"

var (
	bFile        = flag.String("file", "", "Excel/CSV file to import data from")
	bDescription = flag.String("description", "", "Equipment Class description used in import for CSV")
	bSheet       = flag.String("sheet", "", "Name of sheet to import data from")
	bDatasource  = flag.String("datasource", "", "Datasource to bind topics with")

	/* Auth */
	bAuth         = flag.Bool("auth", true, "Authenticate Client")
	sAuthUrl      = flag.String("authUrl", "http://localhost:8090", "URL for Keycloak Auth")
	sRealm        = flag.String("realm", "libre", "Keycloak Realm")
	sClientID     = flag.String("clientId", "libreBaas", "Client ID")
	sClientSecret = flag.String("clientSecret", "", "Client Secret")
	sUser         = flag.String("user", "admin", "Authentication Username")
	sPassword     = flag.String("password", "admin", "Authentication Password")
	sURL          = flag.String("apiUrl", "http://localhost:8080/graphql", "URL for Rhize API")
)

func init() {
	flag.Parse()
}

func main() {

	log.Println(versionString)

	ctx := context.Background()

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	var err error

	if *bAuth {
		client, err = auth.Authenticate(ctx, *sAuthUrl, *sUser, *sPassword, *sRealm, *sClientID, *sClientSecret)
		if err != nil {
			log.Fatalf("Authentication failed: %v", err)
			return
		}
	}

	// Setup Configuration
	config := types.Configuration{
		Context:                   ctx,
		Client:                    client,
		URL:                       sURL,
		FilePath:                  bFile,
		Sheet:                     bSheet,
		EquipmentClassDescription: bDescription,
		Datasource:                bDatasource,
	}

	log.Println("Importing sheet " + *bSheet + " from file " + *bFile)

	importSheet.Import(config)
}

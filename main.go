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

const versionString = "rhize-data-collection-import v1.0.0"

var (
	bFile        = flag.String("file", "", "Excel/CSV file to import data from")
	bDescription = flag.String("description", "", "Equipment Class description used in import for CSV")
	bSheet       = flag.String("sheet", "", "Name of sheet to import data from")

	/* Auth */
	bAuth         = flag.Bool("auth", true, "Authenticate Client")
	sAuthUrl      = flag.String("auth-host", "http://localhost:8090", "URL for Auth")
	sClientID     = flag.String("clientId", "libreBaas", "Client ID")
	sClientSecret = flag.String("clientSecret", "7a7m7IBF0THzyx7K7H4SsgwA5If1xYMO", "Client Secret")
	sUser         = flag.String("user", "admin", "Authentication Username")
	sPassword     = flag.String("password", "admin", "Authentication Password")
	sURL          = flag.String("url", "http://localhost:8080/graphql", "URL for Libre")
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

	if *bAuth {
		client = auth.Authenticate(ctx, *sAuthUrl, *sUser, *sPassword, *sClientID, *sClientSecret)
	}

	// Setup Configuration
	config := types.Configuration{
		Context:                   ctx,
		Client:                    client,
		URL:                       sURL,
		FilePath:                  bFile,
		Sheet:                     bSheet,
		EquipmentClassDescription: bDescription,
	}

	importSheet.Import(config)
}

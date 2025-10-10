package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"rhize-data-collection-import/auth"
	"rhize-data-collection-import/importSheet"
	"rhize-data-collection-import/types"
	"time"

	"github.com/joho/godotenv"
)

const versionString = "rhize-data-collection-import v1.2.1"

var (
	bFile        = flag.String("file", "", "Excel/CSV file to import data from")
	bDescription = flag.String("description", "", "Equipment Class description used in import for CSV")
	bSheet       = flag.String("sheet", "", "Name of sheet to import data from")
	bDatasource  = flag.String("datasource", "", "Datasource to bind topics with")
	bTarget      = flag.String("target", "", "Target type of import data")

	/* Auth */
	bAuth         = flag.Bool("auth", true, "Authenticate Client")
	sAuthUrl      = flag.String("authUrl", "http://localhost:8090", "URL for Keycloak Auth")
	sRealm        = flag.String("realm", "libre", "Keycloak Realm")
	sClientID     = flag.String("clientId", "libreBaas", "Client ID")
	sClientSecret = flag.String("clientSecret", "", "Client Secret")
	sUser         = flag.String("user", "", "Authentication Username")
	sPassword     = flag.String("password", "", "Authentication Password")
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

	// Attempt loading .env
	godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file, running without")
	}

	// Handle Client Secret, User, and Password

	clientSecret := *sClientSecret
	if clientSecret == "" {
		log.Printf("No Client Secret set, loading from .env")
		clientSecret = os.Getenv("RHIZE_OIDC_CLIENT_SECRET")
	}
	user := *sUser
	if user == "" {
		log.Printf("No User set, loading from .env")
		user = os.Getenv("RHIZE_OIDC_USER")
	}
	password := *sPassword
	if password == "" {
		log.Printf("No Password set, loading from .env")
		password = os.Getenv("RHIZE_OIDC_PASSWORD")
	}

	if *bAuth {
		client, err = auth.Authenticate(ctx, *sAuthUrl, user, password, *sRealm, *sClientID, clientSecret)
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
		Target:                    bTarget,
	}

	log.Println("Importing sheet " + *bSheet + " from file " + *bFile)

	importSheet.Import(config)
}

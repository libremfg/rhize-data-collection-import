package commands

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rhize-data-collection-import/adapters"
	"rhize-data-collection-import/adapters/drivers"
	"rhize-data-collection-import/auth"
	"rhize-data-collection-import/types"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Data
	file       string
	ImportData types.ImportData

	// Auth
	bypass       bool
	authUrl      string
	realm        string
	clientId     string
	clientSecret string
	apiUrl       string
	Client       *graphql.Client

	// Command
	RootCmd = &cobra.Command{
		Use:              "",
		Short:            "Rhize Data Collection Import",
		Long:             "Rhize Data Collection Import\n\nSimple utility to import data from a CSV or XLSX.",
		PersistentPreRun: handlePreRun,
		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {
	// Persistent Flags
	RootCmd.PersistentFlags().StringP("sheet", "S", "", "The Excel Sheet to search for data in")
	RootCmd.PersistentFlags().StringP("datasource", "D", "", "The DataSource to bind topics with")

	// File
	RootCmd.Flags().StringVarP(&file, "file", "f", "", "Excel/CSV file to import data from")

	// Authentication
	RootCmd.Flags().BoolVarP(&bypass, "bypass", "b", false, "Bypass Keycloak authentication")
	RootCmd.Flags().StringVarP(&authUrl, "authUrl", "a", "http://localhost:8090", "URL for Keycloak authentication")
	RootCmd.Flags().StringVarP(&realm, "realm", "r", "libre", "Keycloak Realm")
	RootCmd.Flags().StringVarP(&clientId, "clientId", "c", "libreBaas", "Client ID")
	RootCmd.Flags().StringVarP(&clientSecret, "clientSecret", "s", "", "Client Secret")

	// API
	RootCmd.Flags().StringVarP(&apiUrl, "apiUrl", "A", "http://localhost:8080/graphql", "URL for Rhize API")

	// Bind Flags
	viper.BindPFlag("sheet", RootCmd.PersistentFlags().Lookup("sheet"))
	viper.BindPFlag("datasource", RootCmd.PersistentFlags().Lookup("datasource"))

	// Commands
	RootCmd.AddCommand(UnitOfMeasureCmd)
	RootCmd.AddCommand(EquipmentClassCmd)
	RootCmd.AddCommand(EquipmentCmd)

	RootCmd.TraverseChildren = true
}

func handlePreRun(cmd *cobra.Command, args []string) {
	// Handle Client Secret, User, and Password
	if clientSecret == "" {
		log.Printf("No Client Secret set, loading from .env")
		clientSecret = os.Getenv("RHIZE_OIDC_CLIENT_SECRET")
	}

	// Client Setup
	ctx := context.Background()

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	var err error
	if !bypass {
		client, err = auth.Authenticate(ctx, authUrl, realm, clientId, clientSecret)
	}
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	Client = graphql.NewClient(apiUrl, client)

	// Read

	// Determine file type
	// First check that file is set
	if file == "" {
		log.Println("Cannot run without \"file\" set, see usage below: \n")
		return
	}

	var reader adapters.SheetReader

	switch filepath.Ext(file) {
	case ".csv":
		reader = drivers.CSVAdapter{
			Description: viper.GetString("sheet"),
			Datasource:  viper.GetString("datasource"),
		}
	case ".xlsx":
		reader = drivers.XLSXAdapter{
			Sheet:      viper.GetString("sheet"),
			Datasource: viper.GetString("datasource"),
		}
	case "":
		log.Fatalf("Provided file \"%s\" has no extension\n", file)
	default:
		log.Fatalf("Unsupported file type \"%s\"\n", filepath.Ext(file))
	}

	log.Println("Reading data from provided file")
	importData, err := reader.Read(file)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	ImportData = *importData

	// If UoM ID is not set, set it to DataType
	for i := range ImportData.EquipmentClass.Properties {
		if ImportData.EquipmentClass.Properties[i].UnitOfMeasure.ID == "" && ImportData.EquipmentClass.Properties[i].UnitOfMeasure.DataType != "" {
			ImportData.EquipmentClass.Properties[i].UnitOfMeasure.ID = ImportData.EquipmentClass.Properties[i].UnitOfMeasure.DataType
		}
	}
}

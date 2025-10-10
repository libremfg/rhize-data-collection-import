package importSheet

import (
	"log"
	"path/filepath"

	"rhize-data-collection-import/types"

	"github.com/hasura/go-graphql-client"
)

func Import(config types.Configuration) {
	client := graphql.NewClient(*config.URL, config.Client)

	// Determine file type
	var equipmentImportData *ImportData
	var err error

	switch filepath.Ext(*config.FilePath) {
	case ".csv":
		equipmentImportData, err = ReadCSV(*config.FilePath, *config.EquipmentClassDescription)
	case ".xlsx":
		equipmentImportData, err = ReadXLSX(*config.FilePath, *config.Sheet, *config.Datasource)
	default:
		log.Fatalf("Unsupported file extension \"%s\"\n", filepath.Ext(*config.FilePath))
	}
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// If UoM ID is not set, set it to DataType
	for i := range equipmentImportData.EquipmentClassImportData.EquipmentClassProperties {
		if equipmentImportData.EquipmentClassImportData.EquipmentClassProperties[i].UnitOfMeasure.ID == "" && equipmentImportData.EquipmentClassImportData.EquipmentClassProperties[i].UnitOfMeasure.DataType != "" {
			equipmentImportData.EquipmentClassImportData.EquipmentClassProperties[i].UnitOfMeasure.ID = equipmentImportData.EquipmentClassImportData.EquipmentClassProperties[i].UnitOfMeasure.DataType
		}
	}

	log.Println("Adding Imported Unit of Measures")
	UnitOfMeasure(config.Context, client, *&equipmentImportData.EquipmentClassImportData)
	log.Println("Done Imported Unit of Measures")

	log.Println("Adding Imported Equipment model")
	EquipmentModel(config.Context, client, *equipmentImportData)
	log.Println("Done Imported Equipment model")

	log.Println("Done Imported model")
}

// Import
type ImportData struct {
	Datasource               string
	EquipmentImportData      []EquipmentImportData
	EquipmentClassImportData EquipmentClassImportData
}

// Equipment
type EquipmentImportData struct {
	EquipmentName        string
	EquipmentTagBindings []EquipmentTagBindingData
}

type EquipmentTagBindingData struct {
	PropertyID string
	Tag        string
}

// Equipment Class
type EquipmentClassImportData struct {
	EquipmentClassName        string
	EquipmentClassDescription string
	EquipmentClassProperties  []EquipmentClassPropertyData
}

type EquipmentClassPropertyData struct {
	ID            string
	UnitOfMeasure struct {
		ID       string
		DataType string
	}
	Use bool
}

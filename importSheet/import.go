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
	var equipmentImportData *EquipmentImportData
	var err error

	switch filepath.Ext(*config.FilePath) {
	case ".csv":
		equipmentImportData, err = ReadCSV(*config.FilePath, *config.EquipmentClassDescription)
	case ".xlsx":
		equipmentImportData, err = ReadXLSX(*config.FilePath, *config.Sheet)
	default:
		log.Fatalf("Unsupported file extension \"%s\"\n", filepath.Ext(*config.FilePath))
	}
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// If UoM ID is not set, set it to DataType
	for i := range equipmentImportData.EquipmentClassProperties {
		if equipmentImportData.EquipmentClassProperties[i].UnitOfMeasure.ID == "" && equipmentImportData.EquipmentClassProperties[i].UnitOfMeasure.DataType != "" {
			equipmentImportData.EquipmentClassProperties[i].UnitOfMeasure.ID = equipmentImportData.EquipmentClassProperties[i].UnitOfMeasure.DataType
		}
	}

	log.Println("Adding Imported Unit of Measures")
	UnitOfMeasure(config.Context, client, *equipmentImportData)
	log.Println("Done Imported Unit of Measures")

	log.Println("Adding Imported Equipment model")
	EquipmentModel(config.Context, client, *equipmentImportData)
	log.Println("Done Imported Equipment model")

	log.Println("Done Imported model")
}

type EquipmentImportData struct {
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

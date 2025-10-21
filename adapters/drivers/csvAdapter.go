package drivers

import (
	"encoding/csv"
	"os"
	"rhize-data-collection-import/types"
)

type CSVAdapter struct {
	Description string
	Datasource  string
}

func (c CSVAdapter) Read(filePath string) (*types.ImportData, error) {
	// Read file into records
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Handle Equipment Class Data
	equipmentClassData := types.ImportEquipmentClass{
		Label:       records[4][0],
		Description: c.Description,
	}
	equipmentClassPropertyData := make([]types.ImportEquipmentClassProperty, 0)

	for i := 3; i < len(records); i++ {
		equipmentClassPropertyData = append(equipmentClassPropertyData, types.ImportEquipmentClassProperty{
			ID: records[i][1],
			UnitOfMeasure: struct {
				ID       string
				DataType string
			}{
				ID:       records[i][2],
				DataType: records[i][7],
			},
			Use: records[i][12] == "X",
		})
	}

	// Handle Import Data
	data := types.ImportData{
		Datasource:     "",
		EquipmentClass: equipmentClassData,
		Equipment:      nil,
	}

	return &data, nil
}

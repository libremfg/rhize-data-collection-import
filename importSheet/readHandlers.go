package importSheet

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func ReadCSV(filePath string, description string) (*EquipmentImportData, error) {
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

	// Handle data
	data := EquipmentImportData{
		EquipmentClassName:        records[4][0],
		EquipmentClassDescription: description,
	}
	equipmentClassPropertyData := make([]EquipmentClassPropertyData, 0)

	for i := 3; i < len(records); i++ {
		equipmentClassPropertyData = append(equipmentClassPropertyData, EquipmentClassPropertyData{
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

	return &data, nil
}

func ReadXLSX(filePath string, sheet string) (*EquipmentImportData, error) {
	// Read file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Handle data
	equipmentClassName, err := file.GetCellValue(sheet, "A3")
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	data := EquipmentImportData{
		EquipmentClassName:        equipmentClassName,
		EquipmentClassDescription: sheet,
	}
	equipmentClassPropertyData := make([]EquipmentClassPropertyData, 0)

	for i := 3; i < len(rows); i++ {
		row := rows[i]
		if row[1] == "" {
			continue
		}
		equipmentClassPropertyData = append(equipmentClassPropertyData, EquipmentClassPropertyData{
			ID: row[1],
			UnitOfMeasure: struct {
				ID       string
				DataType string
			}{
				ID:       row[2],
				DataType: row[7],
			},
			Use: row[12] == "X",
		})
	}
	data.EquipmentClassProperties = equipmentClassPropertyData

	return &data, nil
}

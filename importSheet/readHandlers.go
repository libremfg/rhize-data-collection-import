package importSheet

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func ReadCSV(filePath string, description string) (*ImportData, error) {
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
	equipmentClassData := EquipmentClassImportData{
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

	// Handle Import Data
	data := ImportData{
		Datasource:               "",
		EquipmentImportData:      nil,
		EquipmentClassImportData: equipmentClassData,
	}

	return &data, nil
}

func ReadXLSX(filePath string, sheet string, datasource string) (*ImportData, error) {
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

	// Handle Equipment Class Data
	equipmentClassName, err := file.GetCellValue(sheet, "A3")
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	equipmentClassData := EquipmentClassImportData{
		EquipmentClassName:        equipmentClassName,
		EquipmentClassDescription: sheet,
	}
	equipmentClassPropertyData := make([]EquipmentClassPropertyData, 0)

	for i := 3; i < len(rows); i++ {
		row := rows[i]
		if row[1] == "" {
			continue
		}
		if len(row) < 13 {
			log.Fatalf("\tERROR: The row %d has insufficient columns, requires 12 has %d. Unable to parse property from spreadsheet. Inspect the spreadsheet for errors and try again.", i, len(row))
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
	equipmentClassData.EquipmentClassProperties = equipmentClassPropertyData

	// Handle Equipment
	equipmentData := make([]EquipmentImportData, 0)

	if datasource != "" {
		for i := 13; i < len(rows[0]); i += 6 {
			if rows[0][i] == "" {
				break
			}
			tagBindings := make([]EquipmentTagBindingData, 0)
			for j := 3; j < len(rows); j++ {
				row := rows[j]
				if len(row) <= i+2 || row[i+2] == "" {
					continue
				}
				// Filter out Equipment Class name at start as it goes unused in binding
				propertyId := row[1][1+len(equipmentClassName):]

				tagBinding := EquipmentTagBindingData{
					PropertyID: propertyId,
					Tag:        row[i+2],
				}
				tagBindings = append(tagBindings, tagBinding)
			}

			data := EquipmentImportData{
				EquipmentName:        rows[0][i],
				EquipmentTagBindings: tagBindings,
			}
			equipmentData = append(equipmentData, data)
			// To-Do: Find a better solution to uneven column amounts than the below solutions
			if i == 13 {
				i++
			}
		}
	}

	// Handle Import Data
	data := ImportData{
		Datasource:               datasource,
		EquipmentImportData:      equipmentData,
		EquipmentClassImportData: equipmentClassData,
	}

	return &data, nil
}

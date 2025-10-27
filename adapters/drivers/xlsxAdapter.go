package drivers

import (
	"log"
	"rhize-data-collection-import/types"

	"github.com/xuri/excelize/v2"
)

type XLSXAdapter struct {
	Sheet      string
	Datasource string
}

func (x XLSXAdapter) Read(filePath string) (*types.ImportData, error) {
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
	equipmentClassName, err := file.GetCellValue(x.Sheet, "A3")
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows(x.Sheet)
	if err != nil {
		return nil, err
	}

	equipmentClassData := types.ImportEquipmentClass{
		Label:       equipmentClassName,
		Description: x.Sheet,
	}
	equipmentClassPropertyData := make([]types.ImportEquipmentClassProperty, 0)

	for i := 3; i < len(rows); i++ {
		row := rows[i]
		if row[1] == "" {
			continue
		}
		if len(row) < 13 {
			log.Fatalf("\tERROR: The row %d has insufficient columns, requires 12 has %d. Unable to parse property from spreadsheet. Inspect the spreadsheet for errors and try again.", i, len(row))
		}
		equipmentClassPropertyData = append(equipmentClassPropertyData, types.ImportEquipmentClassProperty{
			ID: row[1],
			UnitOfMeasure: types.ImportUnitOfMeasure{
				ID:       row[2],
				DataType: row[7],
			},
			Use: row[12] == "X",
		})
	}
	equipmentClassData.Properties = equipmentClassPropertyData

	// Handle Equipment
	equipmentData := make([]types.ImportEquipment, 0)

	if x.Datasource != "" {
		for i := 13; i < len(rows[0]); i += 6 {
			if rows[0][i] == "" {
				break
			}
			tagBindings := make([]types.ImportTagBinding, 0)
			for j := 3; j < len(rows); j++ {
				row := rows[j]
				if len(row) <= i+2 || row[i+2] == "" {
					continue
				}
				// Filter out Equipment Class name at start as it goes unused in binding
				propertyId := row[1][1+len(equipmentClassName):]

				tagBinding := types.ImportTagBinding{
					PropertyID: propertyId,
					Tag:        row[i+2],
				}
				// Optionally add in expressions for tag binding if they exist in comment column
				if len(row) > i+4 && row[i+4] != "" {
					tagBinding.Expression = row[i+4]
				}
				tagBindings = append(tagBindings, tagBinding)
			}

			data := types.ImportEquipment{
				ID:          rows[0][i],
				TagBindings: tagBindings,
			}
			equipmentData = append(equipmentData, data)
			// To-Do: Find a better solution to uneven column amounts than the below solutions
			if i == 13 {
				i++
			}
		}
	}

	// Handle Import Data
	data := types.ImportData{
		Datasource:     x.Datasource,
		EquipmentClass: equipmentClassData,
		Equipment:      equipmentData,
	}

	return &data, nil
}

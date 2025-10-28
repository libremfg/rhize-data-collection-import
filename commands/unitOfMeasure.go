package commands

import (
	"context"
	"errors"
	"log"
	"strings"

	"rhize-data-collection-import/domain"
	"rhize-data-collection-import/types"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
)

var (
	UnitOfMeasureCmd = &cobra.Command{
		Use:     "unitOfMeasure",
		Short:   "Import UoMs from file",
		Aliases: []string{"uom", "units"},
		Run:     importUnitOfMeasure,
	}
)

func importUnitOfMeasure(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log.Println("Starting import for Units of Measure")
	unitOfMeasure(ctx, Client, ImportData.EquipmentClass.Properties)
	log.Println("Finished import for Units of Measure")
}

func unitOfMeasure(ctx context.Context, client *graphql.Client, properties []types.ImportEquipmentClassProperty) {
	// Track added UoM ids
	units := make([]string, 0)

out:
	for _, property := range properties {
		uom := property.UnitOfMeasure.ID
		if uom == "" {
			continue out
		}

		// Check that UoM was not already added
		// Or had an attempt to be added
		for _, unit := range units {
			if unit == uom {
				continue out
			}
		}
		units = append(units, uom)

		dataType, err := convertDataType(property.UnitOfMeasure.DataType)
		if err != nil {
			log.Printf("\tUnknown data type \"%s\", skipping this Unit of Measure: %s\n", property.UnitOfMeasure.DataType, err.Error())
			continue out
		}

		unit := domain.AddUnitOfMeasureInput{
			ID:       uom,
			DataType: &dataType,
		}

		existingUoM, err := types.GetUnitOfMeasure(ctx, client, uom)
		if err != nil {
			log.Printf("\tcould not query unit of measure: %s", err.Error())
			continue out
		}
		operation := "Adding"
		if existingUoM != nil {
			operation = "Updating"
		}

		log.Printf("\t%s UoM for %s", operation, uom)
		err = types.CreateUnitOfMeasure(ctx, client, []domain.AddUnitOfMeasureInput{unit})
		if err != nil {
			log.Printf("\tcould not add unit of measure: %s", err.Error())
		}
	}

}

func convertDataType(inputDataType string) (domain.DataType, error) {
	// Check if DataType exists in Config
	if dataType, ok := DataTypesMap[strings.ToLower(inputDataType)]; ok {
		for i := range domain.AllDataType {
			if strings.ToUpper(dataType) == string(domain.AllDataType[i]) {
				return domain.AllDataType[i], nil
			}
		}
	}

	// Check if DataType already exists in Rhize
	for _, dataType := range domain.AllDataType {
		if strings.ToUpper(inputDataType) == string(dataType) {
			return dataType, nil
		}
	}

	// Try default data types
	var dataType domain.DataType
	switch strings.ToLower(inputDataType) {
	case "number":
		fallthrough
	case "double":
		dataType = domain.DataTypeFloat
	case "uint16":
		dataType = domain.DataTypeUINt16
	case "uint32":
		dataType = domain.DataTypeUINt32
	case "boolean":
		dataType = domain.DataTypeBool
	case "byte":
		dataType = domain.DataTypeByte
	case "datetime":
		dataType = domain.DataTypeDateTime
	case "localizedtext":
		fallthrough
	case "string":
		dataType = domain.DataTypeString
	default:
		return "", errors.New("data type does not exist in Rhize, config, or default types")
	}

	return dataType, nil
}

package types

import (
	"context"
	"net/http"
)

type Configuration struct {
	Context                   context.Context
	Client                    *http.Client
	URL                       *string
	FilePath                  *string
	Sheet                     *string
	EquipmentClassDescription *string
	Datasource                *string
	Target                    *string
}

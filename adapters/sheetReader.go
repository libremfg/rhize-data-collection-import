package adapters

import "rhize-data-collection-import/types"

type SheetReader interface {
	Read(filePath string) (*types.ImportData, error)
}

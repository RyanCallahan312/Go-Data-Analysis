package dto

import "fmt"

// UpdateSheetDTO for incoming sheet request
type UpdateSheetDTO struct {
	FileName  string `json:"fileName"`
	SheetName string `json:"sheetName"`
}

// TextOutput sheet fields to string
func (sheetDto UpdateSheetDTO) TextOutput() string {
	return fmt.Sprintf("File Name: %s,\n\t Sheet Name: %s",
		sheetDto.FileName, sheetDto.SheetName)
}

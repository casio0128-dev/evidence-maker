package conf

import "strings"

/*
	コンフィグの構造体
*/
type Config struct {
	TargetCol string       `json:"targetCol"`
	TargetRow int          `json:"targetStartRow"`
	Offset    int          `json:"offset"`
	Template  TemplateInfo `json:"template"`
}

type TemplateInfo struct {
	FilePath  string `json:"file"`
	SheetName string `json:"sheet"`
}

func (t TemplateInfo) IsFileSpecification() bool {
	return !strings.EqualFold(t.FilePath, "")
}

func (t TemplateInfo) IsSheetSpecification() bool {
	return !strings.EqualFold(t.SheetName, "")
}

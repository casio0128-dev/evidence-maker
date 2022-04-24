package conf

/*
	コンフィグの構造体
*/
type Config struct {
	TargetCol    string `json:"targetCol"`
	TargetRow    int    `json:"targetStartRow"`
	Offset       int    `json:"offset"`
	TemplatePath string `json:"template"`
}

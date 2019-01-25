package config

type LineChartConfig struct {
	Title string `yaml:"title"`
	Data []Data `yaml:"data"`
	Position Position `yaml:"position"`
	RefreshRateMs int `yaml:"refresh-rate-ms"`
	Scale string `yaml:"scale"`
}

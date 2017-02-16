package restuss

type PersistedScan struct {
	ID                   int64  `json:"id"`
	UUID                 string `json:"uuid"`
	Name                 string `json:"name"`
	Enabled              bool   `json:"enabled"`
	Status               string `json:"status"`
	CreationDate         int64  `json:"creation_date"`
	LastModificationDate int64  `json:"last_modification_date"`
	Owner                string `json:"owner"`
}

type Vulnerability struct {
	VulnerabilityIndex int    `json:"vuln_index"`
	Severity           int    `json:"severity"`
	PluginName         string `json:"plugin_name"`
	Count              int    `json:"count"`
	PluginId           int    `json:"plugin_id"`
	PluginFamily       string `json:"plugin_family"`
}

type ScanDetail struct {
	ID              int64
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type ScanSettings struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Targets string `json:"text_targets"`
}

type Scan struct {
	TemplateUUID string       `json:"uuid"`
	Settings     ScanSettings `json:"settings"`
}

type ScanTemplate struct {
	UUID             string `json:"uuid"`
	Name             string `json:"name"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	CloudOnly        bool   `json:"cloud_only"`
	SubscriptionOnly bool   `json:"subscription_only"`
	IsAgent          bool   `json:"is_agent"`
	Info             string `json:"more_info"`
}

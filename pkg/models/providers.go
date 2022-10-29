package models

type CloudProviders struct {
	GCP        GCPConfig        `yaml:"gcp"`
	Azure      AzureConfig      `yaml:"azure"`
	AWS        AWSConfig        `yaml:"aws"`
	CloudArmor CloudArmorConfig `yaml:"cloudarmor"`
}

type GCPConfig struct {
	Disabled  bool   `yaml:"disabled"`
	ProjectID string `yaml:"project_id"`
	Network   string `yaml:"network"`
	Priority  int64  `yaml:"priority"`
	MaxRules  int    `yaml:"max_rules"`
	// Endpoint is used for making calls to a mock server instead of the real Google services endpoints.
	Endpoint string `yaml:"endpoint"`
}

type AzureConfig struct {
	Disabled  bool   `yaml:"disabled"`
	SubscriptID string `yaml:"subscription_id"`
	ResourceGroup string `yaml:"resource_group"`
	Network   string `yaml:"network"`
	Priority  int64  `yaml:"priority"`
	MaxRules  int    `yaml:"max_rules"`
	UserAgent string `yaml:"user_agent"`
	// Endpoint is used for making calls to a mock server instead of the real Google services endpoints.
	Endpoint string `yaml:"endpoint"`
}

type CloudArmorConfig struct {
	Disabled  bool   `yaml:"disabled"`
	ProjectID string `yaml:"project_id"`
	Policy    string `yaml:"policy"`
	Priority  int64  `yaml:"priority"`
	MaxRules  int    `yaml:"max_rules"`
	// Endpoint is used for making calls to a mock server instead of the real Google services endpoints.
	Endpoint string `yaml:"endpoint"`
}

type AWSConfig struct {
	Disabled          bool   `yaml:"disabled"`
	Region            string `yaml:"region"`
	FirewallPolicy    string `yaml:"firewall_policy"`
	Capacity          int    `yaml:"capacity"`
	RuleGroupPriority int64  `yaml:"priority"`
	// Endpoint is used for making calls to a mock server instead of the real AWS services endpoints.
	Endpoint string `yaml:"endpoint"`
}

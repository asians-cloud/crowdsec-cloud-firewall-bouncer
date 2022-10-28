package gcp

import (
	"context"
	"fmt"

	"github.com/asians-cloud/crowdsec-cloud-firewall-bouncer/pkg/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

type Client struct {
	svc      GoogleComputeServiceIface
	project  string
	network  string
	maxRules int
	priority int64
}

const (
	providerName    = "gcp"
	defaultMaxRules = 10
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("provider", providerName)
}

func (c *Client) MaxSourcesPerRule() int {
	return 256
}
func (c *Client) MaxRules() int {
	return c.maxRules
}
func (c *Client) Priority() int64 {
	return c.priority
}

func getProjectIDFromCredentials(config *models.GCPConfig) (string, error) {
	ctx := context.Background()
	credentials, error := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if error != nil {
		return "", error
	}
	if credentials.ProjectID == "" {
		return "", fmt.Errorf("Default credentials does not have a project ID associated")
	}
	return credentials.ProjectID, nil
}

func checkGCPConfig(config *models.GCPConfig) error {
	if config == nil {
		return fmt.Errorf("gcp cloud provider must be specified")
	}
	if config.ProjectID == "" {
		var err error
		config.ProjectID, err = getProjectIDFromCredentials(config)
		if err != nil || config.ProjectID == "" {
			return fmt.Errorf("can't get project id from credentials: %s", err)
		}
	}
	if config.Network == "" {
		return fmt.Errorf("network must be specified in gcp config")
	}
	if config.MaxRules == 0 {
		config.MaxRules = defaultMaxRules
	}
	return nil
}

// NewClient creates a new GCP client
func NewClient(config *models.GCPConfig) (*Client, error) {
	log.Infof("creating client for %s", providerName)
	err := checkGCPConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error while checking GCP config: %s", err)
	}

	return &Client{
		svc:      NewGoogleComputeService(config.Endpoint),
		project:  config.ProjectID,
		network:  config.Network,
		priority: config.Priority,
		maxRules: config.MaxRules,
	}, nil
}

func (c *Client) GetProviderName() string {
	return providerName
}

func (c *Client) GetRules(ruleNamePrefix string) ([]*models.FirewallRule, error) {
	res, err := c.svc.ListFirewallRules(c.project, ruleNamePrefix)
	if err != nil {
		return nil, fmt.Errorf("unable to list firewall rules: %s", err)
	}
	var rules []*models.FirewallRule
	log.Infof("found %d rule(s)", len(res.Items))
	for _, gcpRule := range res.Items {
		log.Infof("%s: %#v", gcpRule.Name, gcpRule.SourceRanges)
		rule := models.FirewallRule{
			Name:         gcpRule.Name,
			SourceRanges: models.ConvertSourceRangesSliceToMap(gcpRule.SourceRanges),
			Priority:     gcpRule.Priority,
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}

func (c *Client) CreateRule(rule *models.FirewallRule) error {
	log.Infof("creating GCP firewall rule %s with %#v", rule.Name, rule.SourceRanges)

	denied := compute.FirewallDenied{
		IPProtocol: "all",
	}

	firewall := compute.Firewall{
		Direction:    "INGRESS",
		Denied:       []*compute.FirewallDenied{&denied},
		Network:      fmt.Sprintf("global/networks/%s", c.network),
		SourceRanges: models.ConvertSourceRangesMapToSlice(rule.SourceRanges),
		Name:         rule.Name,
		Description:  "Blocklist generated by CrowdSec Cloud Firewall Bouncer",
		Priority:     rule.Priority,
	}
	if err := c.svc.InsertFirewallRule(c.project, &firewall); err != nil {
		return fmt.Errorf("unable to create firewall rules %s: %s", rule.Name, err)
	}
	log.Infof("creation of rule %s successful", rule.Name)
	return nil
}

func (c *Client) DeleteRule(rule *models.FirewallRule) error {
	log.Infof("deleting GCP firewall rule %s", rule.Name)
	if err := c.svc.DeleteFirewallRule(c.project, rule.Name); err != nil {
		return fmt.Errorf("unable to delete firewall rule %s: %s", rule.Name, err)
	}
	log.Infof("deletion of rule %s successful", rule.Name)
	return nil
}

func (c *Client) PatchRule(rule *models.FirewallRule) error {
	log.Infof("patching GCP firewall rule %s with %#v", rule.Name, rule.SourceRanges)
	firewallPatchRequest := compute.Firewall{
		SourceRanges: models.ConvertSourceRangesMapToSlice(rule.SourceRanges),
	}
	if err := c.svc.PatchFirewallRule(c.project, rule.Name, &firewallPatchRequest); err != nil {
		return fmt.Errorf("unable to patch firewall rule %s: %s", rule.Name, err)
	}
	log.Infof("patching of rule %s successful", rule.Name)
	return nil
}

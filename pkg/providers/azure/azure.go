package azure

import (
	"context"
	"fmt"
	"strings"
	"github.com/asians-cloud/crowdsec-cloud-firewall-bouncer/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
)


type Client struct {
	svc               network.SecurityGroupsClient
	subscription      string
	GroupName         string
	network           string
	priority          int64
	maxRules          int
	capacity          int
	ruleGroupPriority int64

}

const (
	providerName          = "azure"
	defaultCapacity       = 4000
	defaultPriority int64 = 1
)

func (c *Client) MaxSourcesPerRule() int {
	return c.capacity
}
func (c *Client) MaxRules() int {
	return 1
}

func (c *Client) Priority() int64 {
	return c.ruleGroupPriority
}

func (c *Client) GetProviderName() string {
	return providerName
}

var log *logrus.Entry

func init() {
	log = logrus.WithField("provider", providerName)
}

func assignDefault(config *models.AzureConfig) {
	if config.Capacity == 0 {
		log.Debugf("Setting default rule group capacity (%d)", defaultCapacity)
		config.Capacity = defaultCapacity
	}
	if config.Priority == 0 {
		log.Debugf("Setting default lowest rule group priority (%d)", defaultPriority)
		config.Priority = defaultPriority
	}
}

// NewClient creates a new AWS client
func NewClient(config *models.AzureConfig) (*Client, error) {
	log.Infof("creating client for %s", providerName)

	
	svc := network.NewSecurityGroupsClient(config.SubscriptionID)

	assignDefault(config)

	_ = svc.AddToUserAgent(config.UserAgent)
	
	return &Client{
		svc:               svc,
		subscription:  config.SubscriptionID,
		GroupName: config.ResourceGroup,
		network:  config.Network,
		priority: config.Priority,
		maxRules: config.MaxRules,
	}, nil
}

func (c *Client) GetRules(ruleNamePrefix string) ([]*models.FirewallRule, error) {
	//azure
	//https://docs.microsoft.com/en-us/azure/virtual-network/security-overview
	//https://docs.microsoft.com/en-us/azure/virtual-network/security-overview#security-groups
	//https://docs.microsoft.com/en-us/azure/virtual-network/security-overview#network-security-groups
	var rules []*models.FirewallRule

	//Have trouble to list all rules because I can not get cluster name, which is a part of nsg name
	return rules, nil
}

func (c *Client) CreateRule(rule *models.FirewallRule) error {
	//Pre create rule
	return nil
}

func (c *Client) DeleteRule(rule *models.FirewallRule) error {
	return nil
}

func (c *Client) PatchRule(rule *models.FirewallRule) error {
	log.Infof("patching firewall rule %s with %#v", rule.Name, rule.SourceRanges)
	return nil
}

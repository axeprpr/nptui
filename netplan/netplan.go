package netplan

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	NetplanDir = "/etc/netplan"
	ConfigFile = "/etc/netplan/01-netcfg.yaml"
)

// NetworkConfig represents the netplan configuration structure
type NetworkConfig struct {
	Network Network `yaml:"network"`
}

type Network struct {
	Version   int                       `yaml:"version"`
	Renderer  string                    `yaml:"renderer,omitempty"`
	Ethernets map[string]EthernetConfig `yaml:"ethernets,omitempty"`
}

type EthernetConfig struct {
	DHCP4       bool     `yaml:"dhcp4,omitempty"`
	DHCP6       bool     `yaml:"dhcp6,omitempty"`
	Addresses   []string `yaml:"addresses,omitempty"`
	Gateway4    string   `yaml:"gateway4,omitempty"`
	Gateway6    string   `yaml:"gateway6,omitempty"`
	Routes      []Route  `yaml:"routes,omitempty"`
	Nameservers *DNS     `yaml:"nameservers,omitempty"`
	Optional    bool     `yaml:"optional,omitempty"`
}

type Route struct {
	To  string `yaml:"to"`
	Via string `yaml:"via"`
}

type DNS struct {
	Addresses []string `yaml:"addresses,omitempty"`
	Search    []string `yaml:"search,omitempty"`
}

// LoadConfig loads the netplan configuration
func LoadConfig() (*NetworkConfig, error) {
	// Try to read existing config
	files, err := filepath.Glob(filepath.Join(NetplanDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	// If no files exist, return default config
	if len(files) == 0 {
		return &NetworkConfig{
			Network: Network{
				Version:   2,
				Renderer:  "networkd",
				Ethernets: make(map[string]EthernetConfig),
			},
		}, nil
	}

	// Read the first yaml file
	data, err := ioutil.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	var config NetworkConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Initialize maps if nil
	if config.Network.Ethernets == nil {
		config.Network.Ethernets = make(map[string]EthernetConfig)
	}

	return &config, nil
}

// SaveConfig saves the netplan configuration
func SaveConfig(config *NetworkConfig) error {
	// Ensure netplan directory exists
	if err := os.MkdirAll(NetplanDir, 0755); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write to file
	if err := ioutil.WriteFile(ConfigFile, data, 0600); err != nil {
		return err
	}

	return nil
}

// GetInterfaces returns a list of network interfaces
func GetInterfaces() ([]string, error) {
	interfaces := []string{}
	
	// Read from /sys/class/net
	files, err := ioutil.ReadDir("/sys/class/net")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		name := f.Name()
		// Skip loopback
		if name != "lo" {
			interfaces = append(interfaces, name)
		}
	}

	return interfaces, nil
}

// ApplyConfig applies the netplan configuration
func ApplyConfig() error {
	// Run netplan apply
	cmd := exec.Command("netplan", "apply")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("netplan apply failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// GetInterfaceConfig gets config for a specific interface
func (c *NetworkConfig) GetInterfaceConfig(iface string) *EthernetConfig {
	if config, ok := c.Network.Ethernets[iface]; ok {
		return &config
	}
	return &EthernetConfig{}
}

// SetInterfaceConfig sets config for a specific interface
func (c *NetworkConfig) SetInterfaceConfig(iface string, config EthernetConfig) {
	if c.Network.Ethernets == nil {
		c.Network.Ethernets = make(map[string]EthernetConfig)
	}
	c.Network.Ethernets[iface] = config
}

// GetGateway returns the gateway from either gateway4 or routes
func (ec *EthernetConfig) GetGateway() string {
	// Try old format first
	if ec.Gateway4 != "" {
		return ec.Gateway4
	}
	
	// Try new routes format
	for _, route := range ec.Routes {
		if route.To == "default" {
			return route.Via
		}
	}
	
	return ""
}

// SetGateway sets the gateway using the new routes format
func (ec *EthernetConfig) SetGateway(gateway string) {
	if gateway == "" {
		ec.Routes = nil
		ec.Gateway4 = ""
		return
	}
	
	// Use new routes format
	ec.Routes = []Route{
		{
			To:  "default",
			Via: gateway,
		},
	}
	// Clear old format
	ec.Gateway4 = ""
}

// FormatConfig returns a human-readable string of the config
func (ec *EthernetConfig) FormatConfig() string {
	if ec.DHCP4 {
		return "DHCP (Automatic)"
	}
	
	if len(ec.Addresses) > 0 {
		config := fmt.Sprintf("Static: %s", ec.Addresses[0])
		
		gateway := ec.GetGateway()
		if gateway != "" {
			config += fmt.Sprintf("  Gateway: %s", gateway)
		}
		
		if ec.Nameservers != nil && len(ec.Nameservers.Addresses) > 0 {
			config += fmt.Sprintf("  DNS: %v", ec.Nameservers.Addresses)
		}
		return config
	}
	
	return "Not configured"
}


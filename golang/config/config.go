package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config is the top-level configuration for Faythe's config file.
type Config struct {
	ServerConfig     ServerConfig        `yaml:"server_config"`
	OpenStackConfigs []*OpenStackConfig  `yaml:"openstack_configs,omitempty"`
	StackConfig      []*StackStormConfig `yaml:"stackstorm_configs,omitempty"`
}

// ServerConfig configures values that are used to config Faythe HTTP server
type ServerConfig struct {
	// RemoteHostPattern can define an optional regexp pattern to be matched:
	//
	// - {name} matches anything until the next dot.
	//
	// - {name:pattern} matches the given regexp pattern.
	RemoteHostPattern   string              `yaml:"remote_host_pattern,omitempty"`
	BasicAuthentication BasicAuthentication `yaml:"basic_auth,omitempty"`
}

// BasicAuthentication - HTTP Basic authentication.
type BasicAuthentication struct {
	// Usenname, Password to implement HTTP basic authentication
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

/*
StackStormConfig stores information needed to forward
request to an StackStorm instance.
*/
type StackStormConfig struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	APIKey string `yaml:"api_key"`
}

// StackListOpts allows the filtering and sorting of paginated collections through
// the API.
type StackListOpts struct {
	// ProjectID is the UUID of the project.
	ProjectID string `yaml:"project_id,omitempty"`

	// ID filters the stack list by a stack ID.
	ID string `yaml:"id,omitempty"`

	// Status filters the stack list by a status.
	Status string `yaml:"status,omitempty"`

	// Name filters the stack list by a name.
	Name string `yaml:"name,omitempty"`

	// AllTenants is a bool to show all tenants.
	AllTenants bool `yaml:"all_tenants,omitempty"`

	// Tags lists stacks that contain one or more simple string tags.
	Tags string `yaml:"tags,omitempty"`

	// TagsAny lists stacks that contain one or more simple string tags.
	TagsAny string `yaml:"tags_any,omitempty"`

	// NotTags lists stacks that do not contain one or more simple string tags.
	NotTags string `yaml:"not_tags,omitempty"`

	// NotTagsAny lists stacks that do not contain one or more simple string tags.
	NotTagsAny string `yaml:"not_tags_any,omitempty"`
}

// StackQuery stores information needed to query Heat stacks.
type StackQuery struct {
	// UpdateInterval field is the number of seconds that queries the outputs of stacks
	// that was filters with a given listOpts periodically.
	UpdateInterval time.Duration `yaml:"update_interval"`

	// ListOpts field is the list of Stack list options.
	ListOpts StackListOpts `yaml:"list_opts,omitempty"`
}

// OpenStackConfig stores information needed to authenticate to an OpenStack Cloud.
type OpenStackConfig struct {
	// The name to define which OpenStack
	Name string `yaml:"name"`
	// AuthURL specifies the HTTP endpoint that is required to work with
	// the Identity API of the appropriate version. While it's ultimately needed by
	// all of the identity services, it will often be populated by a provider-level
	// function.
	AuthURL    string `yaml:"auth_url"`
	RegionName string `yaml:"region_name"`

	// Username is required if using Identity V2 API. Consult with your provider's
	// control panel to discover your account's username. In Identity V3, either
	// UserID or a combination of Username and DomainID or DomainName are needed.
	Username string `yaml:"username"`
	UserID   string `yaml:"userid"`

	Password string `yaml:"password"`

	// At most one of DomainID and DomainName must be provided if using Username
	// with Identity V3. Otherwise, either are optional.
	DomainName string `yaml:"domain_name"`
	DomainID   string `yaml:"domain_id"`

	// The ProjectID and ProjectName fields are optional for the Identity V2 API.
	// The same fields are known as project_id and project_name in the Identity
	// V3 API, but are collected as ProjectID and ProjectName here in both cases.
	// Some providers allow you to specify a ProjectName instead of the ProjectId.
	// Some require both. Your provider's authentication policies will determine
	// how these fields influence authentication.
	// If DomainID or DomainName are provided, they will also apply to ProjectName.
	// It is not currently possible to authenticate with Username and a Domain
	// and scope to a Project in a different Domain by using ProjectName. To
	// accomplish that, the ProjectID will need to be provided as the ProjectID
	// option.
	ProjectName string `yaml:"project_id"`
	ProjectID   string `yaml:"project_name"`

	StackQuery StackQuery `yaml:"stack_query"`
}

var (
	// DefaultConfig is the default top-level configuration.
	DefaultConfig = Config{
		ServerConfig: DefaultServerConfig,
	}

	// DefaultServerConfig is the default server configuration.
	DefaultServerConfig = ServerConfig{
		RemoteHostPattern:   ".*",
		BasicAuthentication: BasicAuthentication{},
	}

	// DefaultStackQuery is the default stack query with update
	// interval 30s.
	DefaultStackQuery = StackQuery{
		UpdateInterval: 30 * time.Second,
	}

	// DefaultOpenStackConfig is the default OpenStack configuration.
	DefaultOpenStackConfig = OpenStackConfig{
		StackQuery: DefaultStackQuery,
	}
)

// Load parses the YAML input s into a Config
func Load(s string) (*Config, error) {
	cfg := &Config{}
	// If the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have  to set the DefaultConfig at the entry
	// point as well.
	*cfg = DefaultConfig

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	return cfg, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *ServerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultServerConfig
	type plain ServerConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}

// UnmarshalYAML implemnets the yaml.Unmarshaler interface
func (c *OpenStackConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultOpenStackConfig
	type plain OpenStackConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.Name == "" {
		return errors.New("openstack configuration requires a name")
	}
	if c.RegionName == "" {
		return errors.New("openstack configration requires a region")
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *StackQuery) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultStackQuery
	type plain StackQuery
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.UpdateInterval == 0 {
		c.UpdateInterval = DefaultStackQuery.UpdateInterval
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *StackListOpts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = StackListOpts{}
	type plain StackListOpts
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *StackStormConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = StackStormConfig{}
	type plain StackStormConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.Name == "" {
		return errors.New("stackstorm configuration requires a name")
	}
	if c.Host == "" {
		return errors.New("stackstorm configuration requires host address/host name")
	}
	if c.APIKey == "" {
		return errors.New("stackstorm configuration requires api key")
	}
	return nil
}

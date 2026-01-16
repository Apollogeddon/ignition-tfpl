package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// Client holds the configuration for the Ignition API client
type Client struct {
	HostURL    string
	HTTPClient *retryablehttp.Client
	Token      string
}

// APIErrorResponse represents the structure of error responses from the Ignition API
type APIErrorResponse struct {
	Success       bool     `json:"success"`
	Messages      []string `json:"messages,omitempty"`
	FieldMessages []struct {
		FieldName string   `json:"fieldName"`
		Messages  []string `json:"messages"`
	} `json:"fieldMessages,omitempty"`
	Problem *struct {
		Message    string   `json:"message"`
		StackTrace []string `json:"stacktrace,omitempty"`
	} `json:"problem,omitempty"`
}

func (e *APIErrorResponse) Error() string {
	var msg string
	if e.Problem != nil {
		msg = fmt.Sprintf("API error: %s", e.Problem.Message)
	} else if len(e.Messages) > 0 {
		msg = fmt.Sprintf("API error: %v", e.Messages)
	} else {
		msg = "unknown API error"
	}

	if len(e.FieldMessages) > 0 {
		msg += fmt.Sprintf(" (Field Errors: %v)", e.FieldMessages)
	}
	return msg
}

// ResourceResponse represents a generic resource wrapper in Ignition
type ResourceResponse[T any] struct {
	Module      string `json:"module,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name"`
	Enabled     *bool  `json:"enabled,omitempty"`
	Description string `json:"description,omitempty"`
	Signature   string `json:"signature,omitempty"`
	Config      T      `json:"config"`
}

// ResourceChangesResponse represents the response for create/update/delete operations
type ResourceChangesResponse struct {
	Success bool `json:"success"`
	Changes []struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Collection   string `json:"collection"`
		NewSignature string `json:"newSignature"`
	} `json:"changes"`
}

// DatabaseConfig represents the specific configuration for a database connection
type DatabaseConfig struct {
	Driver     string `json:"driver"`
	Translator string `json:"translator,omitempty"`
	ConnectURL string `json:"connectURL"`
	Username   string `json:"username,omitempty"`
	Password   any    `json:"password,omitempty"` // Can be complex object or null
}

// TagProviderConfig represents the specific configuration for a tag provider
type TagProviderConfig struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// UserSourceProfile represents the profile configuration within a user source
type UserSourceProfile struct {
	Type               string `json:"type"`
	FailoverProfile    string `json:"failoverProfile,omitempty"`
	FailoverMode       string `json:"failoverMode,omitempty"`
	ScheduleRestricted bool   `json:"scheduleRestricted,omitempty"`
}

// UserSourceConfig represents the specific configuration for a user source
type UserSourceConfig struct {
	Profile UserSourceProfile `json:"profile"`
}

// Project represents an Ignition Project
type Project struct {
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Title            string `json:"title,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	Parent           string `json:"parent,omitempty"`
	Inheritable      bool   `json:"inheritable,omitempty"`
	DefaultDB        string `json:"defaultDb,omitempty"`
	TagProvider      string `json:"tagProvider,omitempty"`
	UserSource       string `json:"userSource,omitempty"`
	IdentityProvider string `json:"identityProvider,omitempty"`
}

// AuditProfileSettings represents the type-specific settings for an audit profile
type AuditProfileSettings struct {
	DatabaseName          string `json:"databaseName,omitempty"`
	PruneEnabled          bool   `json:"pruneEnabled,omitempty"`
	AutoCreate            bool   `json:"autoCreate,omitempty"`
	TableName             string `json:"tableName,omitempty"`
	RemoteServer          string `json:"remoteServer,omitempty"`
	RemoteProfile         string `json:"remoteProfile,omitempty"`
	EnableStoreAndForward bool   `json:"enableStoreAndForward,omitempty"`
}

// AuditProfileProfile represents the core profile configuration
type AuditProfileProfile struct {
	Type          string `json:"type"`
	RetentionDays int    `json:"retentionDays,omitempty"`
}

// AuditProfileConfig represents the configuration for an audit profile
type AuditProfileConfig struct {
	Profile  AuditProfileProfile  `json:"profile"`
	Settings AuditProfileSettings `json:"settings"`
}

// AlarmNotificationProfileProfile represents the profile configuration
type AlarmNotificationProfileProfile struct {
	Type string `json:"type"`
}

// AlarmNotificationProfileConfig represents the main config object
type AlarmNotificationProfileConfig struct {
	Profile  AlarmNotificationProfileProfile `json:"profile"`
	Settings map[string]any                  `json:"settings"`
}

// OpcUaConnectionEndpoint represents the endpoint configuration
type OpcUaConnectionEndpoint struct {
	DiscoveryURL   string `json:"discoveryUrl"`
	EndpointURL    string `json:"endpointUrl"`
	SecurityPolicy string `json:"securityPolicy"`
	SecurityMode   string `json:"securityMode"`
}

// OpcUaConnectionSettings represents the settings for an OPC UA connection
type OpcUaConnectionSettings struct {
	Endpoint OpcUaConnectionEndpoint `json:"endpoint"`
}

// OpcUaConnectionProfile represents the profile configuration
type OpcUaConnectionProfile struct {
	Type string `json:"type"` // e.g., "com.inductiveautomation.OpcUaServerType"
}

// OpcUaConnectionConfig represents the main config object
type OpcUaConnectionConfig struct {
	Profile  OpcUaConnectionProfile  `json:"profile"`
	Settings OpcUaConnectionSettings `json:"settings"`
}

// AlarmJournalProfile represents the profile configuration
type AlarmJournalProfile struct {
	Type      string `json:"type"` // e.g., "DATASOURCE"
	QueryOnly bool   `json:"queryOnly,omitempty"`
}

// AlarmJournalSettings represents the union of settings types
type AlarmJournalSettings struct {
	// DATASOURCE
	Datasource string `json:"datasource,omitempty"`
	
	// REMOTE
	RemoteGateway *struct {
		TargetServer  string `json:"targetServer,omitempty"`
		TargetJournal string `json:"targetJournal,omitempty"`
	} `json:"remoteGateway,omitempty"`
	
	// Common / Shared Objects
	Advanced *struct {
		TableName          string `json:"tableName,omitempty"` // DATASOURCE only
		DataTableName      string `json:"dataTableName,omitempty"` // DATASOURCE only
		UseStoreAndForward bool   `json:"useStoreAndForward,omitempty"` // Both
	} `json:"advanced,omitempty"`

	Events *struct {
		MinPriority            string `json:"minPriority,omitempty"`
		StoreShelvedEvents     bool   `json:"storeShelvedEvents,omitempty"`
		StoreFromEnabledChange bool   `json:"storeFromEnabledChange,omitempty"`
	} `json:"events,omitempty"`
}

// AlarmJournalConfig represents the main config object
type AlarmJournalConfig struct {
	Profile  AlarmJournalProfile  `json:"profile"`
	Settings AlarmJournalSettings `json:"settings"`
}

// SMTPProfileSettingsClassic represents the settings for a classic SMTP profile
type SMTPProfileSettingsClassic struct {
	Hostname        string `json:"hostname"`
	Port            int    `json:"port"`
	UseSslPort      bool   `json:"useSslPort"`
	StartTlsEnabled bool   `json:"startTlsEnabled"`
	Username        string `json:"username,omitempty"`
	Password        any    `json:"password,omitempty"`
}

// SMTPProfileSettings represents the union of settings for Email profiles
type SMTPProfileSettings struct {
	Settings *SMTPProfileSettingsClassic `json:"settings,omitempty"`
}

// SMTPProfileProfile represents the profile configuration
type SMTPProfileProfile struct {
	Type string `json:"type"` // "smtp.classic" or "smtp.oauth2"
}

// SMTPProfileConfig represents the main config object for an Email profile
type SMTPProfileConfig struct {
	Profile  SMTPProfileProfile  `json:"profile"`
	Settings SMTPProfileSettings `json:"settings"`
}

// StoreAndForwardMaintenancePolicy represents the maintenance policy for a datastore
type StoreAndForwardMaintenancePolicy struct {
	Action string `json:"action"`
	Limit  struct {
		LimitType string `json:"limitType"`
		Value     int    `json:"value"`
	} `json:"limit"`
}

// StoreAndForwardConfig represents the configuration for a Store and Forward Engine
type StoreAndForwardConfig struct {
	TimeThresholdMs            int                               `json:"timeThresholdMs"`
	ForwardRateMs              int                               `json:"forwardRateMs"`
	ForwardingPolicy           string                            `json:"forwardingPolicy"`
	ForwardingSchedule         string                            `json:"forwardingSchedule,omitempty"`
	IsThirdParty               bool                              `json:"isThirdParty"`
	DataThreshold              int                               `json:"dataThreshold"`
	BatchSize                  int                               `json:"batchSize"`
	ScanRateMs                 int                               `json:"scanRateMs"`
	PrimaryMaintenancePolicy   *StoreAndForwardMaintenancePolicy `json:"primaryMaintenancePolicy,omitempty"`
	SecondaryMaintenancePolicy *StoreAndForwardMaintenancePolicy `json:"secondaryMaintenancePolicy,omitempty"`
}

// IdentityProviderAuthMethod represents an authentication method for an IdP
type IdentityProviderAuthMethod struct {
	Type   string `json:"type"` // "basic" or "badge"
	Config any    `json:"config"`
}

// IdentityProviderInternalConfig represents the "internal" IdP type
type IdentityProviderInternalConfig struct {
	UserSource               string                       `json:"userSource"`
	AuthMethods              []IdentityProviderAuthMethod `json:"authMethods"`
	SessionInactivityTimeout float64                      `json:"sessionInactivityTimeout"`
	SessionExp               float64                      `json:"sessionExp"`
	RememberMeExp            float64                      `json:"rememberMeExp"`
}

// IgnitionSecret represents a secret in Ignition (embedded or referenced)
type IgnitionSecret struct {
	Type string `json:"type"` // "Embedded" or "Referenced"
	Data any    `json:"data"`
}

// IdentityProviderOidcConfig represents the "oidc" IdP type
type IdentityProviderOidcConfig struct {
	ClientId                   string          `json:"clientId"`
	ClientSecret               *IgnitionSecret `json:"clientSecret,omitempty"`
	ProviderId                 string          `json:"providerId"`
	AuthorizationEndpoint      string          `json:"authorizationEndpoint"`
	TokenEndpoint              string          `json:"tokenEndpoint"`
	JsonWebKeysEndpoint        string          `json:"jsonWebKeysEndpoint"`
	JsonWebKeysEndpointEnabled bool            `json:"jsonWebKeysEndpointEnabled"`
	UserInfoEndpoint           string          `json:"userInfoEndpoint,omitempty"`
	EndSessionEndpoint         string          `json:"endSessionEndpoint,omitempty"`
}

// IdentityProviderSamlConfig represents the "saml" IdP type
type IdentityProviderSamlConfig struct {
	IdpEntityId                    string `json:"idpEntityId"`
	SpEntityId                     string `json:"spEntityId,omitempty"`
	SpEntityIdEnabled              bool   `json:"spEntityIdEnabled"`
	AcsBinding                     string `json:"acsBinding"`
	NameIdFormat                   string `json:"nameIdFormat"`
	SsoServiceConfig               struct {
		Uri     string `json:"uri"`
		Binding string `json:"binding"`
	} `json:"ssoServiceConfig"`
	ForceAuthnEnabled              bool     `json:"forceAuthnEnabled"`
	ResponseSignaturesRequired     bool     `json:"responseSignaturesRequired"`
	AssertionSignaturesRequired    bool     `json:"assertionSignaturesRequired"`
	IdpMetadataUrl                 string   `json:"idpMetadataUrl,omitempty"`
	IdpMetadataUrlEnabled          bool     `json:"idpMetadataUrlEnabled"`
	SignatureVerifyingCertificates []string `json:"signatureVerifyingCertificates"`
	SignatureVerifyingKeys         []any    `json:"signatureVerifyingKeys"`
}

// IdentityProviderConfig represents the main config object for an IdP
type IdentityProviderConfig struct {
	Type   string `json:"type"` // "internal", "oidc", "saml"
	Config any    `json:"config"`
}

// GanOutgoingConfig represents the configuration for an outgoing GAN connection
type GanOutgoingConfig struct {
	Host                     string  `json:"host"`
	Port                     int     `json:"port"`
	UseSSL                   bool    `json:"useSSL"`
	PingRateMillis           float64 `json:"pingRateMillis,omitempty"`
	PingTimeoutMillis        float64 `json:"pingTimeoutMillis,omitempty"`
	PingMaxMissed            float64 `json:"pingMaxMissed,omitempty"`
	WsTimeoutMillis          float64 `json:"wsTimeoutMillis,omitempty"`
	HttpConnectTimeoutMillis float64 `json:"httpConnectTimeoutMillis,omitempty"`
	HttpReadTimeoutMillis    float64 `json:"httpReadTimeoutMillis,omitempty"`
	SendThreads              float64 `json:"sendThreads,omitempty"`
	ReceiveThreads           float64 `json:"receiveThreads,omitempty"`
}

// RedundancyConfig represents the configuration for gateway redundancy
type RedundancyConfig struct {
	Role                string `json:"role"` // "Independent", "Backup", "Master"
	ActiveHistoryLevel  string `json:"activeHistoryLevel"`
	JoinWaitTime        int    `json:"joinWaitTime"`
	RecoveryMode        string `json:"recoveryMode"`
	AllowHistoryCleanup bool   `json:"allowHistoryCleanup"`
	GatewayNetworkSetup *struct {
		Host               string  `json:"host,omitempty"`
		Port               int     `json:"port,omitempty"`
		EnableSsl          bool    `json:"enableSsl,omitempty"`
		PingRate           float64 `json:"pingRate,omitempty"`
		PingTimeout        float64 `json:"pingTimeout,omitempty"`
		PingMaxMissed      float64 `json:"pingMaxMissed,omitempty"`
		WebsocketTimeout   float64 `json:"websocketTimeout,omitempty"`
		HttpConnectTimeout float64 `json:"httpConnectTimeout,omitempty"`
		HttpReadTimeout    float64 `json:"httpReadTimeout,omitempty"`
		SendThreads        float64 `json:"sendThreads,omitempty"`
		ReceiveThreads     float64 `json:"receiveThreads,omitempty"`
	} `json:"gatewayNetworkSetup,omitempty"`
}

// GanGeneralSettingsConfig represents general Gateway Network settings
type GanGeneralSettingsConfig struct {
	RequireSSL                  bool    `json:"requireSSL"`
	RequireTwoWayAuth           bool    `json:"requireTwoWayAuth"`
	AllowIncoming               bool    `json:"allowIncoming"`
	SecurityPolicy              string  `json:"securityPolicy"`
	Whitelist                   string  `json:"whitelist,omitempty"`
	AllowedProxyHops            float64 `json:"allowedProxyHops"`
	WebsocketSessionIdleTimeout float64 `json:"websocketSessionIdleTimeout"`
	TempFilesMaxAgeHours        float64 `json:"tempFilesMaxAgeHours"`
}

// NewClient creates a new Ignition API client
func NewClient(host, token string, allowInsecureTLS bool) (*Client, error) {
	rc := retryablehttp.NewClient()
	rc.RetryMax = 10
	rc.Logger = nil // Disable default logging to avoid noise
	rc.HTTPClient.Timeout = 10 * time.Second
	
	if allowInsecureTLS {
		rc.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	c := &Client{
		HTTPClient: rc,
		HostURL:    host,
		Token:      token,
	}

	return c, nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	req, err := retryablehttp.NewRequestWithContext(ctx, method, c.HostURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Ignition-API-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var apiErr APIErrorResponse
		if err := json.Unmarshal(bodyBytes, &apiErr); err == nil && apiErr.Problem != nil {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, bodyBytes)
	}

	return bodyBytes, nil
}

// Generic CRUD methods

func (c *Client) GetResource(ctx context.Context, resourceType, name string, dest any) error {
	path := fmt.Sprintf("/data/api/v1/resources/find/ignition/%s/%s", resourceType, name)
	body, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

func (c *Client) CreateResource(ctx context.Context, resourceType string, item any, dest any) error {
	// API expects an array for creation
	rb, err := json.Marshal([]any{item})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/data/api/v1/resources/ignition/%s", resourceType)
	body, err := c.doRequest(ctx, http.MethodPost, path, rb)
	if err != nil {
		return err
	}

	return c.unmarshalResourceResponse(ctx, "ignition", resourceType, body, dest)
}

func (c *Client) UpdateResource(ctx context.Context, resourceType string, item any, dest any) error {
	// API expects an array for update
	rb, err := json.Marshal([]any{item})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/data/api/v1/resources/ignition/%s", resourceType)
	body, err := c.doRequest(ctx, http.MethodPut, path, rb)
	if err != nil {
		return err
	}

	return c.unmarshalResourceResponse(ctx, "ignition", resourceType, body, dest)
}

func (c *Client) unmarshalResourceResponse(ctx context.Context, module, resourceType string, body []byte, dest any) error {
	if len(body) == 0 {
		return fmt.Errorf("empty response body")
	}

	// Try unmarshaling as ResourceChangesResponse first
	var changes ResourceChangesResponse
	if err := json.Unmarshal(body, &changes); err == nil && len(changes.Changes) > 0 {
		name := changes.Changes[0].Name
		if name == "" {
			return fmt.Errorf("API returned success but no resource name in changes")
		}
		// Since the create/update response doesn't contain the full object, fetch it
		return c.GetResourceWithModule(ctx, module, resourceType, name, dest)
	}

	// Fallback to old behavior: API returns an array or a single object of resources
	if body[0] == '{' {
		return json.Unmarshal(body, dest)
	}

	var rawItems []json.RawMessage
	if err := json.Unmarshal(body, &rawItems); err != nil {
		return err
	}

	if len(rawItems) == 0 {
		return fmt.Errorf("no resources returned from API")
	}

	return json.Unmarshal(rawItems[0], dest)
}

func (c *Client) GetResourceWithModule(ctx context.Context, module, resourceType, name string, dest any) error {
	path := fmt.Sprintf("/data/api/v1/resources/find/%s/%s/%s", module, resourceType, name)
	body, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

func (c *Client) CreateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error {
	// API expects an array for creation
	rb, err := json.Marshal([]any{item})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/data/api/v1/resources/%s/%s", module, resourceType)
	body, err := c.doRequest(ctx, http.MethodPost, path, rb)
	if err != nil {
		return err
	}

	return c.unmarshalResourceResponse(ctx, module, resourceType, body, dest)
}

func (c *Client) UpdateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error {
	// API expects an array for update
	rb, err := json.Marshal([]any{item})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/data/api/v1/resources/%s/%s", module, resourceType)
	body, err := c.doRequest(ctx, http.MethodPut, path, rb)
	if err != nil {
		return err
	}

	return c.unmarshalResourceResponse(ctx, module, resourceType, body, dest)
}

func (c *Client) DeleteResourceWithModule(ctx context.Context, module, resourceType, name, signature string) error {
	path := fmt.Sprintf("/data/api/v1/resources/%s/%s/%s/%s", module, resourceType, name, signature)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

func (c *Client) DeleteResource(ctx context.Context, resourceType, name, signature string) error {
	path := fmt.Sprintf("/data/api/v1/resources/ignition/%s/%s/%s", resourceType, name, signature)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

// Projects

func (c *Client) GetProject(ctx context.Context, name string) (*Project, error) {
	path := fmt.Sprintf("/data/api/v1/projects/find/%s", name)
	body, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var p Project
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) CreateProject(ctx context.Context, p Project) (*Project, error) {
	rb, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(ctx, http.MethodPost, "/data/api/v1/projects", rb)
	if err != nil {
		return nil, err
	}

	return c.waitForProject(ctx, p.Name)
}

func (c *Client) UpdateProject(ctx context.Context, p Project) (*Project, error) {
	rb, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/data/api/v1/projects/%s", p.Name)
	_, err = c.doRequest(ctx, http.MethodPut, path, rb)
	if err != nil {
		return nil, err
	}

	return c.waitForProject(ctx, p.Name)
}

func (c *Client) DeleteProject(ctx context.Context, name string) error {
	path := fmt.Sprintf("/data/api/v1/projects/%s", name)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

func (c *Client) waitForProject(ctx context.Context, name string) (*Project, error) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return c.GetProject(ctx, name)
		case <-ticker.C:
			project, err := c.GetProject(ctx, name)
			if err == nil {
				return project, nil
			}
		}
	}
}

// Database Connections

func (c *Client) GetDatabaseConnection(ctx context.Context, name string) (*ResourceResponse[DatabaseConfig], error) {
	var resp ResourceResponse[DatabaseConfig]
	err := c.GetResource(ctx, "database-connection", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error) {
	var resp ResourceResponse[DatabaseConfig]
	err := c.CreateResource(ctx, "database-connection", db, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error) {
	var resp ResourceResponse[DatabaseConfig]
	err := c.UpdateResource(ctx, "database-connection", db, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteDatabaseConnection(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "database-connection", name, signature)
}

// User Sources

func (c *Client) GetUserSource(ctx context.Context, name string) (*ResourceResponse[UserSourceConfig], error) {
	var resp ResourceResponse[UserSourceConfig]
	err := c.GetResource(ctx, "user-source", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error) {
	var resp ResourceResponse[UserSourceConfig]
	err := c.CreateResource(ctx, "user-source", us, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error) {
	var resp ResourceResponse[UserSourceConfig]
	err := c.UpdateResource(ctx, "user-source", us, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteUserSource(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "user-source", name, signature)
}

// Tag Providers

func (c *Client) GetTagProvider(ctx context.Context, name string) (*ResourceResponse[TagProviderConfig], error) {
	var resp ResourceResponse[TagProviderConfig]
	err := c.GetResource(ctx, "tag-provider", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error) {
	var resp ResourceResponse[TagProviderConfig]
	err := c.CreateResource(ctx, "tag-provider", tp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error) {
	var resp ResourceResponse[TagProviderConfig]
	err := c.UpdateResource(ctx, "tag-provider", tp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteTagProvider(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "tag-provider", name, signature)
}

// Audit Profiles

func (c *Client) GetAuditProfile(ctx context.Context, name string) (*ResourceResponse[AuditProfileConfig], error) {
	var resp ResourceResponse[AuditProfileConfig]
	err := c.GetResource(ctx, "audit-profile", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error) {
	var resp ResourceResponse[AuditProfileConfig]
	err := c.CreateResource(ctx, "audit-profile", ap, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error) {
	var resp ResourceResponse[AuditProfileConfig]
	err := c.UpdateResource(ctx, "audit-profile", ap, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteAuditProfile(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "audit-profile", name, signature)
}

// Alarm Notification Profiles

func (c *Client) GetAlarmNotificationProfile(ctx context.Context, name string) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	var resp ResourceResponse[AlarmNotificationProfileConfig]
	err := c.GetResourceWithModule(ctx, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	var resp ResourceResponse[AlarmNotificationProfileConfig]
	err := c.CreateResourceWithModule(ctx, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", anp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	var resp ResourceResponse[AlarmNotificationProfileConfig]
	err := c.UpdateResourceWithModule(ctx, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", anp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteAlarmNotificationProfile(ctx context.Context, name, signature string) error {
	return c.DeleteResourceWithModule(ctx, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", name, signature)
}

// OPC UA Connections

func (c *Client) GetOpcUaConnection(ctx context.Context, name string) (*ResourceResponse[OpcUaConnectionConfig], error) {
	var resp ResourceResponse[OpcUaConnectionConfig]
	err := c.GetResourceWithModule(ctx, "ignition", "opc-connection", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error) {
	var resp ResourceResponse[OpcUaConnectionConfig]
	err := c.CreateResourceWithModule(ctx, "ignition", "opc-connection", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error) {
	var resp ResourceResponse[OpcUaConnectionConfig]
	err := c.UpdateResourceWithModule(ctx, "ignition", "opc-connection", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteOpcUaConnection(ctx context.Context, name, signature string) error {
	return c.DeleteResourceWithModule(ctx, "ignition", "opc-connection", name, signature)
}

// Alarm Journals

func (c *Client) GetAlarmJournal(ctx context.Context, name string) (*ResourceResponse[AlarmJournalConfig], error) {
	var resp ResourceResponse[AlarmJournalConfig]
	err := c.GetResourceWithModule(ctx, "ignition", "alarm-journal", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error) {
	var resp ResourceResponse[AlarmJournalConfig]
	err := c.CreateResourceWithModule(ctx, "ignition", "alarm-journal", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error) {
	var resp ResourceResponse[AlarmJournalConfig]
	err := c.UpdateResourceWithModule(ctx, "ignition", "alarm-journal", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteAlarmJournal(ctx context.Context, name, signature string) error {
	return c.DeleteResourceWithModule(ctx, "ignition", "alarm-journal", name, signature)
}

// SMTP Profiles

func (c *Client) GetSMTPProfile(ctx context.Context, name string) (*ResourceResponse[SMTPProfileConfig], error) {
	var resp ResourceResponse[SMTPProfileConfig]
	err := c.GetResource(ctx, "email-profile", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error) {
	var resp ResourceResponse[SMTPProfileConfig]
	err := c.CreateResource(ctx, "email-profile", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error) {
	var resp ResourceResponse[SMTPProfileConfig]
	err := c.UpdateResource(ctx, "email-profile", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteSMTPProfile(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "email-profile", name, signature)
}

// Store and Forward

func (c *Client) GetStoreAndForward(ctx context.Context, name string) (*ResourceResponse[StoreAndForwardConfig], error) {
	var resp ResourceResponse[StoreAndForwardConfig]
	err := c.GetResource(ctx, "store-and-forward-engine", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error) {
	var resp ResourceResponse[StoreAndForwardConfig]
	err := c.CreateResource(ctx, "store-and-forward-engine", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error) {
	var resp ResourceResponse[StoreAndForwardConfig]
	err := c.UpdateResource(ctx, "store-and-forward-engine", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteStoreAndForward(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "store-and-forward-engine", name, signature)
}

// Identity Providers

func (c *Client) GetIdentityProvider(ctx context.Context, name string) (*ResourceResponse[IdentityProviderConfig], error) {
	var resp ResourceResponse[IdentityProviderConfig]
	err := c.GetResource(ctx, "identity-provider", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error) {
	var resp ResourceResponse[IdentityProviderConfig]
	err := c.CreateResource(ctx, "identity-provider", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error) {
	var resp ResourceResponse[IdentityProviderConfig]
	err := c.UpdateResource(ctx, "identity-provider", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteIdentityProvider(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "identity-provider", name, signature)
}

// Gateway Network Outgoing

func (c *Client) GetGanOutgoing(ctx context.Context, name string) (*ResourceResponse[GanOutgoingConfig], error) {
	var resp ResourceResponse[GanOutgoingConfig]
	err := c.GetResource(ctx, "gateway-network-outgoing", name, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error) {
	var resp ResourceResponse[GanOutgoingConfig]
	err := c.CreateResource(ctx, "gateway-network-outgoing", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error) {
	var resp ResourceResponse[GanOutgoingConfig]
	err := c.UpdateResource(ctx, "gateway-network-outgoing", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteGanOutgoing(ctx context.Context, name, signature string) error {
	return c.DeleteResource(ctx, "gateway-network-outgoing", name, signature)
}

// Redundancy

func (c *Client) GetRedundancyConfig(ctx context.Context) (*RedundancyConfig, error) {
	path := "/data/api/v1/redundancy/config"
	body, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var config RedundancyConfig
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Client) UpdateRedundancyConfig(ctx context.Context, config RedundancyConfig) error {
	path := "/data/api/v1/redundancy/config"
	rb, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, http.MethodPost, path, rb)
	return err
}

// GAN General Settings

func (c *Client) GetGanGeneralSettings(ctx context.Context) (*ResourceResponse[GanGeneralSettingsConfig], error) {
	var resp ResourceResponse[GanGeneralSettingsConfig]
	// General settings usually don't have a name in the URL for list/get of singletons
	err := c.GetResource(ctx, "gateway-network-settings", "", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateGanGeneralSettings(ctx context.Context, item ResourceResponse[GanGeneralSettingsConfig]) (*ResourceResponse[GanGeneralSettingsConfig], error) {
	var resp ResourceResponse[GanGeneralSettingsConfig]
	err := c.UpdateResource(ctx, "gateway-network-settings", item, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

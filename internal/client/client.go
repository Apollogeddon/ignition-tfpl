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

type Client struct {
	HostURL    string
	HTTPClient *retryablehttp.Client
	Token      string
}

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

type ResourceResponse[T any] struct {
	Module      string `json:"module,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name"`
	Enabled     *bool  `json:"enabled,omitempty"`
	Description string `json:"description,omitempty"`
	Signature   string `json:"signature,omitempty"`
	Config      T      `json:"config"`
}

type ResourceChangesResponse struct {
	Success bool `json:"success"`
	Changes []struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Collection   string `json:"collection"`
		NewSignature string `json:"newSignature"`
	} `json:"changes"`
}

type DatabaseConfig struct {
	Driver     string `json:"driver"`
	Translator string `json:"translator,omitempty"`
	ConnectURL string `json:"connectURL"`
	Username   string `json:"username,omitempty"`
	Password   any    `json:"password,omitempty"`
}

type TagProviderConfig struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type UserSourceProfile struct {
	Type               string `json:"type"`
	FailoverProfile    string `json:"failoverProfile,omitempty"`
	FailoverMode       string `json:"failoverMode,omitempty"`
	ScheduleRestricted bool   `json:"scheduleRestricted,omitempty"`
}

type UserSourceConfig struct {
	Profile UserSourceProfile `json:"profile"`
}

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

type AuditProfileSettings struct {
	DatabaseName          string `json:"databaseName,omitempty"`
	PruneEnabled          bool   `json:"pruneEnabled,omitempty"`
	AutoCreate            bool   `json:"autoCreate,omitempty"`
	TableName             string `json:"tableName,omitempty"`
	RemoteServer          string `json:"remoteServer,omitempty"`
	RemoteProfile         string `json:"remoteProfile,omitempty"`
	EnableStoreAndForward bool   `json:"enableStoreAndForward,omitempty"`
}

type AuditProfileProfile struct {
	Type          string `json:"type"`
	RetentionDays int    `json:"retentionDays,omitempty"`
}

type AuditProfileConfig struct {
	Profile  AuditProfileProfile  `json:"profile"`
	Settings AuditProfileSettings `json:"settings"`
}

type AlarmNotificationProfileProfile struct {
	Type string `json:"type"`
}

type AlarmNotificationProfileConfig struct {
	Profile  AlarmNotificationProfileProfile `json:"profile"`
	Settings map[string]any                  `json:"settings"`
}

type OpcUaConnectionEndpoint struct {
	DiscoveryURL   string `json:"discoveryUrl"`
	EndpointURL    string `json:"endpointUrl"`
	SecurityPolicy string `json:"securityPolicy"`
	SecurityMode   string `json:"securityMode"`
}

type OpcUaConnectionSettings struct {
	Endpoint OpcUaConnectionEndpoint `json:"endpoint"`
}

type OpcUaConnectionProfile struct {
	Type string `json:"type"`
}

type OpcUaConnectionConfig struct {
	Profile  OpcUaConnectionProfile  `json:"profile"`
	Settings OpcUaConnectionSettings `json:"settings"`
}

type AlarmJournalProfile struct {
	Type      string `json:"type"`
	QueryOnly bool   `json:"queryOnly,omitempty"`
}

type AlarmJournalSettings struct {
	Datasource    string `json:"datasource,omitempty"`
	RemoteGateway *struct {
		TargetServer  string `json:"targetServer,omitempty"`
		TargetJournal string `json:"targetJournal,omitempty"`
	} `json:"remoteGateway,omitempty"`
	Advanced *struct {
		TableName          string `json:"tableName,omitempty"`
		DataTableName      string `json:"dataTableName,omitempty"`
		UseStoreAndForward bool   `json:"useStoreAndForward,omitempty"`
	} `json:"advanced,omitempty"`
	Events *struct {
		MinPriority            string `json:"minPriority,omitempty"`
		StoreShelvedEvents     bool   `json:"storeShelvedEvents,omitempty"`
		StoreFromEnabledChange bool   `json:"storeFromEnabledChange,omitempty"`
	} `json:"events,omitempty"`
}

type AlarmJournalConfig struct {
	Profile  AlarmJournalProfile  `json:"profile"`
	Settings AlarmJournalSettings `json:"settings"`
}

type SMTPProfileSettingsClassic struct {
	Hostname        string `json:"hostname"`
	Port            int    `json:"port"`
	UseSslPort      bool   `json:"useSslPort"`
	StartTlsEnabled bool   `json:"startTlsEnabled"`
	Username        string `json:"username,omitempty"`
	Password        any    `json:"password,omitempty"`
}

type SMTPProfileSettings struct {
	Settings *SMTPProfileSettingsClassic `json:"settings,omitempty"`
}

type SMTPProfileProfile struct {
	Type string `json:"type"`
}

type SMTPProfileConfig struct {
	Profile  SMTPProfileProfile  `json:"profile"`
	Settings SMTPProfileSettings `json:"settings"`
}

type StoreAndForwardMaintenancePolicy struct {
	Action string `json:"action"`
	Limit  struct {
		LimitType string `json:"limitType"`
		Value     int    `json:"value"`
	} `json:"limit"`
}

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

type IdentityProviderAuthMethod struct {
	Type   string `json:"type"`
	Config any    `json:"config"`
}

type IdentityProviderInternalConfig struct {
	UserSource               string                       `json:"userSource"`
	AuthMethods              []IdentityProviderAuthMethod `json:"authMethods"`
	SessionInactivityTimeout float64                      `json:"sessionInactivityTimeout"`
	SessionExp               float64                      `json:"sessionExp"`
	RememberMeExp            float64                      `json:"rememberMeExp"`
}

type IgnitionSecret struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

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

type IdentityProviderSamlConfig struct {
	IdpEntityId                    string   `json:"idpEntityId"`
	SpEntityId                     string   `json:"spEntityId,omitempty"`
	SpEntityIdEnabled              bool     `json:"spEntityIdEnabled"`
	AcsBinding                     string   `json:"acsBinding"`
	NameIdFormat                   string   `json:"nameIdFormat"`
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

type IdentityProviderConfig struct {
	Type   string `json:"type"`
	Config any    `json:"config"`
}

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

type RedundancyConfig struct {
	Role                string `json:"role"`
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

type DeviceConfig map[string]any

func NewClient(host, token string, allowInsecureTLS bool) (*Client, error) {
	rc := retryablehttp.NewClient()
	rc.RetryMax = 10
	rc.Logger = nil
	rc.HTTPClient.Timeout = 10 * time.Second
	
	if allowInsecureTLS {
		rc.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{HTTPClient: rc, HostURL: host, Token: token}, nil
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

// Generic Helpers

func (c *Client) GetResource(ctx context.Context, resourceType, name string, dest any) error {
	return c.GetResourceWithModule(ctx, "ignition", resourceType, name, dest)
}

func (c *Client) GetResourceWithModule(ctx context.Context, module, resourceType, name string, dest any) error {
	path := fmt.Sprintf("/data/api/v1/resources/find/%s/%s/%s", module, resourceType, name)
	if name == "" {
		path = fmt.Sprintf("/data/api/v1/resources/find/%s/%s", module, resourceType)
	}
	body, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

func (c *Client) createOrUpdate(ctx context.Context, method, module, resourceType string, item, dest any) error {
	rb, err := json.Marshal([]any{item})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/data/api/v1/resources/%s/%s", module, resourceType)
	body, err := c.doRequest(ctx, method, path, rb)
	if err != nil {
		return err
	}

	return c.unmarshalResourceResponse(ctx, module, resourceType, body, dest)
}

func (c *Client) CreateResource(ctx context.Context, resourceType string, item, dest any) error {
	return c.createOrUpdate(ctx, http.MethodPost, "ignition", resourceType, item, dest)
}

func (c *Client) UpdateResource(ctx context.Context, resourceType string, item, dest any) error {
	return c.createOrUpdate(ctx, http.MethodPut, "ignition", resourceType, item, dest)
}

func (c *Client) CreateResourceWithModule(ctx context.Context, module, resourceType string, item, dest any) error {
	return c.createOrUpdate(ctx, http.MethodPost, module, resourceType, item, dest)
}

func (c *Client) UpdateResourceWithModule(ctx context.Context, module, resourceType string, item, dest any) error {
	return c.createOrUpdate(ctx, http.MethodPut, module, resourceType, item, dest)
}

func (c *Client) DeleteResource(ctx context.Context, resourceType, name, signature string) error {
	return c.DeleteResourceWithModule(ctx, "ignition", resourceType, name, signature)
}

func (c *Client) DeleteResourceWithModule(ctx context.Context, module, resourceType, name, signature string) error {
	path := fmt.Sprintf("/data/api/v1/resources/%s/%s/%s/%s", module, resourceType, name, signature)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

func (c *Client) unmarshalResourceResponse(ctx context.Context, module, resourceType string, body []byte, dest any) error {
	if len(body) == 0 {
		return fmt.Errorf("empty response body")
	}

	var changes ResourceChangesResponse
	if err := json.Unmarshal(body, &changes); err == nil && len(changes.Changes) > 0 {
		return c.GetResourceWithModule(ctx, module, resourceType, changes.Changes[0].Name, dest)
	}

	if body[0] == '{' {
		return json.Unmarshal(body, dest)
	}

	var rawItems []json.RawMessage
	if err := json.Unmarshal(body, &rawItems); err != nil || len(rawItems) == 0 {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return json.Unmarshal(rawItems[0], dest)
}

// Simplified Typed Methods using helpers

func getR[T any](ctx context.Context, c *Client, m, t, n string) (*ResourceResponse[T], error) {
	var r ResourceResponse[T]
	err := c.GetResourceWithModule(ctx, m, t, n, &r)
	return &r, err
}

func (c *Client) GetDatabaseConnection(ctx context.Context, n string) (*ResourceResponse[DatabaseConfig], error) {
	return getR[DatabaseConfig](ctx, c, "ignition", "database-connection", n)
}
func (c *Client) CreateDatabaseConnection(ctx context.Context, i ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error) {
	var r ResourceResponse[DatabaseConfig]
	err := c.CreateResource(ctx, "database-connection", i, &r)
	return &r, err
}
func (c *Client) UpdateDatabaseConnection(ctx context.Context, i ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error) {
	var r ResourceResponse[DatabaseConfig]
	err := c.UpdateResource(ctx, "database-connection", i, &r)
	return &r, err
}
func (c *Client) DeleteDatabaseConnection(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "database-connection", n, s)
}

func (c *Client) GetUserSource(ctx context.Context, n string) (*ResourceResponse[UserSourceConfig], error) {
	return getR[UserSourceConfig](ctx, c, "ignition", "user-source", n)
}
func (c *Client) CreateUserSource(ctx context.Context, i ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error) {
	var r ResourceResponse[UserSourceConfig]
	err := c.CreateResource(ctx, "user-source", i, &r)
	return &r, err
}
func (c *Client) UpdateUserSource(ctx context.Context, i ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error) {
	var r ResourceResponse[UserSourceConfig]
	err := c.UpdateResource(ctx, "user-source", i, &r)
	return &r, err
}
func (c *Client) DeleteUserSource(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "user-source", n, s)
}

func (c *Client) GetTagProvider(ctx context.Context, n string) (*ResourceResponse[TagProviderConfig], error) {
	return getR[TagProviderConfig](ctx, c, "ignition", "tag-provider", n)
}
func (c *Client) CreateTagProvider(ctx context.Context, i ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error) {
	var r ResourceResponse[TagProviderConfig]
	err := c.CreateResource(ctx, "tag-provider", i, &r)
	return &r, err
}
func (c *Client) UpdateTagProvider(ctx context.Context, i ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error) {
	var r ResourceResponse[TagProviderConfig]
	err := c.UpdateResource(ctx, "tag-provider", i, &r)
	return &r, err
}
func (c *Client) DeleteTagProvider(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "tag-provider", n, s)
}

func (c *Client) GetAuditProfile(ctx context.Context, n string) (*ResourceResponse[AuditProfileConfig], error) {
	return getR[AuditProfileConfig](ctx, c, "ignition", "audit-profile", n)
}
func (c *Client) CreateAuditProfile(ctx context.Context, i ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error) {
	var r ResourceResponse[AuditProfileConfig]
	err := c.CreateResource(ctx, "audit-profile", i, &r)
	return &r, err
}
func (c *Client) UpdateAuditProfile(ctx context.Context, i ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error) {
	var r ResourceResponse[AuditProfileConfig]
	err := c.UpdateResource(ctx, "audit-profile", i, &r)
	return &r, err
}
func (c *Client) DeleteAuditProfile(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "audit-profile", n, s)
}

func (c *Client) GetAlarmNotificationProfile(ctx context.Context, n string) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	return getR[AlarmNotificationProfileConfig](ctx, c, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", n)
}
func (c *Client) CreateAlarmNotificationProfile(ctx context.Context, i ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	var r ResourceResponse[AlarmNotificationProfileConfig]
	err := c.CreateResourceWithModule(ctx, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", i, &r)
	return &r, err
}
func (c *Client) UpdateAlarmNotificationProfile(ctx context.Context, i ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	var r ResourceResponse[AlarmNotificationProfileConfig]
	err := c.UpdateResourceWithModule(ctx, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", i, &r)
	return &r, err
}
func (c *Client) DeleteAlarmNotificationProfile(ctx context.Context, n, s string) error {
	return c.DeleteResourceWithModule(ctx, "com.inductiveautomation.alarm-notification", "alarm-notification-profile", n, s)
}

func (c *Client) GetOpcUaConnection(ctx context.Context, n string) (*ResourceResponse[OpcUaConnectionConfig], error) {
	return getR[OpcUaConnectionConfig](ctx, c, "ignition", "opc-connection", n)
}
func (c *Client) CreateOpcUaConnection(ctx context.Context, i ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error) {
	var r ResourceResponse[OpcUaConnectionConfig]
	err := c.CreateResourceWithModule(ctx, "ignition", "opc-connection", i, &r)
	return &r, err
}
func (c *Client) UpdateOpcUaConnection(ctx context.Context, i ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error) {
	var r ResourceResponse[OpcUaConnectionConfig]
	err := c.UpdateResourceWithModule(ctx, "ignition", "opc-connection", i, &r)
	return &r, err
}
func (c *Client) DeleteOpcUaConnection(ctx context.Context, n, s string) error {
	return c.DeleteResourceWithModule(ctx, "ignition", "opc-connection", n, s)
}

func (c *Client) GetAlarmJournal(ctx context.Context, n string) (*ResourceResponse[AlarmJournalConfig], error) {
	return getR[AlarmJournalConfig](ctx, c, "ignition", "alarm-journal", n)
}
func (c *Client) CreateAlarmJournal(ctx context.Context, i ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error) {
	var r ResourceResponse[AlarmJournalConfig]
	err := c.CreateResourceWithModule(ctx, "ignition", "alarm-journal", i, &r)
	return &r, err
}
func (c *Client) UpdateAlarmJournal(ctx context.Context, i ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error) {
	var r ResourceResponse[AlarmJournalConfig]
	err := c.UpdateResourceWithModule(ctx, "ignition", "alarm-journal", i, &r)
	return &r, err
}
func (c *Client) DeleteAlarmJournal(ctx context.Context, n, s string) error {
	return c.DeleteResourceWithModule(ctx, "ignition", "alarm-journal", n, s)
}

func (c *Client) GetSMTPProfile(ctx context.Context, n string) (*ResourceResponse[SMTPProfileConfig], error) {
	return getR[SMTPProfileConfig](ctx, c, "ignition", "email-profile", n)
}
func (c *Client) CreateSMTPProfile(ctx context.Context, i ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error) {
	var r ResourceResponse[SMTPProfileConfig]
	err := c.CreateResource(ctx, "email-profile", i, &r)
	return &r, err
}
func (c *Client) UpdateSMTPProfile(ctx context.Context, i ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error) {
	var r ResourceResponse[SMTPProfileConfig]
	err := c.UpdateResource(ctx, "email-profile", i, &r)
	return &r, err
}
func (c *Client) DeleteSMTPProfile(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "email-profile", n, s)
}

func (c *Client) GetStoreAndForward(ctx context.Context, n string) (*ResourceResponse[StoreAndForwardConfig], error) {
	return getR[StoreAndForwardConfig](ctx, c, "ignition", "store-and-forward-engine", n)
}
func (c *Client) CreateStoreAndForward(ctx context.Context, i ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error) {
	var r ResourceResponse[StoreAndForwardConfig]
	err := c.CreateResource(ctx, "store-and-forward-engine", i, &r)
	return &r, err
}
func (c *Client) UpdateStoreAndForward(ctx context.Context, i ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error) {
	var r ResourceResponse[StoreAndForwardConfig]
	err := c.UpdateResource(ctx, "store-and-forward-engine", i, &r)
	return &r, err
}
func (c *Client) DeleteStoreAndForward(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "store-and-forward-engine", n, s)
}

func (c *Client) GetIdentityProvider(ctx context.Context, n string) (*ResourceResponse[IdentityProviderConfig], error) {
	return getR[IdentityProviderConfig](ctx, c, "ignition", "identity-provider", n)
}
func (c *Client) CreateIdentityProvider(ctx context.Context, i ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error) {
	var r ResourceResponse[IdentityProviderConfig]
	err := c.CreateResource(ctx, "identity-provider", i, &r)
	return &r, err
}
func (c *Client) UpdateIdentityProvider(ctx context.Context, i ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error) {
	var r ResourceResponse[IdentityProviderConfig]
	err := c.UpdateResource(ctx, "identity-provider", i, &r)
	return &r, err
}
func (c *Client) DeleteIdentityProvider(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "identity-provider", n, s)
}

func (c *Client) GetGanOutgoing(ctx context.Context, n string) (*ResourceResponse[GanOutgoingConfig], error) {
	return getR[GanOutgoingConfig](ctx, c, "ignition", "gateway-network-outgoing", n)
}
func (c *Client) CreateGanOutgoing(ctx context.Context, i ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error) {
	var r ResourceResponse[GanOutgoingConfig]
	err := c.CreateResource(ctx, "gateway-network-outgoing", i, &r)
	return &r, err
}
func (c *Client) UpdateGanOutgoing(ctx context.Context, i ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error) {
	var r ResourceResponse[GanOutgoingConfig]
	err := c.UpdateResource(ctx, "gateway-network-outgoing", i, &r)
	return &r, err
}
func (c *Client) DeleteGanOutgoing(ctx context.Context, n, s string) error {
	return c.DeleteResource(ctx, "gateway-network-outgoing", n, s)
}

func (c *Client) GetRedundancyConfig(ctx context.Context) (*RedundancyConfig, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/data/api/v1/redundancy/config", nil)
	if err != nil {
		return nil, err
	}
	var config RedundancyConfig
	return &config, json.Unmarshal(body, &config)
}
func (c *Client) UpdateRedundancyConfig(ctx context.Context, config RedundancyConfig) error {
	rb, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = c.doRequest(ctx, http.MethodPost, "/data/api/v1/redundancy/config", rb)
	return err
}

func (c *Client) GetGanGeneralSettings(ctx context.Context) (*ResourceResponse[GanGeneralSettingsConfig], error) {
	return getR[GanGeneralSettingsConfig](ctx, c, "ignition", "gateway-network-settings", "")
}
func (c *Client) UpdateGanGeneralSettings(ctx context.Context, i ResourceResponse[GanGeneralSettingsConfig]) (*ResourceResponse[GanGeneralSettingsConfig], error) {
	var r ResourceResponse[GanGeneralSettingsConfig]
	err := c.UpdateResource(ctx, "gateway-network-settings", i, &r)
	return &r, err
}

func (c *Client) GetDevice(ctx context.Context, n string) (*ResourceResponse[DeviceConfig], error) {
	return getR[DeviceConfig](ctx, c, "com.inductiveautomation.opcua", "device", n)
}
func (c *Client) CreateDevice(ctx context.Context, i ResourceResponse[DeviceConfig]) (*ResourceResponse[DeviceConfig], error) {
	var r ResourceResponse[DeviceConfig]
	err := c.CreateResourceWithModule(ctx, "com.inductiveautomation.opcua", "device", i, &r)
	return &r, err
}
func (c *Client) UpdateDevice(ctx context.Context, i ResourceResponse[DeviceConfig]) (*ResourceResponse[DeviceConfig], error) {
	var r ResourceResponse[DeviceConfig]
	err := c.UpdateResourceWithModule(ctx, "com.inductiveautomation.opcua", "device", i, &r)
	return &r, err
}
func (c *Client) DeleteDevice(ctx context.Context, n, s string) error {
	return c.DeleteResourceWithModule(ctx, "com.inductiveautomation.opcua", "device", n, s)
}

// Projects (Keep custom logic for wait)

func (c *Client) GetProject(ctx context.Context, name string) (*Project, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/data/api/v1/projects/find/"+name, nil)
	if err != nil {
		return nil, err
	}
	var p Project
	return &p, json.Unmarshal(body, &p)
}

func (c *Client) CreateProject(ctx context.Context, p Project) (*Project, error) {
	rb, _ := json.Marshal(p)
	if _, err := c.doRequest(ctx, http.MethodPost, "/data/api/v1/projects", rb); err != nil {
		return nil, err
	}
	return c.waitForProject(ctx, p.Name)
}

func (c *Client) UpdateProject(ctx context.Context, p Project) (*Project, error) {
	rb, _ := json.Marshal(p)
	if _, err := c.doRequest(ctx, http.MethodPut, "/data/api/v1/projects/"+p.Name, rb); err != nil {
		return nil, err
	}
	return c.waitForProject(ctx, p.Name)
}

func (c *Client) DeleteProject(ctx context.Context, name string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/data/api/v1/projects/"+name, nil)
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
				if p, err := c.GetProject(ctx, name); err == nil {
					return p, nil
				}
			}
		}
	}
	
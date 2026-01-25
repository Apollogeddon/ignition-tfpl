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

type IgnitionClient interface {
	GetResource(ctx context.Context, resourceType, name string, dest any) error
	CreateResource(ctx context.Context, resourceType string, item any, dest any) error
	UpdateResource(ctx context.Context, resourceType string, item any, dest any) error
	DeleteResource(ctx context.Context, resourceType, name, signature string) error
	GetResourceWithModule(ctx context.Context, module, resourceType, name string, dest any) error
	CreateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error
	UpdateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error
	DeleteResourceWithModule(ctx context.Context, module, resourceType, name, signature string) error
	EncryptSecret(ctx context.Context, plaintext string) (*IgnitionSecret, error)
	GetProject(ctx context.Context, name string) (*Project, error)
	CreateProject(ctx context.Context, p Project) (*Project, error)
	UpdateProject(ctx context.Context, p Project) (*Project, error)
	DeleteProject(ctx context.Context, name string) error
	GetDatabaseConnection(ctx context.Context, name string) (*ResourceResponse[DatabaseConfig], error)
	CreateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error)
	UpdateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error)
	DeleteDatabaseConnection(ctx context.Context, name, signature string) error
	GetUserSource(ctx context.Context, name string) (*ResourceResponse[UserSourceConfig], error)
	CreateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error)
	UpdateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error)
	DeleteUserSource(ctx context.Context, name, signature string) error
	GetTagProvider(ctx context.Context, name string) (*ResourceResponse[TagProviderConfig], error)
	CreateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error)
	UpdateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error)
	DeleteTagProvider(ctx context.Context, name, signature string) error
	GetAuditProfile(ctx context.Context, name string) (*ResourceResponse[AuditProfileConfig], error)
	CreateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error)
	UpdateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error)
	DeleteAuditProfile(ctx context.Context, name, signature string) error
	GetAlarmNotificationProfile(ctx context.Context, name string) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	CreateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	UpdateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	DeleteAlarmNotificationProfile(ctx context.Context, name, signature string) error
	GetOpcUaConnection(ctx context.Context, name string) (*ResourceResponse[OpcUaConnectionConfig], error)
	CreateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error)
	UpdateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error)
	DeleteOpcUaConnection(ctx context.Context, name, signature string) error
	GetAlarmJournal(ctx context.Context, name string) (*ResourceResponse[AlarmJournalConfig], error)
	CreateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error)
	UpdateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error)
	DeleteAlarmJournal(ctx context.Context, name, signature string) error
	GetSMTPProfile(ctx context.Context, name string) (*ResourceResponse[SMTPProfileConfig], error)
	CreateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error)
	UpdateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error)
	DeleteSMTPProfile(ctx context.Context, name, signature string) error
	GetStoreAndForward(ctx context.Context, name string) (*ResourceResponse[StoreAndForwardConfig], error)
	CreateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error)
	UpdateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error)
	DeleteStoreAndForward(ctx context.Context, name, signature string) error
	GetIdentityProvider(ctx context.Context, name string) (*ResourceResponse[IdentityProviderConfig], error)
	CreateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error)
	UpdateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error)
	DeleteIdentityProvider(ctx context.Context, name, signature string) error
	GetGanOutgoing(ctx context.Context, name string) (*ResourceResponse[GanOutgoingConfig], error)
	CreateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error)
	UpdateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error)
	DeleteGanOutgoing(ctx context.Context, name, signature string) error
	GetRedundancyConfig(ctx context.Context) (*RedundancyConfig, error)
	UpdateRedundancyConfig(ctx context.Context, config RedundancyConfig) error
	GetGanGeneralSettings(ctx context.Context) (*ResourceResponse[GanGeneralSettingsConfig], error)
	UpdateGanGeneralSettings(ctx context.Context, item ResourceResponse[GanGeneralSettingsConfig]) (*ResourceResponse[GanGeneralSettingsConfig], error)
	GetDevice(ctx context.Context, name string) (*ResourceResponse[DeviceConfig], error)
	CreateDevice(ctx context.Context, item ResourceResponse[DeviceConfig]) (*ResourceResponse[DeviceConfig], error)
	UpdateDevice(ctx context.Context, item ResourceResponse[DeviceConfig]) (*ResourceResponse[DeviceConfig], error)
	DeleteDevice(ctx context.Context, name, signature string) error
}

type Client struct {
	HostURL    string
	HTTPClient *retryablehttp.Client
	Token      string
}

func NewClient(host, token string, allowInsecureTLS bool) (*Client, error) {
	rc := retryablehttp.NewClient()
	rc.RetryMax = 10
	rc.Logger = nil
	rc.ErrorHandler = retryablehttp.PassthroughErrorHandler
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
	defer func() { _ = res.Body.Close() }()

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

func (c *Client) GetProject(ctx context.Context, name string) (*Project, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/data/api/v1/projects/find/"+name, nil)
	if err != nil {
		return nil, err
	}
	var p Project
	return &p, json.Unmarshal(body, &p)
}

func (c *Client) CreateProject(ctx context.Context, p Project) (*Project, error) {
	rb, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	if _, err := c.doRequest(ctx, http.MethodPost, "/data/api/v1/projects", rb); err != nil {
		return nil, err
	}
	return c.waitForProject(ctx, p.Name)
}

func (c *Client) UpdateProject(ctx context.Context, p Project) (*Project, error) {
	rb, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
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
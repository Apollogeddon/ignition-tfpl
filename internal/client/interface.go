package client

import "context"

// IgnitionClient defines the interface for interacting with the Ignition API.
type IgnitionClient interface {
	// Generic
	GetResource(ctx context.Context, resourceType, name string, dest any) error
	CreateResource(ctx context.Context, resourceType string, item any, dest any) error
	UpdateResource(ctx context.Context, resourceType string, item any, dest any) error
	DeleteResource(ctx context.Context, resourceType, name, signature string) error
	
	GetResourceWithModule(ctx context.Context, module, resourceType, name string, dest any) error
	CreateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error
	UpdateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error
	DeleteResourceWithModule(ctx context.Context, module, resourceType, name, signature string) error

	// Secrets
	EncryptSecret(ctx context.Context, plaintext string) (*IgnitionSecret, error)

	// Projects
	GetProject(ctx context.Context, name string) (*Project, error)
	CreateProject(ctx context.Context, p Project) (*Project, error)
	UpdateProject(ctx context.Context, p Project) (*Project, error)
	DeleteProject(ctx context.Context, name string) error

	// Database Connections
	GetDatabaseConnection(ctx context.Context, name string) (*ResourceResponse[DatabaseConfig], error)
	CreateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error)
	UpdateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error)
	DeleteDatabaseConnection(ctx context.Context, name, signature string) error

	// User Sources
	GetUserSource(ctx context.Context, name string) (*ResourceResponse[UserSourceConfig], error)
	CreateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error)
	UpdateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error)
	DeleteUserSource(ctx context.Context, name, signature string) error

	// Tag Providers
	GetTagProvider(ctx context.Context, name string) (*ResourceResponse[TagProviderConfig], error)
	CreateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error)
	UpdateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error)
	DeleteTagProvider(ctx context.Context, name, signature string) error

	// Audit Profiles
	GetAuditProfile(ctx context.Context, name string) (*ResourceResponse[AuditProfileConfig], error)
	CreateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error)
	UpdateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error)
	DeleteAuditProfile(ctx context.Context, name, signature string) error

	// Alarm Notification Profiles
	GetAlarmNotificationProfile(ctx context.Context, name string) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	CreateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	UpdateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	DeleteAlarmNotificationProfile(ctx context.Context, name, signature string) error

	// OPC UA Connections
	GetOpcUaConnection(ctx context.Context, name string) (*ResourceResponse[OpcUaConnectionConfig], error)
	CreateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error)
	UpdateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error)
	DeleteOpcUaConnection(ctx context.Context, name, signature string) error

	// Alarm Journals
	GetAlarmJournal(ctx context.Context, name string) (*ResourceResponse[AlarmJournalConfig], error)
	CreateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error)
	UpdateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error)
	DeleteAlarmJournal(ctx context.Context, name, signature string) error

	// SMTP Profiles
	GetSMTPProfile(ctx context.Context, name string) (*ResourceResponse[SMTPProfileConfig], error)
	CreateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error)
	UpdateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error)
	DeleteSMTPProfile(ctx context.Context, name, signature string) error

	// Store and Forward
	GetStoreAndForward(ctx context.Context, name string) (*ResourceResponse[StoreAndForwardConfig], error)
	CreateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error)
	UpdateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error)
	DeleteStoreAndForward(ctx context.Context, name, signature string) error

	// Identity Providers
	GetIdentityProvider(ctx context.Context, name string) (*ResourceResponse[IdentityProviderConfig], error)
	CreateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error)
	UpdateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error)
	DeleteIdentityProvider(ctx context.Context, name, signature string) error

	// Gateway Network Outgoing
	GetGanOutgoing(ctx context.Context, name string) (*ResourceResponse[GanOutgoingConfig], error)
	CreateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error)
	UpdateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error)
	DeleteGanOutgoing(ctx context.Context, name, signature string) error

	// Redundancy
	GetRedundancyConfig(ctx context.Context) (*RedundancyConfig, error)
	UpdateRedundancyConfig(ctx context.Context, config RedundancyConfig) error

	// GAN General Settings
	GetGanGeneralSettings(ctx context.Context) (*ResourceResponse[GanGeneralSettingsConfig], error)
	UpdateGanGeneralSettings(ctx context.Context, item ResourceResponse[GanGeneralSettingsConfig]) (*ResourceResponse[GanGeneralSettingsConfig], error)
}

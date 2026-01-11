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
}

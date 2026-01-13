package client

import "context"

// MockClient is a mock implementation of IgnitionClient for testing.
type MockClient struct {
	// Generic
	GetResourceFunc        func(ctx context.Context, resourceType, name string, dest any) error
	CreateResourceFunc     func(ctx context.Context, resourceType string, item any, dest any) error
	UpdateResourceFunc     func(ctx context.Context, resourceType string, item any, dest any) error
	DeleteResourceFunc     func(ctx context.Context, resourceType, name, signature string) error
	GetResourceWithModuleFunc    func(ctx context.Context, module, resourceType, name string, dest any) error
	CreateResourceWithModuleFunc func(ctx context.Context, module, resourceType string, item any, dest any) error
	UpdateResourceWithModuleFunc func(ctx context.Context, module, resourceType string, item any, dest any) error
	DeleteResourceWithModuleFunc func(ctx context.Context, module, resourceType, name, signature string) error

	// Projects
	GetProjectFunc    func(ctx context.Context, name string) (*Project, error)
	CreateProjectFunc func(ctx context.Context, p Project) (*Project, error)
	UpdateProjectFunc func(ctx context.Context, p Project) (*Project, error)
	DeleteProjectFunc func(ctx context.Context, name string) error

	// Database Connections
	GetDatabaseConnectionFunc    func(ctx context.Context, name string) (*ResourceResponse[DatabaseConfig], error)
	CreateDatabaseConnectionFunc func(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error)
	UpdateDatabaseConnectionFunc func(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error)
	DeleteDatabaseConnectionFunc func(ctx context.Context, name, signature string) error

	// User Sources
	GetUserSourceFunc    func(ctx context.Context, name string) (*ResourceResponse[UserSourceConfig], error)
	CreateUserSourceFunc func(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error)
	UpdateUserSourceFunc func(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error)
	DeleteUserSourceFunc func(ctx context.Context, name, signature string) error

	// Tag Providers
	GetTagProviderFunc    func(ctx context.Context, name string) (*ResourceResponse[TagProviderConfig], error)
	CreateTagProviderFunc func(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error)
	UpdateTagProviderFunc func(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error)
	DeleteTagProviderFunc func(ctx context.Context, name, signature string) error

	// Audit Profiles
	GetAuditProfileFunc    func(ctx context.Context, name string) (*ResourceResponse[AuditProfileConfig], error)
	CreateAuditProfileFunc func(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error)
	UpdateAuditProfileFunc func(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error)
	DeleteAuditProfileFunc func(ctx context.Context, name, signature string) error

	// Alarm Notification Profiles
	GetAlarmNotificationProfileFunc    func(ctx context.Context, name string) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	CreateAlarmNotificationProfileFunc func(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	UpdateAlarmNotificationProfileFunc func(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error)
	DeleteAlarmNotificationProfileFunc func(ctx context.Context, name, signature string) error

	// OPC UA Connections
	GetOpcUaConnectionFunc    func(ctx context.Context, name string) (*ResourceResponse[OpcUaConnectionConfig], error)
	CreateOpcUaConnectionFunc func(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error)
	UpdateOpcUaConnectionFunc func(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error)
	DeleteOpcUaConnectionFunc func(ctx context.Context, name, signature string) error

	// Alarm Journals
	GetAlarmJournalFunc    func(ctx context.Context, name string) (*ResourceResponse[AlarmJournalConfig], error)
	CreateAlarmJournalFunc func(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error)
	UpdateAlarmJournalFunc func(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error)
	DeleteAlarmJournalFunc func(ctx context.Context, name, signature string) error

	// SMTP Profiles
	GetSMTPProfileFunc    func(ctx context.Context, name string) (*ResourceResponse[SMTPProfileConfig], error)
	CreateSMTPProfileFunc func(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error)
	UpdateSMTPProfileFunc func(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error)
	DeleteSMTPProfileFunc func(ctx context.Context, name, signature string) error

	// Store and Forward
	GetStoreAndForwardFunc    func(ctx context.Context, name string) (*ResourceResponse[StoreAndForwardConfig], error)
	CreateStoreAndForwardFunc func(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error)
	UpdateStoreAndForwardFunc func(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error)
	DeleteStoreAndForwardFunc func(ctx context.Context, name, signature string) error

	// Identity Providers
	GetIdentityProviderFunc    func(ctx context.Context, name string) (*ResourceResponse[IdentityProviderConfig], error)
	CreateIdentityProviderFunc func(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error)
	UpdateIdentityProviderFunc func(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error)
	DeleteIdentityProviderFunc func(ctx context.Context, name, signature string) error

	// Gateway Network Outgoing
	GetGanOutgoingFunc    func(ctx context.Context, name string) (*ResourceResponse[GanOutgoingConfig], error)
	CreateGanOutgoingFunc func(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error)
	UpdateGanOutgoingFunc func(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error)
	DeleteGanOutgoingFunc func(ctx context.Context, name, signature string) error

	// Redundancy
	GetRedundancyConfigFunc    func(ctx context.Context) (*RedundancyConfig, error)
	UpdateRedundancyConfigFunc func(ctx context.Context, config RedundancyConfig) error

	// GAN General Settings
	GetGanGeneralSettingsFunc    func(ctx context.Context) (*ResourceResponse[GanGeneralSettingsConfig], error)
	UpdateGanGeneralSettingsFunc func(ctx context.Context, item ResourceResponse[GanGeneralSettingsConfig]) (*ResourceResponse[GanGeneralSettingsConfig], error)
}

// Generic
func (m *MockClient) GetResource(ctx context.Context, resourceType, name string, dest any) error {
	if m.GetResourceFunc != nil {
		return m.GetResourceFunc(ctx, resourceType, name, dest)
	}
	return nil
}
func (m *MockClient) CreateResource(ctx context.Context, resourceType string, item any, dest any) error {
	if m.CreateResourceFunc != nil {
		return m.CreateResourceFunc(ctx, resourceType, item, dest)
	}
	return nil
}
func (m *MockClient) UpdateResource(ctx context.Context, resourceType string, item any, dest any) error {
	if m.UpdateResourceFunc != nil {
		return m.UpdateResourceFunc(ctx, resourceType, item, dest)
	}
	return nil
}
func (m *MockClient) DeleteResource(ctx context.Context, resourceType, name, signature string) error {
	if m.DeleteResourceFunc != nil {
		return m.DeleteResourceFunc(ctx, resourceType, name, signature)
	}
	return nil
}
func (m *MockClient) GetResourceWithModule(ctx context.Context, module, resourceType, name string, dest any) error {
	if m.GetResourceWithModuleFunc != nil {
		return m.GetResourceWithModuleFunc(ctx, module, resourceType, name, dest)
	}
	return nil
}
func (m *MockClient) CreateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error {
	if m.CreateResourceWithModuleFunc != nil {
		return m.CreateResourceWithModuleFunc(ctx, module, resourceType, item, dest)
	}
	return nil
}
func (m *MockClient) UpdateResourceWithModule(ctx context.Context, module, resourceType string, item any, dest any) error {
	if m.UpdateResourceWithModuleFunc != nil {
		return m.UpdateResourceWithModuleFunc(ctx, module, resourceType, item, dest)
	}
	return nil
}
func (m *MockClient) DeleteResourceWithModule(ctx context.Context, module, resourceType, name, signature string) error {
	if m.DeleteResourceWithModuleFunc != nil {
		return m.DeleteResourceWithModuleFunc(ctx, module, resourceType, name, signature)
	}
	return nil
}

// Projects
func (m *MockClient) GetProject(ctx context.Context, name string) (*Project, error) {
	if m.GetProjectFunc != nil {
		return m.GetProjectFunc(ctx, name)
	}
	return &Project{}, nil
}
func (m *MockClient) CreateProject(ctx context.Context, p Project) (*Project, error) {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(ctx, p)
	}
	return &Project{}, nil
}
func (m *MockClient) UpdateProject(ctx context.Context, p Project) (*Project, error) {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(ctx, p)
	}
	return &Project{}, nil
}
func (m *MockClient) DeleteProject(ctx context.Context, name string) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(ctx, name)
	}
	return nil
}

// Database Connections
func (m *MockClient) GetDatabaseConnection(ctx context.Context, name string) (*ResourceResponse[DatabaseConfig], error) {
	if m.GetDatabaseConnectionFunc != nil {
		return m.GetDatabaseConnectionFunc(ctx, name)
	}
	return &ResourceResponse[DatabaseConfig]{}, nil
}
func (m *MockClient) CreateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error) {
	if m.CreateDatabaseConnectionFunc != nil {
		return m.CreateDatabaseConnectionFunc(ctx, db)
	}
	return &ResourceResponse[DatabaseConfig]{}, nil
}
func (m *MockClient) UpdateDatabaseConnection(ctx context.Context, db ResourceResponse[DatabaseConfig]) (*ResourceResponse[DatabaseConfig], error) {
	if m.UpdateDatabaseConnectionFunc != nil {
		return m.UpdateDatabaseConnectionFunc(ctx, db)
	}
	return &ResourceResponse[DatabaseConfig]{}, nil
}
func (m *MockClient) DeleteDatabaseConnection(ctx context.Context, name, signature string) error {
	if m.DeleteDatabaseConnectionFunc != nil {
		return m.DeleteDatabaseConnectionFunc(ctx, name, signature)
	}
	return nil
}

// User Sources
func (m *MockClient) GetUserSource(ctx context.Context, name string) (*ResourceResponse[UserSourceConfig], error) {
	if m.GetUserSourceFunc != nil {
		return m.GetUserSourceFunc(ctx, name)
	}
	return &ResourceResponse[UserSourceConfig]{}, nil
}
func (m *MockClient) CreateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error) {
	if m.CreateUserSourceFunc != nil {
		return m.CreateUserSourceFunc(ctx, us)
	}
	return &ResourceResponse[UserSourceConfig]{}, nil
}
func (m *MockClient) UpdateUserSource(ctx context.Context, us ResourceResponse[UserSourceConfig]) (*ResourceResponse[UserSourceConfig], error) {
	if m.UpdateUserSourceFunc != nil {
		return m.UpdateUserSourceFunc(ctx, us)
	}
	return &ResourceResponse[UserSourceConfig]{}, nil
}
func (m *MockClient) DeleteUserSource(ctx context.Context, name, signature string) error {
	if m.DeleteUserSourceFunc != nil {
		return m.DeleteUserSourceFunc(ctx, name, signature)
	}
	return nil
}

// Tag Providers
func (m *MockClient) GetTagProvider(ctx context.Context, name string) (*ResourceResponse[TagProviderConfig], error) {
	if m.GetTagProviderFunc != nil {
		return m.GetTagProviderFunc(ctx, name)
	}
	return &ResourceResponse[TagProviderConfig]{}, nil
}
func (m *MockClient) CreateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error) {
	if m.CreateTagProviderFunc != nil {
		return m.CreateTagProviderFunc(ctx, tp)
	}
	return &ResourceResponse[TagProviderConfig]{}, nil
}
func (m *MockClient) UpdateTagProvider(ctx context.Context, tp ResourceResponse[TagProviderConfig]) (*ResourceResponse[TagProviderConfig], error) {
	if m.UpdateTagProviderFunc != nil {
		return m.UpdateTagProviderFunc(ctx, tp)
	}
	return &ResourceResponse[TagProviderConfig]{}, nil
}
func (m *MockClient) DeleteTagProvider(ctx context.Context, name, signature string) error {
	if m.DeleteTagProviderFunc != nil {
		return m.DeleteTagProviderFunc(ctx, name, signature)
	}
	return nil
}

// Audit Profiles
func (m *MockClient) GetAuditProfile(ctx context.Context, name string) (*ResourceResponse[AuditProfileConfig], error) {
	if m.GetAuditProfileFunc != nil {
		return m.GetAuditProfileFunc(ctx, name)
	}
	return &ResourceResponse[AuditProfileConfig]{}, nil
}
func (m *MockClient) CreateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error) {
	if m.CreateAuditProfileFunc != nil {
		return m.CreateAuditProfileFunc(ctx, ap)
	}
	return &ResourceResponse[AuditProfileConfig]{}, nil
}
func (m *MockClient) UpdateAuditProfile(ctx context.Context, ap ResourceResponse[AuditProfileConfig]) (*ResourceResponse[AuditProfileConfig], error) {
	if m.UpdateAuditProfileFunc != nil {
		return m.UpdateAuditProfileFunc(ctx, ap)
	}
	return &ResourceResponse[AuditProfileConfig]{}, nil
}
func (m *MockClient) DeleteAuditProfile(ctx context.Context, name, signature string) error {
	if m.DeleteAuditProfileFunc != nil {
		return m.DeleteAuditProfileFunc(ctx, name, signature)
	}
	return nil
}

// Alarm Notification Profiles
func (m *MockClient) GetAlarmNotificationProfile(ctx context.Context, name string) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	if m.GetAlarmNotificationProfileFunc != nil {
		return m.GetAlarmNotificationProfileFunc(ctx, name)
	}
	return &ResourceResponse[AlarmNotificationProfileConfig]{}, nil
}
func (m *MockClient) CreateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	if m.CreateAlarmNotificationProfileFunc != nil {
		return m.CreateAlarmNotificationProfileFunc(ctx, anp)
	}
	return &ResourceResponse[AlarmNotificationProfileConfig]{}, nil
}
func (m *MockClient) UpdateAlarmNotificationProfile(ctx context.Context, anp ResourceResponse[AlarmNotificationProfileConfig]) (*ResourceResponse[AlarmNotificationProfileConfig], error) {
	if m.UpdateAlarmNotificationProfileFunc != nil {
		return m.UpdateAlarmNotificationProfileFunc(ctx, anp)
	}
	return &ResourceResponse[AlarmNotificationProfileConfig]{}, nil
}
func (m *MockClient) DeleteAlarmNotificationProfile(ctx context.Context, name, signature string) error {
	if m.DeleteAlarmNotificationProfileFunc != nil {
		return m.DeleteAlarmNotificationProfileFunc(ctx, name, signature)
	}
	return nil
}

// OPC UA Connections
func (m *MockClient) GetOpcUaConnection(ctx context.Context, name string) (*ResourceResponse[OpcUaConnectionConfig], error) {
	if m.GetOpcUaConnectionFunc != nil {
		return m.GetOpcUaConnectionFunc(ctx, name)
	}
	return &ResourceResponse[OpcUaConnectionConfig]{}, nil
}
func (m *MockClient) CreateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error) {
	if m.CreateOpcUaConnectionFunc != nil {
		return m.CreateOpcUaConnectionFunc(ctx, item)
	}
	return &ResourceResponse[OpcUaConnectionConfig]{}, nil
}
func (m *MockClient) UpdateOpcUaConnection(ctx context.Context, item ResourceResponse[OpcUaConnectionConfig]) (*ResourceResponse[OpcUaConnectionConfig], error) {
	if m.UpdateOpcUaConnectionFunc != nil {
		return m.UpdateOpcUaConnectionFunc(ctx, item)
	}
	return &ResourceResponse[OpcUaConnectionConfig]{}, nil
}
func (m *MockClient) DeleteOpcUaConnection(ctx context.Context, name, signature string) error {
	if m.DeleteOpcUaConnectionFunc != nil {
		return m.DeleteOpcUaConnectionFunc(ctx, name, signature)
	}
	return nil
}

// Alarm Journals
func (m *MockClient) GetAlarmJournal(ctx context.Context, name string) (*ResourceResponse[AlarmJournalConfig], error) {
	if m.GetAlarmJournalFunc != nil {
		return m.GetAlarmJournalFunc(ctx, name)
	}
	return &ResourceResponse[AlarmJournalConfig]{}, nil
}
func (m *MockClient) CreateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error) {
	if m.CreateAlarmJournalFunc != nil {
		return m.CreateAlarmJournalFunc(ctx, item)
	}
	return &ResourceResponse[AlarmJournalConfig]{}, nil
}
func (m *MockClient) UpdateAlarmJournal(ctx context.Context, item ResourceResponse[AlarmJournalConfig]) (*ResourceResponse[AlarmJournalConfig], error) {
	if m.UpdateAlarmJournalFunc != nil {
		return m.UpdateAlarmJournalFunc(ctx, item)
	}
	return &ResourceResponse[AlarmJournalConfig]{}, nil
}
func (m *MockClient) DeleteAlarmJournal(ctx context.Context, name, signature string) error {
	if m.DeleteAlarmJournalFunc != nil {
		return m.DeleteAlarmJournalFunc(ctx, name, signature)
	}
	return nil
}

// SMTP Profiles
func (m *MockClient) GetSMTPProfile(ctx context.Context, name string) (*ResourceResponse[SMTPProfileConfig], error) {
	if m.GetSMTPProfileFunc != nil {
		return m.GetSMTPProfileFunc(ctx, name)
	}
	return &ResourceResponse[SMTPProfileConfig]{}, nil
}
func (m *MockClient) CreateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error) {
	if m.CreateSMTPProfileFunc != nil {
		return m.CreateSMTPProfileFunc(ctx, item)
	}
	return &ResourceResponse[SMTPProfileConfig]{}, nil
}
func (m *MockClient) UpdateSMTPProfile(ctx context.Context, item ResourceResponse[SMTPProfileConfig]) (*ResourceResponse[SMTPProfileConfig], error) {
	if m.UpdateSMTPProfileFunc != nil {
		return m.UpdateSMTPProfileFunc(ctx, item)
	}
	return &ResourceResponse[SMTPProfileConfig]{}, nil
}
func (m *MockClient) DeleteSMTPProfile(ctx context.Context, name, signature string) error {
	if m.DeleteSMTPProfileFunc != nil {
		return m.DeleteSMTPProfileFunc(ctx, name, signature)
	}
	return nil
}

// Store and Forward
func (m *MockClient) GetStoreAndForward(ctx context.Context, name string) (*ResourceResponse[StoreAndForwardConfig], error) {
	if m.GetStoreAndForwardFunc != nil {
		return m.GetStoreAndForwardFunc(ctx, name)
	}
	return &ResourceResponse[StoreAndForwardConfig]{}, nil
}
func (m *MockClient) CreateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error) {
	if m.CreateStoreAndForwardFunc != nil {
		return m.CreateStoreAndForwardFunc(ctx, item)
	}
	return &ResourceResponse[StoreAndForwardConfig]{}, nil
}
func (m *MockClient) UpdateStoreAndForward(ctx context.Context, item ResourceResponse[StoreAndForwardConfig]) (*ResourceResponse[StoreAndForwardConfig], error) {
	if m.UpdateStoreAndForwardFunc != nil {
		return m.UpdateStoreAndForwardFunc(ctx, item)
	}
	return &ResourceResponse[StoreAndForwardConfig]{}, nil
}
func (m *MockClient) DeleteStoreAndForward(ctx context.Context, name, signature string) error {
	if m.DeleteStoreAndForwardFunc != nil {
		return m.DeleteStoreAndForwardFunc(ctx, name, signature)
	}
	return nil
}

// Identity Providers
func (m *MockClient) GetIdentityProvider(ctx context.Context, name string) (*ResourceResponse[IdentityProviderConfig], error) {
	if m.GetIdentityProviderFunc != nil {
		return m.GetIdentityProviderFunc(ctx, name)
	}
	return &ResourceResponse[IdentityProviderConfig]{}, nil
}
func (m *MockClient) CreateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error) {
	if m.CreateIdentityProviderFunc != nil {
		return m.CreateIdentityProviderFunc(ctx, item)
	}
	return &ResourceResponse[IdentityProviderConfig]{}, nil
}
func (m *MockClient) UpdateIdentityProvider(ctx context.Context, item ResourceResponse[IdentityProviderConfig]) (*ResourceResponse[IdentityProviderConfig], error) {
	if m.UpdateIdentityProviderFunc != nil {
		return m.UpdateIdentityProviderFunc(ctx, item)
	}
	return &ResourceResponse[IdentityProviderConfig]{}, nil
}
func (m *MockClient) DeleteIdentityProvider(ctx context.Context, name, signature string) error {
	if m.DeleteIdentityProviderFunc != nil {
		return m.DeleteIdentityProviderFunc(ctx, name, signature)
	}
	return nil
}

// Gateway Network Outgoing
func (m *MockClient) GetGanOutgoing(ctx context.Context, name string) (*ResourceResponse[GanOutgoingConfig], error) {
	if m.GetGanOutgoingFunc != nil {
		return m.GetGanOutgoingFunc(ctx, name)
	}
	return &ResourceResponse[GanOutgoingConfig]{}, nil
}
func (m *MockClient) CreateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error) {
	if m.CreateGanOutgoingFunc != nil {
		return m.CreateGanOutgoingFunc(ctx, item)
	}
	return &ResourceResponse[GanOutgoingConfig]{}, nil
}
func (m *MockClient) UpdateGanOutgoing(ctx context.Context, item ResourceResponse[GanOutgoingConfig]) (*ResourceResponse[GanOutgoingConfig], error) {
	if m.UpdateGanOutgoingFunc != nil {
		return m.UpdateGanOutgoingFunc(ctx, item)
	}
	return &ResourceResponse[GanOutgoingConfig]{}, nil
}
func (m *MockClient) DeleteGanOutgoing(ctx context.Context, name, signature string) error {
	if m.DeleteGanOutgoingFunc != nil {
		return m.DeleteGanOutgoingFunc(ctx, name, signature)
	}
	return nil
}

// Redundancy
func (m *MockClient) GetRedundancyConfig(ctx context.Context) (*RedundancyConfig, error) {
	if m.GetRedundancyConfigFunc != nil {
		return m.GetRedundancyConfigFunc(ctx)
	}
	return &RedundancyConfig{}, nil
}
func (m *MockClient) UpdateRedundancyConfig(ctx context.Context, config RedundancyConfig) error {
	if m.UpdateRedundancyConfigFunc != nil {
		return m.UpdateRedundancyConfigFunc(ctx, config)
	}
	return nil
}

// GAN General Settings
func (m *MockClient) GetGanGeneralSettings(ctx context.Context) (*ResourceResponse[GanGeneralSettingsConfig], error) {
	if m.GetGanGeneralSettingsFunc != nil {
		return m.GetGanGeneralSettingsFunc(ctx)
	}
	return &ResourceResponse[GanGeneralSettingsConfig]{}, nil
}
func (m *MockClient) UpdateGanGeneralSettings(ctx context.Context, item ResourceResponse[GanGeneralSettingsConfig]) (*ResourceResponse[GanGeneralSettingsConfig], error) {
	if m.UpdateGanGeneralSettingsFunc != nil {
		return m.UpdateGanGeneralSettingsFunc(ctx, item)
	}
	return &ResourceResponse[GanGeneralSettingsConfig]{}, nil
}

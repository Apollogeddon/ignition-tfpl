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

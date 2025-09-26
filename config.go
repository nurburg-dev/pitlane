package pitlane

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Database string
	Password string
}

func NewDBConfig(host, port, username, database, password string) *DBConfig {
	return &DBConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Database: database,
		Password: password,
	}
}

type EngineConfig struct {
	DBConfig *DBConfig
	InitDB   bool
}

func NewEngineConfig(dbc *DBConfig, initDB bool) *EngineConfig {
	return &EngineConfig{
		DBConfig: dbc,
		InitDB:   initDB,
	}
}

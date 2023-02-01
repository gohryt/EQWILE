package database

type (
	Configuration struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}

	Database struct {
		configuration Configuration
	}
)

func Constructor(configuration *Configuration) Database {
	return Database{
		configuration: *configuration,
	}
}

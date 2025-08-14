package config

const (
	EnvProduction = "production"
	EnvDebug      = "debug"
	EnvTest       = "test"
)

type AppConfig struct {
	Environment string `mapstructure:"environment"`

	Db struct {
		Driver           string `mapstructure:"driver"`
		ConnectionString string `mapstructure:"connection"`
		Log              struct {
			Level       string `mapstructure:"level"`
			Output      string `mapstructure:"output"`
			MaxFileSize int    `mapstructure:"max_file_size"`
			MaxDays     int    `mapstructure:"max_days"`
			MaxFileNum  int    `mapstructure:"max_file_num"`
		} `mapstructure:"log"`
	} `mapstructure:"db"`

	Log struct {
		Level       string `mapstructure:"level"`
		Output      string `mapstructure:"output"`
		MaxFileSize int    `mapstructure:"max_file_size"`
		MaxDays     int    `mapstructure:"max_days"`
		MaxFileNum  int    `mapstructure:"max_file_num"`
	} `mapstructure:"log"`

	Blockchain struct {
		Account struct {
			Address        string `mapstructure:"address"`
			PrivateKey     string `mapstructure:"private_key"`
			PrivateKeyFile string `mapstructure:"private_key_file"`
		} `mapstructure:"account"`
	} `mapstructure:"blockchain"`

	Relay struct {
		Api struct {
			Host string `mapstructure:"host"`
			Key  string `mapstructure:"key"`
		} `mapstructure:"api"`
	} `mapstructure:"relay"`

	Tasks struct {
		SyncTaskFeeLogs struct {
			IntervalSeconds uint `mapstructure:"interval_seconds"`
			BatchSize       uint `mapstructure:"batch_size"`
		} `mapstructure:"sync_task_fee_logs"`
	} `mapstructure:"tasks"`
}

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
		HeartbeatLogFile string `mapstructure:"heartbeat_log_file"`
		AlertLogFile     string `mapstructure:"alert_log_file"`
	} `mapstructure:"log"`

	Blockchains map[string]struct {
		RPS         uint64 `mapstructure:"rps"`
		RpcEndpoint string `mapstructure:"rpc_endpoint"`
		GasLimit    uint64 `mapstructure:"gas_limit"`
		GasPrice    uint64 `mapstructure:"gas_price"`
		ChainID     uint64 `mapstructure:"chain_id"`
		Account     struct {
			Address        string `mapstructure:"address"`
			PrivateKey     string `mapstructure:"private_key"`
			PrivateKeyFile string `mapstructure:"private_key_file"`
		} `mapstructure:"account"`
		Contracts struct {
			BenefitAddress string `mapstructure:"benefit_address"`
		} `mapstructure:"contracts"`
		MaxRetries      uint8  `mapstructure:"max_retries"`
		RetryInterval   uint64 `mapstructure:"retry_interval"`
		ReceiptWaitTime uint64 `mapstructure:"receipt_wait_time"`
	} `mapstructure:"blockchains"`

	Relay struct {
		Api struct {
			Host           string `mapstructure:"host"`
			PrivateKey     string `mapstructure:"private_key"`
			PrivateKeyFile string `mapstructure:"private_key_file"`
		} `mapstructure:"api"`
	} `mapstructure:"relay"`

	Tasks struct {
		SyncTaskFeeLogs struct {
			IntervalSeconds            uint   `mapstructure:"interval_seconds"`
			BatchSize                  uint   `mapstructure:"batch_size"`
			MaxTaskFeeAmount           uint64 `mapstructure:"max_task_fee_amount"`
			MaxAddressLogsCountInBatch uint   `mapstructure:"max_address_logs_count_in_batch"`
			MaxNewAddressCountInBatch  uint   `mapstructure:"max_new_address_count_in_batch"`
		} `mapstructure:"sync_task_fee_logs"`
		SyncWithdrawalRequests struct {
			IntervalSeconds uint `mapstructure:"interval_seconds"`
			BatchSize       uint `mapstructure:"batch_size"`
		} `mapstructure:"sync_withdrawal_requests"`
		ProcessWithdrawalRequests struct {
			IntervalSeconds uint   `mapstructure:"interval_seconds"`
			BatchSize       uint   `mapstructure:"batch_size"`
			Timeout         uint64 `mapstructure:"timeout"`
		} `mapstructure:"process_withdrawal_requests"`
	} `mapstructure:"tasks"`
}

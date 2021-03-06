package config

type Config struct {
	// Root is the path to a directory where epm will store cache data
	Root string `toml:"root"`
	// GRPC configuration settings
	GRPC GRPCConfig `toml:"grpc"`
	// DBPath is the path of a database file
	DBPath string `toml:"db_path"`
	// DBTimeout is the amount of time to wait to obtain a database file lock.
	DBTimeout int `toml:"db_timeout"`
	// CachePools stores the configurations of cache pool managers
	CachePools map[string]CachePool `toml:"cache_pools"`
}

// GRPCConfig provides GRPC configuration for the socket
type GRPCConfig struct {
	Address        string `toml:"address"`
	TCPAddress     string `toml:"tcp_address"`
	TCPTLSCert     string `toml:"tcp_tls_cert"`
	TCPTLSKey      string `toml:"tcp_tls_key"`
	UID            int    `toml:"uid"`
	GID            int    `toml:"gid"`
	MaxRecvMsgSize int    `toml:"max_recv_message_size"`
	MaxSendMsgSize int    `toml:"max_send_message_size"`
}

// CachePool provides the configuration for the cache pool
type CachePool struct {
	Type string `toml:"type"`
}

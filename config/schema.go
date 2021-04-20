package config

type Service_ struct {
	Name     string `yaml:name`
	TTL      int64  `yaml:ttl`
	Interval int64  `yaml:interval`
	Address  string `yaml:address`
}

type Logger_ struct {
	Level string `yaml:level`
	Dir   string `yaml:dir`
}

type MongoDB_ struct {
	Address  string `yaml:address`
	Timeout  int    `yaml:timeout`
	User     string `yaml:user`
	Password string `yaml:password`
	DB       string `yaml:db`
}

type Database_ struct {
	MongoDB MongoDB_ `yaml:mongodb`
}

type ConfigSchema_ struct {
	Service Service_ `yaml:service`
	Logger  Logger_  `yaml:logger`
    Database Database_ `yaml:database`
}

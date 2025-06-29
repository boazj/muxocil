package cmd

type Env string

const (
	Development Env = "dev"
	Production  Env = "prod"
)

var (
	AppName            = "muxocil"
	AppDisplayName     = "muxocil"
	AppDisplayNameLong = AppDisplayName
	AppVersion         = "0.1.0"
	Environment        = Development
)

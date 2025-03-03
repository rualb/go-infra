// Package config app config
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-infra/internal/config/consts"
	"go-infra/internal/util/utilconfig"
	xlog "go-infra/internal/util/utillog"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

var (
	AppVersion  = ""
	AppCommit   = ""
	AppDate     = ""
	ShortCommit = ""
)

func dumpVersionAndExitIf() {

	if CmdLine.Version {
		fmt.Printf("version: %s\n", AppVersion)
		fmt.Printf("commit: %s\n", AppCommit)
		fmt.Printf("date: %s\n", AppDate)
		//
		os.Exit(0)
	}

}

type CmdLineConfig struct {
	Config     string
	CertDir    string
	ConfigsDir string
	Env        string
	Name       string
	Version    bool

	SysAPIKey string
	Listen    string
	ListenTLS string
	ListenSys string

	DumpConfig bool
}

const (
	envDevelopment = "development"
	envTesting     = "testing"
	envStaging     = "staging"
	envProduction  = "production"
)

var envNames = []string{
	envDevelopment, envTesting, envStaging, envProduction,
}

var CmdLine = CmdLineConfig{}

// ReadFlags read app flags
func ReadFlags() {

	_ = os.Args
	flag.StringVar(&CmdLine.Config, "config", "", "path to dir with config files")
	flag.StringVar(&CmdLine.CertDir, "cert-dir", "", "path to dir with cert files")
	flag.StringVar(&CmdLine.SysAPIKey, "sys-api-key", "", "sys api key")
	flag.StringVar(&CmdLine.Listen, "listen", "", "listen")
	flag.StringVar(&CmdLine.ListenTLS, "listen-tls", "", "listen TLS")
	flag.StringVar(&CmdLine.ListenSys, "listen-sys", "", "listen sys")
	flag.StringVar(&CmdLine.Env, "env", "", "environment: development, testing, staging, production")
	flag.StringVar(&CmdLine.Name, "name", "", "app name")
	flag.StringVar(&CmdLine.ConfigsDir, "configs-dir", "", "path to dir with configs")

	flag.BoolVar(&CmdLine.Version, "version", false, "app version")

	flag.BoolVar(&CmdLine.DumpConfig, "dump-config", false, "dump config")

	flag.Parse() // dont use from init()

	dumpVersionAndExitIf()

}

type envReader struct {
	envError error
	prefix   string
}

func NewEnvReader() envReader {
	return envReader{prefix: "app_"}
}
func (x *envReader) readEnv(name string) string {
	envName := strings.ToUpper(x.prefix + name) // *nix case-sensitive

	{
		// APP_TITLE
		if envName != "" {
			envValue := os.Getenv(envName)
			if envValue != "" {
				xlog.Info("reading %q value from env: %v = %v", name, envName, envValue)
				return envValue
			}
		}
	}

	{
		// APP_TITLE_FILE
		envNameFile := strings.ToUpper(envName + "_file") //
		filePath := os.Getenv(envNameFile)
		if filePath != "" { // file path
			filePath = filepath.Clean(filePath)
			xlog.Info("reading %q value from file: %v = %v", name, envNameFile, filePath)
			if data, err := os.ReadFile(filePath); err == nil {
				return string(data)
			} else {
				x.envError = err
			}
		}
	}

	return ""
}

func (x *envReader) String(p *string, name string, cmdValue *string) {

	// from cmd
	if cmdValue != nil && *cmdValue != "" {
		xlog.Info("reading %q value from cmd: %v", name, *cmdValue)
		*p = *cmdValue
		return
	}

	// from env
	{
		envValue := x.readEnv(name)
		if envValue != "" {
			*p = envValue
		}
	}

}

func (x *envReader) Bool(p *bool, name string, cmdValue *bool) {

	envName := strings.ToUpper(x.prefix + name) // *nix case-sensitive

	if cmdValue != nil && *cmdValue {
		xlog.Info("reading %q value from cmd: %v", name, *cmdValue)
		*p = *cmdValue
		return
	}
	if envName != "" {
		envValue := os.Getenv(envName)
		if envValue != "" {
			xlog.Info("reading %q value from env: %v = %v", name, envName, envValue)
			*p = envValue == "1" || envValue == "true"
			return
		}
	}
}

func (x *envReader) Float64(p *float64, name string, cmdValue *float64) {

	envName := strings.ToUpper(x.prefix + name) // *nix case-sensitive

	if cmdValue != nil && math.Abs(*cmdValue) > 0.000001 {
		xlog.Info("reading float64 %q value from cmd: %v", name, *cmdValue)
		*p = *cmdValue
		return
	}

	if envName != "" {
		envValue := os.Getenv(envName)
		if envValue != "" {
			xlog.Info("reading float64 %q value from env: %v = %v", name, envName, envValue)

			if v, err := strconv.ParseFloat(envValue, 64); err == nil {
				*p = v
			} else {
				x.envError = err
			}

		}
	}

}
func (x *envReader) Int(p *int, name string, cmdValue *int) {

	envName := strings.ToUpper(x.prefix + name) // *nix case-sensitive

	if cmdValue != nil && *cmdValue != 0 {
		xlog.Info("reading %q value from cmd: %v", name, *cmdValue)
		*p = *cmdValue
		return
	}
	if envName != "" {
		envValue := os.Getenv(envName)
		if envValue != "" {
			xlog.Info("reading %q value from env: %v = %v", name, envName, envValue)

			if v, err := strconv.Atoi(envValue); err == nil {
				*p = v
			} else {
				x.envError = err
			}

		}
	}

}

type Database struct {
	Dialect   string `json:"dialect"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
	User      string `json:"user"`
	Password  string `json:"password"`
	MaxOpen   int    `json:"max_open"`
	MaxIdle   int    `json:"max_idle"`
	IdleTime  int    `json:"idle_time"`
	Migration bool   `json:"migration"`
}

// type AppConfigLog struct {
// 	Level int `json:"level"` // 0=Error 1=Warn 2=Info 3=Debug
// }

type AppConfigMessageGateway struct {
	From     string `json:"from"`
	URL      string `json:"url"`
	Query    string `json:"query"`
	Body     string `json:"body"`
	User     string `json:"credentials"`
	Password string `json:"password"`
	Stdout   bool   `json:"stdout"`
	HTTP     bool   `json:"http"`
}

type AppConfigVault struct {
	VaultAuth map[string]string `json:"auth"` // keyId:keyValue
}

type AppConfigLang struct {
	Langs []string `json:"langs"`
}

type AppConfigMod struct {
	Name  string `json:"-"`
	Env   string `json:"env"` // prod||'' dev stage
	Debug bool   `json:"-"`
	Title string `json:"title"`

	ConfigPath []string `json:"-"` // []string{".", os.Getenv("APP_CONFIG"), flagAppConfig}
}
type AppConfig struct {
	AppConfigMod

	// Log AppConfigLog `json:"logger"`

	Vault AppConfigVault `json:"vault"`

	DB    Database `json:"database"`
	Redis Database `json:"redis"`

	Lang AppConfigLang `json:"lang"`

	SmsGateway   AppConfigMessageGateway `json:"sms_gateway"`
	EmailGateway AppConfigMessageGateway `json:"email_gateway"`

	HTTPTransport AppConfigHTTPTransport `json:"http_transport"`

	HTTPServer AppConfigHTTPServer `json:"http_server"`

	Configs AppConfigConfigs `json:"configs"`
}

func NewAppConfig() *AppConfig {

	res := &AppConfig{

		Lang: AppConfigLang{Langs: []string{"en"}},
		// Log: AppConfigLog{
		// 	Level: consts.LogLevelWarn,
		// },

		DB: Database{
			Dialect:  "postgres",
			Host:     "localhost",
			Port:     "5432",
			Name:     "postgres",
			User:     "postgres",
			Password: "postgres",
			MaxOpen:  0,
			MaxIdle:  0,
			IdleTime: 0,
		},

		Redis: Database{
			Host:     "localhost",
			Port:     "6379",
			Name:     "redis",
			User:     "redis",
			Password: "redis",
		},

		AppConfigMod: AppConfigMod{
			Name:       consts.AppName,
			ConfigPath: []string{},
			Title:      "",
			Env:        "production",
			Debug:      false,
		},

		SmsGateway: AppConfigMessageGateway{
			From:     "",
			URL:      "",
			Query:    "",
			Body:     "",
			User:     "",
			Password: "",
			Stdout:   true,
			HTTP:     true,
		},

		EmailGateway: AppConfigMessageGateway{
			From:     "",
			URL:      "",
			Query:    "",
			Body:     "",
			User:     "",
			Password: "",
			Stdout:   true,
			HTTP:     true,
		},

		HTTPTransport: AppConfigHTTPTransport{},

		HTTPServer: AppConfigHTTPServer{
			ReadTimeout:  0,
			WriteTimeout: 0,
			IdleTimeout:  0,

			RateLimit: 0,
			RateBurst: 0,

			Listen: "127.0.0.1:30780",
			// ListenTLS: "127.0.0.1:30783",
			CertDir: "",

			SysAPIKey: "",
		},

		Configs: AppConfigConfigs{
			Dir: "",
		},
	}

	return res
}

func (x *AppConfig) readEnvName() error {
	reader := NewEnvReader()
	// APP_ENV -env
	reader.String(&x.Env, "env", &CmdLine.Env)
	reader.String(&x.Name, "name", &CmdLine.Name)

	if err := x.validateEnv(); err != nil {
		return err
	}

	configPath := slices.Concat(strings.Split(os.Getenv("APP_CONFIG"), ";"), strings.Split(CmdLine.Config, ";"))
	configPath = slices.Compact(configPath)
	configPath = slices.DeleteFunc(
		configPath,
		func(x string) bool {
			return x == ""
		},
	)

	for i := 0; i < len(configPath); i++ {
		configPath[i] += "/" + x.Name
	}

	// if len(configPath) == 0 {
	// 	configPath = []string{"."} // default
	// }

	if len(configPath) == 0 {
		xlog.Warn("config path is empty")
	} else {
		xlog.Info("config path: %v", configPath)
	}

	x.ConfigPath = configPath

	return nil
}

func (x *AppConfig) readEnvVar() error {
	reader := NewEnvReader()

	// SmsGateway configuration
	reader.String(&x.SmsGateway.From, "sms_gw_from", nil)
	reader.String(&x.SmsGateway.URL, "sms_gw_url", nil)
	reader.String(&x.SmsGateway.Query, "sms_gw_query", nil)
	reader.String(&x.SmsGateway.Body, "sms_gw_body", nil)
	reader.String(&x.SmsGateway.User, "sms_gw_user", nil)
	reader.String(&x.SmsGateway.Password, "sms_gw_password", nil)
	reader.Bool(&x.SmsGateway.Stdout, "sms_gw_stdout", nil)
	reader.Bool(&x.SmsGateway.HTTP, "sms_gw_http", nil)

	// EmailGateway configuration
	reader.String(&x.EmailGateway.From, "email_gw_from", nil)
	reader.String(&x.EmailGateway.URL, "email_gw_url", nil)
	reader.String(&x.EmailGateway.Query, "email_gw_query", nil)
	reader.String(&x.EmailGateway.Body, "email_gw_body", nil)
	reader.String(&x.EmailGateway.User, "email_gw_user", nil)
	reader.String(&x.EmailGateway.Password, "email_gw_password", nil)
	reader.Bool(&x.EmailGateway.Stdout, "email_gw_stdout", nil)
	reader.Bool(&x.EmailGateway.HTTP, "email_gw_http", nil)

	// Database configuration

	reader.String(&x.DB.Dialect, "db_dialect", nil)
	reader.String(&x.DB.Host, "db_host", nil)
	reader.String(&x.DB.Port, "db_port", nil)
	reader.String(&x.DB.Name, "db_name", nil)
	reader.String(&x.DB.User, "db_user", nil)
	reader.String(&x.DB.Password, "db_password", nil)
	reader.Int(&x.DB.MaxOpen, "db_max_open", nil)
	reader.Int(&x.DB.MaxIdle, "db_max_idle", nil)
	reader.Int(&x.DB.IdleTime, "db_idle_time", nil)
	reader.Bool(&x.DB.Migration, "db_migration", nil)

	// General configuration
	reader.String(&x.Title, "title", nil)

	// Http server
	reader.Bool(&x.HTTPServer.AccessLog, "http_access_log", nil)
	reader.Float64(&x.HTTPServer.RateLimit, "http_rate_limit", nil)
	reader.Int(&x.HTTPServer.RateBurst, "http_rate_burst", nil)
	reader.String(&x.HTTPServer.Listen, "http_listen", nil)        // =>listen
	reader.String(&x.HTTPServer.ListenTLS, "http_listen_tls", nil) // =>listen_tls
	reader.Bool(&x.HTTPServer.AutoTLS, "http_auto_tls", nil)
	reader.Bool(&x.HTTPServer.RedirectHTTPS, "http_redirect_https", nil)
	reader.Bool(&x.HTTPServer.RedirectWWW, "http_redirect_www", nil)
	reader.String(&x.HTTPServer.CertDir, "http_cert_dir", &CmdLine.CertDir) // =>cert_dir
	reader.Int(&x.HTTPServer.ReadTimeout, "http_read_timeout", nil)
	reader.Int(&x.HTTPServer.WriteTimeout, "http_write_timeout", nil)
	reader.Int(&x.HTTPServer.IdleTimeout, "http_idle_timeout", nil)
	reader.Int(&x.HTTPServer.ReadHeaderTimeout, "http_read_header_timeout", nil)
	reader.String(&x.HTTPServer.ListenSys, "http_listen_sys", nil)  // =>listen_sys
	reader.String(&x.HTTPServer.SysAPIKey, "http_sys_api_key", nil) // =>sys_api_key

	reader.String(&x.HTTPServer.CertDir, "cert_dir", &CmdLine.CertDir) // short
	reader.String(&x.Configs.Dir, "configs_dir", &CmdLine.ConfigsDir)

	reader.String(&x.HTTPServer.Listen, "listen", &CmdLine.Listen)
	reader.String(&x.HTTPServer.ListenTLS, "listen_tls", &CmdLine.ListenTLS)
	reader.String(&x.HTTPServer.ListenSys, "listen_sys", &CmdLine.ListenSys)

	reader.String(&x.HTTPServer.SysAPIKey, "sys_api_key", &CmdLine.SysAPIKey)

	if reader.envError != nil {
		return reader.envError
	}

	return nil
}

func (x *AppConfig) validateEnv() error {

	if x.Env == "" {
		x.Env = envProduction
	}

	x.Debug = x.Env == envDevelopment
	if !slices.Contains(envNames, x.Env) {
		xlog.Warn("non-standart env name: %v", x.Env)
	}

	return nil

}
func (x AppConfig) validate() error {

	if x.HTTPServer.Listen == "" && x.HTTPServer.ListenTLS == "" {
		return fmt.Errorf("socket Listen and ListenTLS are empty")
	}

	return nil
}

type AppConfigSource struct {
	config *AppConfig
}

func MustNewAppConfigSource() *AppConfigSource {

	res := &AppConfigSource{}

	err := res.Load() // init-load

	if err != nil {
		panic(err)
	}

	return res

}

type AppConfigHTTPTransport struct {
	MaxIdleConns        int `json:"max_idle_conns,omitempty"`
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host,omitempty"`
	IdleConnTimeout     int `json:"idle_conn_timeout,omitempty"`
	MaxConnsPerHost     int `json:"max_conns_per_host,omitempty"`
}

type AppConfigHTTPServer struct {
	AccessLog     bool    `json:"access_log"`
	RateLimit     float64 `json:"rate_limit"`
	RateBurst     int     `json:"rate_burst"`
	Listen        string  `json:"listen"`
	ListenTLS     string  `json:"listen_tls"`
	AutoTLS       bool    `json:"auto_tls"`
	RedirectHTTPS bool    `json:"redirect_https"`
	RedirectWWW   bool    `json:"redirect_www"`

	CertDir string `json:"cert_dir"`

	ReadTimeout       int `json:"read_timeout,omitempty"`        // 5 to 30 seconds
	WriteTimeout      int `json:"write_timeout,omitempty"`       // 10 to 30 seconds, WriteTimeout > ReadTimeout
	IdleTimeout       int `json:"idle_timeout,omitempty"`        // 60 to 120 seconds
	ReadHeaderTimeout int `json:"read_header_timeout,omitempty"` // default get from ReadTimeout

	SysMetrics bool   `json:"sys_metrics"` //
	SysAPIKey  string `json:"sys_api_key"`
	ListenSys  string `json:"listen_sys"`
}
type AppConfigConfigs struct {
	Dir string `json:"dir"`
}

// Load load config
func (x *AppConfigSource) Load() error {

	res := NewAppConfig()

	{
		err := res.readEnvName()
		if err != nil {
			return err
		}
	}

	{
		for i := 0; i < len(res.ConfigPath); i++ {

			dir := res.ConfigPath[i]

			fileName := fmt.Sprintf("config.%s.json", res.Env)

			xlog.Info("loading config from: %v", dir)

			err := utilconfig.LoadConfig(res /*pointer*/, dir, fileName)

			if err != nil {
				return err
			}

		}

	}

	{
		err := res.readEnvVar()
		if err != nil {
			return err
		}

	}

	{
		err := res.validate()
		if err != nil {
			return err
		}
	}

	xlog.Info("config loaded: Name=%v Env=%v Debug=%v ", res.Name, res.Env, res.Debug)

	x.config = res

	if CmdLine.DumpConfig {
		data, _ := json.MarshalIndent(res, "", " ")
		fmt.Println(string(data))
	}

	return nil
}

func (x *AppConfigSource) Config() *AppConfig {

	return x.config

}

// FromJSON from json
func (x *AppConfig) FromJSON(data string) error {

	if data == "" {
		return nil
	}

	err := json.Unmarshal([]byte(data), x)

	if err != nil {
		return err
	}

	return nil
}

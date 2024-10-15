package configs

import (
	"bytes"
	_ "embed"
	"io"

	"gin-api-admin/internal/pkg/env"

	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	MySQL struct {
		Read struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"read"`
		Write struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"write"`
	} `toml:"mysql"`

	Redis struct {
		Addr string `toml:"addr"`
		Pass string `toml:"pass"`
		Db   int    `toml:"db"`
	} `toml:"redis"`

	Mongo struct {
		URI        string `toml:"uri"`
		UserName   string `toml:"username"`
		Password   string `toml:"password"`
		AuthSource string `toml:"authSource"`
	} `toml:"mongo"`

	JWT struct {
		Secret string `toml:"secret"`
	} `toml:"jwt"`

	AES struct {
		Secret string `toml:"secret"`
	} `toml:"aes"`

	RSA struct {
		PublicKey  string `toml:"publicKey"`
		PrivateKey string `toml:"privateKey"`
	} `toml:"rsa"`

	Language struct {
		Local string `toml:"local"`
	} `toml:"language"`
}

var (
	//go:embed dev_configs.toml
	devConfigs []byte

	//go:embed fat_configs.toml
	fatConfigs []byte

	//go:embed uat_configs.toml
	uatConfigs []byte

	//go:embed pro_configs.toml
	proConfigs []byte
)

func init() {
	var r io.Reader

	switch env.Active().Value() {
	case "dev":
		r = bytes.NewReader(devConfigs)
	case "fat":
		r = bytes.NewReader(fatConfigs)
	case "uat":
		r = bytes.NewReader(uatConfigs)
	case "pro":
		r = bytes.NewReader(proConfigs)
	default:
		r = bytes.NewReader(fatConfigs)
	}

	viper.SetConfigType("toml")

	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}
}

func Get() Config {
	return *config
}

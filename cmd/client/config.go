package client

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/transfer"
)

type clientConfig struct {
	ProxyAddr []string      `toml:"proxyAddr"`
	Control   controlConfig `toml:"control"`
	Http      proxyConfig   `toml:"http"`
	Socks5    proxyConfig   `toml:"socks5"`
}

func (c *clientConfig) Addr(t transfer.Type) string {
	switch t {
	case transfer.HTTP:
		return c.Http.ListenAddr
	case transfer.SOCKS5:
		return c.Socks5.ListenAddr
	}
	return ""
}

type controlConfig struct {
	ConfigAddr string `toml:"addr"`
	LogPath    string `toml:"logPath"`
}

type proxyConfig struct {
	Enabled    bool   `toml:"enabled"`
	ListenAddr string `toml:"listenAddr"`
	Status     string `toml:"status"`
}

var (
	singleton  *clientConfig
	configPath = ".seamoom"
)

func Config() *clientConfig {
	if singleton == nil {
		singleton = defaultConfig()
	}
	return singleton
}

func defaultConfig() *clientConfig {
	return &clientConfig{
		ProxyAddr: []string{""},
		Control: controlConfig{
			// is dangerous to open Control page for everyone, do not set value like: ":7777" / "0.0.0.0:7777"
			ConfigAddr: ":7777",
			LogPath:    "web.log",
		},
		Http: proxyConfig{
			false, ":9000", "inactive",
		},
		Socks5: proxyConfig{
			false, ":1080", "inactive",
		},
	}
}

func (c *clientConfig) Save() error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(c); err != nil {
		return err
	}

	fd, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	if err = fd.Truncate(0); err != nil {
		return err
	}

	if _, err = fd.Seek(0, 0); err != nil {
		return err
	}

	_, err = fd.Write(buf.Bytes())
	return err
}

func (c *clientConfig) Load(sg *SigGroup) error {
	if data, err := os.ReadFile(configPath); err == nil {
		// try first start
		err := toml.Unmarshal(data, c)
		if err != nil {
			slog.Debug(consts.CONFIG_NOT_FIND)
			return err
		}
		sg.Detection()
		return nil
	} else {
		return err
	}
}

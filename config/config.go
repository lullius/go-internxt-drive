package config

import (
	"encoding/json"
	"os"
)

const (
	DefaultDriveAPIURL      = "https://api.internxt.com/drive"
	DefaultAuthAPIURL       = "https://api.internxt.com/drive/auth"
	DefaultUsersAPIURL      = "https://api.internxt.com/users"
	DefaultAppCryptoSecret  = "6KYQBP847D4ATSFA"
	DefaultAppCryptoSecret2 = "8Q8VMUE3BJZV87GT"
	DefaultAppMagicIV       = "d139cb9a2cd17092e79e1861cf9d7023"
	DefaultAppMagicSalt     = "38dce0391b49efba88dbc8c39ebf868f0267eb110bb0012ab27dc52a528d61b1d1ed9d76f400ff58e3240028442b1eab9bb84e111d9dadd997982dbde9dbd25e"
)

type Config struct {
	Email            string `json:"email,omitempty"`
	Password         string `json:"password,omitempty"`
	TFA              string `json:"tfa,omitempty"`
	Token            string `json:"token,omitempty"`
	RootFolderID     string `json:"root_folder_id,omitempty"`
	Bucket           string `json:"bucket,omitempty"`
	Mnemonic         string `json:"mnemonic,omitempty"`
	BasicAuthHeader  string `json:"basic_auth_header,omitempty"`
	DriveAPIURL      string `json:"drive_api_url,omitempty"`
	AuthAPIURL       string `json:"auth_api_url,omitempty"`
	UsersAPIURL      string `json:"users_api_url,omitempty"`
	AppCryptoSecret  string `json:"app_crypto_secret,omitempty"`
	AppCryptoSecret2 string `json:"app_crypto_secret2,omitempty"`
	AppMagicIV       string `json:"app_magic_iv,omitempty"`
	AppMagicSalt     string `json:"app_magic_salt,omitempty"`
}

func NewDefault(email, password string) *Config {
	cfg := &Config{
		Email:    email,
		Password: password,
	}
	cfg.applyDefaults()
	return cfg
}

func NewDefault2FA(email, password, tfa string) *Config {
	cfg := &Config{
		Email:    email,
		Password: password,
		TFA:      tfa,
	}
	cfg.applyDefaults()
	return cfg
}

func NewDefaultToken(token string) *Config {
	cfg := &Config{
		Token: token,
	}
	cfg.applyDefaults()
	return cfg
}

func LoadFromJSON(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) SaveToJSON(path string) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, os.ModePerm)
}

func (c *Config) applyDefaults() {
	if c.DriveAPIURL == "" {
		c.DriveAPIURL = DefaultDriveAPIURL
	}
	if c.AuthAPIURL == "" {
		c.AuthAPIURL = DefaultAuthAPIURL
	}
	if c.UsersAPIURL == "" {
		c.UsersAPIURL = DefaultUsersAPIURL
	}
	if c.AppCryptoSecret == "" {
		c.AppCryptoSecret = DefaultAppCryptoSecret
	}
	if c.AppCryptoSecret2 == "" {
		c.AppCryptoSecret2 = DefaultAppCryptoSecret2
	}
	if c.AppMagicIV == "" {
		c.AppMagicIV = DefaultAppMagicIV
	}
	if c.AppMagicSalt == "" {
		c.AppMagicSalt = DefaultAppMagicSalt
	}
}

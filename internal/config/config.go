package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultTimeout     = "30s"
	defaultRetry       = 3
	defaultConcurrency = 4
)

type Profile struct {
	BaseURL     string `mapstructure:"base_url" yaml:"base_url"`
	APIKey      string `mapstructure:"api_key" yaml:"api_key"`
	Timeout     string `mapstructure:"timeout" yaml:"timeout"`
	Retry       int    `mapstructure:"retry" yaml:"retry"`
	Concurrency int    `mapstructure:"concurrency" yaml:"concurrency"`
}

type FileConfig struct {
	Profiles map[string]Profile `mapstructure:"profiles" yaml:"profiles"`
	Current  string             `mapstructure:"current" yaml:"current"`
}

type Manager struct {
	path string
	cfg  FileConfig
}

type Overrides struct {
	ProfileName        string
	BaseURL            string
	APIKey             string
	Timeout            string
	Retry              int
	Concurrency        int
	HasRetry           bool
	HasConcurrency     bool
	HasExplicitProfile bool
}

func DefaultPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir failed: %w", err)
	}
	return filepath.Join(base, "gradmotion", "config.yaml"), nil
}

func NewManager(path string) (*Manager, error) {
	if strings.TrimSpace(path) == "" {
		p, err := DefaultPath()
		if err != nil {
			return nil, err
		}
		path = p
	}
	return &Manager{
		path: path,
		cfg: FileConfig{
			Profiles: map[string]Profile{},
		},
	}, nil
}

func (m *Manager) Path() string {
	return m.path
}

func (m *Manager) Load() error {
	if _, err := os.Stat(m.path); errors.Is(err, os.ErrNotExist) {
		m.cfg = defaultConfig()
		return m.Save()
	}

	v := viper.New()
	v.SetConfigFile(m.path)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config failed: %w", err)
	}

	var cfg FileConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("parse config failed: %w", err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = map[string]Profile{}
	}
	if strings.TrimSpace(cfg.Current) == "" {
		cfg.Current = "prod"
	}
	if _, ok := cfg.Profiles[cfg.Current]; !ok {
		cfg.Profiles[cfg.Current] = defaultProfile()
	}
	m.cfg = cfg
	return nil
}

func (m *Manager) Save() error {
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config dir failed: %w", err)
	}

	v := viper.New()
	v.Set("profiles", m.cfg.Profiles)
	v.Set("current", m.cfg.Current)
	v.SetConfigFile(m.path)
	if err := v.WriteConfigAs(m.path); err != nil {
		return fmt.Errorf("write config failed: %w", err)
	}
	return nil
}

func (m *Manager) CurrentProfileName() string {
	if strings.TrimSpace(m.cfg.Current) == "" {
		return "prod"
	}
	return m.cfg.Current
}

func (m *Manager) SetCurrentProfileName(name string) error {
	if _, ok := m.cfg.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	m.cfg.Current = name
	return nil
}

func (m *Manager) ListProfiles() map[string]Profile {
	out := map[string]Profile{}
	for k, v := range m.cfg.Profiles {
		out[k] = normalizeProfile(v)
	}
	return out
}

func (m *Manager) GetProfile(name string) (Profile, bool) {
	p, ok := m.cfg.Profiles[name]
	if !ok {
		return Profile{}, false
	}
	return normalizeProfile(p), true
}

func (m *Manager) UpsertProfile(name string, p Profile) {
	m.cfg.Profiles[name] = normalizeProfile(p)
}

func (m *Manager) UpdateCurrentProfile(mutator func(*Profile) error) error {
	name := m.CurrentProfileName()
	p, _ := m.GetProfile(name)
	if err := mutator(&p); err != nil {
		return err
	}
	m.UpsertProfile(name, p)
	return nil
}

func (m *Manager) EffectiveProfile(ov Overrides) (string, Profile) {
	name := m.CurrentProfileName()
	if ov.HasExplicitProfile && strings.TrimSpace(ov.ProfileName) != "" {
		name = strings.TrimSpace(ov.ProfileName)
	}

	p, ok := m.GetProfile(name)
	if !ok {
		p = defaultProfile()
	}

	// env overrides
	if v := strings.TrimSpace(os.Getenv("GM_BASE_URL")); v != "" {
		p.BaseURL = v
	}
	if v := strings.TrimSpace(os.Getenv("GM_API_KEY")); v != "" {
		p.APIKey = v
	}
	if v := strings.TrimSpace(os.Getenv("GM_TIMEOUT")); v != "" {
		p.Timeout = v
	}
	if v := strings.TrimSpace(os.Getenv("GM_RETRY")); v != "" {
		if iv, err := strconv.Atoi(v); err == nil {
			p.Retry = iv
		}
	}
	if v := strings.TrimSpace(os.Getenv("GM_CONCURRENCY")); v != "" {
		if iv, err := strconv.Atoi(v); err == nil {
			p.Concurrency = iv
		}
	}

	// flag overrides
	if strings.TrimSpace(ov.BaseURL) != "" {
		p.BaseURL = strings.TrimSpace(ov.BaseURL)
	}
	if strings.TrimSpace(ov.APIKey) != "" {
		p.APIKey = strings.TrimSpace(ov.APIKey)
	}
	if strings.TrimSpace(ov.Timeout) != "" {
		p.Timeout = strings.TrimSpace(ov.Timeout)
	}
	if ov.HasRetry {
		p.Retry = ov.Retry
	}
	if ov.HasConcurrency {
		p.Concurrency = ov.Concurrency
	}

	return name, normalizeProfile(p)
}

func defaultConfig() FileConfig {
	return FileConfig{
		Profiles: map[string]Profile{
			"prod": defaultProfile(),
		},
		Current: "prod",
	}
}

func defaultProfile() Profile {
	return Profile{
		BaseURL:     "http://8.141.22.122:9096",
		Timeout:     defaultTimeout,
		Retry:       defaultRetry,
		Concurrency: defaultConcurrency,
	}
}

func normalizeProfile(p Profile) Profile {
	if strings.TrimSpace(p.Timeout) == "" {
		p.Timeout = defaultTimeout
	}
	if p.Retry <= 0 {
		p.Retry = defaultRetry
	}
	if p.Concurrency <= 0 {
		p.Concurrency = defaultConcurrency
	}
	return p
}

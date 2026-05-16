package config

type ModuleConfig struct {
	LoginPath   string `mapstructure:"login_path" json:"login_path"`
	RefreshPath string `mapstructure:"refresh_path" json:"refresh_path"`
	LogoutPath  string `mapstructure:"logout_path" json:"logout_path"`
}

func (c *ModuleConfig) SetEnvBindings() map[string]string {
	return map[string]string{
		"module.auth.login_path":   "MODULE_AUTH_LOGIN_PATH",
		"module.auth.refresh_path": "MODULE_AUTH_REFRESH_PATH",
		"module.auth.logout_path":  "MODULE_AUTH_LOGOUT_PATH",
	}
}

func (c *ModuleConfig) SetDefaults() map[string]any {
	return map[string]any{
		"module.auth.login_path":   "/auth/token",
		"module.auth.refresh_path": "/auth/refresh",
		"module.auth.logout_path":  "/auth/logout",
	}
}

package auth

import (
	"fmt"
	"time"

	"github.com/KaderKita/data-registry-mod-middleware/config"
	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/adapter/auth/authn"
	"github.com/webcore-go/webcore/app/core"
	"github.com/webcore-go/webcore/app/out"
	appConfig "github.com/webcore-go/webcore/infra/config"
	"github.com/webcore-go/webcore/infra/logger"
)

const (
	ModuleName    = "auth"
	ModuleVersion = "1.0.0"
)

type Module struct {
	config *config.ModuleConfig
	routes []*core.ModuleRoute
	authn  *authn.AuthN
}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return ModuleName
}

func (m *Module) Version() string {
	return ModuleVersion
}

func (m *Module) Dependencies() []string {
	return []string{}
}

func (m *Module) Health(c *fiber.Ctx) error {
	health := map[string]any{
		"status":    "healthy",
		"module":    ModuleName,
		"version":   ModuleVersion,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	return c.JSON(health)
}

func (m *Module) Info(c *fiber.Ctx) error {
	endpoints := []string{}
	for _, endpoint := range m.routes {
		endpoints = append(endpoints, endpoint.Method+" "+endpoint.Path)
	}

	info := map[string]any{
		"name":        ModuleName,
		"version":     ModuleVersion,
		"description": "Modul Login",
		"path":        "/" + ModuleName,
		"endpoints":   endpoints,
		"config":      m.config,
	}
	return c.JSON(info)
}

func (m *Module) Init(ctx *core.AppContext) error {
	m.config = &config.ModuleConfig{}
	if err := appConfig.LoadDefaultConfigModule(m.Name(), m.config); err != nil {
		return err
	}

	if lib, ok := ctx.GetDefaultSingletonInstance("authentication"); ok {
		authn := lib.(*authn.AuthN)
		m.authn = authn
	} else {
		return fmt.Errorf("authentication instance not found")
	}

	m.registerRoutes(ctx.Web)
	logger.Info("Module Auth initialized successfully")

	return nil
}

func (m *Module) Destroy() error {
	return nil
}

func (m *Module) Config() appConfig.Configurable {
	return m.config
}

func (m *Module) Routes() []*core.ModuleRoute {
	return m.routes
}

func (m *Module) Services() map[string]any {
	return map[string]any{
		// "service": nil,
	}
}

func (m *Module) Repositories() map[string]any {
	return map[string]any{
		// "repository": nil,
	}
}

func (m *Module) registerRoutes(root *fiber.App) {
	moduleRoot := root.Group("/" + m.Name())

	m.routes = core.AppendRouteToArray(m.routes, &core.ModuleRoute{
		Method:  "POST",
		Path:    "/token",
		Handler: m.AuthLogin,
		Root:    moduleRoot,
	})

	m.routes = core.AppendRouteToArray(m.routes, &core.ModuleRoute{
		Method:  "POST",
		Path:    "/refresh",
		Handler: m.AuthTokenRefresh,
		Root:    moduleRoot,
	})

	m.routes = core.AppendRouteToArray(m.routes, &core.ModuleRoute{
		Method:  "POST",
		Path:    "/logout",
		Handler: m.AuthLogout,
		Root:    moduleRoot,
	})

	m.routes = core.AppendRouteToArray(m.routes, &core.ModuleRoute{
		Method:  "GET",
		Path:    "/health",
		Handler: m.Health,
		Root:    moduleRoot,
	})

	m.routes = core.AppendRouteToArray(m.routes, &core.ModuleRoute{
		Method:  "GET",
		Path:    "/info",
		Handler: m.Info,
		Root:    moduleRoot,
	})
}

func (m *Module) AuthLogin(c *fiber.Ctx) error {
	userInfo, err := m.authn.Authenticator.Login(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(out.ErrorDetail(
			fiber.StatusBadRequest, 1, "LOGIN_FAILED", "Invalid login", err))
	}

	return c.JSON(userInfo)
}

func (m *Module) AuthTokenRefresh(c *fiber.Ctx) error {
	userInfo, err := m.authn.Authenticator.RefreshToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(out.ErrorDetail(
			fiber.StatusBadRequest, 1, "LOGIN_FAILED", "Invalid login", err))
	}

	return c.JSON(userInfo)
}

func (m *Module) AuthLogout(c *fiber.Ctx) error {
	err := m.authn.Authenticator.Logout(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(out.ErrorDetail(
			fiber.StatusBadRequest, 1, "LOGIN_FAILED", "Invalid login", err))
	}

	return c.JSON(out.SuccessMessage("Logout berhasil"))
}

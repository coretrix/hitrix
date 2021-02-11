package hitrix

import "github.com/sarulabs/di"

const ModeLocal = "local"
const ModeTest = "test"
const ModeProd = "prod"

type AppDefinition struct {
	mode     string
	name     string
	flags    *Flags
	registry *Registry
}

func (app *AppDefinition) Name() string {
	return app.name
}

func (app *AppDefinition) Mode() string {
	return app.mode
}

func (app *AppDefinition) Flags() *Flags {
	return app.flags
}

func (app *AppDefinition) IsInLocalMode() bool {
	return app.IsInMode(ModeLocal)
}

func (app *AppDefinition) IsInTestMode() bool {
	return app.IsInMode(ModeTest)
}

func (app *AppDefinition) IsInProdMode() bool {
	return app.IsInMode(ModeProd)
}
func (app *AppDefinition) IsInMode(mode string) bool {
	return app.mode == mode
}

func serviceApp(app *AppDefinition) *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "app",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return app, nil
		},
	}
}

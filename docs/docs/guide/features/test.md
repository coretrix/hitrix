# Tests
Hitrix provide you test abstract layer that can be used to simulate requests to your api

In your code you can create similar function that makes new instance of your app

```go
func createContextMyApp(t *testing.T, projectName string) *test.Ctx {
	defaultServices := []*service.Definition{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderOrmEngine(),
		//your services here
	}

	return test.CreateContext(t,
		projectName,
		resolvers,
		func(ginEngine *gin.Engine) { middleware.Router(ginEngine) },
		defaultServices,
	)
}

```

Hitrix supports `parallel` tests
In case you want to execute parallel tests you need to set
`PARALLEL_TESTS=true` env var in your IDE config and be sure you don't have set `-p 1` in `Go tool arguments`
In case you want to disable `parallel` tests remove `PARALLEL_TESTS` or set it to `false` and set in `Go tool arguments` value `-p 1`


# Tests
Hitrix provide you test abstract layer that can be used to simulate requests to your graphql api

In your code you can create similar function that makes new instance of your app

```go
func createContextMyApp(t *testing.T, projectName string, resolvers graphql.ExecutableSchema) *test.Ctx {
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

After that you can call queries or mutations

```go
func TestProcessApplePurchaseWithEmail(t *testing.T) {
	type queryRegisterTransactions struct {
		RegisterTransactionsResponse *model.RegisterTransactionsResponse `graphql:"RegisterTransactions(applePurchaseRequest: $applePurchaseRequest)"`
	}

	variables := map[string]interface{}{
		"applePurchaseRequest": model.ApplePurchaseRequest{
			ForceEmail:   false,
		},
	}

	fakeMail := &mailMock.Sender{}
	fakeMail.On("SendTemplate", "hymn@abv.bg").Return(nil)

	got := &queryRegisterTransactions{}
	projectName, resolver := tests.GetWebAPIResolver()
	ctx := tests.CreateContextWebAPI(t, projectName, resolver, &tests.IoCMocks{MailService: fakeMail})

	err := ctx.HandleMutation(got, variables)
	assert.Nil(t, err)

	//...
	fakeMail.AssertExpectations(t)
}
```

Hitrix supports `parallel` tests
In case you want to execute parallel tests you need to set
`PARALLEL_TESTS=true` env var in your IDE config and be sure you don't have set `-p 1` in `Go tool arguments`
In case you want to disable `parallel` tests remove `PARALLEL_TESTS` or set it to `false` and set in `Go tool arguments` value `-p 1`


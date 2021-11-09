# Feature flag service
This service provides you ability to enable and disable different features into your platform

Register the service into your `main.go` file:
```go 
registry.ServiceProviderFeatureFlag()
```

Access the service:
```go
service.DI().FeatureFlag()
```
Please have in your mind that you need to register FeatureFlagEntity

At start of the service hitrix is syncing the feature_flag table with the registered features.
Automatically make inserts and updates if it's needed. The service never remove from the table!

Also you are able to activate and deactivate the features using our dev panel

# Use case

In case you want to enable/disable the whole resolver you can do it in that way
```go
package graph
import (
	"errors"
	"service"
)

func (r *attributeResolver) Values(ctx context.Context, obj *model.Attribute) ([]*model.AttributeValue, error) {
	ormService := service.DI().ORMEngineFromContext(ctx)
	err := service.DI().FeatureFlag().FailIfIsNotActive(ormService, "bundle")
    if err != nil {
		return nil, err
    }
	
	return attribute.ValuesWeb(ctx, obj)
}
```

In case you want to chek if feature is enabled/disabled somewhere in your logic you can use this method

```go
package login
import (
    "errors"
    "service"
)

func Login(ctx context.Context) error {
	ormService := service.DI().ORMEngineFromContext(ctx)
	isActiveBundle := service.DI().FeatureFlag().IsActive(ormService, "bundle")
    if isActiveBundle{
       //your logic here
    }
	//your logic here

	return nil
}
```

There are 2 methods that will help you to run all cron jobs that are related to the feature flag

```go
for _, featureScript := range service.DI().FeatureFlag().GetScriptsSingleInstance(service.DI().OrmEngine()) {
			go b.RunScript(featureScript)
		}
```
and
```go
for _, featureScript := range service.DI().FeatureFlag().GetScriptsMultiInstance(service.DI().OrmEngine()) {
			go b.RunScript(featureScript)
		}
```
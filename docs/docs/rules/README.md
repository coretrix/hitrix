# Documentation
## Localhost tools
### force-alters
If you run your binary with argument `--force-alters` the program will check for DB and RediSearch alters and it will execute them(only in local mode).

`make web-api param=--force-alters`

## Domains
#### Naming convention for our backend domains are:

`[binary name].[env].[project].[domain]`

For prod we are skipping `[env]`
For example for our binary called `web-api` the domain will be `web-api.dev.lys.domain.com` or `web-api.demo.lys.domain.com`

#### Naming convention for our frontend domains are:

`[binary name without suffix api].[env].[project].[domain]`

For prod we are skipping [env]
For example for `web.dev.lys.domain.com` or `web.demo.lys.domain.com`

## Crons (scripts)
What is the differences between `single-instance-cron` and `multi-instance-cron`
- `single-instance-cron` is for crons that cannot scale. Imagine you read something from db every 10min and you update something. If you have more than one instance it's gonna conflict. That's why we gonna create only one pod for it
- `multi-instance-cron`  for crons that can scale. Imagine you read from queue. You can have as much consumers as you want. That`s why we gonna have more pod instances for it

## Naming conventions and rules
1. If you have variable that contains one or more entities you should add a suffix to the variable `Entity/Entities`
   For example productEntity or productEntities
2. Try to avoid `append()`
   Example:
```go
package main

func main()  {
	attributeEntities := make([]*struct{ID int}, 10)
	var someMap = make([]int, len(attributeEntities))
	for i, attributeEntity := range attributeEntities {
		someMap[i] = attributeEntity.ID //here we avoid map because we set the len
	}
}

```
3. If you implement communication with external API or something general that can be valid for every other project like Authentication for example, you can implement it as service in Hitrix.
4. When declaring a variable, use inferred variable declaration syntax:
```go
package main

func main() {
    //Acceptable
    _ = &OrderEntity{}

    //Not acceptable
    var orderEntity OrderEntity
    someFunc(&orderEntity)
}

type OrderEntity struct{
}
func someFunc(entity *OrderEntity)  {
 
}

```

5. For declaring slices, please use `make` to declare your variable:
```go
package main

func main() {
    //Acceptable
    _ = make([]*OrderEntity, 0)

    //Not acceptable
    var _ []*OrderEntity
}
type OrderEntity struct{
}
```

6. When instantiating graphql objects, use methods that are defined in `populate` package. If there is no method for your object, please declare one for your object.
7. Don't use `time.Now()` . Use `ioc.GetClockService().Now()` instead. If you need pointer, use `ioc.GetClockService().NowPointer()`.
8. All entity files should have `_entity.go` suffix, the entity itself should end with `Entity`.
9. In yaml config files we should set env vars only for values that going to be different in different environments(dev/demo/prod) If they are the same we should not use env var, but we can set the value into the yaml file
10. Custom redis indexes can be re-indexed using dev panel but dirty queues needs extra effort from developer side.
- You need to extend the slice into `DevPanelController->GetActionListAction` slice `dirty`
  This step will add new menu in dev panel dashboard.
- Be sure that you added GET url for in into your router
  For example `ginEngine.GET("/dev/mark-as-dirty-price-changed/", devPanel.GetMarkAsDirtyPriceChanged)`
- Your action should look like that
```go
package controller

type DevPanelController struct{
}
func (controller *DevPanelController) GetMarkAsDirtyPriceChanged(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	producer := producers.PriceChangedDirtyAllProducer{}

	err := producer.Produce(ormService)

	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)
		return
	}

	c.JSON(200, gin.H{})
}
```

And your processor should look like that:
```go
package model 
type PriceChangedDirtyAllProducer struct {
}

func (p *PriceChangedDirtyAllProducer) Produce(ormService *beeorm.Engine) error {
	variantEntity := entity.VariantEntity{}
	where := beeorm.NewWhere("1 ORDER BY ID ASC")
	pager := &beeorm.Pager{CurrentPage: 1, PageSize: 1000}
	hasMoreToIndex := true
	for hasMoreToIndex {
		ids := ormService.SearchIDs(where, pager, &variantEntity)

		if len(ids) == 0 {
			break
		}

		if len(ids) < pager.PageSize {
			hasMoreToIndex = false
		}

		ormService.MarkDirty(&entity.VariantEntity{}, redisstream.StreamOrmDirtyPriceChanged, ids...)

		pager.IncrementPage()
	}

	return nil
}
```
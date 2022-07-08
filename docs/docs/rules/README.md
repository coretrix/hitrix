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

#### Naming convention for our REST endpoints domains are:

`[noun]/[noun].../[action]/?params`

Always use `-` as separator. 
Endpoint name is not tied to package name.

`POST` - CREATE entity actions, SEARCH entity actions

`PATCH` - UPDATE entity actions

`GET` - GET entity actions

`DELETE` - DELETE entity actions

`PUT` - not used

Examples:

`GET /profile/payment-info/cards/get/` - gets cards information

`PATCH /profile/payment-info/cards/update/` - updates card information

`POST /profile/payment-info/cards/create/` - creates card

`DELETE /profile/payment-info/cards/delete/` - deletes card


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

## Teamwork rules
We have defined rules that backend and frontend developers should follow to keep good communication and deliver the feature like a one team
1. Read and discuss the epic with the business person
2. Backend and frontend developers together should go through design and define endpoints and request/response structure. Whenever they are ready they should post it as a comment into the backend ticket to be visible that both agreed on it
3. Start implementing the feature
4. Before completing the task they should do following things:
   1. When backend developer is done and all tests pass he should test every endpoint by himself using `swagger` on `dev` environment before complete his ticket
   2. When frontend developer is done he should deploy on `dev` and test the feature very well before complete his ticket 
5. When everything works on `dev` frontend developer is responsible to talk to backend developer and together to deploy on `demo` and go through the flow and verify if it works 
6. Mark the task as completed and inform the business person.
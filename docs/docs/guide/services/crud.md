# CRUD
This service it gives you ability to build gql query and apply different query parameters to the query that should be
used in listing pages

Register the service into your `main.go` file:
```go
registry.ServiceProviderCrud(),
```

Access the service:
```go
service.DI().Crud()
```

### Search vs Filter
Search is used on strings while filtering can be used on wide range of data types.  One important note to remember is that if your column is searchable it can't be filterable. 
 

### defining columns
First you need to define what columns you're going to have and which of them will be Searchable, Sortable or Filterable(user for enum values).
Using this configuration you also define the supported filters that can be applied.

The second step is in your controller(handler) to call few methods from that service that will build the right query for you based on the request.
Crud service supports mysql query builder and redis-search query builder.

Example of the way you can use it:

```go
//defining columns.
func columns() []crud.Column {
    return []crud.Column{
            {
                Key:            "ID",
                Type:           crud.NumberType,
                Label:          "ID",
                Searchable:     false,
                Sortable:       true,
                Visible:        true,
                Filterable:     true,
                FilterValidMap: nil,
            },
            {
                Key:            "Name",
                Type:           crud.StringType,
                Label:          "Name",
                Searchable:     true,
                Sortable:       false,
                Visible:        true,
                Filterable:     false,
                FilterValidMap: nil,
            }
        }
}
```

```go
//listing request using by gql
type ListRequest struct {
    Page     *int                   `json:"Page"`
    PageSize *int                   `json:"PageSize"`
    Filter   map[string]interface{} `json:"Filter"`
    Search   map[string]interface{} `json:"Search"`
    SearchOr map[string]interface{} `json:"SearchOR"`
    Sort     map[string]interface{} `json:"Sort"`
}
```
and at the end your method where you return the response:
```go
cols := columns()

crudService := service.DI().CrudService()

searchParams := crudService.ExtractListParams(cols, crud.ListRequestConvertorFromGQL(userListRequest))
query := crudService.GenerateListRedisSearchQuery(searchParams)

ormService := ioc.GetOrmEngineFromContext(ctx)
var userEntities []*entity.UserEntity

ormService.RedisSearch(&userEntities, query, beeorm.NewPager(searchParams.Page, searchParams.PageSize))

userEntityList := make([]*model.User, len(userEntities))

for i, userEntity := range userEntities {
    userEntityList[i] = populate.UserAdmin(userEntity)
}

return &model.UserList{
    Rows:    userEntityList,
    Total:   len(userEntityList),
    Columns: crud.ColumnConvertorToGQL(cols),
    }, nil
```

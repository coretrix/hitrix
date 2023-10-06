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

### Use Select fields with dependency

Imagine you have these fields with such options, on is the country and for the city filter to work user must select the country first and then select the proper city that exists in the list of that country otherwise this filter will be ignored.

Take london as an example, if we select `United States` as country and `London` as city, then city filter will be ignored, so make sure this is also handled in your front-end.

**Note: If you don't set `FilterDependencyField` or `DataMapStringStringKeyStringValue` this filter will be ignored!**

```go
func columns() []crud.Column {
    return []crud.Column{
    		{
			Key:                    "Country",
			FilterType:             hitrixCrud.SelectTypeStringString,
			Label:                  "Country",
			Searchable:             true,
			Sortable:               false,
			Visible:                true,
			TranslationDataEnabled: true,
			DataStringKeyStringValue: []*hitrixCrud.StringKeyStringValue{
				{Key: "bg", Label: "Bulgaria"},
				{Key: "us", Label: "United States"},
				{Key: "uk", Label: "United Kingdom"},
			},
		},
		{
			Key:                   "City",
			FilterType:            hitrixCrud.SelectTypeStringString,
			Label:                 "City",
			Searchable:            true,
			Sortable:              false,
			Visible:               true,
			FilterDependencyField: "Country",
			DataMapStringStringKeyStringValue: map[string][]*hitrixCrud.StringKeyStringValue{
				"bg": {
					{Key: "sofia", Label: "Sofia"},
					{Key: "plovdiv", Label: "Plovdiv"},
				},
				"us": {
					{Key: "new_york", Label: "New York"},
					{Key: "los_angeles", Label: "Los Angeles"},
				},
				"uk": {
					{Key: "london", Label: "London"},
				},
			},
		},
	}
}
```

### Use CRUD with our export service

You can mix our crud service with our exporter service to add a quick and painless exporting system to your project

In order to do that you can follow these steps

**1 - Making your currently list view function private and creating another function for passing to gin router**

for example here we renamed "List" to "list" and we created a new "ListRequest" function and we're passing it 
to gin router instead of "List"

```go

func Columns() []hitrixCrud.Column {
    return []hitrixCrud.Column{
        {
            Key:        "ID",
            Label:      "ID",
            Searchable: false,
            Sortable:   true,
            Visible:    true,
        },
        {
            Key:        "Name",
            FilterType: hitrixCrud.InputTypeString,
            Label:      "Name",
            Searchable: true,
            Sortable:   false,
            Visible:    true,
        },
    }
}

type ListRow struct {
    ID        uint64
    Name      string
}


func ListRequest(ctx context.Context, request *hitrixCrud.ListRequest) (*city.ResponseDTOList, error) {
	// You should pass this function to the gin router
	return list(service.DI().OrmEngineForContext(ctx), request)
}

func list(ormService *beeorm.Engine, request *hitrixCrud.ListRequest) (*city.ResponseDTOList, error) {
	
	// this is the function that handles the getting and making the payload fot the crud that is going to be used
	// for both list (endpoint) and exporting
	
	cols := Columns()
	crudService := service.DI().Crud()

	searchParams := crudService.ExtractListParams(cols, request)
	query := crudService.GenerateListRedisSearchQuery(searchParams)

	var cityEntities []*entity.CityEntity
	
	// your logic for getting entities from database with the query generated from crud

	rows := make([]ListRow, len(cityEntities))

	// you logic for creating rows for your crud

	return &ResponseDTOList{
		Rows:        rows,
		Total:       int(total),
		Columns:     cols,
		PageContext: getPageContext(ormService),
	}, nil
}
```

**2 - Creating a new handler for exporting the data obtained from "list" function**

```go
func ListExport(ormService *beeorm.Engine, request *hitrixCrud.ListRequest, _ uint64, _ map[string]string) (error, []string, [][]interface{}) {
	
	// This function handles the data for exporting
	
	exportColumns := make([]string, 0) // excel or csv Columns for passing to our exporter service
	allExportData := make([][]interface{}, 0) // data for passing to our exporter service

	pager := beeorm.NewPager(1, 1000)

	for {
		
		// loop for getting all the data out of the list function, bypassing the pagination
		
		request.Page = &pager.CurrentPage
		request.PageSize = &pager.PageSize

		res, err := list(ormService, request)

		if err != nil {
			return err, nil, nil
		}
		
		// GetExporterDataCrud converts any given rows which in example is []city.ListRow to
		// exporter service columns and data as types of []string, [][]interface, exactly the data you need for passing to our
		// exporting system
		
		columns, exportData := hitrixCrud.GetExporterDataCrud(Columns(), res.Rows)

		for _, exportDataRow := range exportData {
			allExportData = append(allExportData, exportDataRow)
		}

		exportColumns = columns

		if len(res.Rows) < pager.PageSize {
			break
		}

		pager.IncrementPage()
	}

	return nil, exportColumns, allExportData
}
```

**3 - Creating the config and passing it to the service**

we pass the `ListExport` function that we created above to the config

```go
const (
    CityExportID                = "cities"
)

var RegisteredExportConfigs = []crud.ExportConfig{
	{
		Handler:          cityView.ListExport,
		ID:               CityExportID,
		AllowedExtraArgs: nil,
		Resource:         "city",
		Permissions:      []string{"read"},
	},
}

RegisterDIGlobalService(
    ...
    hitrixRegistry.ServiceProviderCrud(RegisteredExportConfigs),
    ...
)

```

**4 - Getting export ready data from crud**

Imagine we have these two cities in our database

```json
[
  {
    "ID": 1,
    "Name": "Sofia"
  },
  {
    "ID": 2,
    "Name": "New York"
  }
]
```

We use the crud method to get config/handler and export our data

```go
handler, exists := crudService.GetExportHandler(CityExportID) // getting the handler

err, columns, data := handler(
    ormService,
    &hitrixCrud.ListRequest{},
    1, // user requesting the export
    nil, 
)

print(columns)    // would print []string{"ID", "Name"}
print(data) // would print [][]interface{}{[]interface{"1", "Sofia"}, []interface{"2", "New York"}}

```

and you can easily pass this data to our exporter service, for example:

```go
excelBytes, err := exporterService.XLSXExportToByte("cities", columns, data)
```

and if you like to check the permissions you can get the whole config like this

```go
config, exists := crudService.GetExportConfig(CityExportID) // getting the config
```
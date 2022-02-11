# Dataloaders
What are GraphQL Dataloaders? Please take a look here https://gqlgen.com/reference/dataloaders/

## How to create a dataloader?
1. Create a folder "dataloaders" inside graphql folder
2. Create file dataloaders.go

```go
package dataloaders

//go:generate dataloaden LoaderName uint64 []*/path/entity.SomeModel/Entity
```

3. Generate dataloaders
```bash
go generate ./api/<binary>/graphql/dataloaders/...
```
4. Implement loaders

go generate will create gen.go file that contains loaders skeleton.
```go
func NewLoaders(ctx context.Context) *Loaders {
    return &Loaders{
        VariantsByProductID: LoaderNameLoader{
        maxBatch: 1000,
        wait:     5 * time.Millisecond,
        fetch: func(ids []uint64) ([][]*entity.SomeModel, []error) {
                //Some code here
				return someModelSlice, errors
			},
		}
	}
}
```

5. Create middleware and attach Dataloaders to context
 - Create file dataloader.go
```go
func DataLoaders(ginEngine *gin.Engine) {
	ginEngine.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, key, NewLoaders(ctx))
		c.Request = c.Request.WithContext(ctx)
	})
}
```
 - Add DataLoaders middleware to Gin

6. Add Retriever to resolvers

```go
// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
    DataLoaders dataloaders.Retriever
}
```

7. git add & git commit

## How to use a dataloader?

1. Accessing Dataloaders - dataloaders can be accessed using Retreive function in resolvers.

```go
func (r *someResolver) SomeMethod(ctx context.Context, obj *model.SomeModel) ([]*model.SomeOtherModel, error) {
	if !hitrix.Validate(ctx, nil) {
		return nil, nil
	}

	someOtherModel, err := r.DataLoaders.Retrieve(ctx).VariantsByProductID.Load(obj.ID)

	if err != nil {
		return nil, err
	}

	return someOtherModel, nil
}
```
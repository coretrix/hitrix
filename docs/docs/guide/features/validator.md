# Validator
We support 2 types of validators. One of them is related to graphql, the other one is related to rest.

#### Graphql validator
There are 2 steps that needs to be executed if you want to use this kind of validator

1. Add `directive @validate(rules: String!) on INPUT_FIELD_DEFINITION` into your `schema.graphqls` file

2. Call `ValidateDirective` into your main.go file
```go
config := generated.Config{Resolvers: &graph.Resolver{}, Directives: generated.DirectiveRoot{Validate: hitrix.ValidateDirective()} }

s.RunServer(4001, generated.NewExecutableSchema(config), func(ginEngine *gin.Engine) {
    commonMiddleware.Cors(ginEngine)
    middleware.Router(ginEngine)
})
```

After that you can define the validation rules in that way:
```graphql
input ApplePurchaseRequest {
  ForceEmail: Boolean!
  Name: String
  Email: String @validate(rules: "email") #for rules param you can use everything supported by https://github.com/go-playground/validator validate.Var(value, rules)
  AppleReceipt: String!
}
```

To handle the errors you need to call function `hitrix.Validate(ctx, nil)` in your resolver
```go
func (r *mutationResolver) RegisterTransactions(ctx context.Context, applePurchaseRequest model.ApplePurchaseRequest) (*model.RegisterTransactionsResponse, error) {
    if !hitrix.Validate(ctx, nil) {
        return nil, nil
    }
    // your logic here...
}
```

The function `hitrix.Validate(ctx, nil)` as second param accept callback where you can define your custom validation related to business logic


#### REST validator

You should define tags for every field
```go
type RequestDTOMerchantSave struct {
	StoreID         string `conform:"trim" binding:"required,min=1,max=30"`
	StoreBio        string `conform:"trim" binding:"omitempty,min=5,max=1000"`
	AvatarFileID    *uint64
	ContactPhone    *ContactPhone    `binding:"omitempty"`
	ContactWhatsapp *ContactWhatsapp `binding:"omitempty"`
	ContactWeb      string           `binding:"omitempty,url"`
	ContactTelegram *ContactTelegram `binding:"omitempty"`
	ContactEmail    string           `conform:"trim" binding:"omitempty,email"`
}
```

Using `binding` you can define all rules needed for the particular validation
Using `conform` you can trim the value before validation to be applied


##### Validation notes
	RepatriationAfterTyreBlockInMinutes int               `binding:"numeric,gte=0"`
If you want to support 0 value you should not put `required` tag for the fields because the validator thinks that zero is not a value

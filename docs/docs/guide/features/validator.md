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

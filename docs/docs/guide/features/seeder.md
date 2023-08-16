# Database Seeding
Hitrix supports seeds that can be used to populate database or apply changes on it.

To Seed your database ([wiki](https://en.wikipedia.org/wiki/Database_seeding)),
you can use the seeding feature that is implemented as script (`DBSeedScript`).
This `DBSeedScript` needs to be provided a `Seeds map[string][]Seed` field, Where the string key is the identifier of the project that is used. This measn you can have different project in one code base
The Script can be implemented in your app by making a type that satisfies the `Seed` interface:
```go
type Seed interface {
    Execute(*datalayer.DataLayer)
    Environments() []string
    Name() string
}
```
Example:

`users_seed.go`
```go
type UsersSeed struct{}

func (seed *UsersSeed) Name() string {
    return "UsersSeed"
}

func (seed *UsersSeed) Environments() []string {
    return []string{app.ModeTest, app.ModeLocal, app.ModeDev, app.ModeDemo, app.ModeProd}
}

func (seed *UsersSeed) Execute(ormService *datalayer.DataLayer) {
	// TODO insert a new user entity to the db
}
```
And after your server definition (such as the one in [server.go](example/server.go)),
you could use it to run the seeds just like you would run any script (explained earlier in the [scripts seciton](#running-scripts)) as follows:
```go
// Server Definition
s,_ := hitrix.New()

...

s.RunBackgroundProcess(func(b *hitrix.BackgroundProcessor) {
  // seed database
  go b.RunScript(&hitrixScripts.DBSeedScript{
    SeedsPerProject: map[string][]hitrixScripts.Seed{
        "project_name": {
            &script.UserProfileAttributesSeed{},
        },
    })

  // ... other scripts
})
```

This seed will only run if no seeds has ever been ran with that name (in db table seeder, no row with `'name'='UsersSeed'` )



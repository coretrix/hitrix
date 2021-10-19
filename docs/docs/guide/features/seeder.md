# Database Seeding
Hitrix supports multi-versioned multiple seeds.

To Seed your database ([wiki](https://en.wikipedia.org/wiki/Database_seeding)),
you can use the seeding feature that is implemented as script (`DBSeedScript`).
This `DBSeedScript` needs to be provided a `Seeds map[string]Seed` field, Where the string key is the identifier of the seed that is used to detect seed versioning.
The Script can be implementd in your app by making a type that satisfies the `Seed` interface:
```go
type Seed interface {
	Execute(*beeorm.Engine)
	Version() int
}
```
Example:

`users_seed.go`
```go
type UsersSeed struct{}

func (seed *UsersSeed) Version() int {
	return 1
}

func (seed *UsersSeed) Execute(ormService *beeorm.Engine) {
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
    Seeds: map[string]hitrixScripts.Seed{
      "users-seed": UsersSeed{},
    },
  })

  // ... other scripts
})
```

This seed will only run in the following cases:
* No seeds has ever been ran (in db table settings, no row with `'key'='seeds'` )
* This seed has never been ran
* The seed IS ran before, but with an older version


# Config
This service provides you access to your config file. We support only YAML file.

Register the service into your `main.go` file:
```go
 registry.ServiceProviderConfigDirectory("../config")
```
you should provide the folder where are your config files.

The folder structure should look like that:
```
config
 - app-name
    - config.yaml
 - hitrix.yaml #optional config where you can define some settings related to built-in services like slack service
```

Access the service:
```go
service.DI().Config()
```

## Environment variables in config file
Its good practice to keep your secrets like database credentials and so on out of the repository.
Our advice is to keep them like environment variables and call them into config.yaml file
For example your config can looks like this:
```yaml
orm:
  default:
    mysql: ENV[DEFAULT_MYSQL]
    redis: ENV[DEFAULT_REDIS]
    locker: default
    local_cache: 1000
```
where `DEFAULT_MYSQL` and `DEFAULT_REDIS` are env variables and our framework will automatically replace `ENV[DEFAULT_MYSQL]` and `ENV[DEFAULT_REDIS]` with the right values

If you want to define array of values you should split them by `;` and they will be presented into the yaml file in that way:
```yaml
cors:
    - test1
    - test2
```

If you want to enable the debug for orm you can add this tag `orm_debug: true` on the main level of your config

Also we check if there is .env.XXX file in main config folder where XXX is the value of the APP_MODE.
If there is for example .env.local we are reading those env variables and merge them with config.yaml how we presented above

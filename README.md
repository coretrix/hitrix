# Hitrix

[![checks & tests](https://github.com/coretrix/hitrix/actions/workflows/main.yml/badge.svg)](https://github.com/coretrix/hitrix/actions)
[![codecov](https://codecov.io/gh/coretrix/hitrix/branch/main/graph/badge.svg)](https://codecov.io/gh/coretrix/hitrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/coretrix/hitrix)](https://goreportcard.com/report/github.com/coretrix/hitrix)
[![GPL3 license](https://img.shields.io/badge/license-GPL3-brightgreen.svg)](https://opensource.org/licenses/GPL-3.0)

Hitrix is a framework written in Go based on Gin. 
 It is for web and console applications and provides set of reusable services and features.

### Installation
Please follow our [Official documentation](https://coretrix.github.io/hitrix/)


### The key services and features of Hitrix are:
- Config service
- Dependency injection container 
- ORM service to access database
- Redis Cache layer on top of mysql
- Redis Search layer on top of mysql 
- Redis Stream for queues
- OSS service
- Error logger service
- Mail provider service
- Authentication service
- Feature flag service
- Clockwork service
- Background scripts
- Integration test engine
- Seeder
- Validator
- etc...

------------

#### Hitrix exposes API for special [DEV-PANEL](https://github.com/coretrix/dev-frontend) web frontend where you can manage all your queues, error log, performance, mysql alters and so on. You can fully manage your application using our dev panel

------------
### For Contributors
Every new feature or change already existing feature should be described into the documentation.
To change the documentation you need to execute next steps:

1. Go to `/docs/docs` folder
2. Change the documentation using `.md` files
3. If you want to change sidebar or navbar you can do it in `/docs/docs/.vuepress/config.js`
4. Run it on localhost from, folder `/docs` using `yarn docs:dev`
5. Push your changes

If you create a new service and write documentation for it please follow next template:
1. service definition - what this service can be used for
2. register service - how the service should be registered
3. access service - how the service can be accessed
4. technical documentation - details about the service, examples and different use cases

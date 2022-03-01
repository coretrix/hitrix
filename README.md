[![test status](https://github.com/go-gorm/gorm/workflows/tests/badge.svg?branch=master "test status")](https://github.com/coretrix/hitrix/actions)
[![codecov](https://codecov.io/gh/coretrix/hitrix/branch/main/graph/badge.svg)](https://codecov.io/gh/coretrix/hitrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/coretrix/hitrix)](https://goreportcard.com/report/github.com/coretrix/hitrix)
[![GPL3 license](https://img.shields.io/badge/license-GPL3-brightgreen.svg)](https://opensource.org/licenses/GPL-3.0)

# [Official documentation](https://coretrix.github.io/hitrix/)

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
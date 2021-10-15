# Mail Mandrill service

This service is used for sending transactional emails using Mandrill

Register the service into your `main.go` file:
```go
registry.MailMandrill()
```
and you should register the entity `MailTrackerEntity` into the ORM
Also, you should put your credentials and other configs in your config file

```yml
mandrill:
  api_key: ...
  default_from_email: test@coretrix.tv
  from_namme: coretrix.com
```

Access the service:
```go
service.DI().Mandrill()
```

Some of the functions this service provide are:
```go
	SendTemplate(ormService *beeorm.Engine, message *TemplateMessage) error
	SendTemplateAsync(ormService *beeorm.Engine, message *TemplateMessage) error
	SendTemplateWithAttachments(ormService *beeorm.Engine, message *TemplateAttachmentMessage) error
	SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *TemplateAttachmentMessage) error
```

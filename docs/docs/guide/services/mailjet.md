# Mail Mailjet service

This service is used for sending transactional emails using Mailjet

Register the service into your `main.go` file:
```go
registry.ServiceProviderMail(mail.NewMailjet)
```
and you should register the entity `MailTrackerEntity` into the ORM
Also, you should put your credentials and other configs in your config file

```yml
mailjet:
  api_key_public: ...
  api_key_private: ...
  default_from_email: test@coretrix.tv
  from_name: coretrix.com
  sandbox_mode: false
```

Access the service:
```go
service.DI().Mailjet()
```

Some of the functions this service provide are:
```go
	SendTemplate(ormService *beeorm.Engine, message *TemplateMessage) error
	SendTemplateAsync(ormService *beeorm.Engine, message *TemplateMessage) error
	SendTemplateWithAttachments(ormService *beeorm.Engine, message *TemplateAttachmentMessage) error
	SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *TemplateAttachmentMessage) error
```

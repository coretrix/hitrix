# Mail Mailjet service

This service is used for sending transactional emails using providers like Mailjet or Mandrill

Register the service into your `main.go` file:
```go
registry.ServiceProviderMail(mail.NewMailjet)
```
and you should register the entity `MailTrackerEntity` into the ORM
Also, you should put your credentials and other configs in your config file

```yml
mail:
  mailjet:
    api_key_public: ...
    api_key_private: ...
    default_from_email: test@coretrix.tv
    from_name: coretrix.com
    sandbox_mode: false
  mandrill:
    api_key: ...
    default_from_email: test@coretrix.tv
    from_name: coretrix.com
```

If you set `sandbox_mode=true` we won't send real email to the customer

Access the service:
```go
service.DI().Mail()
```

Some functions this service provide are:
```go
    GetTemplateKeyFromConfig(templateName string) (string, error)
    SendTemplate(ormService *beeorm.Engine, message *Message) error
    SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error
    GetTemplateHTMLCode(templateName string) (string, error)
```

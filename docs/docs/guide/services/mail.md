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
    whitelisted_emails:
      - "@coretrix.com"
  mandrill:
    api_key: ...
    default_from_email: test@coretrix.tv
    from_name: coretrix.com
```

If you set `sandbox_mode=true` we won't send real email to the customer.

BUT we allow sandbox_mode to be disabled only for specific emails or domains.
That's why we created config param called `whitelisted_emails` All emails or domains defined in this config param will be sent to receiver even if sandbox is enabled.

The idea is to be able to test the platform and in the same time to be sure that real emails are not sent to customers.

Access the service:
```go
service.DI().Mail()
```

Some functions this service provide are:
```go
    GetTemplateKeyFromConfig(templateName string) (string, error)
    SendTemplate(ormService *datalayer.DataLayer, message *Message) error
    SendTemplateWithAttachments(ormService *datalayer.DataLayer, message *MessageAttachment) error
    GetTemplateHTMLCode(templateName string) (string, error)
```

# Dynamic link service
This service is used for generating dynamic links, at this moment only Firebase is supported

Register the service into your `main.go` file:
```go
registry.ServiceProviderDynamicLink(),
```

Access the service:
```go
service.DI().DynamicLink()
```

Config sample:

```yml
 api_key: string # required
 dynamic_link_info: # required
   domain_uri_prefix: string # required
   link: string # required
   android_info: # optional
     package_name: string # optional
     fallback_link: string # optional
     min_package_version_code: string # optional
   ios_info: # optional
     bundle_id: string # optional
     fallback_link: string # optional
     custom_scheme: string # optional
     ipad_fallback_link: string # optional
     ipad_bundle_id: string # optional
     app_store_id: string # optional
   navigation_info: # optional
     enable_forced_redirect: boolean # required
   analytics_info: # optional
     google_play_analytics: # optional
       utm_source: string # optional
       utm_medium: string # optional
       utm_campaign: string # optional
       utm_term: string # optional
       utm_content: string # optional
       gcl_id: string # optional
     itunes_connect_analytics: # optional
       at: string # optional
       ct: string # optional
       mt: string # optional
       pt: string # optional
   social_meta_tag_info: # optional
     social_title: string # optional
     social_description: string # optional
     social_image_link: string # optional
 suffix: # optional
   option: string # required, values: "SHORT" or "UNGUESSABLE"
```

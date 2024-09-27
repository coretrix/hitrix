# Validator

#### REST validator

You should define tags for every field
```go
type RequestDTOMerchantSave struct {
	StoreID         string `conform:"trim" binding:"required,min=1,max=30"`
	StoreBio        string `conform:"trim" binding:"omitempty,min=5,max=1000"`
	AvatarFileID    *uint64
	ContactPhone    *ContactPhone    `binding:"omitempty"`
	ContactWhatsapp *ContactWhatsapp `binding:"omitempty"`
	ContactWeb      string           `binding:"omitempty,url"`
	ContactTelegram *ContactTelegram `binding:"omitempty"`
	ContactEmail    string           `conform:"trim" binding:"omitempty,email"`
}
```

Using `binding` you can define all rules needed for the particular validation
Using `conform` you can trim the value before validation to be applied


##### Validation notes
	RepatriationAfterTyreBlockInMinutes int               `binding:"numeric,gte=0"`
If you want to support 0 value you should not put `required` tag for the fields because the validator thinks that zero is not a value

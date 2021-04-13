package helper

import (
	"fmt"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vEn "github.com/go-playground/validator/v10/translations/en"
	"github.com/pariz/gountries"
)

var (
	once      sync.Once
	singleton *Validator
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewValidator() *Validator {
	once.Do(func() {
		validatorInstance := validator.New()
		for ruleName, validatorFunction := range customValidations {
			err := validatorInstance.RegisterValidation(ruleName, validatorFunction)
			if err != nil {
				panic(err)
			}
		}
		english := en.New()
		uni := ut.New(english, english)
		translator, _ := uni.GetTranslator("en")
		_ = vEn.RegisterDefaultTranslations(validatorInstance, translator)
		singleton = &Validator{validator: validatorInstance, translator: translator}
	})
	return singleton
}

func (t *Validator) Validate(field interface{}, rules string) []error {
	err := t.validator.Var(field, rules)
	return t.translateError(err)
}

func (t *Validator) translateError(err error) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(t.translator))
		errs = append(errs, translatedErr)
	}
	return errs
}

// custom validators
var customValidations = map[string]func(validator.FieldLevel) bool{
	"country_code": validateCountryCodeAlpha2,
}

func validateCountryCodeAlpha2(fl validator.FieldLevel) bool {
	query := gountries.New()
	_, err := query.FindCountryByAlpha(fl.Field().String())
	return err == nil
}

package helper

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/coretrix/hitrix/pkg/errors"

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

func (t *Validator) ValidateStruct(s interface{}) error {
	err := t.validator.Struct(s)

	if err != nil {
		var fieldErrors errors.FieldErrors = make(map[string]string)
		validatorErrs := err.(validator.ValidationErrors)
		for _, e := range validatorErrs {
			translatedErr := e.Translate(t.translator)
			fieldErrors[e.Field()] = translatedErr
		}
		return fieldErrors
	}
	return nil
}

func (t *Validator) Engine() interface{} {
	return t.validator
}

func (t *Validator) Validate(field interface{}, rules string) []error {
	err := t.validator.Var(field, rules)
	return t.translateError(err)
}

func NewValidator() *Validator {
	once.Do(func() {
		validatorInstance := validator.New()
		validatorInstance.SetTagName("binding")
		english := en.New()
		uni := ut.New(english, english)
		translator, _ := uni.GetTranslator("en")

		for ruleName, customValidation := range customValidations {
			err := validatorInstance.RegisterValidation(ruleName, customValidation.ValidatorFunction)
			if err != nil {
				panic(err)
			}
			err = validatorInstance.RegisterTranslation(ruleName, translator, func(ut ut.Translator) error {
				return ut.Add(ruleName, customValidation.TranslationMessage, false)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(ruleName, fe.Field())
				return t
			})
			if err != nil {
				panic(err)
			}
		}
		_ = vEn.RegisterDefaultTranslations(validatorInstance, translator)

		singleton = &Validator{validator: validatorInstance, translator: translator}
	})
	return singleton
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

type CustomValidation struct {
	ValidatorFunction  func(validator.FieldLevel) bool
	TranslationMessage string
}

// custom validators
var customValidations = map[string]CustomValidation{
	"country_code_custom": {
		ValidatorFunction:  validateCountryCodeAlpha2,
		TranslationMessage: "not a valid Country Code",
	},
	"password_strength": {
		ValidatorFunction:  validatePasswordStrength(8),
		TranslationMessage: "Not strong enough. Should be more than 8 character, contain at least 1 lowercase, 1 uppercase, 1 number, and 1 special character.",
	},
}

func validateCountryCodeAlpha2(fl validator.FieldLevel) bool {
	query := gountries.New()
	_, err := query.FindCountryByAlpha(fl.Field().String())
	return err == nil
}

func validatePasswordStrength(minLength int) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		pass := fl.Field().String()

		if len(pass) < minLength {
			return false
		}

		ok, _ := regexp.MatchString(`[a-z]+`, pass)
		if !ok {
			return false
		}

		ok, _ = regexp.MatchString(`[A-Z]+`, pass)
		if !ok {
			return false
		}

		ok, _ = regexp.MatchString(`[0-9]+`, pass)
		if !ok {
			return false
		}

		// ref: https://owasp.org/www-community/password-special-characters
		specialChars := " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
		ok, _ = regexp.MatchString("["+specialChars+"]+", pass)

		return ok
	}
}

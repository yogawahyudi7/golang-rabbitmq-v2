package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang-rabbitmq-v2/pkg/logger"
)

type Validator struct {
	validator *validator.Validate
	logger    *logger.Logger
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

func New(log *logger.Logger) *Validator {
	v := validator.New()
	
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validator: v,
		logger:    log,
	}
}

func (v *Validator) Validate(data interface{}) error {
	if err := v.validator.Struct(data); err != nil {
		var validationErrors ValidationErrors
		
		for _, err := range err.(validator.ValidationErrors) {
			validationError := ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: v.getErrorMessage(err),
			}
			validationErrors = append(validationErrors, validationError)
		}
		
		v.logger.WithContext("validator").WithFields(map[string]interface{}{
			"errors": validationErrors,
			"struct": reflect.TypeOf(data).String(),
		}).Error("Validation failed")
		
		return validationErrors
	}
	
	return nil
}

func (v *Validator) getErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()
	value := fmt.Sprintf("%v", err.Value())

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("%s must be no more than %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be no more than %s", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID v4", field)
	case "json":
		return fmt.Sprintf("%s must be valid JSON", field)
	case "base64":
		return fmt.Sprintf("%s must be valid base64", field)
	case "containsany":
		return fmt.Sprintf("%s must contain at least one of the following characters: %s", field, param)
	case "excludes":
		return fmt.Sprintf("%s cannot contain '%s'", field, param)
	case "excludesall":
		return fmt.Sprintf("%s cannot contain any of the following characters: %s", field, param)
	case "excludesrune":
		return fmt.Sprintf("%s cannot contain the character '%s'", field, param)
	default:
		return fmt.Sprintf("%s failed validation for '%s' with value '%s'", field, tag, value)
	}
}
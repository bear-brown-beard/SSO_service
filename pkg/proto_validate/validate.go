package proto_validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	uni *ut.UniversalTranslator
)

type Validate struct {
	validate *validator.Validate
	lang     string
}

func NewValidate() *Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validate{validate: v, lang: "en"}
}

func (v *Validate) Check(req interface{}) error {
	var trans ut.Translator
	if v.lang == "ru" {
		ru := ru.New()
		uni = ut.New(ru, ru)

		trans, _ = uni.GetTranslator("ru")
		ru_translations.RegisterDefaultTranslations(v.validate, trans)
	} else {
		en := en.New()
		uni = ut.New(en, en)

		trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v.validate, trans)
	}

	if err := v.validate.Struct(req); err != nil {
		violations := make([]*errdetails.BadRequest_FieldViolation, 0)

		for _, err := range err.(validator.ValidationErrors) {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       err.Field(),
				Reason:      err.Tag(),
				Description: err.Translate(trans),
			})
		}

		st := status.New(codes.InvalidArgument, "bad request")
		badRequest := &errdetails.BadRequest{
			FieldViolations: violations,
		}

		statusWithDetails, err := st.WithDetails(badRequest)
		if err != nil {
			return status.Error(codes.Internal, "internal server error")
		}

		return statusWithDetails.Err()
	}
	return nil
}
func (v *Validate) CheckByLang(req interface{}, lang string) error {
	var trans ut.Translator
	if lang == "ru" {
		ru := ru.New()
		uni = ut.New(ru, ru)

		trans, _ = uni.GetTranslator("ru")
		ru_translations.RegisterDefaultTranslations(v.validate, trans)
	} else {
		en := en.New()
		uni = ut.New(en, en)

		trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v.validate, trans)
	}

	if err := v.validate.Struct(req); err != nil {
		violations := make([]*errdetails.BadRequest_FieldViolation, 0)

		for _, err := range err.(validator.ValidationErrors) {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       err.Field(),
				Reason:      err.Tag(),
				Description: err.Translate(trans),
			})
		}

		st := status.New(codes.InvalidArgument, "bad request")
		badRequest := &errdetails.BadRequest{
			FieldViolations: violations,
		}

		statusWithDetails, err := st.WithDetails(badRequest)
		if err != nil {
			return status.Error(codes.Internal, "internal server error")
		}

		return statusWithDetails.Err()
	}
	return nil
}

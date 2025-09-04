package controller

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/yiran15/api-server/base/apitypes"
)

var (
	trans    ut.Translator
	validate *Validator
)

type CheckReqInterface interface {
	CheckReq(ctx context.Context, value any) (errMsg string, err error)
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() error {
	v := validator.New()
	zhTrans := zh.New()
	uni := ut.New(zhTrans, zhTrans)
	trans, _ = uni.GetTranslator("zh")
	if err := zh_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return fmt.Errorf("register default translations failed: %w", err)
	}
	if err := registerValidator(v, trans); err != nil {
		return fmt.Errorf("register validator failed: %w", err)
	}
	validate = &Validator{
		validate: v,
	}
	return nil
}

func (receiver *Validator) CheckReq(ctx context.Context, value any) (errMsg string, err error) {
	err = receiver.validate.Struct(value)
	if err == nil {
		return "", nil
	}

	var valErrors validator.ValidationErrors
	if !errors.As(err, &valErrors) {
		return "", errors.New("validate check exception")
	}

	// 使用翻译器
	msgArr := make([]string, 0, len(valErrors))
	for _, e := range valErrors {
		msg := e.Translate(trans)
		msgArr = append(msgArr, msg)
	}

	return strings.Join(msgArr, "; "), nil
}

func registerValidator(v *validator.Validate, trans ut.Translator) error {
	if err := v.RegisterValidation("user_list", userListValidator); err != nil {
		return fmt.Errorf("register user_list validator failed: %w", err)
	}
	err := v.RegisterTranslation("user_list", trans,
		func(ut ut.Translator) error {
			return ut.Add("user_list", "email、mobile 和 name 中最多只能有一个字段非空", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("user_list", fe.Field())
			return t
		},
	)
	if err != nil {
		return fmt.Errorf("register user_list translation failed: %w", err)
	}

	if err := v.RegisterValidation("mobile", mobileValidator); err != nil {
		return fmt.Errorf("register mobile validator failed: %w", err)
	}
	if err := v.RegisterTranslation("mobile", trans,
		func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 必须是有效的中国大陆手机号码", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		},
	); err != nil {
		return fmt.Errorf("register mobile translation failed: %w", err)
	}

	return nil
}

func userListValidator(fl validator.FieldLevel) bool {
	// 获取整个结构体
	user, ok := fl.Parent().Interface().(apitypes.UserListRequest)
	if !ok {
		return false // 如果类型断言失败，返回验证失败
	}

	var count int
	if user.Email != "" {
		count++
	}
	if user.Mobile != "" {
		count++
	}
	if user.Name != "" {
		count++
	}
	return count <= 1
}

var mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func mobileValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return mobileRegex.MatchString(field)
}

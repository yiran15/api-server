package controller

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/yiran15/api-server/base/apitypes"
)

var (
	trans ut.Translator
)

// NewValidator 初始化自定义验证器和翻译器
func NewValidator() error {
	zhTrans := zh.New()
	uni := ut.New(zhTrans, zhTrans)
	trans, _ = uni.GetTranslator("zh")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := zh_translations.RegisterDefaultTranslations(v, trans); err != nil {
			return fmt.Errorf("register default translations failed: %w", err)
		}
		if err := registerValidator(v); err != nil {
			return fmt.Errorf("register validator failed: %w", err)
		}
	}
	return nil
}

// translateErrors 将验证错误翻译成更友好的格式
func translateErrors(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var errMsg []string
	for _, v := range errs.Translate(trans) {
		errMsg = append(errMsg, v)
	}
	return strings.Join(errMsg, "; ")
}

// registerValidator 注册自定义验证器
func registerValidator(v *validator.Validate) error {
	if err := registerUserList(v); err != nil {
		return err
	}
	if err := registerMobile(v); err != nil {
		return err
	}
	return nil
}

var userListValidator validator.Func = func(fl validator.FieldLevel) bool {
	user, ok := fl.Parent().Interface().(apitypes.UserListRequest)
	if !ok {
		return false
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

func registerUserList(v *validator.Validate) error {
	if err := v.RegisterValidation("user_list", userListValidator); err != nil {
		return fmt.Errorf("register user_list validator failed: %w", err)
	}

	if err := v.RegisterTranslation("user_list", trans,
		func(ut ut.Translator) error {
			return ut.Add("user_list", "email、mobile 和 name 中最多只能有一个字段非空", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("user_list", fe.Field())
			return t
		},
	); err != nil {
		return fmt.Errorf("register user_list translation failed: %w", err)
	}
	return nil
}

var mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

var mobileValidator = func(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return mobileRegex.MatchString(field)
}

func registerMobile(v *validator.Validate) error {
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

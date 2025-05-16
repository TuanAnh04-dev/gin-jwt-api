package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var fieldNameMap = map[string]string{
	"Name":     "Tên",
	"Price":    "Giá",
	"Quantity": "Số lượng",
	"Email":    "Email",
	"Password": "Mật khẩu",
	// Thêm các field khác tùy theo project của bạn
}

// ValidateStruct nhận một struct và trả về map lỗi tiếng Việt
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errors[field] = translateError(e)
	}

	return errors
}

// translateError dịch lỗi từ tag sang tiếng Việt
func translateError(fe validator.FieldError) string {
	field := fieldNameMap[fe.Field()]
	if field == "" {
		field = fe.Field()
	}
	switch fe.Tag() {
	case "required":
		return field + " là bắt buộc"
	case "min":
		return field + " phải có tối thiểu " + fe.Param() + " ký tự"
	case "max":
		return field + " không được vượt quá " + fe.Param() + " ký tự"
	case "email":
		return field + " không đúng định dạng email"
	case "gte":
		return field + " phải lớn hơn hoặc bằng " + fe.Param()
	case "gt":
		return field + " phải lớn hơn " + fe.Param()
	default:
		return field + " không hợp lệ"
	}
}

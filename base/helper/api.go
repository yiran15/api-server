package helper

import (
	"errors"
	"fmt"

	"github.com/yiran15/api-server/model"
)

func ValidateRoleApis(reqApis []int64, total int64, apis []*model.Api) error {
	if len(reqApis) == 0 {
		return errors.New("req api ids is empty")
	}
	if len(apis) == 0 {
		return fmt.Errorf("apis not found: %v", reqApis)
	}
	if len(reqApis) == int(total) {
		return nil
	}

	foundApis := make([]int64, 0, total)
	for _, v := range apis {
		foundApis = append(foundApis, v.ID)
	}

	notFoundApis := make([]int64, 0, total)
	for _, api := range reqApis {
		if !InArray(foundApis, api) {
			notFoundApis = append(notFoundApis, api)
		}
	}
	return fmt.Errorf("apis not found: %v", notFoundApis)
}

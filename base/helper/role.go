package helper

import (
	"fmt"

	"github.com/yiran15/api-server/model"
)

func ValidateRoleIds(reqRoleIds []int64, roles []*model.Role, total int64) error {
	// if len(reqRoleIds) == 0 {
	// 	return errors.New("role ids is empty")
	// }
	// if len(roles) == 0 {
	// 	return errors.New("roles is empty")
	// }
	if len(reqRoleIds) == int(total) {
		return nil
	}

	dbRoleIds := make([]int64, 0, total)
	for _, role := range roles {
		dbRoleIds = append(dbRoleIds, role.ID)
	}

	notFoundRoleIds := make([]int64, 0, total)
	for _, roleId := range reqRoleIds {
		if !InArray(dbRoleIds, roleId) {
			notFoundRoleIds = append(notFoundRoleIds, roleId)
		}
	}
	return fmt.Errorf("roles not found: %v", notFoundRoleIds)
}

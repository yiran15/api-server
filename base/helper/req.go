package helper

// RemoveDuplicates 是一个泛型去重函数，接受类型为 T 的切片，其中 T 需满足 comparable 约束。
// 返回去重后的切片，保持原顺序。
func RemoveDuplicates[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}
	// 使用 map 记录已出现的元素
	seen := make(map[T]struct{})
	// 结果切片，保持原顺序
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		// 如果元素未出现过，添加到结果并标记为已出现
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

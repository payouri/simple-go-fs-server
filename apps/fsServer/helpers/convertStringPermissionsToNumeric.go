package helpers

import "fmt"

/*
convertPermissions converts a permission string like "rwxr-xr-x" to a numeric format like 755.
*/
func ConvertStringPermissionsToNumeric(permStr string) string {
	permMap := map[rune]string{
		'r': "4",
		'w': "2",
		'x': "1",
		'-': "0",
	}

	var numericPerm string
	for _, c := range permStr[1:] { // Skip the first character (e.g., 'd' in drwxr-xr-x)
		numericPerm += permMap[c]
	}

	// Split into three parts and sum each part
	user := numericPerm[0:3]
	group := numericPerm[3:6]
	others := numericPerm[6:9]

	userSum := 0
	for _, c := range user {
		userSum += int(c - '0')
	}

	groupSum := 0
	for _, c := range group {
		groupSum += int(c - '0')
	}

	othersSum := 0
	for _, c := range others {
		othersSum += int(c - '0')
	}

	return fmt.Sprintf("%d%d%d", userSum, groupSum, othersSum)
}

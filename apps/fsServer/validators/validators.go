package validators

import (
	"fmt"
	"simple-fs-web-service/constants"
	"strconv"
	"strings"
)

var defaultOffset = 0
var defaultLimit = 50
var defaultSortField = constants.ModifiedDate
var defaultSortOrder = constants.Descending

func ParseSortOrder(sortOrderQueryParam string) (constants.SortOrder, error) {
	switch sortOrderQueryParam {
	case string(constants.Ascending):
		return constants.Ascending, nil
	case string(constants.Descending):
		return constants.Descending, nil
	case "":
		// Par défaut, utiliser Ascending si le paramètre n'est pas fourni
		return defaultSortOrder, nil
	default:
		return "", fmt.Errorf("invalid sort order: %s", sortOrderQueryParam)
	}
}

func ParseSortField(sortFieldQueryParam string) (constants.SortField, error) {
	switch sortFieldQueryParam {
	case string(constants.ModifiedDate):
		return constants.ModifiedDate, nil
	case string(constants.Name):
		return constants.Name, nil
	case "":
		// Par défaut, utiliser ModifiedDate si le paramètre n'est pas fourni
		return defaultSortField, nil
	default:
		return "", fmt.Errorf("invalid sort field: %s", sortFieldQueryParam)
	}
}

func ParseLimit(limitQueryParam string) (int, error) {
	limitQueryParam = strings.TrimSpace(limitQueryParam)
	if limitQueryParam == "" {
		return defaultLimit, nil
	}
	atoiResult, limitErr := strconv.Atoi(limitQueryParam)
	if limitErr != nil {
		return 0, fmt.Errorf("limit is not a valid integer: %s", limitErr.Error())
	}
	return atoiResult, nil
}

func ParseOffset(offsetQueryParam string) (int, error) {
	offsetQueryParam = strings.TrimSpace(offsetQueryParam)
	if offsetQueryParam == "" {
		return defaultOffset, nil
	}
	atoiResult, offsetErr := strconv.Atoi(offsetQueryParam)
	if offsetErr != nil {
		return 0, fmt.Errorf("offset is not a valid integer: %s", offsetErr.Error())
	}
	return atoiResult, nil
}

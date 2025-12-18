package routesHelpers

import (
	"path"
	"strings"
)

func GetRoutePath(routePath string) string {
	return path.Clean(path.Join(BASE_ROUTER_PATH, routePath))
}

func GetPathSegments(requestPath string, routeMountPoint string) []string {
	directoryPath, _ := strings.CutPrefix(requestPath, strings.Replace(routeMountPoint, "*", "", 1))

	return strings.Split(directoryPath, "/")
}

func GetFSFilePath(requestPath string, routeMountPoint string) string {
	return path.Clean(path.Join(GetPathSegments(requestPath, routeMountPoint)...))
}

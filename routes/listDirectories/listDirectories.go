package listDirectories

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	arrayFilter "simple-fs-web-service/arrays/filter"
	"simple-fs-web-service/constants"
	"slices"
	"strings"
)

func getBasePath() string {
	basePath := os.Getenv("BASE_PATH")
	homePath := os.Getenv("HOME")

	if basePath == "" {
		basePath = homePath
	}

	return basePath
}

// convertPermissions converts a permission string like "rwxr-xr-x" to a numeric format like 755.
func permStringToNumeric(permStr string) string {
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

func isValidDirectoryPath(path string) bool {
	if res := filepath.IsAbs(path); res {
		return false
	}
	return path != ""
}

type SerializableFile struct {
	Name         string `json:"name"`
	IsDir        bool   `json:"isDir"`
	BytesSize    int64  `json:"bytesSize"`
	LastModified string `json:"lastModified"`
	Permissions  string `json:"permissions"`
}
type ListDirectoryResponseSuccess struct {
	Files   []SerializableFile
	Total   int
	MaxPage int
	Page    int
}

type ListDirectoryParams struct {
	Path      string
	Limit     int
	Offset    int
	SortField constants.SortField
	SortOrder constants.SortOrder
}

func ListDirectory(params ListDirectoryParams) (ListDirectoryResponseSuccess, error) {
	limit := params.Limit
	offset := params.Offset
	path := params.Path
	sortField := params.SortField
	sortOrder := params.SortOrder
	if path == "" {
		path = "."
	}

	if validationResult := isValidDirectoryPath(path); !validationResult {
		return ListDirectoryResponseSuccess{}, errors.New("invalid path")
	}
	fsStat, lStatError := fs.Lstat(os.DirFS(getBasePath()), path)

	if lStatError != nil {
		return ListDirectoryResponseSuccess{}, lStatError
	}

	if !fsStat.IsDir() {
		return ListDirectoryResponseSuccess{}, errors.New("path is not a directory")
	}

	var files = make([]SerializableFile, 0)

	entries, readDirError := fs.ReadDir(os.DirFS(getBasePath()), path)
	if readDirError != nil {
		return ListDirectoryResponseSuccess{}, errors.New("error while listing directory")
	}
	entries = arrayFilter.Filter(
		entries,
		func(entry fs.DirEntry) bool {
			return !strings.HasPrefix(entry.Name(), ".")
		},
	)

	var items = make(map[string]struct {
		info  fs.FileInfo
		err   error
		entry fs.DirEntry
	})

	var sortError error
	slices.SortFunc(entries, func(entryI fs.DirEntry, entryJ fs.DirEntry) int {
		if sortError != nil {
			return 0
		}
		iInfo, iError := entryI.Info()
		jInfo, jError := entryJ.Info()
		items[entryI.Name()] = struct {
			info  fs.FileInfo
			err   error
			entry fs.DirEntry
		}{
			info:  iInfo,
			err:   iError,
			entry: entryI,
		}
		items[entryJ.Name()] = struct {
			info  fs.FileInfo
			err   error
			entry fs.DirEntry
		}{
			info:  jInfo,
			err:   jError,
			entry: entryJ,
		}
		if sortField == constants.ModifiedDate {
			if iError != nil || jError != nil {
				sortError = errors.New("error while listing directory")
				return 0
			}
			if sortOrder == constants.Descending {
				return int(jInfo.ModTime().Unix() - iInfo.ModTime().Unix())
			}
			if sortOrder == constants.Ascending {
				return int(iInfo.ModTime().Unix() - jInfo.ModTime().Unix())
			}
		}
		if sortOrder == constants.Descending {
			if entryI.Name() < entryJ.Name() {
				return 1
			}
			return -1
		}

		if entryI.Name() > entryJ.Name() {
			return 1
		}
		return -1
	})

	if sortError != nil {
		return ListDirectoryResponseSuccess{}, sortError
	}

	end := offset + limit
	if offset > len(entries) {
		offset = len(entries)
	}
	if end > len(entries) {
		end = len(entries)
	}
	for _, dirEntry := range entries[offset:end] {
		entry := items[dirEntry.Name()]
		info, error := entry.info, entry.err
		if error != nil {
			return ListDirectoryResponseSuccess{}, errors.New("error while listing directory")
		}
		files = append(files, SerializableFile{
			Name:         entry.entry.Name(),
			IsDir:        info.IsDir(),
			BytesSize:    info.Size(),
			Permissions:  permStringToNumeric(info.Mode().String()),
			LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	maxPage := int(math.Ceil(float64(len(entries)) / float64(limit)))
	if maxPage == 0 {
		maxPage = 1
	}
	return ListDirectoryResponseSuccess{
		Files:   files,
		Total:   len(entries),
		MaxPage: int(math.Ceil(float64(len(entries)) / float64(limit))),
		Page:    offset/limit + 1,
	}, nil

}

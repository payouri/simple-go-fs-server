package listDirectoriesRoute

import (
	"errors"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	arrayFilter "simple-fs-web-service/arrays/filter"
	"simple-fs-web-service/constants"
	"simple-fs-web-service/helpers"
	"slices"
	"strings"
)

func isValidDirectoryPath(path string) bool {
	if filepath.IsAbs(path) {
		return false
	}

	return path != ""
}

type ListDirectoryResponseSuccess struct {
	Files   []helpers.FileMetadata
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

	resolvedPath, resolvePathError := helpers.ResolvePath(helpers.GetBasePath(), path)
	if resolvePathError != nil {
		return ListDirectoryResponseSuccess{}, resolvePathError
	}

	isLegalPathResult, isLegalPathError := helpers.IsPathLegal(helpers.GetBasePath(), resolvedPath)
	if isLegalPathError != nil {
		return ListDirectoryResponseSuccess{}, errors.New("unable to check path")
	}
	if !isLegalPathResult {
		return ListDirectoryResponseSuccess{}, errors.New("unauthorized path")
	}

	dirFs := os.DirFS(resolvedPath)

	fsStat, lStatError := fs.Lstat(dirFs, ".")

	if lStatError != nil {
		return ListDirectoryResponseSuccess{}, lStatError
	}

	if !fsStat.IsDir() {
		return ListDirectoryResponseSuccess{}, errors.New("path is not a directory")
	}

	var files = make([]helpers.FileMetadata, 0)

	entries, readDirError := fs.ReadDir(dirFs, ".")
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
	if len(entries) > 1 {
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
	} else if len(entries) == 1 {
		entry := entries[0]
		info, entryError := entry.Info()
		sortError = entryError

		items[entries[0].Name()] = struct {
			info  fs.FileInfo
			err   error
			entry fs.DirEntry
		}{
			info:  info,
			err:   entryError,
			entry: entry,
		}
	}

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
		dirName := dirEntry.Name()
		entry := items[dirName]
		info, entryError, entryData := entry.info, entry.err, entry.entry
		if entryError != nil {
			return ListDirectoryResponseSuccess{}, errors.New("error while listing directory")
		}
		if entryData == nil {
			log.Println(info, entryError)
			continue
		}

		log.Println(dirEntry.Name())
		files = append(files, helpers.OSStatToFileMetadata(strings.Join([]string{resolvedPath, dirName}, "/"), info))
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

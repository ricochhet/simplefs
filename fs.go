package simplefs

import (
	"bufio"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/otiai10/copy"
)

func Combine(path1 string, path2 ...string) string {
	path := append([]string{path1}, path2...)
	return filepath.Join(path...)
}

func FromCwd(path1 ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := append([]string{wd}, path1...)

	return filepath.Join(path...), nil
}

func GetDirectoryName(fileName string) string {
	return filepath.Dir(fileName)
}

func GetFileName(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func GetFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func GetRelativePath(directories ...string) string {
	result := "./" + directories[0]

	for _, dir := range directories[1:] {
		result = path.Join(result, dir)
	}

	return result
}

func TrimPath(input string) string {
	if strings.HasPrefix(input, "./") || strings.HasPrefix(input, ".\\") {
		return input[2:]
	} else if strings.HasPrefix(input, "/") || strings.HasPrefix(input, "\\") {
		return input[1:]
	}

	return input
}

func Copy(a, b string) error {
	if err := copy.Copy(a, b); err != nil {
		return err
	}

	return nil
}

func Exists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func ReadFile(fileName string) ([]byte, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func ReadAllLines(file *os.File) ([]string, error) {
	return Scan(bufio.NewScanner(file))
}

func ReadAllStringLines(input string) ([]string, error) {
	return Scan(bufio.NewScanner(strings.NewReader(input)))
}

func WriteFile(fileName string, data []byte, perm fs.FileMode) error {
	err := os.WriteFile(fileName, data, perm)
	if err != nil {
		return err
	}

	return nil
}

func WriteToFile(file *os.File, entries []string) error {
	for _, entry := range entries {
		if _, err := file.WriteString(entry); err != nil {
			return err
		}
	}

	return nil
}

func OverwriteFile(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}

func Scan(scanner *bufio.Scanner) ([]string, error) {
	entries := []string{}

	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		entries = append(entries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func DeleteDirectory(fileName string) error {
	err := os.RemoveAll(fileName)
	if err != nil {
		return err
	}

	return nil
}

func SortFileNames(files []string) []string {
	sort.Slice(files, func(i, j int) bool {
		parentA := filepath.Dir(files[i])
		parentB := filepath.Dir(files[j])

		if parentA == parentB {
			return filepath.Base(files[i]) < filepath.Base(files[j])
		}

		return parentA < parentB
	})

	return files
}

func GetFiles(filePath string) []string {
	var files []string

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return []string{}
	}

	return SortFileNames(files)
}

func GetDirectories(filePath string) []string {
	var directories []string

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			directories = append(directories, path)
		}

		return nil
	})
	if err != nil {
		return []string{}
	}

	return SortFileNames(directories)
}

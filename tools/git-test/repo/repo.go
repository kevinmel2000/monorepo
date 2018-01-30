package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/lab46/example/gopkg/exec"
	"github.com/lab46/example/gopkg/print"
)

var (
	baseDir = "lab46/example"

	serviceDir = "exservice"
	services   []string
)

// GetRepoDir return repository directory of the project
func GetRepoDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	indexBase := strings.Index(currentDir, baseDir)
	if indexBase == -1 {
		return "", fmt.Errorf("not in any %s repository path", baseDir)
	}
	repoDir := currentDir[:indexBase] + baseDir
	return repoDir, nil
}

// IsServiceDir check if servicedir is exists within a directory
func IsServiceDir(directory string) bool {
	return strings.Contains(directory, serviceDir)
}

// ServiceList return list of service
func ServiceList() ([]string, error) {
	var currentServiceDir string
	repoDir, err := GetRepoDir()
	if err != nil {
		return nil, err
	}

	if repoDir != "" {
		currentServiceDir = path.Join(repoDir, serviceDir)
	} else {
		currentServiceDir = serviceDir
	}

	var services []string
	fileInfo, err := ioutil.ReadDir(currentServiceDir)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfo {
		if info.IsDir() {
			services = append(services, info.Name())
		}
	}
	print.Debug(services)
	return services, nil
}

// GoList return all package list via 'go list ../...' command
func GoList() ([]string, error) {
	cmd := exec.Command("go", "list", "../...")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	trimmedOutput := strings.Trim(string(output), " \n")

	pkgList := strings.Split(trimmedOutput, "\n")
	return pkgList, nil
}

// PathTtype define type of path
type PathTtype int

// types of dir
const (
	FileType PathTtype = iota
	ServiceType
	GopkgType
	NormalType
)

// Dir struct
type Dir struct {
	Name string
	Type PathTtype
}

// FilterDir return a filtered directory
// TODO: need to optimize this function
func FilterDir(dirs []string) (map[string]Dir, error) {
	var (
		serviceDirLength = len(serviceDir)

		// every changed dir must be a unique dir
		// so map is used instead of array
		changedDir = make(map[string]Dir)
	)

	for _, file := range dirs {
		// skip if not in exservice path for now
		if !strings.Contains(file, serviceDir) {
			continue
		}

		// find the last '/' in the directory
		// trim to find root path of sub-directory
		info, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		f := file

		// trim dir if file is detected
		if !info.IsDir() {
			index := strings.LastIndex(file, "/")
			if index != -1 {
				f = file[:index]
			}
		}

		// looking for service directory
		// TODO: separate this to another function
		indexService := strings.Index(f, fmt.Sprintf("%s", serviceDir))
		if indexService != -1 {
			// splitIndex is index of 'serviceDir' found + the length of 'serviceDir' itself
			// the characters after all of that is usually a '/', so we need to +1 the index
			splitIndex := indexService + serviceDirLength
			p := f[:indexService]

			// ex: exservice/ will be splited to '['/', ''] and exservice/bookapp/book to ['/', 'bookapp/', 'book']
			s := strings.SplitAfter(f[splitIndex:], "/")
			if len(s) >= 2 {
				if s[1] != "" {
					serviceName := path.Join(p, serviceDir, strings.TrimSuffix(s[1], "/"), "")
					dir := Dir{
						Name: serviceName,
						Type: ServiceType,
					}
					changedDir[serviceName] = dir
				}
			}
			continue
		}
		dir := Dir{
			Name: f,
			Type: NormalType,
		}
		changedDir[f] = dir
	}
	return changedDir, nil
}

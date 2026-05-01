package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"slices"
)

const (
	FormatError       = "\033[1;31m[ERROR]"
	FormatSuccess     = "\033[1;32m[SUCCESS]"
	FormatTitle       = "\033[1;34m"
	FormatReset       = "\033[0m"
	PackagesFileName  = "drxboot.packages"
)

var Categories = []string{"server", "minimal", "desktop"}
var Repos = []string{"pacman", "aur"}

type PackageMap map[string]map[string][]string

func printHelp() {
	categoriesStr := strings.Join(Categories, ", ")
	fmt.Printf(`%[1]sUSAGE:%[2]s packages <command> [<arguments>]

%[1]sDESCRIPTION:%[2]s
    A package manager helper for Arch Linux that categorizes and searches for packages
    in Pacman and the AUR, maintaining a local registry for system bootstrapping.

%[1]sCOMMANDS:%[2]s
    add <package>           Adds a package. The script will check if it exists in Pacman or
                            the AUR and ask you which list (%[3]s) to add it to.

    remove <package>        Searches for the package across all lists and removes it.

    search <package>        Checks if you already have a package registered in your lists 
                            and tells you exactly where it is located.

    list [filter]           Shows all the saved packages.
                            Optional filters: %[3]s, pacman, aur.

%[1]sENVIRONMENT:%[2]s
    PACKAGES_PATH           Defines a custom directory to store the package list file.
                            Defaults to the user's home directory.

%[1]sEXAMPLES:%[2]s
    packages add neovim
    packages remove firefox
    packages search htop
    packages list aur
`, FormatTitle, FormatReset, categoriesStr)
}

func getFilePath() string {
	filePath := os.Getenv("PACKAGES_PATH")

	if filePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(FormatError, "The user's home directory could not be determined.", FormatReset)
			os.Exit(1)
		}
		filePath = homeDir
	}

	return filepath.Join(filePath, PackagesFileName)
}


func isInPacman(pkgName string) bool {
	cmd := exec.Command("pacman", "-Si", pkgName)
	err := cmd.Run()
	return err == nil
}


func isInAur(pkgName string) bool {
	url := fmt.Sprintf("https://aur.archlinux.org/rpc/?v=5&type=info&arg=%s", pkgName)
	resp, err := http.Get(url)

	if err != nil { return false }
	defer resp.Body.Close()

	var data struct {
		ResultCount int `json:"resultcount"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false
	}

	return data.ResultCount > 0
}


func newPackageMap() PackageMap {
	packages := make(PackageMap)
	for _, category := range Categories {
		packages[category] = make(map[string][]string)
		for _, repository := range Repos {
			packages[category][repository] = []string{}
		}
	}
	return packages
}


func loadAllPackages() PackageMap {
	packages := newPackageMap()
	filePath := getFilePath()

	file, err := os.Open(filePath)
	if err != nil { return packages }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentCat, currentRepo string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		matched := false
		for _, cat := range Categories {
			for _, repo := range Repos {
				if strings.HasPrefix(line, fmt.Sprintf("%s_%s_packages=(", cat, repo)) {
					currentCat = cat
					currentRepo = repo
					matched = true
					break
				}
			}
			if matched { break }
		}

		if matched {
			continue
		} else if line == ")" {
			currentCat = ""
			currentRepo = ""
		} else if currentCat != "" && currentRepo != "" && line != "" && !strings.HasPrefix(line, "#") {
			packages[currentCat][currentRepo] = append(packages[currentCat][currentRepo], line)
		}
	}

	return packages
}


func saveAllPackages(packages PackageMap) {
	filePath := getFilePath()
	dirPath := filepath.Dir(filePath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("%s The directory for storing packages could not be created: %v%s\n", FormatError, err, FormatReset)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("%s The file for storing packages could not be created: %v%s\n", FormatError, err, FormatReset)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString("#!/bin/bash\n\n")

	for _, cat := range Categories {
		for _, repo := range Repos {
			_, _ = fmt.Fprintf(writer, "%s_%s_packages=(\n", cat, repo)

			pkgMap := make(map[string]bool)
			for _, p := range packages[cat][repo] {
				pkgMap[p] = true
			}
			var sortedPackages []string
			for p := range pkgMap {
				sortedPackages = append(sortedPackages, p)
			}
			sort.Strings(sortedPackages)

			for _, pkg := range sortedPackages {
				_, _ = fmt.Fprintf(writer, "\t%s\n", pkg)
			}
			_, _ = writer.WriteString(")\n\n")
		}
	}
	_ = writer.Flush()
}


func findPackageLocation(pkgName string, allPackages PackageMap) (string, string, string, bool) {
	suffixes := []string{"", "-git", "-bin"}
	for _, cat := range Categories {
		for _, repo := range Repos {
			pkgList := allPackages[cat][repo]
			for _, sfx := range suffixes {
				fullPkgName := pkgName + sfx
				if slices.Contains(pkgList, fullPkgName) {
					return cat, repo, fullPkgName, true
				}
			}
		}
	}
	return "", "", "", false
}


func addPackage(pkgName string) {
	allPackages := loadAllPackages()
	cat, repo, fullPkgName, found := findPackageLocation(pkgName, allPackages)

	if found {
		fmt.Printf("%s Package '%s' is already registered in the '%s' list (%s).%s\n", FormatError, fullPkgName, cat, repo, FormatReset)
		return
	}

	var repoType string
	if isInPacman(pkgName) {
		repoType = "pacman"
	} else if isInAur(pkgName) {
		repoType = "aur"
	} else {
		fmt.Printf("%s Package '%s' was not found in Pacman or the AUR.%s\n", FormatError, pkgName, FormatReset)
		return
	}

	fmt.Printf("Package found in: %s%s%s\n", FormatTitle, strings.ToUpper(repoType), FormatReset)

	var listType string
	categoriesPrompt := strings.Join(Categories, "/")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Which list do you want to add it to? (%s): ", categoriesPrompt)
		input, _ := reader.ReadString('\n')
		listType = strings.ToLower(strings.TrimSpace(input))

		if slices.Contains(Categories, listType) {
			break
		}
	}

	allPackages[listType][repoType] = append(allPackages[listType][repoType], pkgName)
	saveAllPackages(allPackages)
	fmt.Printf("%s Package '%s' successfully added to the '%s' list.%s\n", FormatSuccess, pkgName, listType, FormatReset)
}


func removePackage(pkgName string) {
	allPackages := loadAllPackages()
	cat, repo, fullPkgName, found := findPackageLocation(pkgName, allPackages)

	if found {
		var newList []string
		for _, p := range allPackages[cat][repo] {
			if p != pkgName {
				newList = append(newList, p)
			}
		}
		allPackages[cat][repo] = newList

		saveAllPackages(allPackages)
		fmt.Printf("%s Package '%s' was removed from the '%s' list (%s).%s\n", FormatSuccess, fullPkgName, cat, repo, FormatReset)
	} else {
		fmt.Printf("%s Package '%s' was not found in any of your lists.%s\n", FormatError, pkgName, FormatReset)
	}
}


func searchPackage(pkgName string) {
	allPackages := loadAllPackages()
	cat, repo, fullPkgName, found := findPackageLocation(pkgName, allPackages)

	if found {
		fmt.Printf("%s Package '%s' found in list: %s%s (%s)%s\n", FormatSuccess, fullPkgName, FormatReset, strings.ToUpper(cat), strings.ToUpper(repo), FormatReset)
	} else {
		fmt.Printf("%s Package '%s' was not found in any of your lists.%s\n", FormatError, pkgName, FormatReset)
	}
}


func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}


func listPackages(filterArg string) {
	allPackages := loadAllPackages()
	foundAny := false

	for _, cat := range Categories {
		if filterArg == "" || contains(Repos, filterArg) || filterArg == cat {
			for _, repo := range Repos {
				if filterArg == "" || contains(Categories, filterArg) || filterArg == repo {
					pkgs := allPackages[cat][repo]
					if len(pkgs) > 0 {
						fmt.Printf("\n%s%s (%s):%s\n", FormatTitle, strings.ToUpper(cat), strings.ToUpper(repo), FormatReset)
						
						sort.Strings(pkgs)
						rawText := strings.Join(pkgs, "\n")
						
						cmd := exec.Command("column")
						cmd.Stdin = strings.NewReader(rawText)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						
						if err := cmd.Run(); err != nil {
							fmt.Println(rawText)
						}
						foundAny = true
					}
				}
			}
		}
	}

	if !foundAny {
		fmt.Println("No packages to display with that filter.")
	}
}


func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	action := os.Args[1]

	switch action {
	case "add":
		if len(os.Args) < 3 {
			fmt.Printf("%s You must specify a package: packages add <package>%s\n", FormatError, FormatReset)
			return
		}
		addPackage(os.Args[2])

	case "remove":
		if len(os.Args) < 3 {
			fmt.Printf("%s You must specify a package: packages remove <package>%s\n", FormatError, FormatReset)
			return
		}
		removePackage(os.Args[2])

	case "search":
		if len(os.Args) < 3 {
			fmt.Printf("%s You must specify a package: packages search <package>%s\n", FormatError, FormatReset)
			return
		}
		searchPackage(os.Args[2])

	case "list":
		filterArg := ""
		if len(os.Args) > 2 {
			filterArg = strings.ToLower(os.Args[2])
		}

		var validFilters []string
		validFilters = append(validFilters, Categories...)
		validFilters = append(validFilters, Repos...)

		if filterArg != "" && !contains(validFilters, filterArg) {
			fmt.Printf("%s Invalid filter. Use: %s.%s\n", FormatError, strings.Join(validFilters, ", "), FormatReset)
			return
		}
		listPackages(filterArg)

	default:
		fmt.Printf("%s Invalid action.%s\n", FormatError, FormatReset)
		printHelp()
	}
}

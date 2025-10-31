package system

import (
	"os/exec"
	"strings"
)

type PackageManager string

const (
	APT     PackageManager = "apt"
	DNF     PackageManager = "dnf"
	PACMAN  PackageManager = "pacman"
	UNKNOWN PackageManager = "unknown"
)

type PackageInfo struct {
	Name        string
	Version     string
	Description string
	Size        string
	Status      string
}

func DetectPackageManager() PackageManager {
	if _, err := exec.LookPath("apt"); err == nil {
		return APT
	}
	if _, err := exec.LookPath("dnf"); err == nil {
		return DNF
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return PACMAN
	}
	return UNKNOWN
}

func ListPackages(pm PackageManager) ([]PackageInfo, error) {
	switch pm {
	case APT:
		return listAptPackages()
	case DNF:
		return listDnfPackages()
	case PACMAN:
		return listPacmanPackages()
	default:
		return nil, exec.ErrNotFound
	}
}

func SearchPackages(pm PackageManager, query string) ([]PackageInfo, error) {
	switch pm {
	case APT:
		return searchAptPackages(query)
	case DNF:
		return searchDnfPackages(query)
	case PACMAN:
		return searchPacmanPackages(query)
	default:
		return nil, exec.ErrNotFound
	}
}

func listAptPackages() ([]PackageInfo, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Package}\t${Version}\t${Description}\t${Installed-Size}\n")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) >= 2 {
			packages = append(packages, PackageInfo{
				Name:        parts[0],
				Version:     parts[1],
				Description: parts[2],
				Size:        parts[3],
				Status:      "installed",
			})
		}
	}

	return packages, nil
}

func listDnfPackages() ([]PackageInfo, error) {
	cmd := exec.Command("dnf", "list", "installed", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			packages = append(packages, PackageInfo{
				Name:    parts[0],
				Version: parts[1],
				Status:  "installed",
			})
		}
	}

	return packages, nil
}

func listPacmanPackages() ([]PackageInfo, error) {
	cmd := exec.Command("pacman", "-Q")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			packages = append(packages, PackageInfo{
				Name:    parts[0],
				Version: parts[1],
				Status:  "installed",
			})
		}
	}

	return packages, nil
}

func searchAptPackages(query string) ([]PackageInfo, error) {
	cmd := exec.Command("apt-cache", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			packages = append(packages, PackageInfo{
				Name:        parts[0],
				Description: strings.Join(parts[1:], " "),
			})
		}
	}

	return packages, nil
}

func searchDnfPackages(query string) ([]PackageInfo, error) {
	cmd := exec.Command("dnf", "search", query, "--quiet")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			packages = append(packages, PackageInfo{
				Name:    parts[0],
				Version: parts[1],
			})
		}
	}

	return packages, nil
}

func searchPacmanPackages(query string) ([]PackageInfo, error) {
	cmd := exec.Command("pacman", "-Ss", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	lines := strings.Split(string(output), "\n")
	for i := 0; i < len(lines); i += 2 {
		if i+1 < len(lines) {
			nameLine := strings.Fields(lines[i])
			if len(nameLine) > 0 {
				packages = append(packages, PackageInfo{
					Name:        nameLine[0],
					Description: lines[i+1],
				})
			}
		}
	}

	return packages, nil
}

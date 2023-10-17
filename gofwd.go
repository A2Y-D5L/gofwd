package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/Masterminds/semver/v3"
)

const (
	downloadURL = "https://dl.google.com/go/go%s.%s.tar.gz"
	backupTarball = "go_backup.tar.gz"
	installDir    = "/usr/local"
)

type Release struct {
	Version string `json:"version"`
}

func main() {
	if !hasWritePermission(installDir) {
		log.Fatalf("You don't have write permissions to %s. Please run the program with elevated permissions.", installDir)
	}

	localVersion, err := getLocalVersion()
	if err != nil {
		log.Fatalf("Error parsing local Go version: %v", err)
	}

	latestVersion, err := getLatestVersion()
	if err != nil {
		log.Fatalf("Error fetching latest Go version: %v", err)
	}

	if localVersion.Equal(latestVersion) {
		log.Println("Your Go installation is up to date!")
		return
	}

	if !confirmUpdate(localVersion, latestVersion) {
		log.Println("Update canceled.")
		return
	}

	if err := backupCurrentInstallation(); err != nil {
		log.Fatalf("Error backing up Go installation: %v", err)
	}

	tarballPath, err := downloadLatestVersion(latestVersion)
	if err != nil {
		log.Fatalf("Error downloading Go version %s: %v", latestVersion, err)
	}

	if err := installGoFromTarball(tarballPath); err != nil {
		log.Fatalf("Error installing Go from tarball %s: %v", tarballPath, err)
	}

	if err := verifyInstallation(); err != nil {
		log.Fatalf("Installation verification failed: %v", err)
	}
}

func hasWritePermission(dir string) bool {
	_, err := os.OpenFile(dir+"/testwrite", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
	}
	os.Remove(dir + "/testwrite")
	return true
}

func getLocalVersion() (*semver.Version, error) {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to execute 'go version': %v\n", err)
	}

	parts := strings.Split(string(out), " ")
	if len(parts) < 3 {
		fmt.Printf("Unexpected format from 'go version': %v\n", out)
	}
	return semver.NewVersion(parts[2][2:])
}

func getLatestVersion() (*semver.Version, error) {
	resp, err := http.Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Version > releases[j].Version
	})
	latest := releases[0].Version
	fmt.Printf("Latest stable version: %s\n", latest)
	return semver.NewVersion(strings.TrimPrefix(latest, "go"))
}

func confirmUpdate(local, latest *semver.Version) bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("You are about to update Go from version %s to %s. Do you want to continue?", local.String(), latest.String()),
		IsConfirm: true,
	}

	_, err := prompt.Run()

	return err == nil
}

func backupCurrentInstallation() error {
	cmd := exec.Command("tar", "-czf", "go_backup.tar.gz", "-C", "/usr/local", "go")
	return cmd.Run()
}

func downloadLatestVersion(version *semver.Version) (string, error) {
	platform := runtime.GOOS + "-" + runtime.GOARCH
	downloadURL := fmt.Sprintf("https://dl.google.com/go/go%s.%s.tar.gz", version.String(), platform)
	fmt.Printf("Downloading from: %s\n", downloadURL)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create a temporary directory to store the downloaded tarball
	tmpDir, err := os.MkdirTemp("", "gofwd")
	if err != nil {
		return "", err
	}

	tarPath := filepath.Join(tmpDir, fmt.Sprintf("go%s.%s.tar.gz", version.String(), platform))
	out, err := os.Create(tarPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return tarPath, nil
}

func installGoFromTarball(tarballPath string) error {
	defer os.Remove(tarballPath)
	// Extract the tarball
	file, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join("/usr/local", header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}

	return nil
}

func verifyInstallation() error {
	cmd := exec.Command("/usr/local/go/bin/go", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

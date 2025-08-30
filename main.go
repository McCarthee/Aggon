package main

import (
	"archive/zip"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AddonConfig struct {
	Disabled      bool     `json:"disabled,omitempty"`
	Name          string   `json:"name"`
	URL           string   `json:"url"`
	Folder        string   `json:"folder,omitempty"`
	Ignore        []string `json:"ignore,omitempty"`
	Branch        string   `json:"branch,omitempty"`
	Tag           string   `json:"tag,omitempty"`
	LatestRelease bool     `json:"latest_release,omitempty"`
	AssetPattern  string   `json:"asset_pattern,omitempty"`
}

type DirectoryConfig struct {
	Name            string        `json:"name"`
	Path            string        `json:"path"`
	Addons          []AddonConfig `json:"addons"`
	BackupBlacklist []string      `json:"backup_blacklist,omitempty"`
}

type Config []DirectoryConfig

type CacheEntry struct {
	URL          string    `json:"url"`
	Hash         string    `json:"hash"`
	LastModified time.Time `json:"last_modified"`
	Filename     string    `json:"filename"`
}

type CacheIndex map[string]CacheEntry

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	// Check for command line usage
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add":
			if len(os.Args) > 2 && os.Args[2] == "addon" {
				if err := runAddAddonWizard(); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				return
			} else if len(os.Args) > 2 && os.Args[2] == "path" {
				if err := runAddPathWizard(); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				return
			}
		case "format-config":
			if err := formatConfig(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✨ Config formatted successfully!")
			return
		case "--help", "-h":
			printHelp()
			return
		}
	}

	// Run main menu
	runMainMenu()
}

func runMainMenu() {
	for {
		// Load config
		config, err := loadConfig("config.json")
		if err != nil {
			handleConfigError(err)
			return
		}

		// Display menu
		displayMenu(config)

		// Get user choice
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter choice: ")
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			if len(config) > 0 {
				installAllAddons(config)
			} else {
				fmt.Println("⚠ No installation paths configured. Use option 3 first.")
				waitForEnter()
			}
		case "2":
			if err := runAddAddonWizard(); err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter()
			}
		case "3":
			if err := runAddPathWizard(); err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter()
			}
		case "4":
			backupAllAddons(config)
		case "5":
			if err := formatConfig(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("✨ Config formatted successfully!")
			}
			waitForEnter()
		case "q", "quit", "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
			waitForEnter()
		}
	}
}

func displayMenu(config Config) {
	fmt.Print("\033[H\033[2J") // Clear screen

	fmt.Println("🏺 AGGON")
	fmt.Println("World of Warcraft Addon Manager")
	fmt.Println("================================")
	fmt.Println()

	// Config status
	if len(config) == 0 {
		fmt.Println("⚠ No configuration found")
		fmt.Println("Run option 3 to add an installation path first")
		fmt.Println()
	} else {
		fmt.Printf("📁 %d Installation Path(s) Configured\n", len(config))
		fmt.Println()

		// Show configured paths
		for _, dir := range config {
			fmt.Printf("📂 %s (%d addons)\n", dir.Name, len(dir.Addons))
			fmt.Printf("   %s\n", dir.Path)
		}
		fmt.Println()
	}

	// Menu options
	fmt.Println("Menu Options:")
	fmt.Println("─────────────")
	if len(config) > 0 {
		fmt.Println("1. 🚀 Install/Update All Addons")
	}
	fmt.Println("2. ➕ Add New Addon")
	fmt.Println("3. 📁 Add Installation Path")
	fmt.Println("4. 💾 Backup All Addon Directories")
	fmt.Println("5. ✨ Format Config File")
	fmt.Println("q. Quit")
	fmt.Println()
}

func handleConfigError(err error) {
	fmt.Println("🏺 AGGON")
	fmt.Println("========")
	fmt.Println()
	fmt.Println("❌ Configuration Error")
	fmt.Printf("Error: %v\n", err)
	fmt.Println()
	fmt.Println("Creating sample configuration...")

	createSampleConfig()

	fmt.Println("✅ Sample config.json created!")
	fmt.Println("Please edit it with your addon directories and GitHub URLs, then restart.")
	fmt.Println()
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func installAllAddons(config Config) {
	fmt.Print("\033[H\033[2J") // Clear screen

	fmt.Println("🏺 AGGON")
	fmt.Println("========")
	fmt.Println()
	fmt.Println("🚀 Installing/Updating Addons")
	fmt.Println("=============================")
	fmt.Println()

	var successful, failed, disabled, cached int

	for _, dir := range config {
		fmt.Printf("📂 %s\n", dir.Name)
		fmt.Printf("   %s\n", dir.Path)
		fmt.Println()

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dir.Path, 0755); err != nil {
			fmt.Printf("   ❌ Failed to create directory: %v\n", err)
			continue
		}

		// Setup Aggon directories
		aggonDir := filepath.Join(filepath.Dir(dir.Path), "Aggon")
		cacheDir := filepath.Join(aggonDir, "Cache")
		backupDir := filepath.Join(aggonDir, "Backups")

		if err := setupAggonDirectories(aggonDir, cacheDir, backupDir); err != nil {
			fmt.Printf("   ❌ Failed to setup Aggon directories: %v\n", err)
			continue
		}

		// Load cache index
		cacheIndex := loadCacheIndex(cacheDir)

		// Check which addons will actually need changes
		changesNeeded := false
		for _, addon := range dir.Addons {
			if willAddonChange(addon, dir.Path, cacheDir, cacheIndex) {
				changesNeeded = true
				break
			}
		}

		// Create backup only if changes are needed
		if changesNeeded {
			fmt.Printf("   💾 Changes detected - creating backup before installation...\n")
			if err := backupFullDirectory(dir, backupDir); err != nil {
				fmt.Printf("   ⚠️  Pre-installation backup failed: %v\n", err)
			} else {
				fmt.Printf("   ✅ Pre-installation backup created\n")
			}
		} else {
			fmt.Printf("   ℹ️  No changes needed - skipping backup\n")
		}

		// Process each addon
		for _, addon := range dir.Addons {
			if addon.Disabled {
				// Check if addon is currently installed before trying to uninstall
				if addonExists(addon, dir.Path) {
					fmt.Printf("   🗑️  %s - Uninstalling...", addon.Name)
					if err := uninstallAddon(addon, dir.Path); err != nil {
						fmt.Print("\r\033[K")
						fmt.Printf("   ❌ %s - Uninstall Error: %v\n", addon.Name, err)
						failed++
					} else {
						fmt.Print("\r\033[K")
						fmt.Printf("   🗑️  %s - Uninstalled\n", addon.Name)
						disabled++
					}
				} else {
					fmt.Printf("   ⏭️  %s - Already not installed (disabled)\n", addon.Name)
					disabled++
				}
			} else {
									fmt.Printf("   ⏳ %s - Checking for updates...", addon.Name)
				fromCache, err := installAddonWithCache(addon, dir.Path, cacheDir, cacheIndex)
				// Clear the line completely
				fmt.Print("\r\033[K")
				if err != nil {
					fmt.Printf("   ❌ %s - Error: %v\n", addon.Name, err)
					failed++
				} else {
					if fromCache {
						fmt.Printf("   ✅ %s - Up to date (from cache)\n", addon.Name)
						cached++
					} else {
						fmt.Printf("   ✅ %s - Updated successfully\n", addon.Name)
						successful++
					}
				}
			}
		}

		// Save cache index
		saveCacheIndex(cacheDir, cacheIndex)
		fmt.Println()
	}

	// Summary
	fmt.Println("🎉 Installation Complete!")
	fmt.Println("=========================")
	fmt.Printf("✅ %d addons updated\n", successful)
	fmt.Printf("💾 %d addons up to date (cached)\n", cached)
	if failed > 0 {
		fmt.Printf("❌ %d addons failed\n", failed)
	}
	if disabled > 0 {
		fmt.Printf("🗑️  %d addons uninstalled (disabled)\n", disabled)
	}
	fmt.Println()

	waitForEnter()
}

// New function to determine if an addon will actually change
func willAddonChange(addon AddonConfig, targetDir, cacheDir string, cacheIndex CacheIndex) bool {
	// If addon is disabled and currently exists, it will be uninstalled (change)
	if addon.Disabled {
		return addonExists(addon, targetDir)
	}

	// If addon doesn't exist, it will be installed (change)
	if !addonExists(addon, targetDir) {
		return true
	}

	// Check if we have a cached version that's still valid
	cacheKey := getCacheKey(addon)
	entry, exists := cacheIndex[cacheKey]
	if !exists {
		return true // No cache entry means we'll download (change)
	}

	// Check if cached file exists
	cachedFile := filepath.Join(cacheDir, entry.Filename)
	if _, err := os.Stat(cachedFile); err != nil {
		return true // Cache file missing, will download (change)
	}

	// Check if cache is expired and we need to check for updates
	if addon.LatestRelease {
		// For latest releases, check if cache is older than 1 hour
		return time.Since(entry.LastModified) > time.Hour
	} else {
		// For fixed versions, check if cache is older than 24 hours
		return time.Since(entry.LastModified) > 24*time.Hour
	}
}

func runAddAddonWizard() error {
	config, err := loadConfig("config.json")
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	if len(config) == 0 {
		return fmt.Errorf("no installation directories found in config")
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🚀 Add New Addon")
	fmt.Println("================")
	fmt.Println()

	// Select directory
	fmt.Println("Select installation directory:")
	for i, dir := range config {
		fmt.Printf("%d. %s (%s)\n", i+1, dir.Name, dir.Path)
	}
	fmt.Print("Choose directory (1-" + strconv.Itoa(len(config)) + "): ")

	input, _ := reader.ReadString('\n')
	selectedIndex, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || selectedIndex < 1 || selectedIndex > len(config) {
		return fmt.Errorf("invalid directory selection")
	}
	selectedDirIndex := selectedIndex - 1

	// Get addon details
	fmt.Print("Addon Name: ")
	addonName, _ := reader.ReadString('\n')
	addonName = strings.TrimSpace(addonName)
	if addonName == "" {
		return fmt.Errorf("addon name is required")
	}

	fmt.Print("GitHub URL: ")
	githubURL, _ := reader.ReadString('\n')
	githubURL = strings.TrimSpace(githubURL)
	if !strings.Contains(githubURL, "github.com") {
		return fmt.Errorf("must be a valid GitHub URL")
	}

	fmt.Print("Custom Folder Name (optional, press Enter to skip): ")
	folder, _ := reader.ReadString('\n')
	folder = strings.TrimSpace(folder)

	fmt.Print("Ignore Files (comma-separated, optional): ")
	ignoreInput, _ := reader.ReadString('\n')
	ignoreInput = strings.TrimSpace(ignoreInput)

	fmt.Print("Use Latest Release? (y/N): ")
	releaseInput, _ := reader.ReadString('\n')
	latestRelease := strings.ToLower(strings.TrimSpace(releaseInput)) == "y"

	var assetPattern, tag, branch string

	if latestRelease {
		fmt.Print("Asset Pattern (optional, for multiple release assets): ")
		assetPattern, _ = reader.ReadString('\n')
		assetPattern = strings.TrimSpace(assetPattern)
	} else {
		fmt.Print("Specific Tag (optional, press Enter to skip): ")
		tag, _ = reader.ReadString('\n')
		tag = strings.TrimSpace(tag)

		if tag == "" {
			fmt.Print("Specific Branch (optional, press Enter for default): ")
			branch, _ = reader.ReadString('\n')
			branch = strings.TrimSpace(branch)
		}
	}

	// Process ignore files
	var ignoreFiles []string
	if ignoreInput != "" {
		for _, file := range strings.Split(ignoreInput, ",") {
			ignoreFiles = append(ignoreFiles, strings.TrimSpace(file))
		}
	}

	// Create new addon
	newAddon := AddonConfig{
		Name: addonName,
		URL:  githubURL,
	}

	if len(ignoreFiles) > 0 {
		newAddon.Ignore = ignoreFiles
	}
	if folder != "" {
		newAddon.Folder = folder
	}
	if latestRelease {
		newAddon.LatestRelease = true
		if assetPattern != "" {
			newAddon.AssetPattern = assetPattern
		}
	}
	if tag != "" {
		newAddon.Tag = tag
	}
	if branch != "" {
		newAddon.Branch = branch
	}

	// Add to config
	config[selectedDirIndex].Addons = append(config[selectedDirIndex].Addons, newAddon)

	// Save config
	if err := saveConfig("config.json", config); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("✅ Addon added successfully!")
	return nil
}

func backupAllAddons(config Config) {
	fmt.Print("\033[H\033[2J") // Clear screen

	fmt.Println("🏺 AGGON")
	fmt.Println("========")
	fmt.Println()
	fmt.Println("💾 Backing Up All Addon Directories")
	fmt.Println("===================================")
	fmt.Println()

	var successful, failed int

	for _, dir := range config {
		fmt.Printf("📂 %s\n", dir.Name)
		fmt.Printf("   %s\n", dir.Path)
		fmt.Println()

		// Setup Aggon directories
		aggonDir := filepath.Join(filepath.Dir(dir.Path), "Aggon")
		backupDir := filepath.Join(aggonDir, "Backups")

		if err := setupAggonDirectories(aggonDir, backupDir); err != nil {
			fmt.Printf("   ❌ Failed to setup Aggon directories: %v\n", err)
			failed++
			continue
		}

		// Check if addon directory exists
		if _, err := os.Stat(dir.Path); os.IsNotExist(err) {
			fmt.Printf("   ⏭️  Addon directory doesn't exist, skipping\n")
			continue
		}

		fmt.Printf("   💾 Creating full directory backup...\n")
		if err := backupFullDirectory(dir, backupDir); err != nil {
			fmt.Printf("   ❌ Backup failed: %v\n", err)
			failed++
		} else {
			fmt.Printf("   ✅ Backup completed successfully\n")
			successful++
		}
		fmt.Println()
	}

	// Summary
	fmt.Println("🎉 Backup Complete!")
	fmt.Println("===================")
	fmt.Printf("✅ %d directories backed up successfully\n", successful)
	if failed > 0 {
		fmt.Printf("❌ %d directories failed to backup\n", failed)
	}
	fmt.Println()

	waitForEnter()
}

func setupAggonDirectories(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func loadCacheIndex(cacheDir string) CacheIndex {
	indexPath := filepath.Join(cacheDir, "index.json")
	file, err := os.Open(indexPath)
	if err != nil {
		return make(CacheIndex)
	}
	defer file.Close()

	var index CacheIndex
	if err := json.NewDecoder(file).Decode(&index); err != nil {
		return make(CacheIndex)
	}
	if index == nil {
		index = make(CacheIndex)
	}
	return index
}

func saveCacheIndex(cacheDir string, index CacheIndex) error {
	indexPath := filepath.Join(cacheDir, "index.json")
	file, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(index)
}

func getCacheKey(addon AddonConfig) string {
	// Create a unique cache key based on addon configuration
	key := fmt.Sprintf("%s|%s|%s|%s|%v|%s",
		addon.Name, addon.URL, addon.Branch, addon.Tag, addon.LatestRelease, addon.AssetPattern)
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])[:16] // Use first 16 chars of hash
}

func installAddonWithCache(addon AddonConfig, targetDir, cacheDir string, cacheIndex CacheIndex) (bool, error) {
	cacheKey := getCacheKey(addon)

	// Get current download URL
	downloadURL, err := getDownloadURL(addon)
	if err != nil {
		return false, fmt.Errorf("failed to get download URL: %v", err)
	}

	// Check if we have a cached version
	if entry, exists := cacheIndex[cacheKey]; exists {
		cachedFile := filepath.Join(cacheDir, entry.Filename)

		// Check if cached file exists and URL matches
		if _, err := os.Stat(cachedFile); err == nil && entry.URL == downloadURL {
			// Determine if we should check for updates
			shouldUpdate := false
			if addon.LatestRelease {
				// Check for updates if cache is older than 1 hour
				shouldUpdate = time.Since(entry.LastModified) > time.Hour
			} else {
				// For fixed versions, check if cache is older than 24 hours
				shouldUpdate = time.Since(entry.LastModified) > 24*time.Hour
			}

			if !shouldUpdate {
				// Use cached version - but still extract in case files were deleted
				return true, extractZip(cachedFile, targetDir, addon)
			}
		}
	}

	// Download fresh copy
	resp, err := http.Get(downloadURL)
	if err != nil {
		return false, fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Create cache filename
	timestamp := time.Now().Format("20060102-150405")
	cacheFilename := fmt.Sprintf("%s-%s.zip", cacheKey, timestamp)
	cachePath := filepath.Join(cacheDir, cacheFilename)

	// Save to cache
	cacheFile, err := os.Create(cachePath)
	if err != nil {
		return false, fmt.Errorf("failed to create cache file: %v", err)
	}
	defer cacheFile.Close()

	// Copy response to cache file and calculate hash
	hasher := sha256.New()
	writer := io.MultiWriter(cacheFile, hasher)

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		os.Remove(cachePath) // Clean up on error
		return false, fmt.Errorf("failed to save download: %v", err)
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	// Check if this is actually a new version by comparing hashes
	if entry, exists := cacheIndex[cacheKey]; exists && entry.Hash == hash {
		// Same content, just update timestamp and use existing cache
		cacheIndex[cacheKey] = CacheEntry{
			URL:          downloadURL,
			Hash:         hash,
			LastModified: time.Now(),
			Filename:     entry.Filename, // Keep using the old filename
		}
		os.Remove(cachePath) // Remove the duplicate file

		// Extract from existing cache
		existingCachePath := filepath.Join(cacheDir, entry.Filename)
		if _, err := os.Stat(existingCachePath); err == nil {
			return true, extractZip(existingCachePath, targetDir, addon)
		}
		// Fall through to use new file if old one is missing
		cachePath = existingCachePath
		cacheFilename = entry.Filename
	}

	// Update cache index with new file
	cacheIndex[cacheKey] = CacheEntry{
		URL:          downloadURL,
		Hash:         hash,
		LastModified: time.Now(),
		Filename:     cacheFilename,
	}

	// Clean up old cache files for this addon
	cleanupOldCacheFiles(cacheDir, cacheKey, cacheFilename)

	// Extract from cache
	return false, extractZip(cachePath, targetDir, addon)
}

func cleanupOldCacheFiles(cacheDir, cacheKey, currentFilename string) {
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}

	prefix := cacheKey + "-"
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && entry.Name() != currentFilename {
			os.Remove(filepath.Join(cacheDir, entry.Name()))
		}
	}
}

func addonExists(addon AddonConfig, targetDir string) bool {
	if addon.Folder != "" {
		addonPath := filepath.Join(targetDir, addon.Folder)
		_, err := os.Stat(addonPath)
		return err == nil
	}

	addonDirs, err := findAddonDirectories(addon, targetDir)
	if err != nil {
		return false
	}
	return len(addonDirs) > 0
}

// New function to backup entire directory with blacklist support
func backupFullDirectory(dirConfig DirectoryConfig, backupDir string) error {
	// Create backup with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupName := fmt.Sprintf("%s-full-%s.zip", sanitizeFilename(dirConfig.Name), timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	// Create zip file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Get list of addon directories to backup (with blacklist filtering)
	addonDirs, err := getAddonDirectoriesToBackup(dirConfig)
	if err != nil {
		return err
	}

	// Add each addon directory to the zip
	for _, addonPath := range addonDirs {
		addonName := filepath.Base(addonPath)
		err = addDirToZip(zipWriter, addonPath, addonName)
		if err != nil {
			return fmt.Errorf("failed to add %s to backup: %v", addonName, err)
		}
	}

	// Clean up old backups (keep last 5 full backups)
	cleanupOldFullBackups(backupDir, dirConfig.Name)

	return nil
}

// Get list of addon directories, filtering out blacklisted items
func getAddonDirectoriesToBackup(dirConfig DirectoryConfig) ([]string, error) {
	var addonDirs []string

	entries, err := os.ReadDir(dirConfig.Path)
	if err != nil {
		return nil, err
	}

	// Set default blacklist if none specified
	blacklist := dirConfig.BackupBlacklist
	if len(blacklist) == 0 {
		blacklist = getDefaultBlacklist()
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		addonName := entry.Name()

		// Check if addon should be excluded
		if shouldExcludeFromBackup(addonName, blacklist) {
			continue
		}

		addonDirs = append(addonDirs, filepath.Join(dirConfig.Path, addonName))
	}

	return addonDirs, nil
}

// Check if addon should be excluded from backup based on blacklist
func shouldExcludeFromBackup(addonName string, blacklist []string) bool {
	for _, pattern := range blacklist {
		if matchesPattern(addonName, pattern) {
			return true
		}
	}
	return false
}

// Pattern matching with wildcard support
func matchesPattern(text, pattern string) bool {
	// Handle exact match first
	if text == pattern {
		return true
	}

	// Handle wildcard patterns
	if strings.Contains(pattern, "*") {
		// Simple wildcard matching
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			return strings.HasPrefix(text, prefix)
		}

		if strings.HasPrefix(pattern, "*") {
			suffix := strings.TrimPrefix(pattern, "*")
			return strings.HasSuffix(text, suffix)
		}

		// Handle middle wildcards (split on * and check parts)
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			return strings.HasPrefix(text, parts[0]) && strings.HasSuffix(text, parts[1])
		}
	}

	return false
}

// Default blacklist for common system/default addons
func getDefaultBlacklist() []string {
	return []string{
		"Blizzard_*",
		"!BugGrabber",
		"!Swatter",
		".DS_Store",
		"Thumbs.db",
	}
}

// Clean up old full backups (keep last 5)
func cleanupOldFullBackups(backupDir, installName string) {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return
	}

	sanitizedName := sanitizeFilename(installName)
	prefix := sanitizedName + "-full-"

	var backupFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), ".zip") {
			backupFiles = append(backupFiles, entry)
		}
	}

	// Keep only the 5 most recent backups
	if len(backupFiles) > 5 {
		// Sort by name (which includes timestamp) and remove oldest
		for i := 0; i < len(backupFiles)-5; i++ {
			os.Remove(filepath.Join(backupDir, backupFiles[i].Name()))
		}
	}
}

func addDirToZip(zipWriter *zip.Writer, srcDir, baseInZip string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		zipPath := filepath.Join(baseInZip, relPath)
		zipPath = strings.ReplaceAll(zipPath, "\\", "/") // Normalize for zip

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}

func sanitizeFilename(name string) string {
	// Replace invalid filename characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

func runAddPathWizard() error {
	config, err := loadConfig("config.json")
	if err != nil {
		config = Config{}
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📁 Add Installation Path")
	fmt.Println("=======================")
	fmt.Println()

	fmt.Print("Installation Name: ")
	installName, _ := reader.ReadString('\n')
	installName = strings.TrimSpace(installName)
	if installName == "" {
		return fmt.Errorf("installation name is required")
	}

	fmt.Print("Installation Path: ")
	installPath, _ := reader.ReadString('\n')
	installPath = strings.TrimSpace(installPath)
	if installPath == "" {
		return fmt.Errorf("installation path is required")
	}

	// Normalize path (convert backslashes to forward slashes)
	installPath = strings.ReplaceAll(installPath, "\\", "/")

	fmt.Print("Backup Blacklist (comma-separated patterns, optional): ")
	blacklistInput, _ := reader.ReadString('\n')
	blacklistInput = strings.TrimSpace(blacklistInput)

	// Process blacklist
	var blacklist []string
	if blacklistInput != "" {
		for _, pattern := range strings.Split(blacklistInput, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				blacklist = append(blacklist, pattern)
			}
		}
	}

	// Create new directory config
	newDir := DirectoryConfig{
		Name:            installName,
		Path:            installPath,
		Addons:          []AddonConfig{},
		BackupBlacklist: blacklist,
	}

	// Add to config
	config = append(config, newDir)

	// Save config
	if err := saveConfig("config.json", config); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Print("Add an addon to this path now? (y/N): ")
	addAddonInput, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(addAddonInput)) == "y" {
		fmt.Println("✅ Installation path added!")
		fmt.Println()
		return runAddAddonWizard()
	}

	fmt.Println("✅ Installation path added!")
	return nil
}

func printHelp() {
	fmt.Println("🏺 AGGON - WoW Addon Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  aggon                    Start interactive menu")
	fmt.Println("  aggon add addon          Add addon")
	fmt.Println("  aggon add path           Add path")
	fmt.Println("  aggon format-config      Format config file")
	fmt.Println("  aggon --help             Show this help")
}

func waitForEnter() {
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func loadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func createSampleConfig() {
	config := Config{
		{
			Name: "Retail",
			Path: "/path/to/wow/_retail_/Interface/AddOns",
			BackupBlacklist: []string{
				"Blizzard_*",
				"!BugGrabber",
				"!Swatter",
			},
			Addons: []AddonConfig{
				{
					Name: "WeakAuras",
					URL:  "https://github.com/WeakAuras/WeakAuras2",
				},
				{
					Name: "Details",
					URL:  "https://github.com/Tercioo/Details-Damage-Meter",
				},
				{
					Name:   "AdiBags - Mods",
					URL:    "https://github.com/Sattva-108/AdiBags-WoTLK-3.3.5-Mods",
					Ignore: []string{"README.md", ".gitignore", "LICENSE"},
				},
				{
					Name:   "WeakAuras Beta",
					URL:    "https://github.com/WeakAuras/WeakAuras2",
					Branch: "development",
					Folder: "WeakAuras-Beta",
				},
				{
					Name:          "BigWigs",
					URL:           "https://github.com/BigWigsMods/BigWigs",
					LatestRelease: true,
				},
				{
					Disabled: true,
					Name:     "pfQuest - WoTLK",
					URL:      "https://github.com/shagu/pfQuest",
					Folder:   "pfQuest",
				},
			},
		},
		{
			Name: "Classic",
			Path: "/path/to/wow/_classic_/Interface/AddOns",
			BackupBlacklist: []string{
				"Blizzard_*",
			},
			Addons: []AddonConfig{
				{
					Name: "ClassicCastbars",
					URL:  "https://github.com/wardz/ClassicCastbars",
				},
			},
		},
	}

	file, err := os.Create("config.json")
	if err != nil {
		return
	}
	defer file.Close()

	saveConfigFormatted(file, config)
}

func getDownloadURL(addon AddonConfig) (string, error) {
	githubURL := strings.TrimSuffix(addon.URL, "/")

	if !strings.Contains(githubURL, "github.com") {
		return "", fmt.Errorf("not a GitHub URL")
	}

	var downloadURL string

	if addon.LatestRelease {
		releaseURL, err := getLatestReleaseURL(addon, githubURL)
		if err != nil {
			return "", fmt.Errorf("failed to get latest release: %v", err)
		}
		downloadURL = releaseURL
	} else if addon.Tag != "" {
		downloadURL = githubURL + "/archive/refs/tags/" + addon.Tag + ".zip"
	} else if addon.Branch != "" {
		downloadURL = githubURL + "/archive/refs/heads/" + addon.Branch + ".zip"
	} else {
		downloadURL = githubURL + "/archive/refs/heads/main.zip"
		resp, err := http.Head(downloadURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			downloadURL = githubURL + "/archive/refs/heads/master.zip"
		}
	}

	return downloadURL, nil
}

func extractZip(src, dest string, addon AddonConfig) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	var rootFolder string
	if len(reader.File) > 0 {
		rootFolder = strings.Split(reader.File[0].Name, "/")[0]
	}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() || strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}

		relativePath := strings.TrimPrefix(file.Name, rootFolder+"/")
		if relativePath == file.Name {
			relativePath = file.Name
		}

		if shouldIgnoreFile(relativePath, addon.Ignore) {
			continue
		}

		var destPath string
		if addon.Folder != "" {
			destPath = filepath.Join(dest, addon.Folder, relativePath)
		} else {
			destPath = filepath.Join(dest, relativePath)
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func shouldIgnoreFile(filePath string, ignoreList []string) bool {
	if len(ignoreList) == 0 {
		return false
	}

	fileName := filepath.Base(filePath)

	for _, ignorePattern := range ignoreList {
		if fileName == ignorePattern || filePath == ignorePattern || strings.Contains(filePath, ignorePattern) {
			return true
		}
	}

	return false
}

func getLatestReleaseURL(addon AddonConfig, githubURL string) (string, error) {
	parts := strings.Split(githubURL, "/")
	if len(parts) < 5 {
		return "", fmt.Errorf("invalid GitHub URL format")
	}

	owner := parts[3]
	repo := parts[4]
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release info: %v", err)
	}

	if len(release.Assets) == 0 {
		return "", fmt.Errorf("no assets found in latest release")
	}

	if addon.AssetPattern != "" {
		for _, asset := range release.Assets {
			if strings.Contains(strings.ToLower(asset.Name), strings.ToLower(addon.AssetPattern)) {
				return asset.BrowserDownloadURL, nil
			}
		}
		return "", fmt.Errorf("no asset matching pattern '%s' found", addon.AssetPattern)
	}

	if len(release.Assets) > 1 {
		var assetNames []string
		for _, asset := range release.Assets {
			assetNames = append(assetNames, asset.Name)
		}
		return "", fmt.Errorf("multiple assets found, please specify asset_pattern. Available assets: %s", strings.Join(assetNames, ", "))
	}

	return release.Assets[0].BrowserDownloadURL, nil
}

func saveConfig(filename string, config Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return saveConfigFormatted(file, config)
}

func saveConfigFormatted(file *os.File, config Config) error {
	// Custom JSON formatting with 4 spaces indentation and special array handling
	output := "[\n"
	for dirIndex, dir := range config {
		output += "    {\n"
		output += fmt.Sprintf("        \"name\": %q,\n", dir.Name)
		output += fmt.Sprintf("        \"path\": %q,\n", dir.Path)

		// Add backup_blacklist if present
		if len(dir.BackupBlacklist) > 0 {
			if len(dir.BackupBlacklist) <= 5 {
				// Single line format for 5 or fewer items
				var items []string
				for _, item := range dir.BackupBlacklist {
					items = append(items, fmt.Sprintf("%q", item))
				}
				output += fmt.Sprintf("        \"backup_blacklist\": [ %s ],\n", strings.Join(items, ", "))
			} else {
				// Multi-line format for more than 5 items
				output += "        \"backup_blacklist\": [\n"
				for i, item := range dir.BackupBlacklist {
					if i == len(dir.BackupBlacklist)-1 {
						output += fmt.Sprintf("            %q\n", item)
					} else {
						output += fmt.Sprintf("            %q,\n", item)
					}
				}
				output += "        ],\n"
			}
		}

		output += "        \"addons\": [\n"

		for addonIndex, addon := range dir.Addons {
			output += "            {\n"

			var fields []string

			// Field order: disabled (if present), name, url, folder, ignore, branch, tag, latest_release, asset_pattern

			// 1. disabled (only if true)
			if addon.Disabled {
				fields = append(fields, "                \"disabled\": true")
			}

			// 2. name (always present)
			fields = append(fields, fmt.Sprintf("                \"name\": %q", addon.Name))

			// 3. url (always present)
			fields = append(fields, fmt.Sprintf("                \"url\": %q", addon.URL))

			// 4. folder (optional)
			if addon.Folder != "" {
				fields = append(fields, fmt.Sprintf("                \"folder\": %q", addon.Folder))
			}

			// 5. ignore (optional, with special formatting)
			if len(addon.Ignore) > 0 {
				if len(addon.Ignore) <= 5 {
					// Single line format for 5 or fewer items with proper spacing
					var items []string
					for _, item := range addon.Ignore {
						items = append(items, fmt.Sprintf("%q", item))
					}
					fields = append(fields, fmt.Sprintf("                \"ignore\": [ %s ]", strings.Join(items, ", ")))
				} else {
					// Multi-line format for more than 5 items
					ignoreField := "                \"ignore\": [\n"
					for i, item := range addon.Ignore {
						if i == len(addon.Ignore)-1 {
							ignoreField += fmt.Sprintf("                    %q\n", item)
						} else {
							ignoreField += fmt.Sprintf("                    %q,\n", item)
						}
					}
					ignoreField += "                ]"
					fields = append(fields, ignoreField)
				}
			}

			// 6. branch (optional)
			if addon.Branch != "" {
				fields = append(fields, fmt.Sprintf("                \"branch\": %q", addon.Branch))
			}

			// 7. tag (optional)
			if addon.Tag != "" {
				fields = append(fields, fmt.Sprintf("                \"tag\": %q", addon.Tag))
			}

			// 8. latest_release (optional)
			if addon.LatestRelease {
				fields = append(fields, "                \"latest_release\": true")
			}

			// 9. asset_pattern (optional)
			if addon.AssetPattern != "" {
				fields = append(fields, fmt.Sprintf("                \"asset_pattern\": %q", addon.AssetPattern))
			}

			// Join fields with commas
			output += strings.Join(fields, ",\n")
			output += "\n            }"
			if addonIndex < len(dir.Addons)-1 {
				output += ","
			}
			output += "\n"
		}

		output += "        ]\n"
		output += "    }"
		if dirIndex < len(config)-1 {
			output += ","
		}
		output += "\n"
	}
	output += "]\n"

	_, err := file.WriteString(output)
	return err
}

func formatConfig() error {
	config, err := loadConfig("config.json")
	if err != nil {
		return err
	}

	return saveConfig("config.json", config)
}

func uninstallAddon(addon AddonConfig, targetDir string) error {
	if addon.Folder != "" {
		addonPath := filepath.Join(targetDir, addon.Folder)
		if _, err := os.Stat(addonPath); err == nil {
			return os.RemoveAll(addonPath)
		}
		return nil
	}

	addonDirs, err := findAddonDirectories(addon, targetDir)
	if err != nil {
		return fmt.Errorf("failed to find addon directories: %v", err)
	}

	if len(addonDirs) == 0 {
		return nil
	}

	for _, dir := range addonDirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove directory %s: %v", dir, err)
		}
	}

	return nil
}

func findAddonDirectories(addon AddonConfig, targetDir string) ([]string, error) {
	var addonDirs []string

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return addonDirs, nil
		}
		return nil, err
	}

	addonNameLower := strings.ToLower(addon.Name)
	repoName := getRepoName(addon.URL)
	repoNameLower := strings.ToLower(repoName)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()
		dirNameLower := strings.ToLower(dirName)

		if strings.Contains(dirNameLower, addonNameLower) ||
			strings.Contains(dirNameLower, repoNameLower) ||
			strings.Contains(addonNameLower, dirNameLower) {
			addonDirs = append(addonDirs, filepath.Join(targetDir, dirName))
		}
	}

	return addonDirs, nil
}

func getRepoName(githubURL string) string {
	parts := strings.Split(strings.TrimSuffix(githubURL, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

package main

import (
	"time"
)

// DeclarativeConfig represents the complete system configuration
type DeclarativeConfig struct {
	Schema       string                         `json:"schema"`
	Metadata     ConfigMetadata                 `json:"metadata"`
	Installations map[string]InstallationConfig `json:"installations"`
	Addons       map[string]AddonConfig         `json:"addons"`
	Profiles     map[string]ProfileConfig       `json:"profiles"`
	Settings     SystemSettings                 `json:"settings"`
}

// ConfigMetadata contains information about the configuration
type ConfigMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// InstallationConfig defines a WoW installation
type InstallationConfig struct {
	Type    string   `json:"type"`    // retail, classic, custom
	Path    string   `json:"path"`    // Where AddOns folder is located
	Enabled bool     `json:"enabled"` // Whether this installation is active
	Addons  []string `json:"addons"`  // List of addon IDs to install
}

// AddonConfig defines an addon and its source
type AddonConfig struct {
	Source     AddonSource `json:"source"`     // Where to get the addon
	Version    string      `json:"version"`    // Version constraint
	Hash       *string     `json:"hash"`       // Content hash for reproducibility
	Ignore     []string    `json:"ignore"`     // Files to ignore when extracting
	Compatible []string    `json:"compatible"` // Compatible installation types
	Folder     *string     `json:"folder"`     // Custom folder name
}

// AddonSource defines where an addon comes from
type AddonSource struct {
	Type string `json:"type"` // github, url, local
	URL  string `json:"url"`  // Source URL
	Ref  string `json:"ref"`  // Branch, tag, or commit
}

// ProfileConfig defines a profile with specific settings
type ProfileConfig struct {
	Installations []string                   `json:"installations"` // Which installations to use
	Addons        map[string]AddonOverride   `json:"addons"`        // Addon-specific overrides
	Description   string                     `json:"description"`   // Profile description
}

// AddonOverride allows profiles to override addon settings
type AddonOverride struct {
	Ref     *string `json:"ref,omitempty"`     // Override source ref
	Version *string `json:"version,omitempty"` // Override version
	Enabled *bool   `json:"enabled,omitempty"` // Override enabled state
}

// SystemSettings contains global system configuration
type SystemSettings struct {
	AutoUpdate        bool   `json:"auto_update"`        // Automatically update addons
	BackupGenerations int    `json:"backup_generations"` // Number of generations to keep
	ParallelDownloads int    `json:"parallel_downloads"` // Number of concurrent downloads
	VerifyHashes      bool   `json:"verify_hashes"`      // Verify content hashes
	StorePath         string `json:"store_path"`         // Path to content store
	GenerationsPath   string `json:"generations_path"`   // Path to generations
}

// Generation represents a system state snapshot
type Generation struct {
	ID          int                 `json:"id"`
	Timestamp   time.Time           `json:"timestamp"`
	Config      DeclarativeConfig   `json:"config"`
	StateHash   string              `json:"state_hash"`
	Description string              `json:"description"`
	Installations map[string]InstallationState `json:"installations"`
}

// InstallationState tracks the state of an installation
type InstallationState struct {
	Path   string                    `json:"path"`
	Addons map[string]InstalledAddon `json:"addons"`
}

// InstalledAddon represents an addon that's actually installed
type InstalledAddon struct {
	ID           string    `json:"id"`
	Version      string    `json:"version"`
	Hash         string    `json:"hash"`
	StorePath    string    `json:"store_path"`
	InstallPath  string    `json:"install_path"`
	InstalledAt  time.Time `json:"installed_at"`
}

// StoreEntry represents an addon stored in the content store
type StoreEntry struct {
	Hash        string    `json:"hash"`
	Size        int64     `json:"size"`
	StorePath   string    `json:"store_path"`
	SourceURL   string    `json:"source_url"`
	SourceRef   string    `json:"source_ref"`
	CreatedAt   time.Time `json:"created_at"`
	AccessedAt  time.Time `json:"accessed_at"`
}

// BuildPlan represents the changes needed to achieve desired state
type BuildPlan struct {
	CurrentGeneration int                      `json:"current_generation"`
	TargetGeneration  int                      `json:"target_generation"`
	Operations        []Operation              `json:"operations"`
	Downloads         []DownloadOperation      `json:"downloads"`
	Installations     map[string]InstallPlan   `json:"installations"`
}

// Operation represents a single change operation
type Operation struct {
	Type        string      `json:"type"`        // install, uninstall, update, symlink
	Installation string      `json:"installation"` // Which installation
	Addon       string      `json:"addon"`       // Which addon
	From        *string     `json:"from"`        // Current state (for updates)
	To          string      `json:"to"`          // Target state
	Data        interface{} `json:"data"`        // Operation-specific data
}

// DownloadOperation represents a download that needs to happen
type DownloadOperation struct {
	AddonID   string      `json:"addon_id"`
	Source    AddonSource `json:"source"`
	Hash      string      `json:"hash"`
	StorePath string      `json:"store_path"`
}

// InstallPlan represents changes needed for one installation
type InstallPlan struct {
	Path       string              `json:"path"`
	Operations []Operation         `json:"operations"`
	Symlinks   map[string]string   `json:"symlinks"` // addon -> store_path
}

// ReconcileResult contains the results of a reconciliation
type ReconcileResult struct {
	Success       bool              `json:"success"`
	Generation    int               `json:"generation"`
	Operations    int               `json:"operations"`
	Downloaded    int               `json:"downloaded"`
	Installed     int               `json:"installed"`
	Errors        []error           `json:"errors"`
	Duration      time.Duration     `json:"duration"`
}

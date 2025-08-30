package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"
)

// Reconciler handles the reconciliation between desired and actual state
type Reconciler struct {
	Config             DeclarativeConfig
	Store              *Store
	GenerationManager  *GenerationManager
}

// NewReconciler creates a new reconciler
func NewReconciler(config DeclarativeConfig, store *Store, gm *GenerationManager) *Reconciler {
	return &Reconciler{
		Config:            config,
		Store:             store,
		GenerationManager: gm,
	}
}

// Plan creates a build plan to achieve the desired state
func (r *Reconciler) Plan() (*BuildPlan, error) {
	plan := &BuildPlan{
		Operations:    []Operation{},
		Downloads:     []DownloadOperation{},
		Installations: make(map[string]InstallPlan),
	}
	
	// Get current generation (if any)
	current, err := r.GenerationManager.GetCurrent()
	if err == nil {
		plan.CurrentGeneration = current.ID
	}
	
	// Process each enabled installation
	for installID, installConfig := range r.Config.Installations {
		if !installConfig.Enabled {
			continue
		}
		
		installPlan, err := r.planInstallation(installID, installConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to plan installation %s: %v", installID, err)
		}
		
		plan.Installations[installID] = *installPlan
		plan.Operations = append(plan.Operations, installPlan.Operations...)
	}
	
	// Collect required downloads
	for _, installPlan := range plan.Installations {
		for _, op := range installPlan.Operations {
			if op.Type == "install" || op.Type == "update" {
				addonConfig, exists := r.Config.Addons[op.Addon]
				if !exists {
					continue
				}
				
				// Check if we need to download this addon
				hash := r.calculateAddonHash(addonConfig)
				if !r.Store.Exists(hash) {
					download := DownloadOperation{
						AddonID:   op.Addon,
						Source:    addonConfig.Source,
						Hash:      hash,
						StorePath: r.Store.getStorePath(hash),
					}
					plan.Downloads = append(plan.Downloads, download)
				}
			}
		}
	}
	
	return plan, nil
}

// Apply executes a build plan
func (r *Reconciler) Apply(plan *BuildPlan) (*ReconcileResult, error) {
	startTime := time.Now()
	result := &ReconcileResult{
		Success:    true,
		Operations: len(plan.Operations),
		Errors:     []error{},
	}
	
	// Create new generation
	generation, err := r.GenerationManager.Create(r.Config, "Applied declarative configuration")
	if err != nil {
		return result, fmt.Errorf("failed to create generation: %v", err)
	}
	result.Generation = generation.ID
	
	// Execute downloads
	for _, download := range plan.Downloads {
		if err := r.executeDownload(download); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("download %s failed: %v", download.AddonID, err))
			result.Success = false
			continue
		}
		result.Downloaded++
	}
	
	// Execute installation plans
	for installID, installPlan := range plan.Installations {
		if err := r.executeInstallPlan(installID, installPlan, generation); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("installation %s failed: %v", installID, err))
			result.Success = false
			continue
		}
		result.Installed++
	}
	
	// Set as current generation if successful
	if result.Success {
		if err := r.GenerationManager.SetCurrent(generation.ID); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to set current generation: %v", err))
			result.Success = false
		}
	}
	
	result.Duration = time.Since(startTime)
	return result, nil
}

// planInstallation creates a plan for a single installation
func (r *Reconciler) planInstallation(installID string, config InstallationConfig) (*InstallPlan, error) {
	plan := &InstallPlan{
		Path:       config.Path,
		Operations: []Operation{},
		Symlinks:   make(map[string]string),
	}
	
	// Get current state (if any)
	currentState := make(map[string]InstalledAddon)
	if current, err := r.GenerationManager.GetCurrent(); err == nil {
		if installState, exists := current.Installations[installID]; exists {
			currentState = installState.Addons
		}
	}
	
	// Plan operations for each desired addon
	for _, addonID := range config.Addons {
		addonConfig, exists := r.Config.Addons[addonID]
		if !exists {
			return nil, fmt.Errorf("addon %s not found in configuration", addonID)
		}
		
		// Check compatibility
		if !r.isCompatible(addonConfig, config.Type) {
			return nil, fmt.Errorf("addon %s is not compatible with installation type %s", addonID, config.Type)
		}
		
		targetHash := r.calculateAddonHash(addonConfig)
		currentAddon, exists := currentState[addonID]
		
		var operation Operation
		if !exists {
			// New installation
			operation = Operation{
				Type:         "install",
				Installation: installID,
				Addon:        addonID,
				To:           targetHash,
			}
		} else if currentAddon.Hash != targetHash {
			// Update existing
			hash := currentAddon.Hash
			operation = Operation{
				Type:         "update",
				Installation: installID,
				Addon:        addonID,
				From:         &hash,
				To:           targetHash,
			}
		} else {
			// Already up to date, just ensure symlink
			operation = Operation{
				Type:         "symlink",
				Installation: installID,
				Addon:        addonID,
				To:           targetHash,
			}
		}
		
		plan.Operations = append(plan.Operations, operation)
		plan.Symlinks[addonID] = r.Store.getStorePath(targetHash)
	}
	
	// Plan removal of addons that are no longer wanted
	for addonID := range currentState {
		wanted := false
		for _, wantedID := range config.Addons {
			if wantedID == addonID {
				wanted = true
				break
			}
		}
		
		if !wanted {
			hash := currentState[addonID].Hash
			operation := Operation{
				Type:         "uninstall",
				Installation: installID,
				Addon:        addonID,
				From:         &hash,
			}
			plan.Operations = append(plan.Operations, operation)
		}
	}
	
	return plan, nil
}

// executeDownload downloads and stores an addon
func (r *Reconciler) executeDownload(download DownloadOperation) error {
	// This would implement the actual download logic
	// For now, it's a placeholder that simulates the download
	
	fmt.Printf("📥 Downloading %s from %s...\n", download.AddonID, download.Source.URL)
	
	// In a real implementation, this would:
	// 1. Download from the source URL
	// 2. Extract if it's a ZIP
	// 3. Store in the content store
	// 4. Return the actual hash
	
	return nil
}

// executeInstallPlan executes an installation plan
func (r *Reconciler) executeInstallPlan(installID string, plan InstallPlan, generation *Generation) error {
	fmt.Printf("📦 Executing installation plan for %s...\n", installID)
	
	installState := InstallationState{
		Path:   plan.Path,
		Addons: make(map[string]InstalledAddon),
	}
	
	// Create symlinks
	for addonID, storePath := range plan.Symlinks {
		targetPath := filepath.Join(plan.Path, addonID)
		
		if err := r.Store.Link(filepath.Base(storePath), targetPath); err != nil {
			return fmt.Errorf("failed to create symlink for %s: %v", addonID, err)
		}
		
		// Record in installation state
		installState.Addons[addonID] = InstalledAddon{
			ID:          addonID,
			StorePath:   storePath,
			InstallPath: targetPath,
			InstalledAt: time.Now(),
		}
	}
	
	// Update generation with installation state
	generation.Installations[installID] = installState
	
	return nil
}

// calculateAddonHash calculates a deterministic hash for an addon configuration
func (r *Reconciler) calculateAddonHash(config AddonConfig) string {
	// In a real implementation, this would be based on the actual content
	// For now, create a hash from the configuration
	configBytes, _ := json.Marshal(config)
	return fmt.Sprintf("%x", configBytes)[:16]
}

// isCompatible checks if an addon is compatible with an installation type
func (r *Reconciler) isCompatible(addon AddonConfig, installType string) bool {
	for _, compatibleType := range addon.Compatible {
		if compatibleType == installType {
			return true
		}
	}
	return false
}

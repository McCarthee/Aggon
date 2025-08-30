package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GenerationManager manages system state generations
type GenerationManager struct {
	BasePath string
	Store    *Store
}

// NewGenerationManager creates a new generation manager
func NewGenerationManager(basePath string, store *Store) *GenerationManager {
	return &GenerationManager{
		BasePath: basePath,
		Store:    store,
	}
}

// Initialize creates the generations directory structure
func (gm *GenerationManager) Initialize() error {
	generationsPath := filepath.Join(gm.BasePath, "generations")
	return os.MkdirAll(generationsPath, 0755)
}

// Create creates a new generation from the given configuration
func (gm *GenerationManager) Create(config DeclarativeConfig, description string) (*Generation, error) {
	generations, err := gm.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list generations: %v", err)
	}
	
	// Determine next generation ID
	nextID := 1
	if len(generations) > 0 {
		nextID = generations[len(generations)-1].ID + 1
	}
	
	// Create generation
	generation := &Generation{
		ID:          nextID,
		Timestamp:   time.Now(),
		Config:      config,
		Description: description,
		Installations: make(map[string]InstallationState),
	}
	
	// Calculate state hash
	stateHash, err := gm.calculateStateHash(generation)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate state hash: %v", err)
	}
	generation.StateHash = stateHash
	
	// Save generation
	if err := gm.save(generation); err != nil {
		return nil, fmt.Errorf("failed to save generation: %v", err)
	}
	
	return generation, nil
}

// Get retrieves a specific generation
func (gm *GenerationManager) Get(id int) (*Generation, error) {
	generationPath := gm.getGenerationPath(id)
	
	file, err := os.Open(generationPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open generation file: %v", err)
	}
	defer file.Close()
	
	var generation Generation
	if err := json.NewDecoder(file).Decode(&generation); err != nil {
		return nil, fmt.Errorf("failed to decode generation: %v", err)
	}
	
	return &generation, nil
}

// GetCurrent returns the current active generation
func (gm *GenerationManager) GetCurrent() (*Generation, error) {
	currentPath := filepath.Join(gm.BasePath, "generations", "current")
	
	// Read symlink target
	target, err := os.Readlink(currentPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no current generation set")
		}
		return nil, fmt.Errorf("failed to read current symlink: %v", err)
	}
	
	// Extract generation ID from target
	basename := filepath.Base(target)
	parts := strings.Split(basename, "-")
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid current generation symlink")
	}
	
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid generation ID in symlink: %v", err)
	}
	
	return gm.Get(id)
}

// SetCurrent sets the current active generation
func (gm *GenerationManager) SetCurrent(id int) error {
	generation, err := gm.Get(id)
	if err != nil {
		return fmt.Errorf("generation %d not found: %v", id, err)
	}
	
	generationsDir := filepath.Join(gm.BasePath, "generations")
	currentPath := filepath.Join(generationsDir, "current")
	target := gm.getGenerationDirName(generation)
	
	// Remove existing symlink
	if _, err := os.Lstat(currentPath); err == nil {
		if err := os.Remove(currentPath); err != nil {
			return fmt.Errorf("failed to remove current symlink: %v", err)
		}
	}
	
	// Create new symlink
	if err := os.Symlink(target, currentPath); err != nil {
		return fmt.Errorf("failed to create current symlink: %v", err)
	}
	
	return nil
}

// List returns all generations sorted by ID
func (gm *GenerationManager) List() ([]Generation, error) {
	generationsDir := filepath.Join(gm.BasePath, "generations")
	
	entries, err := os.ReadDir(generationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Generation{}, nil
		}
		return nil, fmt.Errorf("failed to read generations directory: %v", err)
	}
	
	var generations []Generation
	
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "current" {
			// Extract generation ID from directory name
			parts := strings.Split(entry.Name(), "-")
			if len(parts) < 1 {
				continue
			}
			
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
			
			generation, err := gm.Get(id)
			if err != nil {
				continue
			}
			
			generations = append(generations, *generation)
		}
	}
	
	// Sort by ID
	sort.Slice(generations, func(i, j int) bool {
		return generations[i].ID < generations[j].ID
	})
	
	return generations, nil
}

// Delete removes a generation (but not if it's current)
func (gm *GenerationManager) Delete(id int) error {
	current, err := gm.GetCurrent()
	if err == nil && current.ID == id {
		return fmt.Errorf("cannot delete current generation %d", id)
	}
	
	generationDir := filepath.Join(gm.BasePath, "generations", gm.getGenerationDirName(&Generation{ID: id}))
	
	return os.RemoveAll(generationDir)
}

// GarbageCollect removes old generations, keeping the specified number
func (gm *GenerationManager) GarbageCollect(keep int) error {
	generations, err := gm.List()
	if err != nil {
		return fmt.Errorf("failed to list generations: %v", err)
	}
	
	if len(generations) <= keep {
		return nil // Nothing to collect
	}
	
	current, err := gm.GetCurrent()
	if err != nil {
		return fmt.Errorf("failed to get current generation: %v", err)
	}
	
	// Keep the most recent generations and the current one
	toDelete := generations[:len(generations)-keep]
	
	for _, gen := range toDelete {
		if gen.ID == current.ID {
			continue // Don't delete current generation
		}
		
		if err := gm.Delete(gen.ID); err != nil {
			return fmt.Errorf("failed to delete generation %d: %v", gen.ID, err)
		}
	}
	
	return nil
}

// save persists a generation to disk
func (gm *GenerationManager) save(generation *Generation) error {
	generationDir := filepath.Join(gm.BasePath, "generations", gm.getGenerationDirName(generation))
	
	// Create generation directory
	if err := os.MkdirAll(generationDir, 0755); err != nil {
		return fmt.Errorf("failed to create generation directory: %v", err)
	}
	
	// Save generation file
	generationPath := filepath.Join(generationDir, "generation.json")
	file, err := os.Create(generationPath)
	if err != nil {
		return fmt.Errorf("failed to create generation file: %v", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(generation); err != nil {
		return fmt.Errorf("failed to encode generation: %v", err)
	}
	
	return nil
}

// getGenerationPath returns the path to a generation's JSON file
func (gm *GenerationManager) getGenerationPath(id int) string {
	dirName := gm.getGenerationDirName(&Generation{ID: id})
	return filepath.Join(gm.BasePath, "generations", dirName, "generation.json")
}

// getGenerationDirName returns the directory name for a generation
func (gm *GenerationManager) getGenerationDirName(generation *Generation) string {
	timestamp := generation.Timestamp.Format("2006-01-02-150405")
	return fmt.Sprintf("%d-%s", generation.ID, timestamp)
}

// calculateStateHash calculates a hash representing the generation's state
func (gm *GenerationManager) calculateStateHash(generation *Generation) (string, error) {
	// Create a deterministic representation of the generation state
	configBytes, err := json.Marshal(generation.Config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %v", err)
	}
	
	installationBytes, err := json.Marshal(generation.Installations)
	if err != nil {
		return "", fmt.Errorf("failed to marshal installations: %v", err)
	}
	
	// Combine and hash
	combined := append(configBytes, installationBytes...)
	hash := fmt.Sprintf("%x", combined)
	
	return hash[:16], nil // Use first 16 chars for brevity
}

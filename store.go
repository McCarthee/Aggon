package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Store manages the content-addressed addon storage
type Store struct {
	BasePath    string
	StorePath   string
	MetadataPath string
}

// NewStore creates a new content store
func NewStore(basePath string) *Store {
	storePath := filepath.Join(basePath, "store")
	metadataPath := filepath.Join(basePath, "metadata")
	
	return &Store{
		BasePath:     basePath,
		StorePath:    storePath,
		MetadataPath: metadataPath,
	}
}

// Initialize creates the store directory structure
func (s *Store) Initialize() error {
	dirs := []string{
		s.BasePath,
		s.StorePath,
		s.MetadataPath,
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	
	return nil
}

// Add stores content in the store and returns the hash
func (s *Store) Add(content io.Reader, metadata StoreEntry) (string, error) {
	// Calculate hash while reading content
	hasher := sha256.New()
	tempFile, err := os.CreateTemp("", "aggon-store-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	// Write to temp file and calculate hash simultaneously
	writer := io.MultiWriter(tempFile, hasher)
	if _, err := io.Copy(writer, content); err != nil {
		return "", fmt.Errorf("failed to write content: %v", err)
	}
	
	hash := hex.EncodeToString(hasher.Sum(nil))
	storePath := s.getStorePath(hash)
	
	// Check if already exists
	if _, err := os.Stat(storePath); err == nil {
		// Already exists, just update metadata
		metadata.AccessedAt = time.Now()
		return hash, s.saveMetadata(hash, metadata)
	}
	
	// Create store directory for this hash
	storeDir := filepath.Dir(storePath)
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create store directory: %v", err)
	}
	
	// Move temp file to store
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %v", err)
	}
	
	if err := os.Rename(tempFile.Name(), storePath); err != nil {
		return "", fmt.Errorf("failed to move to store: %v", err)
	}
	
	// Save metadata
	metadata.Hash = hash
	metadata.StorePath = storePath
	metadata.CreatedAt = time.Now()
	metadata.AccessedAt = time.Now()
	
	return hash, s.saveMetadata(hash, metadata)
}

// Get retrieves content from the store by hash
func (s *Store) Get(hash string) (io.ReadCloser, error) {
	storePath := s.getStorePath(hash)
	
	file, err := os.Open(storePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open store file: %v", err)
	}
	
	// Update access time
	if metadata, err := s.GetMetadata(hash); err == nil {
		metadata.AccessedAt = time.Now()
		s.saveMetadata(hash, metadata)
	}
	
	return file, nil
}

// GetMetadata retrieves metadata for a store entry
func (s *Store) GetMetadata(hash string) (StoreEntry, error) {
	metadataPath := s.getMetadataPath(hash)
	
	file, err := os.Open(metadataPath)
	if err != nil {
		return StoreEntry{}, fmt.Errorf("failed to open metadata: %v", err)
	}
	defer file.Close()
	
	var metadata StoreEntry
	if err := json.NewDecoder(file).Decode(&metadata); err != nil {
		return StoreEntry{}, fmt.Errorf("failed to decode metadata: %v", err)
	}
	
	return metadata, nil
}

// Exists checks if content with given hash exists in store
func (s *Store) Exists(hash string) bool {
	storePath := s.getStorePath(hash)
	_, err := os.Stat(storePath)
	return err == nil
}

// Link creates a symlink from target to store content
func (s *Store) Link(hash, targetPath string) error {
	storePath := s.getStorePath(hash)
	
	// Ensure target directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}
	
	// Remove existing target if it exists
	if _, err := os.Lstat(targetPath); err == nil {
		if err := os.Remove(targetPath); err != nil {
			return fmt.Errorf("failed to remove existing target: %v", err)
		}
	}
	
	// Create symlink
	relPath, err := filepath.Rel(targetDir, storePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %v", err)
	}
	
	if err := os.Symlink(relPath, targetPath); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}
	
	return nil
}

// GarbageCollect removes unused store entries
func (s *Store) GarbageCollect(keepGenerations int) error {
	// This would implement garbage collection logic
	// For now, it's a placeholder
	return nil
}

// ListEntries returns all store entries
func (s *Store) ListEntries() ([]StoreEntry, error) {
	var entries []StoreEntry
	
	err := filepath.Walk(s.MetadataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		
		hash := filepath.Base(path[:len(path)-5]) // Remove .json extension
		if metadata, err := s.GetMetadata(hash); err == nil {
			entries = append(entries, metadata)
		}
		
		return nil
	})
	
	return entries, err
}

// getStorePath returns the file path for storing content with given hash
func (s *Store) getStorePath(hash string) string {
	// Use first 2 chars as subdirectory for better filesystem performance
	subdir := hash[:2]
	filename := hash[2:]
	return filepath.Join(s.StorePath, subdir, filename)
}

// getMetadataPath returns the metadata file path for given hash
func (s *Store) getMetadataPath(hash string) string {
	return filepath.Join(s.MetadataPath, hash+".json")
}

// saveMetadata saves metadata for a store entry
func (s *Store) saveMetadata(hash string, metadata StoreEntry) error {
	metadataPath := s.getMetadataPath(hash)
	
	// Ensure metadata directory exists
	metadataDir := filepath.Dir(metadataPath)
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %v", err)
	}
	
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %v", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %v", err)
	}
	
	return nil
}

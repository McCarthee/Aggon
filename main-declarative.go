package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printDeclarativeHelp()
		return
	}

	command := os.Args[1]
	
	switch command {
	case "switch":
		handleSwitch()
	case "plan":
		handlePlan()
	case "rollback":
		handleRollback()
	case "generations":
		handleGenerations()
	case "test":
		handleTest()
	case "init":
		handleInit()
	case "help", "--help", "-h":
		printDeclarativeHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printDeclarativeHelp()
		os.Exit(1)
	}
}

func handleSwitch() {
	configPath := getConfigPath()
	
	// Load configuration
	config, err := loadDeclarativeConfig(configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	
	// Initialize system
	store := NewStore(config.Settings.StorePath)
	if err := store.Initialize(); err != nil {
		fmt.Printf("❌ Failed to initialize store: %v\n", err)
		os.Exit(1)
	}
	
	gm := NewGenerationManager(config.Settings.GenerationsPath, store)
	if err := gm.Initialize(); err != nil {
		fmt.Printf("❌ Failed to initialize generation manager: %v\n", err)
		os.Exit(1)
	}
	
	// Create reconciler and plan
	reconciler := NewReconciler(*config, store, gm)
	
	fmt.Println("🔄 Planning system changes...")
	plan, err := reconciler.Plan()
	if err != nil {
		fmt.Printf("❌ Failed to create plan: %v\n", err)
		os.Exit(1)
	}
	
	if len(plan.Operations) == 0 {
		fmt.Println("✅ System is already up to date!")
		return
	}
	
	// Show plan summary
	fmt.Printf("📋 Plan Summary:\n")
	fmt.Printf("   Downloads: %d\n", len(plan.Downloads))
	fmt.Printf("   Operations: %d\n", len(plan.Operations))
	fmt.Printf("   Installations: %d\n", len(plan.Installations))
	fmt.Println()
	
	// Apply changes
	fmt.Println("⚡ Applying changes...")
	result, err := reconciler.Apply(plan)
	if err != nil {
		fmt.Printf("❌ Failed to apply changes: %v\n", err)
		os.Exit(1)
	}
	
	// Show results
	if result.Success {
		fmt.Printf("✅ Successfully applied configuration!\n")
		fmt.Printf("   Generation: %d\n", result.Generation)
		fmt.Printf("   Downloaded: %d addons\n", result.Downloaded)
		fmt.Printf("   Installed: %d installations\n", result.Installed)
		fmt.Printf("   Duration: %v\n", result.Duration)
	} else {
		fmt.Printf("❌ Configuration apply failed!\n")
		for _, err := range result.Errors {
			fmt.Printf("   Error: %v\n", err)
		}
		os.Exit(1)
	}
}

func handlePlan() {
	configPath := getConfigPath()
	
	config, err := loadDeclarativeConfig(configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	
	store := NewStore(config.Settings.StorePath)
	gm := NewGenerationManager(config.Settings.GenerationsPath, store)
	reconciler := NewReconciler(*config, store, gm)
	
	fmt.Println("🔄 Creating execution plan...")
	plan, err := reconciler.Plan()
	if err != nil {
		fmt.Printf("❌ Failed to create plan: %v\n", err)
		os.Exit(1)
	}
	
	// Display plan details
	fmt.Printf("📋 Execution Plan\n")
	fmt.Printf("=================\n\n")
	
	if len(plan.Downloads) > 0 {
		fmt.Printf("📥 Downloads (%d):\n", len(plan.Downloads))
		for _, download := range plan.Downloads {
			fmt.Printf("   - %s from %s\n", download.AddonID, download.Source.URL)
		}
		fmt.Println()
	}
	
	if len(plan.Operations) > 0 {
		fmt.Printf("⚙️  Operations (%d):\n", len(plan.Operations))
		for _, op := range plan.Operations {
			switch op.Type {
			case "install":
				fmt.Printf("   + Install %s in %s\n", op.Addon, op.Installation)
			case "update":
				fmt.Printf("   ~ Update %s in %s\n", op.Addon, op.Installation)
			case "uninstall":
				fmt.Printf("   - Uninstall %s from %s\n", op.Addon, op.Installation)
			case "symlink":
				fmt.Printf("   → Link %s in %s\n", op.Addon, op.Installation)
			}
		}
		fmt.Println()
	}
	
	if len(plan.Operations) == 0 {
		fmt.Println("✅ No changes needed - system is up to date!")
	}
}

func handleRollback() {
	// Implementation for rollback command
	fmt.Println("🔄 Rollback functionality - Coming soon!")
}

func handleGenerations() {
	// Implementation for generations command
	fmt.Println("📋 Generations management - Coming soon!")
}

func handleTest() {
	// Implementation for test command
	fmt.Println("🧪 Test configuration - Coming soon!")
}

func handleInit() {
	fmt.Println("🎉 Initializing declarative AGGON...")
	
	// Create default configuration
	config := DeclarativeConfig{
		Schema: "aggon/v2",
		Metadata: ConfigMetadata{
			Name:        "my-wow-setup",
			Version:     "1.0.0",
			Description: "Declarative WoW addon configuration",
		},
		Installations: make(map[string]InstallationConfig),
		Addons:        make(map[string]AddonConfig),
		Profiles:      make(map[string]ProfileConfig),
		Settings: SystemSettings{
			AutoUpdate:        false,
			BackupGenerations: 10,
			ParallelDownloads: 3,
			VerifyHashes:      true,
			StorePath:         ".aggon/store",
			GenerationsPath:   ".aggon/generations",
		},
	}
	
	// Save to file
	configPath := "aggon-declarative.json"
	if err := saveDeclarativeConfig(configPath, &config); err != nil {
		fmt.Printf("❌ Failed to create configuration: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✅ Created %s\n", configPath)
	fmt.Println("   Edit this file to define your addon setup, then run:")
	fmt.Println("   aggon-declarative switch")
}

func getConfigPath() string {
	if len(os.Args) > 2 && os.Args[2] == "--config" && len(os.Args) > 3 {
		return os.Args[3]
	}
	return "aggon-declarative.json"
}

func loadDeclarativeConfig(path string) (*DeclarativeConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()
	
	var config DeclarativeConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}
	
	return &config, nil
}

func saveDeclarativeConfig(path string, config *DeclarativeConfig) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}
	
	return nil
}

func printDeclarativeHelp() {
	fmt.Println("🏺 AGGON v2 - Declarative WoW Addon Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  aggon-declarative <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init                     Create a new declarative configuration")
	fmt.Println("  plan                     Show what changes would be made")
	fmt.Println("  switch                   Apply the configuration")
	fmt.Println("  test                     Test configuration without applying")
	fmt.Println("  rollback [generation]    Rollback to a previous generation")
	fmt.Println("  generations list         List all generations")
	fmt.Println("  help                     Show this help")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --config <file>         Use specific configuration file")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  aggon-declarative init")
	fmt.Println("  aggon-declarative plan")
	fmt.Println("  aggon-declarative switch")
	fmt.Println("  aggon-declarative rollback 5")
}

# 🏺 AGGON - Declarative WoW Addon Manager

A revolutionary declarative addon manager for World of Warcraft, inspired by NixOS principles.

## 🎯 Overview

AGGON provides a declarative approach to World of Warcraft addon management with atomic operations, generation-based rollbacks, and reproducible configurations.

## 📋 Table of Contents

-   [Quick Start](#-quick-start)
-   [Core Concepts](#-core-concepts)
-   [Configuration Guide](#-configuration-guide)
-   [Advanced Usage](#-advanced-usage)
-   [Troubleshooting](#-troubleshooting)
-   [Contributing](#-contributing)

---

## 🚀 Quick Start

### Getting Started

First, compile the project:
```bash
go build -o aggon.exe main.go types.go store.go generations.go reconcile.go
```

Then initialize:
```bash
# Initialize declarative configuration
./aggon.exe init

# Edit the generated configuration
# (See Configuration Guide below)

# Preview changes
./aggon.exe plan

# Apply configuration
./aggon.exe switch
```

---



## 🧬 Core Concepts

### Philosophy

**"Everything should be declarable, reproducible, and rollback-able"**

### Core Concepts

#### 1. **Declarative Configuration**

Instead of imperatively managing addons, you declare the desired state in a configuration file.

#### 2. **Content-Addressed Storage**

Addons are stored immutably with SHA256 hashes, ensuring integrity and deduplication.

#### 3. **Atomic Operations**

Changes are applied atomically - either everything succeeds or nothing changes.

#### 4. **Generation Management**

Every configuration change creates a new "generation" that can be instantly rolled back.

#### 5. **Symlink Farms**

Your addon installations are symlinks pointing to the immutable store.

### Key Benefits

-   🔒 **Immutable**: Addons never get corrupted during updates
-   ⚛️ **Atomic**: All-or-nothing updates prevent broken states
-   🔄 **Reproducible**: Same configuration = identical results
-   ⏪ **Rollback**: Instant return to any previous state
-   🎯 **Multi-Environment**: Profiles for different setups
-   🧹 **Deduplication**: Same addon stored once

### Directory Structure

```
.aggon/
├── store/                    # Content-addressed addon storage
│   ├── ab/c123...           # Addon stored by hash
│   └── de/f456...           # Another addon version
├── generations/              # System state snapshots
│   ├── 1-2024-01-15/        # Generation 1
│   ├── 2-2024-01-16/        # Generation 2
│   └── current -> 2-2024-01-16/
└── metadata/                 # Store metadata
    ├── abc123.json          # Metadata for abc123 hash
    └── def456.json          # Metadata for def456 hash
```

### Commands

```bash
# Initialize new setup
aggon init

# Show planned changes
aggon plan

# Apply configuration
aggon switch

# Test without applying
aggon test

# Rollback to previous generation
aggon rollback

# List all generations
aggon generations list

# Rollback to specific generation
aggon rollback 5

# Get help
aggon help
```

---

## 📝 Configuration Guide

### Complete Configuration Example

```json
{
    "schema": "aggon/v2",
    "metadata": {
        "name": "my-wow-setup",
        "version": "1.0.0",
        "description": "My complete WoW addon setup"
    },
    "installations": {
        "retail": {
            "type": "retail",
            "path": "C:/Games/WoW/_retail_/Interface/AddOns",
            "enabled": true,
            "addons": ["elvui", "weakauras", "details"]
        },
        "classic": {
            "type": "classic",
            "path": "C:/Games/WoW/_classic_/Interface/AddOns",
            "enabled": false,
            "addons": ["questie", "details"]
        },
        "ascension": {
            "type": "custom",
            "path": "C:/Games/Ascension/Interface/Addons",
            "enabled": true,
            "addons": ["elvui-epoch", "pfquest"]
        }
    },
    "addons": {
        "elvui": {
            "source": {
                "type": "github",
                "url": "https://github.com/ElvUI-WotLK/ElvUI",
                "ref": "main"
            },
            "version": "latest",
            "hash": null,
            "ignore": ["README.md", ".gitattributes"],
            "compatible": ["retail", "classic"],
            "folder": null
        },
        "elvui-epoch": {
            "source": {
                "type": "github",
                "url": "https://github.com/Bennylavaa/ElvUI-Epoch",
                "ref": "main"
            },
            "version": "latest",
            "hash": null,
            "ignore": ["README.md"],
            "compatible": ["custom"],
            "folder": "ElvUI"
        },
        "weakauras": {
            "source": {
                "type": "github",
                "url": "https://github.com/WeakAuras/WeakAuras2",
                "ref": "main"
            },
            "version": "latest",
            "hash": null,
            "compatible": ["retail"],
            "folder": null
        },
        "pfquest": {
            "source": {
                "type": "github",
                "url": "https://github.com/shagu/pfQuest",
                "ref": "refs/tags/v3.5"
            },
            "version": "3.5",
            "hash": "sha256:abc123...",
            "ignore": ["LICENSE", "README.md"],
            "compatible": ["custom", "classic"],
            "folder": "pfQuest"
        }
    },
    "profiles": {
        "default": {
            "installations": ["retail"],
            "description": "Standard retail setup"
        },
        "development": {
            "installations": ["retail"],
            "addons": {
                "elvui": {
                    "ref": "development",
                    "version": "beta"
                }
            },
            "description": "Development setup with beta addons"
        },
        "full": {
            "installations": ["retail", "classic", "ascension"],
            "description": "All installations active"
        }
    },
    "settings": {
        "auto_update": false,
        "backup_generations": 10,
        "parallel_downloads": 3,
        "verify_hashes": true,
        "store_path": ".aggon/store",
        "generations_path": ".aggon/generations"
    }
}
```

### Configuration Sections

#### Metadata

Basic information about your configuration.

#### Installations

Define your WoW installations:

-   **type**: `retail`, `classic`, or `custom`
-   **path**: Path to the AddOns directory
-   **enabled**: Whether this installation is active
-   **addons**: List of addon IDs to install

#### Addons

Define individual addons:

-   **source**: Where to get the addon
    -   **type**: Currently only `github` supported
    -   **url**: GitHub repository URL
    -   **ref**: Branch, tag, or commit (`main`, `refs/tags/v1.0`, etc.)
-   **version**: Version constraint or `latest`
-   **hash**: SHA256 hash for reproducibility (optional)
-   **ignore**: Files to exclude when extracting
-   **compatible**: Which installation types can use this addon
-   **folder**: Custom folder name (optional)

#### Profiles

Different configurations for different scenarios:

-   **installations**: Which installations to include
-   **addons**: Addon-specific overrides
-   **description**: Human-readable description

#### Settings

Global system settings:

-   **auto_update**: Automatically check for updates
-   **backup_generations**: Number of generations to keep
-   **parallel_downloads**: Concurrent download limit
-   **verify_hashes**: Verify addon integrity
-   **store_path**: Where to store addons
-   **generations_path**: Where to store generations

---

## 🛠️ Advanced Usage

### Working with Profiles

```bash
# Switch to development profile
aggon switch --profile development

# Apply specific profile
aggon switch --profile full
```

### Generation Management

```bash
# List all generations
aggon generations list

# View specific generation
aggon generations show 5

# Rollback to generation 3
aggon rollback 3

# Clean up old generations (keep last 5)
aggon gc --keep 5
```

### Hash Pinning for Reproducibility

```json
{
    "addons": {
        "elvui": {
            "version": "13.45",
            "hash": "sha256:a1b2c3d4e5f6..."
        }
    }
}
```

### Custom Installation Paths

```json
{
    "installations": {
        "turtle-wow": {
            "type": "custom",
            "path": "D:/Games/TurtleWoW/Interface/AddOns",
            "enabled": true,
            "addons": ["pfquest", "questie"]
        }
    }
}
```

---

## 🚨 Troubleshooting

### Common Issues

#### "Failed to create symlink"

-   **Cause**: Permission issues or existing files
-   **Solution**: Run as administrator or remove conflicting files

#### "Hash mismatch"

-   **Cause**: Addon content changed or corrupted
-   **Solution**: Clear store and re-download: `aggon gc --force`

#### "Generation not found"

-   **Cause**: Generation was deleted or corrupted
-   **Solution**: Check available generations: `aggon generations list`

#### "Configuration parse error"

-   **Cause**: Invalid JSON syntax
-   **Solution**: Validate JSON syntax, check brackets and commas

### Recovery Procedures

#### Broken Installation

```bash
# Rollback to last working generation
aggon rollback

# Or rollback to specific generation
aggon rollback 5
```

#### Corrupted Store

```bash
# Clear store and rebuild
rm -rf .aggon/store
aggon switch
```

#### Reset Everything

```bash
# Nuclear option - complete reset
rm -rf .aggon
aggon init
# Reconfigure and switch
```

---

## 🔍 Debugging

### Verbose Output

```bash
# Enable verbose logging
aggon switch --verbose

# Debug mode
aggon plan --debug
```

### Log Files

-   Generation logs: `.aggon/generations/*/logs/`
-   Store operations: `.aggon/store.log`
-   Error details: `.aggon/errors.log`

---

## 🎯 Best Practices

### Configuration Management

-   ✅ Use version control for your `aggon.json`
-   ✅ Pin addon versions for stability
-   ✅ Use profiles for different environments
-   ✅ Document changes in commit messages

### Maintenance

-   ✅ Regularly clean old generations
-   ✅ Review and update addon versions
-   ✅ Test configurations before applying
-   ✅ Keep backups of working configurations

### Security

-   ✅ Verify addon sources (only trusted GitHub repos)
-   ✅ Enable hash verification
-   ✅ Review changes before applying
-   ✅ Use specific versions instead of "latest" for important addons

---

## 🤝 Contributing

### Development Setup

```bash
git clone <repo>
cd aggon
go build -o aggon.exe main.go types.go store.go generations.go reconcile.go
```

### Running Tests

```bash
go test ./...
```

### Adding Features

1. Create feature branch from `main`
2. Implement changes
3. Add tests
4. Submit pull request

---

## 📜 License

MIT License - See LICENSE file for details.

---

## 🙏 Acknowledgments

-   Inspired by [NixOS](https://nixos.org/) declarative system design
-   Built for the World of Warcraft addon community
-   Special thanks to addon developers maintaining GitHub repositories

---

_Happy addon management! 🏺_

# 🏺 AGGON - World of Warcraft Addon Manager

A modern, declarative addon manager for World of Warcraft, inspired by NixOS principles.

## 🎯 Overview

AGGON comes in two versions:

-   **AGGON v1** (Classic): Traditional imperative addon management
-   **AGGON v2** (Declarative): Revolutionary declarative system with atomic operations and rollbacks

## 📋 Table of Contents

-   [Quick Start](#-quick-start)
-   [AGGON v1 (Classic)](#-aggon-v1-classic)
-   [AGGON v2 (Declarative)](#-aggon-v2-declarative)
-   [Migration Guide](#-migration-guide)
-   [Configuration Reference](#-configuration-reference)
-   [Troubleshooting](#-troubleshooting)
-   [Contributing](#-contributing)

---

## 🚀 Quick Start

### For New Users (Recommended: v2 Declarative)

```bash
# Initialize declarative configuration
./aggon-declarative.exe init

# Edit the generated configuration
# (See Configuration Guide below)

# Preview changes
./aggon-declarative.exe plan

# Apply configuration
./aggon-declarative.exe switch
```

### For Existing Users (v1 Classic)

```bash
# Run interactive menu
./aggon.exe

# Or use command line
./aggon.exe add addon
./aggon.exe add path
```

---

## 🏛️ AGGON v1 (Classic)

The traditional imperative addon manager with interactive menus.

### Features

-   ✅ Interactive menu system
-   ✅ GitHub integration with caching
-   ✅ Basic backup system
-   ✅ Multi-installation support
-   ✅ Clean progress display

### Usage

```bash
# Interactive mode
./aggon.exe

# Command line mode
./aggon.exe add addon          # Add new addon
./aggon.exe add path           # Add installation path
./aggon.exe format-config      # Format configuration
./aggon.exe --help             # Show help
```

### Configuration (config.json)

```json
[
    {
        "name": "Ascension",
        "path": "C:/Games/Ascension/Interface/Addons",
        "addons": [
            {
                "name": "ElvUI-Epoch",
                "url": "https://github.com/Bennylavaa/ElvUI-Epoch",
                "ignore": ["README.md", ".gitattributes"]
            }
        ]
    }
]
```

---

## 🧬 AGGON v2 (Declarative)

A revolutionary declarative system inspired by NixOS principles.

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
aggon-declarative init

# Show planned changes
aggon-declarative plan

# Apply configuration
aggon-declarative switch

# Test without applying
aggon-declarative test

# Rollback to previous generation
aggon-declarative rollback

# List all generations
aggon-declarative generations list

# Rollback to specific generation
aggon-declarative rollback 5

# Get help
aggon-declarative help
```

---

## 📝 Configuration Guide (v2 Declarative)

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

## 🔄 Migration Guide

### From AGGON v1 to v2

1. **Initialize v2 configuration:**

    ```bash
    ./aggon-declarative.exe init
    ```

2. **Convert your v1 config manually:**

    - Copy addon definitions from `config.json`
    - Restructure according to v2 schema
    - Define installations and profiles

3. **Test the configuration:**

    ```bash
    ./aggon-declarative.exe plan
    ```

4. **Apply when ready:**
    ```bash
    ./aggon-declarative.exe switch
    ```

### From Manual Management

1. **Document current setup:**

    - List all installed addons
    - Note their sources (GitHub URLs)
    - Document any custom configurations

2. **Create declarative config:**

    - Use the configuration guide above
    - Define all your addons and installations

3. **Test before applying:**
    ```bash
    ./aggon-declarative.exe plan
    ```

---

## 🛠️ Advanced Usage

### Working with Profiles

```bash
# Switch to development profile
aggon-declarative switch --profile development

# Apply specific profile
aggon-declarative switch --profile full
```

### Generation Management

```bash
# List all generations
aggon-declarative generations list

# View specific generation
aggon-declarative generations show 5

# Rollback to generation 3
aggon-declarative rollback 3

# Clean up old generations (keep last 5)
aggon-declarative gc --keep 5
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
-   **Solution**: Clear store and re-download: `aggon-declarative gc --force`

#### "Generation not found"

-   **Cause**: Generation was deleted or corrupted
-   **Solution**: Check available generations: `aggon-declarative generations list`

#### "Configuration parse error"

-   **Cause**: Invalid JSON syntax
-   **Solution**: Validate JSON syntax, check brackets and commas

### Recovery Procedures

#### Broken Installation

```bash
# Rollback to last working generation
aggon-declarative rollback

# Or rollback to specific generation
aggon-declarative rollback 5
```

#### Corrupted Store

```bash
# Clear store and rebuild
rm -rf .aggon/store
aggon-declarative switch
```

#### Reset Everything

```bash
# Nuclear option - complete reset
rm -rf .aggon
aggon-declarative init
# Reconfigure and switch
```

---

## 🔍 Debugging

### Verbose Output

```bash
# Enable verbose logging
aggon-declarative switch --verbose

# Debug mode
aggon-declarative plan --debug
```

### Log Files

-   Generation logs: `.aggon/generations/*/logs/`
-   Store operations: `.aggon/store.log`
-   Error details: `.aggon/errors.log`

---

## 📊 Comparison: v1 vs v2

| Feature               | AGGON v1              | AGGON v2             |
| --------------------- | --------------------- | -------------------- |
| **Configuration**     | Imperative            | Declarative          |
| **Updates**           | In-place modification | Immutable + symlinks |
| **Rollback**          | Basic backups         | Instant generations  |
| **Reproducibility**   | Limited               | Complete             |
| **Multi-environment** | Manual                | Profiles             |
| **State Management**  | Mutable               | Immutable            |
| **Recovery**          | Manual restore        | Atomic rollback      |
| **Deduplication**     | None                  | Automatic            |
| **Integrity**         | Basic                 | SHA256 verified      |

---

## 🎯 Best Practices

### Configuration Management

-   ✅ Use version control for your `aggon-declarative.json`
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
git checkout declarative-v2
go build -o aggon-declarative.exe main-declarative.go types.go store.go generations.go reconcile.go
```

### Running Tests

```bash
go test ./...
```

### Adding Features

1. Create feature branch from `declarative-v2`
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

# 📝 AGGON Configuration Reference

Complete reference for declarative configuration in AGGON.

## 🏗️ Schema Overview

```json
{
  "schema": "aggon/v2",
  "metadata": { ... },
  "installations": { ... },
  "addons": { ... },
  "profiles": { ... },
  "settings": { ... }
}
```

---

## 📋 Configuration Schema

### Root Level

| Field           | Type   | Required | Description                         |
| --------------- | ------ | -------- | ----------------------------------- |
| `schema`        | string | ✅       | Schema version (must be "aggon/v2") |
| `metadata`      | object | ✅       | Configuration metadata              |
| `installations` | object | ✅       | WoW installation definitions        |
| `addons`        | object | ✅       | Addon definitions                   |
| `profiles`      | object | ❌       | Profile definitions                 |
| `settings`      | object | ✅       | Global settings                     |

---

## 📊 Metadata Section

Describes your configuration.

```json
{
    "metadata": {
        "name": "my-wow-setup",
        "version": "1.0.0",
        "description": "My complete WoW addon configuration"
    }
}
```

| Field         | Type   | Required | Description                                |
| ------------- | ------ | -------- | ------------------------------------------ |
| `name`        | string | ✅       | Human-readable configuration name          |
| `version`     | string | ✅       | Configuration version (semver recommended) |
| `description` | string | ❌       | Description of this configuration          |

---

## 🎮 Installations Section

Defines your WoW installations.

```json
{
    "installations": {
        "retail": {
            "type": "retail",
            "path": "C:/Games/WoW/_retail_/Interface/AddOns",
            "enabled": true,
            "addons": ["elvui", "weakauras"]
        },
        "classic": {
            "type": "classic",
            "path": "C:/Games/WoW/_classic_/Interface/AddOns",
            "enabled": false,
            "addons": ["questie"]
        }
    }
}
```

### Installation Object

| Field     | Type    | Required | Description                                      |
| --------- | ------- | -------- | ------------------------------------------------ |
| `type`    | string  | ✅       | Installation type: `retail`, `classic`, `custom` |
| `path`    | string  | ✅       | Absolute path to AddOns directory                |
| `enabled` | boolean | ✅       | Whether this installation is active              |
| `addons`  | array   | ✅       | List of addon IDs to install                     |

### Installation Types

-   **`retail`**: Modern WoW (Dragonflight, etc.)
-   **`classic`**: WoW Classic (Vanilla, TBC, WotLK)
-   **`custom`**: Private servers, custom clients

### Path Examples

```json
{
    "retail": {
        "path": "C:/Program Files (x86)/World of Warcraft/_retail_/Interface/AddOns"
    },
    "classic": {
        "path": "C:/Program Files (x86)/World of Warcraft/_classic_/Interface/AddOns"
    },
    "ascension": {
        "path": "C:/Games/Ascension Launcher/resources/epoch_live/Interface/Addons"
    },
    "turtle": {
        "path": "D:/Games/Turtle WoW/Interface/AddOns"
    }
}
```

---

## 🧩 Addons Section

Defines individual addons and their sources.

```json
{
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
        }
    }
}
```

### Addon Object

| Field        | Type        | Required | Description                       |
| ------------ | ----------- | -------- | --------------------------------- |
| `source`     | object      | ✅       | Source configuration              |
| `version`    | string      | ✅       | Version constraint                |
| `hash`       | string/null | ❌       | SHA256 hash for reproducibility   |
| `ignore`     | array       | ❌       | Files to ignore during extraction |
| `compatible` | array       | ✅       | Compatible installation types     |
| `folder`     | string/null | ❌       | Custom folder name                |

### Source Object

| Field  | Type   | Required | Description                           |
| ------ | ------ | -------- | ------------------------------------- |
| `type` | string | ✅       | Source type (currently only "github") |
| `url`  | string | ✅       | Repository URL                        |
| `ref`  | string | ✅       | Branch, tag, or commit reference      |

### Source References

```json
{
  "ref": "main"                    // Branch
  "ref": "development"             // Another branch
  "ref": "refs/tags/v1.0.0"       // Specific tag
  "ref": "refs/tags/1.5"          // Tag without 'v' prefix
  "ref": "abc123def456"            // Specific commit hash
}
```

### Version Constraints

```json
{
  "version": "latest"              // Always use latest available
  "version": "1.0.0"              // Specific version
  "version": "^1.0"               // Compatible with 1.x
  "version": "~1.0.0"             // Compatible with 1.0.x
  "version": ">=1.0.0"            // At least version 1.0.0
}
```

### Hash Pinning

For reproducible builds, specify exact content hashes:

```json
{
    "hash": "sha256:a1b2c3d4e5f6789..."
}
```

When `hash` is specified:

-   AGGON verifies downloaded content matches exactly
-   Enables reproducible installations across machines
-   Prevents supply-chain attacks

### Ignore Patterns

Files to exclude during extraction:

```json
{
    "ignore": [
        "README.md", // Exact filename
        ".gitattributes", // Hidden files
        ".github", // Directories
        "LICENSE", // License files
        "docs/", // Documentation folders
        "*.txt", // Glob patterns (future)
        "**/*.md" // Recursive patterns (future)
    ]
}
```

### Custom Folders

Override the default folder name:

```json
{
  "folder": "MyCustomAddonName"    // Use custom name
  "folder": null                   // Use default (from ZIP)
}
```

### Compatibility

Specify which installation types can use this addon:

```json
{
  "compatible": ["retail"]                    // Retail only
  "compatible": ["classic"]                   // Classic only
  "compatible": ["custom"]                    // Private servers only
  "compatible": ["retail", "classic"]         // Both retail and classic
  "compatible": ["retail", "classic", "custom"] // All types
}
```

---

## 🎭 Profiles Section

Profiles allow different configurations for different scenarios.

```json
{
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
            "description": "Development with beta addons"
        },
        "full": {
            "installations": ["retail", "classic", "ascension"],
            "description": "All installations"
        }
    }
}
```

### Profile Object

| Field           | Type   | Required | Description                    |
| --------------- | ------ | -------- | ------------------------------ |
| `installations` | array  | ✅       | Which installations to include |
| `addons`        | object | ❌       | Addon-specific overrides       |
| `description`   | string | ❌       | Human-readable description     |

### Addon Overrides

Profiles can override specific addon settings:

```json
{
    "addons": {
        "elvui": {
            "ref": "development", // Use dev branch
            "version": "beta", // Override version
            "enabled": false // Disable this addon
        },
        "weakauras": {
            "ref": "refs/tags/v5.0.0" // Pin to specific version
        }
    }
}
```

### Profile Usage

```bash
# Switch to specific profile
aggon switch --profile development

# Plan with profile
aggon plan --profile full
```

---

## ⚙️ Settings Section

Global system configuration.

```json
{
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

### Settings Fields

| Field                | Type    | Default              | Description                           |
| -------------------- | ------- | -------------------- | ------------------------------------- |
| `auto_update`        | boolean | false                | Automatically check for addon updates |
| `backup_generations` | integer | 10                   | Number of generations to keep         |
| `parallel_downloads` | integer | 3                    | Concurrent download limit             |
| `verify_hashes`      | boolean | true                 | Verify addon content integrity        |
| `store_path`         | string  | ".aggon/store"       | Path to content store                 |
| `generations_path`   | string  | ".aggon/generations" | Path to generations                   |

### Path Configuration

Paths can be relative or absolute:

```json
{
    "store_path": ".aggon/store", // Relative to config file
    "store_path": "C:/AGGON/Store", // Absolute path
    "store_path": "~/aggon-store", // Home directory (Unix)
    "store_path": "%USERPROFILE%/aggon-store" // Home directory (Windows)
}
```

---

## 📝 Complete Configuration Examples

### Minimal Configuration

```json
{
    "schema": "aggon/v2",
    "metadata": {
        "name": "minimal-setup",
        "version": "1.0.0"
    },
    "installations": {
        "retail": {
            "type": "retail",
            "path": "C:/Games/WoW/_retail_/Interface/AddOns",
            "enabled": true,
            "addons": ["elvui"]
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
            "compatible": ["retail"]
        }
    },
    "profiles": {},
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

### Multi-Installation Setup

```json
{
    "schema": "aggon/v2",
    "metadata": {
        "name": "multi-wow-setup",
        "version": "2.1.0",
        "description": "Retail, Classic, and Ascension setup"
    },
    "installations": {
        "retail": {
            "type": "retail",
            "path": "C:/Games/WoW/_retail_/Interface/AddOns",
            "enabled": true,
            "addons": ["elvui", "weakauras", "details", "bigwigs"]
        },
        "classic": {
            "type": "classic",
            "path": "C:/Games/WoW/_classic_/Interface/AddOns",
            "enabled": true,
            "addons": ["elvui-classic", "questie", "details"]
        },
        "ascension": {
            "type": "custom",
            "path": "C:/Games/Ascension/Interface/Addons",
            "enabled": true,
            "addons": ["elvui-epoch", "pfquest-epoch"]
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
            "compatible": ["retail"],
            "ignore": ["README.md", ".gitattributes"]
        },
        "elvui-classic": {
            "source": {
                "type": "github",
                "url": "https://github.com/ElvUI-WotLK/ElvUI",
                "ref": "classic"
            },
            "version": "latest",
            "compatible": ["classic"],
            "folder": "ElvUI"
        },
        "elvui-epoch": {
            "source": {
                "type": "github",
                "url": "https://github.com/Bennylavaa/ElvUI-Epoch",
                "ref": "main"
            },
            "version": "latest",
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
            "compatible": ["retail"]
        },
        "questie": {
            "source": {
                "type": "github",
                "url": "https://github.com/AeroScripts/QuestieDev",
                "ref": "master"
            },
            "version": "latest",
            "compatible": ["classic"]
        },
        "details": {
            "source": {
                "type": "github",
                "url": "https://github.com/Tercioo/Details-Damage-Meter",
                "ref": "master"
            },
            "version": "latest",
            "compatible": ["retail", "classic"]
        },
        "bigwigs": {
            "source": {
                "type": "github",
                "url": "https://github.com/BigWigsMods/BigWigs",
                "ref": "main"
            },
            "version": "latest",
            "compatible": ["retail"],
            "hash": "sha256:abc123...",
            "ignore": ["README.md", "LICENSE", ".github"]
        },
        "pfquest-epoch": {
            "source": {
                "type": "github",
                "url": "https://github.com/Bennylavaa/pfQuest-epoch",
                "ref": "refs/tags/1.15"
            },
            "version": "1.15",
            "compatible": ["custom"],
            "folder": "pfQuest",
            "hash": "sha256:def456..."
        }
    },
    "profiles": {
        "retail-only": {
            "installations": ["retail"],
            "description": "Retail WoW only"
        },
        "classic-only": {
            "installations": ["classic"],
            "description": "Classic WoW only"
        },
        "ascension-only": {
            "installations": ["ascension"],
            "description": "Ascension only"
        },
        "development": {
            "installations": ["retail"],
            "addons": {
                "elvui": {
                    "ref": "development"
                },
                "weakauras": {
                    "ref": "development"
                }
            },
            "description": "Development versions"
        },
        "stable": {
            "installations": ["retail", "classic"],
            "addons": {
                "elvui": {
                    "ref": "refs/tags/13.45",
                    "hash": "sha256:stable123..."
                }
            },
            "description": "Stable pinned versions"
        }
    },
    "settings": {
        "auto_update": false,
        "backup_generations": 15,
        "parallel_downloads": 5,
        "verify_hashes": true,
        "store_path": "D:/AGGON/Store",
        "generations_path": "D:/AGGON/Generations"
    }
}
```

---

## ✅ Configuration Validation

### Required Fields

AGGON validates your configuration and will report errors for:

-   Missing required fields
-   Invalid schema version
-   Invalid paths
-   Unknown addon IDs in installation addon lists
-   Invalid source URLs
-   Circular dependencies

### Validation Commands

```bash
# Validate configuration without applying
aggon plan

# Test configuration
aggon test

# Verbose validation
aggon plan --verbose
```

### Common Validation Errors

#### "Addon 'xyz' not found"

```json
{
    "installations": {
        "retail": {
            "addons": ["elvui", "nonexistent-addon"] // ❌ Not defined in addons section
        }
    }
}
```

#### "Invalid path"

```json
{
    "installations": {
        "retail": {
            "path": "invalid/path/here" // ❌ Path doesn't exist
        }
    }
}
```

#### "Incompatible addon"

```json
{
    "installations": {
        "retail": {
            "type": "retail",
            "addons": ["classic-only-addon"] // ❌ Addon not compatible with retail
        }
    },
    "addons": {
        "classic-only-addon": {
            "compatible": ["classic"] // Only works with classic
        }
    }
}
```

---

---

## 🎯 Best Practices

### 1. Configuration Management

-   ✅ Store configuration in version control
-   ✅ Use meaningful names and descriptions
-   ✅ Comment complex configurations (when JSON supports it)
-   ✅ Pin important addons with hash verification

### 2. Path Management

-   ✅ Use absolute paths for installations
-   ✅ Verify paths exist before configuration
-   ✅ Use forward slashes even on Windows
-   ✅ Avoid spaces in custom store paths

### 3. Version Management

-   ✅ Use semantic versioning for your config
-   ✅ Pin stable addon versions in production
-   ✅ Use "latest" only for development
-   ✅ Test configurations before deploying

### 4. Security

-   ✅ Only use trusted GitHub repositories
-   ✅ Enable hash verification
-   ✅ Pin versions for critical addons
-   ✅ Review addon sources regularly

### 5. Performance

-   ✅ Set appropriate parallel_downloads
-   ✅ Clean up old generations regularly
-   ✅ Use local store paths for better performance
-   ✅ Group compatible addons efficiently

---

_This completes the comprehensive configuration reference for AGGON v2! 🏺_

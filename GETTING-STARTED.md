# 🚀 Getting Started with AGGON

A step-by-step guide to set up and use the declarative AGGON addon manager.

## 📋 Prerequisites

-   Windows 10/11 (Linux/macOS support planned)
-   World of Warcraft installation(s)
-   Git (optional, for version control)
-   Text editor or IDE

## 🎯 Quick Start (5 Minutes)

### Step 1: Download AGGON

```bash
# Download the latest release
# Or compile from source:
git clone <repo>
cd aggon
git checkout main
go build -o aggon.exe main.go types.go store.go generations.go reconcile.go
```

### Step 2: Initialize Configuration

```bash
# Create initial configuration
./aggon.exe init
```

This creates `config.json` with basic settings.

### Step 3: Configure Your Setup

Edit `config.json` to match your WoW installations:

```json
{
    "schema": "aggon/v2",
    "metadata": {
        "name": "my-wow-setup",
        "version": "1.0.0",
        "description": "My WoW addon configuration"
    },
    "installations": {
        "retail": {
            "type": "retail",
            "path": "C:/Program Files (x86)/World of Warcraft/_retail_/Interface/AddOns",
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
            "compatible": ["retail"],
            "ignore": ["README.md", ".gitattributes"]
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

### Step 4: Test Your Configuration

```bash
# Preview what would happen
./aggon.exe plan
```

Expected output:

```
🔄 Creating execution plan...
📋 Execution Plan
=================

📥 Downloads (1):
   - elvui from https://github.com/ElvUI-WotLK/ElvUI

⚙️  Operations (1):
   + Install elvui in retail
```

### Step 5: Apply Configuration

```bash
# Apply the configuration
./aggon.exe switch
```

Expected output:

```
🔄 Planning system changes...
📋 Plan Summary:
   Downloads: 1
   Operations: 1
   Installations: 1

⚡ Applying changes...
📥 Downloading elvui from https://github.com/ElvUI-WotLK/ElvUI...
📦 Executing installation plan for retail...
✅ Successfully applied configuration!
   Generation: 1
   Downloaded: 1 addons
   Installed: 1 installations
   Duration: 15.2s
```

🎉 **Congratulations!** You now have a declarative WoW addon setup!

---

## 🏗️ Building Your Configuration

### Finding Your WoW Installation Paths

#### Retail WoW

```
Default: C:/Program Files (x86)/World of Warcraft/_retail_/Interface/AddOns
```

#### Classic WoW

```
Default: C:/Program Files (x86)/World of Warcraft/_classic_/Interface/AddOns
```

#### Private Servers

Common locations:

```
Ascension: C:/Games/Ascension Launcher/resources/epoch_live/Interface/Addons
Turtle WoW: D:/Games/Turtle WoW/Interface/AddOns
Warmane: C:/Warmane/Data/Interface/AddOns
```

### Adding More Addons

1. **Find the addon's GitHub repository**
2. **Add to the `addons` section:**

```json
{
    "addons": {
        "weakauras": {
            "source": {
                "type": "github",
                "url": "https://github.com/WeakAuras/WeakAuras2",
                "ref": "main"
            },
            "version": "latest",
            "compatible": ["retail"]
        },
        "details": {
            "source": {
                "type": "github",
                "url": "https://github.com/Tercioo/Details-Damage-Meter",
                "ref": "master"
            },
            "version": "latest",
            "compatible": ["retail", "classic"]
        }
    }
}
```

3. **Add to installation addon list:**

```json
{
    "installations": {
        "retail": {
            "addons": ["elvui", "weakauras", "details"]
        }
    }
}
```

4. **Test and apply:**

```bash
./aggon.exe plan
./aggon.exe switch
```

---

## 🎭 Working with Profiles

Profiles let you switch between different configurations.

### Creating Profiles

```json
{
    "profiles": {
        "raiding": {
            "installations": ["retail"],
            "addons": {
                "weakauras": {
                    "ref": "main"
                }
            },
            "description": "Raiding setup with stable addons"
        },
        "testing": {
            "installations": ["retail"],
            "addons": {
                "elvui": {
                    "ref": "development"
                },
                "weakauras": {
                    "ref": "development"
                }
            },
            "description": "Testing with beta addons"
        }
    }
}
```

### Using Profiles

```bash
# Switch to raiding profile
./aggon.exe switch --profile raiding

# Switch to testing profile
./aggon.exe switch --profile testing

# Return to default configuration
./aggon.exe switch
```

---

## 🔄 Managing Generations

Every configuration change creates a new "generation" - a complete snapshot of your system state.

### Viewing Generations

```bash
# List all generations
./aggon.exe generations list
```

Output:

```
📋 Generations:
   1  2024-01-15 14:30:25  Initial setup
   2  2024-01-15 15:45:12  Added WeakAuras
   3  2024-01-15 16:20:43  Updated to beta versions
 * 4  2024-01-15 17:05:18  Raiding profile (current)
```

### Rolling Back

```bash
# Rollback to previous generation
./aggon.exe rollback

# Rollback to specific generation
./aggon.exe rollback 2

# View what generation 2 contained
./aggon.exe generations show 2
```

### Cleaning Up

```bash
# Keep only last 5 generations
./aggon.exe gc --keep 5

# Force cleanup including current
./aggon.exe gc --force
```

---

## 📋 Common Workflows

### Daily Addon Updates

```bash
# Check what would update
./aggon.exe plan

# Apply updates if satisfied
./aggon.exe switch

# If something breaks, rollback
./aggon.exe rollback
```

### Adding New Installation

1. **Add to installations:**

```json
{
    "installations": {
        "classic": {
            "type": "classic",
            "path": "C:/Games/WoW/_classic_/Interface/AddOns",
            "enabled": true,
            "addons": ["questie", "details"]
        }
    }
}
```

2. **Add compatible addons:**

```json
{
    "addons": {
        "questie": {
            "source": {
                "type": "github",
                "url": "https://github.com/AeroScripts/QuestieDev",
                "ref": "master"
            },
            "version": "latest",
            "compatible": ["classic"]
        }
    }
}
```

3. **Test and apply:**

```bash
./aggon.exe plan
./aggon.exe switch
```

### Switching Between Stable and Beta

1. **Create profiles for each:**

```json
{
    "profiles": {
        "stable": {
            "installations": ["retail"],
            "addons": {
                "elvui": {
                    "ref": "refs/tags/13.45",
                    "hash": "sha256:abc123..."
                }
            },
            "description": "Stable pinned versions"
        },
        "beta": {
            "installations": ["retail"],
            "addons": {
                "elvui": {
                    "ref": "development"
                }
            },
            "description": "Latest development versions"
        }
    }
}
```

2. **Switch between them:**

```bash
# Switch to stable
./aggon.exe switch --profile stable

# Switch to beta for testing
./aggon.exe switch --profile beta

# Rollback if beta has issues
./aggon.exe rollback
```

---

## 🛠️ Troubleshooting Common Issues

### Issue: "Path not found"

**Problem**: Installation path doesn't exist

**Solution**:

1. Verify WoW installation path
2. Create AddOns directory if missing:
    ```bash
    mkdir "C:/Games/WoW/_retail_/Interface/AddOns"
    ```

### Issue: "Permission denied creating symlink"

**Problem**: Windows requires administrator privileges for symlinks

**Solutions**:

1. Run as administrator:

    ```bash
    # Right-click command prompt → "Run as administrator"
    ./aggon.exe switch
    ```

2. Enable Developer Mode (Windows 10/11):
    - Settings → Update & Security → For developers → Developer Mode

### Issue: "Configuration parse error"

**Problem**: Invalid JSON syntax

**Solution**:

1. Validate JSON syntax in editor
2. Check for missing commas, brackets
3. Use online JSON validator

### Issue: "Addon not compatible"

**Problem**: Addon incompatible with installation type

**Solution**:

1. Check addon compatibility:

    ```json
    {
        "addons": {
            "retail-addon": {
                "compatible": ["retail"] // Only works with retail
            }
        }
    }
    ```

2. Don't add to incompatible installations

### Issue: "Hash mismatch"

**Problem**: Downloaded content doesn't match expected hash

**Solutions**:

1. Clear cache and retry:

    ```bash
    rm -rf .aggon/store
    ./aggon.exe switch
    ```

2. Update hash in configuration:
    ```json
    {
        "hash": null // Let AGGON calculate new hash
    }
    ```

---

## 📈 Advanced Configuration

### Multi-Installation Setup

```json
{
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
            "enabled": true,
            "addons": ["elvui-classic", "questie", "details"]
        },
        "ascension": {
            "type": "custom",
            "path": "C:/Games/Ascension/Interface/Addons",
            "enabled": true,
            "addons": ["elvui-epoch", "pfquest"]
        }
    }
}
```

### Version Pinning for Stability

```json
{
    "addons": {
        "elvui": {
            "source": {
                "type": "github",
                "url": "https://github.com/ElvUI-WotLK/ElvUI",
                "ref": "refs/tags/13.45"
            },
            "version": "13.45",
            "hash": "sha256:a1b2c3d4e5f6...",
            "compatible": ["retail"]
        }
    }
}
```

### Custom Folder Names

```json
{
    "addons": {
        "pfquest": {
            "source": {
                "type": "github",
                "url": "https://github.com/shagu/pfQuest",
                "ref": "main"
            },
            "folder": "pfQuest-Custom",
            "compatible": ["classic", "custom"]
        }
    }
}
```

---

## 🔐 Security Best Practices

### 1. Verify Sources

-   ✅ Only use trusted GitHub repositories
-   ✅ Verify repository ownership
-   ✅ Check repository activity and stars

### 2. Pin Important Addons

```json
{
    "addons": {
        "critical-addon": {
            "version": "1.0.0",
            "hash": "sha256:exact-hash-here",
            "ref": "refs/tags/v1.0.0"
        }
    }
}
```

### 3. Enable Verification

```json
{
    "settings": {
        "verify_hashes": true
    }
}
```

### 4. Review Changes

```bash
# Always review before applying
./aggon.exe plan
# Read the output carefully
./aggon.exe switch
```

---

## 💡 Tips and Tricks

### 1. Version Control Your Config

```bash
git init
git add aggon.json
git commit -m "Initial AGGON configuration"
```

### 2. Backup Before Major Changes

```bash
# Create backup
cp aggon.json aggon.json.backup

# Test changes
./aggon.exe plan

# Apply if good, restore if bad
```

### 3. Use Meaningful Descriptions

```json
{
    "metadata": {
        "description": "Raiding setup for Mythic+ with beta WeakAuras"
    },
    "profiles": {
        "raid": {
            "description": "Optimized for raid performance"
        }
    }
}
```

### 4. Organize Complex Configs

```json
{
  "addons": {
    "// UI FRAMEWORK": "comment",
    "elvui": { ... },
    "elvui-addonskins": { ... },

    "// COMBAT": "comment",
    "details": { ... },
    "bigwigs": { ... },

    "// UTILITIES": "comment",
    "weakauras": { ... },
    "postal": { ... }
  }
}
```

---

## 🎓 Learning More

### Documentation

-   [README.md](README.md) - Complete overview
-   [CONFIGURATION.md](CONFIGURATION.md) - Detailed configuration reference

### Getting Help

-   Create GitHub issues for bugs
-   Check existing issues for solutions
-   Join community discussions

### Contributing

-   Fork the repository
-   Create feature branches
-   Submit pull requests
-   Help improve documentation

---

**🎉 You're now ready to master declarative addon management with AGGON v2!**

_Happy addon managing! 🏺_

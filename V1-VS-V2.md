# 🏺 AGGON v1 vs v2: Complete Comparison Guide

Choosing between the classic imperative system and the new declarative system.

## 🎯 Quick Decision Guide

**Choose AGGON v1 if you:**

-   ✅ Want a simple, familiar interface
-   ✅ Prefer interactive menus over configuration files
-   ✅ Have basic addon management needs
-   ✅ Don't need advanced rollback capabilities
-   ✅ Want something that "just works" immediately

**Choose AGGON v2 if you:**

-   ✅ Want reproducible, version-controlled setups
-   ✅ Need atomic operations and safe rollbacks
-   ✅ Manage multiple WoW installations
-   ✅ Want infrastructure-as-code approach
-   ✅ Need profile-based configurations
-   ✅ Value system integrity and reliability

---

## 📊 Feature Comparison

| Feature                | AGGON v1 (Classic)      | AGGON v2 (Declarative)   |
| ---------------------- | ----------------------- | ------------------------ |
| **Configuration**      | Interactive menu + JSON | Declarative JSON file    |
| **Interface**          | Menu-driven             | Command-line with config |
| **Learning Curve**     | Easy                    | Moderate                 |
| **Setup Time**         | 5 minutes               | 15 minutes               |
| **Reproducibility**    | Limited                 | Complete                 |
| **Rollback**           | Basic backups           | Instant generations      |
| **Multi-Installation** | Basic support           | Full support             |
| **Version Control**    | Config file only        | Complete system state    |
| **Atomic Operations**  | No                      | Yes                      |
| **State Management**   | Mutable                 | Immutable                |
| **Deduplication**      | No                      | Automatic                |
| **Hash Verification**  | No                      | Yes                      |
| **Profiles**           | No                      | Yes                      |
| **Recovery**           | Manual restore          | Atomic rollback          |
| **Performance**        | Good                    | Excellent                |
| **Debugging**          | Basic                   | Advanced                 |

---

## 🏛️ Architecture Comparison

### AGGON v1: Traditional Architecture

```
config.json ──► Interactive Menu ──► Direct File Operations
     │                │                      │
     │                │                      ▼
     │                │              WoW/Interface/AddOns/
     │                │                   ├── ElvUI/
     │                │                   ├── WeakAuras/
     │                │                   └── Details/
     │                │
     │                ▼
     └──► Basic Cache ──► Download ──► Extract ──► Install
```

**Characteristics:**

-   Direct file manipulation
-   In-place updates
-   Basic caching
-   Linear process flow

### AGGON v2: Declarative Architecture

```
aggon-declarative.json ──► Build Plan ──► Content Store ──► Symlink Farm
            │                  │              │               │
            │                  │              │               ▼
            │                  │              │       WoW/Interface/AddOns/
            │                  │              │         ├── ElvUI -> store/abc123
            │                  │              │         ├── WeakAuras -> store/def456
            │                  │              │         └── Details -> store/ghi789
            │                  │              │
            │                  │              ▼
            │                  │       .aggon/store/
            │                  │         ├── ab/c123... (ElvUI)
            │                  │         ├── de/f456... (WeakAuras)
            │                  │         └── gh/i789... (Details)
            │                  │
            │                  ▼
            └──► Generation Manager ──► .aggon/generations/
                       │                  ├── 1-2024-01-15/
                       │                  ├── 2-2024-01-16/
                       │                  └── current -> 2-2024-01-16/
                       │
                       ▼
                Atomic Operations ──► Rollback Capability
```

**Characteristics:**

-   Immutable content store
-   Generation-based snapshots
-   Symlink-based installations
-   Atomic state transitions

---

## 🎮 User Experience Comparison

### AGGON v1: Interactive Experience

```bash
$ ./aggon.exe

🏺 AGGON
World of Warcraft Addon Manager
================================

📁 1 Installation Path(s) Configured

📂 Ascension (11 addons)
   C:/Games/Ascension/Interface/Addons

Menu Options:
─────────────
1. 🚀 Install/Update All Addons
2. ➕ Add New Addon
3. 📁 Add Installation Path
4. 💾 Backup All Addon Directories
5. ✨ Format Config File
q. Quit

Enter choice: 1

🚀 Installing/Updating Addons
=============================

📂 Ascension
   C:/Games/Ascension/Interface/Addons

   ✅ ElvUI-Epoch - Up to date (from cache)
   ✅ WeakAuras - Updated successfully
   ...
```

**User Experience:**

-   🟢 Intuitive menu system
-   🟢 Guided workflows
-   🟢 Immediate feedback
-   🟡 Limited customization
-   🔴 No configuration versioning

### AGGON v2: Declarative Experience

```bash
$ ./aggon-declarative.exe init
🎉 Initializing declarative AGGON...
✅ Created aggon-declarative.json

$ ./aggon-declarative.exe plan
🔄 Creating execution plan...
📋 Execution Plan
=================

📥 Downloads (2):
   - elvui from https://github.com/ElvUI-WotLK/ElvUI
   - weakauras from https://github.com/WeakAuras/WeakAuras2

⚙️  Operations (2):
   + Install elvui in retail
   + Install weakauras in retail

$ ./aggon-declarative.exe switch
⚡ Applying changes...
✅ Successfully applied configuration!
   Generation: 1
   Downloaded: 2 addons
   Installed: 1 installations
   Duration: 8.4s
```

**User Experience:**

-   🟢 Infrastructure-as-code approach
-   🟢 Predictable, reproducible results
-   🟢 Powerful rollback capabilities
-   🟡 Requires learning configuration syntax
-   🟡 More setup time initially

---

## 📝 Configuration Comparison

### AGGON v1 Configuration

**File: config.json**

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
            },
            {
                "name": "WeakAuras",
                "url": "https://github.com/WeakAuras/WeakAuras2",
                "folder": "WeakAuras",
                "latest_release": true
            }
        ]
    }
]
```

**Characteristics:**

-   Simple array-based structure
-   Direct addon-to-installation mapping
-   Limited metadata
-   Basic options

### AGGON v2 Configuration

**File: aggon-declarative.json**

```json
{
    "schema": "aggon/v2",
    "metadata": {
        "name": "my-setup",
        "version": "1.0.0",
        "description": "Complete addon configuration"
    },
    "installations": {
        "ascension": {
            "type": "custom",
            "path": "C:/Games/Ascension/Interface/Addons",
            "enabled": true,
            "addons": ["elvui-epoch", "weakauras"]
        }
    },
    "addons": {
        "elvui-epoch": {
            "source": {
                "type": "github",
                "url": "https://github.com/Bennylavaa/ElvUI-Epoch",
                "ref": "main"
            },
            "version": "latest",
            "compatible": ["custom"],
            "ignore": ["README.md", ".gitattributes"]
        },
        "weakauras": {
            "source": {
                "type": "github",
                "url": "https://github.com/WeakAuras/WeakAuras2",
                "ref": "main"
            },
            "version": "latest",
            "compatible": ["retail", "custom"],
            "folder": "WeakAuras"
        }
    },
    "profiles": {
        "minimal": {
            "installations": ["ascension"],
            "addons": {
                "elvui-epoch": { "enabled": false }
            }
        }
    },
    "settings": {
        "verify_hashes": true,
        "backup_generations": 10
    }
}
```

**Characteristics:**

-   Structured, schema-based format
-   Separation of concerns (installations vs addons)
-   Rich metadata and versioning
-   Advanced features (profiles, settings)

---

## 🔄 Workflow Comparison

### AGGON v1 Workflow

```
1. Run ./aggon.exe
2. Navigate interactive menu
3. Add installation paths
4. Add addons one by one
5. Install/update all
6. Manual backup if needed
```

**Typical Day:**

```bash
./aggon.exe
# Select "1. Install/Update All Addons"
# Wait for completion
# Exit menu
```

### AGGON v2 Workflow

```
1. Create/edit aggon-declarative.json
2. Plan changes: ./aggon-declarative.exe plan
3. Apply changes: ./aggon-declarative.exe switch
4. Automatic generation created
5. Rollback if needed: ./aggon-declarative.exe rollback
```

**Typical Day:**

```bash
# Edit configuration file
vim aggon-declarative.json

# Preview changes
./aggon-declarative.exe plan

# Apply if satisfied
./aggon-declarative.exe switch

# Or rollback if issues
./aggon-declarative.exe rollback
```

---

## 🛡️ Safety and Reliability

### AGGON v1 Safety

**Backup System:**

-   ✅ Creates ZIP backups before changes
-   ❌ Manual restore process
-   ❌ Can fail during installation
-   ❌ No atomic operations

**Risk Factors:**

-   🔴 **Partial Updates**: Installation can fail halfway
-   🔴 **File Corruption**: Direct file modification risks
-   🟡 **Recovery Time**: Manual backup restoration
-   🟡 **State Verification**: Limited integrity checking

### AGGON v2 Safety

**Generation System:**

-   ✅ Complete system snapshots
-   ✅ Instant atomic rollback
-   ✅ Never modifies installed files
-   ✅ Content integrity verification

**Risk Mitigation:**

-   🟢 **Atomic Operations**: All-or-nothing changes
-   🟢 **Immutable Storage**: Files never corrupted
-   🟢 **Instant Recovery**: One-command rollback
-   🟢 **State Verification**: SHA256 hash checking

---

## 📈 Performance Comparison

### AGGON v1 Performance

```
Initial Setup:    Very Fast   (2-3 minutes)
Daily Updates:    Fast        (1-2 minutes)
Storage Usage:    Efficient   (no duplication overhead)
Memory Usage:     Low         (simple operations)
Network Usage:    Standard    (download when needed)
```

### AGGON v2 Performance

```
Initial Setup:    Moderate    (5-10 minutes)
Daily Updates:    Very Fast   (30 seconds with cache)
Storage Usage:    Higher      (immutable store overhead)
Memory Usage:     Moderate    (hash calculations)
Network Usage:    Optimized   (smart caching, deduplication)
```

**Performance Notes:**

-   v2 has higher initial overhead but better long-term performance
-   v2's deduplication saves bandwidth over time
-   v2's caching makes subsequent operations much faster

---

## 🎯 Use Case Recommendations

### Choose AGGON v1 for:

#### **Casual Players**

-   Simple addon needs
-   Single WoW installation
-   Infrequent addon changes
-   Prefer point-and-click interfaces

#### **Quick Setup**

-   Need working system immediately
-   Don't want to learn configuration syntax
-   Basic backup needs sufficient

#### **Simple Environments**

-   1-5 addons total
-   Standard WoW installations only
-   No version control requirements

### Choose AGGON v2 for:

#### **Power Users**

-   Multiple WoW installations
-   Complex addon configurations
-   Need reproducible setups
-   Value system reliability

#### **Development/Testing**

-   Testing beta addon versions
-   Need quick environment switching
-   Want rollback capabilities
-   Infrastructure-as-code approach

#### **Team/Guild Management**

-   Sharing addon configurations
-   Standardized raid setups
-   Version-controlled configurations
-   Multiple environment profiles

#### **System Administrators**

-   Managing multiple computers
-   Automated deployment needs
-   Audit trails required
-   Zero-downtime updates

---

## 🔄 Migration Strategies

### From v1 to v2

**Recommended Approach:**

1. **Parallel Setup**: Keep v1 running while setting up v2
2. **Gradual Migration**: Start with one installation
3. **Testing Phase**: Verify all addons work correctly
4. **Full Switchover**: Move all installations to v2
5. **Cleanup**: Remove v1 when confident

**Migration Commands:**

```bash
# Keep v1 config as reference
cp config.json config-v1-backup.json

# Initialize v2
./aggon-declarative.exe init

# Manually convert configuration
# (automatic converter planned for future)

# Test v2 setup
./aggon-declarative.exe plan
./aggon-declarative.exe switch --test

# Apply when ready
./aggon-declarative.exe switch
```

### From Manual to Either Version

**For v1 (Easier):**

1. Document current addons
2. Run `./aggon.exe`
3. Add installation paths via menu
4. Add addons via menu
5. Install all

**For v2 (More Powerful):**

1. Document current addon sources
2. Create declarative configuration
3. Define installations and addons
4. Test with `./aggon-declarative.exe plan`
5. Apply with `./aggon-declarative.exe switch`

---

## 🏆 Final Recommendation

### For Most Users: **Start with AGGON v1**

-   Learn addon management concepts
-   Get familiar with AGGON workflow
-   Migrate to v2 when needs grow

### For Advanced Users: **Go Directly to AGGON v2**

-   Skip the learning curve if comfortable with config files
-   Get the most powerful features immediately
-   Better long-term investment

### For Organizations: **AGGON v2 Only**

-   Infrastructure-as-code approach
-   Version control integration
-   Standardized deployments
-   Audit capabilities

---

## 📚 Learning Path

### AGGON v1 → v2 Migration Path

1. **Master v1 Basics** (1-2 weeks)

    - Understand addon management concepts
    - Learn GitHub integration
    - Practice backup/restore

2. **Learn v2 Concepts** (1 week)

    - Understand declarative configuration
    - Learn generation management
    - Practice with simple setups

3. **Advanced v2 Features** (ongoing)
    - Master profiles and environments
    - Implement version control workflows
    - Optimize for your specific needs

### Direct v2 Learning Path

1. **Configuration Basics** (2-3 days)

    - Learn JSON syntax
    - Understand schema structure
    - Create first working config

2. **Core Operations** (1 week)

    - Master plan/switch/rollback cycle
    - Understand generation management
    - Practice with real addons

3. **Advanced Features** (ongoing)
    - Implement profiles
    - Version control integration
    - Team collaboration workflows

---

**🎯 The bottom line: Both versions are excellent. Choose based on your complexity needs and comfort with declarative systems!**

_Happy addon managing! 🏺_

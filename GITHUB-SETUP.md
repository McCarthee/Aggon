# 🚀 GitHub Repository Setup Instructions

Follow these steps to upload your AGGON project to GitHub with proper branch organization.

## 📋 Current Repository Structure

Your local repository is now properly organized with:

- **`master` branch**: AGGON v1 (Classic) - Interactive, menu-driven system
- **`declarative-v2` branch**: AGGON v2 (Declarative) - NixOS-inspired system

## 🌐 Create GitHub Repository

1. **Go to GitHub** and create a new repository:
   - Repository name: `aggon` (or your preferred name)
   - Description: `🏺 AGGON - World of Warcraft Addon Manager`
   - Make it **Public** (recommended for open source)
   - **DO NOT** initialize with README, .gitignore, or license (we already have them)

2. **Copy the repository URL** (should look like: `https://github.com/yourusername/aggon.git`)

## 📤 Upload to GitHub

Run these commands to push both branches:

```bash
# Add GitHub as remote origin
git remote add origin https://github.com/yourusername/aggon.git

# Push master branch (v1)
git push -u origin master

# Push declarative-v2 branch (v2)
git push -u origin declarative-v2
```

## 🏠 Set Default Branch (Recommended)

After uploading, set the **declarative-v2** branch as the default since it's the more advanced version:

1. Go to your GitHub repository
2. Click **Settings** tab
3. Scroll to **Branches** section
4. Change default branch from `master` to `declarative-v2`
5. Confirm the change

## 📝 GitHub Repository Settings

### Description
```
🏺 AGGON - World of Warcraft Addon Manager with both classic (imperative) and declarative (NixOS-inspired) versions
```

### Topics/Tags
Add these topics to help users find your project:
```
world-of-warcraft, addon-manager, wow, nix-inspired, declarative, go, golang, automation
```

### Branch Protection (Optional)
Consider protecting both branches to require pull requests:
1. Go to **Settings** → **Branches**
2. Add protection rules for `master` and `declarative-v2`
3. Enable "Require pull requests before merging"

## 🎯 Branch Navigation

After setup, users can easily switch between versions:

### For v1 (Classic):
- Default view shows `declarative-v2` with link to master
- Users can click "Switch branches" → `master` for v1

### For v2 (Declarative):  
- Default view shows comprehensive documentation
- Full feature comparison and migration guides available

## 📚 Documentation Strategy

The repository now provides clear documentation paths:

- **`master` branch**: Simple README focusing on v1 usage
- **`declarative-v2` branch**: Comprehensive documentation suite:
  - `README.md` - Complete overview
  - `CONFIGURATION.md` - Detailed configuration reference  
  - `GETTING-STARTED.md` - Step-by-step user guide
  - `V1-VS-V2.md` - Complete comparison guide

## 🎉 Ready to Go!

Your repository is now ready for GitHub with:
- ✅ Clean branch organization
- ✅ Comprehensive documentation
- ✅ Proper .gitignore and licensing
- ✅ Clear navigation between versions
- ✅ Professional presentation

## 🔗 Next Steps

After uploading:
1. Add a nice repository banner/logo
2. Create GitHub releases for both versions
3. Set up GitHub Actions for automated builds (optional)
4. Consider adding issue templates
5. Add contribution guidelines

---

*Ready to share your NixOS-inspired WoW addon manager with the world! 🌟*

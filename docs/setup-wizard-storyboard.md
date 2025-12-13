# Wails Setup - OOBE Storyboard

> An Apple/Microsoft-style Out-of-Box Experience for `wails3 setup`

---

## Design Philosophy

Transform the wizard into an OOBE (Out-of-Box Experience):

- **No footer navigation** - Buttons appear contextually within content
- **Full-screen immersive pages** - Each step is its own moment
- **Progressive disclosure** - Show only what's needed
- **Branching paths** - Adapt flow based on system state
- **Conversational tone** - Guide users naturally through decisions

---

## Page Layout Template

All screens (except Splash) use this consistent layout:

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                         [Icon]                              |
|                                                             |
|                    Large Title Here                         |
|                                                             |
|           Supporting text in a friendly tone                |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |              Content Area                         |   |
|     |           (varies per screen)                     |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                    [ Primary Action ]                       |
|                       secondary link                        |
|                                                             |
+-------------------------------------------------------------+
```

Key principles:
- Theme toggle always visible (top-left)
- Centered content with generous whitespace
- Buttons appear only when relevant
- No persistent footer or step indicators

---

## Flow Diagram

```
                              +----------------+
                              |    Splash      |
                              | "Let's Start"  |
                              +-------+--------+
                                      |
                                      v
                         +------------------------+
                         |   Checking System...   |
                         |   (auto-transition)    |
                         +------------+-----------+
                                      |
                   +------------------+------------------+
                   |                                     |
                   v                                     v
       +-----------------------+           +-----------------------+
       |   All deps OK         |           |   Missing deps        |
       |   "You're all set!"   |           |   "Almost there"      |
       +-----------+-----------+           +-----------+-----------+
                   |                                   |
                   |                                   v
                   |                       +-----------------------+
                   |                       | Install Instructions  |
                   |                       | [Copy] [Re-check]     |
                   |                       +-----------+-----------+
                   |                                   |
                   +<----------------------------------+
                   |
                   v
       +-----------------------+
       |  Cross-Platform?      |
       |  (always shown)       |
       +-----------+-----------+
                   |
       +-----------+-----------+
       |                       |
       v                       v
 +------------+       +------------------------+
 | Not now    |       |   Yes, set this up     |
 +-----+------+       +------------+-----------+
       |                          |
       |              +-----------+-----------+
       |              |                       |
       |              v                       v
       |    +------------------+    +------------------+
       |    | Docker missing?  |    | Docker ready?    |
       |    | -> Install link  |    | -> Build image   |
       |    +--------+---------+    +--------+---------+
       |             |                       |
       |             v                       v
       |    +------------------+    +------------------+
       |    | [Check Again]    |    | Building...      |
       |    | Skip ->          |    | (can continue)   |
       |    +--------+---------+    +--------+---------+
       |             |                       |
       +<------------+-----------------------+
       |
       v
 +-----------------------+
 |  Tell us about you    |
 |  (Author details)     |
 +-----------+-----------+
             |
             v
 +-----------------------+
 |  Project Defaults     |
 |  (Bundle ID, etc)     |
 +-----------+-----------+
             |
             v
 +-----------------------+
 |  Code Signing?        |
 |  [Skip] or [Setup]    |
 +-----------+-----------+
             |
 +-----------+-----------+
 |                       |
 v                       v
+-------------+   +-------------+
| macOS       |   | Windows     |
| Signing     |   | Signing     |
+------+------+   +------+------+
       |                 |
       +--------+--------+
                |
                v
 +-----------------------+
 |     All Done!         |
 |   "Start Building"    |
 +-----------------------+
```

---

## Screen-by-Screen Storyboard

### 1. Splash Screen

**Purpose:** Welcome, set tone, minimal interaction

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                                                             |
|                                                             |
|                       [Wails Logo]                          |
|                     with subtle glow                        |
|                                                             |
|                                                             |
|                   Welcome to Wails                          |
|                                                             |
|           Build beautiful cross-platform apps               |
|                using Go and web tech                        |
|                                                             |
|                                                             |
|                    [ Let's Start ]                          |
|                                                             |
|                                                             |
+-------------------------------------------------------------+
```

**Elements:**
- Theme toggle (top-left) - only UI control
- Logo with glow animation
- Minimal copy
- Single "Let's Start" button

**On click:** Triggers dependency check, transitions to Checking screen

---

### 2. Checking System (Transitional)

**Purpose:** Brief loading while checking dependencies (2-3 seconds)

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                                                             |
|                                                             |
|                        [Spinner]                            |
|                                                             |
|                 Checking your system...                     |
|                                                             |
|               This will only take a moment                  |
|                                                             |
|                                                             |
+-------------------------------------------------------------+
```

**Behavior:**
- No user interaction needed
- Auto-advances when check completes
- Minimum 1.5s display (perceived thoroughness)

---

### 3a. System Ready (Happy Path)

**When:** All required dependencies installed

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                          [Checkmark]                        |
|                         (animated)                          |
|                                                             |
|                      You're all set!                        |
|                                                             |
|              Your system has everything needed              |
|                     to build Wails apps                     |
|                                                             |
|                                                             |
|                       [ Continue ]                          |
|                                                             |
+-------------------------------------------------------------+
```

**Elements:**
- Green animated checkmark
- Simple confirmation message (no list of deps)
- Single "Continue" button

**Next:** Goes to Cross-Platform question

---

### 3b. Missing Dependencies

**When:** Required dependencies not found

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                          [Warning]                          |
|                       (amber icon)                          |
|                                                             |
|                      Almost there!                          |
|                                                             |
|            A few things need to be installed                |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  [X] npm                                          |   |
|     |      Required for frontend tooling                |   |
|     |                                                   |   |
|     |  [X] GTK Development Libraries                    |   |
|     |      Required for Linux builds                    |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|     Run this command to install everything:                 |
|     +---------------------------------------------------+   |
|     |  sudo pacman -S npm gtk3 webkit2gtk      [Copy]   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                      [ Check Again ]                        |
|                                                             |
|                     Continue anyway ->                      |
|                                                             |
+-------------------------------------------------------------+
```

**Dynamic content:**
- Platform-specific package manager (apt/dnf/pacman/brew)
- Combined command when possible
- Individual items if package managers differ

**Actions:**
- Copy button for install command
- "Check Again" - re-runs dependency check
- "Continue anyway" - for advanced users (subtle link)

---

### 4. Cross-Platform Question

**When:** Always shown (Docker status doesn't matter)

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                    [Platform Icons]                         |
|                   Win / Mac / Linux                         |
|                                                             |
|            Build for multiple platforms?                    |
|                                                             |
|         Wails can compile your app for Windows,             |
|        macOS, and Linux from a single machine               |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |   From your {platform} machine, build for:        |   |
|     |                                                   |   |
|     |   - Windows (.exe)                                |   |
|     |   - macOS (.app)                                  |   |
|     |   - Linux (AppImage, .deb, etc)                   |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                   [ Yes, set this up ]                      |
|                                                             |
|                      Not right now ->                       |
|                                                             |
+-------------------------------------------------------------+
```

**Branch logic:**
- Always shown to everyone
- "Yes" -> Checks Docker status, installs if needed, builds image
- "Not right now" -> Author details (faster path)

---

### 5. Setting Up Cross-Platform Builds

**When:** User chose "Yes" to cross-platform builds

This is a single screen that handles all Docker states automatically:

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                        [Docker]                             |
|                                                             |
|             Setting up cross-platform builds                |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  [check] Checking Docker...                       |   |
|     |  [spinner] Installing Docker...                   |   |
|     |  [ ] Building cross-compiler image...             |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|               [ Continue in background ]                    |
|                                                             |
|        The build will complete while you continue           |
|                                                             |
+-------------------------------------------------------------+
```

**States shown in sequence:**
1. Checking Docker... (quick check)
2. Installing Docker... (if not installed - opens install instructions)
3. Starting Docker... (if not running - prompts user)
4. Building image... (progress bar, can continue in background)

**If Docker needs installing:**
```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                        [Docker]                             |
|                                                             |
|                   Install Docker                            |
|                                                             |
|       Cross-platform builds require Docker Desktop          |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |   [Download Docker Desktop]                       |   |
|     |                                                   |   |
|     |   After installing, come back and we'll           |   |
|     |   continue setting up.                            |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                      [ Check Again ]                        |
|                                                             |
|                   Skip for now ->                           |
|                                                             |
+-------------------------------------------------------------+
```

**If Docker not running:**
```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                        [Docker]                             |
|                                                             |
|                   Start Docker                              |
|                                                             |
|          Please start Docker Desktop to continue            |
|                                                             |
|                      [ Check Again ]                        |
|                                                             |
|                   Skip for now ->                           |
|                                                             |
+-------------------------------------------------------------+
```

**Building image (can continue in background):**
```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                        [Docker]                             |
|                                                             |
|             Building cross-compiler image                   |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  [===================>          ] 65%             |   |
|     |                                                   |   |
|     |  This may take several minutes                    |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|               [ Continue in background ]                    |
|                                                             |
+-------------------------------------------------------------+
```

**Key UX:**
- Single flow handles all Docker states
- User can wait OR continue while image builds
- Background status indicator on subsequent screens

---

### 6. Author Details

**Purpose:** Personalize the setup

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                         [User]                              |
|                                                             |
|                Tell us about yourself                       |
|                                                             |
|          This information will be used in your apps         |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  Your Name                                        |   |
|     |  +---------------------------------------------+  |   |
|     |  | Jane Developer                              |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     |  Company (optional)                               |   |
|     |  +---------------------------------------------+  |   |
|     |  | Acme Corp                                   |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                       [ Continue ]                          |
|                                                             |
+-------------------------------------------------------------+
```

**Pre-population:**
- Name from `git config user.name` if available

---

### 7. Project Defaults

**Purpose:** Set app identifier conventions

```
+-------------------------------------------------------------+
|  [Sun/Moon]                           [Docker: 78%]         |
|                                                             |
|                        [Package]                            |
|                                                             |
|                   Project defaults                          |
|                                                             |
|          These will be used when creating new apps          |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  Bundle ID Prefix                                 |   |
|     |  +---------------------------------------------+  |   |
|     |  | com.acme                                    |  |   |
|     |  +---------------------------------------------+  |   |
|     |  Example: com.acme.myapp                          |   |
|     |                                                   |   |
|     |  Default Version                                  |   |
|     |  +---------------------------------------------+  |   |
|     |  | 0.1.0                                       |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     |  Preferred Template                               |   |
|     |  +---------------------------------------------+  |   |
|     |  | React + TypeScript                   [v]    |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                       [ Continue ]                          |
|                                                             |
+-------------------------------------------------------------+
```

**Note:** Docker status indicator (top-right) if build in progress

---

### 8. Code Signing Question

**Purpose:** Ask about app signing (platform-specific)

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                         [Lock]                              |
|                                                             |
|                  Set up code signing?                       |
|                                                             |
|         Code signing lets you distribute apps through       |
|            stores and avoid security warnings               |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |   Code signing is optional during development.    |   |
|     |   You can configure this later in your project.   |   |
|     |                                                   |   |
|     |   Available for your platform:                    |   |
|     |   - macOS Developer ID                            |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                   [ Set up signing ]                        |
|                                                             |
|                         Skip ->                             |
|                                                             |
+-------------------------------------------------------------+
```

**Branch logic:**
- macOS: Show Apple signing options next
- Windows: Show Windows signing options next
- Linux: Skip entirely (no signing needed)

---

### 9a. macOS Code Signing (Optional)

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                        [Apple]                              |
|                                                             |
|                  macOS Code Signing                         |
|                                                             |
|        Required for distribution outside the App Store      |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  Developer ID                                     |   |
|     |  +---------------------------------------------+  |   |
|     |  | Developer ID Application: Jane Dev...       |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     |  Apple ID (for notarization)                      |   |
|     |  +---------------------------------------------+  |   |
|     |  | jane@example.com                            |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     |  Team ID                                          |   |
|     |  +---------------------------------------------+  |   |
|     |  | ABCD1234                                    |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                       [ Continue ]                          |
|                                                             |
+-------------------------------------------------------------+
```

---

### 9b. Windows Code Signing (Optional)

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                       [Windows]                             |
|                                                             |
|                Windows Code Signing                         |
|                                                             |
|           Prevents "Unknown Publisher" warnings             |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  Certificate Path                                 |   |
|     |  +---------------------------------------------+  |   |
|     |  | C:\certs\codesign.pfx               [...]   |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     |  Timestamp Server                                 |   |
|     |  +---------------------------------------------+  |   |
|     |  | http://timestamp.digicert.com               |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                       [ Continue ]                          |
|                                                             |
+-------------------------------------------------------------+
```

---

### 10. Complete Screen

**Purpose:** Celebrate and guide next steps

```
+-------------------------------------------------------------+
|  [Sun/Moon]                                                 |
|                                                             |
|                         [Party]                             |
|                        (animated)                           |
|                                                             |
|                  You're ready to build!                     |
|                                                             |
|          Your development environment is all set up         |
|                                                             |
|     +---------------------------------------------------+   |
|     |                                                   |   |
|     |  Create your first app:                           |   |
|     |  +---------------------------------------------+  |   |
|     |  |  wails3 init -n myapp                [Copy] |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     |  Start developing:                                |   |
|     |  +---------------------------------------------+  |   |
|     |  |  cd myapp && wails3 dev              [Copy] |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     |  Build for production:                            |   |
|     |  +---------------------------------------------+  |   |
|     |  |  wails3 build                        [Copy] |  |   |
|     |  +---------------------------------------------+  |   |
|     |                                                   |   |
|     +---------------------------------------------------+   |
|                                                             |
|                    [ Start Building ]                       |
|                                                             |
|                   Read the documentation                    |
|                                                             |
+-------------------------------------------------------------+
```

**Persistent elements:**
- If Docker build still in progress, show status indicator
- "Start Building" closes the wizard

---

## Branching Scenarios Summary

| Scenario | Flow |
|----------|------|
| **Happy path** | Splash -> Check -> Ready -> Cross-Platform? -> Author -> Defaults -> Done |
| **Missing deps** | Splash -> Check -> Missing -> (install) -> Check -> Ready -> ... |
| **Want cross-platform, Docker ready** | ... -> Cross-Platform? (Yes) -> Building -> Author -> ... |
| **Want cross-platform, no Docker** | ... -> Cross-Platform? (Yes) -> Install Docker -> Check -> Building -> ... |
| **Skip cross-platform** | ... -> Cross-Platform? (Not now) -> Author -> ... |
| **Skip signing** | ... -> Signing? -> Skip -> Done |
| **macOS w/ signing** | ... -> Signing? -> macOS Signing -> Done |
| **Windows w/ signing** | ... -> Signing? -> Windows Signing -> Done |
| **Linux** | ... -> Defaults -> Done (signing skipped) |

---

## State Machine

```typescript
type OOBEStep =
  | 'splash'
  | 'checking'
  | 'deps-ready'
  | 'deps-missing'
  | 'cross-platform-question'
  | 'docker-setup'        // handles install/start/build states
  | 'docker-building'     // can continue in background
  | 'author'
  | 'project-defaults'
  | 'signing-question'
  | 'signing-macos'
  | 'signing-windows'
  | 'complete'
```

---

## Background Processes

### Docker Image Build
- Runs in background while user continues
- Persistent status indicator on subsequent screens
- Blocks wizard close only if still building
- User can wait on build screen OR continue setup

### Dependency Check
- Runs immediately after "Let's Start"
- Shows transitional loading screen
- Auto-advances when complete

---

## Animation Guidelines

- **Page transitions:** Fade + slight slide up (Apple-style)
- **Checkmarks:** Draw animation with scale bounce
- **Progress bars:** Smooth continuous animation
- **Error states:** Gentle shake
- **Buttons:** Subtle glow on hover

---

## Key Differences from Wizard Style

| Wizard | OOBE |
|--------|------|
| Footer with Back/Next | Contextual buttons only |
| Step indicator (1 of 5) | No step counter |
| Linear progression | Branching based on state |
| All screens shown | Skip irrelevant screens |
| Form-heavy | Question-driven |

---

## Implementation Priorities

1. Remove footer navigation from all screens
2. Update Splash to only show theme toggle + "Let's Start"
3. Add "Checking" transitional screen
4. Split dependency results into Ready/Missing paths
5. Add cross-platform question (conditional on Docker)
6. Split Docker into sub-states
7. Add code signing question/screens (conditional on platform)
8. Update Complete screen with next steps

---
title: "Projektstruktur"
description: "Erfahren Sie, wie die Dateien eines typischen Wails-Projekts angeordnet sind"
slug: "guides/dev/project-structure"
sourcePath: "guides/dev/project-structure.md"
---

Diese Seite bietet eine einfache Übersicht über die Dateien, die beim Erstellen eines neuen Projekts mit der Vanilla-Vorlage angelegt werden

```
/
├── main.go                 # Application entry point
├── greetservice.go         # Example backend service exposed to frontend
├── go.mod                  # Go module definition
├── config.yml              # Wails/project configuration
├── Taskfile.yml            # Task runner for dev/build commands
├── build/                  # Packaging and platform-specific assets
│   ├── appicon.png         # Default app icon
│   ├── appicon.icon/       # Icon source files
│   ├── darwin/             # macOS build config
│   ├── windows/            # Windows build config
│   ├── linux/              # Linux build config
│   ├── android/            # Android build config
│   ├── ios/                # iOS build config
│   └── docker/             # Containerized build environment
├── frontend/               # Frontend (Vite + React + TS)
│   ├── index.html          # HTML entry point
│   ├── src/
│   │   └── main.js         # Frontend bootstrap
│   ├── public/             # Static assets
│   ├── dist/               # Built frontend output
│   ├── bindings/           # Auto-generated Go bindings
│   ├── package.json        # Frontend dependencies
│   ├── vite.config.ts      # Vite configuration
│   └── tsconfig.json       # TypeScript configuration
├── bin/                    # Compiled binaries
└── .task/                  # Task runner cache
```

@note{type="info" title="Frontend-Projektdateien"}
Wenn Sie eine andere Startvorlage verwenden, unterscheidet sich der Inhalt von `frontend/src` je nach gewähltem Framework

@end

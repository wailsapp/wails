---
title: "Struktur Proyek"
description: "Pahami tata letak file pada proyek Wails pada umumnya"
slug: "guides/dev/project-structure"
sourcePath: "guides/dev/project-structure.md"
---

Halaman ini adalah referensi sederhana untuk file yang dibuat saat Anda memulai proyek baru menggunakan templat Vanilla

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

@note{type="info" title="File proyek frontend"}
Jika Anda memilih templat awal yang berbeda, isi `frontend/src` akan berbeda tergantung pada framework yang Anda pilih

@end

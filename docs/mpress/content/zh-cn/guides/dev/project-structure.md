---
title: "项目结构"
description: "了解典型 Wails 项目的文件布局"
slug: "guides/dev/project-structure"
sourcePath: "guides/dev/project-structure.md"
---

本页简要介绍使用 Vanilla 模板初始化新项目时创建的文件

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

@note{type="info" title="前端项目文件"}
如果选择其他入门模板，`frontend/src`的内容将因所选框架而异

@end

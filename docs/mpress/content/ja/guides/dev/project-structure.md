---
title: "プロジェクト構成"
description: "一般的な Wails プロジェクトのファイル構成について説明します"
slug: "guides/dev/project-structure"
sourcePath: "guides/dev/project-structure.md"
---

このページは、Vanilla テンプレートを使用して新しいプロジェクトを作成したときに生成されるファイルの簡単なリファレンスです

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

@note{type="info" title="フロントエンドプロジェクトのファイル"}
別のスターターテンプレートを使用する場合、`frontend/src` の内容は選択したフレームワークによって異なります

@end

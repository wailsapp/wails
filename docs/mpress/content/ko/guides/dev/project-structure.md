---
title: "프로젝트 구조"
description: "일반적인 Wails 프로젝트의 파일 구성을 알아봅니다"
slug: "guides/dev/project-structure"
sourcePath: "guides/dev/project-structure.md"
---

이 페이지에서는 Vanilla 템플릿을 사용하여 새 프로젝트를 시작할 때 생성되는 파일을 간단히 설명합니다

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
│   │   └── icons/          # 생성된 hicolor PNG 아이콘 집합
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

`frontend/dist`, `frontend/bindings`, `build/linux/icons` 같은 생성 디렉터리는 해당 빌드 또는 생성 태스크를 처음 실행한 뒤 나타납니다.

@note{type="info" title="프런트엔드 프로젝트 파일"}
다른 시작 템플릿을 사용하면 `frontend/src`의 내용은 선택한 프레임워크에 따라 달라집니다

@end

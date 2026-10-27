---
title: "기여하기"
description: "Wails에 기여하기"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## 기여자 여러분, 환영합니다!

Wails에 대한 기여를 환영합니다! 버그 수정, 기능 추가, 문서 개선 등 어떤 방식으로든 도움을 주시면 감사하겠습니다.

## 기여 방법

### 1. 문제 보고

버그를 발견하셨나요? 다음 내용을 포함하여 [이슈를 등록하세요](https://github.com/wailsapp/wails/issues/new):

- 명확한 설명
- 재현 단계
- 예상 동작과 실제 동작
- 시스템 정보
- 코드 예제

### 2. 문서 개선

수정 PR은 사전 이슈나 실패하는 코드 테스트 없이도 환영합니다.  
다운로드, macOS/Linux/Windows 설치, 실시간 미리 보기 및 변경 사항 제출 방법은 [M-Press로 문서 작성하기](/contributing/documentation/)를 따르세요.

다음과 같은 문서 개선은 언제나 환영합니다:

- 오탈자 및 오류 수정
- 예제 추가
- 설명 명확화
- 콘텐츠 번역

### 3. 코드 제출

풀 리퀘스트를 통해 코드를 기여하세요:

- 버그 수정
- 새 기능
- 성능 개선
- 테스트

### 4. 개선 사항 제안(WEP)

새 기능과 공개 동작 변경에는 Wails Enhancement Proposal(WEP) 절차를 사용합니다. 이 절차를 통해 기능 개발을 투명하게 유지하고 승인된 모든 제안에 구현 담당자가 있도록 보장합니다. 기능 요청 이슈를 등록하지 마세요.

1. 관심도를 파악하려면 선택적으로 GitHub Discussions의 [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) 카테고리나 [Discord](https://discord.gg/JDdSxwjhGf)에 아이디어를 먼저 공유하세요.
2. [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/WEP_TEMPLATE.md)을 `v3/wep/proposals/<proposal name>/proposal.md`에 복사하고 모든 섹션을 작성하세요.
3. 제안서만 포함하고 제목이 `[WEP] <title>`인 초안 풀 리퀘스트를 여세요. 이 PR이 제안을 논의하는 공식 공간입니다.
4. 피드백과 지지를 모으세요(PR의 댓글 및 좋아요 반응). 토론 기간을 최소 2주 두고 제안을 구현할 담당자를 합의하세요.
5. PR을 검토 준비 상태로 표시하세요. 최종 결정은 메인테이너가 내리며, 승인된 제안에는 WEP 번호가 할당되고 병합됩니다.

전체 절차는 [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)에 설명되어 있습니다.

## 시작하기

### 포크 및 복제

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### 소스에서 빌드

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### 테스트 실행

테스트는 최종 스모크 검사가 아니라 변경 사항의 일부입니다. 단위 테스트는 테스트 대상 코드 옆에 두고, 하나의 동작을 여러 입력이나 경계 사례로 검사할 때는 가능하면 테이블 기반 테스트를 사용하세요. 실패 시 어떤 시나리오인지 알 수 있도록 각 사례에 이름을 지정하세요.

새로 추가하거나 변경한 로직은 모두 테스트로 다뤄야 합니다. PR에서 추가하거나 변경한 코드의 Go 문장 커버리지 100%를 목표로 하세요. 저장소 전체의 커버리지 백분율을 변경 사항 테스트의 대체 수단으로 간주해서는 안 됩니다. 예를 들어 특정 OS에서만 발생하는 오류 경로나 실제 하드웨어 없이는 재현하기 어려운 조건처럼 테스트 공백이 타당할 수도 있지만, PR 설명에 해당 공백과 합리적으로 테스트할 수 없는 이유를 명시하세요.

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

통합 테스트 모음, 경쟁 상태 감지 및 전체 CI와 동등한 명령어는 [테스트 및 지속적 통합](/contributing/testing-ci/)을 참조하세요.

## 변경 사항 만들기

### 브랜치 생성

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### 변경 작업 수행

1. Go 규칙에 따라 **코드를 작성하세요**
2. 새 기능에 대한 **테스트를 추가하세요**
3. 필요한 경우 **문서를 업데이트하세요**
4. 문제가 발생하지 않는지 확인하도록 **테스트를 실행하세요**
5. 명확한 메시지로 **변경 사항을 커밋하세요**

### 커밋 지침

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### 풀 리퀘스트 제출

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## 풀 리퀘스트 지침

### 좋은 PR 설명

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### PR 체크리스트

- [ ] 코드가 Go 규칙을 준수함
- [ ] 테스트를 추가/업데이트함
- [ ] 문서를 업데이트함
- [ ] 모든 테스트를 통과함
- [ ] 호환성을 깨는 변경 사항이 없음(또는 문서화함)
- [ ] 커밋 메시지가 명확함

## 코드 지침

### Go 코드 스타일

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### 테스트

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 문서

### 문서 작성

전체 문서 작업 흐름은 [M-Press로 문서 작성하기](/contributing/documentation/)를 따르세요. 여기에서는 고정된 버전, 실행 파일 및 Go 설치 옵션, 편집 가능한 소스 파일, 탐색 메뉴, 번역 검사, `master`을 대상으로 한 PR 제출 방법을 다룹니다.

문서만 변경하는 작업에는 Wails 애플리케이션 빌드나 플랫폼 종속성이 필요하지 않습니다. M-Press는 저장소 루트에서 `mpress dev`을 사용하여 기존 사이트를 미리 보여 줍니다.

### 문서 스타일

- 국제 영어 철자법 사용
- 문제부터 설명
- 실행 가능한 예제 제공
- 문제 해결 방법 포함
- 관련 콘텐츠 상호 참조

## 커뮤니티

### 도움받기

- **Discord:** [커뮤니티에 참여하세요](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions:** 질문하기
- **GitHub Issues:** 버그 보고하기

### 행동 강령

서로 존중하고 포용하며 전문적으로 행동하세요. 우리 모두는 훌륭한 소프트웨어를 함께 만들기 위해 이곳에 있습니다. 자세한 내용은 [행동 강령](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md)을 참조하세요.

## 기여자 표기

기여자는 다음 항목에 이름이 게재됩니다:

- 릴리스 노트
- 기여자 목록
- GitHub 인사이트

Wails에 기여해 주셔서 감사합니다! 🎉

## 다음 단계

@cards{cols="2"}
◆ GitHub 저장소
Wails 저장소를 방문하세요.

[GitHub에서 보기 →](https://github.com/wailsapp/wails)

---
◆ Discord 커뮤니티
커뮤니티에 참여하세요.

[Discord 참여 →](https://discord.gg/JDdSxwjhGf)

---
📖 문서
문서를 읽어보세요.

[문서 둘러보기 →](/quick-start/why-wails/)

@end

---
title: "기여하기"
description: "Wails에 기여하기"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## 기여자 여러분을 환영합니다!

Wails에 대한 기여를 환영합니다! 버그 수정, 기능 추가, 문서 개선 등 어떤 도움이든 감사드립니다.

## 기여 방법

### 1. 문제 보고

버그를 발견하셨나요? 다음 내용을 포함하여 [이슈를 등록하세요](https://github.com/wailsapp/wails/issues/new):

- 명확한 설명
- 재현 단계
- 예상 동작과 실제 동작
- 시스템 정보
- 코드 예제

### 2. 문서 개선

문서 수정 PR은 사전 이슈나 실패하는 코드 테스트 없이도 제출할 수 있습니다.  
M-Press로 변경 사항을 미리 보고 검증하려면 [문서 수정하기](/contributing/documentation/)를 따르세요.

다음과 같은 문서 개선은 언제든 환영합니다:

- 오탈자와 오류 수정
- 예제 추가
- 설명 명확히 하기
- 콘텐츠 번역

### 3. 코드 제출

풀 리퀘스트를 통해 다음과 같은 코드를 기여하세요:

- 버그 수정
- 새 기능
- 성능 개선
- 테스트

### 4. 개선 사항 제안(WEP)

새 기능과 공개 동작의 변경에는 Wails Enhancement Proposal(WEP) 절차를 사용합니다. 이 절차는 기능 개발을 투명하게 유지하고 승인된 모든 제안에 구현 담당자가 있도록 보장합니다. 기능 요청 이슈를 등록하지 마세요.

1. 선택 사항으로, 관심도를 파악하기 위해 GitHub Discussions의 [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) 카테고리나 [Discord](https://discord.gg/JDdSxwjhGf)에서 아이디어를 먼저 공유해도 됩니다.
2. [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md)을 `v3/wep/proposals/<proposal name>/proposal.md`에 복사하고 모든 섹션을 작성하세요.
3. 제안만 포함하는 `[WEP] <title>` 제목의 초안 풀 리퀘스트를 등록하세요. 이 PR이 제안을 논의하는 공식 공간입니다.
4. 피드백과 지지를 모으세요(PR의 댓글과 좋아요 반응). 논의에 최소 2주를 할애하고 제안을 구현할 담당자를 합의하세요.
5. PR을 검토 준비 상태로 표시하세요. 최종 결정은 메인테이너가 내립니다. 승인된 제안에는 WEP 번호가 할당되고 병합됩니다.

전체 절차는 [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)에 문서화되어 있습니다.

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

테스트는 변경 작업의 일부이며, 마지막에 수행하는 단순 확인 절차가 아닙니다. 단위 테스트는 테스트 대상 코드 옆에 두고, 하나의 동작을 여러 입력이나 극단 사례로 검사할 때는 가능한 한 테이블 기반 테스트를 사용하세요. 실패 시 어떤 시나리오인지 알 수 있도록 각 사례에 이름을 지정하세요.

새로 추가하거나 변경한 로직은 완전히 테스트해야 합니다. PR에서 추가하거나 변경한 코드의 Go 구문 커버리지 100%를 목표로 하세요. 저장소 전체의 커버리지 비율을 변경 사항 테스트의 대안으로 간주해서는 안 됩니다. 예를 들어 특정 OS에서만 발생하는 오류 경로나 실제 하드웨어 없이는 재현하기 어려운 조건처럼 커버리지 공백이 타당할 수는 있지만, PR 설명에 해당 공백과 합리적인 테스트가 불가능한 이유를 설명하세요.

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

통합 테스트 스위트, 경쟁 상태 감지 및 CI 전체와 동일한 명령은 [테스트 및 지속적 통합](/contributing/testing-ci/)을 참조하세요.

## 변경 작업

### 브랜치 생성

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### 변경 사항 작성

1. Go 규칙에 따라 **코드를 작성하세요**
2. 새 기능에 대한 **테스트를 추가하세요**
3. 필요한 경우 **문서를 업데이트하세요**
4. 문제가 발생하지 않는지 확인하기 위해 **테스트를 실행하세요**
5. 명확한 메시지와 함께 **변경 사항을 커밋하세요**

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
- [ ] 테스트를 추가하거나 업데이트함
- [ ] 문서를 업데이트함
- [ ] 모든 테스트를 통과함
- [ ] 호환성을 깨는 변경이 없음(또는 문서화함)
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

문서는 M-Press를 사용합니다. `docs/mpress/content/` 아래의 `.md` 파일을 편집한 다음, 저장소 루트에서 미리 보고 유효성을 검사하세요:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### 문서 스타일

- 국제 영어 철자를 사용하세요
- 문제부터 설명하세요
- 실제로 작동하는 예제를 제공하세요
- 문제 해결 방법을 포함하세요
- 관련 콘텐츠를 상호 참조하세요

## 커뮤니티

### 도움받기

- **Discord:** [커뮤니티에 참여하세요](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions:** 질문하세요
- **GitHub Issues:** 버그를 신고하세요

### 행동 강령

서로 존중하고 포용하며 전문적인 태도를 유지하세요. 우리 모두는 훌륭한 소프트웨어를 함께 만들기 위해 이곳에 모였습니다. 자세한 내용은 [행동 강령](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md)을 참조하세요.

## 기여자 표기

기여자는 다음 항목에 이름이 올라갑니다:

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

[Discord 참여하기 →](https://discord.gg/JDdSxwjhGf)

---
📖 문서
문서를 읽어 보세요.

[문서 둘러보기 →](/quick-start/why-wails/)

@end

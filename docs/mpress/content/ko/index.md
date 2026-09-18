---
title: "Go로 데스크톱 앱 빌드하기"
description: "Go와 웹 기술을 사용하는 네이티브 데스크톱 애플리케이션"
banner: {"content":"Wails v3는 현재 베타 버전입니다. \u003ca href=\"https://v2.wails.io\"\u003ev2 문서를 찾고 계신가요?\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"시작하기","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"튜토리얼 보기","variant":"secondary"}],"image":{"alt":"Wails 로고","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"Go와 최신 웹 기술로 아름답고 성능이 뛰어난 애플리케이션을 빌드하세요. 하나의 코드베이스. 세 가지 플랫폼. 브라우저 불필요.*"}
sourcePath: "index.md"
---

<style>
  /* Hero background — the neon "digital Wales" mountain, full-bleed and fixed. */
  body::after {
    content: '';
    position: fixed;
    inset: 0;
    z-index: -1;
    background: url('/digital_wales_master.webp') center center / cover no-repeat;
    opacity: 0.35;
    pointer-events: none;
  }
  
  /* Gradient overlay - highest z-index */
  html::before {
    content: '';
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: 
      linear-gradient(to bottom, var(--sl-color-bg) 0%, transparent 15%, transparent 85%, var(--sl-color-bg) 100%),
      linear-gradient(to right, var(--sl-color-bg) 0%, transparent 10%, transparent 90%, var(--sl-color-bg) 100%);
    z-index: 0;
    pointer-events: none;
  }
  
  
  /* Hero padding */
  .hero {
    padding-top: 3rem !important;
    padding-bottom: 2rem !important;
  }
  
  .small-buttons {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-top: 1rem;
  }
  
  .small-buttons a {
    display: inline-block;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    border-radius: 6px;
    text-decoration: none;
    background: var(--sl-color-gray-6);
    color: var(--sl-color-white);
    border: 1px solid var(--sl-color-gray-5);
    transition: all 0.2s;
  }
  
  .small-buttons a:hover {
    background: var(--sl-color-gray-5);
    border-color: var(--sl-color-accent);
  }
  
  /* Round card corners and add translucent blur effect */
  .mpress-card {
    border-radius: 12px;
    background: rgba(var(--sl-color-gray-6-rgb), 0.6) !important;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
  }

</style>

<span class="browser-footnote" data-browser-footnote>* 정말 원한다면 사용할 수 있습니다</span>
<script>
(() => {
  const words = ['Desktop', 'Server', 'Mobile'];
  let index = 0;
  const init = () => {
    const tagline = document.querySelector('.mpress-frontmatter-hero-tagline');
    const note = document.querySelector('[data-browser-footnote]');
    if (tagline && note && !tagline.contains(note)) tagline.appendChild(note);
    const title = document.querySelector('.mpress-frontmatter-hero-title');
    if (!title || title.querySelector('.morph-word-title')) return;
    if (title.textContent.trim() !== 'Build Desktop Apps with Go') return;
    const titleLine = document.createElement('span');
    titleLine.className = 'morph-title-line';
    const prefix = document.createElement('span');
    prefix.textContent = 'Build';
    const wordWrap = document.createElement('span');
    wordWrap.className = 'morph-word-title';
    const activeWord = document.createElement('span');
    activeWord.textContent = words[0];
    wordWrap.append(activeWord);
    titleLine.append(prefix, wordWrap);
    title.replaceChildren(titleLine, document.createElement('br'), document.createTextNode('Apps with Go'));
    const rotate = () => {
      const word = title.querySelector('.morph-word-title span');
      if (!word) return;
      word.classList.add('morph-word-out');
      setTimeout(() => {
        index = (index + 1) % words.length;
        word.textContent = words[index];
        word.classList.remove('morph-word-out');
        setTimeout(rotate, 3600);
      }, 650);
    };
    setTimeout(rotate, 3600);
  };
  document.readyState === 'loading' ? document.addEventListener('DOMContentLoaded', init) : init();
})();
</script>

## 빠른 시작

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

이제 첫 번째 프로젝트를 시작할 준비가 되었습니다.

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**이제 애플리케이션이 실행 중이며** 핫 리로드와 타입 안전성이 보장되는 Go-to-JS 바인딩을 사용할 수 있습니다.

`wails3 setup`에 문제가 있나요? [수동 설치 가이드](/quick-start/installation/)를 참조하세요.

## 데스크톱 앱을 그대로 모바일 앱으로

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Go 코드를 전혀 변경할 필요가 없습니다.** 동일한 `main.go`, 동일한 서비스, 동일한 프런트엔드를 Wails가 iOS 및 Android용으로 자동 컴파일합니다. 플랫폼별 기능(햅틱, 네이티브 대화 상자, 안전 영역)은 필요할 때 사용할 수 있지만, 앱을 출시하는 데 반드시 필요한 것은 아닙니다.

[모바일 문서 →](/guides/mobile/)

---

## Wails를 선택해야 하는 이유

@cards{cols="2"}
🚀 사용자가 체감하는 성능
- Electron의 150MB에 비해 약 15MB인 바이너리
- 100MB+에 비해 약 10MB인 기본 메모리 사용량
- 2-3초에 비해 &lt;0.5초인 시작 시간
- OS WebView를 사용한 네이티브 렌더링
- 번들 브라우저로 인한 오버헤드 없음

---
⚙ 개발자 경험
- 모든 플랫폼에서 하나의 Go 코드베이스 사용
- React, Vue, Svelte 등 모든 웹 프레임워크 지원
- 개발 중 핫 리로드 지원
- Go를 Javascript에서 쉽게 호출할 수 있도록 바인딩 자동 생성
- 메모리 내 IPC. 네트워크 포트를 사용하지 않음

---
✓ 데스크톱 앱 개발에 필요한 기능 완비
- 수명 주기를 지원하는 다중 창
- 네이티브 메뉴 및 시스템 트레이
- 플랫폼 네이티브 파일 대화 상자
- 시스템 통합 및 바로 가기
- 코드 서명 및 패키징 도구

---
▣ 데스크톱 및 모바일
- Windows, macOS, Linux, iOS, Android
- 동일한 코드베이스, 재작성 불필요
- 모든 플랫폼에서 네이티브 WebView 사용
- 열린 포트 및 localhost 서버 없음
- 필요할 때 플랫폼 기능 사용 가능

@end

## 다음 단계

다음으로 [완전한 애플리케이션을 빌드](/tutorials/03-notes-vanilla/)하거나, [예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 살펴보거나, [API 레퍼런스](/reference/application/)를 확인하세요. v2에서 마이그레이션하시나요? [업그레이드 가이드](/migration/v2-to-v3/)를 참조하세요.

@note{type="info" title="Wails v3 베타"}
Wails v3는 안정적인 데스크톱 API를 제공하는 베타 소프트웨어입니다. 이미 여러 팀이 프로덕션 환경에서 사용하고 있지만, 3.0을 위한 최종 마무리 작업이 진행 중이므로 배포하기 전에 철저히 테스트하는 것이 좋습니다.

@end

@cards{cols="1"}
♥ Wails 개발 후원하기
Wails는 개발자가 개발자를 위해 만든 무료 오픈 소스입니다. Wails로 훌륭한 애플리케이션을 빌드하는 데 도움을 받으셨다면 지속적인 개발을 후원해 주세요.

여러분의 후원은 프로젝트를 유지하고 문서를 개선하며 전체 커뮤니티에 도움이 되는 새로운 기능을 개발하는 데 보탬이 됩니다.

[후원자 되기 →](https://github.com/sponsors/leaanthony)

@end

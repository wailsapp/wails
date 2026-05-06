#!/usr/bin/env python3
"""
Wails v3 documentation translator using Ollama streaming API.
Usage:
  python3 translate-docs.py --locale all
  python3 translate-docs.py --locale zh-cn
  python3 translate-docs.py --locale ja,ko
"""
import os
import sys
import json
import argparse
import hashlib
import re
import time
from pathlib import Path

try:
    import requests
except ImportError:
    print("Installing requests...")
    os.system("pip3 install requests -q")
    import requests

OLLAMA_BASE = "http://ai-master.taileaa27f.ts.net:11434"
MODEL = "qwen3.6:35b"

ZAI_BASE = os.environ.get("ZAI_BASE", "https://api.z.ai/api/coding/paas/v4")
ZAI_API_KEY = os.environ.get("ZAI_API_KEY", "")
ZAI_MODEL = os.environ.get("ZAI_MODEL", "glm-5.1")

# Active backend — overridden by --zai-model CLI flag
TRANSLATE_BACKEND = "ollama"

DOCS_ROOT = Path(__file__).parent.parent / "src" / "content" / "docs"
CACHE_DIR = Path(__file__).parent.parent / ".translation-cache"

LOCALE_NAMES = {
    "zh-cn": "Simplified Chinese (Mandarin)",
    "zh-tw": "Traditional Chinese (Taiwan)",
    "ja": "Japanese",
    "ko": "Korean",
    "ru": "Russian",
    "fr": "French",
    "pt": "Portuguese (Brazilian)",
    "de": "German",
}

ALL_LOCALES = list(LOCALE_NAMES.keys())

# Priority files to translate - shorter files first for broader coverage
PRIORITY_FILES = [
    "index.mdx",
    "quick-start/why-wails.mdx",
    "quick-start/next-steps.mdx",          # 34 lines
    "status.mdx",                           # 33 lines
    "getting-started/installation.mdx",    # 114 lines
    "feedback.mdx",                        # 87 lines
    "credits.mdx",                         # 63 lines
    "community/links.md",
    "community/templates.md",
    "faq.mdx",                             # 260 lines
    "quick-start/first-app.mdx",           # 270 lines
    "getting-started/your-first-app.mdx",  # 273 lines
    "quick-start/installation.mdx",        # 430 lines - large
    "concepts/architecture.mdx",
    "concepts/bridge.mdx",
    "concepts/lifecycle.mdx",
    "concepts/manager-api.mdx",
    "concepts/build-system.mdx",
    "contributing.mdx",
]

# Paths within the docs that have locale-specific translations.
# When a file is translated, its path is added here so that internal
# links in OTHER translated files can be rewritten to the locale-prefixed version.
# Keys are source-relative paths (no leading /), values are the locale-relative path
# (same value, since translated files mirror the source structure).
# This set is updated dynamically as files are translated during a run.
ALWAYS_LOCALIZE_PATHS = {
    # These core pages are always translated first; links to them should always
    # be locale-prefixed. Updated at runtime via build_translated_paths().
}

# Paths in source docs that map to a DIFFERENT target path in locale docs.
# e.g. the primary installation CTA links to quick-start/installation but we
# translate getting-started/installation, so the link should be rewritten.
PATH_REMAP = {
    "/quick-start/installation": "/getting-started/installation",
}

SYSTEM_PROMPT = """You are a professional technical documentation translator.
Translate the given MDX/Markdown documentation accurately, completely, and naturally in {lang}.

CRITICAL INSTRUCTION: Ensure the translation covers the ENTIRE source text. Do not truncate output.
Every heading, paragraph, list item, and sentence in the source must appear in the translation.

STRICT RULES — violating any of these will cause the page to break or render incorrectly:

1. TRANSLATION QUALITY:
   - Translate ALL prose text naturally in {lang}. Do not be overly literal.
   - Maintain consistent terminology for technical concepts throughout the document.
   - PRESERVE ALL product names, proper nouns, and brand names exactly as written in the source.
     These must NEVER be translated, transliterated, altered, or "corrected":
     Wails, Electron, Go, npm, Xcode, WebView2, React, Vue, Svelte, macOS, Windows, Linux,
     Discord, TypeScript, JavaScript — and any other brand names or tool names.
   - Keep technical loanwords in their original English form unless a well-established {lang}
     translation exists (e.g. "hot reload", "WebView", "binary" may stay as-is).
   - Software release terms: "Ship" and "shipping" (in the context of releasing software)
     should be translated as the equivalent of "release/publish" (e.g. German: "Veröffentlichen",
     French: "Publier", Japanese: "リリース"), NOT as "deliver/send" (e.g. NOT "Liefern").

2. PRESERVE frontmatter structure exactly:
   - Keep all YAML keys in English (title:, description:, link:, icon:, etc.)
   - Translate ONLY the string values of: title, description, tagline, text, label, alt, content (banner content)
   - Copy all other frontmatter values (links, icons, variants, booleans) UNCHANGED

3. CODE BLOCKS — COPY VERBATIM, CHARACTER FOR CHARACTER:
   Everything between ``` fences is CODE and must be copied without any change whatsoever.
   This rule has NO exceptions:
   - Shell commands, Go code, configuration files — do NOT translate
   - Code comments (// ..., # ..., /* ... */) — do NOT translate
     WRONG: `# Wails installieren`  CORRECT: `# Install Wails`
     Comments inside code blocks are part of the code — never translate them.
   - d2 diagram definitions — the ENTIRE block including quoted string labels like
     "Your UI\\n(React/Vue/etc)" must be copied exactly.
     WRONG: "Ihre UI\\n(React/Vue/etc)"  CORRECT: "Your UI\\n(React/Vue/etc)"
     String labels inside d2 blocks are diagram code, not prose — never translate them.
   - JSX/MDX component names and props (e.g. <Tabs>, <TabItem label="...">) — do NOT translate
   - import statements — do NOT translate
   - Inline code (`backtick content`) — copy exactly, including CLI commands and paths

4. URLs AND PATHS — copy unchanged:
   - URLs (http://, https://) — copy exactly
   - File paths and directory names — copy exactly

5. FRONTMATTER DELIMITERS — CRITICAL:
   - The document MUST begin with exactly one line containing only: ---
   - The frontmatter block MUST end with exactly one line containing only: ---
   - Do NOT output two --- lines at the start (wrong: ---\\n---\\ntitle:)
   - Do NOT add any --- at the end of the document body

6. Return ONLY the translated document. No preamble, no explanation, no markdown fences around the whole document.

7. Preserve all blank lines, heading levels (#, ##, ###), list structure (-, *), indentation,
   line breaks, and MDX component structure exactly as in the source."""

CHUNK_SIZE = 8000  # chars; files larger than this are translated in chunks

SUMMARIZE_PROMPT = """The following is a translated excerpt from Wails v3 documentation.
Write a 2-3 sentence context note for the translator handling the next excerpt:
- What topics/sections were covered
- Key terminology choices (how specific technical terms were rendered in {lang})
Keep it concise — it will be injected into the next translation prompt.

---
{translated}
---"""

CHUNK_TRANSLATE_PROMPT = """Translate the following {chunk_label} of a Wails v3 documentation page to {lang}.
Apply all rules strictly. Return only the translated text.
{context_block}
{content}"""


def file_hash(content: str) -> str:
    return hashlib.md5(content.encode()).hexdigest()[:12]


def split_into_chunks(content: str, max_size: int = CHUNK_SIZE) -> list[str]:
    """Split MDX content at paragraph/heading boundaries, never inside code blocks."""
    if len(content) <= max_size:
        return [content]

    chunks = []
    start = 0

    while start < len(content):
        if len(content) - start <= max_size:
            chunks.append(content[start:])
            break

        lo = start + max_size // 2
        hi = min(len(content), start + max_size)
        window = content[lo:hi]

        # Rough check: odd number of ``` lines before lo means we're inside a code block
        pre = content[start:lo]
        in_code = pre.count('\n```') % 2 == 1

        split_pos = None
        if not in_code:
            # Prefer last heading boundary (\n\n## ...) in window
            for m in reversed(list(re.finditer(r'\n\n(?=#+\s)', window))):
                split_pos = lo + m.end()
                break

            if split_pos is None:
                # Fall back to last paragraph break
                last = window.rfind('\n\n')
                if last >= 0:
                    split_pos = lo + last + 2

        if split_pos is None:
            # Last resort: last newline before hi
            last_nl = content.rfind('\n', start, hi)
            split_pos = (last_nl + 1) if last_nl > start else hi

        chunks.append(content[start:split_pos])
        start = split_pos

    return chunks


def strip_frontmatter(content: str) -> str:
    """Remove a leading frontmatter block if the model added one to a continuation chunk."""
    content = content.lstrip("\n")
    if content.startswith("---"):
        end = content.find("\n---", 3)
        if end >= 0:
            return content[end + 4:].lstrip("\n")
    return content


def summarize_chunk(translated: str, lang: str) -> str:
    """Ask the model for a brief context summary of a translated chunk."""
    prompt = SUMMARIZE_PROMPT.format(translated=translated[:4000], lang=lang)
    payload = {
        "model": MODEL,
        "messages": [{"role": "user", "content": prompt}],
        "stream": False,
        "think": False,
        "options": {"temperature": 0.1, "num_predict": 300},
    }
    try:
        resp = requests.post(f"{OLLAMA_BASE}/api/chat", json=payload, timeout=(10, 60))
        resp.raise_for_status()
        return resp.json().get("message", {}).get("content", "").strip()
    except Exception:
        return ""


def load_cache(locale: str) -> dict:
    cache_file = CACHE_DIR / f"{locale}.json"
    if cache_file.exists():
        try:
            return json.loads(cache_file.read_text())
        except Exception:
            return {}
    return {}


def save_cache(locale: str, cache: dict):
    CACHE_DIR.mkdir(exist_ok=True)
    cache_file = CACHE_DIR / f"{locale}.json"
    cache_file.write_text(json.dumps(cache, indent=2, ensure_ascii=False))


def build_translated_paths(locale: str, files: list) -> set:
    """
    Return the set of root-relative paths (e.g. '/quick-start/next-steps')
    that have translated versions for this locale, based on the files list.
    This drives automatic link rewriting in post_process().
    """
    paths = set()
    for src_path in files:
        rel = src_path.relative_to(DOCS_ROOT)
        # Strip extension to get the URL path
        url_path = "/" + str(rel.with_suffix(""))
        # Also check if the locale file actually exists (cached from prior run)
        out_path = DOCS_ROOT / locale / rel
        if out_path.exists() or True:  # optimistic: assume we'll translate it
            paths.add(url_path)
    return paths


def post_process(content: str, locale: str, rel_path: Path, translated_paths: set) -> str:
    """
    Apply deterministic post-processing to every translated file:

    1. Fix malformed double---- frontmatter (model artifact)
    2. Remove trailing --- artifact (model artifact)
    3. Fix relative asset paths (locale files are one dir deeper than source root)
    4. Rewrite internal links for translated pages to locale-prefixed versions
    5. Remap CTA paths that point to a different-named translated file
    """
    # 1. Fix double---- frontmatter: model sometimes outputs "---\n---\ntitle:"
    if content.startswith("---\n---\n"):
        content = content[4:]  # strip the spurious leading "---\n"

    # 2. Remove trailing --- artifact
    stripped = content.rstrip()
    if stripped.endswith("\n---"):
        content = stripped[:-4].rstrip() + "\n"

    # 3. Fix relative asset paths
    #    Source root: docs/src/content/docs/index.mdx → ../../assets = docs/src/assets ✓
    #    Locale root: docs/src/content/docs/de/index.mdx → ../../assets = docs/src/content/assets ✗
    #    Fix: ../../assets/ → ../../../assets/
    content = content.replace("../../assets/", "../../../assets/")

    # 4 & 5. Link rewriting
    #
    # For each translated path, rewrite:
    #   - Markdown links: [text](/path) → [text](/<locale>/path)
    #   - Frontmatter link: values: link: /path → link: /<locale>/path
    #   - href="/path" in JSX
    #
    # Also apply PATH_REMAP before prefixing:
    #   /quick-start/installation → /<locale>/getting-started/installation
    #   (because we translate getting-started/installation, not quick-start/installation)
    #
    # ONLY rewrite paths that have a translated version. Leave other root-relative
    # links as-is so they fall back to the English page via Starlight locale fallback.

    def make_locale_path(src_path: str) -> str | None:
        """Return locale-prefixed path if src_path should be rewritten, else None."""
        # Apply remap first
        remapped = PATH_REMAP.get(src_path, src_path)
        # Check if the remapped path (without extension) is in translated_paths
        # translated_paths contains paths like /quick-start/next-steps (no extension)
        if remapped in translated_paths:
            return f"/{locale}{remapped}"
        return None

    def rewrite_markdown_link(m: re.Match) -> str:
        text, path, suffix = m.group(1), m.group(2), m.group(3)
        new_path = make_locale_path(path)
        if new_path:
            return f"[{text}]({new_path}{suffix if suffix != ')' else ''})"
        return m.group(0)

    # Markdown links: [text](/path) or [text](/path/sub)
    content = re.sub(
        r'\[([^\]]*)\]\((/[a-zA-Z0-9/_-]+)(/?)\)',
        rewrite_markdown_link,
        content
    )

    # Frontmatter link: /path
    def rewrite_fm_link(m: re.Match) -> str:
        key, path, tail = m.group(1), m.group(2), m.group(3)
        new_path = make_locale_path(path)
        if new_path:
            return f"{key}{new_path}{tail}"
        return m.group(0)

    content = re.sub(
        r'(link:\s*)(/[a-zA-Z0-9/_-]+)(\s*$)',
        rewrite_fm_link,
        content,
        flags=re.MULTILINE
    )

    # href="/path" in JSX
    def rewrite_href(m: re.Match) -> str:
        path = m.group(1)
        new_path = make_locale_path(path)
        if new_path:
            return f'href="{new_path}"'
        return m.group(0)

    content = re.sub(r'href="(/[a-zA-Z0-9/_-]+)"', rewrite_href, content)

    return content


def translate_with_ollama(content: str, lang: str, file_path: str,
                          context_summary: str = "", is_continuation: bool = False) -> str:
    """Translate content using Ollama chat API (thinking disabled)."""
    system = SYSTEM_PROMPT.format(lang=lang)

    context_block = ""
    if context_summary:
        context_block = (
            "Previous document context (already translated — maintain consistent terminology):\n"
            f"{context_summary}\n"
        )

    chunk_label = "continuation" if is_continuation else "page"
    user_msg = CHUNK_TRANSLATE_PROMPT.format(
        chunk_label=chunk_label,
        lang=lang,
        context_block=context_block,
        content=content,
    )

    payload = {
        "model": MODEL,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user_msg},
        ],
        "stream": True,
        "think": False,
        "options": {
            "temperature": 0.2,
            "top_p": 0.9,
            "num_predict": 16384,
        }
    }

    response_parts = []
    try:
        resp = requests.post(
            f"{OLLAMA_BASE}/api/chat",
            json=payload,
            stream=True,
            timeout=(10, 600)
        )
        resp.raise_for_status()

        for line in resp.iter_lines():
            if not line:
                continue
            try:
                data = json.loads(line)
            except json.JSONDecodeError:
                continue

            msg = data.get("message", {})
            token = msg.get("content", "")
            if token:
                response_parts.append(token)

            if data.get("done"):
                break

    except requests.exceptions.Timeout:
        print(f"  TIMEOUT translating {file_path}", flush=True)
        return None
    except Exception as e:
        print(f"  ERROR translating {file_path}: {e}", flush=True)
        return None

    result = "".join(response_parts).strip()

    # Strip any leading/trailing markdown code fences the model might add
    if result.startswith("```") and result.endswith("```"):
        lines = result.split("\n")
        result = "\n".join(lines[1:-1])

    # Strip any leading "---" explanatory text before the frontmatter
    if not result.startswith("---") and "---" in result:
        idx = result.index("---")
        result = result[idx:]

    return result if result else None


def translate_with_zai(content: str, lang: str, file_path: str,
                       context_summary: str = "", is_continuation: bool = False) -> str | None:
    """Translate content using z.ai OpenAI-compatible API."""
    if not ZAI_API_KEY:
        print("  ERROR: ZAI_API_KEY not set", flush=True)
        return None

    system = SYSTEM_PROMPT.format(lang=lang)
    context_block = ""
    if context_summary:
        context_block = (
            "Previous document context (already translated — maintain consistent terminology):\n"
            f"{context_summary}\n"
        )
    chunk_label = "continuation" if is_continuation else "page"
    user_msg = CHUNK_TRANSLATE_PROMPT.format(
        chunk_label=chunk_label,
        lang=lang,
        context_block=context_block,
        content=content,
    )

    try:
        resp = requests.post(
            f"{ZAI_BASE}/chat/completions",
            headers={
                "Authorization": f"Bearer {ZAI_API_KEY}",
                "Content-Type": "application/json",
            },
            json={
                "model": ZAI_MODEL,
                "messages": [
                    {"role": "system", "content": system},
                    {"role": "user", "content": user_msg},
                ],
                "temperature": 0.2,
                "max_tokens": 16000,
            },
            timeout=360,
        )
        resp.raise_for_status()
        data = resp.json()
    except requests.exceptions.Timeout:
        print(f"  TIMEOUT (z.ai) translating {file_path}", flush=True)
        return None
    except Exception as e:
        print(f"  ERROR (z.ai) translating {file_path}: {e}", flush=True)
        return None

    choice = data.get("choices", [{}])[0]
    result = choice.get("message", {}).get("content", "").strip()

    if not result:
        usage = data.get("usage", {})
        print(f"  WARNING: empty content from {ZAI_MODEL} (usage={usage}) for {file_path}", flush=True)
        return None

    # Strip any wrapping markdown fences
    if result.startswith("```") and result.endswith("```"):
        lines = result.split("\n")
        result = "\n".join(lines[1:-1])

    # Strip any preamble before frontmatter
    if not result.startswith("---") and "---" in result:
        idx = result.index("---")
        result = result[idx:]

    return result if result else None


def _translate_chunk(content: str, lang: str, file_path: str,
                     context_summary: str = "", is_continuation: bool = False) -> str | None:
    """Dispatch to the active translation backend."""
    if TRANSLATE_BACKEND == "zai":
        return translate_with_zai(content, lang, file_path, context_summary, is_continuation)
    return translate_with_ollama(content, lang, file_path, context_summary, is_continuation)


def translate_doc(content: str, lang: str, file_path: str) -> str | None:
    """Translate a document, splitting into chunks with context passing for large files."""
    chunks = split_into_chunks(content)

    if len(chunks) == 1:
        return _translate_chunk(content, lang, file_path)

    print(f"    Splitting into {len(chunks)} chunks", flush=True)
    parts = []
    context = ""

    for i, chunk in enumerate(chunks):
        is_last = i == len(chunks) - 1
        translated_chunk = _translate_chunk(
            chunk, lang, file_path,
            context_summary=context,
            is_continuation=(i > 0),
        )
        if translated_chunk is None:
            return None

        # Strip any spurious frontmatter the model adds to continuation chunks
        if i > 0:
            translated_chunk = strip_frontmatter(translated_chunk)

        parts.append(translated_chunk)

        if not is_last:
            context = summarize_chunk(translated_chunk, lang)
            preview = context[:80].replace('\n', ' ')
            print(f"    Chunk {i+1}/{len(chunks)} done → context: {preview}...", flush=True)
        else:
            print(f"    Chunk {i+1}/{len(chunks)} done", flush=True)

    return "".join(parts)


def translate_file(
    src_path: Path,
    locale: str,
    lang: str,
    cache: dict,
    translated_paths: set,
    force: bool = False,
) -> tuple[bool, str]:
    """
    Translate a single file and apply post-processing.
    Returns (success, status) where status is 'translated', 'cached', or 'failed'.
    """
    content = src_path.read_text(encoding="utf-8")
    h = file_hash(content)
    cache_key = str(src_path.relative_to(DOCS_ROOT))

    # Check cache (skipped when force=True)
    if not force and cache_key in cache and cache[cache_key].get("hash") == h:
        rel = src_path.relative_to(DOCS_ROOT)
        out_path = DOCS_ROOT / locale / rel
        if out_path.exists():
            return True, "cached"

    print(f"  Translating {src_path.relative_to(DOCS_ROOT)} → {locale}...", flush=True)

    translated = translate_doc(content, lang, str(src_path))
    if not translated:
        return False, "failed"

    rel = src_path.relative_to(DOCS_ROOT)
    out_path = DOCS_ROOT / locale / rel
    out_path.parent.mkdir(parents=True, exist_ok=True)

    # Apply all post-processing fixes
    translated = post_process(translated, locale, rel, translated_paths)

    out_path.write_text(translated, encoding="utf-8")

    cache[cache_key] = {"hash": h, "translated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())}
    return True, "translated"


def get_source_files() -> list[Path]:
    """Get priority docs files that exist."""
    files = []
    for rel in PRIORITY_FILES:
        p = DOCS_ROOT / rel
        if p.exists():
            files.append(p)

    # Add any remaining files not in priority list
    all_files = set(
        p for p in DOCS_ROOT.rglob("*.mdx")
        if not any(part in str(p) for part in ["blog", "changelog", "whats-new", "showcase"])
        and p.parts[len(DOCS_ROOT.parts)] not in ALL_LOCALES
    ) | set(
        p for p in DOCS_ROOT.rglob("*.md")
        if not any(part in str(p) for part in ["blog", "changelog", "whats-new", "showcase"])
        and p.parts[len(DOCS_ROOT.parts)] not in ALL_LOCALES
    )

    # Add non-priority files at the end
    for p in sorted(all_files):
        if p not in files:
            files.append(p)

    return files


def run_locale(locale: str, files: list[Path], max_files: int = None, force: bool = False):
    lang = LOCALE_NAMES[locale]
    cache = load_cache(locale)

    # Build the set of paths that will be translated in this run (for link rewriting)
    files_to_process = files[:max_files] if max_files else files
    translated_paths = build_translated_paths(locale, files_to_process)

    print(f"\n{'='*60}", flush=True)
    print(f"Translating to {locale} ({lang})", flush=True)
    print(f"{'='*60}", flush=True)

    translated = 0
    cached = 0
    failed = 0
    count = 0

    for src_path in files:
        if max_files and count >= max_files:
            break
        count += 1

        success, status = translate_file(src_path, locale, lang, cache, translated_paths, force=force)
        if success:
            if status == "translated":
                translated += 1
                print(f"  ✓ {src_path.relative_to(DOCS_ROOT)}", flush=True)
            else:
                cached += 1
                print(f"  (cached) {src_path.relative_to(DOCS_ROOT)}", flush=True)
        else:
            failed += 1
            print(f"  ✗ FAILED: {src_path.relative_to(DOCS_ROOT)}", flush=True)

        # Save cache after each file
        save_cache(locale, cache)

    print(f"\nLocale {locale}: {translated} translated, {cached} cached, {failed} failed", flush=True)
    return {"locale": locale, "translated": translated, "cached": cached, "failed": failed}


def main():
    global TRANSLATE_BACKEND, ZAI_MODEL, MODEL

    parser = argparse.ArgumentParser(description="Translate Wails v3 docs")
    parser.add_argument("--locale", default="all", help="Locale(s) to translate (comma-separated or 'all')")
    parser.add_argument("--max-files", type=int, default=None, help="Max files per locale")
    parser.add_argument("--zai-model", default=None,
                        help="Use z.ai instead of Ollama; specify model name (e.g. glm-5.1, glm-4.7)")
    parser.add_argument("--ollama-model", default=None,
                        help="Override Ollama model (e.g. hf.co/Jackrong/Qwen3.5-9B-GLM5.1-Distill-v1-GGUF:Q4_K_M)")
    parser.add_argument("--force", action="store_true",
                        help="Re-translate even if cached (for model comparison)")
    args = parser.parse_args()

    if args.zai_model:
        TRANSLATE_BACKEND = "zai"
        ZAI_MODEL = args.zai_model
        print(f"Backend: z.ai ({ZAI_MODEL})", flush=True)
    else:
        if args.ollama_model:
            MODEL = args.ollama_model
        print(f"Backend: Ollama ({MODEL})", flush=True)

    if args.locale == "all":
        locales = ALL_LOCALES
    else:
        locales = [l.strip() for l in args.locale.split(",")]
        for l in locales:
            if l not in LOCALE_NAMES:
                print(f"Unknown locale: {l}. Valid: {', '.join(ALL_LOCALES)}")
                sys.exit(1)

    files = get_source_files()
    print(f"Found {len(files)} source files to translate", flush=True)
    print(f"Locales: {', '.join(locales)}", flush=True)
    if args.force:
        print("Force mode: bypassing cache", flush=True)

    results = []
    for locale in locales:
        result = run_locale(locale, files, args.max_files, force=args.force)
        results.append(result)

    print("\n" + "="*60, flush=True)
    print("TRANSLATION SUMMARY", flush=True)
    print("="*60, flush=True)
    print(f"{'Locale':<10} {'Translated':>12} {'Cached':>8} {'Failed':>8}", flush=True)
    for r in results:
        print(f"{r['locale']:<10} {r['translated']:>12} {r['cached']:>8} {r['failed']:>8}", flush=True)

    total_files = sum(r['translated'] + r['cached'] for r in results)
    print(f"\nTotal files processed: {total_files}", flush=True)

    # Ensure all translated locales are registered in astro.config.mjs
    print("\nUpdating astro.config.mjs locale registrations...", flush=True)
    ensure_astro_locales(locales)


ASTRO_CONFIG = Path(__file__).parent.parent / "astro.config.mjs"

# Locale metadata for astro.config.mjs registration
LOCALE_ASTRO = {
    "zh-cn": ('"zh-cn"', "简体中文",          "zh-CN"),
    "zh-tw": ('"zh-tw"', "繁體中文",          "zh-TW"),
    "ja":    ("ja",      "日本語",             "ja"),
    "ko":    ("ko",      "한국어",             "ko"),
    "ru":    ("ru",      "Русский",           "ru"),
    "fr":    ("fr",      "Français",          "fr"),
    "pt":    ("pt",      "Português (Brasil)", "pt-BR"),
    "de":    ("de",      "Deutsch",           "de"),
}


def ensure_astro_locales(locales: list[str]):
    """Add any missing locales to the Starlight locales block in astro.config.mjs."""
    if not ASTRO_CONFIG.exists():
        return
    config = ASTRO_CONFIG.read_text(encoding="utf-8")
    changed = False
    for locale in locales:
        meta = LOCALE_ASTRO.get(locale)
        if not meta:
            continue
        key, label, lang = meta
        # Check if locale key is already present
        if f"{key}:" in config or f"{key} :" in config:
            continue
        # Insert before the closing brace of the locales block
        entry = f'        {key}: {{ label: "{label}", lang: "{lang}", dir: "ltr" }},\n'
        config = config.replace(
            '        root: { label: "English", lang: "en", dir: "ltr" },\n      },',
            f'        root: {{ label: "English", lang: "en", dir: "ltr" }},\n{entry}      }},'
        )
        # If the above pattern already matched and was replaced in a previous iteration,
        # look for the locale block closing differently
        if entry not in config:
            # Find the locales closing brace and insert before it
            marker = "      },\n      plugins:"
            if marker in config and entry not in config:
                config = config.replace(marker, f"{entry}{marker}", 1)
        changed = True
        print(f"  Registered locale {locale} ({label}) in astro.config.mjs", flush=True)
    if changed:
        ASTRO_CONFIG.write_text(config, encoding="utf-8")


if __name__ == "__main__":
    main()

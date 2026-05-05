#!/usr/bin/env python3
"""
Wails v3 translation QA scorer.
Usage:
  python3 translate-qa.py --locale all
  python3 translate-qa.py --locale zh-cn
  python3 translate-qa.py --locale zh-cn --ai-verify
  python3 translate-qa.py --locale all --ai-verify --ai-model z1-mini
"""
import os
import sys
import json
import argparse
import re
from pathlib import Path

try:
    import requests
except ImportError:
    os.system("pip3 install requests -q")
    import requests

DOCS_ROOT = Path(__file__).parent.parent / "src" / "content" / "docs"
QA_DIR = Path(__file__).parent.parent / ".translation-qa"

LOCALE_NAMES = {
    "zh-cn": "Simplified Chinese",
    "zh-tw": "Traditional Chinese",
    "ja": "Japanese",
    "ko": "Korean",
    "ru": "Russian",
    "fr": "French",
    "pt": "Portuguese",
    "de": "German",
}

ALL_LOCALES = list(LOCALE_NAMES.keys())

# Characters expected in each locale (basic sanity check)
LOCALE_CHAR_RANGES = {
    "zh-cn": (r'[一-鿿]', "Chinese characters"),
    "zh-tw": (r'[一-鿿]', "Chinese characters"),
    "ja": (r'[぀-ヿ一-鿿]', "Japanese characters"),
    "ko": (r'[가-힯]', "Korean characters"),
    "ru": (r'[Ѐ-ӿ]', "Cyrillic characters"),
    "fr": (r'[a-zA-ZÀ-ÿ]', None),  # Latin, harder to distinguish
    "pt": (r'[a-zA-ZÀ-ÿ]', None),
    "de": (r'[a-zA-ZÄÖÜäöüß]', None),
}

# z.ai Coding Plan API — dedicated coding endpoint per https://docs.z.ai/devpack/overview
# GLM Coding Plan requires https://api.z.ai/api/coding/paas/v4 (NOT the general paas/v4)
ZAI_BASE = os.environ.get("ZAI_BASE", "https://api.z.ai/api/coding/paas/v4")
ZAI_API_KEY = os.environ.get("ZAI_API_KEY", "")
ZAI_DEFAULT_MODEL = "GLM-5-Turbo"

# Max characters of body text to send to AI verifier (keeps cost low)
AI_VERIFY_MAX_CHARS = 3000


def extract_frontmatter(content: str) -> tuple[str, str]:
    """Split content into frontmatter and body."""
    if not content.startswith("---"):
        return "", content
    end = content.find("---", 3)
    if end == -1:
        return "", content
    return content[3:end].strip(), content[end+3:].strip()


def count_code_blocks(content: str) -> int:
    return len(re.findall(r'```[\s\S]*?```', content))


def count_inline_code(content: str) -> int:
    return len(re.findall(r'`[^`]+`', content))


def strip_code_blocks(text: str) -> str:
    """Remove code blocks from text (for AI sampling — we don't want it scoring code)."""
    text = re.sub(r'```[\s\S]*?```', '', text)
    text = re.sub(r'`[^`\n]+`', '', text)
    return text.strip()


def ai_verify_translation(src_body: str, tgt_body: str, locale: str, lang_name: str, model: str) -> dict:
    """
    Call z.ai to verify translation quality.
    Returns {"score": float|None, "issues": list[str], "ai_used": bool}.
    Score is None if the API call failed or key is missing.
    """
    if not ZAI_API_KEY:
        return {"score": None, "issues": ["ZAI_API_KEY not set — skipping AI verification"], "ai_used": False}

    # Strip code blocks so the AI focuses on prose quality
    src_sample = strip_code_blocks(src_body)[:AI_VERIFY_MAX_CHARS]
    tgt_sample = strip_code_blocks(tgt_body)[:AI_VERIFY_MAX_CHARS]

    if len(src_sample) < 100 or len(tgt_sample) < 100:
        return {"score": None, "issues": ["Not enough prose to AI-verify"], "ai_used": False}

    prompt = f"""You are a professional translation quality assessor. Rate this {lang_name} translation of English technical documentation.

Score 0.0–1.0 based on:
- Accuracy: same meaning conveyed?
- Fluency: reads naturally in {lang_name}?
- Completeness: nothing missing or spuriously added?
- Technical terms: product names, UI labels, and technical jargon handled correctly?

English source (prose excerpt, code removed):
---
{src_sample}
---

{lang_name} translation (prose excerpt, code removed):
---
{tgt_sample}
---

Reply with JSON only — no prose, no markdown fences:
{{"score": <0.0-1.0>, "issues": ["specific issue 1", "specific issue 2"]}}

If the translation is good, return an empty issues list."""

    try:
        resp = requests.post(
            f"{ZAI_BASE}/chat/completions",
            headers={
                "Authorization": f"Bearer {ZAI_API_KEY}",
                "Content-Type": "application/json",
            },
            json={
                "model": model,
                "messages": [{"role": "user", "content": prompt}],
                "temperature": 0.1,
                "max_tokens": 2000,  # GLM-5-Turbo uses ~300 reasoning tokens before output
            },
            timeout=45,
        )
        resp.raise_for_status()

        raw = resp.json()["choices"][0]["message"]["content"].strip()
        # Strip markdown fences if model wraps in ```json ... ```
        if raw.startswith("```"):
            lines = raw.split("\n")
            raw = "\n".join(lines[1:-1])
        result = json.loads(raw)
        ai_score = float(result.get("score", 0.5))
        ai_issues = [f"[AI] {i}" for i in result.get("issues", [])]
        return {"score": ai_score, "issues": ai_issues, "ai_used": True}

    except Exception as e:
        return {"score": None, "issues": [f"[AI] Verification error: {e}"], "ai_used": False}


def check_file_pair(src_path: Path, tgt_path: Path, locale: str,
                    ai_verify: bool = False, ai_model: str = ZAI_DEFAULT_MODEL) -> dict:
    """Score a translated file against its source."""
    src = src_path.read_text(encoding="utf-8")
    tgt = tgt_path.read_text(encoding="utf-8")

    issues = []
    score = 1.0

    # 1. Check file is not empty
    if len(tgt.strip()) < 50:
        return {"score": 0.0, "issues": ["File is empty or too short"]}

    # 2. Check file is not just a copy of source (not translated)
    if src.strip() == tgt.strip():
        return {"score": 0.0, "issues": ["File appears to be an untranslated copy of source"]}

    # 3. Check translated file has content (not much shorter than source)
    src_len = len(src)
    tgt_len = len(tgt)
    ratio = tgt_len / src_len if src_len > 0 else 0
    # CJK languages can be shorter, European languages similar length
    min_ratio = 0.3 if locale in ["zh-cn", "zh-tw", "ja", "ko"] else 0.5
    max_ratio = 3.0
    if ratio < min_ratio:
        issues.append(f"Translated file is suspiciously short (ratio: {ratio:.2f})")
        score -= 0.3
    elif ratio > max_ratio:
        issues.append(f"Translated file is suspiciously long (ratio: {ratio:.2f})")
        score -= 0.1

    # 4. Check frontmatter is preserved
    src_fm, src_body = extract_frontmatter(src)
    tgt_fm, tgt_body = extract_frontmatter(tgt)

    if src_fm and not tgt_fm:
        issues.append("Frontmatter missing in translation")
        score -= 0.2

    if src_fm and tgt_fm:
        # Check YAML keys are preserved
        src_keys = set(re.findall(r'^(\w+):', src_fm, re.MULTILINE))
        tgt_keys = set(re.findall(r'^(\w+):', tgt_fm, re.MULTILINE))
        missing_keys = src_keys - tgt_keys
        if missing_keys:
            issues.append(f"Missing frontmatter keys: {missing_keys}")
            score -= 0.15

    # 5. Check code blocks are preserved
    src_code_blocks = re.findall(r'```[\s\S]*?```', src)
    tgt_code_blocks = re.findall(r'```[\s\S]*?```', tgt)
    if len(src_code_blocks) != len(tgt_code_blocks):
        issues.append(f"Code block count mismatch: src={len(src_code_blocks)}, tgt={len(tgt_code_blocks)}")
        score -= 0.15

    # 6. Check code block contents are not translated
    for i, (sc, tc) in enumerate(zip(src_code_blocks, tgt_code_blocks)):
        # Extract the code content (strip language marker)
        sc_content = re.sub(r'^```\w*\n?', '', sc).rstrip('`').strip()
        tc_content = re.sub(r'^```\w*\n?', '', tc).rstrip('`').strip()
        if sc_content != tc_content:
            issues.append(f"Code block {i+1} content was modified")
            score -= 0.1
            break  # Only report once

    # 7. Check MDX imports are preserved
    src_imports = re.findall(r'^import\s+.+$', src, re.MULTILINE)
    tgt_imports = re.findall(r'^import\s+.+$', tgt, re.MULTILINE)
    if len(src_imports) != len(tgt_imports):
        issues.append(f"Import count mismatch: src={len(src_imports)}, tgt={len(tgt_imports)}")
        score -= 0.1

    # 8. Check locale-specific characters appear (for non-Latin locales)
    char_pattern, char_desc = LOCALE_CHAR_RANGES.get(locale, (None, None))
    if char_pattern and char_desc:
        if not re.search(char_pattern, tgt_body):
            issues.append(f"No {char_desc} found in translation body")
            score -= 0.4

    # 9. Check URLs are preserved
    src_urls = re.findall(r'https?://\S+', src)
    tgt_urls = re.findall(r'https?://\S+', tgt)
    if src_urls and len(tgt_urls) < len(src_urls) * 0.7:
        issues.append(f"URL count mismatch: src={len(src_urls)}, tgt={len(tgt_urls)}")
        score -= 0.1

    # 10. Check for malformed double-frontmatter artifact (model bug)
    if tgt.startswith("---\n---\n"):
        issues.append("Malformed double frontmatter (---\\n---\\n) — post-processing should have fixed this")
        score -= 0.3

    # 11. Check for trailing --- artifact (model bug)
    if tgt.rstrip().endswith("\n---"):
        issues.append("Trailing --- artifact at end of file — post-processing should have fixed this")
        score -= 0.1

    # 12. Check relative asset paths are correct for locale depth
    #     Source uses ../../assets/; locale files should use ../../../assets/
    #     Use string replacement to avoid false positive: ../../../assets/ contains ../../assets/ as substring
    tgt_without_good_path = tgt.replace('../../../assets/', '')
    if '../../assets/' in tgt_without_good_path:
        issues.append("Incorrect relative asset path ../../assets/ (should be ../../../assets/ for locale files)")
        score -= 0.2

    # 13. Check that internal links to translated pages use locale prefix
    #     Heuristic: if the source has a /quick-start/installation link and the
    #     target still has that exact link without a locale prefix, flag it.
    src_internal = set(re.findall(r'\]\((/[a-z][a-zA-Z0-9/_-]+)\)', src))
    tgt_internal = set(re.findall(r'\]\((/[a-z][a-zA-Z0-9/_-]+)\)', tgt))
    # Links that appear in both source and target unchanged could mean missing locale rewrite
    unchanged_links = src_internal & tgt_internal
    # Only flag if the link is to a page that is commonly translated (heuristic check)
    commonly_translated = {"/quick-start/next-steps", "/getting-started/installation",
                           "/quick-start/why-wails", "/status"}
    missing_rewrites = unchanged_links & commonly_translated
    if missing_rewrites:
        issues.append(f"Links to translated pages missing locale prefix: {sorted(missing_rewrites)}")
        score -= 0.05 * len(missing_rewrites)

    # 14. Check d2 diagram string labels are not translated
    #     d2 uses "Quoted strings" as node labels — these are part of the diagram code
    #     and must not be translated even though they look like prose.
    src_d2_blocks = re.findall(r'```d2[\s\S]*?```', src)
    tgt_d2_blocks = re.findall(r'```d2[\s\S]*?```', tgt)
    for i, (sd, td) in enumerate(zip(src_d2_blocks, tgt_d2_blocks)):
        src_labels = re.findall(r'"([^"]+)"', sd)
        tgt_labels = re.findall(r'"([^"]+)"', td)
        if src_labels != tgt_labels:
            changed = [(s, t) for s, t in zip(src_labels, tgt_labels) if s != t]
            if changed:
                issues.append(
                    f"d2 diagram string labels translated in block {i+1}: "
                    + ", ".join(f'"{s}" → "{t}"' for s, t in changed[:3])
                )
                score -= 0.15

    heuristic_score = max(0.0, min(1.0, score))

    # 15. AI verification (optional, calls z.ai API)
    ai_result = None
    if ai_verify:
        lang_name = LOCALE_NAMES.get(locale, locale)
        ai_result = ai_verify_translation(src_body, tgt_body, locale, lang_name, ai_model)
        if ai_result["score"] is not None:
            issues.extend(ai_result["issues"])
            # Blend: heuristic 60%, AI 40%
            combined = heuristic_score * 0.6 + ai_result["score"] * 0.4
            return {
                "score": round(combined, 3),
                "heuristic_score": round(heuristic_score, 3),
                "ai_score": round(ai_result["score"], 3),
                "ai_used": True,
                "issues": issues,
            }
        else:
            # AI unavailable — fall back to heuristic only, but surface the error
            issues.extend(ai_result["issues"])

    return {
        "score": round(heuristic_score, 3),
        "heuristic_score": round(heuristic_score, 3),
        "ai_score": None,
        "ai_used": False,
        "issues": issues,
    }


def score_locale(locale: str, ai_verify: bool = False, ai_model: str = ZAI_DEFAULT_MODEL) -> dict:
    locale_dir = DOCS_ROOT / locale
    if not locale_dir.exists():
        return {"locale": locale, "error": "Locale directory not found", "avg_score": 0.0, "files": []}

    files = list(locale_dir.rglob("*.mdx")) + list(locale_dir.rglob("*.md"))
    if not files:
        return {"locale": locale, "error": "No translated files found", "avg_score": 0.0, "files": []}

    results = []
    for tgt_path in sorted(files):
        rel = tgt_path.relative_to(locale_dir)
        src_path = DOCS_ROOT / rel
        if not src_path.exists():
            continue
        result = check_file_pair(src_path, tgt_path, locale, ai_verify=ai_verify, ai_model=ai_model)
        result["file"] = str(rel)
        results.append(result)

    if not results:
        return {"locale": locale, "error": "No matching source files found", "avg_score": 0.0, "files": []}

    avg = sum(r["score"] for r in results) / len(results)
    low_quality = [r for r in results if r["score"] < 0.75]

    return {
        "locale": locale,
        "avg_score": round(avg, 3),
        "files_scored": len(results),
        "low_quality_count": len(low_quality),
        "low_quality": low_quality,
        "files": results,
    }


def main():
    parser = argparse.ArgumentParser(description="QA score Wails v3 translations")
    parser.add_argument("--locale", default="all")
    parser.add_argument("--json", action="store_true", help="Output full JSON")
    parser.add_argument("--ai-verify", action="store_true",
                        help="Use z.ai LLM to verify translation accuracy (requires ZAI_API_KEY env var)")
    parser.add_argument("--ai-model", default=ZAI_DEFAULT_MODEL,
                        help=f"z.ai model to use for verification (default: {ZAI_DEFAULT_MODEL})")
    args = parser.parse_args()

    if args.ai_verify and not ZAI_API_KEY:
        print("⚠ --ai-verify requires ZAI_API_KEY environment variable to be set.", flush=True)
        print(f"  Set it with: export ZAI_API_KEY=<your-key>", flush=True)
        print(f"  Using heuristic scoring only.", flush=True)

    if args.locale == "all":
        locales = ALL_LOCALES
    else:
        locales = [l.strip() for l in args.locale.split(",")]

    all_results = {}
    ai_label = f" (+ z.ai/{args.ai_model})" if args.ai_verify and ZAI_API_KEY else ""
    print(f"\nQA Scoring Results{ai_label}", flush=True)
    print("="*70, flush=True)
    print(f"{'Locale':<10} {'Avg Score':>10} {'Files':>7} {'Low Quality':>12}", flush=True)
    print("-"*70, flush=True)

    for locale in locales:
        result = score_locale(locale, ai_verify=args.ai_verify, ai_model=args.ai_model)
        all_results[locale] = result
        if "error" in result:
            print(f"{locale:<10} {'ERROR':>10}  {result['error']}", flush=True)
        else:
            flag = " ⚠" if result["low_quality_count"] > 0 else " ✓"
            print(
                f"{locale:<10} {result['avg_score']:>10.3f} {result['files_scored']:>7} {result['low_quality_count']:>12}{flag}",
                flush=True
            )
            if result["low_quality"]:
                for lq in result["low_quality"]:
                    ai_note = f" [AI: {lq['ai_score']:.3f}]" if lq.get("ai_score") is not None else ""
                    print(f"  → LOW ({lq['score']:.3f}{ai_note}): {lq['file']}", flush=True)
                    for issue in lq["issues"]:
                        print(f"    - {issue}", flush=True)

    print("="*70, flush=True)

    if args.json:
        print(json.dumps(all_results, indent=2))

    import datetime
    today = datetime.date.today().strftime("%Y-%m-%d")

    # Save per-locale reports to .translation-qa/
    QA_DIR.mkdir(exist_ok=True)
    for locale, result in all_results.items():
        per_locale_path = QA_DIR / f"{locale}-{today}.json"
        per_locale_path.write_text(json.dumps(result, indent=2, ensure_ascii=False))

    # Save combined results
    combined_path = QA_DIR / f"all-{today}.json"
    combined_path.write_text(json.dumps(all_results, indent=2, ensure_ascii=False))
    print(f"\nQA results saved to: {QA_DIR}/", flush=True)

    return all_results


if __name__ == "__main__":
    main()

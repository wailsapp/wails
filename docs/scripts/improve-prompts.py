#!/usr/bin/env python3
"""
3-iteration prompt improvement loop for Wails v3 doc translation.

Tests current prompt against AI-improved variants, A/B compares outputs,
and writes the winning prompt back to translate-docs.py.

Usage:
  python3 improve-prompts.py
  python3 improve-prompts.py --locale de --iterations 3
  python3 improve-prompts.py --locale de,ja --dry-run
"""
import os, sys, json, re, time, argparse, textwrap
from pathlib import Path

try:
    import requests
except ImportError:
    os.system("pip3 install requests -q")
    import requests

OLLAMA_BASE = "http://ai-master.taileaa27f.ts.net:11434"
TRANSLATE_MODEL = "qwen3.6:35b"
EVAL_MODEL = "qwen3.6:35b"  # same model; evaluate mode via different prompt

DOCS_ROOT = Path(__file__).parent.parent / "src" / "content" / "docs"
WORK_DIR = Path(__file__).parent.parent / ".prompt-improvement"

LOCALE_NAMES = {
    "de": "German",
    "zh-cn": "Simplified Chinese",
    "zh-tw": "Traditional Chinese",
    "ja": "Japanese",
    "ko": "Korean",
    "ru": "Russian",
    "fr": "French",
    "pt": "Portuguese",
}

# ── Test files ───────────────────────────────────────────────────────────────
# why-wails.mdx: prose-heavy, d2 diagram, known issues from PR review
# installation.mdx: code-heavy, technical commands, frontmatter-rich
TEST_FILES = [
    "quick-start/why-wails.mdx",
    "getting-started/installation.mdx",
]

# ── Known issues from PR review (seed for evaluator context) ─────────────────
KNOWN_ISSUES = """
From PR review of previous translations, these systemic issues were found:
- Product name "Wails" was altered in Portuguese (became "Wais") — proper nouns must be preserved exactly
- d2 diagram code block content was translated in German ("Your UI" → "Ihre UI") — ALL code block content must remain in English
- Korean milliseconds translated to wrong word ("하마초" instead of "밀리초")
- zh-tw: direction description error ("向後端" instead of "向前端" — "to the frontend")
- Frontmatter double-dash artifact (---\\n--- at start) appeared in multiple locales
- Trailing --- artifact appeared at end of files
These suggest the prompt needs stronger rules about: (1) product name preservation,
(2) ALL code block content being untranslatable including d2 diagrams,
(3) precision in technical direction terms.
"""

# ── V0: Current production prompt ────────────────────────────────────────────
PROMPT_V0 = {
    "name": "V0 (baseline)",
    "system": """You are a professional technical documentation translator.
Translate the given MDX/Markdown documentation accurately.

STRICT RULES — violating any of these will cause the page to break:

1. Translate ALL prose text naturally in {lang}. Do not be overly literal.

2. PRESERVE frontmatter structure exactly:
   - Keep all YAML keys in English (title:, description:, link:, icon:, etc.)
   - Translate ONLY the string values of: title, description, tagline, text, label, alt, content (banner content)
   - Copy all other frontmatter values (links, icons, variants, booleans) UNCHANGED

3. NEVER translate ANY of the following — copy them character-for-character:
   - Code blocks (``` ```) — including ALL content inside them
   - Inline code (`code`)
   - Code comments (// ..., # ..., /* ... */) inside code blocks
   - JSX/MDX component names and props (e.g. <Tabs>, <TabItem label="...">)
   - import statements
   - URLs (http://, https://)
   - File paths and CLI commands
   - Variable names and function names
   - d2 diagram definitions

4. FRONTMATTER DELIMITERS — CRITICAL:
   - The document MUST begin with exactly one line containing only: ---
   - The frontmatter block MUST end with exactly one line containing only: ---
   - Do NOT output two --- lines at the start
   - Do NOT add any --- at the end of the document body

5. Return ONLY the translated document. No preamble, no explanation, no markdown fences around the whole document.

6. Preserve all blank lines, heading levels, and list structure exactly as in the source.""",

    "user": """Translate this Wails v3 documentation page to {lang}.
Apply all rules strictly. Return only the translated document.

---
{content}
---""",
}

# ── Ollama call ───────────────────────────────────────────────────────────────

def ollama_call(system: str, user: str, model: str = TRANSLATE_MODEL,
                temperature: float = 0.2, max_tokens: int = 8192,
                think: bool = False) -> str | None:
    payload = {
        "model": model,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        "stream": True,
        "think": think,
        "options": {"temperature": temperature, "top_p": 0.9, "num_predict": max_tokens},
    }
    parts = []
    try:
        resp = requests.post(f"{OLLAMA_BASE}/api/chat", json=payload,
                             stream=True, timeout=(10, 600))
        resp.raise_for_status()
        for line in resp.iter_lines():
            if not line:
                continue
            try:
                data = json.loads(line)
            except json.JSONDecodeError:
                continue
            token = data.get("message", {}).get("content", "")
            if token:
                parts.append(token)
            if data.get("done"):
                break
    except Exception as e:
        print(f"  Ollama error: {e}", flush=True)
        return None
    return "".join(parts).strip() or None


# ── Translation ───────────────────────────────────────────────────────────────

def translate(content: str, lang: str, prompt: dict, file_label: str) -> str | None:
    system = prompt["system"].replace("{lang}", lang)
    user = prompt["user"].replace("{lang}", lang).replace("{content}", content)
    print(f"    → translating {file_label}...", flush=True)
    result = ollama_call(system, user, temperature=0.2)
    if not result:
        return None
    # Strip wrapping markdown fences
    if result.startswith("```") and result.endswith("```"):
        result = "\n".join(result.split("\n")[1:-1])
    # Strip leading text before frontmatter
    if not result.startswith("---") and "---" in result:
        result = result[result.index("---"):]
    return result


# ── Evaluation ────────────────────────────────────────────────────────────────

EVAL_SYSTEM = """You are a professional translation quality assessor and prompt engineer.
You evaluate machine translations of technical documentation and identify how the TRANSLATION PROMPT
should be improved to prevent recurring issues.

Be precise, concise, and actionable. Focus on SYSTEMIC issues (patterns that would affect many files),
not one-off word choices."""

EVAL_USER = """Evaluate this {lang} translation of an English technical documentation page.

## Known systemic issues from prior translation runs:
{known_issues}

## English source:
```
{source}
```

## {lang} translation:
```
{translation}
```

Respond with JSON:
{{
  "score": <0.0-1.0>,
  "accuracy_score": <0.0-1.0>,
  "fluency_score": <0.0-1.0>,
  "technical_score": <0.0-1.0>,
  "issues": [
    {{"severity": "critical|major|minor", "description": "...", "example": "source snippet → translated snippet"}}
  ],
  "prompt_improvements": [
    "Specific rule or wording to add/change in the system prompt to prevent this issue"
  ],
  "summary": "One sentence overall assessment"
}}"""


def evaluate(source: str, translation: str, lang: str) -> dict:
    user = EVAL_USER.format(
        lang=lang,
        known_issues=KNOWN_ISSUES,
        source=source[:3000],
        translation=translation[:3000],
    )
    print(f"    → evaluating...", flush=True)
    result = ollama_call(EVAL_SYSTEM, user, model=EVAL_MODEL, temperature=0.1, max_tokens=2048)
    if not result:
        return {"score": 0.5, "issues": [], "prompt_improvements": [], "summary": "Evaluation failed"}
    # Strip markdown fences
    if result.startswith("```"):
        result = "\n".join(result.split("\n")[1:-1])
    # Find JSON object
    m = re.search(r'\{[\s\S]+\}', result)
    if not m:
        return {"score": 0.5, "issues": [], "prompt_improvements": [], "summary": result[:200]}
    try:
        return json.loads(m.group(0))
    except json.JSONDecodeError:
        return {"score": 0.5, "issues": [], "prompt_improvements": [result[:500]], "summary": "Parse error"}


# ── Prompt evolution ──────────────────────────────────────────────────────────

IMPROVE_SYSTEM = """You are a prompt engineer specializing in machine translation of technical documentation.
Given a current translation prompt and a list of issues found in its outputs, produce an IMPROVED version
of the system prompt that will prevent those issues while preserving what works well.

Output ONLY the improved system prompt text — no explanation, no markdown fences, no commentary."""

IMPROVE_USER = """Current system prompt:
---
{current_prompt}
---

Issues found across translated files (most critical first):
{issues}

Produce an improved version of the system prompt that directly addresses these issues.
The prompt must still cover all the original rules but should be clearer and more specific
where the issues indicate ambiguity."""


def generate_improved_prompt(current_prompt: dict, all_issues: list[str]) -> dict:
    issues_text = "\n".join(f"- {i}" for i in all_issues)
    user = IMPROVE_USER.format(
        current_prompt=current_prompt["system"],
        issues=issues_text,
    )
    print(f"  → generating improved prompt...", flush=True)
    improved_system = ollama_call(IMPROVE_SYSTEM, user, temperature=0.3, max_tokens=4096)
    if not improved_system:
        return current_prompt

    return {
        "name": f"V{current_prompt['name'][1]} improved",
        "system": improved_system,
        "user": current_prompt["user"],  # user template stays same
    }


# ── Heuristic checks ──────────────────────────────────────────────────────────

def heuristic_score(src: str, tgt: str, locale: str) -> float:
    score = 1.0
    if len(tgt.strip()) < 50:
        return 0.0
    if src.strip() == tgt.strip():
        return 0.0
    ratio = len(tgt) / max(len(src), 1)
    if ratio < (0.3 if locale in ["zh-cn", "zh-tw", "ja", "ko"] else 0.5):
        score -= 0.3
    if tgt.startswith("---\n---\n"):
        score -= 0.3
    if tgt.rstrip().endswith("\n---"):
        score -= 0.1
    src_blocks = re.findall(r'```[\s\S]*?```', src)
    tgt_blocks = re.findall(r'```[\s\S]*?```', tgt)
    if len(src_blocks) != len(tgt_blocks):
        score -= 0.15
    for sc, tc in zip(src_blocks, tgt_blocks):
        sc_c = re.sub(r'^```\w*\n?', '', sc).rstrip('`').strip()
        tc_c = re.sub(r'^```\w*\n?', '', tc).rstrip('`').strip()
        if sc_c != tc_c:
            score -= 0.1
            break
    if re.search(r'(?<!\.)\.\.\/\.\.\/assets\/', tgt):
        score -= 0.2
    return max(0.0, min(1.0, score))


# ── Main loop ─────────────────────────────────────────────────────────────────

def run_iteration(prompt: dict, test_files: list[Path], locale: str, iteration: int) -> dict:
    lang = LOCALE_NAMES[locale]
    iter_dir = WORK_DIR / locale / f"iter{iteration}"
    iter_dir.mkdir(parents=True, exist_ok=True)

    results = {}
    for src_path in test_files:
        rel = str(src_path.relative_to(DOCS_ROOT))
        source = src_path.read_text(encoding="utf-8")

        out_path = iter_dir / src_path.name
        # Re-use cached if already translated (supports --dry-run reruns)
        if out_path.exists():
            print(f"  [cached] {rel}", flush=True)
            translation = out_path.read_text(encoding="utf-8")
        else:
            translation = translate(source, lang, prompt, rel)
            if not translation:
                results[rel] = {"error": "translation failed"}
                continue
            out_path.write_text(translation, encoding="utf-8")

        h_score = heuristic_score(source, translation, locale)
        eval_result = evaluate(source, translation, lang)

        combined = h_score * 0.4 + eval_result.get("score", 0.5) * 0.6
        results[rel] = {
            "heuristic": round(h_score, 3),
            "ai_score": round(eval_result.get("score", 0.5), 3),
            "accuracy": round(eval_result.get("accuracy_score", 0.5), 3),
            "fluency": round(eval_result.get("fluency_score", 0.5), 3),
            "technical": round(eval_result.get("technical_score", 0.5), 3),
            "combined": round(combined, 3),
            "issues": eval_result.get("issues", []),
            "prompt_improvements": eval_result.get("prompt_improvements", []),
            "summary": eval_result.get("summary", ""),
        }
        print(f"    scores: heuristic={h_score:.3f} ai={eval_result.get('score',0):.3f} combined={combined:.3f}", flush=True)
        print(f"    summary: {eval_result.get('summary','')[:120]}", flush=True)

    return results


def extract_all_improvements(iter_results: dict) -> list[str]:
    """Collect all unique prompt improvement suggestions from an iteration."""
    seen = set()
    improvements = []
    for file_result in iter_results.values():
        if "error" in file_result:
            continue
        for imp in file_result.get("prompt_improvements", []):
            if imp not in seen:
                seen.add(imp)
                improvements.append(imp)
        for issue in file_result.get("issues", []):
            if issue.get("severity") == "critical":
                desc = issue.get("description", "")
                ex = issue.get("example", "")
                key = desc[:80]
                if key not in seen:
                    seen.add(key)
                    improvements.append(f"Critical: {desc}" + (f" (e.g. {ex})" if ex else ""))
    return improvements


def avg_combined(iter_results: dict) -> float:
    scores = [r["combined"] for r in iter_results.values() if "combined" in r]
    return sum(scores) / len(scores) if scores else 0.0


def print_comparison_table(all_iters: list[tuple[str, dict]]):
    print("\n" + "="*80, flush=True)
    print("A/B COMPARISON TABLE", flush=True)
    print("="*80, flush=True)
    # Collect all file names
    files = []
    for _, results in all_iters:
        for f in results:
            if f not in files:
                files.append(f)

    header = f"{'File':<40}" + "".join(f" {name:>12}" for name, _ in all_iters)
    print(header, flush=True)
    print("-"*80, flush=True)
    for f in files:
        row = f"{f[:38]:<40}"
        for _, results in all_iters:
            r = results.get(f, {})
            val = f"{r.get('combined', 0):.3f}" if "combined" in r else "  N/A"
            row += f" {val:>12}"
        print(row, flush=True)

    print("-"*80, flush=True)
    avg_row = f"{'AVERAGE':<40}"
    for _, results in all_iters:
        avg_row += f" {avg_combined(results):>12.3f}"
    print(avg_row, flush=True)
    print("="*80, flush=True)


def apply_winning_prompt(winning_prompt: dict):
    """Patch translate-docs.py with the winning system prompt."""
    script_path = Path(__file__).parent / "translate-docs.py"
    content = script_path.read_text(encoding="utf-8")

    # Find existing SYSTEM_PROMPT assignment and replace it
    pattern = r'(SYSTEM_PROMPT\s*=\s*""")([\s\S]*?)(""")'
    m = re.search(pattern, content)
    if not m:
        print("  ✗ Could not find SYSTEM_PROMPT in translate-docs.py — manual update needed", flush=True)
        return False

    new_content = content[:m.start()] + f'SYSTEM_PROMPT = """{winning_prompt["system"]}"""' + content[m.end():]
    script_path.write_text(new_content, encoding="utf-8")
    print(f"  ✓ Patched translate-docs.py with winning prompt", flush=True)
    return True


def main():
    parser = argparse.ArgumentParser(description="3-iteration translation prompt improvement loop")
    parser.add_argument("--locale", default="de", help="Comma-separated locales to test (default: de)")
    parser.add_argument("--iterations", type=int, default=3, help="Number of iterations (default: 3)")
    parser.add_argument("--dry-run", action="store_true", help="Use cached translations if available")
    parser.add_argument("--no-patch", action="store_true", help="Don't write winning prompt to translate-docs.py")
    args = parser.parse_args()

    locales = [l.strip() for l in args.locale.split(",") if l.strip() in LOCALE_NAMES]
    if not locales:
        print(f"Unknown locale(s). Valid: {', '.join(LOCALE_NAMES)}")
        sys.exit(1)

    test_file_paths = [DOCS_ROOT / f for f in TEST_FILES if (DOCS_ROOT / f).exists()]
    WORK_DIR.mkdir(parents=True, exist_ok=True)

    print(f"\n{'='*80}", flush=True)
    print(f"TRANSLATION PROMPT IMPROVEMENT LOOP", flush=True)
    print(f"Locales: {', '.join(locales)} | Files: {len(test_file_paths)} | Iterations: {args.iterations}", flush=True)
    print(f"{'='*80}\n", flush=True)

    all_locale_results = {}

    for locale in locales:
        lang = LOCALE_NAMES[locale]
        print(f"\n{'─'*60}", flush=True)
        print(f"Locale: {locale} ({lang})", flush=True)
        print(f"{'─'*60}", flush=True)

        current_prompt = PROMPT_V0.copy()
        current_prompt["name"] = "V0"
        iteration_history = []  # [(prompt_name, results)]

        for iteration in range(args.iterations):
            print(f"\n  [Iteration {iteration+1}/{args.iterations}] Prompt: {current_prompt['name']}", flush=True)

            results = run_iteration(current_prompt, test_file_paths, locale, iteration)
            iteration_history.append((current_prompt["name"], results))

            avg = avg_combined(results)
            print(f"  Average combined score: {avg:.3f}", flush=True)

            if iteration < args.iterations - 1:
                # Collect issues and generate improved prompt for next iteration
                improvements = extract_all_improvements(results)
                print(f"  Found {len(improvements)} improvement suggestions:", flush=True)
                for imp in improvements[:5]:
                    print(f"    · {textwrap.shorten(imp, 100)}", flush=True)

                next_prompt = generate_improved_prompt(current_prompt, improvements)
                next_prompt["name"] = f"V{iteration+1}"
                current_prompt = next_prompt

                # Save improved prompt to disk for inspection
                prompt_path = WORK_DIR / locale / f"prompt_v{iteration+1}.txt"
                prompt_path.write_text(next_prompt["system"], encoding="utf-8")
                print(f"  Saved improved prompt to: {prompt_path}", flush=True)

        # Comparison
        print_comparison_table(iteration_history)

        # Pick winning prompt (highest avg combined score)
        best_idx = max(range(len(iteration_history)), key=lambda i: avg_combined(iteration_history[i][1]))
        best_name, best_results = iteration_history[best_idx]
        print(f"\n  Winner: {best_name} (avg combined = {avg_combined(best_results):.3f})", flush=True)

        # Save full report
        report = {
            "locale": locale,
            "test_files": [str(f.relative_to(DOCS_ROOT)) for f in test_file_paths],
            "iterations": [
                {"prompt": name, "avg_combined": round(avg_combined(r), 3), "results": r}
                for name, r in iteration_history
            ],
            "winner": best_name,
        }
        report_path = WORK_DIR / f"{locale}-report.json"
        report_path.write_text(json.dumps(report, indent=2, ensure_ascii=False), encoding="utf-8")
        print(f"  Report saved to: {report_path}", flush=True)
        all_locale_results[locale] = report

        # Apply winning prompt to translate-docs.py (use last iteration's prompt as winner)
        # The winning prompt is the one from the best scoring iteration
        # We need to find which prompt object corresponds to best_idx
        # Prompts are: iter 0 = V0, iter 1 = V1, iter 2 = V2
        # The prompt used in iteration i is saved at prompt_v{i}.txt (for i>0) or is PROMPT_V0
        if not args.no_patch and best_idx > 0:
            winning_prompt_path = WORK_DIR / locale / f"prompt_v{best_idx}.txt"
            if winning_prompt_path.exists():
                winning_prompt = {
                    "system": winning_prompt_path.read_text(encoding="utf-8"),
                    "user": PROMPT_V0["user"],
                }
                print(f"\nApplying winning prompt ({best_name}) to translate-docs.py...", flush=True)
                apply_winning_prompt(winning_prompt)
        elif not args.no_patch and best_idx == 0:
            print(f"\n  V0 (baseline) is already the best — no changes to translate-docs.py needed.", flush=True)

    print(f"\n{'='*80}", flush=True)
    print("IMPROVEMENT LOOP COMPLETE", flush=True)
    print(f"Results saved to: {WORK_DIR}/", flush=True)
    print(f"{'='*80}\n", flush=True)

    return all_locale_results


if __name__ == "__main__":
    main()

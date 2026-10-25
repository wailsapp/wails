---
title: "Release-Richtlinie"
description: "Wie Wails v3 seine Nightly- und Weekly-Releases versioniert, wie Maintainer sie erstellen und warum Benutzer Versionen festschreiben sollten."
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

> **Status: Entwurf zur Prüfung durch den Verantwortlichen.** Mit **(Entscheidung ausstehend)** gekennzeichnete Punkte sind nicht bestätigt.

## Das Modell in einer Tabelle

Wails v3 veröffentlicht **zeitgesteuert, nicht funktionsgesteuert**. Nach der `v3.0.0`-Basisversion:

| Release | Tag | Rhythmus | GitHub-Release | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| Nightly | `v3.<next>.0-nightly.<YYYYMMDD>` (UTC-Datum; bei erneuter Ausführung am selben Tag wird `.N` angehängt) | automatisiert, täglich, nur wenn `UNRELEASED_CHANGELOG.md` Inhalt enthält | Pre-Release, **nicht** Latest | dieselbe Version ohne `v`, Dist-Tag `next` |
| Weekly | `v3.<next>.0` | automatisiert, sonntags um 20:00 Uhr Australia/Sydney **(Tag nicht bestätigt)** | reguläres Release, Latest **(Entscheidung ausstehend)** | dieselbe Version, Dist-Tag `latest` |
| Hotfix | `v3.<current>.<N>`, `N > 0` | selten, manuell | reguläres Release, Latest | dieselbe Version, Dist-Tag `latest` |
| Major | `v4.0.0` | nur bei inkompatiblen Änderungen | reguläres Release | `latest` |

`<next>` ist die höchste **veröffentlichte** v3-Minor-Version plus eins. Ein Nightly ist eine Vorschau auf das **nächste Weekly**. Daher basiert es auf dem höchsten Release-Tag; Pre-Releases werden ignoriert: Nach einem Hotfix `v3.1.1` ist das nächste Nightly weiterhin `v3.2.0-nightly.*`.

**Eine Erhöhung der Minor-Version ist kein Feature-Meilenstein.** `v3.4.0` ist „das Weekly-Release nach `v3.3.0`“. Lesen Sie die Release Notes statt der Versionsnummer, um zu sehen, was sich geändert hat.

Die Richtlinie wird anhand der **Form** der Version ausgewählt: Ein `-nightly.`-Suffix kennzeichnet ein Nightly; jedes andere Suffix (`-alpha`, `-beta`, `-rc`) kennzeichnet einen älteren Build aus der Zeit vor GA; ohne Suffix handelt es sich um ein Release.

### Synchronisierte Versionen

Die CLI, das Go-Modul und das npm-Paket `@wailsio/runtime` verwenden **dieselbe Version** aus `v3/internal/version/version.txt` (`3.2.0-nightly.20261102` auf npm für den Tag `v3.2.0-nightly.20261102`). Mit `wails3 init` erstellte Projekte schreiben `@wailsio/runtime` auf die eigene Version der CLI fest; Entwicklungs-Builds der CLI greifen ersatzweise auf `latest` zurück.

## Durchgängiges Beispiel

| Datum | Ereignis | Tag | `go install ...@latest` wird aufgelöst zu |
| --- | --- | --- | --- |
| So., 25. Okt. | GA, manuelles Release | `v3.0.0` | `v3.0.0` |
| Mo., 26. bis Sa., 31. Okt. | Nightlies | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| So., 1. Nov., 20:00 Uhr Sydney | Weekly | `v3.1.0` | `v3.1.0` |
| Mo., 2. Nov. | Nightly (Vorschau auf das nächste Weekly) | `v3.2.0-nightly.20261102` | `v3.1.0` |
| Mi., 4. Nov. | Hotfix, manuell, aus einem Release-Branch | `v3.1.1` | `v3.1.1` |
| So., 8. Nov. | Weekly | `v3.2.0` | `v3.2.0` |

Ein Weekly wird immer oberhalb aller Nightlies mit derselben Versionsnummer einsortiert, und ein Nightly wird niemals zu `@latest`.

## Was Benutzer ausführen

| Ziel | Befehl |
| --- | --- |
| Neuestes Release (Weekly oder Hotfix) | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` oder `wails3 update cli` |
| **Reproduzierbar (empfohlen)** | einen exakten Tag festschreiben: `v3.1.0` in `go.mod`, CI und Dokumentation |
| Ein Nightly | dessen exakter Tag (`@v3.2.0-nightly.20261102`), `wails3 update cli -nightly` |
| Neuester Build beliebiger Art | `wails3 update cli -pre` |
| Runtime-npm-Paket | `@wailsio/runtime` (`latest` = Weekly/Hotfix, `next` = Nightly) |

Warum dies funktioniert: Go löst `@latest` zum höchsten Tag **ohne Pre-Release-Suffix** auf und berücksichtigt Pre-Releases nur, wenn kein Release vorhanden ist. Das GitHub-Kontrollkästchen „pre-release“ ist für den Go-Proxy nicht sichtbar. Daher tragen Nightlies einen SemVer-Pre-Release-Suffix, statt sich darauf zu verlassen. Schreiben Sie Versionen in Produktionsumgebungen immer fest.

## Branches

| Branch | Zweck |
| --- | --- |
| `master` | Entwicklung und Nightlies. |
| `releases/**` | Kurzlebige Branches für einmalige Releases: das `v3.0.0`-GA-Release und Hotfixes. Es gelten dieselben erforderlichen CI-Prüfungen wie für `master`. |

Weeklies werden aus `master` erstellt, jedoch **nicht von dessen Spitze**: Der Workflow wählt den neuesten Commit seit dem vorherigen Weekly aus, für den alle erforderlichen Prüfungen (`v3 Build Results`, `Run Go Tests`, `Bundled runtime matches source`) bestanden wurden, und versieht genau diesen Commit mit einem Tag. Ist keiner erfolgreich, wird nichts veröffentlicht. Nightlies werden aus `master` erstellt und geben nur eine Warnung aus, wenn HEAD nicht erfolgreich geprüft wurde.

## So erstellen Sie ein Release

Release-Berechtigungen sind auf die Projektleitung beschränkt (siehe Governance-Dokument). Mitwirkende erstellen oder veröffentlichen niemals Tags; sie fügen Einträge zu `v3/UNRELEASED_CHANGELOG.md` hinzu.

Alles wird über **Nightly + Weekly Release v3** (`.github/workflows/nightly-release-v3.yml`) ausgeführt:

- **Nightly (automatisch)** um 15:00 Uhr UTC. Wird übersprungen, wenn das Changelog für unveröffentlichte Änderungen leer ist.
- **Weekly (automatisch)** um 20:00 Uhr Sydney-Zeit. GitHub-Cron unterstützt nur UTC. Daher werden zwei Cron-Jobs (um 09:00 und 10:00 Uhr UTC) geplant, und ein vorgeschalteter Prüfjob fährt nur fort, wenn es in Sydney 20:00 Uhr ist; dabei wird die Sommerzeit berücksichtigt. Das Weekly passiert das Release-Gate, wählt den neuesten erfolgreich geprüften Commit aus, erhöht die Minor-Version, versieht diesen Commit mit einem Tag, ohne das Changelog zu ändern, und kopiert anschließend die Versionsdateien zurück auf `master`.
- **Manuell.** Führen Sie den Workflow mit `channel`, `dry_run` (standardmäßig aktiviert) und optional `release_version` aus. Ein Probelauf gibt Version, Art, GitHub-Flags, npm-Dist-Tag und Hinweise aus, ohne GitHub oder npm zu verändern.
- **GA oder Hotfix.** Legen Sie `ref: releases/...` und `release_version` (`v3.0.0`, `v3.1.1`) fest. Der Release-Commit wird per Cherry-Pick zurück nach `master` übernommen.

Lokaler Probelauf: `cd v3/tasks/release && go run release.go --dry-run --channel weekly`; `go run release.go --print-policy v3.2.0-nightly.20261102` zeigt die abgeleiteten Flags an.

### Sicherheitsmechanismen

- Eine automatisch ermittelte Version muss gemäß semantischer Versionsreihenfolge strikt größer sein als jedes vorhandene Tag (Vorabversion < Release-Version), andernfalls verweigert das Tool den Vorgang.
- Es wird jeweils nur ein Release-Job ausgeführt (gemeinsame Parallelitätsgruppe).
- Die Veröffentlichung auf npm wird durch das Release-**Tag** ausgelöst und baut den mit dem Tag versehenen Commit.

## Release-Sperre

Der **wöchentliche Release** (niemals der Nightly-Build) verweigert das Setzen eines Tags, solange eines der folgenden Elemente vorhanden ist:

1. offene Issues oder Pull Requests mit dem Label `P0` oder `release-blocker`;
2. offene Issues mit dem Label `regression`, die seit dem letzten wöchentlichen Release erstellt wurden;
3. offene Dependabot-Warnungen mit kritischem oder hohem Schweregrad, offene Code-Scanning-Warnungen mit kritischem oder hohem Schweregrad sowie alle offenen Secret-Scanning-Warnungen, **jedoch nur für ausgelieferten Code**. Warnungen, die sich ausschließlich unter `v3/examples/`, `v2/examples/`, `website/`, `docs/`, `.github/` oder in einem beliebigen `examples/`-Verzeichnis befinden, werden ignoriert (`SHIPPED_ALERT_IGNORE` in `v3/tasks/releasegate/main.go`, **vom Verantwortlichen zu prüfen**);
4. Sicherheitshinweise des Repositorys, die sich noch im Status `triage` befinden.

**Wer eine Sperre setzen darf.** Ein Label zählt nur, wenn das zugehörige neueste `labeled`-Ereignis (ohne späteres `unlabeled`) von einem Mitglied des `wailsapp/maintainers`-Teams oder von einem Mitwirkenden mit `maintain`-/`admin`-Berechtigung ausgelöst wurde. Wird dasselbe Label von einer anderen Person angewendet, wird es ignoriert und protokolliert, sodass Außenstehende einen Release weder sperren noch entsperren können. Das automatisch angewendete Label `high-priority` wird absichtlich **nicht** verwendet. Bot-Konten wie `taliesin-ai` bestehen die Zulassungsprüfung, wenn sie über `maintain` verfügen **(Entscheidung ausstehend: Bots zulassen oder ausschließen)**.

**Verhalten bei Fehlern.** Die Sperre schlägt **restriktiv** fehl: Wenn eine Warnungs-API oder die Zulassungsliste nicht gelesen werden kann, wird der wöchentliche Release abgebrochen. Bei einem Fehler setzt er kein Tag, schreibt die Liste einschließlich Links in die Job-Zusammenfassung und aktualisiert einen Kommentar im einzigen Issue „Weekly release status“, sofern die Repository-Variable `WEEKLY_STATUS_ISSUE` gesetzt ist.

**Vom Verantwortlichen auf GitHub vorzunehmende Einrichtung** (nichts wird automatisch erstellt): Erstelle die Labels `P0`, `release-blocker` und `regression`; erstelle das Issue „Weekly release status“ und setze `WEEKLY_STATUS_ISSUE` auf dessen Nummer; gewähre `WAILS_REPO_TOKEN` Lesezugriff auf Issues, Pull Requests, Administration (Berechtigung für Mitwirkende), Dependabot-Warnungen, Code-Scanning-Warnungen, Secret-Scanning-Warnungen, Sicherheitshinweise des Repositorys und Organisationsmitglieder sowie Schreibzugriff auf Issues für den Statuskommentar. Das standardmäßige `GITHUB_TOKEN` kann diese Daten nicht lesen.

**Außerkraftsetzung.** `override_reason` (Workflow-Dispatch-Eingabe) überspringt die Sperre. Der Grund, der ausführende Akteur und jede übersprungene Sperre werden in die Zusammenfassung und das Status-Issue geschrieben. Nur der Release-Verantwortliche sollte diese Option verwenden. GitHub kann eine Eingabe nicht auf eine einzelne Person beschränken. Führe den Release-Job daher in einer Umgebung mit erforderlichen Prüfern aus (empfohlen), damit eine Außerkraftsetzung ausdrücklich genehmigt werden muss.

## Verwandte Themen

- [Tests und Continuous Integration](/contributing/testing-ci/)
- [Triage für Maintainer](/contributing/maintainer-triage/)

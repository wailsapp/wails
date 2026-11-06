---
title: "Kira"
description: "Ein nativer macOS-Desktopclient, der AWS – ECS, RDS, S3, DynamoDB und mehr – in einer einzigen tastaturgesteuerten Oberfläche vereint"
slug: "community/showcase/kira"
sourcePath: "community/showcase/kira.md"
---

**[Kira](https://kira.thiennguyen.dev)** ist ein **nativer macOS-Desktopclient für AWS** – *„AWS ohne Reibungsverluste“*. Die mit **Go, Wails und React + TypeScript** entwickelte Anwendung bündelt AWS-Vorgänge in einer einzigen tastaturgesteuerten App. Authentifizieren Sie sich einmal über AWS SSO und verwalten Sie die Infrastruktur über viele Konten hinweg, ohne zwischen der Konsole, der CLI und zahlreichen Datenbanktools wechseln zu müssen.

Am Anfang stand ein persönliches Ärgernis: Nur um eine Änderung bereitzustellen, musste man ständig zwischen Browser-Tabs, `aws`-Befehlen und einem separaten SQL-Client wechseln – das war umständlich. Kira vereint diese Arbeitsabläufe in einem schnellen, nativen Fenster: Melden Sie sich an, wählen Sie ein Konto aus, und alles, was Sie benötigen, ist nur einen Tastendruck entfernt.

![Kiras kontoübergreifende Übersicht – Produktions-, Staging- und Entwicklungskonten an einem Ort](/assets/showcase-images/kira_screenshot_1.png)

## Wichtigste Funktionen

- **AWS SSO für mehrere Konten** – Melden Sie sich einmal an und wechseln Sie anschließend zentral zwischen Konten und Regionen
- **ECS** – Durchsuchen Sie Cluster, Services und Tasks; stellen Sie Services erneut bereit, skalieren Sie sie oder setzen Sie sie zurück; zeigen Sie Task-Definitionen an, überwachen Sie Service-Metriken und öffnen Sie über ECS Exec eine interaktive Shell
- **Datenbanken** – Führen Sie SQL-Abfragen für RDS aus, fragen Sie DynamoDB ab und durchsuchen Sie es per Scan, und stellen Sie Verbindungen zu PostgreSQL, MySQL und Redshift her – mit sicheren SSH-Tunneln und im macOS-Schlüsselbund gespeicherten Anmeldedaten
- **S3** – Navigieren Sie durch Buckets und Präfixe; zeigen Sie Vorschauen von Objekten an, laden Sie sie hoch oder herunter, kopieren, benennen Sie sie um oder löschen Sie sie; und erstellen Sie Ordner
- **Secrets Manager** – Listen Sie Secrets auf und rufen Sie deren Werte für das aktive Konto ab
- **CloudWatch Logs** – Verfolgen und durchsuchen Sie Logstreams in Echtzeit
- **Smart Query** – Optionale KI-gestützte SQL-Generierung auf Basis der `claude`-CLI
- **Erweiterungen** – Installieren Sie benutzerdefinierte `.kext`-Bundles, die Aktionsschaltflächen hinzufügen, deren Funktionen durch kleine Go-Skripte bereitgestellt werden
- **Schnelle Navigation** – Ein globales Tastenkürzel zum Aufrufen, eine `Cmd+K`-Befehlspalette und `kira://`-Deep-Links

## Im Detail

Überwachen Sie Ihre ECS-Services live – Task-Zustand, CPU- und Arbeitsspeicherauslastung sowie Bereitstellungsstatus – und stellen Sie sie erneut bereit, skalieren Sie sie oder setzen Sie sie zurück, ohne die Liste zu verlassen.

![Kira beim Durchsuchen von ECS-Services mit Live-Metriken](/assets/showcase-images/kira_screenshot_17.png)

Durchsuchen Sie S3 wie mit einem Dateimanager. Zeigen Sie Objektvorschauen an, prüfen Sie Metadaten und Versionen und laden Sie Objekte direkt hoch oder herunter, benennen Sie sie um oder löschen Sie sie.

![Kiras S3-Objektbrowser mit Objektvorschau und Metadaten](/assets/showcase-images/kira_screenshot_5.png)

Wird als signiertes, notarisiertes `.dmg` für macOS vertrieben.

[Kira besuchen](https://kira.thiennguyen.dev) | [Dokumentation lesen](https://docs.kira.thiennguyen.dev)

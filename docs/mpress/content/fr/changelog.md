---
title: "Journal des modifications"
description: "Historique des versions et notes de publication de Wails v3"
slug: "changelog"
sourcePath: "changelog.md"
---

Légende :

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- Toutes les modifications notables apportées à ce projet seront consignées dans ce fichier.

Le format repose sur [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), et ce projet respecte la [gestion sémantique des versions](https://semver.org/spec/v2.0.0.html).

- `Added` pour les nouvelles fonctionnalités.
- `Changed` pour les modifications apportées aux fonctionnalités existantes.
- `Deprecated` pour les fonctionnalités qui seront bientôt supprimées.
- `Removed` pour les fonctionnalités désormais supprimées.
- `Fixed` pour toutes les corrections de bogues.
- `Security` en cas de vulnérabilités.

_/

/_   * VEUILLEZ NE PAS METTRE À JOUR CE FICHIER *   Les mises à jour devraient être ajoutées à `v3/UNRELEASED_CHANGELOG.md`   Merci ! _/

## [Non publié]

## v3.0.0-beta.21 - 2026-09-13

## Ajouts

- Publication de la documentation v3 de Wails avec M-Press dans la [PR](https://github.com/wailsapp/wails/pull/6116) par @leaanthony

## Corrections

- Analyse des valeurs de slug JSON dans le frontmatter MPD pour générer le journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/6118) par @leaanthony
- L’outil de mise à jour efface les variables d’environnement auxiliaires et relance la cible d’origine après l’échec d’une sauvegarde dans la [PR](https://github.com/wailsapp/wails/pull/6080) par @cnmax
- Démarrage du gestionnaire de signaux par défaut pendant App.Run dans la [PR](https://github.com/wailsapp/wails/pull/6098) par @leaanthony
- Sous Windows, le menu gère les menus nil, libère les ressources remplacées et redessine la barre de menus dans la [PR](https://github.com/wailsapp/wails/pull/6112) par @taliesin-ai
- Rétablissement de la création de packages MSIX pour les nouveaux projets utilisant une configuration YAML partagée dans la [PR](https://github.com/wailsapp/wails/pull/6115) par @leaanthony
- Correction de l’échec du chargement des bindings JavaScript et TypeScript générés lorsque les créateurs de modèles génériques font référence à des déclarations auxiliaires ultérieures, et prévention des dépassements de pile lors de la création de modèles génériques mutuellement dépendants (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Modifié

- Mise à jour des liens de présentation de Clave vers le site web et le dépôt actuels dans la [PR](https://github.com/wailsapp/wails/pull/6082) par @01xR4in

## Corrigé

- Annulation des requêtes de ressources Windows interrompues, y compris celles des workers, tout en préservant les gestionnaires keepalive lors de la navigation. Transmission des contextes de requête natifs via l’enveloppe de l’application sur les plateformes Apple. (#5963, #5969)
- Conservation des entrées du journal des modifications lors de push concurrents grâce à de nouvelles tentatives, dans la [PR](https://github.com/wailsapp/wails/pull/6094) par @leaanthony
- Correction de l’échec de `go mod vendor` avec `pattern arm64/WebView2Loader.dll: no matching files found` sur toutes les plateformes, en supprimant les incorporations qui faisaient référence à des binaires jamais distribués dans le module ; cette modification corrige [#5782](https://github.com/wailsapp/wails/issues/5782) et [#5376](https://github.com/wailsapp/wails/issues/5376), dans la [PR](https://github.com/wailsapp/wails/pull/6031) par @Grantmartin2002

## Supprimé

- Suppression de la prise en charge du chargeur WebView2 natif, remplacé par le chargeur en Go pur. Cette modification supprime les binaires `WebView2Loader.dll` incorporés et la dépendance `github.com/jchv/go-winloader`. Le tag de compilation `native_webview2loader` reste accepté et ne provoque plus d’erreur, mais n’a aucun effet sur les compilations v3, dans la [PR](https://github.com/wailsapp/wails/pull/6031) par @Grantmartin2002
- Suppression des tags de compilation inutilisés et de l’option FPS du guide de l’API macOS dans la [PR](https://github.com/wailsapp/wails/pull/6097) par @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Ajouté

- Protection des API macOS privées par des tags de compilation afin d’en réserver l’utilisation aux activations explicites — consultez la [documentation](https://v3.wails.io/features/browser/integration), la [documentation](https://v3.wails.io/features/environment/info), la [documentation](https://v3.wails.io/features/windows/basics), la [documentation](https://v3.wails.io/features/windows/frameless), la [documentation](https://v3.wails.io/features/windows/notch-windows), la [documentation](https://v3.wails.io/features/windows/options), la [documentation](https://v3.wails.io/guides/build/macos), la [documentation](https://v3.wails.io/guides/build/private-macos-apis) et la [documentation](https://v3.wails.io/reference/overview), dans la [PR](https://github.com/wailsapp/wails/pull/6087) par @leaanthony

## Corrigé

- Rejet des requêtes d’exécution dépassant 64 Mio avec le code HTTP 413 dans la [PR](https://github.com/wailsapp/wails/pull/6091) par @leaanthony

## Sécurité

- Renforcement de la sécurité des origines MCP et de l’accès distant grâce à une authentification par jeton dans la [PR](https://github.com/wailsapp/wails/pull/6092) par @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Corrigé

- Correction d’une fuite de mémoire de Calloc sous Linux et Darwin grâce à l’utilisation de récepteurs de pointeur dans la [PR](https://github.com/wailsapp/wails/pull/6083) par @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Corrigé

- Windows : un `GetRequest` ayant échoué ou nil dans le gestionnaire WebResourceRequested n’arrête plus le processus (`log.Fatal` / panique due au déréférencement de nil) — la requête est désormais abandonnée et consignée, dans la [PR](https://github.com/wailsapp/wails/pull/6006) par @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Modifié

- Demande du mot de passe de notarisation dans une nouvelle fenêtre de terminal, dans la [PR](https://github.com/wailsapp/wails/pull/6029) par @leaanthony

## Corrigé

- Gestion correcte des types de clic sur l’icône de la zone de notification sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5919) par @ChewbaccaCookie
- Suppression par la CI des dépôts apt Microsoft inutilisés avant la mise à jour, dans la [PR](https://github.com/wailsapp/wails/pull/6041) par @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Corrigé

- Augmentation à 60 secondes du délai d’expiration de l’incorporation de WebView2 dans la [PR](https://github.com/wailsapp/wails/pull/6043) par @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Corrigé

- Nommage correct des frappes de touches Contrôle-lettre sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/6032) par @taliesin-ai
- Correction des icônes ICO de la zone de notification et adaptation au thème de la barre des tâches sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/6016) par @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Corrigé

- Poursuite du traitement des tâches du thread principal sous macOS pendant l’exécution d’une boucle modale, dans la [PR](https://github.com/wailsapp/wails/pull/6026) par @leaanthony
- Possibilité d’échec du stockage sécurisé mobile, avec refus par défaut en cas d’échec, dans la [PR](https://github.com/wailsapp/wails/pull/5923) par @mortenolsrud
- Exécution des hooks d’événements de l’application même lorsqu’aucun écouteur n’est enregistré, dans la [PR](https://github.com/wailsapp/wails/pull/5999) par @archy-rock3t-cloud
- Correction de fautes de frappe dans les commentaires et la documentation localisée, dans la [PR](https://github.com/wailsapp/wails/pull/6023) par @haoku123
- Suppression des binaires macOS précompilés ajoutés au dépôt sous `v3/examples` dans la [PR](https://github.com/wailsapp/wails/pull/6025) par @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Ajouté

- Ajout de fenêtres de notification macOS pour l’encoche, avec un exemple de cycle de vie et de télémétrie — consultez la [documentation](https://v3.wails.io/features/windows/notch-windows) dans la [PR](https://github.com/wailsapp/wails/pull/6010) par @leaanthony
- Ajout de la prise en charge des fenêtres NSPanel de macOS, avec de nouvelles options et une intégration native — voir la [documentation](https://v3.wails.io/features/windows/options) dans la [PR](https://github.com/wailsapp/wails/pull/6008) par @leaanthony

## Corrections

- Empêche le blocage de SQLite Prepare lors d’appels simultanés dans la [PR](https://github.com/wailsapp/wails/pull/5998) par @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Suppressions

- Suppression de l’ancien outil de suivi de l’implémentation dans la documentation, dans la [PR](https://github.com/wailsapp/wails/pull/6005) par @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Corrections

- Correction de l’hôte Linux GTK4 qui ignorait les arguments de lancement liés aux protocoles personnalisés et aux associations de fichiers, dans la [PR](https://github.com/wailsapp/wails/pull/6000) par @midagedev
- Prise en charge correcte des lignes supprimées et des corrections provenant de la même source lors de la validation du journal des modifications, dans la [PR](https://github.com/wailsapp/wails/pull/5993) par @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Ajouts

- Ajout du serveur MCP sécurisé wails3 pour la gestion de projets assistée par des agents, dans la [PR](https://github.com/wailsapp/wails/pull/5896) par @leaanthony
- Ajout d’une documentation sur les modèles dans les liaisons — voir la [documentation](https://v3.wails.io/features/bindings/models) dans la [PR](https://github.com/wailsapp/wails/pull/5988) par @taliesin-ai
- Prise en charge de l’installation avec rpm-ostree sur les systèmes Linux atomiques, dans la [PR](https://github.com/wailsapp/wails/pull/5987) par @leaanthony
- Ajout de la génération et de la publication natives d’un graphique hebdomadaire de l’historique des étoiles — voir la [documentation](https://v3.wails.io/credits), la [documentation](https://v3.wails.io/de/credits), la [documentation](https://v3.wails.io/fr/credits), la [documentation](https://v3.wails.io/id/credits), la [documentation](https://v3.wails.io/ja/credits), la [documentation](https://v3.wails.io/ko/credits), la [documentation](https://v3.wails.io/pt/credits), la [documentation](https://v3.wails.io/ru/credits), la [documentation](https://v3.wails.io/zh-cn/credits), la [documentation](https://v3.wails.io/zh-tw/credits) dans la [PR](https://github.com/wailsapp/wails/pull/5986) par @leaanthony
- Ajout du paquet mac réservé à Darwin pour résoudre les ressources du bundle d’application — voir la [documentation](https://v3.wails.io/guides/build/macos) dans la [PR](https://github.com/wailsapp/wails/pull/5965) par @leaanthony
- Ajout d’une page de présentation de Condui et d’une entrée dans l’index — voir la [documentation](https://v3.wails.io/community/showcase/condui) et la [documentation](https://v3.wails.io/community/showcase) dans la [PR](https://github.com/wailsapp/wails/pull/5962) par @mgueregath
- Ajout d’une page de présentation de Redis Viewer comprenant des captures d’écran et un lien vers le projet — voir la [documentation](https://v3.wails.io/community/showcase) et la [documentation](https://v3.wails.io/community/showcase/redisviewer) dans la [PR](https://github.com/wailsapp/wails/pull/5984) par @redisviewer

## Modifications

- Mise à jour des indicateurs d’application GTK vers G<em>APPLICATION</em>NON_UNIQUE sous Linux, dans la [PR](https://github.com/wailsapp/wails/pull/5971) par @overlordtm
- Journalisation des événements de fenêtre manquants au niveau débogage plutôt qu’au niveau avertissement, dans la [PR](https://github.com/wailsapp/wails/pull/5914) par @julianstorer

## Corrections

- Permet aux raccourcis clavier macOS enregistrés d’avoir priorité sur la webview, dans la [PR](https://github.com/wailsapp/wails/pull/5902) par @julianstorer
- Réparation des liens rompus dans la barre latérale de la documentation, dans la [PR](https://github.com/wailsapp/wails/pull/5937) par @northes
- Annulation des contextes des requêtes de ressources sous macOS et iOS lorsque WebKit interrompt la tâche correspondante du schéma personnalisé (#5963)
- Prise en charge du message WindowSetFullscreenButtonEnabled dans la [PR](https://github.com/wailsapp/wails/pull/5976) par @archy-rock3t-cloud
- Importation de Fragment dans le modèle preact-ts afin de corriger l’échec de la compilation, dans la [PR](https://github.com/wailsapp/wails/pull/5979) par @haoku123
- Empêche les anciennes applications GTK3 limitées aux services de planter lorsque la détection des écrans s’exécute avant qu’une fenêtre active ou un affichage ne soit disponible (#5966)
- Permet aux exécutions de publication avec une version explicite de continuer lorsque la section non publiée du journal des modifications est vide (#5977)

## Sécurité

- Mise à jour des fichiers de verrouillage nanoid du site web vers la version corrigée 3.3.18 afin de résoudre les avis de sécurité, dans la [PR](https://github.com/wailsapp/wails/pull/5985) par @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Ajouts

- Ajout de la génération d’URL de documentation aux entrées automatiques du journal des modifications, dans la [PR](https://github.com/wailsapp/wails/pull/5957) par @taliesin-ai
- Ajout de Streams : des flux d’octets bidirectionnels entre Go et JavaScript, avec le modèle de programmation WebSocket et sans socket d’écoute. Déclarez un flux en Go avec `app.HandleStream(name, handler)`, puis connectez-vous depuis le frontend avec `Stream(name)`, qui renvoie un objet de forme `WebSocket`. Le trafic Go→JS est transporté par une requête d’interrogation maintenue par fenêtre via le serveur de ressources, et le trafic JS→Go par une requête POST normale ; aucun port TCP n’est ouvert et rien ne transite par `evaluateJavaScript`. Dans les builds serveur (`-tags server`), le même gestionnaire est à la place servi via un véritable WebSocket, si bien que le code de l’application est identique dans tous les types de builds. Par @leaanthony
- Déplacement de l’entrée du journal des modifications relative à la boîte aux lettres vers la section Non publié, dans la [PR](https://github.com/wailsapp/wails/pull/5935) par @leaanthony

## Modifications

- Mise à jour de la génération automatique de la barre latérale de la documentation et de la dérivation du type des auteurs du blog, dans la [PR](https://github.com/wailsapp/wails/pull/5938) par @leaanthony

## Corrections

- L’initialisation de WebView2 utilise une échéance et une pompe à messages, dans la [PR](https://github.com/wailsapp/wails/pull/5952) par @leaanthony
- Le test des cookies WebView2 est ignoré dans l’intégration continue sauf activation explicite, et son exécution est verrouillée sur le thread actuel du système d’exploitation, dans la [PR](https://github.com/wailsapp/wails/pull/5951) par @leaanthony
- Les générateurs de menus Windows rétablissent les identifiants de commande des éléments parents de sous-menus, dans la [PR](https://github.com/wailsapp/wails/pull/5944) par @gilad-ch
- Alignement de l’image officielle de compilation croisée sur la version minimale de GTK 4.14+ prise en charge sous Linux (#5928)
- Configuration du projet Xcode pour iOS afin de conserver les indicateurs hérités de l’éditeur de liens et d’ajouter -ObjC, dans la [PR](https://github.com/wailsapp/wails/pull/5915) par @mortenolsrud
- Correction du renouvellement excessif des connexions TCP dans le proxy de ressources `wails3 dev` avec les frontends volumineux, qui pouvait épuiser les ports éphémères de l’hôte et faire échouer des processus sans rapport avec l’application en produisant l’erreur `EADDRNOTAVAIL`
- Mise en file d’attente du code JavaScript des événements de chaque fenêtre afin d’assurer leur distribution ordonnée et la contre-pression, dans la [PR](https://github.com/wailsapp/wails/pull/5934) par @leaanthony

## Suppressions

- Suppression du pipeline de publication des binaires pour ordinateur : les versions v3 sont uniquement publiées sous forme de tags, et la CLI `wails3` s’installe avec `go install`. Suppression de `release-v3.yml` et de l’étape nocturne qui le déclenchait, dans la [PR](https://github.com/wailsapp/wails/pull/5946) par @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Ajouts

- Ajout d’une préférence de lecture automatique sur macOS pour désactiver l’exigence d’une action de l’utilisateur lors de la lecture de contenus multimédias dans la [PR](https://github.com/wailsapp/wails/pull/5512) par @Eyalm321
- Déplacement de l’entrée du journal des modifications relative à la boîte aux lettres vers la section « Non publié » dans la [PR](https://github.com/wailsapp/wails/pull/5935) par @leaanthony

## Modifié

- L’animation de zoom de macOS utilise CADisplayLink ou NSTimer pour gagner en fluidité dans la [PR](https://github.com/wailsapp/wails/pull/5945) par @savely-krasovsky

## Corrigé

- Configuration du projet Xcode iOS afin de conserver les indicateurs hérités de l’éditeur de liens et d’ajouter -ObjC dans la [PR](https://github.com/wailsapp/wails/pull/5915) par @mortenolsrud
- Correction du renouvellement excessif des connexions TCP dans le proxy de ressources `wails3 dev` avec les interfaces volumineuses, qui pouvait épuiser les ports éphémères de l’hôte et provoquer l’échec de processus sans rapport avec l’application avec `EADDRNOTAVAIL`
- Mise en file d’attente du JavaScript des événements de chaque fenêtre afin d’assurer un envoi ordonné et la contre-pression dans la [PR](https://github.com/wailsapp/wails/pull/5934) par @leaanthony

### Ajouté

- Implémentation d’une boîte aux lettres FIFO asynchrone générique pour l’envoi ordonné des événements dans la [PR](https://github.com/wailsapp/wails/pull/5851) par @savely-krasovsky et @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Ajouté

- Implémentation d’un stockage borné côté hôte pour les événements surdimensionnés et d’un envoi JavaScript ordonné dans la [PR](https://github.com/wailsapp/wails/pull/5930) par @leaanthony
- Implémentation du rebond de l’icône dans le Dock macOS lors du clignotement d’une fenêtre dans la [PR](https://github.com/wailsapp/wails/pull/5921) par @julianstorer

## Corrigé

- Lors du vidage, le serveur de ressources conserve les erreurs de détection du type de contenu et les préfixes non écrits dans la [PR](https://github.com/wailsapp/wails/pull/5931) par @leaanthony
- Prévention du plantage des applications macOS lors du remplacement du menu de l’application depuis un rappel Wails
- Correction des menus natifs illisibles sous Windows 10 1809 / Windows Server 2019 (build 17763). Les exportations uxtheme du mode sombre étaient conditionnées au build 18334 ; l’activation explicite du mode sombre au niveau de l’application ne s’exécutait donc jamais sur ces hôtes : l’arrière-plan du menu était peint en sombre, mais Windows continuait à dessiner le texte du menu avec le thème clair, produisant du texte sombre sur un arrière-plan sombre. Les ordinaux existent à partir de 17763 ; la condition correspond donc désormais à cette version.
- Correction de `w32.GetStockObject`, qui appelait `GetDeviceCaps` au lieu de `GetStockObject` et renvoyait ainsi 0 pour chaque objet prédéfini.
- Amélioration de la gestion et du signalement des erreurs de téléchargement du programme d’amorçage WebView2 dans la [PR](https://github.com/wailsapp/wails/pull/5924) par @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Corrigé

- L’activation des applications macOS respecte désormais la politique d’activation uniquement pour les applications ordinaires dans la [PR](https://github.com/wailsapp/wails/pull/5897) par @julianstorer
- Protection des fenêtres GTK non initialisées dans les builds Linux dans la [PR](https://github.com/wailsapp/wails/pull/5898) par @julianstorer
- Définition explicite d’une couleur d’arrière-plan opaque pour les fenêtres WebKit sous Linux avant le chargement de l’URL dans la [PR](https://github.com/wailsapp/wails/pull/5899) par @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Modifié

- Les tâches de build Android ciblent arm64 par défaut et deploy-emulator sélectionne l’architecture de l’hôte dans la [PR](https://github.com/wailsapp/wails/pull/5890) par @mortenolsrud

## Corrigé

- Conservation de l’état de zoom des fenêtres macOS pendant leur déplacement et réduction des animations dans la [PR](https://github.com/wailsapp/wails/pull/5900) par @leaanthony
- Correction du build Windows en mode serveur par l’ajout de `!server` à la contrainte de build `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Ajouté

- Documentation de l’achèvement de la vérification bêta de la phase 10 dans les détails d’implémentation, dans la [PR](https://github.com/wailsapp/wails/pull/5881) par @leaanthony

## Corrigé

- Transmission du handle de fenêtre à l’API du mode sombre de Windows et validation des arguments dans la [PR](https://github.com/wailsapp/wails/pull/5877) par @leaanthony
- Centralisation de la résolution de l’état des boutons de la barre de titre des fenêtres macOS sans cadre dans la [PR](https://github.com/wailsapp/wails/pull/5870) par @taliesin-ai
- Prévention de l’affichage illisible du texte des menus natifs lorsqu’une application Windows demande le mode sombre alors que le thème des applications Windows est clair. Le menu utilise désormais l’arrière-plan natif clair correspondant jusqu’à ce que Windows puisse afficher le texte du menu adapté au mode sombre.
- Correction des menus natifs illisibles sous Windows 10 1809 / Windows Server 2019 (build 17763). Les exportations uxtheme du mode sombre étaient conditionnées au build 18334 ; l’activation explicite du mode sombre au niveau de l’application ne s’exécutait donc jamais sur ces hôtes : l’arrière-plan du menu était peint en sombre, mais Windows continuait à dessiner le texte du menu avec le thème clair, produisant du texte sombre sur un arrière-plan sombre. Les ordinaux existent à partir de 17763 ; la condition correspond donc désormais à cette version.

## v3.0.0-beta.2 - 2026-08-02

## Modifié

- Passage de la v3 de la phase alpha à la phase bêta
- Documentation des valeurs par défaut intelligentes de la zone de notification système et du masquage automatique des fenêtres contextuelles, avec une couverture des régressions pour la sélection du gestionnaire de clics (#5840).
- L’outil de mise à jour GitHub exclut par défaut les ressources des programmes d’installation Windows dans la [PR](https://github.com/wailsapp/wails/pull/5861) par @leaanthony
- Ajout de la prise en charge des angles arrondis, carrés ou de rayon personnalisé pour les fenêtres macOS sans cadre dans la [PR](https://github.com/wailsapp/wails/pull/5866) par @leaanthony

## Corrigé

- Signalement des dimensions réelles des fenêtres GTK4 et émission, depuis la surface configurée, des événements de redimensionnement et de changement d’état d’agrandissement, de réduction et de plein écran (#5830).
- Correction du plantage de WebKit sous Linux lors de l’envoi de Blob ou de FormData dans des requêtes fetch dans la [PR](https://github.com/wailsapp/wails/pull/5854) par @taliesin-ai
- Le shim de fetch transmet undefined lorsque les en-têtes Blob/FormData sont absents dans la [PR](https://github.com/wailsapp/wails/pull/5865) par @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Ajouté

## Modifié

- Ajout de la prise en charge des angles arrondis, carrés ou de rayon personnalisé pour les fenêtres macOS sans cadre dans la [PR](https://github.com/wailsapp/wails/pull/5866) par @leaanthony

## Corrigé

- Le shim de fetch transmet undefined lorsque les en-têtes Blob/FormData sont absents dans la [PR](https://github.com/wailsapp/wails/pull/5865) par @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Ajouté

- Ajout de la prise en charge de la création de paquets DMG pour macOS, avec de nouvelles options et tâches de compilation, dans la [PR](https://github.com/wailsapp/wails/pull/5857) par @leaanthony

## Modifié

- Le programme de mise à jour GitHub exclut désormais par défaut les ressources des programmes d’installation Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5861) par @leaanthony

## Corrigé

- Correction du plantage de WebKit sous Linux lors de l’envoi d’un Blob ou de FormData dans des requêtes fetch, dans la [PR](https://github.com/wailsapp/wails/pull/5854) par @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Ajouté

- Implémentation des actions de double-clic sur la barre de titre de macOS permettant d’agrandir ou de réduire la fenêtre, dans la [PR](https://github.com/wailsapp/wails/pull/5853) par @taliesin-ai

## Modifié

- Mise à jour du tutoriel du service QR pour utiliser NewServiceWithOptions et ajout d’espacement, dans la [PR](https://github.com/wailsapp/wails/pull/5849) par @jeongkyu

## Corrigé

- Maintien de la réactivité de WKWebView pendant le zoom sous macOS, dans la [PR](https://github.com/wailsapp/wails/pull/5856) par @leaanthony
- Correction des requêtes de taille de fenêtre GTK4 et émission des événements de redimensionnement, d’agrandissement, de réduction et de changement d’état du plein écran depuis le `GdkSurface` configuré.

## v3.0.0-alpha2.119 - 2026-07-27

## Corrigé

- Mise à jour de la documentation en plusieurs langues afin d’y inclure des diagrammes d’architecture, dans la [PR](https://github.com/wailsapp/wails/pull/5833) par @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Ajouté

- Ajout de chemins par défaut pour les fichiers d’entrée et de sortie de la génération d’icônes, dans la [PR](https://github.com/wailsapp/wails/pull/5825) par @taliesin-ai
- Ajout des modules d’entrée sources à sideEffects dans le package.json du paquet d’exécution, dans la [PR](https://github.com/wailsapp/wails/pull/5797) par @savely-krasovsky
- Ajout d’une section sur la licence et la provenance au guide de contribution, dans la [PR](https://github.com/wailsapp/wails/pull/5816) par @taliesin-ai

## Corrigé

- Application de CSS GTK4 limité aux fenêtres sans cadre afin de supprimer l’arrondi des bordures, dans la [PR](https://github.com/wailsapp/wails/pull/5800) par @savely-krasovsky
- Gestion sans erreur fatale des échecs de récupération de la position du curseur sous Windows pour les menus contextuels et l’énumération des écrans, dans la [PR](https://github.com/wailsapp/wails/pull/5789) par @wayneforrest
- Sous macOS, la boîte de dialogue d’ouverture de fichier filtre désormais correctement les extensions et valide les fichiers autorisés selon leur suffixe, dans la [PR](https://github.com/wailsapp/wails/pull/5678) par @phergul
- Protection de l’initialisation du mode sombre de Windows contre les appels d’API nil, dans la [PR](https://github.com/wailsapp/wails/pull/5793) par @roachadam
- Correction de l’échec de compilation en 32 bits du programme de mise à jour : la constante `maxArchiveTotalSize` (2 Gio) dépassait la capacité du type `int` de la plateforme lorsqu’elle était transmise à `fmt.Errorf` sur `GOARCH=386`. Elle est désormais explicitement typée `int64`.
- Correction d’une panique due à un pointeur nil au démarrage lorsqu’une fenêtre utilise une barre de titre sombre, ou sombre selon le système, sur les versions de Windows qui ne chargent pas les API uxtheme du mode sombre, telles que Windows 10 1809 / Windows Server 2019 (build 17763). Les appels à `AllowDarkModeForWindow` lors de la configuration du thème de la fenêtre sont désormais protégés contre les valeurs nil, comme l’est déjà `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Ajouté

- Implémentation d’une logique personnalisée de test de position pour les zones non clientes sous Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5462) par @savely-krasovsky

## Modifié

- Configuration de la détection de la mise à l’échelle du moniteur par WebView2 en fonction de UseVisualHosting, dans la [PR](https://github.com/wailsapp/wails/pull/5761) par @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Ajouté

- Mise à jour de la FAQ pour la recentrer sur les fonctionnalités et les recommandations de Wails v3, dans la [PR](https://github.com/wailsapp/wails/pull/5763) par @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Corrigé

- Correction de `Menu.Update()` qui ne reconstruisait pas le menu natif sous Linux avec GTK4 (#5659, problème diagnostiqué et corrigé indépendamment par @puneetdixit200 dans #5539)
- Correction du plantage lors de l’énumération des écrans macOS après un changement d’affichage, en copiant les chaînes d’identifiant et de nom des écrans et en prenant un instantané de leur nombre (#5565, problème diagnostiqué et corrigé indépendamment par @x-haose dans #5584)
- Correction d’un plantage sous Windows lorsque `WM_ERASEBKGND` peint un arrière-plan uni pendant une transition de réduction ou de restauration au cours de laquelle `GetClientRect` renvoie nil (protection signalée par @sinspired dans #5636)
- Correction de l’erreur des liaisons frontend, qui était toujours analysée comme du texte, par @mbaklor dans #5690
- Correction de l’échec de compilation sous Windows lors de l’utilisation de l’étiquette de compilation `server`, dû à l’absence, dans les fichiers d’interface graphique Windows, de la contrainte de compilation `!server` déjà présente dans leurs équivalents macOS et Linux (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Ajouté

- Implémentation du protocole Update Manifest et du fournisseur de point de terminaison, dans la [PR](https://github.com/wailsapp/wails/pull/5720) par @taliesin-ai

## Modifié

- Intégration de la liaison `webview2` au module v3 sous le nom `v3/internal/webview2`, avec suppression du module autonome, de ses workflows nocturnes de publication et de synchronisation, ainsi que des changements successifs de version dans go.mod, puisque v3 en est le seul consommateur, dans la [PR](https://github.com/wailsapp/wails/pull/5711) par @taliesin-ai

## Corrigé

- Déplacement de la détection de la mise à l’échelle du moniteur par WebView2 et du correctif de resynchronisation de l’hôte lors d’un changement de DPI vers la section « Non publié », dans la [PR](https://github.com/wailsapp/wails/pull/5750) par @taliesin-ai
- Mise à jour du marshaling COM de WebView2 pour les paramètres float64 et BOOL, dans la [PR](https://github.com/wailsapp/wails/pull/5741) par @wayneforrest
- Prévention d’une panique et d’un déréférencement nil lors de la mise à jour et de la destruction de l’icône de la zone de notification Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5703) par @wayneforrest
- Correction des fenêtres masquées qui ne se remasquaient pas correctement sous Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5743) par @wayneforrest
- Synchronisation de la visibilité du contrôleur WebView2 avec la réduction, l’agrandissement et la restauration de la fenêtre, dans la [PR](https://github.com/wailsapp/wails/pull/5742) par @wayneforrest

### Corrigé

- Réactivation de la détection de la mise à l’échelle du moniteur par WebView2 et exécution de la resynchronisation de l’hôte uniquement lors d’un changement de DPI, dans la [PR](https://github.com/wailsapp/wails/pull/5734) par @taliesin-ai, à partir du correctif validé par @randalmurphal, avec vérification de la cause première par @eleclin et tests matériels par @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Ajouts

- Ajout d’un avertissement lors de la création d’un AAB de publication sans `ANDROID_KEYSTORE_FILE` défini (Google Play rejette les bundles signés pour le débogage), ainsi que de la documentation sur la création et la signature des App Bundles dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai
- Ajout de la documentation Why Wails dans plusieurs langues dans la [PR](https://github.com/wailsapp/wails/pull/5739) par @taliesin-ai
- Prise en charge du mappage de Go time.Time vers JS Date ou string dans les bindings dans la [PR](https://github.com/wailsapp/wails/pull/5398) par @fbbdev
- Ajout de tâches de création de paquets Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) pour leur soumission au Play Store — les tâches APK restent disponibles pour les tests locaux ou sur émulateur — dans la [PR](https://github.com/wailsapp/wails/pull/5728) par @mortenolsrud (corrige [#5726](https://github.com/wailsapp/wails/issues/5726))
- Ajout de cibles de tâches pour les appareils Android physiques et reprise des autorisations d’accès à la caméra et à la localisation dans la [PR](https://github.com/wailsapp/wails/pull/5735) par @taliesin-ai

## Modifications

- Mise à niveau de `webview2` vers la version v1.0.28 ([notes de publication](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Mise à niveau du modèle Android `compileSdk`/`targetSdk` de 34 vers 35, version exigée par Google Play pour les nouvelles soumissions d’applications, dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai

## Corrections

- Correction des masques d’avatars préintégrés dans sponsorkit dans la [PR](https://github.com/wailsapp/wails/pull/5745) par @leaanthony
- Correction de la création automatique d’AVD Android, qui sélectionnait une image système ou une version de cmdline-tools incorrecte en raison du tri lexicographique des versions, dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai
- Correction de l’assistant de configuration, qui suggérait une version obsolète du NDK Android (désormais 26.3.11579264, conformément à l’exigence documentée), dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai
- Mise à jour de la documentation française relative à SvelteKit et aux options dans la [PR](https://github.com/wailsapp/wails/pull/5744) par @leaanthony
- Correction d’une erreur SIGSEGV lors de l’énumération des écrans sous macOS pendant les changements d’affichage dans la [PR](https://github.com/wailsapp/wails/pull/5516) par @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Ajouts

- Ajout d’un générateur SVG des contributeurs écrit en Go et mise à jour des pages de crédits de la documentation et du site web dans la [PR](https://github.com/wailsapp/wails/pull/5724) par @taliesin-ai
- Prise en charge du mappage de Go time.Time vers JS Date ou string dans les bindings dans la [PR](https://github.com/wailsapp/wails/pull/5398) par @fbbdev

## Modifications

- Remplacement du pipeline d’images des sponsors basé sur Node par un générateur Go dans la [PR](https://github.com/wailsapp/wails/pull/5719) par @taliesin-ai

## Corrections

- Correction du script d’installation des dépendances des ressources de build Android dans la [PR](https://github.com/wailsapp/wails/pull/5729) par @taliesin-ai
- Rejet du caractère de contrôle U+0085 (NEXT LINE) dans `ValidateAndSanitizeURL`, ce qui complète la couverture des caractères d’espacement par le validateur d’URL
- Recalcul du cadre DWM lors d’un changement de DPI pour les fenêtres sans cadre dans la [PR](https://github.com/wailsapp/wails/pull/4785) par @leaanthony
- Correction de la détection des zones de dépôt par glisser-déposer, qui échouait sous Windows avec une mise à l’échelle différente de 100 %, dans la [PR](https://github.com/wailsapp/wails/pull/4632) par @yulesxoxo
- Ajout d’une gestion explicite de la mémoire Objective-C pour les objets Cocoa utilisés dans les boîtes de dialogue, menus, icônes de barre d’état et notifications sous Darwin, dans la [PR](https://github.com/wailsapp/wails/pull/5714) par @taliesin-ai
- Correction de bogues du backend CGO sous Linux et de problèmes liés à l’icône de la zone de notification dans la [PR](https://github.com/wailsapp/wails/pull/5718) par @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Ajouts

- Ajout de HappyTools à la vitrine de la communauté dans la [PR](https://github.com/wailsapp/wails/pull/5061) par @Aliuyanfeng
- Ajout de la prise en charge de la langue indonésienne et d’une documentation complète dans la [PR](https://github.com/wailsapp/wails/pull/5643) par @triadmoko
- Ajout de l’option DisableMenu à WindowsWindow dans la [PR](https://github.com/wailsapp/wails/pull/4813) par @leaanthony

## Modifications

- Mise à jour du modèle Taskfile et de la CLI afin de répartir les tâches de build et de création de paquets selon GOOS et ARCH dans la [PR](https://github.com/wailsapp/wails/pull/5617) par @leaanthony

## Corrections

- Correction d’un problème lié au regroupement des fenêtres macOS en onglets dans la [PR](https://github.com/wailsapp/wails/pull/5708) par @taliesin-ai

## Suppressions

- Suppression des fichiers MDX traduits en allemand des sections sur la contribution, les fonctionnalités et les guides dans la [PR](https://github.com/wailsapp/wails/pull/5702) par @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Ajouts

- Implémentation du rechargement et du rechargement forcé de WebView sous macOS, ainsi que d’un mécanisme de récupération après l’arrêt du processus WebContent, dans la [PR](https://github.com/wailsapp/wails/pull/5129) par @wayneforrest
- Ajout d’une documentation complète en allemand sur la contribution, les fonctionnalités et les guides dans la [PR](https://github.com/wailsapp/wails/pull/5396) par @leaanthony
- Amélioration des notifications avec le son, les pièces jointes, la planification et une API de mise à jour des notifications, dans la [PR](https://github.com/wailsapp/wails/pull/5333) par @popaprozac

## Corrections

- Recalcul du cadre DWM lors d’un changement de DPI pour les fenêtres sans cadre dans la [PR](https://github.com/wailsapp/wails/pull/4785) par @leaanthony
- Correction de la détection des zones de dépôt par glisser-déposer, qui échouait sous Windows avec une mise à l’échelle différente de 100 %, dans la [PR](https://github.com/wailsapp/wails/pull/4632) par @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Ajouts

- Ajout d’exemples de code à la documentation d’EventsEmit dans la [PR](https://github.com/wailsapp/wails/pull/5026) par @iamhabbeboy
- Ajout d’une option d’hébergement visuel WebView2 sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/5380) par @MerIijn
- Ajout de Klustr à la documentation de la vitrine de la communauté dans la [PR](https://github.com/wailsapp/wails/pull/5536) par @SametKUM
- Ajout de Kira à la vitrine de la communauté, avec de nouvelles pages et une entrée dans le journal des modifications, dans la [PR](https://github.com/wailsapp/wails/pull/5685) par @thiennguyen93
- Ajout d’une section de commentaires au guide du service MCP dans la [PR](https://github.com/wailsapp/wails/pull/5694) par @taliesin-ai

## Modifications

- Le mode serveur dispose désormais d’une compilation de production de premier ordre, cohérente avec les tâches de compilation pour ordinateur de bureau (#5693). Par défaut, `task build:server` compile un binaire de production (`-tags server,production`, `-trimpath`, symboles supprimés) et accepte `DEV=true` (serveur de développement), `OBFUSCATED=true` (garble) et `EXTRA_TAGS`. `task run:server` exécute un serveur de développement. `Dockerfile.server` / `task build:docker` compilent d’abord le serveur de production (`-tags server,production`) et le frontend de production ; par défaut, l’image utilise une compilation statique en Go pur sur distroless/static, avec `CGO_ENABLED`, `GO_IMAGE` et `RUNTIME_IMAGE` exposés comme arguments de compilation substituables pour les applications CGO.

## Corrections

- Empêcher un plantage lors de la fermeture d’une fenêtre contenant des appels asynchrones en attente dans la [PR](https://github.com/wailsapp/wails/pull/4435) par @leaanthony
- Empêcher l’activation de la fenêtre lors de l’ouverture d’applications masquées sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/5249) par @leaanthony
- Garantir que les métadonnées des requêtes WebKit, l’achèvement des réponses et la gestion des flux de corps s’exécutent sur le thread principal GTK dans la [PR](https://github.com/wailsapp/wails/pull/5668) par @taliesin-ai
- Corriger `Menu.Update()`, qui ne reconstruisait pas le menu natif sous Linux avec GTK4 (#5659, problème diagnostiqué et corrigé indépendamment par @puneetdixit200 dans #5539)
- Corriger le plantage lors de l’énumération des écrans macOS après un changement d’affichage, en copiant les chaînes d’identifiant et de nom des écrans et en prenant un instantané de leur nombre (#5565, problème diagnostiqué et corrigé indépendamment par @x-haose dans #5584)
- Corriger le contenu WebView2 qui rétrécissait puis disparaissait après le déplacement d’une fenêtre entre des moniteurs utilisant des DPI différents sous Windows, en réappliquant les limites du contrôleur dans le gestionnaire `WM_DPICHANGED`, comme lors de la resynchronisation des DPI à la sortie de la réduction (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Ajouts

- Ajouter des raccourcis clavier globaux (à l’échelle du système) via `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Les raccourcis se déclenchent même lorsque l’application n’a pas le focus. Implémentation native propre à chaque plateforme, sans dépendance tierce : raccourcis Carbon sous macOS, `RegisterHotKey` sous Windows, `XGrabKey` sous X11 et interface de raccourcis globaux du portail de bureau XDG sous Wayland.
- Ajouter un serveur MCP intégré : un serveur Model Context Protocol qui démarre automatiquement lorsque l’application est compilée avec la balise `mcp`, afin de permettre aux agents LLM de tester et de contrôler une application Wails en cours d’exécution — contrôle des fenêtres, inspection du DOM, évaluation de JavaScript, appels de méthodes liées, événements et simulation des entrées de souris et de clavier, représentées par un curseur animé à l’écran. Aucun code utilisateur n’est nécessaire : la balise `mcp` est ajoutée automatiquement par `wails3 build`/`wails3 dev` lorsque `WAILS_MCP=1` est défini. Configuration entièrement réalisée au moyen de variables d’environnement (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Corrections

- Corriger `Menu.Update()`, qui ne reconstruisait pas le menu natif sous Linux avec GTK4 (#5659, problème diagnostiqué et corrigé indépendamment par @puneetdixit200 dans #5539)
- Corriger le plantage lors de l’énumération des écrans macOS après un changement d’affichage, en copiant les chaînes d’identifiant et de nom des écrans et en prenant un instantané de leur nombre (#5565, problème diagnostiqué et corrigé indépendamment par @x-haose dans #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Ajouts

- Ajouter la documentation expérimentale de Wake avec une navigation dans la barre latérale dans la [PR](https://github.com/wailsapp/wails/pull/5613) par @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Modifications

- Mettre à niveau `webview2` vers v1.0.27.
  - ci(webview2) : corriger la compilation de publication (compilation croisée pour Windows + go.sum complet) (#5671)\

  **Diff complète :** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Retirer go vet de la compilation croisée du workflow de publication de webview2 dans la [PR](https://github.com/wailsapp/wails/pull/5672) par @taliesin-ai
- Mettre à jour le modèle OpenRouter du journal des modifications automatique vers google/gemini-2.5-flash-lite dans la [PR](https://github.com/wailsapp/wails/pull/5670) par @taliesin-ai
- Mettre à niveau `webview2` vers v1.0.26.

### Corrections

- **Récupérer après des erreurs COM transitoires à l’exécution au lieu de quitter** (#5658, #5580). Auparavant, `Chromium.errorCallback` appelait `os.Exit(1)` pour *toute* erreur COM ; un incident récupérable après le démarrage arrêtait donc toute l’application. Les chemins d’exécution (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) consignent désormais l’erreur et récupèrent. En particulier, un message web mal formé ou non fiable dans `MessageReceived` est maintenant ignoré au lieu d’arrêter le processus. Cela corrige la catégorie de plantages lors du passage entre des moniteurs utilisant des DPI différents (#5544, #5650). Les erreurs dans les chemins de création de l’environnement ou du contrôleur restent fatales.\

**Diff complète :** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Corrections

- Corriger le workflow release-webview2 afin qu’il gère correctement les fichiers go.sum dans la [PR](https://github.com/wailsapp/wails/pull/5671) par @taliesin-ai
- Corriger les mises à jour des menus Linux avec GTK4 en effaçant puis en reconstruisant le menu natif dans la [PR](https://github.com/wailsapp/wails/pull/5659) par @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Ajouts

- Ajouter `application.System` pour détecter la plateforme à l’exécution depuis du code partagé : `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (la balise de compilation `server`) et `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` pour tester directement une cible unique. Il se compile pour chaque cible, ce qui permet de créer des branches conditionnelles sans balises de compilation. Des fonctions auxiliaires équivalentes pour le frontend (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) sont disponibles dans `@wailsio/runtime`
- Ajouter un guide « Utiliser d’autres frameworks frontend » expliquant comment placer votre propre projet Vite dans `frontend/` (couvre Solid, Preact, Lit, SvelteKit, Qwik, Angular, etc.)
- L’assistant `wails3 setup` vérifie désormais la chaîne d’outils mobiles (iOS/Android) — Xcode et l’environnement d’exécution du simulateur iOS, le JDK, le SDK/NDK Android et l’émulateur — avec, le cas échéant, une installation en un clic et des correctifs de configuration du shell pouvant être copiés
- Les projets générés incluent un fichier `frontend/.npmrc` qui définit le paramètre `minimum-release-age` à 7 jours afin de réduire l’exposition aux paquets récemment publiés (et potentiellement compromis). Ce paramètre est respecté par pnpm et bun, et ignoré sans conséquence par npm

## Modifications

- Repenser tous les modèles de démarrage intégrés avec une nouvelle apparence de bannière principale représentant une montagne au néon (web, iOS et Android)
- **TypeScript est désormais le langage par défaut des modèles de démarrage et utilise le nom de modèle sans suffixe.** `wails3 init` (sans `-t`) génère la structure d’un projet TypeScript ; `-t vanilla`, `-t react`, `-t vue` et `-t svelte` utilisent TypeScript, tandis que leurs variantes JavaScript sont disponibles sous `-t vanilla-js`, `-t react-js`, `-t vue-js` et `-t svelte-js`. Les modèles intégrés déclarent leur langage avec `typescript:` dans `template.yaml` ; les modèles de la communauté utilisant le suffixe `-ts` continuent de fonctionner comme solution de repli
- Repenser l’assistant `wails3 setup` avec le thème « Wails numérique » au néon (effet translucide de verre dépoli sur fond montagneux)

## Corrigé

- Correction d’un plantage sous Windows lors de la restauration d’une application restée réduite assez longtemps pour que WebView2 soit suspendu ou que son processus de rendu/GPU soit recyclé. La resynchronisation du DPI lors de la réduction/restauration (#5544) ne touche désormais le contrôleur WebView2 que si le DPI de la fenêtre a réellement changé, ce qui évite les appels COM fatals vers un contrôleur suspendu dans le cas courant d’une restauration sans changement de DPI (#5605)
- Correction de plantages natifs `SIGABRT`/`SIGSEGV` répétés (généralement dans `g_object_unref` pendant la boucle principale GTK) dans les applications Linux exécutées longtemps et soumises à de fréquents chargements de ressources ou de médias. Le serveur de ressources finalisait les `WebKitURISchemeRequest` depuis des goroutines de travail et appelait ainsi des fonctions WebKit2GTK non sûres entre threads en dehors du thread principal GTK ; la finalisation (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) s’exécute désormais dans le thread principal. Cette modification complète le correctif partiel de #5566. Elle concerne les builds GTK3 et GTK4/WebKitGTK 6.0 (#5631, #5557)
- Correction d’un `fatal error: invalid pointer found on stack` intermittent dans `setupSignalHandlers` sous Linux/GTK3. Les identifiants de fenêtre transmis comme `user_data` de signal étaient conservés dans une variable locale Go de type `unsafe.Pointer` ; le ramasse-miettes interrompait donc l’exécution lorsqu’il analysait cette valeur (qui n’est pas un pointeur) pendant une copie de pile. Côté Go, l’identifiant est désormais conservé sous forme d’entier (`uintptr_t`), ce qui rétroporte vers l’ancien chemin GTK3 le même correctif que #4958 avait appliqué au chemin GTK4 (où les fonctions de signal C avaient été converties en `uintptr_t` afin d’éliminer les erreurs `-race`/checkptr) (#5631)

## Supprimé

- Suppression des modèles de démarrage `react-swc`, `preact`, `lit`, `solid`, `qwik` et `sveltekit`, ainsi que de leurs variantes `-ts`. L’ensemble intégré pris en charge comprend désormais `vanilla`, `react`, `vue` et `svelte`, tous en TypeScript par défaut, avec des variantes JavaScript `-js`. Tout autre framework peut toujours être utilisé en [fournissant votre propre frontend](https://v3.wails.io/guides/dev/frontend-frameworks) ou au moyen d’un modèle personnalisé

## v3.0.0-alpha2.104 - 2026-06-18

## Corrigé

- Correction d’un plantage iOS (SIGABRT) lorsqu’une méthode liée d’un service Go renvoie une chaîne vide. Le générateur de réponses de ressources iOS vérifiait le pointeur du corps avec `buf != nil` au lieu de vérifier sa longueur ; un corps de longueur nulle provoquait donc une panique de `&buf[0]`. Il vérifie désormais la longueur, comme les générateurs de réponses pour ordinateur

## v3.0.0-alpha2.103 - 2026-06-15

## Modifié

- Déplacement des fonctionnalités natives iOS et Android vers les gestionnaires de plateforme : appelez-les via `application.IOS.*` et `application.Android.*` (par exemple `application.IOS.Haptic("medium")` et `application.Android.Share(payload)`) plutôt qu’au moyen des anciennes fonctions libres `application.IOS*`/`application.Android*` (#5602)
- Renommage des événements du pont mobile : les événements multiplateformes utilisent désormais le préfixe `common:*` (par exemple `common:haptic` et `common:location`), tandis que les événements propres à une plateforme utilisent `ios:*` / `android:*` (par exemple `ios:backgroundTask` et `android:foregroundService`) ; le préfixe `native:*` n’est plus utilisé (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Ajouté

- Ajout d’un assistant `wails3 setup` expérimental pour la configuration interactive des projets et la vérification des dépendances
- Ajout de l’option `--json` à `wails3 doctor` pour produire une sortie lisible par une machine
- Ajout d’une section sur l’état de la signature à la commande `wails3 doctor`

## Corrigé

- Correction de la détection de npm sous Linux afin de vérifier le PATH en plus du gestionnaire de paquets

## v3.0.0-alpha.101 - 2026-06-13

## Ajouté

- iOS : boîtes de dialogue natives pour les messages (UIAlertController) et boîtes de dialogue pour ouvrir un fichier, plusieurs fichiers ou un répertoire (UIDocumentPickerViewController) ; les boîtes de dialogue d’enregistrement renvoient une erreur explicite
- iOS : prise en charge du presse-papiers via UIPasteboard
- iOS : métriques d’écran réelles via UIScreen (points, pixels, facteur d’échelle et zone de travail tenant compte de la zone sûre)
- iOS : builds pour appareils (`IOS_PLATFORM=device`), prise en charge de l’identité de signature du code, du profil de provisionnement et des droits, empaquetage `.ipa` et `deploy-device` via devicectl
- iOS : version minimale d’iOS configurable (`ios.minIOSVersion` dans build/config.yml)
- iOS : `wails3 doctor` indique la disponibilité de Xcode et du SDK iOS sous macOS
- iOS : événements système — les informations sur la batterie, le réseau, le thème, le verrouillage de l’écran et le manque de mémoire sont exposées sous forme d’événements d’application `events.IOS.*` ainsi que d’événements d’application `events.Common.*` indépendants de la plateforme
- iOS : pont natif des fonctionnalités mobiles (`application.IOS*` exporté) — feuille de partage, ouverture d’URL, maintien de l’appareil éveillé, lampe torche, marges de la zone sûre, luminosité, informations sur l’application, verrouillage de l’orientation, barre d’état, biométrie (Face ID/Touch ID), notifications locales et stockage sécurisé dans le trousseau Keychain
- iOS : capteurs et matériel — retour haptique, géolocalisation ponctuelle, accéléromètre, proximité, synthèse vocale, informations de stockage, état de l’alimentation et de la batterie, état du réseau, marges du clavier et détection des captures d’écran
- iOS : documentation (IOS.md et guide sur le site de documentation)
- Android : boîtes de dialogue natives pour les messages (AlertDialog) et boîtes de dialogue pour ouvrir un ou plusieurs fichiers (Storage Access Framework, importés sous forme de copies dans le cache) ; les boîtes de dialogue d’ouverture de répertoire et d’enregistrement renvoient une erreur explicite
- Android : prise en charge du presse-papiers via ClipboardManager
- Android : métriques d’écran réelles via WindowMetrics/DisplayMetrics (dp, pixels, facteur d’échelle et zone de travail tenant compte des barres système)
- Android : méthodes du runtime pour le retour haptique (`Android.Haptics.Vibrate`), les informations sur l’appareil (`Android.Device.Info`) et les messages toast (`Android.Toast.Show`)
- Android : événements de cycle de vie typés (`events.Android.*`, générés depuis events.txt), avec `ActivityCreated` associé à `Common.ApplicationStarted`
- Android : la chaîne de build produit des APK de débogage et de publication installables (`android:run`, `android:package`, `android:package:fat`) ; la signature des versions publiées utilise par défaut le magasin de clés de débogage, ou un véritable magasin de clés fourni au moyen des variables d’environnement `ANDROID_KEYSTORE_*`
- Android : `wails3 doctor` indique les SDK Android, NDK et JDK
- Android : événements système — les informations sur la batterie, le réseau, le thème, le verrouillage de l’écran et le manque de mémoire sont exposées sous forme d’événements d’application `events.Android.*` ainsi que d’événements d’application `events.Common.*` indépendants de la plateforme
- Android : pont natif des fonctionnalités mobiles (`application.Android*` exporté) — partage, ouverture d’URL, maintien de l’appareil éveillé, lampe torche, marges de la zone sûre, luminosité, informations sur l’application, verrouillage de l’orientation, barre d’état, biométrie (BiometricPrompt), notifications locales et stockage sécurisé avec EncryptedSharedPreferences
- Android : capteurs et matériel — retour haptique, géolocalisation ponctuelle, accéléromètre, proximité, synthèse vocale, informations de stockage, état de l’alimentation et de la batterie, état du réseau, marges du clavier et blocage des captures d’écran avec FLAG_SECURE
- Android : documentation (ANDROID.md et guide sur le site de documentation)
- Exemple : l’exemple complet `mobile` s’enrichit d’onglets Mobile et Matériel qui présentent le pont des fonctionnalités natives sous iOS et Android (les onglets en forme de pilule passent sur plusieurs lignes)
- Mobile : batterie — l’accéléromètre, le capteur de proximité, la lampe torche et l’horloge périodique de l’exemple sont suspendus lorsque l’application passe en arrière-plan, puis rétablis à son retour (Android maintient le processus en cours d’exécution en arrière-plan, tandis que la lampe torche est un état matériel qui persiste sous iOS) ; sous Android, les récepteurs d’événements système ne sont enregistrés que lorsque l’application est au premier plan
- iOS : capture par caméra — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → un événement `native:capture` contenant une miniature encodée en base64)
- iOS : exécution en arrière-plan — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (une fenêtre d’exécution de tâche en arrière-plan UIApplication) et un paramètre `ios.backgroundModes` configurable (build/config.yml), qui permet au modèle de génération d’ajouter `UIBackgroundModes` au fichier Info.plist généré
- Android : capture par caméra — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (caméra système via FileProvider → un événement `native:capture`)
- Android : service de premier plan — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (un `WailsForegroundService` associé à une notification permanente maintient le processus actif pour les longues tâches en arrière-plan)
- Exemple : un onglet Caméra illustrant la capture de photos et de vidéos ainsi que l’exécution en arrière-plan (service de premier plan sous Android, fenêtre d’exécution de tâche en arrière-plan sous iOS)

## Corrections

- Correction de l’échec systématique de `getUserMedia` avec `NotAllowedError` sous Linux : WebKitGTK refuse les demandes d’autorisation que personne ne traite, et le signal `permission-request` n’était pas connecté. La caméra et le microphone sont désormais gérés au moyen d’une nouvelle table multiplateforme `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), prise en compte sous Linux (WebKitGTK) comme sous Windows (WebView2). Sous Linux, où aucune invite native n’existe, la caméra et le microphone sont autorisés par défaut (ce qui rétablit `getUserMedia`) et peuvent être désactivés avec `PermissionDeny` (#5552)
- iOS : `GOOS=ios` se compile de nouveau (export de `events.IOS`, stubs des noms de méthodes mobiles), tout comme les builds utilisant le tag de production (correction des tags de build dans pkg/application et plusieurs services)
- iOS : les événements Go→JS et ExecJS fonctionnent désormais — la page ne se charge plus deux fois au démarrage et la négociation `wails:runtime:ready` ne peut plus être perdue
- iOS : `ApplicationDidFinishLaunching`/`ApplicationStarted` n’entrent plus en concurrence avec le démarrage de l’application ; suppression de l’attente fixe de 2 secondes au démarrage
- iOS : correction d’une fuite de chaîne C à chaque exécution JavaScript de Go vers JS
- iOS : `hasListeners` reflète désormais l’enregistrement réel des écouteurs
- iOS : la journalisation de débogage du framework est exclue des builds de production lors de la compilation
- Android : `GOOS=android` se compile de nouveau — définition de `events.Android`, suppression du tableau d’écouteurs `events_android.go` hors limites, ajout du stub de nom de méthode mobile et arrêt de l’inclusion des fichiers Linux de bureau (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) dans les builds Android
- Android : les liaisons JS→Go fonctionnent désormais — la WebView ne pouvant pas transmettre les corps des requêtes POST `fetch()` à `shouldInterceptRequest`, les appels du runtime passent par un transport JavascriptInterface (`nativeHandleRuntimeCall`) au lieu de provoquer un plantage en raison d’un corps de requête nil
- Android : les appels `Screens.*` du runtime renvoient des données réelles — le ScreenManager est désormais alimenté au démarrage (il n’était jamais raccordé, si bien que `GetAll` renvoyait nil)
- Android : la journalisation de débogage du framework est exclue des builds de production lors de la compilation et passe par logcat sous le tag `Wails` dans les builds de débogage
- Android : véritable registre `hasListeners`, gestion des références et des exceptions JNI, et cycle de vie de page à chargement unique (sans double navigation)
- Correction de l’échec de `wails3 generate bindings` avec « Access is denied » sous Windows lorsque le serveur de développement Vite est actif, en synchronisant les fichiers générés dans le répertoire de sortie au lieu de renommer un répertoire par-dessus celui-ci (#5515)
- Correction d’un plantage fatal intermittent sous macOS lors de la lecture des informations d’écran après un changement d’affichage : l’identifiant et le nom de l’écran contenaient des pointeurs vers des tampons `UTF8String` à libération automatique, susceptibles d’être libérés avant leur copie par Go (utilisation après libération). Les chaînes font désormais l’objet d’un `strdup`, puis sont libérées après leur conversion ; l’énumération des écrans s’exécute dans un pool de libération automatique explicite afin de ne plus provoquer de fuite lorsqu’elle est appelée depuis des goroutines Go (#5556)
- Correction d’un SIGSEGV intermittent sous Linux lorsque l’assetserver ferme un `WebKitURISchemeRequest` : le dernier `g_object_unref` s’exécutait dans la goroutine de l’assetserver et finalisait ainsi un GObject WebKit hors du thread principal de GTK. L’opération unref est désormais transférée vers le contexte principal de GTK au moyen de `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Ajouts

- Extension de `MacWebviewPreferences` avec des options de configuration WKWebView supplémentaires : `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` et `ApplicationNameForUserAgent` (#5549)

## Corrections

- Correction de l’échec de `wails3 generate bindings` avec « Access is denied » sous Windows lorsque le serveur de développement Vite est actif, en synchronisant les fichiers générés dans le répertoire de sortie au lieu de renommer un répertoire par-dessus celui-ci (#5561)
- Correction des événements JS de redimensionnement qui ne se déclenchaient pas pour les fenêtres sans cadre sous Linux ; correction de la détection du bord de la barre de défilement pour les fenêtres sans cadre (#5368)
- Correction de l’échec du programme de mise à jour sous Windows avec « invalid cross-device link » lorsque le répertoire temporaire se trouve sur un volume différent de celui du répertoire d’installation (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Corrections

- Correction de l’échec de `wails3 generate bindings` avec « Access is denied » sous Windows lorsque le serveur de développement Vite est actif, en synchronisant les fichiers générés dans le répertoire de sortie au lieu de renommer un répertoire par-dessus celui-ci (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Corrections

- Correction du blocage de l’interface WebKit sous Linux durant les périodes d’inactivité (par exemple lorsque l’inspecteur est ouvert), en cessant d’imposer `SA_ONSTACK` à `SIGUSR1`, ce qui perturbait la synchronisation du thread du ramasse-miettes de JavaScriptCore (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Ajouts

- Ajout d’une page consacrée au débogage et à l’utilisation de `runtime/trace`

## Modifications

- Suppression de quelques importations `_ "embed"` inutiles afin de nettoyer légèrement le code

## Corrections

- Correction des contraintes de largeur et de hauteur minimales qui n’étaient plus appliquées après la restauration d’une fenêtre agrandie sous Windows (#4593)
- Correction de la traversée des clics de souris en mode plein écran avec les options de fenêtre Frameless + Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Ajouts

- Ajout de la prise en charge de l’obfuscation avec Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)) : identifiants stables des méthodes de liaison, intégration au build et au Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), ainsi que des tags de structure JSON sur chaque charge utile exposée au runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), afin que le format d’échange résiste au renommage des champs exportés effectué par Garble.

## v3.0.0-alpha.95 - 2026-05-20

## Ajouts

- Ajout de la page manquante sur la structure du projet

## Modifications

- Documentation : modification de deux diagrammes de la page consacrée à l’architecture afin d’utiliser des diagrammes de séquence pour un affichage plus clair
- Documentation : ajout d’une note indiquant que l’installation de D2 est requise au préalable pour l’exécution

## Corrections

- Correction de `wails3 generate appimage` avec GTK4 par défaut : l’outil de création de paquets détecte désormais la pile GTK à partir du binaire avant de rechercher les fichiers d’exécution. Il sélectionne ainsi `libwebkitgtkinjectedbundle.so` (sous `webkitgtk-6.0/`) pour les builds GTK4 et `libwebkit2gtkinjectedbundle.so` (sous `webkit2gtk-4.1/`) pour les builds `-tags gtk3`. La détection de `.relr.dyn` vérifie également `libgtk-4.so.1` afin que la suppression des symboles soit correctement désactivée sur les chaînes d’outils modernes, quelle que soit la pile. (#5475)
- Correction de l’échec de `wails3 generate appimage` lors de son invocation avec un `-builddir` relatif : l’outil de création de paquets résout désormais dès le départ `-binary`, `-icon`, `-desktopfile`, `-builddir` et `-outputdir` en chemins absolus, afin que le `s.CD` effectué en cours de traitement n’interrompe ni la goroutine de téléchargement d’AppRun ni la vérification effectuée avec `ldd` après la copie.
- Correction de l’échec de `wails3 generate appimage` lors du déplacement de l’AppImage finale vers `-outputdir` lorsque le champ `Name=` du fichier desktop ne correspond pas au nom de base du binaire : l’outil de création de paquets force désormais le plugin appimage de linuxdeploy, au moyen de la variable d’environnement `OUTPUT`, à écrire l’AppImage dans `<binary>-<arch>.AppImage` plutôt que sous le nom dérivé du fichier desktop.
- Correction de `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` et `Common.SystemDidWake` qui ne se déclenchaient plus sous Linux après que la pile GTK4 + WebKitGTK 6.0 est devenue la pile par défaut dans alpha.93. La nouvelle implémentation `application_linux.go` `run()` par défaut n’appelait ni `setupCommonEvents()`, qui transmet les événements `Linux.*` à leurs équivalents `Common.*`, ni `monitorPowerEvents()`. L’assistant DBus de surveillance de l’alimentation est désormais partagé entre les chemins de build GTK3 et GTK4 au moyen de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Corrections

- Correction de `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` et `Common.SystemDidWake` qui ne se déclenchaient plus sous Linux après que la pile GTK4 + WebKitGTK 6.0 est devenue la pile par défaut dans alpha.93. La nouvelle implémentation `application_linux.go` `run()` par défaut n’appelait ni `setupCommonEvents()`, qui transmet les événements `Linux.*` à leurs équivalents `Common.*`, ni `monitorPowerEvents()`. L’assistant DBus de surveillance de l’alimentation est désormais partagé entre les chemins de build GTK3 et GTK4 au moyen de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Ajouts

- Ajout de `XDG_SESSION_TYPE` à la sortie de `wails3 doctor` sous Linux par @leaanthony

## Corrections

- Correction du plantage du menu de fenêtre sous Wayland, causé par l’accès d’appmenu-gtk-module à une fenêtre non réalisée (#4769), par @leaanthony
- Correction du plantage de l’application GTK lorsque son nom contient des caractères non valides, tels que des espaces ou des parenthèses, par @leaanthony
- Correction de l’erreur « mémoire insuffisante » lors de l’initialisation du glisser-déposer sous Windows (#4701), par @overlordtm
- Correction d’une condition de concurrence dans le stockage des rappels du thread principal, due à l’utilisation incorrecte de RLock pour supprimer une entrée de la map (Linux, macOS et iOS) (#4424), par @leaanthony
- Correction de la gestion des variables lors de la transmission d’arguments de ligne de commande aux tâches. Les variables CLI indiquées sous forme de paires KEY=VALUE sont désormais correctement initialisées et propagées pendant toute l’exécution des tâches.
- Correction du conflit lié à NSWindowZoomButton sous macOS : `MaximiseButtonState` et `FullscreenButtonState` appliquent désormais l’état le plus restrictif au démarrage comme pendant l’exécution ; aucune des deux méthodes de définition ne peut plus remplacer silencieusement l’autre (#5319)
- Correction d’un ensemble de bogues préexistants dans l’ancien chemin de build GTK3 (`-tags gtk3`), détectés par CodeRabbit sur #5463 : les lancements par association de fichiers n’ignorent plus les gestionnaires de démarrage ; `getTheme` vérifie désormais les limites et les types ; `appName` ne libère plus la mémoire gérée par GLib ; `clipboardGet` ne provoque plus de fuite de l’objet `gchar*` renvoyé par GTK ; `Calloc` utilise désormais des récepteurs de type pointeur, et `NewCalloc` renvoie `*Calloc`, afin que le pool suive réellement les allocations ; `zoomOut` utilise l’inverse de `zoomInFactor` au lieu d’un multiplicateur négatif qui imposait la limite 1.0 ; `execJS` réutilise le nom de monde vide préalloué au lieu de provoquer la fuite d’un objet `C.CString("")` à chaque appel ; un `fmt.Println` de développement a été supprimé de `menuItem.setAccelerator`. Résout #5465.
- Correction de la même fuite liée à un récepteur par valeur `Calloc` dans le chemin de build GTK4 par défaut (`linux_cgo.go`) : l’emploi de récepteurs de type pointeur et de `NewCalloc() *Calloc` permet désormais de suivre et de libérer réellement les allocations `c.String(...)` propres à chaque fenêtre.

## v3.0.0-alpha.92 - 2026-05-15

## Ajouts

- Modification des Taskfiles afin de permettre le choix du gestionnaire de paquets frontend au moyen de l’option `PACKAGE_MANAGER`
- Enrichissement des données de modèle avec `{{.Opn}}` et `{{.Cls}}` afin de rendre l’écriture de modèles Taskfile plus prévisible

## Modifications

- Modification de certains Taskfiles existants afin qu’ils utilisent `{{.Opn}} and {{.Cls}}`

## Corrections

- Correction d’une erreur fatale `concurrent map read and map write` à l’exécution dans `linuxSystemTray` lorsque le menu de la zone de notification est mis à jour pendant que le panneau le lit.
- Utilisation de `log` au lieu de `fmt` pour la sortie des erreurs et traces de pile WebView2, afin que les messages ne soient pas perdus lorsque l’application s’exécute sans console attachée sous Windows.

## v3.0.0-alpha.91 - 2026-05-12

## Modifications

- Mise à jour du SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5414) par `@github-actions[bot]`
- **RUPTURE DE COMPATIBILITÉ (macOS) :** normalisation du système de coordonnées macOS afin que `GetScreens`, `Position` et `SetPosition` utilisent tous le même espace : des points logiques, avec l’axe Y orienté vers le bas et `(0,0)` placé dans l’angle supérieur gauche de l’écran principal. Cela correspond à Windows, à GTK ainsi qu’aux API publiques d’Electron et du Web. Les écrans physiquement situés au-dessus de l’écran principal signalent désormais une valeur `Bounds.Y` négative, auparavant positive, et les valeurs `Position()`/`SetPosition()` sont désormais exprimées en points logiques plutôt qu’en `points × primaryScale`. La conversion aller-retour `Position()` → `SetPosition()` reste préservée ; les valeurs absolues enregistrées par les builds alpha antérieurs et les contournements calculés manuellement, par exemple la multiplication par `primaryScale` ou l’inversion de Y par rapport à la hauteur d’un écran, devront être mis à jour. Résout [#5117](https://github.com/wailsapp/wails/issues/5117).

## Corrections

- Validation défensive du nom et de la longueur du corps des signaux DBus afin d’éviter les paniques dans [PR](https://github.com/wailsapp/wails/pull/5416) par @leaanthony
- Correction d’un problème de sécurité mémoire dans la gestion des menus GTK sous Linux dans [PR](https://github.com/wailsapp/wails/pull/5363) par @leaanthony
- Détection des GPU NVIDIA et désactivation du moteur de rendu DMA-BUF sous Linux dans [PR](https://github.com/wailsapp/wails/pull/5295) par @leaanthony
- Correction de la conversion Y entre écrans par `SetPosition` sous macOS : utilisation de la hauteur de l’écran principal comme référence globale afin que les fenêtres soient placées à la bonne position sur les moniteurs décalés verticalement par rapport à l’écran principal dans [#5117](https://github.com/wailsapp/wails/issues/5117)
- Correction du modèle de PR Git afin qu’il pointe vers la bonne URL de commentaires dans [PR](https://github.com/wailsapp/wails/pull/5109) par @wayneforrest
- Correction d’une série de plantages de la zone de notification Windows liés à `SetMenu`, provoqués par un appel système `DestroyMenu` défectueux qui transmettait quatre arguments au lieu d’un. Chaque appel renvoyait donc FALSE et ne libérait rien. Libération également des handles HMENU et HBITMAP (y compris ceux alloués à l’exécution via `MenuItem.SetBitmap`) lors de la reconstruction des menus, réinitialisation des tables obsolètes de cases à cocher et de boutons radio dans `Win32Menu.Update`, et suppression dans `systemtray.updateMenu` d’un appel `Update()` redondant qui doublait les allocations. Les applications qui s’exécutent longtemps dans la zone de notification ne présentent désormais plus de fuite d’objets GDI/USER à chaque reconstruction du menu.

## v3.0.0-alpha.90 - 2026-05-11

## Ajouts

- Ajout d’un nom d’application configurable pour l’agent utilisateur WKWebView sur macOS dans la [PR](https://github.com/wailsapp/wails/pull/5261) par @vinhvoit225
- Ajout de la dépendance indirecte github.com/coder/websocket à l’exemple gin-service dans la [PR](https://github.com/wailsapp/wails/pull/5400) par @taliesin-ai
- Ajout de la prise en charge de la comparaison par égalité profonde aux tests des ressources de compilation dans la [PR](https://github.com/wailsapp/wails/pull/5402) par @leaanthony

## Modifications

- Regroupement de la sortie de compilation dans le répertoire assets dans la [PR](https://github.com/wailsapp/wails/pull/5401) par @taliesin-ai
- Mise à jour du SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5399) par `@github-actions[bot]`

## Corrections

- Utilisation d’un objet de notification pour le message d’instance unique sur macOS dans la [PR](https://github.com/wailsapp/wails/pull/5289) par @overlordtm
- Regroupement par lots des fonctions de rappel Windows afin d’éviter la perte de promesses sous forte charge dans la [PR](https://github.com/wailsapp/wails/pull/5383) par @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Ajouts

- Ajout de la tâche go<em>test</em>results pour agréger les résultats des tests Go dans la [PR](https://github.com/wailsapp/wails/pull/5316) par @leaanthony

## Modifications

- Découpage conditionnel des charges utiles RPC volumineuses en requêtes POST fragmentées dans la [PR](https://github.com/wailsapp/wails/pull/5369) par @leaanthony
- Mise à niveau de Vite de 5.x.x vers 8.0.0 dans tous les modèles frontend dans la [PR](https://github.com/wailsapp/wails/pull/5386) par @leaanthony
- Migration de la configuration du port du serveur de développement Vite vers des variables d’environnement dans la [PR](https://github.com/wailsapp/wails/pull/5365) par @leaanthony
- Configuration du serveur de développement Vite pour qu’il écoute sur 127.0.0.1 dans tous les modèles, dans la [PR](https://github.com/wailsapp/wails/pull/5361) par @leaanthony
- Mise à jour du SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5384) par `@github-actions[bot]`

## Corrections

- Assainissement des ébauches de modèle Info.plist lors de la mise à jour de build-assets dans la [PR](https://github.com/wailsapp/wails/pull/5312) par @leaanthony
- Correction de l’état obsolète des menus macOS en appliquant de manière synchrone les mutateurs d’éléments de menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) sur le thread principal, ce qui élimine la condition de concurrence `dispatch_async` qui entraînait l’affichage de l’état précédent lors de la réouverture rapide des menus (#5002)
- Exclusion des fichiers `*_test.go` en mode développement afin d’éviter les recompilations inutiles dans la [PR](https://github.com/wailsapp/wails/pull/5203) par @leaanthony
- Prévention d’une erreur de segmentation dans Menu.Update() lorsque l’application n’est pas en cours d’exécution, dans la [PR](https://github.com/wailsapp/wails/pull/5291) par @wucm667
- Utilisation de lastSizeWParam pour conditionner le redessin de la barre de menus sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/5382) par @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Modifications

- Modification de HiddenOnTaskbar afin d’utiliser WS<em>EX</em>TOOLWINDOW dans la [PR](https://github.com/wailsapp/wails/pull/5371) par @leaanthony
- Réorganisation des dépendances et suppression de la directive replace de webview2 dans go.mod, dans la [PR](https://github.com/wailsapp/wails/pull/5370) par @atterpac
- Mise à jour du SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5358) par `@github-actions[bot]`

## Corrections

- Suppression des alias génériques d’indirection et regroupement des types de clés de tables dans la [PR](https://github.com/wailsapp/wails/pull/5331) par @fbbdev

## Suppressions

- Suppression du workflow PR-master et de ses tâches de documentation, de tests Go et de saut des tests, dans la [PR](https://github.com/wailsapp/wails/pull/5377) par @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Ajouts

- Ajout de la documentation en coréen pour Wails v3 dans la [PR](https://github.com/wailsapp/wails/pull/5352) par @leaanthony
- Ajout de la documentation en français sur l’installation et le démarrage rapide dans la [PR](https://github.com/wailsapp/wails/pull/5354) par @leaanthony
- Ajout de la documentation en portugais sur le démarrage rapide, les concepts et la communauté dans la [PR](https://github.com/wailsapp/wails/pull/5355) par @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Ajouts

- Ajout de la traduction française de la documentation dans la [PR](https://github.com/wailsapp/wails/pull/5328) par @leaanthony
- Ajout de la langue allemande au site de documentation dans la [PR](https://github.com/wailsapp/wails/pull/5343) par @leaanthony

## Modifications

- Enregistrement des 8 langues traduites dans la configuration de la documentation, dans la [PR](https://github.com/wailsapp/wails/pull/5347) par @leaanthony
- Mise à jour de divers fichiers liés à Windows pour WebView2 dans la [PR](https://github.com/wailsapp/wails/pull/5317) par @leaanthony

## Corrections

- Séparation de la répartition des boîtes de dialogue entre GTK3 et GTK4 sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/5340) par @leaanthony
- Exécution garantie des fonctions de rappel des boîtes de dialogue sur le thread GTK, ce qui corrige les erreurs de segmentation, dans la [PR](https://github.com/wailsapp/wails/pull/5339) par @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Ajouts

- Ajout de l’URL du modèle de PR au dépôt dans la [PR](https://github.com/wailsapp/wails/pull/5179) par @leaanthony
- Ajout de la documentation en allemand pour Wails v3 dans la [PR](https://github.com/wailsapp/wails/pull/5330) par @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Ajouts

- Ajout d’une option permettant de désactiver la sortie du mode plein écran avec la touche Échap sur macOS dans la [PR](https://github.com/wailsapp/wails/pull/5307) par @leaanthony
- Ajout d’une option permettant de désactiver la sortie du mode plein écran avec la touche Échap sur macOS dans la [PR](https://github.com/wailsapp/wails/pull/5310) par @leaanthony
- Ajout de la documentation de Pausa à la vitrine communautaire dans la [PR](https://github.com/wailsapp/wails/pull/5288) par @yuseferi

## Modifié

- Mise à jour du SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5308) par `@github-actions[bot]`
- Mise à jour de la commande de génération des icônes pour gérer les plateformes non prises en charge dans la [PR](https://github.com/wailsapp/wails/pull/5309) par @leaanthony
- Remplacement de l’API booléenne de plein écran par le type ternaire ButtonState et implémentation des liaisons propres aux plateformes dans la [PR](https://github.com/wailsapp/wails/pull/5224) par @leaanthony

## Corrigé

- Protection des opérations de gestion du focus de WebView2 contre un état de contrôleur nil dans la [PR](https://github.com/wailsapp/wails/pull/5315) par @leaanthony
- Mise à jour du workflow GitHub Actions afin qu’il référence correctement la branche de base de la PR dans la [PR](https://github.com/wailsapp/wails/pull/5313) par @leaanthony
- Fichiers `*_test.go` ignorés en mode développement afin d’éviter les reconstructions inutiles dans la [PR](https://github.com/wailsapp/wails/pull/5203) par @leaanthony
- Prévention d’une erreur de segmentation de Menu.Update() lorsque l’application n’est pas en cours d’exécution dans la [PR](https://github.com/wailsapp/wails/pull/5291) par @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Ajouté

- Ajout de l’indicateur InstallScope et d’une option de build pour une installation par machine ou par utilisateur dans la [PR](https://github.com/wailsapp/wails/pull/5094) par @symball
- Ajout à BrowserWindow d’une méthode SetScreen sans effet afin de satisfaire l’interface Window dans la [PR](https://github.com/wailsapp/wails/pull/5294) par @leaanthony

## Corrigé

- Détection des GPU NVIDIA et désactivation du moteur de rendu DMA-BUF sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/5295) par @leaanthony
- Correction du modèle de PR Git afin qu’il pointe vers la bonne URL de commentaires dans la [PR](https://github.com/wailsapp/wails/pull/5109) par @wayneforrest
- Correction d’une série de plantages `SetMenu` de la zone de notification Windows causés par un appel système `DestroyMenu` défectueux qui transmettait quatre arguments au lieu d’un, de sorte que chaque appel renvoyait FALSE et ne libérait rien. Libération également des handles HMENU et HBITMAP, y compris ceux alloués à l’exécution via `MenuItem.SetBitmap`, lors de la reconstruction des menus, réinitialisation des mappages obsolètes des cases à cocher et boutons radio dans `Win32Menu.Update`, et suppression dans `systemtray.updateMenu` d’un appel `Update()` redondant qui doublait les allocations. Les applications de zone de notification exécutées pendant de longues périodes ne provoquent désormais plus de fuite d’objets GDI/USER à chaque reconstruction de menu.

## v3.0.0-alpha.82 - 2026-05-01

## Corrigé

- Correction de la génération du fichier desktop afin de gérer correctement le nom utilisé dans ce fichier, dans la [PR](https://github.com/wailsapp/wails/pull/5232) par @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Modifié

- Ajustement à 15:00 UTC de l’horaire des versions nightly dans la [PR](https://github.com/wailsapp/wails/pull/5286) par @leaanthony

## Corrigé

- Correction des valeurs de Screen Bounds, WorkArea et Size réduites de moitié sur les Mac Retina -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Modifié

- Mise à jour des dépendances de la documentation et des chargeurs de collections de contenu dans la [PR](https://github.com/wailsapp/wails/pull/5285) par @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Ajouté

- Octroi de l’autorisation actions: write à la tâche trigger-release dans la [PR](https://github.com/wailsapp/wails/pull/5270) par @leaanthony

## Modifié

- La tâche de publication utilise désormais la branche master par défaut et met à jour la formulation du journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/5283) par @leaanthony
- Mise à jour du workflow auto-changelog afin d’utiliser la dernière version dans la [PR](https://github.com/wailsapp/wails/pull/5282) par @leaanthony
- Amélioration de l’efficacité des workflows par l’ajout de filtres de chemins et la suppression des workflows inutilisés dans la [PR](https://github.com/wailsapp/wails/pull/5280) par @leaanthony
- Mise à jour de la documentation afin que les liens des exemples fassent référence à la branche master dans la [PR](https://github.com/wailsapp/wails/pull/5274) par @leaanthony
- Mise à jour de la documentation et des exemples pour la v3 dans la [PR](https://github.com/wailsapp/wails/pull/5272) par @leaanthony

## Corrigé

- Amélioration du proxy inverse avec une logique de nouvelle tentative et l’utilisation forcée d’IPv4 pour le développement dans la [PR](https://github.com/wailsapp/wails/pull/5265) par @AkagiYui
- Réécriture du workflow de déclenchement du journal des modifications non publié dans la [PR](https://github.com/wailsapp/wails/pull/5281) par @leaanthony

## Supprimé

- Suppression des scripts de test shell destinés à différents types de tests dans la [PR](https://github.com/wailsapp/wails/pull/5267) par @leaanthony
- Suppression du workflow de déploiement de la documentation v3-alpha et de l’enregistrement CNAME dans la [PR](https://github.com/wailsapp/wails/pull/5266) par @leaanthony

### Ajouté

- Ajout de l’entrée Routage frontend à la navigation de la barre latérale dans la [PR](https://github.com/wailsapp/wails/pull/5196) par @leaanthony
- Ajout d’un guide sur le routage frontend avec des recommandations propres à chaque framework dans la [PR](https://github.com/wailsapp/wails/pull/5185) par @leaanthony
- Ajout de la prise en charge des feuilles modales (macOS)
- Mise à niveau de la version de ghw pour améliorer la prise en charge des appareils Apple par @leaanthony (#4977)
- Ajout de la méthode `GetBadge` au service du Dock
- Ajout de l’indicateur `-tags` à la commande `wails3 build` pour transmettre des tags de build Go personnalisés (par exemple, `wails3 build -tags gtk4`) (#4957)
- Ajout de la documentation sur la génération automatique des énumérations dans le générateur de liaisons, avec une page Énumérations dédiée et une entrée dans la barre latérale (#4972)
- Ajout de l’indicateur `-tags` à la commande `wails3 build` pour transmettre des tags de build Go personnalisés (par exemple, `wails3 build -tags gtk4`) (#4957)
- Ajout dans `v3/examples/web-apis/` d’exemples d’API Web illustrant 41 API de navigateur, notamment Stockage (localStorage, sessionStorage, IndexedDB, Cache API), Réseau (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Médias (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Appareil (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Performances (Performance API, Mutation Observer, Intersection/Resize Observer), Interface utilisateur (Web Components, Pointer Events, Selection, Dialog, Drag and Drop), et bien plus encore
- Ajout d’un exemple de vérificateur de compatibilité des API WebView (`v3/examples/webview-api-check/`) qui teste plus de 200 API de navigateur sur différentes plateformes
- Ajout du package `internal/libpath` pour rechercher les chemins des bibliothèques natives sous Linux, avec recherche parallèle, mise en cache et prise en charge de Flatpak/Snap/Nix
- **En cours :** ajout de la prise en charge expérimentale de WebKitGTK 6.0 / GTK4 sous Linux, disponible via `-tags gtk4` (GTK3/WebKit2GTK 4.1 reste utilisé par défaut)
- Remarque : avec les gestionnaires de fenêtres en mosaïque (par exemple Hyprland ou Sway), les opérations de réduction et d’agrandissement peuvent ne pas fonctionner comme prévu, car le gestionnaire de fenêtres contrôle la géométrie des fenêtres
- Ajout d’instructions pour mettre en place des **gestionnaires à exécution unique** dans la documentation sur l’**écoute des événements en JavaScript**, par @AbdelhadiSeddar
- Ajout à `WebviewWindowOptions` de l’option `UseApplicationMenu`, qui permet aux fenêtres sous Windows et Linux d’hériter du menu d’application défini via `app.Menu.Set()`, par @leaanthony
- Ajout de la prise en charge des fichiers `.icon` (format Apple Icon Composer) pour générer des icônes Liquid Glass et des catalogues de ressources sous macOS (#4934), par @wimaha
- Ajout d’un mode serveur expérimental pour les déploiements sans interface graphique ou sur le Web (`-tags server`). Il permet d’exécuter les applications Wails en tant que serveurs HTTP sans dépendances d’interface graphique natives. Compilez avec `wails3 task build:server`. Consultez `examples/server` pour plus de détails.
- Ajout du paquet `internal/libpath` pour rechercher les chemins des bibliothèques natives sous Linux, avec recherche parallèle, mise en cache et prise en charge de Flatpak, Snap et Nix
- Ajout à `MacWindow` de l’option `CollectionBehavior` pour contrôler le comportement des fenêtres entre les Spaces de macOS et le mode plein écran (#4756), par @leaanthony
- Ajout de tests unitaires pour pkg/application, par @leaanthony
- Ajout de la prise en charge des protocoles personnalisés à la création de paquets MSIX, par @leaanthony
- Ajout de la détection de l’environnement de bureau sous Linux, [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Ajout de la méthode `Window.Print()` à l’environnement d’exécution JavaScript pour ouvrir la boîte de dialogue d’impression depuis le frontend (#4290), par @leaanthony
- Ajout de `XDG_SESSION_TYPE` à la sortie de `wails3 doctor` sous Linux, par @leaanthony
- Ajout d’événements supplémentaires de changement de chargement WebKit2 sous Linux : `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896), par @leaanthony
- Ajout de `XDG_SESSION_TYPE` à la sortie de `wails3 doctor` sous Linux, par @leaanthony
- Génération du fichier `.desktop` pendant la compilation Linux, et non plus uniquement lors de la création du paquet (#4575)
- Ajout d’une documentation sur les dépendances d’exécution Linux, avec les noms de paquets propres à chaque distribution et des exemples de création de paquets nfpm (#4339), par @leaanthony
- Ajout des informations sur la version du pilote NVIDIA à la sortie de `wails3 doctor` sous Linux, par @leaanthony
- Ajout de l’origine au gestionnaire de messages bruts, par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4710)
- Ajout de la prise en charge des liens universels sous macOS, par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4712)
- Refactorisation de la couche de transport des liaisons, par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4702)
- Ajout d’identifiants aria-label aux modèles helloworld afin que l’application d’exemple puisse être facilement testée par des clients de test Appium, par @chinenual dans la [PR](https://github.com/wailsapp/wails/pull/4760)
- Ajout de l’origine au gestionnaire de messages bruts, par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4710)
- Ajout de la prise en charge des liens universels sous macOS, par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4712)
- Refactorisation de la couche de transport des liaisons, par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4702)
- Événements typés, par @fbbdev et @ianvs dans [#4633](https://github.com/wailsapp/wails/pull/4633)
- Ajout de l’exemple `systray-clock` illustrant une icône de zone de notification sans fenêtre, avec mise à jour en direct de l’infobulle (#4653).
- Ajout d’un modèle de protocole NSIS pour Windows, par @Tolfx dans #4510
- Ajout de tests pour build-assets, par @Tolfx dans #4510
- macOS : affichage des contrôles de fenêtre natifs dans la barre des menus, dans [#4588](https://github.com/wailsapp/wails/pull/4588), par @nidib
- Ajout d’un service Dock macOS permettant de masquer ou d’afficher l’icône de l’application dans le Dock, par @popaprozac dans la [PR](https://github.com/wailsapp/wails/pull/4451)
- Ajout d’un service Dock macOS permettant de masquer ou d’afficher l’icône de l’application dans le Dock, par @popaprozac dans la [PR](https://github.com/wailsapp/wails/pull/4451)
- Ajout de la prise en charge de l’effet Liquid Glass natif sous macOS avec NSGlassEffectView (macOS 15.0+) et repli sur NSVisualEffectView, ainsi que d’options complètes de personnalisation des matériaux, par @leaanthony dans [#4534](https://github.com/wailsapp/wails/pull/4534)
- Assainissement des URL du navigateur, par @leaanthony dans [#4500](https://github.dev/wailsapp/wails/pull/4500). Basé sur [#4484](https://github.com/wailsapp/wails/pull/4484) par @APShenkin.
- Ajout de la protection du contenu sous Windows et macOS par [@leaanthony](https://github.com/leaanthony), à partir du travail initial de [@Taiterbase](https://github.com/Taiterbase) dans cette [PR](https://github.com/wailsapp/wails/pull/4241)
- Ajout de la prise en charge de la transmission de variables CLI aux commandes Task via les alias `wails3 build` et `wails3 package` (#4422), par @leaanthony dans la [PR](https://github.com/wailsapp/wails/pull/4488)
- Prise en charge des zones de dépôt, avec des événements fournissant les données des éléments déposés, par [@atterpac](https://github.com/atterpac) dans [#4318](https://github.com/wailsapp/wails/pull/4318)
- Ajout de `AdditionalLaunchArgs` aux options de `WindowsWindow` afin de transmettre des arguments de ligne de commande supplémentaires au navigateur WebView2, dans la [PR](https://github.com/wailsapp/wails/pull/4467)
- Ajout de l’exécution automatique de go mod tidy après wails init, par [@triadmoko](https://github.com/triadmoko) dans la [PR](https://github.com/wailsapp/wails/pull/4286)
- Fonctionnalité Snap Assist de Windows, par @leaanthony dans la [PR](https://github.dev/wailsapp/wails/pull/4463)
- Ajout de `AdditionalLaunchArgs` aux options de `WindowsWindow` afin de transmettre des arguments de ligne de commande supplémentaires au navigateur WebView2, dans la [PR](https://github.com/wailsapp/wails/pull/4467)
- Ajout de l’exécution automatique de go mod tidy après wails init, par [@triadmoko](https://github.com/triadmoko) dans la [PR](https://github.com/wailsapp/wails/pull/4286)
- Fonctionnalité Snap Assist de Windows, par @leaanthony dans la [PR](https://github.dev/wailsapp/wails/pull/4463)
- Ajout de l’implémentation Windows de `getAccentColor` par [@almas-x](https://github.com/almas-x) dans la [PR](https://github.com/wailsapp/wails/pull/4427)
- Ajout de l’implémentation Windows de `getAccentColor` par [@almas-x](https://github.com/almas-x) dans la [PR](https://github.com/wailsapp/wails/pull/4427)
- Menus et barre des menus adaptés au thème sombre sous Windows. Par @leaanthony dans [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Renommage des services intégrés pour rendre les liaisons JS/TS plus explicites, par @popaprozac dans la [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` permet d’obtenir la couleur d’accentuation du système de l’utilisateur. Fonctionne sous macOS. Par [@etesam913](https://github.com/etesam913)
- Ajout de l’API `window.ToggleFrameless()` par [@atterpac](https://github.com/atterpac) dans [#4137](https://github.com/wailsapp/wails/pull/4137)
- Ajout des dépendances de compilation propres à chaque distribution pour Linux par @leaanthony dans la [PR](https://github.com/wailsapp/wails/pull/4345)
- Ajout d’un guide sur les bindings par @atterpac dans la [PR](https://github.com/wailsapp/wails/pull/4404)
- **Organisation de l’infrastructure de test** : déplacement des fichiers de test Docker vers le répertoire dédié `test/docker/`, avec des images optimisées et une fiabilité de compilation accrue, par [@leaanthony](https://github.com/leaanthony) dans [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Amélioration des modèles de gestion des ressources** : ajout, dans les exemples, d’un nettoyage approprié des gestionnaires d’événements et d’une gestion des goroutines tenant compte du contexte, par [@leaanthony](https://github.com/leaanthony) dans [#4359](https://github.com/wailsapp/wails/pull/4359)
- Prise en charge de la compilation d’AppImage pour aarch64 par [@AkshayKalose](https://github.com/AkshayKalose) dans [#3981](https://github.com/wailsapp/wails/pull/3981)
- Ajout d’une section de diagnostic à `wails doctor` par [@leaanthony](https://github.com/leaanthony)
- Ajout de la fenêtre au contexte lors de l’appel d’une méthode de service par [@leaanthony](https://github.com/leaanthony)
- Ajout de l’exemple `window-call` pour montrer comment déterminer quelle fenêtre appelle un service, par [@leaanthony](https://github.com/leaanthony)
- Nouveau guide sur les menus par [@leaanthony](https://github.com/leaanthony)
- Amélioration de la gestion des paniques par [@leaanthony](https://github.com/leaanthony)
- Nouveau guide sur les menus par [@leaanthony](https://github.com/leaanthony)
- Ajout de commentaires de documentation pour l’API Service par [@fbbdev](https://github.com/fbbdev) dans [#4024](https://github.com/wailsapp/wails/pull/4024)
- Ajout de la fonction `application.NewServiceWithOptions` pour initialiser les services avec une configuration supplémentaire, par [@leaanthony](https://github.com/leaanthony) dans [#4024](https://github.com/wailsapp/wails/pull/4024)
- Amélioration du contrôle des menus par [@FalcoG](https://github.com/FalcoG) et [@leaanthony](https://github.com/leaanthony) dans [#4031](https://github.com/wailsapp/wails/pull/4031)
- Ajout de documentation par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de l’annulation des événements dans les écouteurs d’événements standard par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de `Hide`, `Show` et `Destroy` pour la zone de notification, par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de `SetTooltip` pour la zone de notification, par [@leaanthony](https://github.com/leaanthony). Idée originale de [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Affichage du chemin du package dans les avertissements du générateur de bindings concernant les types non pris en charge, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la prise en charge des alias génériques dans le générateur de bindings par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la prise en charge de l’option JSON `omitzero` dans le générateur de bindings par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la directive `//wails:ignore` pour empêcher la génération de bindings pour certaines méthodes de service, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la directive `//wails:internal` sur les services et les modèles afin d’autoriser les types exportés en Go, mais pas en JS/TS, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la prise en charge, dans le générateur de bindings, des constantes dont le type est un alias afin de permettre les énumérations faiblement typées, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de tests du générateur de bindings pour les fonctionnalités de Go 1.24 par [@fbbdev](https://github.com/fbbdev) dans [#4068](https://github.com/wailsapp/wails/pull/4068)
- Ajout de la prise en charge de macOS 15 « Sequoia » à `OSInfo.Branding` afin d’améliorer la détection de la version du système d’exploitation dans [#4065](https://github.com/wailsapp/wails/pull/4065)
- Ajout du hook `PostShutdown` pour exécuter du code personnalisé une fois le processus d’arrêt terminé, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout de la structure `FatalError` pour permettre la détection des erreurs fatales dans les gestionnaires d’erreurs personnalisés, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Standardisation et documentation de l’ordre de démarrage et d’arrêt des services par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout d’un banc de test pour la séquence de démarrage et d’arrêt de l’application, ainsi que de tests de démarrage et d’arrêt des services, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout de la méthode `RegisterService` pour enregistrer des services après la création de l’application, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout du champ `MarshalError` aux options de l’application et des services pour personnaliser la gestion des erreurs lors des appels de bindings, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout d’un wrapper de promesse annulable qui propage les demandes d’annulation dans les chaînes de promesses, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Ajout de la possibilité de lier l’annulation d’un appel de binding à un `AbortSignal`, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Prise en charge des attributs `data-wml-*` pour WML, en plus des attributs `wml-*` habituels, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la méthode `Configure` à tous les services pour permettre une configuration tardive ou une reconfiguration dynamique, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Lorsqu’il n’est pas configuré, le service `fileserver` envoie une réponse 503 Service Unavailable, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Lorsqu’il n’est pas configuré, le service `kvstore` fournit par défaut un magasin clé-valeur en mémoire, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Load` au service `kvstore` pour recharger les données depuis le fichier après une modification de la configuration, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Clear` au service `kvstore` pour supprimer toutes les clés, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout du type `Level` au service `log` afin de fournir les constantes de niveau de journalisation côté JS, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Log` au service `log` pour définir dynamiquement le niveau de journalisation, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Lorsqu’il n’est pas configuré, le service `sqlite` fournit par défaut une base de données en mémoire, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Close` au service `sqlite` pour permettre la fermeture manuelle de la base de données, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la prise en charge de l’annulation pour les méthodes de requête du service `sqlite`, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la prise en charge des requêtes préparées au service `sqlite`, avec des liaisons JS, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Prise en charge de Gin par [Lea Anthony](https://github.com/leaanthony) dans la [PR](https://github.com/wailsapp/wails/pull/3537), basée sur le travail initial de [@AnalogJ](https://github.com/AnalogJ) dans cette [PR](https://github.com/wailsapp/wails/pull/3537)
- Correction de l’activation permanente de l’enregistrement automatique et de l’enregistrement automatique des mots de passe, par [@oSethoum](https://github.com/osethoum) dans [#4134](https://github.com/wailsapp/wails/pull/4134)
- Ajout de `SetMenu()` aux fenêtres afin de permettre la définition d’un menu sur une fenêtre, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la prise en charge des notifications par [@popaprozac](https://github.com/popaprozac) dans [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Ajout de la prise en charge des associations de fichiers sur mac, par [@wimaha](https://github.com/wimaha) dans [#4177](https://github.com/wailsapp/wails/pull/4177)
- Ajout de `wails3 tool version` pour l’incrémentation des versions sémantiques, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la prise en charge des badges sous macOS et Windows, par [@popaprozac](https://github.com/popaprozac) dans [#](https://github.com/wailsapp/wails/pull/4234)
- Ajout de la prise en charge des événements enregistrés et strictement typés, par [@fbbdev](https://github.com/fbbdev) et [@IanVS](https://github.com/IanVS) dans [#4161](https://github.com/wailsapp/wails/pull/4161)
- Ajout de la possibilité d’enregistrer des hooks pour les événements personnalisés, par [@fbbdev](https://github.com/fbbdev) et [@IanVS](https://github.com/IanVS) dans [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` pour ouvrir le gestionnaire de fichiers du système au chemin `path`, avec mise en surbrillance facultative via `selectFile`, par [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Nouvel indicateur `-git` pour la commande `wails3 init`, par [@leaanthony](https://github.com/leaanthony)
- Nouvelle commande `wails3 generate webview2bootstrapper`, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la méthode `init()` à l’environnement d’exécution pour permettre son initialisation manuelle, par [@leaanthony](https://github.com/leaanthony)
- Ajout de l’option `WindowDidMoveDebounceMS` à WindowOptions de Window, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la fonctionnalité d’instance unique par [@leaanthony](https://github.com/leaanthony). Basée sur la [PR v2](https://github.com/wailsapp/wails/pull/2951) de @APshenkin.
- Commande `wails3 generate template` par [@leaanthony](https://github.com/leaanthony)
- Commande `wails3 releasenotes` par [@leaanthony](https://github.com/leaanthony)
- Commande `wails3 update cli` par [@leaanthony](https://github.com/leaanthony)
- Option `-clean` pour la commande `wails3 generate bindings`, par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de la création d’AppImage Linux pour aarch64 (arm64), par [@AkshayKalose](https://github.com/AkshayKalose) dans [#3981](https://github.com/wailsapp/wails/pull/3981)
- Ajout d’un hyperlien pour le sponsor par @ansxuman dans [#3958](https://github.com/wailsapp/wails/pull/3958)
- Prise en charge de la création de paquets Linux deb, rpm et Arch Linux par
- Ajout de la prise en charge des builds et paquets universels Darwin par
- Ajout de la documentation sur les événements au site web par
- Modèles pour sveltekit et sveltekit-ts configurés pour le développement sans SSR
- Mise à jour des ressources de build à l’aide de la nouvelle commande `wails3 update build-assets` par
- Exemple permettant de tester l’API HTML de glisser-déposer par
- Prise en charge des associations de fichiers par [leaanthony](https://github.com/leaanthony) dans
- Nouvelle commande `wails3 generate runtime` par
- Nouvelle option `InitialPosition` permettant d’indiquer si la fenêtre doit être centrée ou
- Ajout des méthodes `Path` et `Paths` au paquet `application` par
- Ajout des options Windows `GeneralAutofillEnabled` et `PasswordAutosaveEnabled`
- Ajout de la possibilité de récupérer la fenêtre qui appelle une méthode de service par
- Ajout des options `EnabledFeatures` et `DisabledFeatures` pour Webview2 par
- ⊞ Nouveau système DIP pour une meilleure prise en charge des moniteurs à haute densité de pixels par
- ⊞ Option de nom de classe de fenêtre par [windom](https://github.com/windom/) dans
- Les services ont été étendus afin de fournir des fonctionnalités de plugin. Par
- 🐧 Événements WindowDidMove / WindowDidResize dans
- ⊞ Événement WindowDidResize dans
-  Ajout de l’événement ApplicationShouldHandleReopen afin de pouvoir gérer le dock
-  Ajout de getPrimaryScreen/getScreens à l’implémentation par @tmclane dans
-  Ajout d’une option permettant d’afficher la barre d’outils en mode plein écran sous macOS par
- 🐧 Ajout de la logique onKeyPress pour convertir une pression de touche Linux en accélérateur
- 🐧 Ajout de la tâche `run:linux` par
- Exportation de la méthode `SetIcon` par [@almas-x](https://github.com/almas-x) dans
- Amélioration de `OnShutdown` par [@almas-x](https://github.com/almas-x) dans
- Restauration de la méthode `ToggleMaximise` dans l’interface `Window` par
- Ajout d’informations supplémentaires à `Environment()`. Par @leaanthony dans
- Exposition de la méthode `WebviewWindow.IsFocused` sur l’interface `Window` par
- Prise en charge de plusieurs événements déclencheurs séparés par des espaces dans le système WML par
- Ajout d’exports ESM depuis le script d’environnement d’exécution JS intégré par
- Ajout d’un indicateur au générateur de liaisons pour utiliser le script d’environnement d’exécution JS intégré au lieu de
- Implémentation de `setIcon` sous Linux par [@abichinger](https://github.com/abichinger)
- Ajout de l’indicateur `-port` à la commande dev et prise en charge de la variable d’environnement
- Ajout de tests pour les appels de méthodes liées par
- ⊞ ajout de `SetIgnoreMouseEvents` pour une fenêtre déjà créée par
-  Ajout de la possibilité de définir le niveau d’empilement (ordre) d’une fenêtre par

### Corrections

- Correction des valeurs `Screen.Bounds`, `WorkArea` et `Size` divisées par deux sur les Mac Retina, en convertissant les valeurs en points de NSScreen en pixels physiques dans les champs `Physical*`, et renseignement des propriétés de premier niveau `Screen.X`/`Y` afin que la détection des contacts entre plusieurs moniteurs et le positionnement dans la zone de travail soient corrects dans la [PR](https://github.com/wailsapp/wails/pull/5168) par @wayneforrest
- Correction d’un accès concurrent non synchronisé aux données dans ScreenManager qui provoque un interblocage de WebKit DisplayLink lors d’un changement de configuration d’affichage (par exemple, le branchement à chaud d’un moniteur externe pendant la mise en veille ou la sortie de veille)
- Définition directe de CFBundleIconName sur appicon lorsque Assets.car existe dans la [PR](https://github.com/wailsapp/wails/pull/5154) par @symball
- Correction de `wails3 doctor` qui indiquait des paquets WebKitGTK incorrects sous Fedora, openSUSE, Arch et NixOS — les entrées de repli 4.0 ont été supprimées, car la v3 requiert l’API 4.1 à la compilation (#5071)
- Correction du nom du paquet webkit2gtk dans le diagnostic d’openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, nom correct du paquet openSUSE) (#5071)
- Correction de l’erreur `Unexpected token '<'` lorsque `/wails/custom.js` est absent en mode de développement desktop. Ajout d’un gestionnaire 404 explicite pour `/wails/custom.js` et d’une validation de `Content-Type` insensible à la casse dans `loadOptionalScript`, afin d’empêcher l’injection de pages HTML de repli de SPA en tant que JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Correction de l’état de surbrillance du menu de la zone de notification sous macOS : l’icône affiche désormais l’état sélectionné lorsque le menu est ouvert (#4910)
- Correction de la fenêtre attachée à la zone de notification qui apparaissait derrière les autres fenêtres sous macOS : elle utilise désormais le niveau de fenêtre contextuelle approprié (#4910)
- Correction des exemples d’importation `@wailsio/runtime` erronés dans toute la documentation (#4989)
- Correction de l’impossibilité de réduire une fenêtre sans cadre sous darwin (#4294)
- Correction des blocages de 20-30 minutes pendant `wails3 build` et `wails3 dev`, en excluant `node_modules/` de la vérification de go-task visant à déterminer si les fichiers sont à jour. Auparavant, le motif glob `sources: "**/*"` amenait go-task à énumérer et à calculer la somme de contrôle de chaque fichier dans `node_modules/` (50000-100000 fichiers, voire davantage, avec des dépendances lourdes comme MUI), ce qui était particulièrement lent sous Windows/NTFS (#4939)
- Correction de l’échec de compilation avec GTK4 dû à une collision du typedef C `Screen` avec X11 Xlib.h (#4957)
- Harmonisation des méthodes de badge du Dock sous macOS
- Correction de l’application de `InvisibleTitleBarHeight` à toutes les fenêtres macOS au lieu des seules fenêtres sans cadre ou à barre de titre transparente (#4960)
- Correction des secousses et tremblements de la fenêtre lors de son redimensionnement depuis les coins supérieurs avec `InvisibleTitleBarHeight` activé, en empêchant le début du déplacement à proximité des bords de la fenêtre (#4960)
- Correction de la génération des types mappés dont les clés sont des valeurs d’énumération dans les liaisons JS/TS (#4437) par @fbbdev
- Correction du glisser-déposer de fichiers qui ne fonctionnait pas sous Windows avec une mise à l’échelle d’affichage différente de 100 %
- Correction du glisser-déposer interne HTML5 qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Windows
- Correction des coordonnées de dépôt de fichiers exprimées dans le mauvais espace de pixels sous Windows (confusion entre pixels physiques et pixels CSS)
- Correction du manque de fiabilité du glisser-déposer de fichiers avec les effets de survol sous Linux
- Correction du glisser-déposer interne HTML5 qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Linux
- Correction de l’affichage et du masquage d’une fenêtre sous Linux/GTK4, qui la restauraient parfois à l’état réduit, grâce à l’utilisation de `gtk_window_present()` (#4957)
- Correction de la lecture et de la définition de la position d’une fenêtre sous Linux/GTK4, qui renvoyaient toujours 0,0, grâce à l’ajout d’une prise en charge conditionnelle de X11 via `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Correction de la taille maximale des fenêtres qui n’était pas respectée sous Linux/GTK4, grâce à l’ajout d’un plafonnement de la taille fondé sur les signaux pour remplacer la fonction `gtk_window_set_geometry_hints` supprimée (#4957)
- Correction de la mise à l’échelle selon la résolution sous Linux/GTK4, grâce à l’implémentation d’un calcul correct de PhysicalBounds et de la prise en charge de la mise à l’échelle fractionnaire via `gdk_monitor_get_scale` (GTK 4.14+)
- Correction de la duplication des éléments de menu lors de la création de nouvelles fenêtres sous Linux/GTK4
- Correction de la génération des types mappés dont les clés sont des valeurs d’énumération dans les liaisons JS/TS (#4437) par @fbbdev
- Correction du glisser-déposer de fichiers qui ne fonctionnait pas sous Windows avec une mise à l’échelle d’affichage différente de 100 %
- Correction du glisser-déposer interne HTML5 qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Windows
- Correction des coordonnées de dépôt de fichiers exprimées dans le mauvais espace de pixels sous Windows (confusion entre pixels physiques et pixels CSS)
- Correction du manque de fiabilité du glisser-déposer de fichiers avec les effets de survol sous Linux
- Correction du glisser-déposer interne HTML5 qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Linux
- Correction de la mise à l’échelle selon la résolution sous Linux/GTK4, grâce à l’implémentation d’un calcul correct de PhysicalBounds et de la prise en charge de la mise à l’échelle fractionnaire via `gdk_monitor_get_scale` (GTK 4.14+)
- Correction de la duplication des éléments de menu lors de la création de nouvelles fenêtres sous Linux/GTK4
- Correction de la génération des types mappés dont les clés sont des valeurs d’énumération dans les liaisons JS/TS (#4437) par @fbbdev
- Correction du problème de « fenêtres fantômes » sous macOS, dû à l’utilisation des API AppKit en dehors du thread principal dans App.Window.Current() (#4947) par @wimaha
- Correction de l’élément HTML `<input type="file">` qui ne fonctionnait pas sous macOS, grâce à l’implémentation de WKUIDelegate runOpenPanelWithParameters (#4862)
- Correction du glisser-déposer natif de fichiers qui ne fonctionnait pas avec le module npm `@wailsio/runtime` sous macOS/Linux (#4953) par @leaanthony
- Correction de la génération des liaisons pour les alias de types entre paquets (#4578) par @fbbdev
- Correction du plantage d’OpenFileDialog sous Linux dû à une violation de la sécurité des threads de GTK (#3683) par @ddmoney420
- Correction du plantage SIGSEGV lors de l’appel de `Focus()` sur une fenêtre masquée ou détruite (#4890) par @ddmoney420
- Correction d’une panique potentielle lors de la définition d’une icône ou d’une image bitmap vide sous Linux (#4923) par @ddmoney420
- Correction du plantage d’ErrorDialog lors de son appel depuis une liaison de service sous macOS (#3631) par @leaanthony
- Affichage des menus sous Windows dans `v3\examples\dialogs` par @ndianabasi
- Correction d’une condition de concurrence provoquant une TypeError pendant le rechargement de la page (#4872) par @ddmoney420
- Correction de la sortie erronée des tests du générateur de liaisons, en supprimant l’état global de la méthode `Collector.IsVoidAlias()` (#4941) par @fbbdev
- Correction du sélecteur de fichiers `<input type="file">` qui ne fonctionnait pas sous macOS (#4862) par @leaanthony
- Corrige l’utilisation de systèmes de coordonnées incohérents par `Position()` et `SetPosition()` sous macOS, qui entraînait un décalage de la position de la fenêtre lors de l’enregistrement et de la restauration de l’état (#4816), par @leaanthony
- Corrige l’erreur « Access is denied » de SetProcessDpiAwarenessContext lorsque la gestion de la mise à l’échelle DPI est déjà définie dans le manifeste de l’application (#4803)
- Met à jour la page de documentation sur les raccourcis clavier et corrige le type du paramètre de rappel de `KeyBinding.Add`, par @ndianabasi
- Corrige la documentation sur la génération de liaisons personnalisées : il faut utiliser `-d String` au lieu de `-o String`
- Corrige la non-suppression des éléments enfants d’un menu lors de l’appel à `menu.Update()`
- Corrige les références obsolètes à l’API Manager dans la documentation (31 fichiers utilisent désormais le nouveau modèle, comme `app.Window.New()`, `app.Event.Emit()`, etc.), par @leaanthony
- Corrige un plantage sous Linux lors d’une panique dans les méthodes Go liées à JavaScript, dû au remplacement des gestionnaires de signaux par WebKit (#3965), par @leaanthony
- Corrige l’absence d’effet de SaveFileDialog.SetFilename() sous Linux (#4841), par @samstanier
- Corrige l’affichage de coordonnées de dépôt égales à undefined dans l’exemple de glisser-déposer
- Corrige l’échec de création du paquet d’application macOS lorsque APP_NAME contient des espaces (problème de développement des accolades)
- Corrige une panique d’indice hors limites sous Windows lors de l’appel de méthodes de service (rétablissement de la version antérieure à goccy/go-json)
- Corrige le glisser-déposer de fichiers qui ne fonctionnait pas sous Windows lorsque la mise à l’échelle de l’affichage différait de 100 %
- Corrige le glisser-déposer interne HTML5 qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Windows
- Corrige les coordonnées de dépôt de fichiers exprimées dans le mauvais espace de pixels sous Windows (confusion entre pixels physiques et pixels CSS)
- Corrige le fonctionnement peu fiable du glisser-déposer de fichiers avec les effets de survol sous Linux
- Corrige le glisser-déposer interne HTML5 qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Linux
- Met à jour toutes les commandes des fichiers Taskfile.yml pour tous les systèmes d’exploitation afin de prendre en charge les espaces dans des variables telles que `APP_NAME`, par @ndianabasi
- Corrige une erreur d’argument de commande lors de l’exécution de la tâche 'build:universal:lipo:go' sous Linux, par @wux1an
- Corrige l’erreur Docker « undefined symbol: **<em>ubsan</em>handle_xxxxxxx » lors de l’exécution de 'wails3 build GOOS=darwin GOARCH=arm64' sous Linux, par @wux1an
- Regroupe la documentation sur les protocoles personnalisés et ajoute des sections sur les Universal Links, par @leaanthony
- Corrige le plantage du menu de la zone de notification sous Windows lors de clics répétés sur l’icône, en ajoutant une protection contre les appels simultanés à TrackPopupMenuEx (#4151), par @leaanthony
- Empêche le plantage de l’application lorsque systray.Run() est appelé avant app.Run(), par @leaanthony
- Corrige un plantage sous macOS lors du basculement de la visibilité d’une fenêtre avec Hide()/Show() lorsque ApplicationShouldTerminateAfterLastWindowClosed est activé (#4389), par @leaanthony
- Corrige une fuite de mémoire dans les menus contextuels sous macOS et Windows lors de leur ouverture répétée (#4012), par @leaanthony
- Corrige la non-réutilisation des ressources natives des menus contextuels sous macOS, qui entraînait la création d’un nouveau menu à chaque affichage (#4012), par @leaanthony
- Corrige le clic sur l’icône du Dock sous macOS qui n’affichait pas les fenêtres masquées lorsque l’application était démarrée avec `Hidden: true` (#4583), par @leaanthony
- Corrige l’absence d’ouverture de la boîte de dialogue d’impression sous macOS, due à un type de pointeur de fenêtre incorrect dans l’appel CGO (#4290), par @leaanthony
- Corrige un plantage du menu de fenêtre sous Wayland provoqué par l’accès d’appmenu-gtk-module à une fenêtre non réalisée (#4769), par @leaanthony
- Corrige le plantage de l’application GTK lorsque le nom de l’application contient des caractères non valides (espaces, parenthèses, etc.), par @leaanthony
- Corrige l’erreur « not enough memory » lors de l’initialisation du glisser-déposer sous Windows (#4701), par @overlordtm
- Corrige l’ouverture du mauvais répertoire par l’explorateur de fichiers sous Linux, due à un échappement incorrect de l’URI (#4397), par @leaanthony
- Corrige l’échec de compilation d’une AppImage sur les distributions Linux modernes (Arch, Fedora 39+, Ubuntu 24.04+) en détectant automatiquement les sections ELF `.relr.dyn` et en désactivant la suppression des symboles (#4642), par @leaanthony
- Corrige `wails doctor`, qui signalait à tort les paquets webkit comme installés sur Fedora et les systèmes basés sur DNF (#4457), par @leaanthony
- Corrige le comportement par défaut de `config.yml`, qui exécutait `wails3 dev` avec une version de production, par @mbaklor
- Corrige les échecs de compilation provoqués par les stubs de services iOS qui importaient un paquet inexistant, par @leaanthony
- Corrige la journalisation structurée dans les méthodes debug/info, qui provoquait des erreurs « no formatting directives », par @leaanthony
- Supprime les instructions temporaires d’affichage de débogage incluses par erreur lors de la fusion de la prise en charge des plateformes mobiles, par @leaanthony
- Corrige un plantage de WebKitGTK sous Wayland avec les processeurs graphiques NVIDIA (erreur de protocole 71) en désactivant automatiquement le moteur de rendu DMA-BUF, par @leaanthony
- Corrige la valeur alpha ignorée dans `application.WebviewWindowOptions.BackgroundColour` sous Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Corrige l’icône de la zone de notification sous Windows, qui n’utilisait pas par défaut celle de l’application en l’absence d’icône personnalisée (#4704)
- Suit la propriété des `HICON` afin de ne détruire que les handles créés par l’utilisateur, ce qui évite les plantages lors du redémarrage de l’Explorateur (#4653).
- Libère l’écouteur du thème système de Windows et les icônes conservées dans la zone de notification lors de la destruction, afin d’éliminer les fuites de goroutines et de contextes de périphérique (#4653).
- Tronque les infobulles de la zone de notification à 127 unités UTF-16 pour éviter d’altérer les paires de substitution et les glyphes multioctets (#4653).
- Corrige l’échec de la tâche de création du paquet Windows (#4667)
- Corrige la variable appicon d’AppImage sous Linux dans le fichier Taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige l’erreur de compilation sous Windows provoquée par la modification de signature dans go-webview2 v1.0.22 (#4513, #4645)
- Corrige la variable appicon d’AppImage sous Linux dans le fichier Taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige l’itération sur les protocoles de desktop.tmpl sous Linux en remplaçant `<.Info.Protocol>` par `<.Protocol>`, par @Tolfx dans #4510
- Corrige l’erreur de redéfinition dans la démonstration de Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Corrige les mises à jour du menu de la zone de notification sous Linux [#4604](https://github.com/wailsapp/wails/issues/4604), par [@JackDoan](https://github.com/JackDoan)
- Corrige l’apparition d’une fenêtre blanche sous Windows lors de la création d’une fenêtre masquée, par @leaanthony dans [#4612](https://github.com/wailsapp/wails/pull/4612)
- Corrige le chemin d’importation du paquet de notifications dans la documentation, par @rxliuli dans [#4617](https://github.com/wailsapp/wails/pull/4617)
- Correction du glisser-déposer qui ne fonctionnait pas avec le paquet npm @wailsio/runtime (#4489), par @leaanthony dans #4616
- Windows : correction du scintillement de la fenêtre au démarrage et de l’affichage incorrect des fenêtres masquées dans la [PR](https://github.com/wailsapp/wails/pull/4600), par @leaanthony.
- Correction des problèmes de taille des fenêtres lors de leur agrandissement sous Wayland (https://github.com/wailsapp/wails/issues/4429), par [@samstanier](https://github.com/samstanier)
- Correction des problèmes de taille des fenêtres lors de leur agrandissement sous Wayland (https://github.com/wailsapp/wails/issues/4429), par [@samstanier](https://github.com/samstanier)
- Correction de l’erreur de redéfinition dans la démonstration de Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Correction d’un problème pouvant provoquer un plantage d’AssetServer sous macOS dans [#4576](https://github.com/wailsapp/wails/pull/4576), par @jghiloni
- Correction du problème de compilation avec NextJs. Corrigé dans [#4585](https://github.com/wailsapp/wails/pull/4585) par @rev42
- Correction des pipelines de publication nocturne dans [#4597](https://github.com/wailsapp/wails/pull/4597), par @riadafridishibly
- Correction de l’erreur de redéfinition dans la démonstration de Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Correction d’un problème pouvant provoquer un plantage d’AssetServer sous macOS dans [#4576](https://github.com/wailsapp/wails/pull/4576), par @jghiloni
- Correction du problème de compilation avec NextJs. Corrigé dans [#4585](https://github.com/wailsapp/wails/pull/4585) par @rev42
- Correction des pipelines de publication nocturne dans [#4597](https://github.com/wailsapp/wails/pull/4597), par @riadafridishibly
- Correction de l’erreur de redéfinition dans la démonstration de Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Correction de SetBackgroundColour sous Windows par @PPTGamer dans la [PR](https://github.com/wailsapp/wails/pull/4492)
- Mise à jour de la documentation pour refléter les changements issus de la refactorisation de l’API Manager, par @yulesxoxo dans la [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Correction de la variable appicon du fichier .desktop Linux dans le Taskfile Linux, [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Mise à jour de la documentation pour refléter les changements issus de la refactorisation de l’API Manager, par @yulesxoxo dans la [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Correction du déréférencement d’un pointeur nil sous Windows, signalé dans [#4456](https://github.com/wailsapp/wails/issues/4456), par @leaanthony dans [#4460](https://github.com/wailsapp/wails/pull/4460)
- Ajout de la prise en charge de `allowsBackForwardNavigationGestures` dans WKWebView sous macOS afin d’activer les gestes de navigation par balayage à deux doigts (#1857)
- Correction du problème empêchant onClick de fonctionner pour les éléments de menu initialement désactivés, par @leaanthony dans la [PR #4469](https://github.com/wailsapp/wails/pull/4469). Merci à @IanVS pour l’enquête initiale.
- Correction du serveur Vite qui n’était pas arrêté après l’échec de la compilation (#4403)
- Correction de la panique lors de la fermeture ou de l’annulation d’un `SaveFileDialog` sous Windows. Corrigé dans la [PR](https://github.com/wailsapp/wails/pull/4284) par @hkhere
- Correction du glisser-déposer au niveau HTML sous Windows, par [@mbaklor](https://github.com/mbaklor) dans [#4259](https://github.com/wailsapp/wails/pull/4259)
- Ajout de la prise en charge de `allowsBackForwardNavigationGestures` dans WKWebView sous macOS afin d’activer les gestes de navigation par balayage à deux doigts (#1857)
- Correction du problème empêchant onClick de fonctionner pour les éléments de menu initialement désactivés, par @leaanthony dans la [PR #4469](https://github.com/wailsapp/wails/pull/4469). Merci à @IanVS pour l’enquête initiale.
- Correction du serveur Vite qui n’était pas arrêté après l’échec de la compilation (#4403)
- Correction de l’analyse des notifications sous Windows par @popaprozac dans la [PR](https://github.com/wailsapp/wails/pull/4450)
- Correction de la commande doctor afin qu’elle vérifie les dépendances du SDK Windows, par [@kodumulo](https://github.com/kodumulo) dans [#4390](https://github.com/wailsapp/wails/issues/4390)
- Correction d’un déréférencement de pointeur nil dans processURLRequest sur Mac, par [@etesam913](https://github.com/etesam913) dans [#4366](https://github.com/wailsapp/wails/pull/4366)
- Correction d’un bogue sous Linux qui empêchait l’utilisation de boîtes de dialogue filtrées, par [@bh90210](https://github.com/bh90210) dans [#4287](https://github.com/wailsapp/wails/pull/4287)
- Correction des problèmes du menu Édition sous Windows et Linux, par [@leaanthony](https://github.com/leaanthony) dans [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Mise à jour de la version système minimale dans les fichiers .plist macOS, de 10.13.0 à 10.15.0, par [@AkshayKalose](https://github.com/AkshayKalose) dans [#3981](https://github.com/wailsapp/wails/pull/3981)
- Correction du problème de saut d’identifiant de fenêtre, par [@leaanthony](https://github.com/leaanthony)
- Correction du problème de menu nil lors de l’appel à RegisterContextMenu, par [@leaanthony](https://github.com/leaanthony)
- Correction des cycles de dépendances dans la sortie du générateur de liaisons, par [@fbbdev](https://github.com/fbbdev) dans [#4001](https://github.com/wailsapp/wails/pull/4001)
- Correction des erreurs d’utilisation avant définition dans la sortie du générateur de liaisons, par [@fbbdev](https://github.com/fbbdev) dans [#4001](https://github.com/wailsapp/wails/pull/4001)
- Transmission des indicateurs de compilation au générateur de liaisons, par [@fbbdev](https://github.com/fbbdev) dans [#4023](https://github.com/wailsapp/wails/pull/4023)
- Remplacement des barres obliques inverses par des barres obliques dans les chemins du Taskfile Windows afin de garantir son fonctionnement sur les plateformes autres que Windows, par [@leaanthony](https://github.com/leaanthony)
- Correction des événements Mac et JavaScript sur Mac, par [@leaanthony](https://github.com/leaanthony)
- Correction de l’interblocage des événements sous macOS, par [@leaanthony](https://github.com/leaanthony)
- Correction d’une erreur `Parameter incorrect` lors de l’initialisation d’une fenêtre sous Windows lorsque du HTML était fourni sans JavaScript, par [@leaanthony](https://github.com/leaanthony)
- Correction de la taille du préfixe de réponse utilisé pour détecter le type de contenu dans le serveur de ressources, par [@fbbdev](https://github.com/fbbdev) dans [#4049](https://github.com/wailsapp/wails/pull/4049)
- Correction du traitement des réponses autres que 404 sur le chemin de l’index racine dans le serveur de ressources, par [@fbbdev](https://github.com/fbbdev) dans [#4049](https://github.com/wailsapp/wails/pull/4049)
- Correction d’un comportement indéfini dans le générateur de liaisons lors du test des propriétés de types génériques, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la sortie du générateur de liaisons pour les modèles dont le type sous-jacent n’a pas les mêmes propriétés que le type enveloppe nommé, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la sortie du générateur de liaisons pour les types de clés de map et le prétraitement, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la sortie du générateur de liaisons pour les structs qui implémentent des interfaces de sérialisation, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la détection des cycles de types impliquant des types génériques dans le générateur de liaisons, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction des références non valides à des modèles non exportés dans la sortie du générateur de liaisons par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Déplacement du code injecté à la fin des fichiers de service par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la gestion des erreurs provenant des opérations de fermeture de fichiers dans le générateur de liaisons par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Suppression des avertissements pour les services qui définissent des méthodes de cycle de vie ou HTTP, mais aucune autre méthode liée, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction des modèles autres que React qui n’affichaient pas le pied de page Hello World avec le jeu de couleurs système clair par [@marcus-crane](https://github.com/marcus-crane) dans [#4056](https://github.com/wailsapp/wails/pull/4056)
- Correction des éléments de menu masqués sous macOS par [@leaanthony](https://github.com/leaanthony)
- Correction de la gestion et du formatage des erreurs dans les processeurs de messages par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Correction de l’arrêt des services qui était ignoré lors de la fermeture de l’application par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Exécution garantie des mises à jour des menus sur le thread principal par [@leaanthony](https://github.com/leaanthony)
- Le mécanisme de déplacement et de redimensionnement est désormais plus robuste et correspond plus étroitement au comportement attendu de la plateforme, grâce à [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Correction de [#4097](https://github.com/wailsapp/wails/issues/4097) : Webpack/Angular supprimait le code d’initialisation de l’environnement d’exécution, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Correction des éléments de menu initialement masqués par [@IanVS](https://github.com/IanVS) dans [#4116](https://github.com/wailsapp/wails/pull/4116)
- Correction d’assetFileServer, qui ne servait pas les fichiers `.html` lors d’une requête sans extension lorsque `[request]` n’existe pas, mais que `[request].html` existe
- Correction des chemins de génération des icônes par [@robin-samuel](https://github.com/robin-samuel) dans [#4125](https://github.com/wailsapp/wails/pull/4125)
- Correction de la non-émission des événements `fullscreen`, `unfullscreen`, `unminimise` et `unmaximise` par [@oSethoum](https://github.com/osethoum) dans [#4130](https://github.com/wailsapp/wails/pull/4130)
- Correction de l’erreur NSIS due à un préfixe incorrect dans la version par défaut de la configuration par [@robin-samuel](https://github.com/robin-samuel) dans [#4126](https://github.com/wailsapp/wails/pull/4126)
- Correction de la fonction Dialogs du runtime, qui renvoyait des chemins échappés sous Windows, par [TheGB0077](https://github.com/TheGB0077) dans [#4188](https://github.com/wailsapp/wails/pull/4188)
- Correction du chemin de détection de WebView2 dans HKCU par [@leaanthony](https://github.com/leaanthony).
- Correction d’un problème de saisie sous macOS par [@leaanthony](https://github.com/leaanthony).
- Correction du nom du fichier de tâche de génération des icônes Windows par [@yulesxoxo](https://github.com/yulesxoxo) dans [#4219](https://github.com/wailsapp/wails/pull/4219).
- Correction d’un problème de transparence des fenêtres sans cadre par [@leaanthony](https://github.com/leaanthony), à partir du travail de @kron.
- Correction des appels de gestion du focus lorsqu’une fenêtre est désactivée ou réduite, par [@leaanthony](https://github.com/leaanthony), à partir du travail de @kron.
- Correction des icônes de la zone de notification qui ne réapparaissaient pas après le redémarrage de la barre des tâches par [@leaanthony](https://github.com/leaanthony), à partir du travail de @kron.
- Correction de fallbackResponseWriter, qui n’implémentait pas Flush(), dans [#4245](https://github.com/wailsapp/wails/pull/4245)
- Correction de fallbackResponseWriter, qui n’implémentait pas Flush(), par [@superDingda] dans [#4236](https://github.com/wailsapp/wails/issues/4236)
- Correction des plantages lors de la fermeture d’une fenêtre macOS alors qu’un appel asynchrone à une fonction Go liée était en attente, par [@joshhardy](https://github.com/joshhardy) dans [#4354](https://github.com/wailsapp/wails/pull/4354)
- Correction d’une condition de concurrence au démarrage du mode d’efficacité de Windows par [@leaanthony](https://github.com/leaanthony)
- Correction de la libération des handles d’icônes Windows par [@leaanthony](https://github.com/leaanthony).
- Correction de `OpenFileManager` sous Windows par [@PPTGamer](https://github.com/PPTGamer) dans [#4375](https://github.com/wailsapp/wails/pull/4375).
- Correction des options de largeur minimale/maximale sous Linux par @atterpac dans [#3979](https://github.com/wailsapp/wails/pull/3979)
- Correction des définitions de types des modèles TypeScript grâce à une mise à jour de la version npm par @atterpac dans [#3966](https://github.com/wailsapp/wails/pull/3966)
- Correction de la référence CSS du modèle SvelteKit par @atterpac dans [#3945](https://github.com/wailsapp/wails/pull/3945)
- Exécution garantie sur le thread principal des principaux rappels de run() pour les fenêtres par [@leaanthony](https://github.com/leaanthony)
- Correction des exemples de sélection de répertoire dans les boîtes de dialogue par [@leaanthony](https://github.com/leaanthony)
- Création d’une nouvelle page d’erreur en chinois lorsque index.html est absent par [@leaanthony](https://github.com/leaanthony)
-  Exécution garantie du rappel `windowDidBecomeKey` sur le thread principal par [@leaanthony](https://github.com/leaanthony)
-  Prise en charge du plein écran pour les fenêtres sans cadre par [@leaanthony](https://github.com/leaanthony)
-  Amélioration de la logique de destruction des fenêtres par [@leaanthony](https://github.com/leaanthony)
-  Correction de la logique de positionnement des fenêtres lorsqu’elles sont rattachées aux icônes de la zone de notification par [@leaanthony](https://github.com/leaanthony)
-  Prise en charge du plein écran pour les fenêtres sans cadre par [@leaanthony](https://github.com/leaanthony)
- Correction de la gestion des événements par [@leaanthony](https://github.com/leaanthony)
- Correction de la logique de fermeture des fenêtres par [@leaanthony](https://github.com/leaanthony)
- Le fichier de tâches commun génère désormais par défaut les liaisons TypeScript pour les modèles TypeScript, grâce à [@leaanthony](https://github.com/leaanthony)
- Correction de la fermeture de l’application lors de la réception du message WM_CLOSE lorsqu’aucune fenêtre n’est ouverte ou que seule une icône de la zone de notification est présente, par [@mmalcek](https://github.com/mmalcek) dans [#3990](https://github.com/wailsapp/wails/pull/3990)
- Correction de la compilation avec garble par @5aaee9 dans [#3192](https://github.com/wailsapp/wails/pull/3192)
- Correction des compilations NSIS sous Windows par [@leaanthony](https://github.com/leaanthony)
- Correction de l’interblocage dans la boîte de dialogue Linux de sélection multiple, causé par un élément non fermé
- Correction du nettoyage multiplateforme des fichiers .syso pendant la compilation Windows par
- Correction de la compilation de l’AppImage amd64 par @atterpac dans
- Correction de la mise à jour des ressources de compilation par @ansxuman dans
- Correction de l’implémentation de `OnClick` et `OnRightClick` pour la zone de notification sous Linux par @atterpac
- Correction de `AlwaysOnTop` qui ne fonctionnait pas sur Mac par
-  Correction de `application.NewEditMenu` qui incluait un doublon
- 🐧 Correction de la compilation pour aarch64
- ⊞ Correction des éléments de menu de type bouton radio par
- Correction de l’erreur lors de la création d’une application .app exécutable sur macOS lorsque « name » et « outputfilename »
- Correction d’un bogue lié à l’utilisation de customEventProcessor dans l’exemple de glisser-déposer par
- 🐧 Correction de l’erreur de compilation sous Linux introduite par l’ajout d’IgnoreMouseEvents par
- ⊞ Correction d’un bogue de génération du fichier d’icône syso par
- 🐧 Intégration d’un correctif permettant une exécution native sous Wayland, provenant de
- Ne pas lier les méthodes de service internes dans
- ⊞ Correction d’une panique au démarrage de la zone de notification dans
- Ne pas lier les méthodes de service internes dans
- ⊞ Correction d’une panique au démarrage de la zone de notification dans
- Refactorisation majeure des éléments de menu et de la gestion des événements. Pour l’instant, elle améliore principalement macOS. Par
- Correction des tests après la refactorisation des extensions et des événements dans
- ⊞ Correction de l’avertissement `Failed to unregister class Chrome_WidgetWin_0`. Par
- Problèmes liés aux modules
- Correction de la transmission des événements de redimensionnement par [atterpac](https://github.com/atterpac) dans
- 🐧 Correction d’une erreur de gestion du thème sous NixOS par
- Correction de l’installation d’un projet entre différents volumes sous Windows par
- Correction du CSS du modèle React afin d’afficher le pied de page par
- Correction des processus zombies en mode développement grâce à la mise à jour vers la dernière version de refresh
- Correction de la récupération du fichier WebKit pour AppImage par [Atterpac](https://github.com/atterpac)
- Correction de la vérification des paquets apt par Doctor, par [Atterpac](https://github.com/Atterpac) dans
- Correction du blocage de l’application à sa fermeture (Darwin) par @5aaee9 dans
- Correction des couleurs d’arrière-plan des exemples sous Windows par
- Correction des menus contextuels par défaut par [mmghv](https://github.com/mmghv) dans
- Correction des valeurs hexadécimales des touches fléchées sous Darwin par
- Rétablissement du fonctionnement du glisser-déposer sous Windows. Ajouté par
- Correction d’un bogue sous Linux dans Doctor lorsque l’utilisateur ne dispose pas des pilotes appropriés
- Correction de la mise à l’échelle DPI au démarrage (Windows). Modification par [@almas-x](https://github.com/almas-x) dans
- Correction de la ligne de remplacement dans `go.mod` afin d’utiliser des chemins relatifs. Corrige les chemins Windows contenant
- Correction de la gestion des clics dans la zone de notification sous macOS lorsqu’aucune fenêtre n’y est associée par
- Correction de l’échec de la compilation sous Windows dû à une option inconnue par
- Correction du plantage sous Windows lors d’un clic gauche sur l’icône de la zone de notification en l’absence d’un
- Correction de la baseURL incorrecte lors de la seconde ouverture d’une fenêtre par @5aaee9 dans la PR
- Correction de l’ordre des branches if dans la méthode `WebviewWindow.Restore` par
- Calcul correct de `startURL` sur plusieurs appels à `GetStartURL` lorsque
- Correction du type JS de la structure `Screen` afin qu’il corresponde à son équivalent Go par
- Correction de la méthode `WML.Reload` afin d’assurer le nettoyage approprié des événements enregistrés
- Correction de la fermeture immédiate du menu contextuel personnalisé sous Linux par
- Correction du chemin de sortie et de l’extension des fichiers de modèles produits par la liaison
- Correction des chemins d’importation des fichiers de modèles dans le code JS produit par la liaison
- Correction du glisser-déposer sur certaines distributions Linux par
- Correction de la tâche manquante pour macOS lors de l’utilisation de `wails3 task dev` par
- Correction de l’enregistrement d’événements qui provoquait une affectation dans une map nil par
- Correction de la désérialisation des paramètres des méthodes liées par
- Correction de la gestion des valeurs de retour multiples des méthodes liées par
- Correction de la détection par Doctor d’une installation de npm qui n’a pas été effectuée avec le gestionnaire de paquets système
- Correction de l’absence de MicrosoftEdgeWebview2Setup.exe. Merci à
- Correction d’un plantage aléatoire sous Linux dû à la gestion de l’identifiant de fenêtre par @leaanthony. D’après
- Correction du plantage de systemTray.setIcon sous Linux par
- Correction garantissant l’application du cadre de fenêtre dès le premier appel de la fonction `setFrameless` sur

### Modifié

- **RUPTURE DE COMPATIBILITÉ** : les clés de map dans les liaisons JS/TS générées sont désormais marquées comme facultatives afin de refléter fidèlement la sémantique des maps Go. En TypeScript, l’accès aux valeurs d’une map renvoie désormais `T | undefined` au lieu de `T`, ce qui nécessite des vérifications de valeur nulle ou des assertions (#4943), par `@fbbdev`
- Remplacement de l’utilisation de `Event` par `Events` conformément aux modifications apportées à `@wailsio/runtime`, avec adaptation des appels de fonction dans la documentation de `Features/Events/Event System`, par @AbdelhadiSeddar
- Déplacement de `EnabledFeatures`, `DisabledFeatures` et `AdditionalBrowserArgs` depuis les options propres à chaque fenêtre vers `Options.Windows` au niveau de l’application (#4559), par @leaanthony
- Mise à jour du README de l’exemple `Drag N Drop`, en soulignant que cet exemple illustre `Internal Drag and Drop`, par @ndianabasi
- Passage de divers journaux de débogage du niveau Info au niveau Debug (par @mbaklor)
- **RUPTURE DE COMPATIBILITÉ :** renommage de `EnableDragAndDrop` en `EnableFileDrop` dans les options de fenêtre
- **RUPTURE DE COMPATIBILITÉ :** renommage de `DropZoneDetails` en `DropTargetDetails` dans le contexte d’événement
- **RUPTURE DE COMPATIBILITÉ :** renommage de la méthode `DropZoneDetails()` en `DropTargetDetails()` sur `WindowEventContext`
- **RUPTURE DE COMPATIBILITÉ :** suppression de l’événement `WindowDropZoneFilesDropped` ; utilisez `WindowFilesDropped` à la place
- **RUPTURE DE COMPATIBILITÉ :** remplacer l’attribut HTML `data-wails-dropzone` par `data-file-drop-target`
- **RUPTURE DE COMPATIBILITÉ :** remplacer la classe CSS de survol `wails-dropzone-hover` par `file-drop-target-active`
- **RUPTURE DE COMPATIBILITÉ :** supprimer les options `DragEffect`, `OnEnterEffect` et `OnOverEffect` de Windows (elles faisaient partie de l’interface IDropTarget supprimée)
- Adopter goccy/go-json pour l’ensemble du traitement JSON à l’exécution (liaisons de méthodes, événements, requêtes de la vue web, notifications et kvstore), ce qui améliore les performances de 21-63 % et réduit les allocations mémoire de 40-60 %
- Optimiser l’agencement de la structure BoundMethod et mettre en cache l’indicateur isVariadic afin de réduire le surcoût de chaque appel
- Utiliser un tampon d’arguments alloué sur la pile pour les méthodes comportant `<=8` arguments afin d’éviter les allocations sur le tas
- Optimiser la collecte des résultats des appels de méthodes afin d’éviter l’allocation d’une tranche pour une valeur de retour unique
- Utiliser sync.Map pour le cache des types MIME afin d’améliorer les performances en accès concurrent
- Utiliser un pool de tampons pour lire le corps des requêtes du transport HTTP
- Allouer à la demande le canal CloseNotify dans le détecteur de type de contenu afin de réduire les allocations par requête
- Supprimer la journalisation de débogage CSS du serveur de ressources
- Étendre la table de correspondance des extensions de types MIME afin de couvrir plus de 50 formats web courants (polices, audio, vidéo, etc.)
- Mettre à jour la documentation des options `X/Y` de Window @ruhuang2001
- Mettre à jour la documentation de `Frontend Runtime` en ajoutant davantage d’options pour générer les liaisons du frontend, par @ndianabasi
- Mettre à jour la page de documentation du serveur de ressources de Wails v3, par @ndianabasi
- **RUPTURE DE COMPATIBILITÉ** : supprimer les fonctions de boîte de dialogue au niveau du paquet (`application.InfoDialog()`, `application.QuestionDialog()`, etc.). Utiliser à la place le gestionnaire `app.Dialog` : `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Mettre la documentation des boîtes de dialogue en conformité avec l’API réelle : utiliser `app.Dialog.*`, `AddButton()` avec des fonctions de rappel (et non `SetButtons()`), `SetDefaultButton(*Button)` (et non une chaîne), `AddFilter()` (et non `SetFilters()`), `SetFilename()` (et non `SetDefaultFilename()`) et `app.Dialog.OpenFile().CanChooseDirectories(true)` pour sélectionner un dossier
- **RUPTURE DE COMPATIBILITÉ** : les builds de production sont désormais utilisés par défaut. Pour créer des builds de développement, définir `DEV=true` dans vos Taskfiles. Générer un nouveau projet pour les exemples, par @leaanthony
- Lors de l’émission d’un événement personnalisé avec zéro ou un argument de données, la valeur des données est affectée directement au champ Data sans être encapsulée dans une tranche, par [@fbbdev](https://github.com/fbbdev) dans [#4633](https://github.com/wailsapp/wails/pull/4633)
- Sous Windows, les icônes de zone de notification respectent désormais `SystemTray.Show()`/`Hide()` en basculant `NIS_HIDDEN`, ce qui permet aux applications de disparaître réellement puis de réapparaître (#4653).
- L’enregistrement des icônes de zone de notification réutilise les icônes résolues, définit `NOTIFYICON_VERSION_4` une seule fois et active `NIF_SHOWTIP` afin que les infobulles soient rétablies après le redémarrage de l’Explorateur (#4653).
- macOS : utiliser `visibleFrame` plutôt que `frame` pour centrer la fenêtre en excluant les zones de la barre de menus et du Dock
- macOS : utiliser `visibleFrame` plutôt que `frame` pour centrer la fenêtre en excluant les zones de la barre de menus et du Dock
- Lors de l’exécution de `wails3 update build-assets` avec le paramètre `-config`, les valeurs définies au moyen des paramètres `-product*` sont
- `window.NativeWindowHandle()` -> `window.NativeWindow()`, par @leaanthony dans [#4471](https://github.com/wailsapp/wails/pull/4471)
- Remanier la gestion interne des fenêtres, par @leaanthony dans [#4471](https://github.com/wailsapp/wails/pull/4471)
- Suppression de `application.WindowIDKey` et `application.WindowNameKey` (remplacés par `application.WindowKey`), par [@leaanthony](https://github.com/leaanthony)
- ContextMenuData renvoie désormais une chaîne plutôt que any, par [@leaanthony](https://github.com/leaanthony)
- Dans les liaisons JS/TS, les champs de classe dont le type est un tableau de longueur fixe sont désormais initialisés avec la longueur attendue au lieu d’être vides, par [@fbbdev](https://github.com/fbbdev) dans [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData renvoie désormais une chaîne plutôt que any, par [@leaanthony](https://github.com/leaanthony)
- `application.NewService` n’accepte plus les options comme paramètre facultatif (utiliser `application.NewServiceWithOptions` à la place), par [@leaanthony](https://github.com/leaanthony) dans [#4024](https://github.com/wailsapp/wails/pull/4024)
- Suppression de la dépendance `nanoid`, par [@leaanthony](https://github.com/leaanthony)
- Mise à jour de l’exemple Window pour les styles de fenêtre mica, acrylique et à onglets, par [@leaanthony](https://github.com/leaanthony)
- Dans les liaisons JS/TS, les fichiers de modèles `internal.js/ts` ont été supprimés ; tous les modèles se trouvent désormais dans `models.js/ts`, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dans les liaisons JS/TS, les types nommés ne sont jamais rendus comme des alias d’autres types nommés ; l’ancien comportement est désormais limité aux alias, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dans les liaisons JS/TS en mode classe, les champs de structure dont le type est un paramètre de type sont marqués comme facultatifs et ne sont jamais initialisés automatiquement, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Supprimer ESLint des modèles, par [@IanVS](https://github.com/IanVS) dans [#4059](https://github.com/wailsapp/wails/pull/4059)
- Mettre à jour l’année de copyright vers 2025, par [@IanVS](https://github.com/IanVS) dans [#4037](https://github.com/wailsapp/wails/pull/4037)
- Ajouter la documentation d’event.Sender, par [@IanVS](https://github.com/IanVS) dans [#4075](https://github.com/wailsapp/wails/pull/4075)
- Prise en charge de Go 1.24, par [@leaanthony](https://github.com/leaanthony)
- Les hooks `ServiceStartup` sont désormais invoqués lors de l’appel de `App.Run`, et non dans `application.New`, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Les erreurs de `ServiceStartup` sont désormais renvoyées par `App.Run` au lieu de mettre fin au processus, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Les appels de liaisons et de boîtes de dialogue depuis JS sont désormais rejetés avec des objets d’erreur plutôt qu’avec des chaînes, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Amélioration du positionnement du menu de l’icône de zone de notification sous Windows, par [@leaanthony](https://github.com/leaanthony)
- Le runtime JS a été porté vers TypeScript, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Le runtime s’initialise dès son importation, sans qu’il soit nécessaire d’attendre le chargement de la fenêtre, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Le runtime n’exporte plus de méthode init. Un import avec effets de bord peut servir à l’initialiser, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Les méthodes liées renvoient désormais une `CancellablePromise` qui est rejetée avec une `CancelError` en cas d’annulation. Le résultat réel de l’appel est ignoré, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Les types de services intégrés sont désormais systématiquement appelés `Service`, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Les fonctions de création de services intégrés avec options sont désormais systématiquement appelées `NewWithConfig`, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- La méthode `Select` du service `sqlite` s’appelle désormais `Query` par souci de cohérence avec les API Go, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Modèles : déplacement du runtime dans "dependencies" et organisation des fichiers package.json, par [@IanVS](https://github.com/IanVS) dans [#4133](https://github.com/wailsapp/wails/pull/4133)
- Création et signature ad hoc des bundles d’application en développement afin d’activer certaines API macOS, par [@popaprozac](https://github.com/popaprozac) dans [#4171](https://github.com/wailsapp/wails/pull/4171)
- Déplacement des ressources de compilation vers des répertoires propres à chaque plateforme, par [@leaanthony](https://github.com/leaanthony)
- Déplacement et renommage des Taskfiles dans des répertoires propres à chaque plateforme, par [@leaanthony](https://github.com/leaanthony)
- Expérience nettement améliorée lorsque `index.html` est absent, par [@leaanthony](https://github.com/leaanthony)
- [Windows] Amélioration des performances de réduction et de restauration, par [@leaanthony](https://github.com/leaanthony). Basé sur la [PR](https://github.com/wailsapp/wails/pull/3955) d’origine de [562589540](https://github.com/562589540)
- Suppression de l’option `ShouldClose` (enregistrez plutôt un hook pour events.Common.WindowClosing), par [@leaanthony](https://github.com/leaanthony)
- [Windows] Réduction du scintillement à l’ouverture d’une fenêtre, par [@leaanthony](https://github.com/leaanthony)
- Suppression de `Window.Destroy`, car cette fonction était destinée à un usage interne, par [@leaanthony](https://github.com/leaanthony)
- Renommage des événements `WindowClose` en `WindowClosing`, par [@leaanthony](https://github.com/leaanthony)
- Les builds frontend utilisent désormais l’environnement vite "development" ou "production" selon le type de build, par [@leaanthony](https://github.com/leaanthony)
- Mise à jour vers go-webview2 v1.19, par [@leaanthony](https://github.com/leaanthony)
- Utilisation garantie du fork de taskfile, par @leaanthony
- Mise à jour du fork de Taskfile pour corriger les problèmes de version lors de l’installation avec
- Utilisation du fork de Taskfile pour corriger les problèmes de version lors de l’installation avec
- `service.OnStartup` arrête désormais l’application en cas d’erreur et exécute
- Refactorisation des messages de clic de la zone de notification afin de mieux correspondre aux interactions utilisateur, par
- Intégration de `all:frontend/dist` aux ressources pour prendre en charge les frameworks qui génèrent
- Refactorisation de Taskfile par [leaanthony](https://github.com/leaanthony) dans
- Mise à niveau vers `go-webview2` v1.0.16 par
- Correction du type `Screen` afin d’inclure `ID` et non `Id`, par
- Mise à jour de la version de Wails dans `go.mod.tmpl` pour prendre en charge `application.ServiceOptions`, par
- Correction de la détermination du nom du service par [windom](https://github.com/windom/) dans
- mkdocs serve utilise désormais docker, par [leaanthony](https://github.com/leaanthony)
- Regroupement de la configuration de développement dans `config.yml`, par
- La boîte de dialogue de la zone de notification utilise désormais par défaut l’icône de l’application si elle est disponible (Windows), par
- Meilleur rapport sur le GPU et la mémoire sous macOS, par
- Suppression de `WebviewGpuIsDisabled` et de `EnableFraudulentWebsiteWarnings`
- Modification de l’API Events : `On`/`Emit` -> événements utilisateur, `OnApplicationEvent` ->
- Correction de l’API Events sous Linux par [TheGB0077](https://github.com/TheGB0077) dans
- [CI] amélioration des actions et possibilité de les exécuter également dans les forks et
- Renommage de `AbsolutePosition()` en `Position()`, par
- Mise à jour de la dépendance WebKit Linux vers webkit2gtk-4.1 au lieu de webkitgtk2-4.0 afin de
- Le script du runtime JS fourni est désormais un module ESM : les balises script qui l’importent
- Le paquet `@wailsio/runtime` ne publie pas son API sur le `window.wails`
- Le module `@wailsio/runtime/src/window` de l’API Window expose désormais l’élément englobant
- L’API Window JS a été mise à jour pour correspondre à l’actuel `WebviewWindow` Go
- Le générateur de liaisons utilise désormais les appels par ID par défaut. L’option CLI `-id`
- Nouvelle organisation du code de liaison : les fichiers de sortie étaient auparavant organisés dans des dossiers
- Le champ de structure `application.Options.Bind` a été renommé en
- Nouvelle syntaxe pour lier des services : les instances de service doivent désormais être encapsulées dans un
- Désactivation de l’indicateur d’activité dans un environnement sans terminal ou CI, par

### Supprimé

- **RUPTURE** : suppression de `EnabledFeatures`, `DisabledFeatures` et `AdditionalLaunchArgs` des options `WindowsWindow` propres à chaque fenêtre. Utilisez plutôt les options `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` et `Options.Windows.AdditionalBrowserArgs` au niveau de l’application. Ces indicateurs s’appliquent globalement à l’environnement WebView2 partagé (#4559), par @leaanthony
- Suppression de l’implémentation native de `IDropTarget` sous Windows au profit de l’approche fondée sur JavaScript (conforme au comportement de la v2)
- Suppression de la dépendance github.com/wailsapp/mimetype au profit d’une table d’extensions étendue et de la fonction http.DetectContentType de la bibliothèque standard, ce qui réduit la taille du binaire d’environ 1.2MB
- Suppression de la dépendance gopkg.in/ini.v1 grâce à l’implémentation d’un analyseur minimal de fichiers .desktop pour l’explorateur de fichiers Linux, soit un gain d’environ 45KB
- Suppression de samber/lo du code d’exécution grâce à l’utilisation du paquet slices de la bibliothèque standard de Go 1.21+ et d’un nombre minimal de fonctions auxiliaires internes, ce qui économise environ 310 Ko
- Suppression des instructions printf de débogage du gestionnaire de schémas d’URL de Darwin (#4834)
- **RUPTURE DE COMPATIBILITÉ** : suppression de l’événement `linux:WindowLoadChanged` ; utilisez désormais `linux:WindowLoadFinished` pour détecter la fin du chargement de la WebView (#3896), par @leaanthony

### Ruptures de compatibilité

- **Refactorisation de l’API des gestionnaires** : réorganisation de l’API de l’application, auparavant plate, en gestionnaires structurés afin d’améliorer l’organisation du code et la facilité de découverte, par [@leaanthony](https://github.com/leaanthony) dans [#4359](https://github.com/wailsapp/wails/pull/4359)
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Renommage des méthodes de Service : `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`, par [@leaanthony](https://github.com/leaanthony)
- Déplacement des méthodes `Path` et `Paths` vers le paquet `application`, par [@leaanthony](https://github.com/leaanthony)
- Le menu de l’application est désormais réservé à macOS, par [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## Ajouté

## Corrigé

## v3.0.0-alpha.77 - 2026-04-18

## Corrigé

## v3.0.0-alpha.76 - 2026-04-17

## Corrigé

## v3.0.0-alpha.75 - 2026-04-16

## Corrigé

## v3.0.0-alpha.74 - 2026-03-01

## Ajouté

## Corrigé

## v3.0.0-alpha.73 - 2026-02-27

## Corrigé

## v3.0.0-alpha.72 - 2026-02-16

## Corrigé

## v3.0.0-alpha.71 - 2026-02-10

## Ajouté

## Corrigé

## v3.0.0-alpha.70 - 2026-02-09

## Ajouté

## Corrigé

## v3.0.0-alpha.69 - 2026-02-08

## Ajouté

## Corrigé

## v3.0.0-alpha.68 - 2026-02-07

## Ajouté

## Modifié

## Corrigé

## v3.0.0-alpha.67 - 2026-02-04

## Ajouté

## Modifié

## Corrigé

## v3.0.0-alpha.66 - 2026-02-03

## Ajouté

## Modifié

## Corrigé

## Supprimé

## v3.0.0-alpha.65 - 2026-02-01

## Ajouté

## v3.0.0-alpha.64 - 2026-01-26

## Ajouté

## v3.0.0-alpha.63 - 2026-01-25

## Corrigé

## v3.0.0-alpha.62 - 2026-01-22

## Corrigé

## v3.0.0-alpha.61 - 2026-01-20

## Corrigé

## v3.0.0-alpha.60 - 2026-01-14

## Corrigé

## v3.0.0-alpha.59 - 2026-01-11

## Modifié

## v3.0.0-alpha.58 - 2026-01-09

## Corrigé

## v3.0.0-alpha.57 - 2026-01-05

## Modifié

## Corrigé

## v3.0.0-alpha.56 - 2026-01-04

## Ajouté

## Modifié

## Corrigé

## Supprimé

## v3.0.0-alpha.55 - 2026-01-02

## Modifié

## Corrigé

## Supprimé

## v3.0.0-alpha.54 - 2025-12-29

## Ajouté

## Corrigé

## Supprimé

## v3.0.0-alpha.53 - 2025-12-27

## Ajouté

## Corrigé

## v3.0.0-alpha.52 - 2025-12-26

## Corrigé

## v3.0.0-alpha.51 - 2025-12-23

## Corrigé

## v3.0.0-alpha.50 - 2025-12-21

## Modifié

## v3.0.0-alpha.49 - 2025-12-18

## Modifié

## v3.0.0-alpha.48 - 2025-12-16

## Ajouté

## Modifié

## Corrigé

## v3.0.0-alpha.47 - 2025-12-15

## Ajouté

## Corrigé

## v3.0.0-alpha.46 - 2025-12-14

## Ajouté

## Supprimé

## v3.0.0-alpha.45 - 2025-12-13

## Ajouté

## Corrigé

## v3.0.0-alpha.44 - 2025-12-12

## Ajouté

## Modifié

## Corrigé

## v3.0.0-alpha.43 - 2025-12-11

## Ajouté

## v3.0.0-alpha.42 - 2025-12-10

## Ajouté

## v3.0.0-alpha.41 - 2025-11-23

## Corrigé

## v3.0.0-alpha.40 - 2025-11-13

## Corrigé

## v3.0.0-alpha.39 - 2025-11-12

## Ajouté

## Modifié

## v3.0.0-alpha.38 - 2025-11-04

## Ajouté

## Modifié

## Corrigé

## v3.0.0-alpha.37 - 2025-11-02

## Corrigé

## v3.0.0-alpha.36 - 2025-10-15

## Corrigé

## v3.0.0-alpha.35 - 2025-10-14

## Corrigé

## v3.0.0-alpha.34 - 2025-10-06

## Ajouté

## Corrigé

## v3.0.0-alpha.33 - 2025-10-04

## Corrigé

## v3.0.0-alpha.32 - 2025-10-02

## Corrigé

## v3.0.0-alpha.31 - 2025-09-27

## Corrigé

## v3.0.0-alpha.30 - 2025-09-26

## Corrigé

## v3.0.0-alpha.29 - 2025-09-25

## Ajouté

## Modifié

## Corrigé

## v3.0.0-alpha.29 - 2025-09-25

## Ajouté

## Modifié

## Corrigé

## v3.0.0-alpha.27 - 2025-09-07

## Corrigé

## v3.0.0-alpha.26 - 2025-08-24

## Ajouté

## v3.0.0-alpha.25 - 2025-08-16

## Modifié

ne sont plus ignorés et remplacent la valeur de configuration.

## v3.0.0-alpha.24 - 2025-08-13

## Ajouté

## v3.0.0-alpha.23 - 2025-08-11

## Corrigé

## v3.0.0-alpha.22 - 2025-08-10

## Ajouté

## Modifié

+ Correction des dépendances trop générales des paquets Linux et des dépendances RPM obsolètes.

## v3.0.0-alpha.21 - 2025-08-07

## Corrigé

## v3.0.0-alpha.20 - 2025-08-06

## Corrigé

## v3.0.0-alpha.19 - 2025-08-05

## Ajouté

## Corrigé

## v3.0.0-alpha.18 - 2025-08-03

## Ajouté

## Corrigé

## v3.0.0-alpha.17 - 2025-07-31

## Corrigé

## v3.0.0-alpha.16 - 2025-07-25

## Ajouté

## v3.0.0-alpha.15 - 2025-07-25

## Ajouté

## v3.0.0-alpha.14 - 2025-07-25

## Ajouté

## v3.0.0-alpha.12 - 2025-07-15

### Ajouté

### Corrigé

## v3.0.0-alpha.11 - 2025-07-12

## Ajouté

## v3.0.0-alpha.10 - 2025-07-06

### Modifications incompatibles

### Ajouté

### Corrigé

### Modifié

## v3.0.0-alpha.9 - 2025-01-13

### Ajouté

### Corrigé

### Modifié

## v3.0.0-alpha.8.3 - 2024-12-07

### Modifié

## v3.0.0-alpha.8.2 - 2024-12-07

### Modifié

`go install` par @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### Modifié

`go install` par @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### Ajouté

@atterpac dans [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) dans   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) dans   [#3867](https://github.com/wailsapp/wails/pull/3867)   par [atterpac](https://github.com/atterpac) dans   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) dans   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positionnée aux coordonnées X/Y indiquées par   [leaanthony](https://github.com/leaanthony) dans   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) et   [leaanthony](https://github.com/leaanthony) dans   [#3823](https://github.com/wailsapp/wails/pull/3823)   par [leaanthony](https://github.com/leaanthony) dans   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) dans   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Modifié

`service.OnShutdown`pour tous les services précédemment démarrés par @atterpac dans   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac dans [#3907](https://github.com/wailsapp/wails/pull/3907)   les sous-dossiers par @atterpac dans   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) dans   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) dans   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (remplacé par les options `EnabledFeatures` et `DisabledFeatures`) par   [leaanthony](https://github.com/leaanthony)

### Corrigé

variable de canal par @michael-freling dans   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) dans   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   dans [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) dans   [#3841](https://github.com/wailsapp/wails/pull/3841)   rôle `PasteAndMatchStyle` dans le menu d’édition sous Darwin par   [johnmccabe](https://github.com/johnmccabe) dans   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) dans   [#3854](https://github.com/wailsapp/wails/pull/3854) par   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   sont différents, par @nickisworking dans   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Ajouté

[mmghv](https://github.com/mmghv) dans   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) et   [leaanthony](https://github.com/leaanthony) dans   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Modifié

Événements d’application `OnWindowEvent` -> Événements de fenêtre, par   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   les branches préfixées par `v3/` ou `v3-` par   [stendler](https://github.com/stendler) dans   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Corrigé

[etesam913](https://github.com/etesam913) dans   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) dans   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) dans   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) dans   [#3614](https://github.com/wailsapp/wails/pull/3614) par   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) par   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) par   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) par   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) par   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) par   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Corrigé

## v3.0.0-alpha.5 - 2024-07-30

### Ajouté

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   clic sur l’icône par @5aaee9 dans [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) dans   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   dans [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) dans   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) dans   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   basé sur la [PR](https://github.com/wailsapp/wails/pull/2044) de @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   le paquet npm par [@fbbdev](https://github.com/fbbdev) dans   [#3334](https://github.com/wailsapp/wails/pull/3334)   dans [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` par [@abichinger](https://github.com/abichinger) dans   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) dans   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) dans   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) dans   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Corrigé

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) dans   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) dans   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) dans   [#3477](https://github.com/wailsapp/wails/pull/3477)   par [Atterpac](https://github.com/atterpac) dans   [#3320](https://github.com/wailsapp/wails/pull/3320).   dans [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) dans   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) dans   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) dans la   [PR](https://github.com/wailsapp/wails/pull/3039)   installé. Ajouté par [@pylotlight](https://github.com/pylotlight) dans la   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   espaces — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) dans la PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) dans la PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   fenêtre attachée [tw1nk](https://github.com/tw1nk) dans la PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) dans   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` est présent.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   écouteurs par [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) dans   [#3330](https://github.com/wailsapp/wails/pull/3330)   générateur par [@fbbdev](https://github.com/fbbdev) dans   [#3334](https://github.com/wailsapp/wails/pull/3334)   générateur par [@fbbdev](https://github.com/fbbdev) dans   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) dans   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) dans   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) dans   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) dans   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) dans   [#3431](https://github.com/wailsapp/wails/pull/3431)   par [@pekim](https://github.com/pekim) dans   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) par   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) dans   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows par [@bruxaodev](https://github.com/bruxaodev/) dans   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Modifications

[mmghv](https://github.com/mmghv) dans   [#3611](https://github.com/wailsapp/wails/pull/3611)   prise en charge d’Ubuntu 24.04 LTS par [atterpac](https://github.com/atterpac) dans   [#3461](https://github.com/wailsapp/wails/pull/3461)   doit posséder l’attribut `type="module"`. Par   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   objet et ne démarre pas le système WML. Cette modification améliore   l’encapsulation. Si vous le souhaitez, vous pouvez démarrer manuellement le système WML en appelant   la nouvelle méthode `WML.Enable`. Le script d’exécution JS intégré effectue toujours automatiquement les deux   opérations. Par [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   objet window comme export par défaut. Il n’est désormais plus possible d’importer   des méthodes individuelles au moyen de la syntaxe ESM d’import nommé ou d’espace de noms.   API. Certaines méthodes ont changé de nom ou de prototype, notamment : `Screen`   devient `GetScreen` ; `GetZoomLevel`/`SetZoomLevel` deviennent `GetZoom`/`SetZoom` ;   `GetZoom`, `Width` et `Height` renvoient désormais directement les valeurs au lieu de les encapsuler   dans des objets. Par [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   a été supprimé. Utilisez l’option de CLI `-names` pour revenir aux appels par nom.   Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   nommés d’après leur paquet conteneur ; les chemins d’import Go complets sont désormais utilisés,   y compris le chemin du module. Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   appel à `application.NewService`. Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) dans   [#3574](https://github.com/wailsapp/wails/pull/3574)

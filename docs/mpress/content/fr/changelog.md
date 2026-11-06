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

Le format repose sur [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), et ce projet respecte la [gestion sémantique de version](https://semver.org/spec/v2.0.0.html).

- `Added` pour les nouvelles fonctionnalités.
- `Changed` pour les modifications apportées aux fonctionnalités existantes.
- `Deprecated` pour les fonctionnalités qui seront bientôt supprimées.
- `Removed` pour les fonctionnalités désormais supprimées.
- `Fixed` pour toutes les corrections de bogues.
- `Security` en cas de vulnérabilités.

_/

/_   * MIS À JOUR AUTOMATIQUEMENT *   Les notes de publication sont générées à partir des PR fusionnées. Ne modifiez pas ce fichier dans les PR de contribution ordinaires et ne recréez pas le fichier `v3/UNRELEASED_CHANGELOG.md` supprimé. Consultez `AGENTS.md` pour la maintenance du journal des modifications. _/

## [Non publié]

## v3.2.0-nightly.20261102 - 2026-11-02

## Modifié
- Utilisation de DisableOverscroll au lieu de DisableBounce dans la [PR](https://github.com/wailsapp/wails/pull/6077) par @lethusang

## v3.2.0-nightly.20261101 - 2026-11-01

## Ajouté
- Ajout de backends expérimentaux sans CGO pour macOS et Linux dans la [PR](https://github.com/wailsapp/wails/pull/6521) par @taliesin-ai
- Builds YAML natifs, empaquetage portable et signature dans la [PR](https://github.com/wailsapp/wails/pull/6531) par @taliesin-ai

## Corrigé
- Correction de divers problèmes de build nocturne de v3 et nettoyage des dépendances dans la [PR](https://github.com/wailsapp/wails/pull/6519) par @taliesin-ai
- Réparation de la validation des versions de v3 et de la classification des fixtures dans la [PR](https://github.com/wailsapp/wails/pull/6533) par @taliesin-ai
- Résolution des échecs de validation multiplateforme des versions dans la [PR](https://github.com/wailsapp/wails/pull/6538) par @taliesin-ai

## v3.1.0 - 2026-11-01

## Ajouté
- Ajout de backends expérimentaux sans CGO pour macOS et Linux dans la [PR](https://github.com/wailsapp/wails/pull/6521) par @taliesin-ai
- Builds YAML natifs, empaquetage portable et signature dans la [PR](https://github.com/wailsapp/wails/pull/6531) par @taliesin-ai

## Modifié
- Non-copie des métadonnées de bureau ignorées dans la [PR](https://github.com/wailsapp/wails/pull/6391) par @taliesin-ai
- Formatage des longueurs de contenu sans fmt dans la [PR](https://github.com/wailsapp/wails/pull/6393) par @taliesin-ai
- Sélection des racines d’index sans collecter les candidats dans la [PR](https://github.com/wailsapp/wails/pull/6404) par @taliesin-ai
- Réutilisation des tampons de copie des fichiers non positionnables dans la [PR](https://github.com/wailsapp/wails/pull/6408) par @taliesin-ai
- Non-allocation d’un lecteur pour détecter les préfixes dans la [PR](https://github.com/wailsapp/wails/pull/6410) par @taliesin-ai
- Lecture de détection MIME ignorée pour les extensions connues dans la [PR](https://github.com/wailsapp/wails/pull/6414) par @taliesin-ai
- Prévention de la croissance de la file d’attente et des oscillations de réduction dans la [PR](https://github.com/wailsapp/wails/pull/6415) par @taliesin-ai
- Non-allocation des contextes de requête ignorés dans la [PR](https://github.com/wailsapp/wails/pull/6419) par @taliesin-ai
- Préparation des journaux de débogage désactivés ignorée dans la [PR](https://github.com/wailsapp/wails/pull/6420) par @taliesin-ai
- Allocation différée des tampons de détection MIME dans la [PR](https://github.com/wailsapp/wails/pull/6428) par @taliesin-ai
- Réduction des allocations liées au formatage des hyperliens dans la [PR](https://github.com/wailsapp/wails/pull/6443) par @taliesin-ai
- Suppression des allocations temporaires liées au formatage des puces dans la [PR](https://github.com/wailsapp/wails/pull/6449) par @taliesin-ai
- Réduction de la conversion des UUID en chaîne à une seule allocation dans la [PR](https://github.com/wailsapp/wails/pull/6453) par @taliesin-ai
- Réduction des allocations liées au formatage des handles d’instruction dans la [PR](https://github.com/wailsapp/wails/pull/6454) par @taliesin-ai
- Élimination de l’allocation temporaire du condensat des UUID fondés sur un nom dans la [PR](https://github.com/wailsapp/wails/pull/6455) par @taliesin-ai
- Non-copie des lignes incomplètes de sortie de commande dans la [PR](https://github.com/wailsapp/wails/pull/6460) par @taliesin-ai
- Copie en flux continu des exécutables pendant la préparation des paquets dans la [PR](https://github.com/wailsapp/wails/pull/6472) par @taliesin-ai
- Non-réanalyse des lignes de sortie incomplètes dans la [PR](https://github.com/wailsapp/wails/pull/6473) par @taliesin-ai
- Réutilisation des motifs compilés de routage des commandes dans la [PR](https://github.com/wailsapp/wails/pull/6475) par @taliesin-ai
- Suppression des recherches de métadonnées pendant la découverte du runtime dans la [PR](https://github.com/wailsapp/wails/pull/6482) par @taliesin-ai
- Réduction des allocations de la page de repli localisée dans la [PR](https://github.com/wailsapp/wails/pull/6485) par @taliesin-ai
- Non-recompilation des motifs pour chaque processeur dans la [PR](https://github.com/wailsapp/wails/pull/6490) par @taliesin-ai

## Corrigé
- Implémentation de SetTooltip pour la zone de notification dans la [PR](https://github.com/wailsapp/wails/pull/6017) par @fan711
- Lecture du mode sombre depuis le portail des paramètres dans la [PR](https://github.com/wailsapp/wails/pull/6072) par @TechHutTV
- Possibilité pour Window.SetZoom de descendre sous 1.0 dans la [PR](https://github.com/wailsapp/wails/pull/6166) par @APshenkin
- Distribution à chaque écouteur Go dans l’ordre d’émission dans la [PR](https://github.com/wailsapp/wails/pull/6275) par @DevLumuz
- Lecture des hooks d’événement sous hookLock dans Emit dans la [PR](https://github.com/wailsapp/wails/pull/6276) par @DevLumuz
- Annonce systématique de la propriété Menu de la zone de notification afin qu’une icône créée après Run ne soit pas ignorée dans la [PR](https://github.com/wailsapp/wails/pull/6284) par @AlbinoGeek
- Mise en attente du travail du thread principal envoyé avant Run jusqu’au démarrage de Run dans la [PR](https://github.com/wailsapp/wails/pull/6287) par @AlbinoGeek
- Utilisation de l’URL de ressource correcte dans le fournisseur de mise à jour GitHub dans la [PR](https://github.com/wailsapp/wails/pull/6304) par @gavidroselj
- Conservation du corps de la requête lorsque WebView2 omet Content-Length dans la [PR](https://github.com/wailsapp/wails/pull/6326) par @strraaa
- Signalement des écrans macOS et GTK dans l’espace DIP (limites avec plusieurs densités de pixels) dans la [PR](https://github.com/wailsapp/wails/pull/6346) par @julianstorer
- Une notification ignorée n’est plus signalée comme un clic dans la [PR](https://github.com/wailsapp/wails/pull/6367) par @bercedev
- Arrêt des processus de développement à la fermeture du terminal dans la [PR](https://github.com/wailsapp/wails/pull/6371) par @xmasdev
- Association des exécutions du journal des modifications aux commits issus d’une fusion par squash lors de la vérification des preuves dans la [PR](https://github.com/wailsapp/wails/pull/6373) par @taliesin-ai
- Maintien pendant une seconde de l’anti-rebond appliqué au redémarrage de la barre des tâches dans la [PR](https://github.com/wailsapp/wails/pull/6386) par @taliesin-ai
- Conservation des contraintes de taille lors d’un passage en plein écran qui commence avec ces contraintes dans la [PR](https://github.com/wailsapp/wails/pull/6387) par @taliesin-ai
- Non-redimensionnement d’une fenêtre en plein écran ou agrandie pour respecter une nouvelle contrainte dans la [PR](https://github.com/wailsapp/wails/pull/6388) par @taliesin-ai
- ne pas afficher de boîte de dialogue avant Run dans l’exemple liquid-glass dans la [PR](https://github.com/wailsapp/wails/pull/6389) par @taliesin-ai
- une boîte de dialogue demandée avant Run ne fait rien au lieu d’attendre indéfiniment dans la [PR](https://github.com/wailsapp/wails/pull/6390) par @taliesin-ai
- utiliser la solution de repli lorsque l’entrée de bureau ne contient aucune commande dans la [PR](https://github.com/wailsapp/wails/pull/6392) par @taliesin-ai
- détecter le type du contenu actuel sans conserver les noms de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6394) par @taliesin-ai
- reconnaître les extensions de ressources dont la casse est mixte dans la [PR](https://github.com/wailsapp/wails/pull/6395) par @taliesin-ai
- libérer les tampons consommés par la détection du type de contenu dans la [PR](https://github.com/wailsapp/wails/pull/6396) par @taliesin-ai
- libérer les références ringqueue supprimées dans la [PR](https://github.com/wailsapp/wails/pull/6397) par @taliesin-ai
- respecter les suppressions d’en-têtes avant la sélection de la solution de repli dans la [PR](https://github.com/wailsapp/wails/pull/6398) par @taliesin-ai
- valider les en-têtes de la réponse de repli avant le vidage dans la [PR](https://github.com/wailsapp/wails/pull/6399) par @taliesin-ai
- ne détecter que les fichiers d’index dont le nom correspond exactement dans la [PR](https://github.com/wailsapp/wails/pull/6400) par @taliesin-ai
- éviter les corps et les lectures intégrales pour les requêtes HEAD non repositionnables dans la [PR](https://github.com/wailsapp/wails/pull/6401) par @taliesin-ai
- remplir le préfixe de détection MIME sur plusieurs lectures courtes dans la [PR](https://github.com/wailsapp/wails/pull/6402) par @taliesin-ai
- valider les réponses non repositionnables vides ayant réussi dans la [PR](https://github.com/wailsapp/wails/pull/6403) par @taliesin-ai
- éviter de dupliquer les cibles d’URI développées dans la [PR](https://github.com/wailsapp/wails/pull/6405) par @taliesin-ai
- préserver les chemins lors du développement des champs dans la [PR](https://github.com/wailsapp/wails/pull/6406) par @taliesin-ai
- respecter les répertoires de recherche de fichiers desktop de XDG dans la [PR](https://github.com/wailsapp/wails/pull/6407) par @taliesin-ai
- valider les en-têtes du détecteur lors d’un vidage à vide dans la [PR](https://github.com/wailsapp/wails/pull/6409) par @taliesin-ai
- ignorer les répertoires qui masquent des entrées de bureau dans la [PR](https://github.com/wailsapp/wails/pull/6411) par @taliesin-ai
- préserver les cibles trouvées dans les chemins d’exécutables dans la [PR](https://github.com/wailsapp/wails/pull/6412) par @taliesin-ai
- libérer les captures de la fonction de rappel après le déclenchement du minuteur dans la [PR](https://github.com/wailsapp/wails/pull/6413) par @taliesin-ai
- éviter les correspondances erronées de bibliothèques ldconfig dans la [PR](https://github.com/wailsapp/wails/pull/6416) par @taliesin-ai
- restreindre les correspondances de bibliothèques de repli dans le système de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6417) par @taliesin-ai
- éliminer la situation de compétition lors de l’invalidation du cache dans la [PR](https://github.com/wailsapp/wails/pull/6418) par @taliesin-ai
- respecter la configuration facultative du journal dans la [PR](https://github.com/wailsapp/wails/pull/6421) par @taliesin-ai
- éviter la mutation partagée de l’URL de démarrage et les ports obsolètes dans la [PR](https://github.com/wailsapp/wails/pull/6422) par @taliesin-ai
- arrêter les nouvelles tentatives du proxy de développement lorsque la requête est annulée dans la [PR](https://github.com/wailsapp/wails/pull/6423) par @taliesin-ai
- renvoyer 404 lorsque des ressources intégrées sont absentes dans la [PR](https://github.com/wailsapp/wails/pull/6424) par @taliesin-ai
- gérer les requêtes HEAD visant le runtime intégré dans la [PR](https://github.com/wailsapp/wails/pull/6425) par @taliesin-ai
- préserver les chemins échappés lors du routage des services dans la [PR](https://github.com/wailsapp/wails/pull/6426) par @taliesin-ai
- ordonner correctement les identifiants numériques de préversion surdimensionnés dans la [PR](https://github.com/wailsapp/wails/pull/6427) par @taliesin-ai
- rejeter les erreurs HTTP avant d’enregistrer les téléchargements d’outils dans la [PR](https://github.com/wailsapp/wails/pull/6429) par @taliesin-ai
- respecter l’échéance de lancement pendant la détection des applications de bureau dans la [PR](https://github.com/wailsapp/wails/pull/6430) par @taliesin-ai
- préserver les corps de requête entre les nouvelles tentatives du proxy de développement dans la [PR](https://github.com/wailsapp/wails/pull/6432) par @taliesin-ai
- fermer les corps de réponse HTTP dans tous les chemins de retour dans la [PR](https://github.com/wailsapp/wails/pull/6434) par @taliesin-ai
- préserver les erreurs d’annulation dans les recherches parallèles de bibliothèques dans la [PR](https://github.com/wailsapp/wails/pull/6435) par @taliesin-ai
- rejeter les réponses d’étiquettes GitHub ayant échoué avant leur décodage dans la [PR](https://github.com/wailsapp/wails/pull/6436) par @taliesin-ai
- préserver les chemins littéraux de modèles contenant des crochets dans la [PR](https://github.com/wailsapp/wails/pull/6437) par @taliesin-ai
- installer wails3 une seule fois et préserver les erreurs du programme d’installation dans la [PR](https://github.com/wailsapp/wails/pull/6438) par @taliesin-ai
- ignorer les entrées et les en-têtes masqués dans les blocs de commentaires dans la [PR](https://github.com/wailsapp/wails/pull/6439) par @taliesin-ai
- signaler les erreurs de recherche de version stable dans la [PR](https://github.com/wailsapp/wails/pull/6440) par @taliesin-ai
- afficher les messages en texte brut sans interpréter leur format dans la [PR](https://github.com/wailsapp/wails/pull/6441) par @taliesin-ai
- éviter les paniques d’écriture après le chargement de la valeur JSON null dans la [PR](https://github.com/wailsapp/wails/pull/6442) par @taliesin-ai
- renvoyer des instantanés lors de la lecture de l’ensemble du magasin dans la [PR](https://github.com/wailsapp/wails/pull/6444) par @taliesin-ai
- propager les erreurs après l’itération partielle des lignes dans la [PR](https://github.com/wailsapp/wails/pull/6446) par @taliesin-ai
- conserver tels quels les chemins de destination de l’extraction dans la [PR](https://github.com/wailsapp/wails/pull/6447) par @taliesin-ai
- libérer les handles des instructions préparées lors de la fermeture et de la réouverture dans la [PR](https://github.com/wailsapp/wails/pull/6448) par @taliesin-ai
- gérer la préparation en situation de compétition avec une fermeture ou une réouverture dans la [PR](https://github.com/wailsapp/wails/pull/6450) par @taliesin-ai
- respecter les requêtes HEAD pour la page d’index de repli dans la [PR](https://github.com/wailsapp/wails/pull/6451) par @taliesin-ai
- exiger des handles d’instruction hexadécimaux complets dans la [PR](https://github.com/wailsapp/wails/pull/6452) par @taliesin-ai
- préserver la base de données en mémoire par défaut en cas d’accès concurrents dans la [PR](https://github.com/wailsapp/wails/pull/6456) par @taliesin-ai
- séparer stderr de la sortie HEAD analysée dans la [PR](https://github.com/wailsapp/wails/pull/6457) par @taliesin-ai
- générer les ressources du runtime à partir des bundles intégrés dans la [PR](https://github.com/wailsapp/wails/pull/6458) par @taliesin-ai
- détecter les marqueurs TODO, FIXME et TBD dans la [PR](https://github.com/wailsapp/wails/pull/6459) par @taliesin-ai
- valider les réponses de l’API des notes de version dans la [PR](https://github.com/wailsapp/wails/pull/6461) par @taliesin-ai
- respecter les répertoires des tâches lors des vérifications d’état dans la [PR](https://github.com/wailsapp/wails/pull/6462) par @taliesin-ai
- empêcher la perte de données lors de la copie d’un fichier sur lui-même dans la [PR](https://github.com/wailsapp/wails/pull/6463) par @taliesin-ai
- inclure le répertoire de la tâche dans l’identité du cache dans la [PR](https://github.com/wailsapp/wails/pull/6464) par @taliesin-ai
- sonder les ports une seule fois par vérification dans la [PR](https://github.com/wailsapp/wails/pull/6465) par @taliesin-ai
- éviter une panique lorsque des entrées de cache sont décodées avec la valeur null dans la [PR](https://github.com/wailsapp/wails/pull/6466) par @taliesin-ai
- incrémenter les préversions numériques sans limites d’entier dans la [PR](https://github.com/wailsapp/wails/pull/6467) par @taliesin-ai
- empêcher les vérifications historiques de masquer les absences dans le cache dans la [PR](https://github.com/wailsapp/wails/pull/6468) par @taliesin-ai
- exécuter les tâches par défaut dans le répertoire de base de l’exécuteur dans la [PR](https://github.com/wailsapp/wails/pull/6469) par @taliesin-ai
- préserver les signes de pourcentage littéraux dans les commandes de l’explorateur dans la [PR](https://github.com/wailsapp/wails/pull/6470) par @taliesin-ai
- aligner les répertoires des préconditions et des commandes de tâche dans la [PR](https://github.com/wailsapp/wails/pull/6471) par @taliesin-ai
- remplacer les liens symboliques par des fichiers ordinaires générés dans la [PR](https://github.com/wailsapp/wails/pull/6474) par @taliesin-ai
- conserver les DMG existants en cas d’échec de la compilation dans la [PR](https://github.com/wailsapp/wails/pull/6476) par @taliesin-ai
- préserver les expressions régulières entre guillemets et les barres obliques inverses des chemins dans la [PR](https://github.com/wailsapp/wails/pull/6477) par @taliesin-ai
- valider les dimensions des icônes générées avant l’allocation des images dans la [PR](https://github.com/wailsapp/wails/pull/6478) par @taliesin-ai
- faire passer les chaînes de commandes et les redirections par le shell dans la [PR](https://github.com/wailsapp/wails/pull/6479) par @taliesin-ai
- conserver les paquets Linux précédents en cas d’échec de la compilation dans la [PR](https://github.com/wailsapp/wails/pull/6480) par @taliesin-ai
- utiliser une exécution de secours par le shell propre à chaque plateforme dans la [PR](https://github.com/wailsapp/wails/pull/6481) par @taliesin-ai
- servir les ressources dont le chemin commence par plusieurs barres obliques dans la [PR](https://github.com/wailsapp/wails/pull/6483) par @taliesin-ai
- reconnaître les versions stables comportant des métadonnées de compilation dans la [PR](https://github.com/wailsapp/wails/pull/6484) par @taliesin-ai
- signaler avec exactitude les échecs de génération des ressources dans la [PR](https://github.com/wailsapp/wails/pull/6486) par @taliesin-ai
- empêcher les conditions de concurrence lors du remplacement du rapporteur actif dans la [PR](https://github.com/wailsapp/wails/pull/6487) par @taliesin-ai
- sérialiser les enregistrements simultanés des fichiers de cache dans la [PR](https://github.com/wailsapp/wails/pull/6488) par @taliesin-ai
- prendre en charge plusieurs candidats d’égalité dans les modèles dans la [PR](https://github.com/wailsapp/wails/pull/6489) par @taliesin-ai
- normaliser les tables de catégories enregistrées qui valent null dans la [PR](https://github.com/wailsapp/wails/pull/6491) par @taliesin-ai
- afficher les diagnostics des commandes shell variables ayant échoué dans la [PR](https://github.com/wailsapp/wails/pull/6492) par @taliesin-ai
- exclure les chaînes de commandes shell du cache Go implicite dans la [PR](https://github.com/wailsapp/wails/pull/6493) par @taliesin-ai
- sérialiser la persistance des catégories sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/6494) par @taliesin-ai
- mettre en cache les sorties Go comportant des espaces entre guillemets ou échappés dans la [PR](https://github.com/wailsapp/wails/pull/6495) par @taliesin-ai
- conserver les icônes Mac en cas d’échec de l’encodage dans la [PR](https://github.com/wailsapp/wails/pull/6496) par @taliesin-ai
- cesser d’ignorer les indicateurs de nombre de tests Go non valides dans la [PR](https://github.com/wailsapp/wails/pull/6497) par @taliesin-ai
- conserver les permissions de la source lors de nouvelles copies de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6498) par @taliesin-ai
- corriger divers problèmes de compilation nightly de v3 et nettoyer les dépendances dans la [PR](https://github.com/wailsapp/wails/pull/6519) par @taliesin-ai
- réparer la validation des versions de v3 et la classification des fixtures dans la [PR](https://github.com/wailsapp/wails/pull/6533) par @taliesin-ai
- résoudre les échecs de validation des versions sur plusieurs plateformes dans la [PR](https://github.com/wailsapp/wails/pull/6538) par @taliesin-ai

## v3.1.0-nightly.20261028 - 2026-10-28

## Modifié
- Accélérer la génération du journal des modifications en réutilisant les expressions régulières compilées dans la [PR](https://github.com/wailsapp/wails/pull/6490) par @taliesin-ai
- Réduire les allocations lors du service des pages de secours localisées dans la [PR](https://github.com/wailsapp/wails/pull/6485) par @taliesin-ai
- Accélérer la détection du runtime AppImage en évitant les consultations inutiles de métadonnées dans la [PR](https://github.com/wailsapp/wails/pull/6482) par @taliesin-ai
- Réutiliser les modèles compilés de routage des commandes pour améliorer les performances de wake dans la [PR](https://github.com/wailsapp/wails/pull/6475) par @taliesin-ai
- Améliorer les performances en évitant d’analyser de nouveau les lignes de sortie inachevées dans la [PR](https://github.com/wailsapp/wails/pull/6473) par @taliesin-ai
- Copier les exécutables en flux continu pendant la préparation des paquets MSIX afin de réduire l’utilisation de la mémoire et d’accélérer les compilations dans la [PR](https://github.com/wailsapp/wails/pull/6472) par @taliesin-ai
- Éviter de copier les lignes inachevées de sortie des commandes afin d’améliorer les performances dans la [PR](https://github.com/wailsapp/wails/pull/6460) par @taliesin-ai
- Réduire les allocations lors de la génération d’UUID basés sur un nom dans la [PR](https://github.com/wailsapp/wails/pull/6455) par @taliesin-ai
- Réduire les allocations lors du formatage des handles d’instructions SQLite dans la [PR](https://github.com/wailsapp/wails/pull/6454) par @taliesin-ai
- Réduire la conversion des UUID en chaînes à une seule allocation dans la [PR](https://github.com/wailsapp/wails/pull/6453) par @taliesin-ai
- Réduire les allocations lors du formatage des entrées à puces du journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/6449) par @taliesin-ai
- Réduire les allocations lors du formatage des hyperliens du terminal dans la [PR](https://github.com/wailsapp/wails/pull/6443) par @taliesin-ai
- Réduire les allocations lors du service des ressources en ne détectant le type MIME des tampons qu’à la demande dans la [PR](https://github.com/wailsapp/wails/pull/6428) par @taliesin-ai
- Ignorer la préparation des journaux de débogage dans le serveur de ressources lorsque leur journalisation est désactivée dans la [PR](https://github.com/wailsapp/wails/pull/6420) par @taliesin-ai
- Réduire les allocations lors du traitement des requêtes de la webview en évitant les contextes de requête ensuite ignorés dans la [PR](https://github.com/wailsapp/wails/pull/6419) par @taliesin-ai
- Accélérer le service des ressources en omettant la détection du type MIME pour les extensions de fichier connues dans la [PR](https://github.com/wailsapp/wails/pull/6414) par @taliesin-ai
- Réduire les allocations lors du service des ressources en évitant d’utiliser un lecteur pour détecter le type du contenu dans la [PR](https://github.com/wailsapp/wails/pull/6410) par @taliesin-ai
- Réduire les allocations lors du service de fichiers de ressources ne permettant pas le déplacement dans la [PR](https://github.com/wailsapp/wails/pull/6408) par @taliesin-ai
- Accélérer la sélection de la racine de l’index du serveur de ressources en évitant de collecter les candidats dans la [PR](https://github.com/wailsapp/wails/pull/6404) par @taliesin-ai
- Améliorer les performances du serveur de ressources en formatant les longueurs de contenu sans fmt dans la [PR](https://github.com/wailsapp/wails/pull/6393) par @taliesin-ai
- Réduire les copies en mémoire lors de la lecture des métadonnées de l’explorateur de fichiers du bureau dans la [PR](https://github.com/wailsapp/wails/pull/6391) par @taliesin-ai

## Corrigé
- Conserver les permissions du fichier source lors de sa copie vers de nouveaux fichiers de destination dans la [PR](https://github.com/wailsapp/wails/pull/6498) par @taliesin-ai
- Empêcher le redimensionnement des fenêtres en plein écran ou agrandies lorsque les contraintes de taille minimale ou maximale changent dans la [PR](https://github.com/wailsapp/wails/pull/6388) par @taliesin-ai
- Publier les mises à jour planifiées de l’historique des étoiles du dépôt wailsapp/wails dans la [PR](https://github.com/wailsapp/wails/pull/6500) par @taliesin-ai
- Cesser d’ignorer les indicateurs de nombre de tests Go non valides dans wake dans la [PR](https://github.com/wailsapp/wails/pull/6497) par @taliesin-ai
- Conserver les icônes Mac en cas d’échec de l’encodage dans la [PR](https://github.com/wailsapp/wails/pull/6496) par @taliesin-ai
- Corriger les échecs de mise en cache de compilation pour les sorties Go comportant des espaces entre guillemets ou échappés dans la [PR](https://github.com/wailsapp/wails/pull/6495) par @taliesin-ai
- Sérialiser la persistance des catégories de notifications sous Linux afin d’éviter les conditions de concurrence lors d’écritures simultanées dans la [PR](https://github.com/wailsapp/wails/pull/6494) par @taliesin-ai
- Exclure les chaînes de commandes shell du cache de compilation Go implicite dans la [PR](https://github.com/wailsapp/wails/pull/6493) par @taliesin-ai
- Corriger l’échec du chargement des tables de catégories de notifications lorsque les données enregistrées valent null dans la [PR](https://github.com/wailsapp/wails/pull/6491) par @taliesin-ai
- Afficher les diagnostics des commandes shell variables ayant échoué dans la [PR](https://github.com/wailsapp/wails/pull/6492) par @taliesin-ai
- Corriger la correspondance des modèles de wake lorsque plusieurs candidats d’égalité existent dans la [PR](https://github.com/wailsapp/wails/pull/6489) par @taliesin-ai
- Sérialiser les enregistrements simultanés du fichier de cache afin d’éviter toute écriture corrompue dans le cache dans la [PR](https://github.com/wailsapp/wails/pull/6488) par @taliesin-ai
- Éviter les conditions de concurrence lors du remplacement du rapporteur actif dans la [PR](https://github.com/wailsapp/wails/pull/6487) par @taliesin-ai
- Signaler correctement les échecs de génération des ressources dans la CLI dans la [PR](https://github.com/wailsapp/wails/pull/6486) par @taliesin-ai
- Reconnaître dans la CLI les versions stables qui incluent des métadonnées de build dans la [PR](https://github.com/wailsapp/wails/pull/6484) par @taliesin-ai
- Servir les ressources dont le chemin de requête contient plusieurs barres obliques initiales dans la [PR](https://github.com/wailsapp/wails/pull/6483) par @taliesin-ai
- Corriger l’exécution de secours via le shell pour wake sur chaque plateforme dans la [PR](https://github.com/wailsapp/wails/pull/6481) par @taliesin-ai
- Conserver les paquets Linux précédents en cas d’échec du build dans la [PR](https://github.com/wailsapp/wails/pull/6480) par @taliesin-ai
- Faire passer les chaînes de commandes et les redirections par le shell dans la [PR](https://github.com/wailsapp/wails/pull/6479) par @taliesin-ai
- Conserver les expressions régulières entre guillemets et les barres obliques inverses des chemins lors de l’analyse des commandes wake dans la [PR](https://github.com/wailsapp/wails/pull/6477) par @taliesin-ai
- Conserver les fichiers DMG existants en cas d’échec du packaging des builds dans la [PR](https://github.com/wailsapp/wails/pull/6476) par @taliesin-ai
- Remplacer les liens symboliques par des fichiers ordinaires générés dans la sortie des bindings dans la [PR](https://github.com/wailsapp/wails/pull/6474) par @taliesin-ai
- Conserver les signes de pourcentage littéraux dans les commandes de l’explorateur Linux dans la [PR](https://github.com/wailsapp/wails/pull/6470) par @taliesin-ai
- Aligner les répertoires des préconditions et des commandes de tâches afin qu’ils soient résolus de manière cohérente dans la [PR](https://github.com/wailsapp/wails/pull/6471) par @taliesin-ai
- Exécuter les tâches par défaut dans le répertoire de base de l’exécuteur dans la [PR](https://github.com/wailsapp/wails/pull/6469) par @taliesin-ai
- Empêcher les vérifications héritées de masquer les absences dans le cache dans la [PR](https://github.com/wailsapp/wails/pull/6468) par @taliesin-ai
- Corriger la numérotation des préversions dans la CLI afin qu’elle s’incrémente sans être limitée par les dépassements de capacité d’entier dans la [PR](https://github.com/wailsapp/wails/pull/6467) par @taliesin-ai
- Éviter une panique lorsque des entrées du cache sont décodées en null dans la [PR](https://github.com/wailsapp/wails/pull/6466) par @taliesin-ai
- Ne sonder les ports qu’une seule fois par vérification dans la CLI dans la [PR](https://github.com/wailsapp/wails/pull/6465) par @taliesin-ai
- Inclure le répertoire de la tâche dans l’identité du cache afin d’éviter des résultats wake obsolètes dans la [PR](https://github.com/wailsapp/wails/pull/6464) par @taliesin-ai
- Éviter toute perte de données lorsque le packaging copie un fichier sur lui-même dans la [PR](https://github.com/wailsapp/wails/pull/6463) par @taliesin-ai
- Faire en sorte que les vérifications d’état respectent les répertoires des tâches dans la [PR](https://github.com/wailsapp/wails/pull/6462) par @taliesin-ai
- Détecter les marqueurs TODO, FIXME et TBD dans les entrées du journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/6459) par @taliesin-ai
- Générer les ressources du runtime à partir des bundles intégrés dans la CLI dans la [PR](https://github.com/wailsapp/wails/pull/6458) par @taliesin-ai
- Séparer stderr de la sortie HEAD analysée dans les commandes git dans la [PR](https://github.com/wailsapp/wails/pull/6457) par @taliesin-ai
- Préserver la base de données SQLite en mémoire par défaut lors des accès simultanés dans la [PR](https://github.com/wailsapp/wails/pull/6456) par @taliesin-ai
- Renvoyer des réponses correctes aux requêtes HEAD visant la page d’index de secours dans la [PR](https://github.com/wailsapp/wails/pull/6451) par @taliesin-ai
- Gérer les conditions de concurrence entre la préparation d’une instruction SQLite et la fermeture ou la réouverture de la base de données dans la [PR](https://github.com/wailsapp/wails/pull/6450) par @taliesin-ai
- Libérer les handles des instructions préparées lors de la fermeture ou de la réouverture de la base de données SQLite dans la [PR](https://github.com/wailsapp/wails/pull/6448) par @taliesin-ai
- Conserver tels quels les chemins de destination de l’extraction lors de l’utilisation de gosod dans la [PR](https://github.com/wailsapp/wails/pull/6447) par @taliesin-ai
- Propager les erreurs des requêtes SQLite après une itération partielle sur les lignes dans la [PR](https://github.com/wailsapp/wails/pull/6446) par @taliesin-ai
- Renvoyer des instantanés pour les lectures de l’ensemble du stockage dans kvstore dans la [PR](https://github.com/wailsapp/wails/pull/6444) par @taliesin-ai
- Éviter les paniques lors des écritures dans kvstore après le chargement d’une valeur JSON null dans la [PR](https://github.com/wailsapp/wails/pull/6442) par @taliesin-ai
- Afficher les messages simples du terminal sans interpréter les verbes de formatage dans la [PR](https://github.com/wailsapp/wails/pull/6441) par @taliesin-ai
- Signaler les erreurs lors de la recherche de la dernière version stable dans l’outil de mise à jour dans la [PR](https://github.com/wailsapp/wails/pull/6440) par @taliesin-ai
- Ignorer les entrées masquées et les en-têtes dans les blocs de commentaires du journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/6439) par @taliesin-ai
- N’installer wails3 qu’une seule fois et signaler les erreurs du programme d’installation dans l’outil de mise à jour dans la [PR](https://github.com/wailsapp/wails/pull/6438) par @taliesin-ai
- Conserver tels quels les chemins de modèles entre crochets dans gosod dans la [PR](https://github.com/wailsapp/wails/pull/6437) par @taliesin-ai
- Rejeter les réponses en échec pour les tags GitHub avant de les décoder dans l’outil de mise à jour dans la [PR](https://github.com/wailsapp/wails/pull/6436) par @taliesin-ai
- Conserver les erreurs d’annulation pendant les recherches parallèles de bibliothèques dans la [PR](https://github.com/wailsapp/wails/pull/6435) par @taliesin-ai
- Fermer le corps des réponses HTTP sur chaque chemin de retour afin d’éviter les fuites de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6434) par @taliesin-ai
- Conserver le corps des requêtes lorsque le proxy de développement les réessaie dans la [PR](https://github.com/wailsapp/wails/pull/6432) par @taliesin-ai
- Respecter le délai limite de lancement lors de la détection de l’explorateur de fichiers du bureau dans la [PR](https://github.com/wailsapp/wails/pull/6430) par @taliesin-ai
- Rejeter les erreurs HTTP avant d’enregistrer les téléchargements des outils AppImage dans la [PR](https://github.com/wailsapp/wails/pull/6429) par @taliesin-ai
- Ordonner correctement les identifiants numériques de préversion surdimensionnés lors des comparaisons semver dans la [PR](https://github.com/wailsapp/wails/pull/6427) par @taliesin-ai
- Conserver les chemins échappés lors du routage des services dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6426) par @taliesin-ai
- Gérer les requêtes HEAD visant le runtime inclus dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6425) par @taliesin-ai
- Renvoyer 404 lorsque des ressources incluses sont introuvables, au lieu d’une réponse incorrecte, dans la [PR](https://github.com/wailsapp/wails/pull/6424) par @taliesin-ai
- Arrêter les nouvelles tentatives du proxy de développement lorsqu’une requête est annulée dans la [PR](https://github.com/wailsapp/wails/pull/6423) par @taliesin-ai
- Éviter la mutation partagée de l’URL de démarrage et les ports obsolètes dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6422) par @taliesin-ai
- Respecter la configuration facultative du logger dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6421) par @taliesin-ai
- Supprimer la condition de concurrence lors de l’invalidation du cache des chemins de bibliothèques dans la [PR](https://github.com/wailsapp/wails/pull/6418) par @taliesin-ai
- Éviter les fausses correspondances de bibliothèques lors de la recherche dans la sortie de ldconfig dans la [PR](https://github.com/wailsapp/wails/pull/6416) par @taliesin-ai
- Empêcher la file d’attente du serveur de ressources d’alterner entre croissance et réduction sous charge dans la [PR](https://github.com/wailsapp/wails/pull/6415) par @taliesin-ai
- Libérer les captures du callback après le déclenchement du temporisateur d’anti-rebond afin d’éviter de conserver des références obsolètes dans la [PR](https://github.com/wailsapp/wails/pull/6413) par @taliesin-ai
- Conserver les cibles de l’explorateur de fichiers trouvées dans les chemins d’exécutables dans la [PR](https://github.com/wailsapp/wails/pull/6412) par @taliesin-ai
- Ignorer les répertoires qui masquent des entrées de bureau dans l’explorateur de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6411) par @taliesin-ai
- Corriger l’absence d’envoi des en-têtes du serveur de ressources lorsque le corps de la réponse est vide dans la [PR](https://github.com/wailsapp/wails/pull/6409) par @taliesin-ai
- Respecter les répertoires XDG de recherche des fichiers desktop lors de la localisation des icônes d’application dans la [PR](https://github.com/wailsapp/wails/pull/6407) par @taliesin-ai
- Préserver les chemins de fichiers lors du développement des champs dans l’explorateur de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6406) par @taliesin-ai
- Éviter de dupliquer les URI cibles développées dans l’explorateur de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6405) par @taliesin-ai
- Valider les réponses de ressources vides réussies et non repositionnables au lieu de les laisser en attente dans la [PR](https://github.com/wailsapp/wails/pull/6403) par @taliesin-ai
- Remplir correctement le préfixe de détection MIME malgré des lectures courtes dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6402) par @taliesin-ai
- Éviter de lire intégralement le corps des réponses aux requêtes HEAD non repositionnables dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6401) par @taliesin-ai
- Ne servir que les fichiers d’index correspondant exactement lors de la résolution des requêtes de répertoire dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6400) par @taliesin-ai
- Corriger les notifications Linux ignorées qui étaient signalées comme des clics sur les serveurs ne prenant pas en charge les actions dans la [PR](https://github.com/wailsapp/wails/pull/6367) par @bercedev
- Arrêter les processus de développement à la fermeture du terminal dans la [PR](https://github.com/wailsapp/wails/pull/6371) par @xmasdev
- Maintenir l’anti-rebond du redémarrage de la barre des tâches Windows pendant une seconde avant d’effacer l’indicateur dans la [PR](https://github.com/wailsapp/wails/pull/6386) par @taliesin-ai
- Corriger la perte des en-têtes de réponse de repli dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6399) par @taliesin-ai
- Respecter les suppressions d’en-têtes avant la sélection de la réponse de repli dans le serveur de ressources dans la [PR](https://github.com/wailsapp/wails/pull/6398) par @taliesin-ai
- Libérer les références supprimées de la file circulaire dans le serveur de ressources afin d’éviter les fuites de mémoire dans la [PR](https://github.com/wailsapp/wails/pull/6397) par @taliesin-ai
- Libérer les tampons consommés par le détecteur de contenu dans le serveur de ressources afin de réduire l’utilisation de la mémoire dans la [PR](https://github.com/wailsapp/wails/pull/6396) par @taliesin-ai
- Reconnaître les extensions de fichiers de ressources dont la casse est mixte lors du service des ressources dans la [PR](https://github.com/wailsapp/wails/pull/6395) par @taliesin-ai
- Détecter correctement le type de contenu dans le serveur de ressources sans conserver les noms de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6394) par @taliesin-ai
- Utiliser une valeur par défaut lorsqu’une entrée de bureau ne contient aucune commande dans la [PR](https://github.com/wailsapp/wails/pull/6392) par @taliesin-ai
- Empêcher l’exemple liquid-glass d’afficher une boîte de dialogue avant le démarrage de l’application sur les systèmes autres que macOS dans la [PR](https://github.com/wailsapp/wails/pull/6389) par @taliesin-ai
- Conserver les contraintes de taille de la fenêtre lorsque le mode plein écran est activé depuis une taille contrainte sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/6387) par @taliesin-ai
- Corriger les boîtes de dialogue appelées avant Run qui se bloquaient indéfiniment au lieu de rendre la main dans la [PR](https://github.com/wailsapp/wails/pull/6390) par @taliesin-ai

## Sécurité
- Valider la taille des icônes générées avant l’allocation de l’image afin d’éviter une utilisation excessive de la mémoire dans la [PR](https://github.com/wailsapp/wails/pull/6478) par @taliesin-ai
- Valider les réponses de l’API des notes de version avant de les utiliser dans la [PR](https://github.com/wailsapp/wails/pull/6461) par @taliesin-ai
- Exiger des handles d’instruction hexadécimaux complets dans le plugin SQLite dans la [PR](https://github.com/wailsapp/wails/pull/6452) par @taliesin-ai
- Restreindre les correspondances de repli de la bibliothèque du système de fichiers afin d’empêcher le chargement de bibliothèques non prévues dans la [PR](https://github.com/wailsapp/wails/pull/6417) par @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Modifications
- Relier la page 404 aux pages correspondantes de la documentation Wails v2 sur v2.wails.io dans la [PR](https://github.com/wailsapp/wails/pull/6366) par @taliesin-ai
- Déplacer le site de documentation v2 vers v2.wails.io dans la [PR](https://github.com/wailsapp/wails/pull/6365) par @taliesin-ai
- Servir la documentation v3 depuis wails.io et la documentation v2 depuis v2.wails.io dans la [PR](https://github.com/wailsapp/wails/pull/6364) par @taliesin-ai
- Mettre à jour les notes de version de v3.0.0 et corriger l’affirmation de l’avis de sécurité dans la [PR](https://github.com/wailsapp/wails/pull/6363) par @taliesin-ai

## Corrections
- Exécuter séquentiellement les rappels des écouteurs d’événements Go dans l’ordre d’émission dans la [PR](https://github.com/wailsapp/wails/pull/6275) par @DevLumuz
- Toujours publier la propriété Menu de la zone de notification afin que les icônes créées après Run ne soient pas ignorées sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/6284) par @AlbinoGeek
- Conserver les tâches du thread principal soumises avant Run jusqu’au démarrage de l’application dans la [PR](https://github.com/wailsapp/wails/pull/6287) par @AlbinoGeek
- Utiliser l’URL correcte de l’API des ressources lors du téléchargement des ressources d’une version par le programme de mise à jour GitHub dans la [PR](https://github.com/wailsapp/wails/pull/6304) par @gavidroselj
- Corriger les limites et les positions des écrans dans les configurations macOS et Linux à valeurs DPI mixtes dans la [PR](https://github.com/wailsapp/wails/pull/6346) par @julianstorer
- Préserver le corps des requêtes sous Windows lorsque WebView2 omet l’en-tête Content-Length dans la [PR](https://github.com/wailsapp/wails/pull/6326) par @strraaa
- Corriger la condition de concurrence lors de l’émission des hooks d’événement en prenant un instantané des hooks sous verrou dans la [PR](https://github.com/wailsapp/wails/pull/6276) par @DevLumuz
- Corriger `SystemTray.SetTooltip()`, qui ne faisait rien sous Linux : la description de l’infobulle StatusNotifierItem est désormais définie (#6017)
- Exécuter la vérification des preuves du journal des modifications dans le workflow de publication hebdomadaire dans la [PR](https://github.com/wailsapp/wails/pull/6374) par @taliesin-ai
- Corriger la vérification des preuves du journal des modifications de version afin d’associer les exécutions aux commits de fusion par squash dans la [PR](https://github.com/wailsapp/wails/pull/6373) par @taliesin-ai
- Corriger la détection du mode sombre sous Linux, qui renvoyait le mode clair sur un bureau sombre jusqu’à la modification du jeu de couleurs, ainsi que `linux:SystemThemeChanged`, qui ne se déclenchait jamais avec le backend GTK4. Les deux backends partagent désormais une même implémentation reposant sur l’espace de noms du portail `org.freedesktop.appearance` (#6072)
- `Window.SetZoom` ne relève plus à 1.0 les facteurs inférieurs à 1.0 pour les fenêtres avec cadre (Windows, macOS, Linux) : WebView2 `ZoomFactor`, WKWebView `magnification` et le niveau de zoom de WebKitGTK acceptent toute valeur positive, et l’option de fenêtre `Zoom` transmettait déjà ces valeurs. Les fenêtres sans cadre conservent le minimum de 1.0, car leurs zones de déplacement et de redimensionnement ne sont pas adaptées au zoom ([#4590](https://github.com/wailsapp/wails/issues/4590)) ; les facteurs non positifs, NaN et Inf sont désormais ignorés au lieu d’être relevés à 1.0 — par @APshenkin

## v3.0.0 - 2026-10-26

## Modifications
- Les préversions nocturnes sont publiées sous la forme `vX.Y.0-nightly.YYYYMMDD` et donnent un aperçu de la prochaine version hebdomadaire. Les versions hebdomadaires incrémentent la version mineure (`vX.Y.0`), afin que `@latest` de Go soit toujours résolu vers une version hebdomadaire ou corrective ; la version du paquet npm `@wailsio/runtime` correspond toujours à celle de Wails, les versions nocturnes utilisant le dist-tag `nightly`
- Les modèles de projet épinglent `@wailsio/runtime` à la version de la CLI qui les a générés
- Wails v3 nécessite Go 1.26 ou une version ultérieure (Go 1.26.10 ou 1.27.3 et les versions ultérieures incluent le correctif de sécurité `crypto/x509`)

## Corrections
- Corriger `wails3 update cli`, qui ne trouvait pas les versions v3 : il ne prenait en compte que les tags v2 et effectuait l’installation depuis un chemin incorrect
- Corriger le plantage des applications Windows avec `panic: SetTextColor failed` lors du rendu de la barre de menus de l’application sur certaines versions de Windows : la barre de menus est dessinée dans un contexte de périphérique valide ou utilise le rendu par défaut, `UAHMENUITEM.Position` correspond désormais au `int` C et le rendu du menu ne provoque jamais de panique (#6311)
- Protéger la mise à jour préparée tant qu’un redémarrage est en attente : un téléchargement ou un second redémarrage lancé pendant l’arrêt ne peut ni supprimer la charge utile que l’assistant de permutation attend d’installer, ni entrer en concurrence pour celle-ci, par @julianstorer
- Corriger le plantage des applications Windows dans `applyDarkMode` lorsque l’initialisation de `pkg/w32` déchargeait `uxtheme.dll` alors que ses pointeurs vers des fonctions ordinales étaient encore utilisés : le module reste désormais chargé pendant toute la durée de vie du processus (#6356)
- Valider la charge utile préparée avant l’échange, puis l’application installée avant la suppression de la sauvegarde, afin qu’une mise à jour endommagée soit annulée au lieu de remplacer une application fonctionnelle, par @julianstorer

## Supprimé
- Supprimer l’intégration de Crowdin après la synchronisation finale des traductions dans la [PR](https://github.com/wailsapp/wails/pull/5736), par @github-actions[bot]

## Sécurité
- Bloquer les messages réservés `wails:` provenant de cadres enfants d’origine différente sous Linux (GTK3 et GTK4) en injectant un script de protection dans tous les cadres, ce qui atténue GHSA-2p5f-2gx8-mh8m
- Mettre à jour `github.com/klauspost/compress` vers 1.18.7 (GHSA-259r-337f-4rfw) et `golang.org/x/image` vers 0.45.0, puis compiler avec Go 1.26.10 afin d’intégrer le correctif de `crypto/x509` (GO-2026-6809)
- Mettre à jour `vite` vers 8.2.4 ou une version ultérieure dans le paquet d’exécution, l’assistant de configuration, les exemples et les modèles de projet (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Limitations connues
- **Linux :** la protection des cadres enfants ne peut pas empêcher un cadre enfant de *même origine* (par exemple un iframe `srcdoc` ou `about:blank`) d’atteindre le gestionnaire de messages natif par l’intermédiaire de son parent, car WebKitGTK ne fournit au gestionnaire aucune information sur le cadre. Les cadres d’origine différente et les cadres en bac à sable dépourvus de `allow-same-origin` sont bloqués. N’intégrez pas de contenu non fiable de même origine dans une fenêtre Wails.
- **Android (expérimental) :** le pont JavaScript `window.wails` est exposé à tous les cadres de la WebView, sans vérification du cadre principal ni de l’origine. Le contenu d’un iframe intégré, y compris celui d’une origine différente, peut appeler des méthodes liées et envoyer des messages réservés `wails:`. N’intégrez pas de contenu tiers non fiable dans la WebView de l’application. Ce comportement diffère de macOS et d’iOS, qui ignorent les messages `wails:` provenant de cadres enfants. Un correctif est prévu pour 3.0.1.
- **Toutes les plateformes :** les appels de méthodes liées provenant de cadres enfants transitent par HTTP vers `/wails/runtime`, et non par le gestionnaire de messages natif ; la protection ne les restreint donc pas. Considérez que tout iframe que vous intégrez a accès à vos méthodes Go liées.

## v3.0.0-beta.34 - 2026-10-20

## Corrigé
- Corriger la WebView Windows qui ne récupérait pas le focus lors de l’activation de la fenêtre, ainsi que la touche Tab qui restait confinée dans la WebView au lieu de déplacer le focus vers la fenêtre native puis de le ramener (#3783)
- Corriger les structures de tables virtuelles COM `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` et d’autres interfaces dérivées) qui omettaient la chaîne de méthodes héritées, ce qui dirigeait les appels vers le mauvais emplacement de la table virtuelle (#6129)
- Corriger les boîtes de dialogue de question Windows qui ignoraient les libellés personnalisés des boutons et ne déclenchaient jamais leurs fonctions de rappel ; les boutons personnalisés sont désormais affichés avec une TaskDialog (avec repli sur MessageBox et mise en correspondance positionnelle des fonctions de rappel lorsque comctl32 v6 n’est pas disponible) (#5794)
- Corriger le blocage de `wails3 dev` pendant 60 s avec Vite : le contrôle automatique de disponibilité du frontend utilise désormais une sonde TCP au lieu d’une requête HTTP GET qui n’envoyait aucun en-tête `Accept` et recevait une réponse 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Maintenir le fonctionnement des menus de la zone de notification lorsqu’ils sont mis à jour pendant leur ouverture dans la [PR](https://github.com/wailsapp/wails/pull/6301), par @taliesin-ai
- Gérer correctement les erreurs COM transitoires de WebView2 pendant les opérations d’exécution dans la [PR](https://github.com/wailsapp/wails/pull/6302), par @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Ajouté
- Le générateur de liaisons produit désormais un module `bindings/events` qui exporte, pour chaque événement détectable enregistré avec `application.RegisterEvent` au moyen d’un nom constant à la compilation, une constante de nom typée (par exemple `UserLogin` pour `user:login`), utilisable avec `Events.On`/`Emit` sans le plugin Vite (#5645)

## Corrigé
- Corriger les boîtes de dialogue de question Windows qui ignoraient les libellés personnalisés des boutons et ne déclenchaient jamais leurs fonctions de rappel ; les boutons personnalisés sont désormais affichés avec une TaskDialog (avec repli sur MessageBox et mise en correspondance positionnelle des fonctions de rappel lorsque comctl32 v6 n’est pas disponible) (#5794)
- Corrige le blocage de `wails3 dev` pendant 60 s avec Vite : le contrôle automatique de disponibilité du frontend utilise désormais une sonde TCP au lieu d’une requête HTTP GET qui n’envoyait aucun en-tête `Accept` et recevait une réponse 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Maintient le fonctionnement des menus de la zone de notification lorsqu’ils sont mis à jour pendant qu’ils sont ouverts dans [PR](https://github.com/wailsapp/wails/pull/6301) par @taliesin-ai
- Gère correctement les erreurs COM transitoires de WebView2 lors des opérations d’exécution dans [PR](https://github.com/wailsapp/wails/pull/6302) par @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Ajouts
- Ajoute des en-têtes de section natifs aux menus macOS pour macOS 14+ — consultez le [site de documentation](https://v3.wails.io) (20 pages mises à jour) dans [PR](https://github.com/wailsapp/wails/pull/6297) par @leaanthony

## Corrections
- Rend SetRelativePosition relatif à la zone de travail de l’écran sous macOS dans [PR](https://github.com/wailsapp/wails/pull/6299) par @leaanthony
- L’utilitaire de mise à jour relance la cible inchangée après l’expiration du délai d’attente du processus parent ou l’échec du remplacement dans [PR](https://github.com/wailsapp/wails/pull/6298) par @taliesin-ai
- Fige pnpm à la version 10 pour la documentation versionnée dans [PR](https://github.com/wailsapp/wails/pull/6293) par @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Ajouts
- Expose l’export du sous-chemin `@wailsio/runtime/cancellable` afin que `CancellablePromise` puisse être importé sans évaluer le point d’entrée racine du runtime (#6153)

## Modifications
- Nécessite Go 1.26 ou une version ultérieure : le module v3 déclare désormais `go 1.26.0`, et les projets générés par `wails3 init` et `wails3 service init` déclarent la même version minimale

## Corrections
- Corrige les variations de la sortie du générateur de liaisons lors de la compilation avec Go 1.27, où les types `encoding/json` tels que `Marshaler`, `Unmarshaler` et `RawMessage` devenaient des alias de types `encoding/json/v2` et `encoding/json/jsontext` : ils sont de nouveau rendus sous la forme `any` et aucun modèle n’est généré pour `encoding/json/v2` (#6292)
- Valide avant le clonage les URL de modèles distants transmises à `wails3 init -t` : seules les URL `https://`, `http://`, `ssh://`, `git://` et `user@host:path` sont acceptées, les valeurs sont transmises à git après `--`, et git est limité à ces transports pendant le clonage
- Rejette, dans le transport HTTP du runtime, les valeurs d’en-tête `x-wails-window-id` hors limites au lieu de les tronquer silencieusement en l’identifiant d’une autre fenêtre
- Renforce l’extraction des archives par l’utilitaire de mise à jour contre les attaques zip-slip : les chemins des entrées sont désormais vérifiés avec `filepath.IsLocal`, ainsi qu’au moyen d’un contrôle garantissant que le chemin résultant reste contenu dans le répertoire cible. Des tests couvrent les entrées imbriquées `..`, celles comportant des barres obliques inverses ou des chemins absolus, ainsi que celles permettant une écriture indirecte par lien symbolique dans les archives zip et tar.gz
- Corrige le plantage au démarrage, hors de Debian et Ubuntu, des AppImages provenant de `wails3 generate appimage` avec le message « Unable to spawn a new child process: … WebKitNetworkProcess ». La version intégrée de WebKit exécute désormais les processus auxiliaires fournis dans l’AppImage. Le bac à sable bubblewrap de WebKit reste activé
- Corrige l’échec de compilation des builds `-tags server` pour Android et iOS, ce qui résout [#6271](https://github.com/wailsapp/wails/issues/6271)
- Corrige l’échec de permutation de l’utilitaire de mise à jour sous Linux lorsque le répertoire temporaire se trouve sur un autre système de fichiers, tel qu’un tmpfs `/tmp`, en préparant les fichiers à côté du binaire installé et en les copiant entre les systèmes de fichiers
- Empêche l’utilitaire de mise à jour de tenter de remplacer le montage en lecture seule d’une AppImage au lieu du fichier `.AppImage`
- Empêche l’acceptation, depuis des iframes sous macOS et iOS, des messages réservés de contrôle de fenêtre `wails:`. Ils ne sont désormais acceptés que depuis le cadre principal de la fenêtre ; les autres messages provenant d’iframes parviennent toujours à `RawMessageHandler` avec les informations sur leur origine
- Arrête l’envoi de `WindowDidUpdate` et `WindowWillUpdate` sous macOS : `NSWindowDidUpdateNotification` se déclenche à chaque passage de la boucle d’événements, si bien que la page recevait environ 120 appels `evaluateJavaScript:` vides par seconde ; sous macOS, chacun d’eux dans 27.0 consommait l’activation transitoire du clic de l’utilisateur. `navigator.clipboard.writeText()` et les autres API soumises à une activation rejetaient alors les véritables clics avec `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), par @APshenkin
- Définit la langue par défaut du fichier info.json de Windows sur en-US (0409) afin qu’Explorer affiche les détails de version (#6210)
- Corrige la désactivation du moteur de rendu DMA-BUF sous Linux sur les ordinateurs portables équipés d’un GPU dédié NVIDIA aux côtés d’un GPU intégré AMD ou Intel : sous Wayland, il n’est désormais désactivé que lorsque WebKit effectue le rendu sur le GPU NVIDIA (parce que l’application lui a été déléguée ou que le compositeur s’y exécute), et non plus chaque fois que le module nvidia est chargé, ce qui ralentissait considérablement le rendu sur le GPU intégré. Les sessions X11 conservent l’ancien comportement, car le GPU de rendu de Xorg est défini par sa configuration
- Gère les échecs de renommage entre périphériques dans le programme de mise à jour Linux grâce à une copie intermédiaire de secours, dans la [PR](https://github.com/wailsapp/wails/pull/6135) par @john-okeefe
- Corrige `WebviewWindow.SetMenu` sous Windows, qui ne redessinait pas la barre de menus à moins qu’un thème de barre de menus ne soit configuré, ce qui corrige [#6103](https://github.com/wailsapp/wails/issues/6103)
- Corrige `WebviewWindow.SetMenu(nil)` sous Windows, qui provoquait une panique au lieu de rester sans effet comme sous macOS, ce qui corrige [#6104](https://github.com/wailsapp/wails/issues/6104)
- Préserve les entrées du journal des modifications lors de poussées concurrentes grâce à des nouvelles tentatives, dans la [PR](https://github.com/wailsapp/wails/pull/6094) par @leaanthony
- Corrige l’échec de `go mod vendor` avec `pattern arm64/WebView2Loader.dll: no matching files found` sur toutes les plateformes en supprimant les incorporations qui faisaient référence à des binaires jamais distribués dans le module ; corrige [#5782](https://github.com/wailsapp/wails/issues/5782) et [#5376](https://github.com/wailsapp/wails/issues/5376), dans la [PR](https://github.com/wailsapp/wails/pull/6031) par @Grantmartin2002
- Empêche la destruction en double des fenêtres sous macOS, dans la [PR](https://github.com/wailsapp/wails/pull/6159) par @charlesbalderson
- Gère WM_CONTEXTMENU afin d’ouvrir le menu de la zone de notification par un clic droit, dans la [PR](https://github.com/wailsapp/wails/pull/6162) par @wzk2025
- Corrige les cinq variantes de `InvokeSync*` qui bloquaient indéfiniment leur appelant lorsque la fonction de rappel distribuée provoquait une panique ; corrige [#6107](https://github.com/wailsapp/wails/issues/6107), dans la [PR](https://github.com/wailsapp/wails/pull/6108) par @RALIST
- Reconstruit en toute sécurité les sous-menus GTK3 détruits, dans la [PR](https://github.com/wailsapp/wails/pull/6236) par @leaanthony
- Implémente sous Windows une couche de transfert pour les workers imbriqués, dans la [PR](https://github.com/wailsapp/wails/pull/6261) par @taliesin-ai
- Corrige `OpenFile()` sous Windows, qui ignorait `PromptForMultipleSelection()` lorsque `CanChooseDirectories(true)` : les utilisateurs peuvent désormais sélectionner plusieurs dossiers à la fois, comme sous macOS et Linux (#6263)
- Corrige sous windows/arm64 les arguments `double` passés par valeur aux méthodes COM de WebView2 — `ZoomFactor`, `RasterizationScale`, le `Expires` des cookies, ainsi que l’échelle, la taille de page et les marges des paramètres d’impression étaient concernés (les méthodes de définition des paramètres d’impression renvoyaient auparavant `ErrDoubleArgUnsupported`). Les appels concernés passent désormais par un trampoline `runtime.cgocall` qui charge `d0`, ce que le chemin d’appel système de Go ne peut pas faire (golang.org/issue/62583), par @APshenkin

## Sécurité
- Le mode serveur (`-tags server`) n’accepte désormais les appels `/wails/runtime` que sous forme de requêtes `POST` provenant de la propre origine du serveur ou d’une origine répertoriée dans `ServerOptions.WebSocketOriginPatterns` ; les autres méthodes reçoivent la réponse HTTP 405 et les autres origines de navigateur la réponse HTTP 403. Les déploiements derrière un proxy inverse qui réécrit l’en-tête `Host` doivent ajouter leur origine publique à `WebSocketOriginPatterns`
- Le mode serveur (`-tags server`) lié à une adresse de bouclage (par défaut `localhost`, `127.0.0.0/8` ou `::1`) rejette désormais les requêtes avec la réponse HTTP 403, sauf si l’en-tête `Host` vaut `localhost`, contient une adresse IP littérale, correspond à l’hôte de liaison ou à un hôte conforme à `ServerOptions.WebSocketOriginPatterns`. Cela bloque les attaques par rebinding DNS contre `/wails/runtime` et les points de terminaison WebSocket
- Sous Windows, les notifications ne se rabattent plus sur PowerShell lorsque l’API de notifications toast WinRT échoue ; l’erreur est renvoyée à la place

## v3.0.0-beta.30 - 2026-10-15

## Corrigé
- Empêche le programme de mise à jour de tenter de remplacer le montage en lecture seule d’une AppImage au lieu du fichier `.AppImage`
- Arrête l’envoi de `WindowDidUpdate` et `WindowWillUpdate` sous macOS : `NSWindowDidUpdateNotification` se déclenche à chaque passage de la boucle d’événements, si bien que la page recevait environ 120 appels `evaluateJavaScript:` vides par seconde ; sous macOS, chacun d’eux dans 27.0 consommait l’activation transitoire du clic de l’utilisateur. `navigator.clipboard.writeText()` et les autres API soumises à une activation rejetaient alors les véritables clics avec `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), par @APshenkin
- Définit la langue par défaut du fichier info.json de Windows sur en-US (0409) afin qu’Explorer affiche les détails de version (#6210)
- Gère les échecs de renommage entre périphériques dans le programme de mise à jour Linux grâce à une copie intermédiaire de secours, dans la [PR](https://github.com/wailsapp/wails/pull/6135) par @john-okeefe
- Corrige `WebviewWindow.SetMenu` sous Windows, qui ne redessinait pas la barre de menus à moins qu’un thème de barre de menus ne soit configuré, ce qui corrige [#6103](https://github.com/wailsapp/wails/issues/6103)
- Corrige `WebviewWindow.SetMenu(nil)` sous Windows, qui provoquait une panique au lieu de rester sans effet comme sous macOS, ce qui corrige [#6104](https://github.com/wailsapp/wails/issues/6104)
- Préserve les entrées du journal des modifications lors de poussées concurrentes grâce à des nouvelles tentatives, dans la [PR](https://github.com/wailsapp/wails/pull/6094) par @leaanthony
- Corrige l’échec de `go mod vendor` avec `pattern arm64/WebView2Loader.dll: no matching files found` sur toutes les plateformes en supprimant les incorporations qui faisaient référence à des binaires jamais distribués dans le module ; corrige [#5782](https://github.com/wailsapp/wails/issues/5782) et [#5376](https://github.com/wailsapp/wails/issues/5376), dans la [PR](https://github.com/wailsapp/wails/pull/6031) par @Grantmartin2002
- Empêche la destruction en double des fenêtres sous macOS, dans la [PR](https://github.com/wailsapp/wails/pull/6159) par @charlesbalderson
- Gère WM_CONTEXTMENU afin d’ouvrir le menu de la zone de notification par un clic droit, dans la [PR](https://github.com/wailsapp/wails/pull/6162) par @wzk2025
- Corrige les cinq variantes de `InvokeSync*` qui bloquaient indéfiniment leur appelant lorsque la fonction de rappel distribuée provoquait une panique ; corrige [#6107](https://github.com/wailsapp/wails/issues/6107), dans la [PR](https://github.com/wailsapp/wails/pull/6108) par @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Ajouté
- Sous macOS, les raccourcis globaux tiennent désormais compte de la disposition du clavier : sur les dispositions autres que QWERTY (AZERTY, QWERTZ, Dvorak, ...), les accélérateurs alphabétiques sont associés à la touche portant la lettre correspondante via `UCKeyTranslate`, plutôt qu’à une position physique fixe du clavier QWERTY.

## Corrigé
- Initialise le menu natif avant de le définir sous Linux, dans la [PR](https://github.com/wailsapp/wails/pull/6233) par @abhicodes-007
- Corrige les régressions de préparation à la disponibilité générale sous macOS dans la gestion de la sortie et l’assistant de configuration, dans la [PR](https://github.com/wailsapp/wails/pull/6237) par @taliesin-ai
- Corrige les enveloppes COM `AddRef`/`Release` de WebView2 sous Windows, qui signalaient à tort comme un échec une ancienne dernière erreur du thread. Elles renvoient désormais le compteur de références, de sorte qu’une libération ne peut plus interrompre l’application ni faire échouer une requête de ressource (#3931)
- Corrige `ICoreWebView2Deferral.Release`, qui appelait l’emplacement `AddRef` et entraînait une fuite du report de chaque requête de ressource différée sous Windows (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Corrigé
- Conserve la sortie des applications macOS en dehors des fournisseurs de fichiers, dans la [PR](https://github.com/wailsapp/wails/pull/6225) par @taliesin-ai
- Finalise l’assistant de configuration et affiche les informations système, dans la [PR](https://github.com/wailsapp/wails/pull/6224) par @taliesin-ai
- Améliore la gestion des erreurs de récupération de WebView2 et le nettoyage des ressources sous Windows, dans la [PR](https://github.com/wailsapp/wails/pull/6222) par @leaanthony
- Amélioration de la récupération de WebView2 sous Windows afin d’éviter les plantages de l’application et l’abandon de contrôleurs dans la [PR](https://github.com/wailsapp/wails/pull/6220) par @taliesin-ai
- Correction des recommandations relatives à l’autorisation des médias sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/6221) par @taliesin-ai
- Correction du clic droit sur une icône de la zone de notification Linux, qui exécutait le gestionnaire de clic en plus d’ouvrir le menu, ce qui faisait basculer une fenêtre associée à chaque clic droit (#6018)
- Maintien en vie des gestionnaires de rappel WebView2 jusqu’à leur libération native sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/6184) par @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Corrections
- Correction de l’option `Permissions`, qui était ignorée sous macOS 12 et versions ultérieures : la méthode de capture multimédia `WKUIDelegate` est désormais implémentée, de sorte que `PermissionAllow` et `PermissionDeny` s’appliquent aux demandes d’accès à la caméra et au microphone comme sous Linux et Windows. `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` restent nécessaires, de même que les autorisations correspondantes pour les appareils lorsque l’application est exécutée dans un bac à sable

## v3.0.0-beta.26 - 2026-09-25

## Ajouts
- Le mécanisme de surveillance applique des contrôles de disponibilité HTTP au serveur de développement frontend dans la [PR](https://github.com/wailsapp/wails/pull/6150) par @atterpac

## Corrections
- Récupération après les défaillances de processus WebView2 au lieu de laisser une fenêtre vide dans la [PR](https://github.com/wailsapp/wails/pull/6002) par @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Corrections
- Arrêt de `wails3 dev` et de ses processus en arrière-plan lors de la fermeture de l’application principale (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Sécurité
- Restriction des autorisations de GITHUB_TOKEN dans les workflows de compilation et de publication afin de respecter le principe du moindre privilège
- Correction des traversées de chemin et des échappements par lien symbolique dans le middleware de ressources de l’exemple d’écran
- Maintien du fonctionnement de l’installation des dépendances par l’assistant de configuration, tout en la limitant aux commandes de gestionnaires de paquets approuvées

## v3.0.0-beta.23 - 2026-09-16

## Ajouts

- Possibilité de remplacer l’identifiant d’application GTK sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/5972) par @overlordtm

## Corrections

- Sous Linux, l’instance unique réserve le nom de bus sous l’identifiant propre à l’application dans la [PR](https://github.com/wailsapp/wails/pull/5973) par @overlordtm
- Annulation des requêtes de schéma personnalisé abandonnées sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/6126) par @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Ajouts

- Achèvement des traductions et portage des contenus multimédias de la communauté — consultez le [site de documentation](https://v3.wails.io) (1812 pages mises à jour) dans la [PR](https://github.com/wailsapp/wails/pull/6122) par @leaanthony
- Documentation des boutons Wails intégrables, avec différents styles et différentes tailles, dans la [PR](https://github.com/wailsapp/wails/pull/6120) par @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Ajouts

- Diffusion de la documentation de Wails v3 avec M-Press dans la [PR](https://github.com/wailsapp/wails/pull/6116) par @leaanthony

## Corrections

- Analyse des valeurs JSON de slug dans le frontmatter MPD pour générer le journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/6118) par @leaanthony
- Le programme de mise à jour efface les variables d’environnement de l’utilitaire et relance la cible d’origine après un échec de sauvegarde dans la [PR](https://github.com/wailsapp/wails/pull/6080) par @cnmax
- Démarrage du gestionnaire de signaux par défaut pendant App.Run dans la [PR](https://github.com/wailsapp/wails/pull/6098) par @leaanthony
- Sous Windows, le menu prend en charge les menus nil, libère les ressources remplacées et redessine la barre de menus dans la [PR](https://github.com/wailsapp/wails/pull/6112) par @taliesin-ai
- Rétablissement de la création de paquets MSIX pour les nouveaux projets utilisant la configuration YAML partagée dans la [PR](https://github.com/wailsapp/wails/pull/6115) par @leaanthony
- Correction de l’échec du chargement des bindings JavaScript et TypeScript générés lorsque des créateurs de modèles génériques font référence à des déclarations d’utilitaires ultérieures, et prévention des débordements de pile lors de la création de modèles génériques mutuellement dépendants (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Modifications

- Mise à jour des liens de présentation de Clave vers le site Web et le dépôt actuels dans la [PR](https://github.com/wailsapp/wails/pull/6082) par @01xR4in

## Corrections

- Annulation des requêtes de ressources Windows abandonnées, y compris celles des workers, tout en préservant les gestionnaires keepalive lors de la navigation. Transmission des contextes de requête natifs via l’enveloppe de l’application sur les plateformes Apple. (#5963, #5969)
- Conservation des entrées du journal des modifications lors de push concurrents grâce à de nouvelles tentatives dans la [PR](https://github.com/wailsapp/wails/pull/6094) par @leaanthony
- Correction de l’échec de `go mod vendor` avec `pattern arm64/WebView2Loader.dll: no matching files found` sur toutes les plateformes, en supprimant les éléments incorporés qui faisaient référence à des binaires jamais distribués dans le module, ce qui corrige [#5782](https://github.com/wailsapp/wails/issues/5782) et [#5376](https://github.com/wailsapp/wails/issues/5376), dans la [PR](https://github.com/wailsapp/wails/pull/6031) par @Grantmartin2002

## Suppressions

- Suppression de la prise en charge du chargeur WebView2 natif, remplacé par le chargeur en Go pur. Cela supprime les binaires `WebView2Loader.dll` incorporés et la dépendance `github.com/jchv/go-winloader`. Le tag de compilation `native_webview2loader` est toujours accepté et ne provoque plus d’erreur, mais il n’a aucun effet sur les compilations v3, dans la [PR](https://github.com/wailsapp/wails/pull/6031) par @Grantmartin2002
- Suppression des tags de compilation inutilisés et de l’option FPS du guide de l’API macOS dans la [PR](https://github.com/wailsapp/wails/pull/6097) par @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Ajouts

- Protection des API macOS privées derrière des tags de compilation pour une utilisation facultative — consultez la [documentation](https://v3.wails.io/features/browser/integration), la [documentation](https://v3.wails.io/features/environment/info), la [documentation](https://v3.wails.io/features/windows/basics), la [documentation](https://v3.wails.io/features/windows/frameless), la [documentation](https://v3.wails.io/features/windows/notch-windows), la [documentation](https://v3.wails.io/features/windows/options), la [documentation](https://v3.wails.io/guides/build/macos), la [documentation](https://v3.wails.io/guides/build/private-macos-apis) et la [documentation](https://v3.wails.io/reference/overview), dans la [PR](https://github.com/wailsapp/wails/pull/6087) par @leaanthony

## Corrections

- Rejet des requêtes d’exécution dépassant 64 Mio avec le code HTTP 413 dans la [PR](https://github.com/wailsapp/wails/pull/6091) par @leaanthony

## Sécurité

- Renforcement de la sécurité des origines MCP et de l’accès distant grâce à l’authentification par jeton dans la [PR](https://github.com/wailsapp/wails/pull/6092) par @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Corrections

- Correction de la fuite de mémoire de Calloc sous Linux et Darwin grâce à l’utilisation de récepteurs de pointeur dans la [PR](https://github.com/wailsapp/wails/pull/6083) par @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Corrections

- Windows : un `GetRequest` ayant échoué ou nil dans le gestionnaire WebResourceRequested n’arrête plus le processus (`log.Fatal` / panique due au déréférencement d’une valeur nil) — la requête est désormais abandonnée et journalisée dans la [PR](https://github.com/wailsapp/wails/pull/6006) par @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Modifications

- Demande du mot de passe de notarisation dans une nouvelle fenêtre de terminal dans la [PR](https://github.com/wailsapp/wails/pull/6029) par @leaanthony

## Corrigé

- Gestion correcte des types de clic sur l’icône de la zone de notification sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5919) par @ChewbaccaCookie
- La CI supprime les dépôts apt Microsoft inutilisés avant la mise à jour dans la [PR](https://github.com/wailsapp/wails/pull/6041) par @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Corrigé

- Augmentation à 60 secondes du délai d’expiration de l’intégration de WebView2 dans la [PR](https://github.com/wailsapp/wails/pull/6043) par @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Corrigé

- Nommage correct des frappes Contrôle-lettre sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/6032) par @taliesin-ai
- Correction des icônes ICO de la zone de notification et adaptation au thème de la barre des tâches sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/6016) par @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Corrigé

- Poursuite du traitement des tâches du thread principal sous macOS pendant l’exécution d’une boucle modale dans la [PR](https://github.com/wailsapp/wails/pull/6026) par @leaanthony
- Le stockage sécurisé mobile peut désormais échouer et se ferme de manière sécurisée en cas d’échec dans la [PR](https://github.com/wailsapp/wails/pull/5923) par @mortenolsrud
- Exécution des hooks d’événements de l’application même lorsqu’aucun écouteur n’est enregistré dans la [PR](https://github.com/wailsapp/wails/pull/5999) par @archy-rock3t-cloud
- Correction de fautes dans les commentaires et la documentation localisée dans la [PR](https://github.com/wailsapp/wails/pull/6023) par @haoku123
- Suppression des binaires macOS précompilés validés sous `v3/examples` dans la [PR](https://github.com/wailsapp/wails/pull/6025) par @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Ajouté

- Ajout de fenêtres de notification pour l’encoche de macOS, avec un exemple de cycle de vie et de télémétrie — voir la [documentation](https://v3.wails.io/features/windows/notch-windows) dans la [PR](https://github.com/wailsapp/wails/pull/6010) par @leaanthony
- Ajout de la prise en charge des fenêtres NSPanel sous macOS, avec de nouvelles options et une intégration native — voir la [documentation](https://v3.wails.io/features/windows/options) dans la [PR](https://github.com/wailsapp/wails/pull/6008) par @leaanthony

## Corrigé

- Prévention du blocage de SQLite Prepare lors d’appels simultanés dans la [PR](https://github.com/wailsapp/wails/pull/5998) par @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Supprimé

- Suppression de l’outil obsolète de suivi de l’implémentation dans la documentation, dans la [PR](https://github.com/wailsapp/wails/pull/6005) par @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Corrigé

- Correction de l’hôte Linux GTK4 qui ignorait les arguments de lancement d’un protocole personnalisé et d’une association de fichiers dans la [PR](https://github.com/wailsapp/wails/pull/6000) par @midagedev
- Gestion correcte des lignes supprimées et des corrections provenant de la même source lors de la validation du journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/5993) par @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Ajouté

- Ajout du serveur sécurisé wails3 mcp pour la gestion de projets assistée par des agents dans la [PR](https://github.com/wailsapp/wails/pull/5896) par @leaanthony
- Ajout d’une documentation sur les modèles dans les bindings — voir la [documentation](https://v3.wails.io/features/bindings/models) dans la [PR](https://github.com/wailsapp/wails/pull/5988) par @taliesin-ai
- Prise en charge de l’installation avec rpm-ostree sur les systèmes Linux atomiques dans la [PR](https://github.com/wailsapp/wails/pull/5987) par @leaanthony
- Ajout de la génération et de la publication natives et hebdomadaires du graphique de l’historique des étoiles — voir la [documentation](https://v3.wails.io/credits), la [documentation](https://v3.wails.io/de/credits), la [documentation](https://v3.wails.io/fr/credits), la [documentation](https://v3.wails.io/id/credits), la [documentation](https://v3.wails.io/ja/credits), la [documentation](https://v3.wails.io/ko/credits), la [documentation](https://v3.wails.io/pt/credits), la [documentation](https://v3.wails.io/ru/credits), la [documentation](https://v3.wails.io/zh-cn/credits) et la [documentation](https://v3.wails.io/zh-tw/credits) dans la [PR](https://github.com/wailsapp/wails/pull/5986) par @leaanthony
- Ajout du package mac réservé à Darwin pour résoudre les ressources du bundle de l’application — voir la [documentation](https://v3.wails.io/guides/build/macos) dans la [PR](https://github.com/wailsapp/wails/pull/5965) par @leaanthony
- Ajout d’une page de présentation de Condui et d’une entrée dans l’index — voir la [documentation](https://v3.wails.io/community/showcase/condui) et la [documentation](https://v3.wails.io/community/showcase) dans la [PR](https://github.com/wailsapp/wails/pull/5962) par @mgueregath
- Ajout d’une page de présentation de Redis Viewer avec des captures d’écran et un lien vers le projet — voir la [documentation](https://v3.wails.io/community/showcase) et la [documentation](https://v3.wails.io/community/showcase/redisviewer) dans la [PR](https://github.com/wailsapp/wails/pull/5984) par @redisviewer

## Modifié

- Mise à jour des indicateurs de l’application GTK vers G<em>APPLICATION</em>NON_UNIQUE sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/5971) par @overlordtm
- Journalisation des événements de fenêtre manquants au niveau débogage plutôt qu’au niveau avertissement dans la [PR](https://github.com/wailsapp/wails/pull/5914) par @julianstorer

## Corrigé

- Les raccourcis clavier enregistrés sous macOS peuvent désormais prendre le pas sur la webview dans la [PR](https://github.com/wailsapp/wails/pull/5902) par @julianstorer
- Réparation des liens rompus dans la barre latérale de la documentation dans la [PR](https://github.com/wailsapp/wails/pull/5937) par @northes
- Annulation des contextes de requêtes de ressources sous macOS et iOS lorsque WebKit abandonne la tâche correspondante du schéma personnalisé (#5963)
- Gestion du message WindowSetFullscreenButtonEnabled dans la [PR](https://github.com/wailsapp/wails/pull/5976) par @archy-rock3t-cloud
- Importation de Fragment dans le modèle preact-ts pour résoudre l’échec de compilation dans la [PR](https://github.com/wailsapp/wails/pull/5979) par @haoku123
- Prévention du plantage des anciennes applications GTK3 uniquement composées de services lorsque la détection des écrans s’exécute avant qu’une fenêtre ou un affichage actif soit disponible (#5966)
- Les exécutions de publication avec une version explicite peuvent désormais se poursuivre lorsque le journal des modifications non publiées est vide (#5977)

## Sécurité

- Mise à jour des fichiers de verrouillage nanoid du site web vers la version corrigée 3.3.18 afin de résoudre les avis de sécurité dans la [PR](https://github.com/wailsapp/wails/pull/5985) par @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Ajouté

- Ajout de la génération des URL de documentation aux entrées automatiques du journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/5957) par @taliesin-ai
- Ajout de Streams : des flux d’octets bidirectionnels entre Go et JavaScript, avec le modèle de programmation WebSocket et sans socket d’écoute. Déclarez un flux en Go avec `app.HandleStream(name, handler)` et connectez-vous depuis le frontend avec `Stream(name)`, qui renvoie un objet de la forme `WebSocket`. Le sens Go→JS utilise une interrogation maintenue par fenêtre via le serveur de ressources, tandis que le sens JS→Go utilise une requête POST normale ; aucun port TCP n’est lié et rien ne transite par `evaluateJavaScript`. Dans les builds serveur (`-tags server`), le même gestionnaire est fourni à la place via un véritable WebSocket, de sorte que le code de l’application reste identique dans tous les builds. Par @leaanthony
- Déplacement de l’entrée du journal des modifications relative à la boîte aux lettres vers la section Non publié dans la [PR](https://github.com/wailsapp/wails/pull/5935) par @leaanthony

## Modifié

- Mise à jour de la génération automatique de la barre latérale de la documentation et de la dérivation du type des auteurs du blog dans la [PR](https://github.com/wailsapp/wails/pull/5938) par @leaanthony

## Corrigé

- L’initialisation de WebView2 utilise une échéance et une boucle de messages dans la [PR](https://github.com/wailsapp/wails/pull/5952) par @leaanthony
- Le test des cookies WebView2 est ignoré dans la CI sauf activation explicite et son exécution est verrouillée sur le thread actuel du système d’exploitation dans la [PR](https://github.com/wailsapp/wails/pull/5951) par @leaanthony
- Les générateurs de menus Windows rétablissent les identifiants de commande des éléments parents de sous-menus dans la [PR](https://github.com/wailsapp/wails/pull/5944) par @gilad-ch
- Aligner l’image officielle de compilation croisée sur le niveau de référence de prise en charge de GTK 4.14+ sous Linux (#5928)
- Configurer le projet Xcode iOS pour conserver les indicateurs hérités de l’éditeur de liens et ajouter -ObjC dans la [PR](https://github.com/wailsapp/wails/pull/5915) par @mortenolsrud
- Corriger le renouvellement excessif des connexions TCP dans le proxy de ressources `wails3 dev` avec les interfaces volumineuses, qui pouvait épuiser les ports éphémères de l’hôte et provoquer l’échec de processus sans rapport avec `EADDRNOTAVAIL`
- Mettre en file d’attente le JavaScript des événements de chaque fenêtre afin d’assurer leur distribution ordonnée et la contre-pression dans la [PR](https://github.com/wailsapp/wails/pull/5934) par @leaanthony

## Supprimé

- Supprimer le pipeline de publication des binaires pour ordinateur de bureau : les versions v3 se limitent désormais aux balises et la CLI `wails3` s’installe avec `go install`. Supprime `release-v3.yml` ainsi que l’étape nocturne qui le déclenchait dans la [PR](https://github.com/wailsapp/wails/pull/5946) par @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Ajouté

- Ajouter une préférence de lecture automatique sous macOS pour désactiver l’obligation d’une action de l’utilisateur avant la lecture de contenus multimédias dans la [PR](https://github.com/wailsapp/wails/pull/5512) par @Eyalm321
- Déplacer l’entrée de journal des modifications relative à la boîte aux lettres vers la section Non publié dans la [PR](https://github.com/wailsapp/wails/pull/5935) par @leaanthony

## Modifié

- L’animation de zoom sous macOS utilise CADisplayLink ou NSTimer pour offrir de meilleures performances dans la [PR](https://github.com/wailsapp/wails/pull/5945) par @savely-krasovsky

## Corrigé

- Configurer le projet Xcode iOS pour conserver les indicateurs hérités de l’éditeur de liens et ajouter -ObjC dans la [PR](https://github.com/wailsapp/wails/pull/5915) par @mortenolsrud
- Corriger le renouvellement excessif des connexions TCP dans le proxy de ressources `wails3 dev` avec les interfaces volumineuses, qui pouvait épuiser les ports éphémères de l’hôte et provoquer l’échec de processus sans rapport avec `EADDRNOTAVAIL`
- Mettre en file d’attente le JavaScript des événements de chaque fenêtre afin d’assurer leur distribution ordonnée et la contre-pression dans la [PR](https://github.com/wailsapp/wails/pull/5934) par @leaanthony

### Ajouté

- Implémenter une boîte aux lettres FIFO asynchrone générique pour la distribution ordonnée des événements dans la [PR](https://github.com/wailsapp/wails/pull/5851) par @savely-krasovsky et @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Ajouté

- Implémenter un stockage borné côté hôte pour les événements surdimensionnés ainsi que leur distribution JavaScript ordonnée dans la [PR](https://github.com/wailsapp/wails/pull/5930) par @leaanthony
- Implémenter le rebond de l’icône dans le Dock macOS lors du clignotement d’une fenêtre dans la [PR](https://github.com/wailsapp/wails/pull/5921) par @julianstorer

## Corrigé

- Le serveur de ressources conserve les erreurs de détection du type de contenu et les préfixes non écrits lors du vidage dans la [PR](https://github.com/wailsapp/wails/pull/5931) par @leaanthony
- Empêcher les applications macOS de planter lors du remplacement du menu de l’application depuis un rappel Wails
- Corriger les menus natifs illisibles sous Windows 10 1809 / Windows Server 2019 (build 17763). Les exportations uxtheme du mode sombre étaient conditionnées au build 18334 ; l’activation explicite du mode sombre au niveau de l’application ne s’exécutait donc jamais sur ces hôtes : l’arrière-plan du menu était peint en sombre, mais Windows continuait à dessiner le texte du menu avec le thème clair, ce qui produisait du texte sombre sur un fond sombre. Les ordinaux existent à partir de 17763 ; la condition correspond désormais à cette version.
- Corriger `w32.GetStockObject`, qui appelait `GetDeviceCaps` au lieu de `GetStockObject` et renvoyait ainsi 0 pour chaque objet prédéfini.
- Améliorer la gestion et le signalement des erreurs de téléchargement du programme d’amorçage WebView2 dans la [PR](https://github.com/wailsapp/wails/pull/5924) par @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Corrigé

- Sous macOS, l’activation de l’application respecte désormais la politique d’activation uniquement pour les applications standard dans la [PR](https://github.com/wailsapp/wails/pull/5897) par @julianstorer
- Protéger les fenêtres GTK non initialisées dans les builds Linux dans la [PR](https://github.com/wailsapp/wails/pull/5898) par @julianstorer
- Définir une couleur d’arrière-plan opaque explicite pour les fenêtres WebKit sous Linux avant le chargement de l’URL dans la [PR](https://github.com/wailsapp/wails/pull/5899) par @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Modifié

- Les tâches de build Android ciblent arm64 par défaut et deploy-emulator sélectionne l’architecture de l’hôte dans la [PR](https://github.com/wailsapp/wails/pull/5890) par @mortenolsrud

## Corrigé

- Conserver l’état de zoom des fenêtres macOS pendant leur déplacement et réduire les animations dans la [PR](https://github.com/wailsapp/wails/pull/5900) par @leaanthony
- Corriger le build Windows en mode serveur en ajoutant `!server` à la contrainte de build `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Ajouté

- Documenter l’achèvement de la vérification bêta de la phase 10 dans les détails d’implémentation, dans la [PR](https://github.com/wailsapp/wails/pull/5881) par @leaanthony

## Corrigé

- Transmettre le handle de fenêtre à l’API de mode sombre de Windows et valider les arguments dans la [PR](https://github.com/wailsapp/wails/pull/5877) par @leaanthony
- Centraliser la résolution de l’état des boutons de barre de titre pour les fenêtres sans cadre sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5870) par @taliesin-ai
- Éviter que le texte des menus natifs devienne illisible lorsqu’une application Windows demande le mode sombre alors que le thème des applications Windows est clair. Le menu utilise désormais l’arrière-plan natif clair correspondant jusqu’à ce que Windows puisse afficher le texte du menu en sombre.
- Corriger les menus natifs illisibles sous Windows 10 1809 / Windows Server 2019 (build 17763). Les exportations uxtheme du mode sombre étaient conditionnées au build 18334 ; l’activation explicite du mode sombre au niveau de l’application ne s’exécutait donc jamais sur ces hôtes : l’arrière-plan du menu était peint en sombre, mais Windows continuait à dessiner le texte du menu avec le thème clair, ce qui produisait du texte sombre sur un fond sombre. Les ordinaux existent à partir de 17763 ; la condition correspond désormais à cette version.

## v3.0.0-beta.2 - 2026-08-02

## Modifié

- Faire passer v3 du stade alpha au stade bêta
- Documenter les valeurs par défaut intelligentes de la zone de notification et le masquage automatique des menus contextuels, avec des tests de non-régression couvrant la sélection du gestionnaire de clics (#5840).
- Le programme de mise à jour GitHub exclut par défaut les ressources d’installation Windows dans la [PR](https://github.com/wailsapp/wails/pull/5861) par @leaanthony
- Ajouter la prise en charge des angles arrondis, droits et à rayon personnalisé pour les fenêtres sans cadre sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5866) par @leaanthony

## Corrigé

- Signaler les dimensions en temps réel des fenêtres GTK4 et émettre depuis la surface configurée les événements de redimensionnement, de maximisation, de minimisation et de passage en plein écran (#5830).
- Corriger le plantage de WebKit sous Linux lors de l’envoi d’un Blob ou de FormData dans les requêtes fetch dans la [PR](https://github.com/wailsapp/wails/pull/5854) par @taliesin-ai
- Le shim de fetch transmet undefined lorsque les en-têtes Blob/FormData sont absents dans la [PR](https://github.com/wailsapp/wails/pull/5865) par @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Ajouté

## Modifié

- Ajouter la prise en charge des angles arrondis, droits et à rayon personnalisé pour les fenêtres sans cadre sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5866) par @leaanthony

## Corrigé

- Le shim de fetch transmet undefined lorsque les en-têtes Blob/FormData sont absents dans la [PR](https://github.com/wailsapp/wails/pull/5865) par @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Ajouté

- Ajout de la prise en charge de la création de paquets DMG pour macOS, avec de nouvelles options et tâches de build, dans la [PR](https://github.com/wailsapp/wails/pull/5857) par @leaanthony

## Modifié

- L’outil de mise à jour GitHub exclut par défaut les ressources des programmes d’installation Windows dans la [PR](https://github.com/wailsapp/wails/pull/5861) par @leaanthony

## Corrigé

- Correction du plantage de WebKit sous Linux lors de l’envoi de données Blob ou FormData dans des requêtes fetch, dans la [PR](https://github.com/wailsapp/wails/pull/5854) par @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Ajouté

- Implémentation sous macOS des actions de maximisation ou de réduction lors d’un double-clic sur la barre de titre, dans la [PR](https://github.com/wailsapp/wails/pull/5853) par @taliesin-ai

## Modifié

- Mise à jour du tutoriel sur le service QR afin d’utiliser NewServiceWithOptions et ajout d’espacement dans la [PR](https://github.com/wailsapp/wails/pull/5849) par @jeongkyu

## Corrigé

- Maintien de la réactivité de WKWebView pendant le zoom sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5856) par @leaanthony
- Correction des requêtes de taille de fenêtre avec GTK4 et émission des événements de redimensionnement, de maximisation, de réduction et de changement d’état plein écran depuis le `GdkSurface` configuré.

## v3.0.0-alpha2.119 - 2026-07-27

## Corrigé

- Mise à jour de la documentation en plusieurs langues pour y inclure des diagrammes d’architecture, dans la [PR](https://github.com/wailsapp/wails/pull/5833) par @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Ajouté

- Ajout de chemins par défaut pour les entrées et les sorties de la génération d’icônes dans la [PR](https://github.com/wailsapp/wails/pull/5825) par @taliesin-ai
- Ajout des modules d’entrée sources à sideEffects dans le package.json du runtime, dans la [PR](https://github.com/wailsapp/wails/pull/5797) par @savely-krasovsky
- Ajout d’une section sur la licence et la provenance au guide de contribution dans la [PR](https://github.com/wailsapp/wails/pull/5816) par @taliesin-ai

## Corrigé

- Application d’une CSS sans cadre limitée à GTK4 afin de supprimer le rayon des bordures, dans la [PR](https://github.com/wailsapp/wails/pull/5800) par @savely-krasovsky
- Gestion sans erreur fatale des échecs de récupération de la position du curseur sous Windows pour les menus contextuels et l’énumération des écrans, dans la [PR](https://github.com/wailsapp/wails/pull/5789) par @wayneforrest
- Sous macOS, la boîte de dialogue d’ouverture de fichier filtre correctement les extensions et valide les fichiers autorisés d’après leur suffixe, dans la [PR](https://github.com/wailsapp/wails/pull/5678) par @phergul
- Protection de l’initialisation du mode sombre sous Windows contre les appels à une API nil, dans la [PR](https://github.com/wailsapp/wails/pull/5793) par @roachadam
- Correction de l’échec du build 32 bits dans l’outil de mise à jour : la constante `maxArchiveTotalSize` (2 Gio) dépassait la capacité du type `int` de la plateforme lorsqu’elle était transmise à `fmt.Errorf` sur `GOARCH=386`. Elle est désormais explicitement typée `int64`.
- Correction d’une panique due à un pointeur nil au démarrage lorsqu’une fenêtre utilise une barre de titre sombre (ou sombre selon le système) dans des builds Windows qui ne chargent pas les API uxtheme du mode sombre, comme Windows 10 1809 / Windows Server 2019 (build 17763). Les appels à `AllowDarkModeForWindow` effectués lors de la configuration du thème de la fenêtre sont désormais protégés contre les valeurs nil, comme c’était déjà le cas dans `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Ajouté

- Implémentation d’une logique personnalisée de test de positionnement pour les zones non clientes sous Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5462) par @savely-krasovsky

## Modifié

- Configuration de la détection de l’échelle du moniteur par WebView2 en fonction de UseVisualHosting, dans la [PR](https://github.com/wailsapp/wails/pull/5761) par @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Ajouté

- Mise à jour de la documentation de la FAQ pour la recentrer sur les fonctionnalités et les recommandations de Wails v3, dans la [PR](https://github.com/wailsapp/wails/pull/5763) par @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Corrigé

- Correction de `Menu.Update()` qui ne reconstruisait pas le menu natif sous Linux avec GTK4 (#5659, problème diagnostiqué et corrigé indépendamment par @puneetdixit200 dans #5539)
- Correction d’un plantage lors de l’énumération des écrans macOS après un changement d’affichage, en copiant les chaînes d’identifiant et de nom des écrans et en prenant un instantané de leur nombre (#5565, problème diagnostiqué et corrigé indépendamment par @x-haose dans #5584)
- Correction d’un plantage sous Windows lorsque `WM_ERASEBKGND` peint un arrière-plan uni pendant une transition de réduction/restauration au cours de laquelle `GetClientRect` renvoie nil (protection signalée par @sinspired dans #5636)
- Correction par @mbaklor, dans #5690, d’une erreur qui faisait que les erreurs des bindings frontend étaient toujours analysées comme du texte
- Correction de l’échec du build Windows lors de l’utilisation du tag de build `server`, dû à l’absence, dans les fichiers de l’interface graphique Windows, de la contrainte de build `!server` déjà présente dans leurs équivalents macOS et Linux (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Ajouté

- Implémentation du protocole Update Manifest et du fournisseur de points de terminaison dans la [PR](https://github.com/wailsapp/wails/pull/5720) par @taliesin-ai

## Modifié

- Intégration du binding `webview2` au module v3 sous le nom `v3/internal/webview2`, avec suppression du module autonome, de ses workflows de publication/synchronisation nightly et du jeu de versions dans go.mod (v3 en est le seul consommateur), dans la [PR](https://github.com/wailsapp/wails/pull/5711) par @taliesin-ai

## Corrigé

- Déplacement de la détection de l’échelle du moniteur par WebView2 et du correctif de resynchronisation de l’hôte lors d’un changement de DPI vers la section Non publié, dans la [PR](https://github.com/wailsapp/wails/pull/5750) par @taliesin-ai
- Mise à jour du marshaling COM de WebView2 pour les paramètres float64 et BOOL, dans la [PR](https://github.com/wailsapp/wails/pull/5741) par @wayneforrest
- Prévention des paniques et des déréférencements de nil lors de la mise à jour et de la destruction des icônes de la zone de notification Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5703) par @wayneforrest
- Correction des fenêtres masquées qui ne se remasquaient pas correctement sous Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5743) par @wayneforrest
- Synchronisation de la visibilité du contrôleur WebView2 avec la réduction, la maximisation et la restauration de la fenêtre, dans la [PR](https://github.com/wailsapp/wails/pull/5742) par @wayneforrest

### Corrigé

- Réactivation de la détection de l’échelle du moniteur par WebView2 et conditionnement de la resynchronisation de l’hôte au changement de DPI, dans la [PR](https://github.com/wailsapp/wails/pull/5734) par @taliesin-ai, d’après le correctif validé par @randalmurphal, avec vérification de la cause racine par @eleclin et tests matériels par @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Ajouté

- Ajout d’un avertissement lors de la création d’un AAB de publication sans que `ANDROID_KEYSTORE_FILE` soit défini (Google Play rejette les bundles signés pour le débogage), et documentation de la création et de la signature des paquets App Bundle, dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai
- Ajout de la documentation « Pourquoi Wails » en plusieurs langues dans la [PR](https://github.com/wailsapp/wails/pull/5739) par @taliesin-ai
- Prise en charge du mappage de Go time.Time vers JS Date ou string dans les bindings, dans la [PR](https://github.com/wailsapp/wails/pull/5398) par @fbbdev
- Ajout de tâches de création de paquets Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) pour la soumission au Play Store — les tâches APK restent destinées aux tests locaux et sur émulateur dans la [PR](https://github.com/wailsapp/wails/pull/5728) par @mortenolsrud (corrige [#5726](https://github.com/wailsapp/wails/issues/5726))
- Ajout de cibles de tâches pour les appareils Android physiques et reprise des autorisations d’accès à la caméra et à la localisation dans la [PR](https://github.com/wailsapp/wails/pull/5735) par @taliesin-ai

## Modifications

- Mise à niveau de `webview2` vers v1.0.28 ([notes de version](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Mise à niveau des éléments `compileSdk`/`targetSdk` du modèle Android de 34 vers 35, version exigée par Google Play pour les nouvelles soumissions d’applications, dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai

## Corrections

- Correction des masques d’avatar préintégrés dans sponsorkit dans la [PR](https://github.com/wailsapp/wails/pull/5745) par @leaanthony
- Correction de la création automatique des AVD Android, qui sélectionnait une image système ou une version de cmdline-tools incorrecte en raison du tri lexicographique des versions, dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai
- Correction de l’assistant de configuration, qui suggérait une version obsolète du NDK Android (désormais 26.3.11579264, conformément à l’exigence documentée), dans la [PR](https://github.com/wailsapp/wails/pull/5730) par @taliesin-ai
- Mise à jour de la documentation française relative à SvelteKit et aux options dans la [PR](https://github.com/wailsapp/wails/pull/5744) par @leaanthony
- Correction d’une erreur SIGSEGV lors de l’énumération des écrans sous macOS pendant les changements d’affichage dans la [PR](https://github.com/wailsapp/wails/pull/5516) par @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Ajouts

- Ajout d’un générateur SVG des contributeurs écrit en Go et mise à jour des pages de crédits de la documentation et du site web dans la [PR](https://github.com/wailsapp/wails/pull/5724) par @taliesin-ai
- Prise en charge du mappage de Go time.Time vers JS Date ou string dans les liaisons dans la [PR](https://github.com/wailsapp/wails/pull/5398) par @fbbdev

## Modifications

- Remplacement du pipeline d’images des sponsors basé sur Node par un générateur Go dans la [PR](https://github.com/wailsapp/wails/pull/5719) par @taliesin-ai

## Corrections

- Correction du script d’installation des dépendances des ressources de compilation Android dans la [PR](https://github.com/wailsapp/wails/pull/5729) par @taliesin-ai
- Rejet du caractère de contrôle U+0085 (NEXT LINE) dans `ValidateAndSanitizeURL`, ce qui complète la prise en charge des caractères d’espacement par le validateur d’URL
- Recalcul du cadre DWM lors d’un changement de DPI pour les fenêtres sans cadre dans la [PR](https://github.com/wailsapp/wails/pull/4785) par @leaanthony
- Correction de la détection des zones de dépôt DnD, qui échouait sous Windows avec une mise à l’échelle différente de 100 %, dans la [PR](https://github.com/wailsapp/wails/pull/4632) par @yulesxoxo
- Ajout d’une gestion explicite de la mémoire Objective-C pour les objets Cocoa utilisés dans les boîtes de dialogue, les menus, la zone de notification et les notifications sous Darwin, dans la [PR](https://github.com/wailsapp/wails/pull/5714) par @taliesin-ai
- Correction de bogues du backend CGO sous Linux et de problèmes liés à la zone de notification dans la [PR](https://github.com/wailsapp/wails/pull/5718) par @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Ajouts

- Ajout de HappyTools à la vitrine de la communauté dans la [PR](https://github.com/wailsapp/wails/pull/5061) par @Aliuyanfeng
- Ajout de la prise en charge de la langue indonésienne et d’une documentation complète dans la [PR](https://github.com/wailsapp/wails/pull/5643) par @triadmoko
- Ajout de l’option DisableMenu à WindowsWindow dans la [PR](https://github.com/wailsapp/wails/pull/4813) par @leaanthony

## Modifications

- Mise à jour du modèle Taskfile et de la CLI afin de répartir les tâches de compilation et de création de paquets selon GOOS et ARCH dans la [PR](https://github.com/wailsapp/wails/pull/5617) par @leaanthony

## Corrections

- Correction d’un problème lié au regroupement des fenêtres macOS en onglets dans la [PR](https://github.com/wailsapp/wails/pull/5708) par @taliesin-ai

## Suppressions

- Suppression des fichiers MDX traduits en allemand des sections consacrées aux contributions, aux fonctionnalités et aux guides dans la [PR](https://github.com/wailsapp/wails/pull/5702) par @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Ajouts

- Implémentation du rechargement normal et forcé de WebView sous macOS, et ajout d’un mécanisme de récupération après l’arrêt du processus WebContent dans la [PR](https://github.com/wailsapp/wails/pull/5129) par @wayneforrest
- Ajout d’une documentation allemande complète sur les contributions, les fonctionnalités et les guides dans la [PR](https://github.com/wailsapp/wails/pull/5396) par @leaanthony
- Amélioration des notifications avec le son, les pièces jointes, la planification et l’API de mise à jour dans la [PR](https://github.com/wailsapp/wails/pull/5333) par @popaprozac

## Corrections

- Recalcul du cadre DWM lors d’un changement de DPI pour les fenêtres sans cadre dans la [PR](https://github.com/wailsapp/wails/pull/4785) par @leaanthony
- Correction de la détection des zones de dépôt DnD, qui échouait sous Windows avec une mise à l’échelle différente de 100 %, dans la [PR](https://github.com/wailsapp/wails/pull/4632) par @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Ajouts

- Ajout d’exemples de code à la documentation d’EventsEmit dans la [PR](https://github.com/wailsapp/wails/pull/5026) par @iamhabbeboy
- Ajout d’une option d’hébergement visuel WebView2 sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/5380) par @MerIijn
- Ajout de Klustr à la documentation de la vitrine de la communauté dans la [PR](https://github.com/wailsapp/wails/pull/5536) par @SametKUM
- Ajout de Kira à la vitrine de la communauté, avec de nouvelles pages et une entrée dans le journal des modifications, dans la [PR](https://github.com/wailsapp/wails/pull/5685) par @thiennguyen93
- Ajout d’une section de commentaires au guide du service MCP dans la [PR](https://github.com/wailsapp/wails/pull/5694) par @taliesin-ai

## Modifications

- Le mode serveur dispose désormais d’une compilation de production de premier ordre, cohérente avec les tâches de compilation pour ordinateur de bureau (#5693). `task build:server` compile par défaut un binaire de production (`-tags server,production`, `-trimpath`, symboles supprimés) et accepte `DEV=true` (serveur de développement), `OBFUSCATED=true` (garble) et `EXTRA_TAGS`. `task run:server` exécute un serveur de développement. `Dockerfile.server` / `task build:docker` compilent d’abord le serveur de production (`-tags server,production`) et le frontend de production. Par défaut, l’image utilise une compilation statique en Go pur sur distroless/static, et `CGO_ENABLED`, `GO_IMAGE` et `RUNTIME_IMAGE` sont exposés comme arguments de compilation substituables pour les applications CGO.

## Corrections

- Prévention d’un plantage lors de la fermeture d’une fenêtre contenant des appels asynchrones en attente dans la [PR](https://github.com/wailsapp/wails/pull/4435) par @leaanthony
- Prévention de l’activation d’une fenêtre lors de l’ouverture d’applications masquées sous Windows dans la [PR](https://github.com/wailsapp/wails/pull/5249) par @leaanthony
- Exécution sur le thread principal de GTK du traitement des métadonnées des requêtes WebKit, de l’achèvement des réponses et du flux de corps dans la [PR](https://github.com/wailsapp/wails/pull/5668) par @taliesin-ai
- Correction de `Menu.Update()`, qui ne reconstruisait pas le menu natif sous Linux avec GTK4 (#5659, problème diagnostiqué et corrigé indépendamment par @puneetdixit200 dans #5539)
- Correction d’un plantage lors de l’énumération des écrans macOS après un changement d’affichage, en copiant les chaînes d’identifiant et de nom des écrans et en prenant un instantané de leur nombre (#5565, problème diagnostiqué et corrigé indépendamment par @x-haose dans #5584)
- Correction du contenu WebView2 qui rétrécissait puis disparaissait après le déplacement d’une fenêtre entre des moniteurs aux DPI différents sous Windows, en réappliquant les limites du contrôleur dans le gestionnaire `WM_DPICHANGED`, sur le modèle de la resynchronisation des DPI à la sortie de la réduction (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Ajouts

- Ajout de raccourcis clavier globaux (à l’échelle du système) via `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Les raccourcis se déclenchent même lorsque l’application n’a pas le focus. Ils sont implémentés nativement sur chaque plateforme, sans dépendance tierce : raccourcis Carbon sur macOS, `RegisterHotKey` sous Windows, `XGrabKey` sous X11 et interface de raccourcis globaux du portail de bureau XDG sous Wayland.
- Ajout d’un serveur MCP intégré : un serveur Model Context Protocol qui démarre automatiquement lorsque l’application est compilée avec l’étiquette `mcp`, afin de permettre aux agents LLM de tester et de contrôler une application Wails en cours d’exécution — contrôle des fenêtres, inspection du DOM, évaluation de JavaScript, appels de méthodes liées, événements et simulation des entrées de la souris et du clavier, représentées par un curseur animé à l’écran. Aucun code utilisateur n’est requis : l’étiquette `mcp` est automatiquement ajoutée par `wails3 build`/`wails3 dev` lorsque `WAILS_MCP=1` est défini. La configuration s’effectue entièrement au moyen de variables d’environnement (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Corrigé

- Correction de `Menu.Update()` qui ne reconstruisait pas le menu natif sous Linux avec GTK4 (#5659, problème diagnostiqué et corrigé indépendamment par @puneetdixit200 dans #5539)
- Correction d’un plantage lors de l’énumération des écrans macOS après un changement d’affichage, en copiant les chaînes d’identifiant et de nom des écrans et en prenant un instantané de leur nombre (#5565, problème diagnostiqué et corrigé indépendamment par @x-haose dans #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Ajouté

- Ajout de la documentation expérimentale de Wake, avec navigation dans la barre latérale, dans la [PR](https://github.com/wailsapp/wails/pull/5613) par @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Modifié

- Mise à niveau de `webview2` vers v1.0.27.
  - ci(webview2) : correction de la compilation de publication (compilation croisée pour Windows + go.sum complet) (#5671)\

  **Diff complet :** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Suppression de go vet de la compilation croisée du workflow de publication de webview2 dans la [PR](https://github.com/wailsapp/wails/pull/5672) par @taliesin-ai
- Mise à jour du modèle OpenRouter du journal des modifications automatique vers google/gemini-2.5-flash-lite dans la [PR](https://github.com/wailsapp/wails/pull/5670) par @taliesin-ai
- Mise à niveau de `webview2` vers v1.0.26.

### Corrections

- **Récupération après des erreurs COM transitoires à l’exécution au lieu de quitter** (#5658, #5580). Auparavant, `Chromium.errorCallback` appelait `os.Exit(1)` pour *toute* erreur COM : une défaillance temporaire récupérable après le démarrage arrêtait donc toute l’application. Les chemins d’exécution (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) journalisent désormais l’erreur et poursuivent l’exécution. En particulier, un message web malformé ou non fiable dans `MessageReceived` est désormais ignoré au lieu d’arrêter le processus. Cela corrige la catégorie de plantages liés au passage entre des moniteurs aux résolutions différentes (#5544, #5650). Les erreurs sur les chemins de création de l’environnement ou du contrôleur restent fatales.\

**Diff complet :** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Corrigé

- Correction du workflow release-webview2 pour traiter correctement les fichiers go.sum dans la [PR](https://github.com/wailsapp/wails/pull/5671) par @taliesin-ai
- Correction des mises à jour des menus GTK4 sous Linux en effaçant puis en reconstruisant le menu natif dans la [PR](https://github.com/wailsapp/wails/pull/5659) par @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Ajouté

- Ajout de `application.System` pour détecter la plateforme à l’exécution depuis du code partagé : `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (l’étiquette de compilation `server`) et `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` pour tester directement une cible unique. Il est compilé pour chaque cible, ce qui permet de créer des branches conditionnelles sans étiquettes de compilation. Des fonctions d’assistance correspondantes pour le frontend (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) sont disponibles dans `@wailsio/runtime`
- Ajout d’un guide « Utiliser d’autres frameworks frontend » qui explique comment placer votre propre projet Vite dans `frontend/` (Solid, Preact, Lit, SvelteKit, Qwik, Angular, etc.)
- L’assistant `wails3 setup` vérifie désormais la chaîne d’outils mobiles (iOS/Android) — Xcode et l’environnement d’exécution du simulateur iOS, le JDK, le SDK/NDK Android et l’émulateur — avec, le cas échéant, une installation en un clic et des correctifs copiables pour la configuration du shell
- Les projets générés incluent un fichier `frontend/.npmrc` qui définit un `minimum-release-age` de 7 jours afin de réduire l’exposition aux paquets récemment publiés, potentiellement compromis (pris en charge par pnpm et bun ; ignoré sans effet par npm)

## Modifié

- Refonte de tous les modèles de démarrage intégrés avec une nouvelle présentation principale représentant une montagne au néon (web, iOS et Android)
- **TypeScript est désormais le langage par défaut des modèles de démarrage et utilise le nom de modèle sans suffixe.** `wails3 init` (sans `-t`) génère la structure d’un projet TypeScript ; `-t vanilla`, `-t react`, `-t vue` et `-t svelte` utilisent TypeScript, tandis que leurs variantes JavaScript sont disponibles sous les noms `-t vanilla-js`, `-t react-js`, `-t vue-js` et `-t svelte-js`. Les modèles intégrés déclarent leur langage avec `typescript:` dans `template.yaml` ; les modèles communautaires utilisant le suffixe `-ts` continuent de fonctionner comme solution de repli
- Refonte de l’assistant `wails3 setup` avec le thème « Wails numérique » au néon (effet de verre dépoli dynamique sur un arrière-plan montagneux)

## Corrigé

- Correction d’un plantage sous Windows lors de la restauration d’une application restée réduite assez longtemps pour que WebView2 soit suspendu ou que son processus de rendu/GPU soit recyclé. La resynchronisation du DPI lors de la réduction/restauration (#5544) ne modifie désormais le contrôleur WebView2 que si le DPI de la fenêtre a réellement changé, évitant ainsi les appels COM fatals vers un contrôleur suspendu lors d’une restauration courante sans changement de DPI (#5605)
- Correction de plantages natifs répétés dans `SIGABRT`/`SIGSEGV` (généralement dans `g_object_unref` pendant la boucle principale GTK) dans les applications Linux exécutées longtemps et soumises à de fréquents chargements de ressources ou de médias. Le serveur de ressources terminait les `WebKitURISchemeRequest`s depuis des goroutines de travail, appelant ainsi des fonctions WebKit2GTK non sûres pour les threads en dehors du thread principal de GTK ; la finalisation (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) s’exécute désormais sur le thread principal. Cette modification complète la correction partielle de #5566. Le problème affecte les compilations 6.0 GTK3 et GTK4/WebKitGTK (#5631, #5557)
- Correction d’une erreur intermittente `fatal error: invalid pointer found on stack` dans `setupSignalHandlers` sous Linux/GTK3. Les identifiants de fenêtre transmis comme `user_data` de signal étaient conservés dans une variable locale Go de type `unsafe.Pointer` ; le ramasse-miettes interrompait donc le programme lorsqu’il analysait cette valeur (qui n’est pas un pointeur) pendant une copie de pile. L’identifiant conserve désormais un type entier (`uintptr_t`) côté Go. Cela rétroporte vers l’ancien chemin GTK3 la correction que #4958 avait appliquée au chemin GTK4, où les fonctions de signal C avaient été remplacées par `uintptr_t` afin d’éliminer les erreurs `-race`/checkptr (#5631)

## Supprimé

- Suppression des modèles de démarrage `react-swc`, `preact`, `lit`, `solid`, `qwik` et `sveltekit`, ainsi que de leurs variantes `-ts`. L’ensemble intégré pris en charge se compose désormais de `vanilla`, `react`, `vue` et `svelte` — chacun utilise TypeScript par défaut et propose une variante JavaScript `-js`. Tout autre framework peut toujours être utilisé en [apportant votre propre frontend](https://v3.wails.io/guides/dev/frontend-frameworks) ou au moyen d’un modèle personnalisé

## v3.0.0-alpha2.104 - 2026-06-18

## Corrigé

- Correction d’un plantage iOS (SIGABRT) lorsqu’une méthode liée d’un service Go renvoie une chaîne vide. Le générateur de réponse des ressources iOS vérifiait le pointeur du corps avec `buf != nil` au lieu de vérifier sa longueur ; un corps de longueur nulle provoquait donc une panique de `&buf[0]`. Il vérifie désormais la longueur, comme les générateurs pour ordinateur de bureau

## v3.0.0-alpha2.103 - 2026-06-15

## Modifié

- Déplacement des fonctionnalités natives d’iOS et d’Android vers des gestionnaires de plateforme : appelez-les via `application.IOS.*` et `application.Android.*` (par exemple `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`) au lieu des anciennes fonctions libres `application.IOS*`/`application.Android*` (#5602)
- Renommer les événements du pont mobile : les événements multiplateformes utilisent désormais le préfixe `common:*` (par exemple `common:haptic` et `common:location`), tandis que les événements propres à une plateforme utilisent `ios:*` / `android:*` (par exemple `ios:backgroundTask` et `android:foregroundService`) ; le préfixe `native:*` n’est plus utilisé (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Ajouts

- Ajouter l’assistant expérimental `wails3 setup` pour configurer les projets de manière interactive et vérifier les dépendances
- Ajouter l’option `--json` à `wails3 doctor` pour produire une sortie lisible par une machine
- Ajouter une section sur l’état de la signature à la commande `wails3 doctor`

## Corrections

- Corriger la détection de npm sous Linux afin qu’elle vérifie le PATH en plus du gestionnaire de paquets

## v3.0.0-alpha.101 - 2026-06-13

## Ajouts

- iOS : boîtes de dialogue natives pour les messages (UIAlertController) et pour l’ouverture d’un fichier, de plusieurs fichiers ou d’un répertoire (UIDocumentPickerViewController) ; les boîtes de dialogue d’enregistrement renvoient une erreur explicite
- iOS : prise en charge du presse-papiers via UIPasteboard
- iOS : métriques d’écran réelles via UIScreen (points, pixels, échelle et zone de travail tenant compte de la zone de sécurité)
- iOS : builds pour appareils (`IOS_PLATFORM=device`), prise en charge de l’identité de signature du code, du profil de provisionnement et des autorisations, empaquetage `.ipa` et `deploy-device` via devicectl
- iOS : version minimale d’iOS configurable (`ios.minIOSVersion` dans build/config.yml)
- iOS : `wails3 doctor` indique la disponibilité de Xcode et du SDK iOS sous macOS
- iOS : les événements système concernant la batterie, le réseau, le thème, le verrouillage de l’écran et le manque de mémoire sont exposés sous forme d’événements d’application `events.IOS.*` et `events.Common.*` indépendants de la plateforme
- iOS : pont natif vers les fonctionnalités mobiles (`application.IOS*` exporté) — feuille de partage, ouverture d’URL, maintien de l’appareil éveillé, lampe torche, marges de la zone de sécurité, luminosité, informations sur l’application, verrouillage de l’orientation, barre d’état, biométrie (Face ID/Touch ID), notifications locales et stockage sécurisé dans le Keychain
- iOS : capteurs et matériel — retour haptique, géolocalisation ponctuelle, accéléromètre, proximité, synthèse vocale, informations sur le stockage, état de l’alimentation et de la batterie, état du réseau, marges du clavier et détection des captures d’écran
- iOS : documentation (IOS.md et guide sur le site de documentation)
- Android : boîtes de dialogue natives pour les messages (AlertDialog) et pour l’ouverture d’un ou de plusieurs fichiers (Storage Access Framework, avec importation sous forme de copies dans le cache) ; les boîtes de dialogue d’ouverture de répertoire et d’enregistrement renvoient une erreur explicite
- Android : prise en charge du presse-papiers via ClipboardManager
- Android : métriques d’écran réelles via WindowMetrics/DisplayMetrics (dp, pixels, échelle et zone de travail tenant compte des barres système)
- Android : méthodes d’exécution pour le retour haptique (`Android.Haptics.Vibrate`), les informations sur l’appareil (`Android.Device.Info`) et les notifications toast (`Android.Toast.Show`)
- Android : événements de cycle de vie typés (`events.Android.*`, générés à partir de events.txt), avec `ActivityCreated` associé à `Common.ApplicationStarted`
- Android : la chaîne de build produit des APK de débogage et de publication installables (`android:run`, `android:package`, `android:package:fat`) ; la signature des publications utilise par défaut le magasin de clés de débogage, ou un véritable magasin de clés au moyen des variables d’environnement `ANDROID_KEYSTORE_*`
- Android : `wails3 doctor` indique les SDK Android, NDK et JDK
- Android : les événements système concernant la batterie, le réseau, le thème, le verrouillage de l’écran et le manque de mémoire sont exposés sous forme d’événements d’application `events.Android.*` et `events.Common.*` indépendants de la plateforme
- Android : pont natif vers les fonctionnalités mobiles (`application.Android*` exporté) — partage, ouverture d’URL, maintien de l’appareil éveillé, lampe torche, marges de la zone de sécurité, luminosité, informations sur l’application, verrouillage de l’orientation, barre d’état, biométrie (BiometricPrompt), notifications locales et stockage sécurisé dans EncryptedSharedPreferences
- Android : capteurs et matériel — retour haptique, géolocalisation ponctuelle, accéléromètre, proximité, synthèse vocale, informations sur le stockage, état de l’alimentation et de la batterie, état du réseau, marges du clavier et blocage des captures d’écran avec FLAG_SECURE
- Android : documentation (ANDROID.md et guide sur le site de documentation)
- Exemple : l’exemple complet `mobile` comporte désormais des onglets Mobile et Hardware qui présentent le pont vers les fonctionnalités natives sous iOS et Android (les onglets en forme de pastille passent sur plusieurs lignes)
- Mobile : batterie — l’accéléromètre, le capteur de proximité, la lampe torche et l’horloge périodique de l’exemple sont suspendus lorsque l’application passe en arrière-plan, puis rétablis à son retour (Android laisse le processus s’exécuter en arrière-plan et, sous iOS, la lampe torche constitue un état matériel persistant) ; les récepteurs d’événements système Android ne sont enregistrés que lorsque l’application est au premier plan
- iOS : capture avec la caméra — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → événement `native:capture` contenant une miniature encodée en base64)
- iOS : exécution en arrière-plan — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (une fenêtre de tâche en arrière-plan UIApplication) et paramètre `ios.backgroundModes` configurable dans build/config.yml, qui insère `UIBackgroundModes` dans le modèle du fichier Info.plist généré
- Android : capture avec la caméra — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (caméra système via FileProvider → événement `native:capture`)
- Android : service au premier plan — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (un `WailsForegroundService` accompagné d’une notification permanente maintient le processus actif pour les tâches longues en arrière-plan)
- Exemple : un onglet Camera présentant la capture de photos et de vidéos ainsi que l’exécution en arrière-plan (service au premier plan sous Android, fenêtre de tâche en arrière-plan sous iOS)

## Corrections

- Corriger l’échec systématique de `getUserMedia` avec `NotAllowedError` sous Linux : WebKitGTK refuse les demandes d’autorisation qu’aucun gestionnaire ne traite, et le signal `permission-request` n’était pas connecté. La caméra et le microphone sont désormais gérés au moyen d’une nouvelle table multiplateforme `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), respectée sous Linux (WebKitGTK) comme sous Windows (WebView2). Sous Linux, qui ne fournit aucune invite native, la caméra et le microphone sont autorisés par défaut (ce qui rétablit `getUserMedia`) et peuvent être désactivés avec `PermissionDeny` (#5552)
- iOS : `GOOS=ios` compile de nouveau (`events.IOS` exporté et stubs de noms de méthodes mobiles), tout comme les builds portant la balise de production (correction des balises de build dans pkg/application et plusieurs services)
- iOS : les événements Go→JS et ExecJS fonctionnent désormais — la page ne se charge plus deux fois au démarrage et la négociation `wails:runtime:ready` ne peut plus être perdue
- iOS : `ApplicationDidFinishLaunching`/`ApplicationStarted` ne sont plus en concurrence avec le démarrage de l’application ; suppression de l’attente fixe de 2 secondes au démarrage
- iOS : correction d’une fuite de chaîne C à chaque exécution de JavaScript de Go vers JS
- iOS : `hasListeners` reflète désormais l’enregistrement réel des écouteurs
- iOS : la journalisation de débogage du framework est exclue des builds de production lors de la compilation
- Android : `GOOS=android` compile de nouveau — définition de `events.Android`, suppression du tableau d’écouteurs `events_android.go` dont les accès dépassaient les limites, ajout du stub de nom de méthode mobile et exclusion des fichiers Linux pour ordinateur de bureau (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) des builds Android
- Android : les liaisons JS→Go fonctionnent désormais — la WebView ne pouvant pas transmettre les corps des requêtes POST `fetch()` à `shouldInterceptRequest`, les appels d’exécution passent par un transport JavascriptInterface (`nativeHandleRuntimeCall`) au lieu de provoquer un plantage à cause d’un corps de requête nil
- Android : les appels d’exécution `Screens.*` renvoient des données réelles — le ScreenManager est désormais rempli au démarrage (il n’avait jamais été connecté, de sorte que `GetAll` renvoyait nil)
- Android : la journalisation de débogage du framework est exclue des builds de production lors de la compilation et, dans les builds de débogage, elle est acheminée vers logcat sous l’étiquette `Wails`
- Android : véritable registre `hasListeners`, gestion des références et des exceptions JNI, et cycle de vie de page à chargement unique (sans double navigation)
- Corriger l’échec de `wails3 generate bindings` avec « Access is denied » sous Windows lorsque le serveur de développement Vite est en cours d’exécution, en synchronisant les fichiers générés dans le répertoire de sortie au lieu de les renommer par-dessus celui-ci (#5515)
- Correction d’un plantage fatal intermittent sur macOS lors de la lecture des informations d’écran après un changement d’affichage : l’identifiant et le nom de l’écran stockaient des pointeurs vers des tampons `UTF8String` libérés automatiquement, qui pouvaient l’être avant leur copie par Go (utilisation après libération). Les chaînes sont désormais copiées avec `strdup`, puis libérées après conversion, et l’énumération des écrans s’exécute dans un pool de libération automatique explicite afin de ne plus provoquer de fuite lorsqu’elle est appelée depuis des goroutines Go (#5556)
- Correction d’une erreur SIGSEGV intermittente sous Linux lorsque l’assetserver ferme un `WebKitURISchemeRequest` : le dernier `g_object_unref` s’exécutait dans la goroutine de l’assetserver et finalisait un GObject WebKit en dehors du thread principal de GTK. L’opération unref est désormais transférée vers le contexte principal de GTK au moyen de `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Ajouts

- Extension de `MacWebviewPreferences` avec des options de configuration WKWebView supplémentaires : `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` et `ApplicationNameForUserAgent` (#5549)

## Corrections

- Correction de l’échec de `wails3 generate bindings` avec le message « Access is denied » sous Windows lorsque le serveur de développement Vite est en cours d’exécution, en synchronisant les fichiers générés dans le répertoire de sortie au lieu de les renommer par-dessus celui-ci (#5561)
- Correction des événements de redimensionnement JS qui ne se déclenchaient pas pour les fenêtres sans cadre sous Linux, ainsi que de la détection des bords de la barre de défilement pour ces fenêtres (#5368)
- Correction de l’échec du programme de mise à jour sous Windows avec le message « invalid cross-device link » lorsque le répertoire temporaire se trouve sur un volume différent de celui du répertoire d’installation (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Corrections

- Correction de l’échec de `wails3 generate bindings` avec le message « Access is denied » sous Windows lorsque le serveur de développement Vite est en cours d’exécution, en synchronisant les fichiers générés dans le répertoire de sortie au lieu de les renommer par-dessus celui-ci (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Corrections

- Correction du blocage de l’interface WebKit sous Linux lorsqu’elle est inactive (par exemple, lorsque l’inspecteur est ouvert), en cessant de forcer `SA_ONSTACK` sur `SIGUSR1`, ce qui perturbait la synchronisation du thread du ramasse-miettes de JavaScriptCore (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Ajouts

- Ajout d’une page de débogage et de l’utilisation de `runtime/trace`

## Modifications

- Suppression de certaines importations `_ "embed"` inutiles afin de nettoyer légèrement le code

## Corrections

- Correction des contraintes de largeur et de hauteur minimales qui n’étaient plus appliquées après la restauration d’une fenêtre agrandie sous Windows (#4593)
- Correction de la transmission des clics de souris à travers une fenêtre en mode plein écran avec les options Frameless + Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Ajouts

- Ajout de la prise en charge de l’obfuscation Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)) : identifiants stables des méthodes de liaison, intégration à la compilation et au Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), ainsi que des balises de structure JSON sur chaque charge utile exposée au runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), afin que le format d’échange résiste au renommage par Garble des champs exportés.

## v3.0.0-alpha.95 - 2026-05-20

## Ajouts

- Ajout de la page manquante sur la structure du projet

## Modifications

- Documentation : remplacement de quelques diagrammes de la page sur l’architecture par des diagrammes de séquence pour un affichage plus clair
- Documentation : ajout d’une note indiquant que l’installation de D2 est requise avant l’exécution

## Corrections

- Correction de `wails3 generate appimage` avec la configuration GTK4 par défaut : le bundler détecte désormais la pile GTK à partir du binaire avant de rechercher les fichiers d’exécution. Il sélectionne ainsi `libwebkitgtkinjectedbundle.so` (sous `webkitgtk-6.0/`) pour les compilations GTK4 et `libwebkit2gtkinjectedbundle.so` (sous `webkit2gtk-4.1/`) pour les compilations `-tags gtk3`. La sonde `.relr.dyn` vérifie également `libgtk-4.so.1`, de sorte que la suppression des symboles est correctement désactivée avec les chaînes d’outils modernes, quelle que soit la pile. (#5475)
- Correction de l’échec de `wails3 generate appimage` lorsqu’il est appelé avec un `-builddir` relatif : le bundler résout désormais dès le départ `-binary`, `-icon`, `-desktopfile`, `-builddir` et `-outputdir` en chemins absolus, afin que le `s.CD` exécuté en cours de processus ne perturbe ni la goroutine de téléchargement d’AppRun ni la sonde `ldd` exécutée après la copie.
- Correction de l’échec de `wails3 generate appimage` lors du déplacement de l’AppImage finale vers `-outputdir` lorsque le champ `Name=` du fichier desktop ne correspond pas au nom de base du binaire : le bundler force désormais le plugin appimage de linuxdeploy, au moyen de la variable d’environnement `OUTPUT`, à écrire l’AppImage dans `<binary>-<arch>.AppImage` plutôt que sous le nom dérivé du fichier desktop.
- Correction de `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` et `Common.SystemDidWake` qui ne se déclenchaient plus sous Linux après que la pile GTK4 + WebKitGTK 6.0 est devenue la configuration par défaut dans alpha.93. Le nouveau `application_linux.go` `run()` par défaut n’appelait ni `setupCommonEvents()`, qui transmet les événements `Linux.*` à leurs équivalents `Common.*`, ni `monitorPowerEvents()`. L’utilitaire DBus de surveillance de l’alimentation est désormais partagé entre les chemins de compilation GTK3 et GTK4 au moyen de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Corrections

- Correction de `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` et `Common.SystemDidWake` qui ne se déclenchaient plus sous Linux après que la pile GTK4 + WebKitGTK 6.0 est devenue la configuration par défaut dans alpha.93. Le nouveau `application_linux.go` `run()` par défaut n’appelait ni `setupCommonEvents()`, qui transmet les événements `Linux.*` à leurs équivalents `Common.*`, ni `monitorPowerEvents()`. L’utilitaire DBus de surveillance de l’alimentation est désormais partagé entre les chemins de compilation GTK3 et GTK4 au moyen de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Ajouts

- Ajout de `XDG_SESSION_TYPE` à la sortie de `wails3 doctor` sous Linux par @leaanthony

## Corrections

- Correction du plantage du menu de fenêtre sous Wayland, causé par l’accès d’appmenu-gtk-module à une fenêtre non réalisée (#4769), par @leaanthony
- Correction du plantage de l’application GTK lorsque son nom contient des caractères non valides (espaces, parenthèses, etc.), par @leaanthony
- Correction de l’erreur « not enough memory » lors de l’initialisation du glisser-déposer sous Windows (#4701), par @overlordtm
- Correction d’une condition de concurrence dans le stockage des rappels du thread principal, due à l’utilisation incorrecte de RLock pour supprimer un élément de la map (Linux, macOS, iOS) (#4424), par @leaanthony
- Correction de la gestion des variables lors de la transmission d’arguments de ligne de commande aux tâches. Les variables de la CLI indiquées sous forme de paires KEY=VALUE sont désormais correctement initialisées et propagées pendant toute l’exécution des tâches.
- Correction du conflit NSWindowZoomButton sous macOS : `MaximiseButtonState` et `FullscreenButtonState` appliquent désormais l’état le plus restrictif au démarrage comme pendant l’exécution ; aucun des deux setters ne peut plus remplacer silencieusement l’autre (#5319)
- Correction d’un ensemble de bogues préexistants dans l’ancien chemin de compilation GTK3 (`-tags gtk3`), détectés par CodeRabbit dans la demande #5463 : les lancements par association de fichiers n’ignorent plus les gestionnaires de démarrage ; `getTheme` vérifie désormais les limites et les types ; `appName` ne libère plus de mémoire appartenant à GLib ; `clipboardGet` ne laisse plus fuir le `gchar*` renvoyé par GTK ; `Calloc` utilise désormais des receveurs de type pointeur et `NewCalloc` renvoie `*Calloc`, afin que le pool suive réellement les allocations ; `zoomOut` utilise l’inverse de `zoomInFactor` au lieu d’un multiplicateur négatif qui ramenait la valeur à 1.0 ; `execJS` réutilise le nom de monde vide préalloué au lieu de laisser fuir un `C.CString("")` à chaque appel ; un `fmt.Println` de développement a été supprimé de `menuItem.setAccelerator`. Résout le problème #5465.
- Corrige la même fuite liée au récepteur par valeur de `Calloc` dans le chemin de compilation GTK4 par défaut (`linux_cgo.go`) : des récepteurs par pointeur et `NewCalloc() *Calloc` permettent désormais de suivre et de libérer réellement les allocations `c.String(...)` propres à chaque fenêtre.

## v3.0.0-alpha.92 - 2026-05-15

## Ajouté

- Modifie les Taskfiles afin de permettre le choix du gestionnaire de paquets frontend utilisé au moyen de l’option `PACKAGE_MANAGER`
- Enrichit les données des modèles avec `{{.Opn}}` et `{{.Cls}}` afin de rendre l’écriture des modèles de Taskfile plus prévisible

## Modifié

- Modifie certains Taskfiles existants afin qu’ils utilisent `{{.Opn}} and {{.Cls}}`

## Corrigé

- Corrige une erreur fatale d’exécution `concurrent map read and map write` dans `linuxSystemTray` lorsque le menu de la zone de notification est mis à jour pendant que le panneau le lit.
- Utilise `log` au lieu de `fmt` pour afficher les erreurs et les traces de pile de WebView2, afin que les messages ne soient pas perdus lorsque l’application s’exécute sous Windows sans console attachée.

## v3.0.0-alpha.91 - 2026-05-12

## Modifié

- Met à jour le SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5414) par `@github-actions[bot]`
- **RUPTURE DE COMPATIBILITÉ (macOS) :** normalise le système de coordonnées de macOS afin que `GetScreens`, `Position` et `SetPosition` utilisent tous le même espace : des points logiques, avec l’axe Y orienté vers le bas et `(0,0)` dans l’angle supérieur gauche de l’écran principal. Ce comportement correspond à celui de Windows, de GTK ainsi que des API publiques d’Electron et du Web. Les écrans physiquement situés au-dessus de l’écran principal signalent désormais une valeur `Bounds.Y` négative (auparavant positive), et les valeurs `Position()`/`SetPosition()` sont maintenant exprimées en points logiques plutôt qu’en `points × primaryScale`. La conversion aller-retour `Position()` → `SetPosition()` est préservée ; les valeurs absolues consignées par les versions alpha antérieures ou les solutions de contournement calculées manuellement (par exemple, multiplication par `primaryScale` ou inversion de Y par rapport à la hauteur d’un écran) devront être mises à jour. Résout [#5117](https://github.com/wailsapp/wails/issues/5117).

## Corrigé

- Valide de manière défensive le nom et la longueur du corps des signaux DBus afin d’éviter les paniques, dans la [PR](https://github.com/wailsapp/wails/pull/5416) par @leaanthony
- Corrige un problème de sûreté de la mémoire dans la gestion des menus GTK sous Linux, dans la [PR](https://github.com/wailsapp/wails/pull/5363) par @leaanthony
- Détecte les GPU NVIDIA et désactive le moteur de rendu DMA-BUF sous Linux, dans la [PR](https://github.com/wailsapp/wails/pull/5295) par @leaanthony
- Corrige la conversion Y entre écrans de `SetPosition` sous macOS : utilise la hauteur de l’écran principal comme référence globale afin que les fenêtres soient placées correctement sur les moniteurs décalés verticalement par rapport à l’écran principal, dans [#5117](https://github.com/wailsapp/wails/issues/5117)
- Corrige le modèle de PR Git afin qu’il pointe vers la bonne URL de retour d’expérience, dans la [PR](https://github.com/wailsapp/wails/pull/5109) par @wayneforrest
- Corrige une série de plantages `SetMenu` de la zone de notification Windows provoqués par un appel système `DestroyMenu` défectueux qui transmettait quatre arguments au lieu d’un : chaque appel renvoyait donc FALSE et ne libérait rien. Libère également les handles HMENU et HBITMAP, y compris ceux alloués à l’exécution par `MenuItem.SetBitmap`, lors de la reconstruction des menus, réinitialise les tables obsolètes de cases à cocher et de boutons radio dans `Win32Menu.Update`, et supprime dans `systemtray.updateMenu` un appel `Update()` redondant qui doublait les allocations. Les applications utilisant durablement la zone de notification ne génèrent désormais plus de fuite d’objets GDI/USER à chaque reconstruction du menu.

## v3.0.0-alpha.90 - 2026-05-11

## Ajouté

- Ajoute un nom d’application configurable à l’agent utilisateur de WKWebView sous macOS, dans la [PR](https://github.com/wailsapp/wails/pull/5261) par @vinhvoit225
- Ajoute la dépendance indirecte github.com/coder/websocket à l’exemple gin-service, dans la [PR](https://github.com/wailsapp/wails/pull/5400) par @taliesin-ai
- Ajoute la prise en charge des comparaisons par égalité profonde aux tests des ressources de compilation, dans la [PR](https://github.com/wailsapp/wails/pull/5402) par @leaanthony

## Modifié

- Regroupe la sortie de compilation dans le répertoire des ressources, dans la [PR](https://github.com/wailsapp/wails/pull/5401) par @taliesin-ai
- Met à jour le SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5399) par `@github-actions[bot]`

## Corrigé

- Utilise un objet de notification pour le message d’instance unique sous macOS, dans la [PR](https://github.com/wailsapp/wails/pull/5289) par @overlordtm
- Regroupe les rappels Windows par lots afin d’éviter la perte de promesses sous forte charge, dans la [PR](https://github.com/wailsapp/wails/pull/5383) par @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Ajouté

- Ajoute la tâche go<em>test</em>results pour agréger les résultats des tests Go, dans la [PR](https://github.com/wailsapp/wails/pull/5316) par @leaanthony

## Modifié

- Répartit sous certaines conditions les charges utiles RPC volumineuses entre plusieurs requêtes POST fragmentées, dans la [PR](https://github.com/wailsapp/wails/pull/5369) par @leaanthony
- Met à niveau Vite de 5.x.x vers 8.0.0 dans tous les modèles frontend, dans la [PR](https://github.com/wailsapp/wails/pull/5386) par @leaanthony
- Migre la configuration du port du serveur de développement Vite vers des variables d’environnement, dans la [PR](https://github.com/wailsapp/wails/pull/5365) par @leaanthony
- Configure le serveur de développement Vite pour qu’il se lie à 127.0.0.1 dans tous les modèles, dans la [PR](https://github.com/wailsapp/wails/pull/5361) par @leaanthony
- Met à jour le SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5384) par `@github-actions[bot]`

## Corrigé

- Assainit les ébauches du modèle Info.plist lors de la mise à jour des ressources de compilation, dans la [PR](https://github.com/wailsapp/wails/pull/5312) par @leaanthony
- Corrige l’état obsolète des menus macOS en appliquant les fonctions de mutation des éléments de menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) de manière synchrone sur le thread principal, ce qui élimine la condition de concurrence `dispatch_async` qui amenait les menus à afficher l’état précédent lorsqu’ils étaient rouverts rapidement (#5002)
- Ignore les fichiers `*_test.go` en mode développement afin d’éviter les recompilations inutiles, dans la [PR](https://github.com/wailsapp/wails/pull/5203) par @leaanthony
- Empêche une erreur de segmentation de Menu.Update() lorsque l’application n’est pas en cours d’exécution, dans la [PR](https://github.com/wailsapp/wails/pull/5291) par @wucm667
- Utilise lastSizeWParam pour conditionner le redessin de la barre de menus sous Windows, dans la [PR](https://github.com/wailsapp/wails/pull/5382) par @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Modifié

- Modifie HiddenOnTaskbar afin d’utiliser WS<em>EX</em>TOOLWINDOW, dans la [PR](https://github.com/wailsapp/wails/pull/5371) par @leaanthony
- Réorganise les dépendances et supprime la directive replace de webview2 dans go.mod, dans la [PR](https://github.com/wailsapp/wails/pull/5370) par @atterpac
- Met à jour le SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5358) par `@github-actions[bot]`

## Corrigé

- Supprime les alias d’indirection génériques et regroupe les types de clés de tables, dans la [PR](https://github.com/wailsapp/wails/pull/5331) par @fbbdev

## Supprimé

- Supprime le workflow PR-master, ainsi que la documentation, les tests Go et les tests ignorés, dans la [PR](https://github.com/wailsapp/wails/pull/5377) par @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Ajouté

- Ajoute la documentation en coréen pour Wails v3, dans la [PR](https://github.com/wailsapp/wails/pull/5352) par @leaanthony
- Ajoute la documentation en français sur l’installation et le démarrage rapide, dans la [PR](https://github.com/wailsapp/wails/pull/5354) par @leaanthony
- Ajoute la documentation en portugais sur le démarrage rapide, les concepts et la communauté, dans la [PR](https://github.com/wailsapp/wails/pull/5355) par @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Ajouts

- Ajout de la traduction française de la documentation dans la [PR](https://github.com/wailsapp/wails/pull/5328) par @leaanthony
- Ajout de la langue allemande au site de documentation dans la [PR](https://github.com/wailsapp/wails/pull/5343) par @leaanthony

## Modifications

- Enregistrement des 8 langues traduites dans la configuration de la documentation dans la [PR](https://github.com/wailsapp/wails/pull/5347) par @leaanthony
- Mise à jour de divers fichiers liés à Windows pour WebView2 dans la [PR](https://github.com/wailsapp/wails/pull/5317) par @leaanthony

## Corrections

- Séparation de la répartition des boîtes de dialogue entre GTK3 et GTK4 sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/5340) par @leaanthony
- Exécution garantie des fonctions de rappel des boîtes de dialogue sur le thread GTK, ce qui corrige des erreurs de segmentation, dans la [PR](https://github.com/wailsapp/wails/pull/5339) par @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Ajouts

- Ajout de l’URL du modèle de PR au dépôt dans la [PR](https://github.com/wailsapp/wails/pull/5179) par @leaanthony
- Ajout de la documentation allemande pour Wails v3 dans la [PR](https://github.com/wailsapp/wails/pull/5330) par @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Ajouts

- Ajout d’une option permettant d’empêcher la touche Échap de quitter le mode plein écran sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5307) par @leaanthony
- Ajout d’une option permettant d’empêcher la touche Échap de quitter le mode plein écran sous macOS dans la [PR](https://github.com/wailsapp/wails/pull/5310) par @leaanthony
- Ajout de Pausa à la vitrine de la communauté dans la documentation, dans la [PR](https://github.com/wailsapp/wails/pull/5288) par @yuseferi

## Modifications

- Mise à jour du SVG des sponsors dans la [PR](https://github.com/wailsapp/wails/pull/5308) par `@github-actions[bot]`
- Mise à jour de la commande de génération d’icônes pour gérer les plateformes non prises en charge dans la [PR](https://github.com/wailsapp/wails/pull/5309) par @leaanthony
- Remplacement de l’API booléenne du plein écran par le type ternaire ButtonState et implémentation des liaisons propres aux plateformes dans la [PR](https://github.com/wailsapp/wails/pull/5224) par @leaanthony

## Corrections

- Protection des opérations de focus de WebView2 contre un état de contrôleur nil dans la [PR](https://github.com/wailsapp/wails/pull/5315) par @leaanthony
- Mise à jour du workflow GitHub Actions afin qu’il référence correctement la branche de base de la PR dans la [PR](https://github.com/wailsapp/wails/pull/5313) par @leaanthony
- Fichiers `*_test.go` ignorés en mode développement pour éviter les reconstructions inutiles dans la [PR](https://github.com/wailsapp/wails/pull/5203) par @leaanthony
- Prévention d’une erreur de segmentation de Menu.Update() lorsque l’application n’est pas en cours d’exécution dans la [PR](https://github.com/wailsapp/wails/pull/5291) par @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Ajouts

- Ajout de l’indicateur InstallScope et d’une option de compilation pour l’installation par machine ou par utilisateur dans la [PR](https://github.com/wailsapp/wails/pull/5094) par @symball
- Ajout à BrowserWindow d’une méthode SetScreen sans effet afin de satisfaire l’interface Window dans la [PR](https://github.com/wailsapp/wails/pull/5294) par @leaanthony

## Corrections

- Détection des GPU NVIDIA et désactivation du moteur de rendu DMA-BUF sous Linux dans la [PR](https://github.com/wailsapp/wails/pull/5295) par @leaanthony
- Correction du modèle de PR Git afin qu’il pointe vers la bonne URL de commentaires dans la [PR](https://github.com/wailsapp/wails/pull/5109) par @wayneforrest
- Correction d’une série de plantages de la zone de notification Windows `SetMenu` causés par un appel système `DestroyMenu` défectueux qui transmettait quatre arguments au lieu d’un, si bien que chaque appel renvoyait FALSE sans rien libérer. Libération également des handles HMENU et HBITMAP, y compris ceux alloués à l’exécution via `MenuItem.SetBitmap`, lors de la reconstruction des menus, réinitialisation des tables obsolètes de cases à cocher et de boutons radio dans `Win32Menu.Update`, et suppression dans `systemtray.updateMenu` d’un appel `Update()` redondant qui doublait les allocations. Les applications utilisant longtemps la zone de notification ne provoquent plus de fuite d’objets GDI/USER à chaque reconstruction du menu.

## v3.0.0-alpha.82 - 2026-05-01

## Corrections

- Correction de la génération du fichier desktop afin de gérer correctement le nom du bureau dans la [PR](https://github.com/wailsapp/wails/pull/5232) par @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Modifications

- Modification de la planification des versions nocturnes à 15:00 UTC dans la [PR](https://github.com/wailsapp/wails/pull/5286) par @leaanthony

## Corrections

- Correction des valeurs Screen Bounds, WorkArea et Size divisées par deux sur les Mac Retina -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Modifications

- Mise à jour des dépendances de la documentation et des chargeurs de collections de contenu dans la [PR](https://github.com/wailsapp/wails/pull/5285) par @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Ajouts

- Octroi de l’autorisation actions: write à la tâche trigger-release dans la [PR](https://github.com/wailsapp/wails/pull/5270) par @leaanthony

## Modifications

- La tâche de publication utilise désormais la branche master par défaut et met à jour la formulation du journal des modifications dans la [PR](https://github.com/wailsapp/wails/pull/5283) par @leaanthony
- Mise à jour du workflow auto-changelog afin d’utiliser la dernière version dans la [PR](https://github.com/wailsapp/wails/pull/5282) par @leaanthony
- Amélioration de l’efficacité des workflows grâce à l’ajout de filtres de chemins et à la suppression des workflows inutilisés dans la [PR](https://github.com/wailsapp/wails/pull/5280) par @leaanthony
- Mise à jour de la documentation afin que les liens vers les exemples référencent la branche master dans la [PR](https://github.com/wailsapp/wails/pull/5274) par @leaanthony
- Mise à jour de la documentation et des exemples pour v3 dans la [PR](https://github.com/wailsapp/wails/pull/5272) par @leaanthony

## Corrections

- Amélioration du proxy inverse avec une logique de nouvelle tentative et l’utilisation forcée d’IPv4 pour le développement dans la [PR](https://github.com/wailsapp/wails/pull/5265) par @AkagiYui
- Réécriture du workflow de déclenchement du journal des modifications non publiées dans la [PR](https://github.com/wailsapp/wails/pull/5281) par @leaanthony

## Suppressions

- Suppression des scripts de test shell destinés à divers usages de test dans la [PR](https://github.com/wailsapp/wails/pull/5267) par @leaanthony
- Suppression du workflow de déploiement de la documentation v3-alpha et de l’enregistrement CNAME dans la [PR](https://github.com/wailsapp/wails/pull/5266) par @leaanthony

### Ajouts

- Ajout de l’entrée Routage frontend à la barre de navigation latérale dans la [PR](https://github.com/wailsapp/wails/pull/5196) par @leaanthony
- Ajout d’un guide sur le routage frontend avec des recommandations propres à chaque framework dans la [PR](https://github.com/wailsapp/wails/pull/5185) par @leaanthony
- Ajout de la prise en charge des feuilles modales (macOS)
- Mise à niveau de la version de ghw pour améliorer la prise en charge des appareils Apple par @leaanthony (#4977)
- Ajout de la méthode `GetBadge` au service du Dock
- Ajout de l’indicateur `-tags` à la commande `wails3 build` pour transmettre des balises de compilation Go personnalisées (par exemple, `wails3 build -tags gtk4`) (#4957)
- Ajout d’une documentation sur la génération automatique des énumérations dans le générateur de liaisons, comprenant une page Énumérations dédiée et une entrée dans la barre de navigation latérale (#4972)
- Ajout de l’indicateur `-tags` à la commande `wails3 build` pour transmettre des balises de compilation Go personnalisées (par exemple, `wails3 build -tags gtk4`) (#4957)
- Ajout d’exemples d’API Web dans `v3/examples/web-apis/` présentant plus de 41 API de navigateur, notamment pour le stockage (localStorage, sessionStorage, IndexedDB, Cache API), le réseau (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), les médias (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), les appareils (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), les performances (Performance API, Mutation Observer, Intersection/Resize Observer), l’interface utilisateur (Web Components, Pointer Events, Selection, Dialog, glisser-déposer), et bien plus encore
- Ajout d’un exemple de vérificateur de compatibilité des API WebView (`v3/examples/webview-api-check/`), qui teste plus de 200 API de navigateur sur différentes plateformes
- Ajout du paquet `internal/libpath` pour rechercher les chemins des bibliothèques natives sous Linux, avec recherche parallèle, mise en cache et prise en charge de Flatpak/Snap/Nix
- **Travail en cours :** ajout de la prise en charge expérimentale de WebKitGTK 6.0 / GTK4 sous Linux, disponible via `-tags gtk4` (GTK3/WebKit2GTK 4.1 reste utilisé par défaut)
- Remarque : avec les gestionnaires de fenêtres en mosaïque (par exemple Hyprland ou Sway), les opérations de réduction et d’agrandissement peuvent ne pas fonctionner comme prévu, car le gestionnaire de fenêtres contrôle la géométrie des fenêtres
- Ajout aux docs de **l’écoute des événements en JavaScript** d’une explication par @AbdelhadiSeddar sur la création de **gestionnaires à usage unique**
- Ajout par @leaanthony de l’option `UseApplicationMenu` à `WebviewWindowOptions`, permettant aux fenêtres sous Windows/Linux d’hériter du menu d’application défini via `app.Menu.Set()`
- Ajout par @wimaha de la prise en charge des fichiers `.icon` (format Apple Icon Composer) pour générer des icônes Liquid Glass et des catalogues de ressources (macOS) (#4934)
- Ajout d’un mode serveur expérimental pour les déploiements sans interface graphique ou sur le Web (`-tags server`). Il permet d’exécuter des applications Wails comme serveurs HTTP sans dépendances d’interface graphique native. Compilez avec `wails3 task build:server`. Consultez `examples/server` pour plus de détails.
- Ajout du paquet `internal/libpath` pour rechercher les chemins des bibliothèques natives sous Linux, avec recherche parallèle, mise en cache et prise en charge de Flatpak/Snap/Nix
- Ajout par @leaanthony de l’option `CollectionBehavior` à `MacWindow` pour contrôler le comportement des fenêtres entre les Spaces macOS et le mode plein écran (#4756)
- Ajout par @leaanthony de tests unitaires pour pkg/application
- Ajout par @leaanthony de la prise en charge des protocoles personnalisés dans la création de paquets MSIX
- Ajout de la détection de l’environnement de bureau sous Linux [PR nº 4797](https://github.com/wailsapp/wails/pull/4797)
- Ajout par @leaanthony de la méthode `Window.Print()` à l’environnement d’exécution JavaScript pour ouvrir la boîte de dialogue d’impression depuis le frontend (#4290)
- Ajout par @leaanthony de `XDG_SESSION_TYPE` à la sortie de `wails3 doctor` sous Linux
- Ajout par @leaanthony d’événements WebKit2 supplémentaires signalant les changements de chargement sous Linux : `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896)
- Ajout par @leaanthony de `XDG_SESSION_TYPE` à la sortie de `wails3 doctor` sous Linux
- Génération du fichier `.desktop` pendant la compilation Linux, et non plus seulement pendant la création du paquet (#4575)
- Ajout par @leaanthony d’une documentation sur les dépendances d’exécution sous Linux, avec les noms de paquets propres à chaque distribution et des exemples de création de paquets nfpm (#4339)
- Ajout par @leaanthony des informations sur la version du pilote NVIDIA à la sortie de `wails3 doctor` sous Linux
- Ajout de l’origine au gestionnaire de messages bruts par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4710)
- Ajout de la prise en charge des liens universels sous macOS par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4712)
- Remaniement de la couche de transport des liaisons par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4702)
- Ajout d’identifiants aria-label aux modèles helloworld afin que l’application d’exemple puisse être facilement testée par des clients de test Appium, par @chinenual dans la [PR](https://github.com/wailsapp/wails/pull/4760)
- Ajout de l’origine au gestionnaire de messages bruts par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4710)
- Ajout de la prise en charge des liens universels sous macOS par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4712)
- Remaniement de la couche de transport des liaisons par @APshenkin dans la [PR](https://github.com/wailsapp/wails/pull/4702)
- Événements typés par @fbbdev et @ianvs dans [#4633](https://github.com/wailsapp/wails/pull/4633)
- Ajout de l’exemple `systray-clock` présentant une icône de zone de notification sans fenêtre, avec mise à jour dynamique de l’infobulle (#4653).
- Ajout par @Tolfx d’un modèle de protocole NSIS pour Windows dans #4510
- Ajout par @Tolfx de tests pour build-assets dans #4510
- macOS : affichage des commandes de fenêtre natives dans la barre des menus, par @nidib dans [#4588](https://github.com/wailsapp/wails/pull/4588)
- Ajout par @popaprozac d’un service Dock macOS permettant de masquer ou d’afficher l’icône de l’application dans le Dock, dans la [PR](https://github.com/wailsapp/wails/pull/4451)
- Ajout par @popaprozac d’un service Dock macOS permettant de masquer ou d’afficher l’icône de l’application dans le Dock, dans la [PR](https://github.com/wailsapp/wails/pull/4451)
- Ajout par @leaanthony de la prise en charge native de l’effet Liquid Glass sous macOS avec NSGlassEffectView (macOS 15.0+) et repli sur NSVisualEffectView, ainsi que d’options complètes de personnalisation des matériaux, dans [#4534](https://github.com/wailsapp/wails/pull/4534)
- Assainissement des URL du navigateur par @leaanthony dans [#4500](https://github.dev/wailsapp/wails/pull/4500). Basé sur [#4484](https://github.com/wailsapp/wails/pull/4484) par @APShenkin.
- Ajout de la protection du contenu sous Windows/Mac par [@leaanthony](https://github.com/leaanthony), à partir du travail initial de [@Taiterbase](https://github.com/Taiterbase) dans cette [PR](https://github.com/wailsapp/wails/pull/4241)
- Ajout par @leaanthony de la prise en charge de la transmission de variables CLI aux commandes Task via les alias `wails3 build` et `wails3 package` (#4422), dans la [PR](https://github.com/wailsapp/wails/pull/4488)
- Prise en charge des zones de dépôt dont les événements fournissent les données des éléments déposés, par [@atterpac](https://github.com/atterpac) dans [#4318](https://github.com/wailsapp/wails/pull/4318)
- Ajout de `AdditionalLaunchArgs` aux options de `WindowsWindow` afin de pouvoir transmettre des arguments de ligne de commande supplémentaires au navigateur WebView2, dans la [PR](https://github.com/wailsapp/wails/pull/4467)
- Ajout de l’exécution automatique de go mod tidy après wails init par [@triadmoko](https://github.com/triadmoko) dans la [PR](https://github.com/wailsapp/wails/pull/4286)
- Fonctionnalité Snap Assist de Windows par @leaanthony dans la [PR](https://github.dev/wailsapp/wails/pull/4463)
- Ajout de `AdditionalLaunchArgs` aux options de `WindowsWindow` afin de pouvoir transmettre des arguments de ligne de commande supplémentaires au navigateur WebView2, dans la [PR](https://github.com/wailsapp/wails/pull/4467)
- Ajout de l’exécution automatique de go mod tidy après wails init par [@triadmoko](https://github.com/triadmoko) dans la [PR](https://github.com/wailsapp/wails/pull/4286)
- Fonctionnalité Snap Assist de Windows par @leaanthony dans la [PR](https://github.dev/wailsapp/wails/pull/4463)
- Ajout de l’implémentation Windows de `getAccentColor` par [@almas-x](https://github.com/almas-x) dans la [PR](https://github.com/wailsapp/wails/pull/4427)
- Ajout de l’implémentation Windows de `getAccentColor` par [@almas-x](https://github.com/almas-x) dans la [PR](https://github.com/wailsapp/wails/pull/4427)
- Menus et barre des menus en thème sombre sous Windows. Par @leaanthony dans [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Renommage des services intégrés afin de rendre les liaisons JS/TS plus explicites, par @popaprozac dans la [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` permet d’obtenir la couleur d’accentuation du système de l’utilisateur. Fonctionne sous macOS. Par [@etesam913](https://github.com/etesam913)
- Ajout de l’API `window.ToggleFrameless()` par [@atterpac](https://github.com/atterpac) dans [#4137](https://github.com/wailsapp/wails/pull/4137)
- Ajout par @leaanthony des dépendances de compilation propres à chaque distribution Linux dans la [PR](https://github.com/wailsapp/wails/pull/4345)
- Ajout d’un guide sur les liaisons par @atterpac dans la [PR](https://github.com/wailsapp/wails/pull/4404)
- **Infrastructure de test réorganisée** : déplacement des fichiers de test Docker vers le répertoire dédié `test/docker/`, avec des images optimisées et une fiabilité accrue des builds, par [@leaanthony](https://github.com/leaanthony) dans [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Amélioration des modèles de gestion des ressources** : ajout dans les exemples d’un nettoyage correct des gestionnaires d’événements et d’une gestion des goroutines tenant compte du contexte, par [@leaanthony](https://github.com/leaanthony) dans [#4359](https://github.com/wailsapp/wails/pull/4359)
- Prise en charge des builds AppImage pour aarch64 par [@AkshayKalose](https://github.com/AkshayKalose) dans [#3981](https://github.com/wailsapp/wails/pull/3981)
- Ajout d’une section sur les diagnostics à `wails doctor` par [@leaanthony](https://github.com/leaanthony)
- Ajout de la fenêtre au contexte lors de l’appel d’une méthode de service par [@leaanthony](https://github.com/leaanthony)
- Ajout d’un exemple `window-call` montrant comment déterminer quelle fenêtre appelle un service, par [@leaanthony](https://github.com/leaanthony)
- Nouveau guide sur les menus par [@leaanthony](https://github.com/leaanthony)
- Amélioration de la gestion des paniques par [@leaanthony](https://github.com/leaanthony)
- Nouveau guide sur les menus par [@leaanthony](https://github.com/leaanthony)
- Ajout de commentaires de documentation pour l’API Service par [@fbbdev](https://github.com/fbbdev) dans [#4024](https://github.com/wailsapp/wails/pull/4024)
- Ajout de la fonction `application.NewServiceWithOptions` pour initialiser les services avec une configuration supplémentaire, par [@leaanthony](https://github.com/leaanthony) dans [#4024](https://github.com/wailsapp/wails/pull/4024)
- Amélioration du contrôle des menus par [@FalcoG](https://github.com/FalcoG) et [@leaanthony](https://github.com/leaanthony) dans [#4031](https://github.com/wailsapp/wails/pull/4031)
- Documentation supplémentaire par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de l’annulation des événements dans les écouteurs d’événements standard par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de `Hide`, `Show` et `Destroy` dans la zone de notification par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de `SetTooltip` dans la zone de notification par [@leaanthony](https://github.com/leaanthony). Idée originale de [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Affichage du chemin du paquet dans les avertissements du générateur de liaisons concernant les types non pris en charge, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la prise en charge des alias génériques par le générateur de liaisons, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la prise en charge de l’indicateur JSON `omitzero` par le générateur de liaisons, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la directive `//wails:ignore` pour empêcher la génération de liaisons pour certaines méthodes de service, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la directive `//wails:internal` sur les services et les modèles afin d’autoriser des types exportés en Go, mais pas en JS/TS, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de la prise en charge des constantes de type alias par le générateur de liaisons afin de permettre des énumérations faiblement typées, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ajout de tests du générateur de liaisons pour les fonctionnalités 1.24 de Go, par [@fbbdev](https://github.com/fbbdev) dans [#4068](https://github.com/wailsapp/wails/pull/4068)
- Ajout de la prise en charge de macOS 15 « Sequoia » à `OSInfo.Branding` pour améliorer la détection de la version du système d’exploitation dans [#4065](https://github.com/wailsapp/wails/pull/4065)
- Ajout du hook `PostShutdown` pour exécuter du code personnalisé une fois le processus d’arrêt terminé, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout de la structure `FatalError` pour permettre la détection des erreurs fatales dans les gestionnaires d’erreurs personnalisés, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Standardisation et documentation de l’ordre de démarrage et d’arrêt des services par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout d’un banc de test pour la séquence de démarrage et d’arrêt de l’application, ainsi que de tests de démarrage et d’arrêt des services, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout de la méthode `RegisterService` pour enregistrer des services après la création de l’application, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout du champ `MarshalError` dans les options de l’application et des services pour personnaliser la gestion des erreurs lors des appels de liaison, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ajout d’un wrapper de promesse annulable qui propage les demandes d’annulation dans les chaînes de promesses, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Ajout de la possibilité de lier l’annulation d’un appel de liaison à un `AbortSignal`, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Prise en charge des attributs `data-wml-*` pour WML en plus des attributs `wml-*` habituels, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la méthode `Configure` à tous les services pour permettre une configuration tardive ou une reconfiguration dynamique, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Lorsqu’il n’est pas configuré, le service `fileserver` envoie une réponse 503 Service Unavailable, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Lorsqu’il n’est pas configuré, le service `kvstore` fournit par défaut un magasin clé-valeur en mémoire, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Load` au service `kvstore` pour recharger les données depuis le fichier après une modification de la configuration, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Clear` au service `kvstore` pour supprimer toutes les clés, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout du type `Level` au service `log` pour fournir des constantes de niveau de journalisation côté JS, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Log` au service `log` pour définir dynamiquement le niveau de journalisation, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Lorsqu’il n’est pas configuré, le service `sqlite` fournit par défaut une base de données en mémoire, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la méthode `Close` au service `sqlite` pour fermer manuellement la base de données, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la prise en charge de l’annulation pour les méthodes de requête du service `sqlite`, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ajout de la prise en charge des instructions préparées au service `sqlite` avec des liaisons JS, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Prise en charge de Gin par [Lea Anthony](https://github.com/leaanthony) dans la [PR](https://github.com/wailsapp/wails/pull/3537), d’après le travail initial de [@AnalogJ](https://github.com/AnalogJ) dans cette [PR](https://github.com/wailsapp/wails/pull/3537)
- Correction de l’enregistrement automatique et de l’enregistrement automatique des mots de passe qui étaient toujours activés, par [@oSethoum](https://github.com/osethoum) dans [#4134](https://github.com/wailsapp/wails/pull/4134)
- Ajout de `SetMenu()` aux fenêtres pour permettre de leur attribuer un menu, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la prise en charge des notifications par [@popaprozac](https://github.com/popaprozac) dans [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Ajout de la prise en charge des associations de fichiers sur Mac par [@wimaha](https://github.com/wimaha) dans [#4177](https://github.com/wailsapp/wails/pull/4177)
- Ajout de `wails3 tool version` pour l’incrémentation des versions sémantiques, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la prise en charge des badges sous macOS et Windows par [@popaprozac](https://github.com/popaprozac) dans [#](https://github.com/wailsapp/wails/pull/4234)
- Ajout de la prise en charge des événements enregistrés et strictement typés par [@fbbdev](https://github.com/fbbdev) et [@IanVS](https://github.com/IanVS) dans [#4161](https://github.com/wailsapp/wails/pull/4161)
- Ajout de la possibilité d’enregistrer des hooks pour les événements personnalisés par [@fbbdev](https://github.com/fbbdev) et [@IanVS](https://github.com/IanVS) dans [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` pour ouvrir le gestionnaire de fichiers système au chemin `path`, avec une mise en surbrillance facultative au moyen de `selectFile`, par [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Nouvelle option `-git` pour la commande `wails3 init`, par [@leaanthony](https://github.com/leaanthony)
- Nouvelle commande `wails3 generate webview2bootstrapper` par [@leaanthony](https://github.com/leaanthony)
- Ajout de la méthode `init()` au runtime pour permettre son initialisation manuelle, par [@leaanthony](https://github.com/leaanthony)
- Ajout de l’option `WindowDidMoveDebounceMS` aux WindowOptions des fenêtres, par [@leaanthony](https://github.com/leaanthony)
- Ajout de la fonctionnalité d’instance unique par [@leaanthony](https://github.com/leaanthony). Basée sur la PR [v2](https://github.com/wailsapp/wails/pull/2951) de @APshenkin.
- Commande `wails3 generate template` par [@leaanthony](https://github.com/leaanthony)
- Commande `wails3 releasenotes` par [@leaanthony](https://github.com/leaanthony)
- Commande `wails3 update cli` par [@leaanthony](https://github.com/leaanthony)
- Option `-clean` pour la commande `wails3 generate bindings`, par [@leaanthony](https://github.com/leaanthony)
- Prise en charge de la création d’AppImage Linux pour aarch64 (arm64) par [@AkshayKalose](https://github.com/AkshayKalose) dans [#3981](https://github.com/wailsapp/wails/pull/3981)
- Ajout d’un lien hypertexte vers le sponsor par @ansxuman dans [#3958](https://github.com/wailsapp/wails/pull/3958)
- Prise en charge de la création de paquets Linux deb, rpm et Arch Linux par
- Ajout de la prise en charge des builds et paquets universels Darwin par
- Ajout de la documentation sur les événements au site web par
- Modèles pour sveltekit et sveltekit-ts configurés pour le développement sans SSR
- Mise à jour des ressources de build à l’aide de la nouvelle commande `wails3 update build-assets` par
- Exemple permettant de tester l’API HTML Drag and Drop par
- Prise en charge des associations de fichiers par [leaanthony](https://github.com/leaanthony) dans
- Nouvelle commande `wails3 generate runtime` par
- Nouvelle option `InitialPosition` permettant d’indiquer si la fenêtre doit être centrée ou
- Ajout des méthodes `Path` et `Paths` au paquet `application` par
- Ajout des options Windows `GeneralAutofillEnabled` et `PasswordAutosaveEnabled`
- Ajout de la possibilité de récupérer la fenêtre qui appelle une méthode de service par
- Ajout des options `EnabledFeatures` et `DisabledFeatures` pour WebView2 par
- ⊞ Nouveau système DIP pour une meilleure prise en charge des moniteurs à haute densité de pixels par
- ⊞ Option de nom de classe de fenêtre par [windom](https://github.com/windom/) dans
- Extension des services pour fournir des fonctionnalités de plugin. Par
- 🐧 Événements WindowDidMove et WindowDidResize dans
- ⊞ Événement WindowDidResize dans
-  Ajout de l’événement ApplicationShouldHandleReopen afin de pouvoir gérer le Dock
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
- Ajout d’exports ESM depuis le script du runtime JS intégré par
- Ajout d’une option au générateur de bindings pour utiliser le script du runtime JS intégré au lieu de
- Implémentation de `setIcon` sous Linux par [@abichinger](https://github.com/abichinger)
- Ajout de l’option `-port` à la commande dev et prise en charge de la variable d’environnement
- Ajout de tests pour les appels de méthodes liées par
- ⊞ Ajout de `SetIgnoreMouseEvents` pour une fenêtre déjà créée par
-  Ajout de la possibilité de définir le niveau d’empilement (ordre) d’une fenêtre par

### Corrections

- Correction des valeurs `Screen.Bounds`, `WorkArea` et `Size` divisées par deux sur les Mac Retina, en convertissant les valeurs en points de NSScreen en pixels de périphérique dans les champs `Physical*` et en renseignant les valeurs `Screen.X`/`Y` de premier niveau, afin que la détection du contact entre plusieurs moniteurs et le positionnement dans la zone de travail soient corrects, dans la [PR](https://github.com/wailsapp/wails/pull/5168) par @wayneforrest
- Correction d’une concurrence critique dans ScreenManager qui provoquait un interblocage de WebKit DisplayLink lors d’une modification de la configuration des écrans (par exemple, le branchement à chaud d’un moniteur externe pendant une mise en veille ou une sortie de veille)
- Définition directe de CFBundleIconName sur appicon lorsque Assets.car existe, dans la [PR](https://github.com/wailsapp/wails/pull/5154) par @symball
- Correction de `wails3 doctor` qui signalait des paquets WebKitGTK incorrects sous Fedora, openSUSE, Arch et NixOS — les entrées de repli 4.0 ont été supprimées, car v3 nécessite l’API 4.1 lors de la compilation (#5071)
- Correction du nom du paquet webkit2gtk vérifié par doctor sous openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, qui est le nom correct du paquet openSUSE) (#5071)
- Corrige l’erreur `Unexpected token '<'` lorsque `/wails/custom.js` est absent en mode de développement desktop. Ajoute un gestionnaire 404 explicite pour `/wails/custom.js` ainsi qu’une validation de `Content-Type` insensible à la casse dans `loadOptionalScript`, afin d’empêcher l’injection des pages HTML de repli des SPA en tant que JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Corrige l’état de surbrillance du menu de la zone de notification sous macOS : l’icône apparaît désormais sélectionnée lorsque le menu est ouvert (#4910)
- Corrige l’affichage de la fenêtre attachée à la zone de notification derrière d’autres fenêtres sous macOS : elle utilise désormais le niveau approprié pour une fenêtre contextuelle (#4910)
- Corrige les exemples d’importation `@wailsio/runtime` incorrects dans toute la documentation (#4989)
- Corrige l’impossibilité de réduire une fenêtre sans cadre sous Darwin (#4294)
- Corrige les blocages de 20-30 minutes pendant `wails3 build` et `wails3 dev` en excluant `node_modules/` de la vérification d’actualisation de go-task. Auparavant, le motif glob `sources: "**/*"` amenait go-task à énumérer et à calculer la somme de contrôle de chaque fichier dans `node_modules/` (plus de 50000-100000 fichiers, avec de lourdes dépendances telles que MUI), une opération particulièrement lente sous Windows/NTFS (#4939)
- Corrige l’échec de compilation avec GTK4 provoqué par un conflit entre le typedef C `Screen` et X11 Xlib.h (#4957)
- Harmonise les méthodes de badge du Dock sous macOS
- Corrige l’application de `InvisibleTitleBarHeight` à toutes les fenêtres macOS au lieu des seules fenêtres sans cadre ou à barre de titre transparente (#4960)
- Corrige les secousses et tremblements de la fenêtre lors de son redimensionnement depuis les coins supérieurs avec `InvisibleTitleBarHeight` activé, en empêchant le démarrage du déplacement près des bords de la fenêtre (#4960)
- Corrige la génération des types mappés dont les clés sont des énumérations dans les liaisons JS/TS (#4437), par @fbbdev
- Corrige le glisser-déposer de fichiers sous Windows, qui ne fonctionnait pas avec une mise à l’échelle d’affichage différente de 100 %
- Corrige le glisser-déposer interne HTML5, qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Windows
- Corrige les coordonnées de dépôt de fichiers sous Windows, qui utilisaient le mauvais espace de pixels (pixels physiques au lieu de pixels CSS)
- Corrige le manque de fiabilité du glisser-déposer de fichiers sous Linux avec des effets de survol
- Corrige le glisser-déposer interne HTML5, qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Linux
- Corrige l’affichage et le masquage d’une fenêtre sous Linux/GTK4, qui la restauraient parfois à l’état réduit, grâce à l’utilisation de `gtk_window_present()` (#4957)
- Corrige la lecture et la définition de la position d’une fenêtre sous Linux/GTK4, qui renvoyaient toujours 0,0, en ajoutant une prise en charge conditionnelle de X11 via `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Corrige l’absence d’application de la taille maximale des fenêtres sous Linux/GTK4 en remplaçant la fonction `gtk_window_set_geometry_hints` supprimée par un plafonnement de taille fondé sur les signaux (#4957)
- Corrige la mise à l’échelle selon le DPI sous Linux/GTK4 en implémentant un calcul correct de PhysicalBounds et la prise en charge de la mise à l’échelle fractionnaire via `gdk_monitor_get_scale` (GTK 4.14+)
- Corrige la duplication des éléments de menu lors de la création de nouvelles fenêtres sous Linux/GTK4
- Corrige la génération des types mappés dont les clés sont des énumérations dans les liaisons JS/TS (#4437), par @fbbdev
- Corrige le glisser-déposer de fichiers sous Windows, qui ne fonctionnait pas avec une mise à l’échelle d’affichage différente de 100 %
- Corrige le glisser-déposer interne HTML5, qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Windows
- Corrige les coordonnées de dépôt de fichiers sous Windows, qui utilisaient le mauvais espace de pixels (pixels physiques au lieu de pixels CSS)
- Corrige le manque de fiabilité du glisser-déposer de fichiers sous Linux avec des effets de survol
- Corrige le glisser-déposer interne HTML5, qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Linux
- Corrige la mise à l’échelle selon le DPI sous Linux/GTK4 en implémentant un calcul correct de PhysicalBounds et la prise en charge de la mise à l’échelle fractionnaire via `gdk_monitor_get_scale` (GTK 4.14+)
- Corrige la duplication des éléments de menu lors de la création de nouvelles fenêtres sous Linux/GTK4
- Corrige la génération des types mappés dont les clés sont des énumérations dans les liaisons JS/TS (#4437), par @fbbdev
- Corrige le problème de « fenêtres fantômes » sous macOS, provoqué par l’accès aux API AppKit en dehors du thread principal dans App.Window.Current() (#4947), par @wimaha
- Corrige l’élément HTML `<input type="file">`, qui ne fonctionnait pas sous macOS, en implémentant WKUIDelegate runOpenPanelWithParameters (#4862)
- Corrige le glisser-déposer natif de fichiers, qui ne fonctionnait pas lors de l’utilisation du module npm `@wailsio/runtime` sous macOS/Linux (#4953), par @leaanthony
- Corrige la génération des liaisons pour les alias de types inter-paquets (#4578), par @fbbdev
- Corrige le plantage d’OpenFileDialog sous Linux dû à une violation de la sécurité des threads GTK (#3683), par @ddmoney420
- Corrige le plantage SIGSEGV lors de l’appel de `Focus()` sur une fenêtre masquée ou détruite (#4890), par @ddmoney420
- Corrige une panique potentielle lors de la définition d’une icône ou d’une image bitmap vide sous Linux (#4923), par @ddmoney420
- Corrige le plantage d’ErrorDialog sous macOS lorsqu’il est appelé depuis une liaison de service (#3631), par @leaanthony
- Permet l’affichage des menus sous Windows dans `v3\examples\dialogs`, par @ndianabasi
- Corrige une condition de concurrence qui provoquait une TypeError lors du rechargement de la page (#4872), par @ddmoney420
- Corrige la sortie incorrecte des tests du générateur de liaisons en supprimant l’état global de la méthode `Collector.IsVoidAlias()` (#4941), par @fbbdev
- Corrige le sélecteur de fichiers `<input type="file">`, qui ne fonctionnait pas sous macOS (#4862), par @leaanthony
- Corrige l’utilisation de systèmes de coordonnées incohérents par `Position()` et `SetPosition()` sous macOS, qui entraînait un décalage de la position de la fenêtre lors de l’enregistrement et de la restauration de son état (#4816), par @leaanthony
- Corrige l’erreur « Access is denied » de SetProcessDpiAwarenessContext lorsque la prise en charge du DPI est déjà définie dans le manifeste de l’application (#4803)
- Met à jour la page de documentation sur les raccourcis clavier et corrige le type du paramètre de rappel de `KeyBinding.Add`, par @ndianabasi
- Corrige la documentation sur la génération d’une liaison personnalisée : il faut utiliser `-d String` au lieu de `-o String`
- Corrige le menu, qui n’effaçait pas ses éléments enfants lors de `menu.Update()`
- Corrige les références obsolètes à l’API Manager dans la documentation (les fichiers 31 utilisent désormais le nouveau modèle, tel que `app.Window.New()`, `app.Event.Emit()`, etc.), par @leaanthony
- Corrige le plantage de Linux en cas de panique dans les méthodes Go liées à JS, dû au remplacement des gestionnaires de signaux par WebKit (#3965), par @leaanthony
- Corrige SaveFileDialog.SetFilename(), qui n’avait aucun effet sous Linux (#4841), par @samstanier
- Corrige les coordonnées de dépôt affichées comme undefined dans l’exemple de glisser-déposer
- Corrige l’échec de création du paquet d’application macOS lorsque APP_NAME contient des espaces (problème de développement des accolades)
- Corrige une panique d’indice hors limites sous Windows lors de l’appel de méthodes de service (rétablissement de la version antérieure de goccy/go-json)
- Corrige le glisser-déposer de fichiers sous Windows, qui ne fonctionnait pas avec une mise à l’échelle d’affichage différente de 100 %
- Corrige le glisser-déposer interne HTML5, qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Windows
- Corrige les coordonnées de dépôt de fichiers sous Windows, qui utilisaient le mauvais espace de pixels (pixels physiques au lieu de pixels CSS)
- Corrige le manque de fiabilité du glisser-déposer de fichiers sous Linux avec des effets de survol
- Corrige le glisser-déposer interne HTML5, qui ne fonctionnait plus lorsque le dépôt de fichiers était activé sous Linux
- Met à jour toutes les commandes des fichiers Taskfile.yml pour tous les systèmes d’exploitation afin de prendre en charge les espaces dans des variables telles que `APP_NAME`, par @ndianabasi
- Corrige l’erreur d’argument de commande lors de l’exécution de la tâche 'build:universal:lipo:go' sous Linux, par @wux1an
- Correction de l’erreur Docker « undefined symbol: **<em>ubsan</em>handle_xxxxxxx » lors de l’exécution de 'wails3 build GOOS=darwin GOARCH=arm64' sous Linux, par @wux1an
- Regroupement de la documentation sur les protocoles personnalisés et ajout de sections sur les Universal Links, par @leaanthony
- Correction du plantage du menu de la zone de notification Windows lors de clics répétés sur l’icône, grâce à l’ajout d’une protection contre les appels simultanés à TrackPopupMenuEx (#4151), par @leaanthony
- Prévention du plantage de l’application lorsque systray.Run() est appelé avant app.Run(), par @leaanthony
- Correction d’un plantage sous macOS lors du basculement de la visibilité d’une fenêtre avec Hide()/Show() lorsque ApplicationShouldTerminateAfterLastWindowClosed est activé (#4389), par @leaanthony
- Correction d’une fuite de mémoire dans les menus contextuels sous macOS et Windows lors de leur ouverture répétée (#4012), par @leaanthony
- Correction de la non-réutilisation des ressources natives des menus contextuels sous macOS, qui entraînait la création d’un nouveau menu à chaque affichage (#4012), par @leaanthony
- Correction du clic sur l’icône du Dock macOS qui n’affichait pas les fenêtres masquées lorsque l’application était lancée avec `Hidden: true` (#4583), par @leaanthony
- Correction de la boîte de dialogue d’impression qui ne s’ouvrait pas sous macOS en raison d’un type de pointeur de fenêtre incorrect dans l’appel CGO (#4290), par @leaanthony
- Correction du plantage du menu de fenêtre sous Wayland, causé par l’accès d’appmenu-gtk-module à une fenêtre non réalisée (#4769), par @leaanthony
- Correction du plantage de l’application GTK lorsque son nom contient des caractères non valides (espaces, parenthèses, etc.), par @leaanthony
- Correction de l’erreur « not enough memory » lors de l’initialisation du glisser-déposer sous Windows (#4701), par @overlordtm
- Correction de l’ouverture du mauvais répertoire par l’explorateur de fichiers sous Linux en raison d’un échappement incorrect de l’URI (#4397), par @leaanthony
- Correction de l’échec de construction d’une AppImage sur les distributions Linux modernes (Arch, Fedora 39+, Ubuntu 24.04+) grâce à la détection automatique des sections ELF `.relr.dyn` et à la désactivation du dépouillement des symboles (#4642), par @leaanthony
- Correction de `wails doctor` qui signalait à tort les paquets WebKit comme installés sur Fedora et les systèmes basés sur DNF (#4457), par @leaanthony
- Correction du comportement par défaut de `config.yml`, qui exécutait `wails3 dev` avec une version de production, par @mbaklor
- Correction des stubs de services iOS qui provoquaient des échecs de construction en important un paquet inexistant, par @leaanthony
- Correction de la journalisation structurée dans les méthodes debug/info, qui provoquait des erreurs « no formatting directives », par @leaanthony
- Suppression des instructions temporaires d’affichage de débogage incluses accidentellement lors de la fusion des plateformes mobiles, par @leaanthony
- Correction du plantage de WebKitGTK sous Wayland avec les GPU NVIDIA (erreur de protocole 71) grâce à la désactivation automatique du moteur de rendu DMA-BUF, par @leaanthony
- Correction de la valeur alpha ignorée dans `application.WebviewWindowOptions.BackgroundColour` sous Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Correction de l’icône de la zone de notification Windows qui n’utilisait pas par défaut celle de l’application lorsqu’aucune icône personnalisée n’était fournie (#4704)
- Suivi de la propriété des `HICON` afin de ne détruire que les handles créés par l’utilisateur, ce qui évite les plantages lors de leur recyclage par l’Explorateur (#4653).
- Libération de l’écouteur du thème système Windows et des icônes conservées dans la zone de notification lors de la destruction, afin d’éliminer les fuites de goroutines et de contextes de périphérique (#4653).
- Troncature des infobulles de la zone de notification à 127 unités UTF-16 afin d’éviter d’altérer les paires de substitution et les glyphes multioctets (#4653).
- Correction de l’échec de la tâche de création du paquet Windows (#4667)
- Correction de la variable appicon de l’AppImage Linux dans le Taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Correction de l’erreur de construction sous Windows causée par la modification de la signature de go-webview2 v1.0.22 (#4513, #4645)
- Correction de la variable appicon de l’AppImage Linux dans le Taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Correction de la plage de protocoles dans desktop.tmpl sous Linux par la suppression de `<.Info.Protocol>` à `<.Protocol>`, par @Tolfx dans #4510
- Correction de l’erreur de redéfinition dans la démonstration Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Correction des mises à jour du menu de la zone de notification sous Linux [#4604](https://github.com/wailsapp/wails/issues/4604), par [@JackDoan](https://github.com/JackDoan)
- Correction de la fenêtre blanche qui apparaissait sous Windows lors de la création d’une fenêtre masquée, par @leaanthony dans [#4612](https://github.com/wailsapp/wails/pull/4612)
- Correction du chemin d’importation du paquet de notifications dans la documentation, par @rxliuli dans [#4617](https://github.com/wailsapp/wails/pull/4617)
- Correction du glisser-déposer qui ne fonctionnait pas avec le paquet npm @wailsio/runtime (#4489), par @leaanthony dans #4616
- Windows : correction du scintillement de la fenêtre au démarrage et de l’affichage incorrect des fenêtres masquées dans [PR](https://github.com/wailsapp/wails/pull/4600), par @leaanthony.
- Correction des problèmes de maximisation de la taille des fenêtres sous Wayland (https://github.com/wailsapp/wails/issues/4429) par [@samstanier](https://github.com/samstanier)
- Correction des problèmes de maximisation de la taille des fenêtres sous Wayland (https://github.com/wailsapp/wails/issues/4429) par [@samstanier](https://github.com/samstanier)
- Correction de l’erreur de redéfinition dans la démonstration Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Correction d’un problème pouvant faire planter AssetServer sous macOS dans [#4576](https://github.com/wailsapp/wails/pull/4576), par @jghiloni
- Correction du problème de compilation lors de la construction avec NextJs, dans [#4585](https://github.com/wailsapp/wails/pull/4585), par @rev42
- Correction des pipelines de publication nocturne dans [#4597](https://github.com/wailsapp/wails/pull/4597), par @riadafridishibly
- Correction de l’erreur de redéfinition dans la démonstration Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Correction d’un problème pouvant faire planter AssetServer sous macOS dans [#4576](https://github.com/wailsapp/wails/pull/4576), par @jghiloni
- Correction du problème de compilation lors de la construction avec NextJs, dans [#4585](https://github.com/wailsapp/wails/pull/4585), par @rev42
- Correction des pipelines de publication nocturne dans [#4597](https://github.com/wailsapp/wails/pull/4597), par @riadafridishibly
- Correction de l’erreur de redéfinition dans la démonstration Liquid Glass dans [#4542](https://github.com/wailsapp/wails/pull/4542), par @Etesam913
- Correction de SetBackgroundColour sous Windows, par @PPTGamer dans [PR](https://github.com/wailsapp/wails/pull/4492)
- Mise à jour de la documentation pour refléter les modifications issues de la refactorisation de l’API Manager, par @yulesxoxo dans [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Correction de la variable appicon du fichier .desktop Linux dans le Taskfile Linux [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Mise à jour de la documentation pour refléter les modifications issues de la refactorisation de l’API Manager, par @yulesxoxo dans [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Correction du déréférencement d’un pointeur nil sous Windows, signalé dans [#4456](https://github.com/wailsapp/wails/issues/4456), par @leaanthony dans [#4460](https://github.com/wailsapp/wails/pull/4460)
- Ajout de la prise en charge de `allowsBackForwardNavigationGestures` dans WKWebView sous macOS afin d’activer les gestes de navigation par balayage à deux doigts (#1857)
- Correction d’un problème empêchant onClick de fonctionner pour les éléments de menu initialement désactivés, par @leaanthony dans [PR #4469](https://github.com/wailsapp/wails/pull/4469). Merci à @IanVS pour l’analyse initiale.
- Correction du serveur Vite qui n’était pas arrêté lorsque la construction échouait (#4403)
- Correction d’une panique lors de la fermeture ou de l’annulation d’un `SaveFileDialog` sous Windows, dans [PR](https://github.com/wailsapp/wails/pull/4284), par @hkhere
- Correction du glisser-déposer au niveau HTML sous Windows par [@mbaklor](https://github.com/mbaklor) dans [#4259](https://github.com/wailsapp/wails/pull/4259)
- Ajout de la prise en charge de `allowsBackForwardNavigationGestures` dans WKWebView sous macOS afin d’activer les gestes de navigation par balayage à deux doigts (#1857)
- Correction d’un problème qui empêchait onClick de fonctionner pour les éléments de menu initialement définis comme désactivés, par @leaanthony dans la [PR #4469](https://github.com/wailsapp/wails/pull/4469). Merci à @IanVS pour l’investigation initiale.
- Correction du serveur Vite qui n’était pas arrêté en cas d’échec de la compilation (#4403)
- Correction de l’analyse des notifications sous Windows par @popaprozac dans la [PR](https://github.com/wailsapp/wails/pull/4450)
- Correction de la commande doctor afin qu’elle vérifie les dépendances du SDK Windows, par [@kodumulo](https://github.com/kodumulo) dans [#4390](https://github.com/wailsapp/wails/issues/4390)
- Correction d’un déréférencement de pointeur nil dans processURLRequest sur Mac par [@etesam913](https://github.com/etesam913) dans [#4366](https://github.com/wailsapp/wails/pull/4366)
- Correction d’un bogue sous Linux qui empêchait l’utilisation de boîtes de dialogue filtrées, par [@bh90210](https://github.com/bh90210) dans [#4287](https://github.com/wailsapp/wails/pull/4287)
- Correction de problèmes du menu Édition sous Windows et Linux par [@leaanthony](https://github.com/leaanthony) dans [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Mise à jour de la version minimale du système dans les fichiers .plist de macOS, de 10.13.0 à 10.15.0, par [@AkshayKalose](https://github.com/AkshayKalose) dans [#3981](https://github.com/wailsapp/wails/pull/3981)
- Correction du problème de saut d’identifiants de fenêtre par [@leaanthony](https://github.com/leaanthony)
- Correction du problème de menu nil lors de l’appel à RegisterContextMenu par [@leaanthony](https://github.com/leaanthony)
- Correction des cycles de dépendances dans la sortie du générateur de bindings par [@fbbdev](https://github.com/fbbdev) dans [#4001](https://github.com/wailsapp/wails/pull/4001)
- Correction des erreurs d’utilisation avant définition dans la sortie du générateur de bindings par [@fbbdev](https://github.com/fbbdev) dans [#4001](https://github.com/wailsapp/wails/pull/4001)
- Transmission des indicateurs de compilation au générateur de bindings par [@fbbdev](https://github.com/fbbdev) dans [#4023](https://github.com/wailsapp/wails/pull/4023)
- Remplacement des barres obliques inverses par des barres obliques dans les chemins du Taskfile Windows afin d’assurer son fonctionnement sur les plateformes autres que Windows, par [@leaanthony](https://github.com/leaanthony)
- Correction des événements JS Mac + Mac par [@leaanthony](https://github.com/leaanthony)
- Correction d’un interblocage des événements sous macOS par [@leaanthony](https://github.com/leaanthony)
- Correction d’une erreur `Parameter incorrect` lors de l’initialisation d’une fenêtre sous Windows lorsque du HTML est fourni sans JS, par [@leaanthony](https://github.com/leaanthony)
- Correction de la taille du préfixe de réponse utilisé pour détecter le type de contenu dans le serveur de ressources, par [@fbbdev](https://github.com/fbbdev) dans [#4049](https://github.com/wailsapp/wails/pull/4049)
- Correction du traitement des réponses non-404 sur le chemin de l’index racine dans le serveur de ressources, par [@fbbdev](https://github.com/fbbdev) dans [#4049](https://github.com/wailsapp/wails/pull/4049)
- Correction d’un comportement indéfini dans le générateur de bindings lors du test des propriétés de types génériques, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la sortie du générateur de bindings pour les modèles dont le type sous-jacent ne possède pas les mêmes propriétés que le wrapper nommé, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la sortie du générateur de bindings pour les types de clés de map et le prétraitement, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la sortie du générateur de bindings pour les structs qui implémentent des interfaces de sérialisation, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction de la détection des cycles de types impliquant des types génériques dans le générateur de bindings, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction des références non valides à des modèles non exportés dans la sortie du générateur de bindings, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Déplacement du code injecté à la fin des fichiers de service par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction du traitement des erreurs provenant des opérations de fermeture de fichiers dans le générateur de bindings, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Suppression des avertissements pour les services qui définissent des méthodes de cycle de vie ou HTTP, mais aucune autre méthode liée, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Correction des modèles autres que React qui n’affichaient pas le pied de page Hello World lorsque le jeu de couleurs clair du système était utilisé, par [@marcus-crane](https://github.com/marcus-crane) dans [#4056](https://github.com/wailsapp/wails/pull/4056)
- Correction des éléments de menu masqués sous macOS par [@leaanthony](https://github.com/leaanthony)
- Correction du traitement et du formatage des erreurs dans les processeurs de messages, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Correction de l’arrêt des services qui était omis lors de la fermeture de l’application, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Exécution des mises à jour du menu sur le thread principal par [@leaanthony](https://github.com/leaanthony)
- Le mécanisme de déplacement et de redimensionnement est désormais plus robuste et correspond davantage au comportement attendu de la plateforme, grâce à [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Correction de [#4097](https://github.com/wailsapp/wails/issues/4097) : Webpack/Angular éliminait le code d’initialisation du runtime, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Correction des éléments de menu initialement masqués par [@IanVS](https://github.com/IanVS) dans [#4116](https://github.com/wailsapp/wails/pull/4116)
- Correction d’assetFileServer, qui ne servait pas les fichiers `.html` lors d’une requête sans extension lorsque `[request]` n’existait pas, mais que `[request].html` existait
- Correction des chemins de génération des icônes par [@robin-samuel](https://github.com/robin-samuel) dans [#4125](https://github.com/wailsapp/wails/pull/4125)
- Correction des événements `fullscreen`, `unfullscreen`, `unminimise` et `unmaximise` qui n’étaient pas émis, par [@oSethoum](https://github.com/osethoum) dans [#4130](https://github.com/wailsapp/wails/pull/4130)
- Correction d’une erreur NSIS causée par un préfixe incorrect sur la version par défaut dans la configuration, par [@robin-samuel](https://github.com/robin-samuel) dans [#4126](https://github.com/wailsapp/wails/pull/4126)
- Correction de la fonction runtime Dialogs qui renvoyait des chemins échappés sous Windows, par [TheGB0077](https://github.com/TheGB0077) dans [#4188](https://github.com/wailsapp/wails/pull/4188)
- Correction du chemin de détection de Webview2 dans HKCU par [@leaanthony](https://github.com/leaanthony).
- Correction d’un problème de saisie sous macOS par [@leaanthony](https://github.com/leaanthony).
- Correction du nom du fichier de tâches de génération des icônes Windows par [@yulesxoxo](https://github.com/yulesxoxo) dans [#4219](https://github.com/wailsapp/wails/pull/4219).
- Correction d’un problème de transparence des fenêtres sans cadre par [@leaanthony](https://github.com/leaanthony), à partir du travail de @kron.
- Correction des appels de mise au premier plan lorsque la fenêtre est désactivée ou réduite, par [@leaanthony](https://github.com/leaanthony), à partir du travail de @kron.
- Correction des icônes de zone de notification qui ne réapparaissaient pas après le redémarrage de la barre des tâches, par [@leaanthony](https://github.com/leaanthony), à partir du travail de @kron.
- Correction de fallbackResponseWriter qui n’implémentait pas Flush() dans [#4245](https://github.com/wailsapp/wails/pull/4245)
- Correction de fallbackResponseWriter qui n’implémentait pas Flush(), par [@superDingda], dans [#4236](https://github.com/wailsapp/wails/issues/4236)
- Correction des plantages survenant à la fermeture d’une fenêtre macOS lorsqu’un appel asynchrone en attente vers une fonction Go liée était en cours, par [@joshhardy](https://github.com/joshhardy), dans [#4354](https://github.com/wailsapp/wails/pull/4354)
- Correction d’une condition de concurrence au démarrage du mode Efficacité de Windows, par [@leaanthony](https://github.com/leaanthony)
- Correction du nettoyage des handles d’icône sous Windows, par [@leaanthony](https://github.com/leaanthony).
- Correction de `OpenFileManager` sous Windows, par [@PPTGamer](https://github.com/PPTGamer), dans [#4375](https://github.com/wailsapp/wails/pull/4375).
- Correction des options de largeur minimale et maximale sous Linux, par @atterpac, dans [#3979](https://github.com/wailsapp/wails/pull/3979)
- Ajout des définitions de types des modèles TypeScript via une mise à jour de la version npm, par @atterpac, dans [#3966](https://github.com/wailsapp/wails/pull/3966)
- Correction de la référence CSS du modèle SvelteKit, par @atterpac, dans [#3945](https://github.com/wailsapp/wails/pull/3945)
- Garantie que les rappels principaux de run() de la fenêtre sont appelés sur le thread principal, par [@leaanthony](https://github.com/leaanthony)
- Correction des exemples de boîte de dialogue de sélection de répertoire, par [@leaanthony](https://github.com/leaanthony)
- Création d’une nouvelle page d’erreur en chinois lorsque index.html est absent, par [@leaanthony](https://github.com/leaanthony)
-  Garantie que le rappel `windowDidBecomeKey` s’exécute sur le thread principal, par [@leaanthony](https://github.com/leaanthony)
-  Prise en charge du plein écran pour les fenêtres sans cadre, par [@leaanthony](https://github.com/leaanthony)
-  Amélioration de la logique de destruction des fenêtres, par [@leaanthony](https://github.com/leaanthony)
-  Correction de la logique de positionnement des fenêtres associées à des icônes de la zone de notification, par [@leaanthony](https://github.com/leaanthony)
-  Prise en charge du plein écran pour les fenêtres sans cadre, par [@leaanthony](https://github.com/leaanthony)
- Correction de la gestion des événements, par [@leaanthony](https://github.com/leaanthony)
- Correction de la logique de fermeture des fenêtres, par [@leaanthony](https://github.com/leaanthony)
- Le Taskfile commun génère désormais par défaut des liaisons TypeScript pour les modèles TypeScript, par [@leaanthony](https://github.com/leaanthony)
- Correction de la fermeture de l’application lors de la réception du message WM_CLOSE lorsqu’aucune fenêtre n’est ouverte ou que seule la zone de notification est utilisée, par [@mmalcek](https://github.com/mmalcek), dans [#3990](https://github.com/wailsapp/wails/pull/3990)
- Correction de la compilation avec garble, par @5aaee9, dans [#3192](https://github.com/wailsapp/wails/pull/3192)
- Correction des compilations NSIS sous Windows, par [@leaanthony](https://github.com/leaanthony)
- Correction d’un interblocage dans la boîte de dialogue Linux de sélection multiple, causé par un élément non fermé
- Correction du nettoyage multiplateforme des fichiers .syso pendant la compilation Windows, par
- Correction de la compilation de l’AppImage amd64, par @atterpac, dans
- Correction de la mise à jour des ressources de compilation, par @ansxuman, dans
- Correction de l’implémentation de `OnClick` et `OnRightClick` pour la zone de notification sous Linux, par @atterpac
- Correction de `AlwaysOnTop` qui ne fonctionnait pas sur Mac, par
-  Correction de `application.NewEditMenu` qui incluait un doublon
- 🐧 Correction de la compilation pour aarch64
- ⊞ Correction des éléments de menu de groupe de boutons radio, par
- Correction de l’erreur lors de la création d’une application .app exécutable sous macOS lorsque 'name' et 'outputfilename'
- Correction d’un bogue lié à l’utilisation de customEventProcessor dans l’exemple de glisser-déposer, par
- 🐧 Correction de l’erreur de compilation sous Linux introduite par l’ajout d’IgnoreMouseEvents, par
- ⊞ Correction d’un bogue de génération du fichier d’icône syso, par
- 🐧 Intégration d’une correction permettant une exécution native sous Wayland, provenant de
- Ne pas lier les méthodes de service internes dans
- ⊞ Correction d’une panique au démarrage de la zone de notification dans
- Ne pas lier les méthodes de service internes dans
- ⊞ Correction d’une panique au démarrage de la zone de notification dans
- Refactorisation majeure des éléments de menu et de la gestion des événements. Pour l’instant, elle améliore principalement macOS. Par
- Correction des tests après la refactorisation des plugins et des événements dans
- ⊞ Correction de l’avertissement `Failed to unregister class Chrome_WidgetWin_0`. Par
- Problèmes de modules
- Correction de la transmission des événements de redimensionnement, par [atterpac](https://github.com/atterpac), dans
- 🐧 Correction d’une erreur de gestion du thème sous NixOS, par
- Correction de l’installation d’un projet entre différents volumes sous Windows, par
- Correction du CSS du modèle React afin d’afficher le pied de page, par
- Correction des processus zombies en mode développement grâce à la mise à jour vers la dernière version de refresh
- Correction de la récupération du fichier WebKit de l’AppImage, par [Atterpac](https://github.com/atterpac)
- Correction de la vérification des paquets apt par Doctor, par [Atterpac](https://github.com/Atterpac), dans
- Correction du blocage de l’application lors de sa fermeture sous Darwin, par @5aaee9, dans
- Correction des couleurs d’arrière-plan des exemples sous Windows, par
- Correction des menus contextuels par défaut, par [mmghv](https://github.com/mmghv), dans
- Correction des valeurs hexadécimales des touches fléchées sous Darwin, par
- Activation du glisser-déposer sous Windows. Ajouté par
- Correction d’un bogue sous Linux dans Doctor lorsque l’utilisateur ne dispose pas des pilotes appropriés
- Correction de la mise à l’échelle DPI au démarrage sous Windows. Modification par [@almas-x](https://github.com/almas-x), dans
- Correction de la ligne de remplacement dans `go.mod` afin d’utiliser des chemins relatifs. Corrige les chemins Windows contenant
- Correction de la gestion des clics sur l’icône de la zone de notification sous macOS lorsqu’aucune fenêtre n’y est associée, par
- Correction de l’échec de la compilation Windows causé par une option inconnue, par
- Correction du plantage sous Windows lors d’un clic gauche sur l’icône de la zone de notification en l’absence d’un
- Correction d’une baseURL incorrecte lors de la deuxième ouverture d’une fenêtre, par @5aaee9, dans la PR
- Correction de l’ordre des branches if dans la méthode `WebviewWindow.Restore`, par
- Calcul correct de `startURL` sur plusieurs appels de `GetStartURL` lorsque
- Correction du type JS de la structure `Screen` afin qu’il corresponde à son équivalent Go, par
- Correction de la méthode `WML.Reload` afin de garantir le nettoyage correct des événements enregistrés
- Correction de la fermeture immédiate du menu contextuel personnalisé sous Linux, par
- Correction du chemin de sortie et de l’extension des fichiers de modèles produits par la liaison
- Correction des chemins d’importation des fichiers de modèles dans le code JS produit par la liaison
- Correction du glisser-déposer sur certaines distributions Linux, par
- Correction de la tâche manquante sous macOS lors de l’utilisation de `wails3 task dev`, par
- Correction de l’enregistrement d’événements qui provoquait une affectation dans une map nil, par
- Correction de la désérialisation des paramètres des méthodes liées, par
- Corriger la gestion de plusieurs valeurs de retour des méthodes liées par
- Corriger la détection par doctor de npm lorsqu’il n’est pas installé avec le gestionnaire de paquets système
- Corriger l’absence de MicrosoftEdgeWebview2Setup.exe. Merci à
- Corriger le plantage aléatoire sous Linux dû à la gestion des identifiants de fenêtre, par @leaanthony. Basé sur
- Corriger le plantage de systemTray.setIcon sous Linux, par
- Faire en sorte que le cadre de la fenêtre soit appliqué dès le premier appel de la fonction `setFrameless` sur

### Modifications

- **RUPTURE DE COMPATIBILITÉ** : les clés de map dans les liaisons JS/TS générées sont désormais marquées comme facultatives afin de refléter fidèlement la sémantique des maps Go. En TypeScript, l’accès aux valeurs d’une map renvoie désormais `T | undefined` au lieu de `T`, ce qui nécessite des vérifications de valeur nulle ou des assertions (#4943), par `@fbbdev`
- Remplacer l’utilisation de `Event` par `Events` conformément aux modifications apportées à `@wailsio/runtime`, et adapter les appels de fonction correspondants dans la documentation de `Features/Events/Event System`, par @AbdelhadiSeddar
- Déplacer `EnabledFeatures`, `DisabledFeatures` et `AdditionalBrowserArgs` des options propres à chaque fenêtre vers `Options.Windows` au niveau de l’application (#4559), par @leaanthony
- Mettre à jour le README de l’exemple `Drag N Drop` et souligner que cet exemple illustre `Internal Drag and Drop`, par @ndianabasi
- Faire passer divers journaux de débogage du niveau Info au niveau Debug (par @mbaklor)
- **RUPTURE DE COMPATIBILITÉ :** renommer `EnableDragAndDrop` en `EnableFileDrop` dans les options de fenêtre
- **RUPTURE DE COMPATIBILITÉ :** renommer `DropZoneDetails` en `DropTargetDetails` dans le contexte d’événement
- **RUPTURE DE COMPATIBILITÉ :** renommer la méthode `DropZoneDetails()` en `DropTargetDetails()` sur `WindowEventContext`
- **RUPTURE DE COMPATIBILITÉ :** supprimer l’événement `WindowDropZoneFilesDropped` et utiliser `WindowFilesDropped` à la place
- **RUPTURE DE COMPATIBILITÉ :** remplacer l’attribut HTML `data-wails-dropzone` par `data-file-drop-target`
- **RUPTURE DE COMPATIBILITÉ :** remplacer la classe CSS de survol `wails-dropzone-hover` par `file-drop-target-active`
- **RUPTURE DE COMPATIBILITÉ :** supprimer les options `DragEffect`, `OnEnterEffect` et `OnOverEffect` de Windows (elles faisaient partie de l’interface IDropTarget supprimée)
- Utiliser goccy/go-json pour tous les traitements JSON à l’exécution (liaisons de méthodes, événements, requêtes de webview, notifications et kvstore), afin d’améliorer les performances de 21-63 % et de réduire les allocations mémoire de 40-60 %
- Optimiser l’organisation de la structure BoundMethod et mettre en cache l’indicateur isVariadic afin de réduire le surcoût de chaque appel
- Utiliser un tampon d’arguments alloué sur la pile pour les méthodes comportant `<=8` arguments afin d’éviter les allocations sur le tas
- Optimiser la collecte des résultats des appels de méthode afin d’éviter l’allocation d’une slice pour une seule valeur de retour
- Utiliser sync.Map pour le cache des types MIME afin d’améliorer les performances en accès concurrent
- Utiliser un pool de tampons pour lire le corps des requêtes du transport HTTP
- Allouer paresseusement le canal CloseNotify dans le détecteur de type de contenu afin de réduire les allocations par requête
- Supprimer la journalisation CSS de débogage du serveur de ressources
- Étendre la table des extensions de types MIME afin de couvrir plus de 50 formats web courants (polices, audio, vidéo, etc.)
- Mettre à jour la documentation des options `X/Y` de Window, par @ruhuang2001
- Mettre à jour la documentation de `Frontend Runtime` en ajoutant davantage d’options de génération des liaisons frontend, par @ndianabasi
- Mettre à jour la page de documentation du serveur de ressources v3 de Wails, par @ndianabasi
- **RUPTURE DE COMPATIBILITÉ** : supprimer les fonctions de boîte de dialogue au niveau du paquet (`application.InfoDialog()`, `application.QuestionDialog()`, etc.). Utiliser à la place le gestionnaire `app.Dialog` : `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Mettre la documentation des boîtes de dialogue en conformité avec l’API réelle : utiliser `app.Dialog.*`, `AddButton()` avec des fonctions de rappel (et non `SetButtons()`), `SetDefaultButton(*Button)` (et non une chaîne), `AddFilter()` (et non `SetFilters()`), `SetFilename()` (et non `SetDefaultFilename()`) et `app.Dialog.OpenFile().CanChooseDirectories(true)` pour sélectionner un dossier
- **RUPTURE DE COMPATIBILITÉ** : les builds de production sont désormais utilisés par défaut. Pour créer des builds de développement, définissez `DEV=true` dans vos Taskfiles. Générer un nouveau projet pour les exemples, par @leaanthony
- Lors de l’émission d’un événement personnalisé avec zéro ou un argument de données, la valeur des données est affectée directement au champ Data sans être encapsulée dans une slice, par [@fbbdev](https://github.com/fbbdev) dans [#4633](https://github.com/wailsapp/wails/pull/4633)
- Sous Windows, les icônes de zone de notification respectent désormais `SystemTray.Show()`/`Hide()` en basculant `NIS_HIDDEN`, ce qui permet aux applications de disparaître réellement, puis de réapparaître (#4653).
- L’enregistrement des icônes de zone de notification réutilise les icônes résolues, définit `NOTIFYICON_VERSION_4` une seule fois et active `NIF_SHOWTIP` afin que les infobulles soient rétablies après le redémarrage de l’Explorateur (#4653).
- macOS : utiliser `visibleFrame` au lieu de `frame` pour centrer les fenêtres en excluant les zones de la barre des menus et du Dock
- macOS : utiliser `visibleFrame` au lieu de `frame` pour centrer les fenêtres en excluant les zones de la barre des menus et du Dock
- Lors de l’exécution de `wails3 update build-assets` avec le paramètre `-config`, les valeurs définies au moyen des paramètres `-product*` sont
- `window.NativeWindowHandle()` -> `window.NativeWindow()`, par @leaanthony dans [#4471](https://github.com/wailsapp/wails/pull/4471)
- Remanier la gestion interne des fenêtres, par @leaanthony dans [#4471](https://github.com/wailsapp/wails/pull/4471)
- Suppression de `application.WindowIDKey` et `application.WindowNameKey` (remplacés par `application.WindowKey`), par [@leaanthony](https://github.com/leaanthony)
- ContextMenuData renvoie désormais une chaîne au lieu de any, par [@leaanthony](https://github.com/leaanthony)
- Dans les liaisons JS/TS, les champs de classe dont le type est un tableau de longueur fixe sont désormais initialisés à la longueur attendue au lieu d’être vides, par [@fbbdev](https://github.com/fbbdev) dans [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData renvoie désormais une chaîne au lieu de any, par [@leaanthony](https://github.com/leaanthony)
- `application.NewService` n’accepte plus les options comme paramètre facultatif (utilisez `application.NewServiceWithOptions` à la place), par [@leaanthony](https://github.com/leaanthony) dans [#4024](https://github.com/wailsapp/wails/pull/4024)
- Suppression de la dépendance `nanoid`, par [@leaanthony](https://github.com/leaanthony)
- Mise à jour de l’exemple Window pour les styles de fenêtre mica, acrylique et à onglets, par [@leaanthony](https://github.com/leaanthony)
- Dans les liaisons JS/TS, les fichiers de modèles `internal.js/ts` ont été supprimés ; tous les modèles se trouvent désormais dans `models.js/ts`, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dans les liaisons JS/TS, les types nommés ne sont jamais rendus comme des alias d’autres types nommés ; l’ancien comportement est désormais réservé aux alias, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dans les liaisons JS/TS en mode classe, les champs de structure dont le type est un paramètre de type sont marqués comme facultatifs et ne sont jamais initialisés automatiquement, par [@fbbdev](https://github.com/fbbdev) dans [#4045](https://github.com/wailsapp/wails/pull/4045)
- Supprimer ESLint des modèles, par [@IanVS](https://github.com/IanVS) dans [#4059](https://github.com/wailsapp/wails/pull/4059)
- Mettre à jour l’année de copyright vers 2025, par [@IanVS](https://github.com/IanVS) dans [#4037](https://github.com/wailsapp/wails/pull/4037)
- Ajout de la documentation de event.Sender par [@IanVS](https://github.com/IanVS) dans [#4075](https://github.com/wailsapp/wails/pull/4075)
- Prise en charge de Go 1.24 par [@leaanthony](https://github.com/leaanthony)
- Les hooks `ServiceStartup` sont désormais invoqués lors de l’appel à `App.Run`, et non dans `application.New`, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Les erreurs `ServiceStartup` sont désormais renvoyées par `App.Run` au lieu de mettre fin au processus, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Les appels de liaisons et de boîtes de dialogue depuis JS sont désormais rejetés avec des objets d’erreur plutôt qu’avec des chaînes, par [@fbbdev](https://github.com/fbbdev) dans [#4066](https://github.com/wailsapp/wails/pull/4066)
- Amélioration du positionnement du menu de la zone de notification sous Windows par [@leaanthony](https://github.com/leaanthony)
- Le runtime JS a été porté vers TypeScript par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Le runtime s’initialise dès son importation, sans qu’il soit nécessaire d’attendre le chargement de la fenêtre, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Le runtime n’exporte plus de méthode init. Une importation pour ses effets de bord permet de l’initialiser, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Les méthodes liées renvoient désormais une `CancellablePromise` qui est rejetée avec une `CancelError` en cas d’annulation. Le résultat réel de l’appel est ignoré, par [@fbbdev](https://github.com/fbbdev) dans [#4100](https://github.com/wailsapp/wails/pull/4100)
- Les types de services intégrés sont désormais systématiquement nommés `Service`, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Les fonctions de création de services intégrés avec options sont désormais systématiquement nommées `NewWithConfig`, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- La méthode `Select` du service `sqlite` se nomme désormais `Query` par souci de cohérence avec les API Go, par [@fbbdev](https://github.com/fbbdev) dans [#4067](https://github.com/wailsapp/wails/pull/4067)
- Modèles : déplacement du runtime vers « dependencies » et réorganisation des fichiers package.json par [@IanVS](https://github.com/IanVS) dans [#4133](https://github.com/wailsapp/wails/pull/4133)
- Création et signature ad hoc des bundles d’application en développement afin d’activer certaines API macOS, par [@popaprozac](https://github.com/popaprozac) dans [#4171](https://github.com/wailsapp/wails/pull/4171)
- Déplacement des ressources de build vers des répertoires propres à chaque plateforme par [@leaanthony](https://github.com/leaanthony)
- Déplacement et renommage des Taskfiles vers des répertoires propres à chaque plateforme par [@leaanthony](https://github.com/leaanthony)
- Expérience nettement améliorée lorsque `index.html` est absent, par [@leaanthony](https://github.com/leaanthony)
- [Windows] Amélioration des performances de réduction et de restauration par [@leaanthony](https://github.com/leaanthony). Basé sur la [PR](https://github.com/wailsapp/wails/pull/3955) d’origine de [562589540](https://github.com/562589540)
- Suppression de l’option `ShouldClose` (enregistrez plutôt un hook pour events.Common.WindowClosing) par [@leaanthony](https://github.com/leaanthony)
- [Windows] Réduction du scintillement à l’ouverture d’une fenêtre par [@leaanthony](https://github.com/leaanthony)
- Suppression de `Window.Destroy`, car cette fonction était destinée à un usage interne, par [@leaanthony](https://github.com/leaanthony)
- Renommage des événements `WindowClose` en `WindowClosing` par [@leaanthony](https://github.com/leaanthony)
- Les builds du frontend utilisent désormais l’environnement vite « development » ou « production » selon le type de build, par [@leaanthony](https://github.com/leaanthony)
- Mise à jour vers go-webview2 v1.19 par [@leaanthony](https://github.com/leaanthony)
- Garantie de l’utilisation du fork de taskfile par @leaanthony
- Mise à jour du fork de Taskfile pour corriger les problèmes de version lors de l’installation avec
- Utilisation du fork de Taskfile pour corriger les problèmes de version lors de l’installation avec
- `service.OnStartup` arrête désormais l’application en cas d’erreur et exécute
- Refactorisation des messages de clic de la zone de notification afin de mieux correspondre aux interactions utilisateur, par
- Intégration des ressources incluant `all:frontend/dist` pour prendre en charge les frameworks qui génèrent
- Refactorisation du Taskfile par [leaanthony](https://github.com/leaanthony) dans
- Mise à niveau vers `go-webview2` v1.0.16 par
- Correction du type `Screen` pour inclure `ID` et non `Id`, par
- Mise à jour de la version Wails de `go.mod.tmpl` pour prendre en charge `application.ServiceOptions`, par
- Correction de la détermination du nom de service par [windom](https://github.com/windom/) dans
- mkdocs serve utilise désormais docker, par [leaanthony](https://github.com/leaanthony)
- Regroupement de la configuration de développement dans `config.yml` par
- La boîte de dialogue de la zone de notification utilise désormais par défaut l’icône de l’application si elle est disponible (Windows), par
- Amélioration des informations sur le GPU et la mémoire sous macOS, par
- Suppression de `WebviewGpuIsDisabled` et de `EnableFraudulentWebsiteWarnings`
- Modification de l’API Events : `On`/`Emit` -> événements utilisateur, `OnApplicationEvent` ->
- Correction de l’API Events sous Linux par [TheGB0077](https://github.com/TheGB0077) dans
- [CI] Amélioration des actions et possibilité de les exécuter également dans les forks et
- Renommage de `AbsolutePosition()` en `Position()` par
- Mise à jour de la dépendance WebKit sous Linux vers webkit2gtk-4.1 au lieu de webkitgtk2-4.0 afin de
- Le script du runtime JS inclus est désormais un module ESM : les balises script qui l’importent
- Le paquet `@wailsio/runtime` ne publie pas son API sur `window.wails`
- Le module de l’API Window `@wailsio/runtime/src/window` expose désormais le conteneur
- L’API JS des fenêtres a été mise à jour pour correspondre à l’actuelle API Go `WebviewWindow`
- Le générateur de liaisons utilise désormais par défaut des appels par ID. L’option de CLI `-id`
- Nouvelle organisation du code des liaisons : les fichiers de sortie étaient auparavant organisés dans des dossiers
- Le champ de structure `application.Options.Bind` a été renommé en
- Nouvelle syntaxe pour lier les services : les instances de service doivent désormais être enveloppées dans un
- Désactivation de l’indicateur de progression dans un environnement sans terminal ou de CI, par

### Supprimé

- **RUPTURE DE COMPATIBILITÉ** : suppression de `EnabledFeatures`, `DisabledFeatures` et `AdditionalLaunchArgs` des options `WindowsWindow` propres à chaque fenêtre. Utilisez plutôt les options `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` et `Options.Windows.AdditionalBrowserArgs` au niveau de l’application. Ces indicateurs s’appliquent globalement à l’environnement WebView2 partagé (#4559), par @leaanthony
- Suppression de l’implémentation native de `IDropTarget` sous Windows au profit d’une approche basée sur JavaScript (correspond au comportement de v2)
- Suppression de la dépendance github.com/wailsapp/mimetype au profit d’une table d’extensions étendue et de http.DetectContentType de la bibliothèque standard, réduisant la taille du binaire d’environ 1.2 Mo
- Suppression de la dépendance gopkg.in/ini.v1 grâce à l’implémentation d’un analyseur minimal de fichiers .desktop pour l’explorateur de fichiers Linux, soit un gain d’environ 45 Ko
- Suppression de samber/lo du code d’exécution grâce à l’utilisation du paquet slices de la bibliothèque standard de Go 1.21+ et d’utilitaires internes minimaux, soit un gain d’environ 310 Ko
- Suppression des instructions printf de débogage du gestionnaire de schémas d’URL Darwin (#4834)
- **RUPTURE DE COMPATIBILITÉ** : suppression de l’événement `linux:WindowLoadChanged`. Utilisez désormais `linux:WindowLoadFinished` pour détecter la fin du chargement de la WebView (#3896), par @leaanthony

### Ruptures de compatibilité

- **Refactorisation de l’API des gestionnaires** : réorganisation de l’API de l’application, d’une structure à plat vers des gestionnaires structurés, afin d’améliorer l’organisation du code et la facilité de découverte, par [@leaanthony](https://github.com/leaanthony) dans [#4359](https://github.com/wailsapp/wails/pull/4359)
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
- Méthodes de Service renommées : `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`, par [@leaanthony](https://github.com/leaanthony)
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

+ Correction des dépendances de paquets Linux trop larges et des dépendances RPM obsolètes.

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

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Modifié

`service.OnShutdown`pour tous les services précédents qui ont démarré, par @atterpac dans   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac dans [#3907](https://github.com/wailsapp/wails/pull/3907)   sous-dossiers par @atterpac dans   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) dans   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) dans   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (remplacé par les options `EnabledFeatures` et `DisabledFeatures`) par   [leaanthony](https://github.com/leaanthony)

### Corrigé

variable de canal par @michael-freling dans   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) dans   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   dans [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) dans   [#3841](https://github.com/wailsapp/wails/pull/3841)   rôle `PasteAndMatchStyle` dans le menu d’édition sous Darwin, par   [johnmccabe](https://github.com/johnmccabe) dans   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) dans   [#3854](https://github.com/wailsapp/wails/pull/3854) par   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   sont différents, par @nickisworking dans   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Ajouté

[mmghv](https://github.com/mmghv) dans   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) et   [leaanthony](https://github.com/leaanthony) dans   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Modifié

Événements d’application `OnWindowEvent` -> événements de fenêtre, par   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   branches préfixées par `v3/` ou `v3-`, par   [stendler](https://github.com/stendler) dans   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Corrigé

[etesam913](https://github.com/etesam913) dans   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) dans   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) dans   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) dans   [#3614](https://github.com/wailsapp/wails/pull/3614) par   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) par   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) par   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) par   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) par   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) par   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Corrigé

## v3.0.0-alpha.5 - 2024-07-30

### Ajouté

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   clic sur l’icône par @5aaee9 dans [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) dans   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   dans [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) dans   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) dans   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   basé sur la [PR](https://github.com/wailsapp/wails/pull/2044) de @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   le paquet npm par [@fbbdev](https://github.com/fbbdev) dans   [#3334](https://github.com/wailsapp/wails/pull/3334)   dans [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` par [@abichinger](https://github.com/abichinger) dans   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) dans   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) dans   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) dans   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Corrigé

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) dans   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) dans   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) dans   [#3477](https://github.com/wailsapp/wails/pull/3477)   par [Atterpac](https://github.com/atterpac) dans   [#3320](https://github.com/wailsapp/wails/pull/3320).   dans [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) dans   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) dans   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) dans   [PR](https://github.com/wailsapp/wails/pull/3039)   installé. Ajouté par [@pylotlight](https://github.com/pylotlight) dans   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   espaces — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) dans la PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) dans la PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   fenêtre attachée [tw1nk](https://github.com/tw1nk) dans la PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) dans   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` est présent.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   écouteurs par [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) dans   [#3330](https://github.com/wailsapp/wails/pull/3330)   générateur par [@fbbdev](https://github.com/fbbdev) dans   [#3334](https://github.com/wailsapp/wails/pull/3334)   générateur par [@fbbdev](https://github.com/fbbdev) dans   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) dans   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) dans   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) dans   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) dans   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) dans   [#3431](https://github.com/wailsapp/wails/pull/3431)   par [@pekim](https://github.com/pekim) dans   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) par   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) dans   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows par [@bruxaodev](https://github.com/bruxaodev/) dans   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Modifié

[mmghv](https://github.com/mmghv) dans   [#3611](https://github.com/wailsapp/wails/pull/3611)   prise en charge d’Ubuntu 24.04 LTS par [atterpac](https://github.com/atterpac) dans   [#3461](https://github.com/wailsapp/wails/pull/3461)   doit posséder l’attribut `type="module"`. Par   [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   objet et ne démarre pas le système WML. Cette modification améliore   l’encapsulation. Si nécessaire, le système WML peut être démarré manuellement en appelant   la nouvelle méthode `WML.Enable`. Le script d’exécution JS fourni effectue toujours les deux   opérations automatiquement. Par [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   objet window comme exportation par défaut. Il n’est désormais plus possible d’importer   des méthodes individuelles avec la syntaxe ESM d’importation nommée ou d’espace de noms.   API. Certaines méthodes ont changé de nom ou de prototype, notamment : `Screen`   devient `GetScreen` ; `GetZoomLevel`/`SetZoomLevel` deviennent `GetZoom`/`SetZoom` ;   `GetZoom`, `Width` et `Height` renvoient désormais directement les valeurs au lieu de les encapsuler   dans des objets. Par [@fbbdev](https://github.com/fbbdev) dans   [#3295](https://github.com/wailsapp/wails/pull/3295)   a été supprimé. Utilisez l’option de CLI `-names` pour revenir aux appels par nom.   Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   nommés d’après leur paquet conteneur ; les chemins d’importation Go complets sont désormais utilisés,   y compris le chemin du module. Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   appel à `application.NewService`. Par [@fbbdev](https://github.com/fbbdev) dans   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) dans   [#3574](https://github.com/wailsapp/wails/pull/3574)

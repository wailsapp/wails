---
title: "Wails v3.1.0 : plus de Go, moins d’obstacles à la compilation"
description: "Compilations expérimentales d’applications de bureau sans CGO et Wake : une méthode plus simple pour compiler, empaqueter et signer les applications Wails."
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

La semaine a été bien remplie chez Wails : compilations sans CGO, configuration de projet simplifiée et moins de difficultés pour l’empaquetage et la signature.

Wails v3.1.0 inaugure nos versions fonctionnelles hebdomadaires, qui vous apporteront régulièrement de nouvelles possibilités. Les correctifs de bogues continueront d’être publiés dans des mises à jour correctives nocturnes dès qu’ils seront prêts.

## Correctifs et optimisations de la semaine

Cette semaine, le framework, la CLI, les outils de compilation et la CI ont bénéficié de **100 correctifs, de 21 améliorations des performances et de quatre changements renforçant la sécurité**. Parmi les principales améliorations figurent un positionnement plus fiable des fenêtres sur les configurations à plusieurs densités de pixels, des correctifs pour la zone de notification et le mode sombre sous Linux, un service des ressources plus léger et une consommation de mémoire réduite pendant l’empaquetage MSIX. [Découvrez quelques-uns de ces changements ci-dessous](#quelques-points-forts-de-la-semaine).

## Compilations expérimentales sans CGO

Purego permet à Go d’appeler les bibliothèques système natives sans CGo, la passerelle habituellement utilisée par Go vers le code C. Nous avons commencé à l’étudier dès son apparition, mais la maintenance de deux backends était difficile. Purego ayant gagné en maturité et l’IA générative rendant cette maintenance plus facile à gérer, nous sommes prêts à proposer des backends expérimentaux pour macOS et Linux.

Vous pouvez désormais compiler et effectuer une compilation croisée d’applications de bureau sans compilateur C :

```sh
WAILS_PUREGO=1 wails3 build
```

Après avoir généré les ressources de votre frontend, vous pouvez également compiler directement avec Go :

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

Le backend macOS utilise AppKit et WebKit ; Linux prend en charge GTK4/WebKitGTK 6.0 et GTK3. Les utilisateurs ont toujours besoin des bibliothèques d’exécution natives, tandis que l’empaquetage et la signature ont leurs propres prérequis. CGo reste utilisé par défaut pendant que nous recueillons vos retours.

[Détails de l’implémentation et des plateformes](https://github.com/wailsapp/wails/pull/6521)

## Pourquoi Wake ?

Wails v2 masquait une trop grande partie du processus de compilation. Dans v3, les Taskfiles ont tout rendu accessible, mais la plupart des utilisateurs souhaitent seulement quelques personnalisations sans avoir à maintenir un système de compilation.

Wake regroupe les choix courants dans `wails.yaml` : applications, outils CLI complémentaires, serveurs, cibles et empaquetage. Task reste notre système de compilation principal et demeure disponible lorsque vous avez besoin de sa souplesse.

Wake prend désormais en charge le rechargement à chaud, les compilations incrémentielles, la mise en cache des sorties et l’empaquetage en parallèle. Vous disposez déjà de fichiers binaires ? `--skip-build` les empaquette sans que Go ni les outils frontend soient installés.

Activez-le avec `WAILS_EXP_USE_WAKE=1`.

### Une prise en main guidée

Lorsque Wake est activé, `wails3 init -ui` vous guide dans le choix du backend, de la cible et des paquets, vérifie vos outils et enregistre la configuration du projet.

@image{src="/assets/screenshots/wake/application.png" alt="L’assistant de projet Wake présentant les choix d’application de bureau et de backend, dont l’option expérimentale Pure Go." expand}

Choisissez les formats de paquets et les options de signature pour chaque cible :

@image{src="/assets/screenshots/wake/packages.png" alt="Sélection des paquets Wake pour Linux, présentant les formats AppImage, DEB, RPM, Void et Arch Linux ainsi que l’option permettant de signer les sorties de publication." expand}

La vérification des outils ci-dessous présente une compilation Linux standard :

@image{src="/assets/screenshots/wake/tools.png" alt="Wake vérifiant le compilateur C, Go, npm et pkg-config pour une compilation Linux standard." expand}

### Que puis-je produire depuis ma machine ?

La compilation, l’empaquetage et la signature ont des prérequis différents. Voici ce que vous pouvez faire depuis chaque système hôte :

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| Opération | Hôte Windows | Hôte macOS | Hôte Linux | Prérequis |
| Compiler des applications de bureau avec purego | ✓ | ✓ | ✓ | Go et le backend expérimental purego pour les cibles macOS/Linux |
| Créer des paquets MSIX | ✓ | ✓ | ✓ | Application compilée ; aucun SDK Windows requis |
| Créer des programmes d’installation NSIS | ✓ | ✓ | ✓ | Application compilée et `makensis` |
| Créer des applications macOS et des DMG | ✓ | ✓ | ✓ | Application compilée ; Go se charge de créer le paquet |
| Créer des paquets DEB, RPM, Arch et Void | ✓ | ✓ | ✓ | Application compilée |
| Créer des AppImages | ✓ | ✓ | ✓ | Runtime AppImage de la cible ; arborescence du runtime d’interface graphique Linux pour l’empaquetage croisé |
| Signer des paquets Windows avec une identité exportée | ✓ | ✓ | ✓ | Certificat RSA PFX exporté et clé privée |
| Signer des applications Apple et des DMG avec une identité exportée | ✓ | ✓ | ✓ | Identité RSA PKCS#12 exportée ; consultez la remarque sur la validation ci-dessous |
| Signer des paquets DEB, RPM et Arch | ✓ | ✓ | ✓ | Clé OpenPGP exportée |
| Soumettre des applications à la notarisation Apple et joindre les tickets | ✓ | ✓ | ✓ | Clé d’API d’équipe Apple et application correctement signée |
| Utiliser `codesign` d’Apple ou les profils du trousseau | – | ✓ | – | macOS et l’identité de signature appropriée |
@end

Les compilations CGo nécessitent toujours une chaîne d’outils C. Les hooks personnalisés qui modifient les DMG nécessitent macOS et `hdiutil`.

### Signature et notarisation

Wake prend en charge la signature des paquets Windows, Apple et Linux, y compris les signatures Apple ad hoc. La signature multiplateforme utilise des identités exportées ; les identités stockées dans un magasin de certificats ou un trousseau nécessitent toujours les outils natifs. Les identifiants utilisateur évitent de placer les clés privées dans la configuration du projet. La signature des paquets AppImage et Void n’est pas encore incluse.

La notarisation Apple s’effectue avec Go au moyen d’une clé d’API d’équipe. Nous avons testé depuis Linux la soumission et l’ajout du ticket, et Gatekeeper a accepté le résultat. Cette application a été signée avec [`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html) d’Apple ; il reste à vérifier qu’Apple accepte les signatures Developer ID générées par Wake avec Go.

### Passer à Wake

La migration prend en charge les neuf modèles intégrés. Elle est irréversible : l’activation de `wails.yaml` remplace le routage du Taskfile. Les commandes personnalisées et les substitutions non prises en charge sont rejetées afin d’être examinées.

[Fonctionnalités, prérequis et détails de migration de Wake](https://github.com/wailsapp/wails/pull/6531)

## Quelques points forts de la semaine

- **Les fenêtres s’affichent au bon endroit sur les écrans à densités de pixels différentes.** Les limites et les positions des écrans sont désormais correctement gérées sur les configurations macOS et Linux dont les écrans utilisent des facteurs d’échelle différents. [#6346](https://github.com/wailsapp/wails/pull/6346)
- **La zone de notification et la détection du thème fonctionnent plus fiablement sous Linux.** Les icônes de zone de notification créées après le démarrage de l’application ne sont plus ignorées, les infobulles sont renseignées et la détection du mode sombre lit dès le démarrage le paramètre d’apparence du bureau.
- **Les fonctions de rappel des événements Go conservent leur ordre pour chaque écouteur.** Chaque écouteur s’exécute séquentiellement dans l’ordre d’émission ; différents écouteurs peuvent s’exécuter simultanément, tandis que l’ordre des émissions simultanées n’est pas défini. La liste des hooks d’événement est copiée sous verrou afin d’éviter une course aux données. [#6275](https://github.com/wailsapp/wails/pull/6275), [#6276](https://github.com/wailsapp/wails/pull/6276)
- **Le cache de compilation prend en charge davantage de projets réels.** Les espaces entre guillemets ou échappés n’entraînent plus d’échecs de correspondance dans le cache des sorties Go, et les enregistrements simultanés du cache sont sérialisés afin d’éviter la corruption des fichiers. [#6495](https://github.com/wailsapp/wails/pull/6495), [#6488](https://github.com/wailsapp/wails/pull/6488)
- **Le service des ressources effectue moins de travail inutile.** La détection du type MIME est ignorée pour les extensions de fichiers connues, et lorsque la journalisation de débogage est désactivée, les messages destinés à être supprimés ne sont plus préparés. Ces deux changements réduisent le travail effectué sur le chemin de traitement des requêtes. [#6414](https://github.com/wailsapp/wails/pull/6414), [#6420](https://github.com/wailsapp/wails/pull/6420)
- **L’empaquetage MSIX évite les tampons de la taille des fichiers binaires.** Les exécutables sont écrits dans le paquet sous forme de flux au lieu d’être entièrement chargés en mémoire. Auparavant, la copie d’un fichier représentatif de 16 Mio allouait environ 16.8 Mo par opération. [#6472](https://github.com/wailsapp/wails/pull/6472)

## Essayez ces fonctionnalités

Ces deux fonctionnalités sont expérimentales. Essayez une compilation multiplateforme, empaquetez votre application et dites-nous où vous rencontrez des difficultés. Merci à toutes les personnes qui testent, signalent des problèmes et contribuent.

Nous remercions tout particulièrement beesophie, membre de notre communauté Discord, qui a créé [Cook](https://github.com/SimTXNG/cook), un remarquable système de compilation qui nous a révélé certaines lacunes du nôtre. Voir ce que Cook rendait possible nous a aidés à déterminer comment simplifier la compilation avec Wails et a véritablement inspiré ce travail.

Merci, Beesophie !

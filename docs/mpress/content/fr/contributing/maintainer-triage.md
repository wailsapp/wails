---
title: "Triage par les mainteneurs"
description: "Trier de manière cohérente les bogues, les signalements concernant la documentation et les WEP"
slug: "contributing/maintainer-triage"
sourcePath: "contributing/maintainer-triage.md"
---

## Objectif

Les issues consignent les bogues reproductibles et les problèmes de documentation. Une pull request de WEP (Wails Enhancement Proposal) consigne une fonctionnalité proposée ou une modification proposée du comportement public.

## Triage des issues

- Vérifiez qu’un rapport de bogue indique une version publiée, la plateforme, les étapes permettant de reproduire le problème, le comportement attendu, le comportement réel et la sortie de `wails3 doctor`.
- Attribuez aux rapports confirmés le libellé `Bug` ainsi que les libellés de version et de plateforme appropriés. Demandez un exemple minimal permettant de reproduire le problème si nécessaire.
- Laissez ouverts les signalements concernant la documentation lorsqu’ils identifient un défaut concret de celle-ci ; encouragez l’auteur du signalement à proposer une PR s’il peut apporter la modification.
- Redirigez les demandes de fonctionnalités vers le guide des WEP, puis fermez-les. Le workflow de redirection automatisé traite les issues nouvellement étiquetées comme des améliorations ; employez la même formulation pour les issues plus anciennes.
- Déplacez les questions et les demandes d’assistance vers GitHub Discussions ou Discord.

## Triage des WEP

1. Vérifiez que la PR est un brouillon intitulé `[WEP] <title>` et qu’elle contient uniquement le WEP et les documents complémentaires.
2. Vérifiez qu’elle utilise le modèle de WEP, désigne une personne chargée de l’implémentation et traite de la compatibilité, des plateformes, des tests, de la maintenance ainsi que de la sécurité et de la confidentialité.
3. Conservez les discussions techniques dans la PR du WEP. Les discussions fournissent un contexte utile, mais ne constituent pas le registre des décisions.
4. Consignez la décision des mainteneurs dans un commentaire de la PR — accepté, rejeté ou retiré — accompagnée d’une brève justification.
5. Pour un WEP accepté, attribuez-lui son numéro, mettez à jour l’index des WEP, fusionnez la PR du WEP et exigez que les PR d’implémentation contiennent un lien vers celui-ci.

## Issues d’amélioration existantes

Ne supprimez pas discrètement les anciennes issues d’amélioration. Pour chaque demande toujours pertinente, ajoutez le commentaire de redirection et fermez-la ; les contributeurs intéressés peuvent ouvrir un WEP. Fermez les doublons en ajoutant un lien vers le WEP ou la décision existants.

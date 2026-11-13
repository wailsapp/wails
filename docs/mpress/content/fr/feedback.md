---
title: "Retours"
description: "Comment faire part de vos retours et signaler des problèmes concernant Wails v3"
slug: "feedback"
sourcePath: "feedback.md"
---

Vos retours sont les bienvenus et vivement encouragés ! Avant de créer un ticket ou une discussion, vérifiez s'il en existe déjà sur le même sujet. Voici les différentes manières de contribuer :

@tabs
[Bugs]
Si vous trouvez un bug, veuillez [ouvrir un ticket](https://github.com/wailsapp/wails/issues/new/choose) sur GitHub à l'aide du modèle de rapport de bug.

- Décrivez clairement le bug à l'aide d'un exemple simple et reproductible. Si la documentation n'indique pas clairement ce qui *devrait* se produire, mentionnez-le dans le rapport.
- Incluez la sortie de `wails3 doctor` dans votre rapport.
- Si le bug correspond à un comportement qui n'est pas conforme à la documentation actuelle, effectuez également les opérations suivantes :
  - Mettez à jour un exemple existant dans le répertoire `v3/examples` ou créez-en un nouveau qui illustre clairement le problème.
  - Ouvrez une [PR](https://github.com/wailsapp/wails/pulls) faisant référence au ticket.


@note{type="caution"}
*Rappelez-vous* qu'un comportement inattendu n'est pas nécessairement un bug : il peut simplement ne pas correspondre à vos attentes. Dans ce cas, utilisez `Suggestions`.

@end

Vous pouvez également discuter des bugs dans le canal [#v3](https://discord.gg/bdj28QNHmT) sur Discord.

[Correctifs]
Si vous disposez d'un correctif pour un bug ou d'une amélioration de la documentation, veuillez :

- Ouvrir une pull request dans le [dépôt Wails](https://github.com/wailsapp/wails) en suivant les [directives de contribution](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md).
- Mentionnez tous les tickets associés dans la description de la PR.

[Améliorations]
Les nouvelles fonctionnalités et les modifications du comportement public doivent être proposées au moyen d'une pull request à l'état de brouillon contenant une **WEP (Wails Enhancement Proposal)**, et non d'un ticket de demande de fonctionnalité.

- Consultez le [processus WEP](https://github.com/wailsapp/wails/blob/v3.0.0/v3/wep/README.md).
- Copiez le modèle, puis ouvrez une PR à l'état de brouillon intitulée `[WEP] <title>`, contenant uniquement la WEP et les documents complémentaires.
- Vous pouvez d'abord discuter de manière informelle d'une idée dans les [discussions GitHub](https://github.com/wailsapp/wails/discussions) ou dans le canal Discord [#v3](https://discord.gg/bdj28QNHmT), mais une PR de WEP est requise pour qu'un mainteneur prenne une décision.

[Votes favorables]
- Manifestez votre soutien aux bugs, aux WEP et aux discussions à l'aide de la réaction :thumbsup: sur GitHub.
- Veuillez *ne pas* vous contenter d'ajouter des commentaires tels que « +1 » ou « moi aussi ».
- Ajoutez toutefois un commentaire si vous avez une contribution pertinente à apporter, par exemple « ce bug affecte également les builds ARM » ou « Une autre approche consisterait à… ».

@end

Vous trouverez [ici](https://github.com/orgs/wailsapp/projects/6) les problèmes connus et les travaux en cours.

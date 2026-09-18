---
title: "Retours"
description: "Comment faire part de vos retours et signaler des problèmes concernant Wails v3"
slug: "feedback"
sourcePath: "feedback.md"
---

Vos retours sont les bienvenus et nous vous encourageons à nous en faire part ! Avant de créer un problème ou une discussion, recherchez ceux qui existent déjà. Voici les différentes manières de contribuer :

@tabs
[Bugs]
Si vous trouvez un bug, veuillez [ouvrir un problème](https://github.com/wailsapp/wails/issues/new/choose) sur GitHub à l’aide du modèle de rapport de bug.

- Décrivez clairement le bug à l’aide d’un exemple simple et reproductible. Si la documentation n’indique pas clairement ce qui *devrait* se produire, mentionnez-le dans le rapport.
- Incluez la sortie de `wails3 doctor` dans votre rapport.
- Si le bug correspond à un comportement qui ne concorde pas avec la documentation actuelle, procédez également comme suit :
  - Mettez à jour un exemple existant dans le répertoire `v3/examples`, ou créez-en un qui illustre clairement le problème.
  - Ouvrez une [PR](https://github.com/wailsapp/wails/pulls) faisant référence au problème.


@note{type="caution"}
*N’oubliez pas* qu’un comportement inattendu n’est pas nécessairement un bug : il se peut simplement que le résultat ne corresponde pas à vos attentes. Dans ce cas, utilisez `Suggestions`.

@end

Vous pouvez également discuter des bugs dans le canal [#v3](https://discord.gg/bdj28QNHmT) sur Discord.

[Correctifs]
Si vous disposez d’un correctif pour un bug ou d’une amélioration de la documentation, veuillez procéder comme suit :

- Ouvrez une pull request dans le [dépôt Wails](https://github.com/wailsapp/wails) en suivant les [consignes de contribution](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md).
- Référencez tous les problèmes associés dans la description de la PR.

[Améliorations]
Les nouvelles fonctionnalités et les modifications du comportement public doivent être proposées au moyen d’une pull request en brouillon contenant une **WEP (proposition d’amélioration de Wails)**, et non d’un problème de demande de fonctionnalité.

- Consultez le [processus WEP](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).
- Copiez le modèle, puis ouvrez une PR en brouillon intitulée `[WEP] <title>` qui ne contient que la WEP et les documents à l’appui.
- Vous pouvez commencer par discuter d’une idée de manière informelle dans les [discussions GitHub](https://github.com/wailsapp/wails/discussions) ou sur le canal Discord [#v3](https://discord.gg/bdj28QNHmT), mais une PR de WEP est requise pour qu’un responsable de maintenance prenne une décision.

[Votes favorables]
- Montrez votre soutien aux bugs, aux WEP et aux discussions à l’aide de la réaction :thumbsup: sur GitHub.
- Veuillez *ne pas* ajouter simplement des commentaires tels que « +1 » ou « moi aussi ».
- Ajoutez toutefois un commentaire si vous avez une contribution pertinente, par exemple « ce bug affecte également les builds ARM » ou « Une autre approche consisterait à… ».

@end

Vous trouverez [ici](https://github.com/orgs/wailsapp/projects/6) les problèmes connus et les travaux en cours.

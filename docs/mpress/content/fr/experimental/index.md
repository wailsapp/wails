---
title: "Fonctionnalités expérimentales"
description: "Un espace consacré aux expérimentations en cours dans Wails v3 : ce qu’elles sont, pourquoi elles existent et comment donner votre avis."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Ici se trouvent des expérimentations"}
Tout ce qui figure dans cette section est, par définition, expérimental. Ces fonctionnalités sont facultatives, désactivées par défaut et peuvent changer de forme, être renommées ou être entièrement supprimées d’une version à l’autre. Ne fondez aucun élément essentiel sur elles sans être prêt à vous adapter à de fréquents changements.

@end

## Ce que signifie « expérimental »

Wails publie ses expérimentations au grand jour. Une expérimentation est une idée qui nous paraît suffisamment prometteuse pour vous être proposée en avant-première, mais sur laquelle nous ne nous sommes pas encore pleinement engagés. Nous la publions *parce que* nous voulons tirer les enseignements d’une utilisation réelle avant de décider si elle doit devenir une composante permanente et prise en charge de Wails.

Cela signifie que tout ce qui figure dans cette section présente les caractéristiques suivantes :

- **Elle est facultative.** Les expérimentations ne modifient jamais le comportement par défaut de `wails3`. Vous les activez délibérément, généralement à l’aide d’une variable d’environnement ou d’une option de compilation. Lorsqu’elles sont désactivées, votre flux de travail existant ne change en rien.
- **Elle pourrait ne pas perdurer.** Certaines expérimentations deviennent des fonctionnalités stables. D’autres sont remaniées au point de devenir méconnaissables ou sont abandonnées. Nous préférons essayer publiquement des idées et en tirer rapidement des enseignements plutôt que de ne publier que ce dont nous sommes déjà certains.
- **L’API n’est pas figée.** Les noms, les options, les valeurs par défaut et le comportement peuvent évoluer d’une version à l’autre pendant qu’une expérimentation prend forme. Les notes de version signaleront ces changements, mais ne vous attendez pas aux garanties de stabilité offertes par les fonctionnalités stables.

## Nous voulons connaître votre avis

C’est le point essentiel. La survie des expérimentations dépend des retours des personnes qui les utilisent réellement. Si vous en essayez une, nous voulons vraiment savoir :

- A-t-elle fonctionné pour votre projet ? Où a-t-elle échoué ?
- Était-elle plus rapide, plus claire ou plus agréable à utiliser, ou le changement n’en valait-il pas la peine ?
- Quelles conditions devraient être réunies pour que vous l’utilisiez par défaut ?

Les retours les plus utiles sont concrets : ce que vous avez exécuté, le résultat attendu et ce qui s’est réellement produit. Chaque expérimentation dispose de son propre fil dans la catégorie **Expérimentations** de GitHub Discussions :

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Discussions sur les expérimentations" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Trouvez le fil consacré à l’expérimentation que vous utilisez et indiquez-nous comment elle s’est comportée, ce qui n’a pas fonctionné ou ce qui manque."}
@end

## Expérimentations actuelles

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Un outil alternatif d’exécution des compilations, adapté à Wails et compatible avec vos Taskfiles existants. Compilations incrémentielles plus rapides, sortie structurée et exécution parallèle par défaut."}
@linkcard{title="Contrôle par LLM (MCP)" href="/guides/mcp-service/" description="Un serveur Model Context Protocol intégré qui permet aux agents LLM d’inspecter, de tester et de piloter une application Wails en cours d’exécution, sans aucun code utilisateur ; il s’active à l’aide d’un tag de compilation."}
@end

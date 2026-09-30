---
title: "Corriger la documentation"
description: "Soumettez une PR de correction pour la documentation de Wails v3 à l’aide de M-Press."
sourcePath: "contributing/documentation.md"
---

Les PR de correction sont les bienvenues. Corrigez les fautes de frappe, les liens rompus, les exemples obsolètes, les explications peu claires ou les traductions. Pour une correction portant uniquement sur la documentation, vous n’avez pas besoin de créer un ticket ni de disposer d’un test de code en échec.

## Prévisualiser localement

Créez un fork de [wailsapp/wails](https://github.com/wailsapp/wails/fork), clonez votre fork, puis créez une branche à partir de `master`.

Installez la version épinglée du générateur de documentation :

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

Vous pouvez également télécharger un binaire vérifié depuis la [version M-Press v1.0.17](https://github.com/leaanthony/mpress/releases/tag/v1.0.17).

Depuis la racine du dépôt Wails :

```sh
mpress version
mpress dev
```

Modifiez les fichiers sources `.md` dans `docs/mpress/content/`. L’anglais est la langue par défaut et se trouve directement dans ce répertoire. Les traductions existantes se trouvent dans des dossiers de langue tels que `fr/` et `id/`. La prévisualisation est régénérée à chaque enregistrement.

Conservez le bloc de métadonnées situé en haut de chaque page ainsi que les composants `@...` / `@end` appariés. Les paragraphes ordinaires, les titres, les listes et le code délimité peuvent être modifiés comme du texte. Ne modifiez pas les fichiers générés dans `docs/mpress/site/`.

## Vérifier la correction

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Vérifiez la page modifiée dans le navigateur et exécutez tous les exemples de code que vous avez modifiés. Pour les corrections de traduction, comparez l’intégralité du passage modifié avec la version anglaise. Conservez intacts les commandes, les noms d’API, les liens, les exemples de code et les connexions des diagrammes. Employez un langage technique naturel, préservez les exigences et les réserves, et traduisez aussi bien le texte que les libellés visibles des diagrammes, les libellés de navigation et les descriptions d’images.

Chaque langue publiée doit proposer une traduction complète de chaque page anglaise. N’utilisez ni textes indicatifs en anglais ni pages de repli. Une correction apportée à une seule traduction peut ne modifier que cette langue. Si vous changez le sens du texte anglais, mettez à jour les pages correspondantes dans les autres langues publiées ; il n’est pas nécessaire de régénérer les pages sans rapport.

## Soumettre une pull request

Ouvrez une PR vers `master`. Décrivez le problème, expliquez votre correction et indiquez les vérifications que vous avez effectuées. Ajoutez des captures d’écran pour les modifications visibles de la mise en page, ainsi que les détails de plateforme et de version pour les exemples de code modifiés.

Les identifiants Cloudflare et les services privés ne sont pas nécessaires. Les vérifications publiques des PR génèrent et valident le site statique sans identifiants de déploiement.

Pour les modifications de code et les propositions de fonctionnalités, consultez [Contribuer à Wails](/contributing/). Pour le fonctionnement interne, consultez la [Présentation technique](/contributing/overview/).

---
title: "Configuration de l’UAC sous Windows"
description: "Configurez le contrôle de compte d’utilisateur (UAC) pour votre application Wails sous Windows"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

Plateformes concernées : <span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

Le contrôle de compte d’utilisateur (UAC) de Windows détermine les privilèges d’exécution de votre application Wails. Par défaut, les applications Wails v3 incluent une configuration UAC explicite dans leur manifeste Windows, ce qui garantit un comportement cohérent sur différentes machines.

## Niveaux d’exécution UAC

Les applications Windows peuvent demander différents niveaux d’exécution dans leur fichier manifeste. Wails v3 inclut automatiquement une configuration UAC avec un niveau d’exécution par défaut que vous pouvez personnaliser selon les besoins de votre application.

### Niveaux d’exécution disponibles

| Niveau | Description | Cas d’utilisation |
| --- | --- | --- |
| `asInvoker` | S’exécute avec les mêmes privilèges que le processus parent | Valeur par défaut pour la plupart des applications |
| `highestAvailable` | S’exécute avec les privilèges les plus élevés dont dispose l’utilisateur | Applications susceptibles de nécessiter des privilèges élevés |
| `requireAdministrator` | Nécessite toujours des privilèges d’administrateur | Utilitaires système, programmes d’installation |

### Configuration par défaut

Les applications Wails v3 incluent une configuration UAC par défaut dans leur manifeste Windows :

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

Cette configuration garantit que votre application :

- S’exécute avec les mêmes privilèges que le processus qui la lance
- Ne nécessite pas d’élévation des privilèges par défaut
- Fonctionne de manière cohérente sur différentes machines
- Ne déclenche pas d’invite UAC pour les utilisateurs standard

## Personnalisation de la configuration UAC

Comme Wails v3 encourage les utilisateurs à personnaliser leurs ressources de compilation, vous pouvez modifier la configuration UAC en éditant directement votre modèle de manifeste Windows.

### Localisation du modèle de manifeste

Le modèle de manifeste Windows se trouve à l’emplacement suivant :

```
build/windows/wails.exe.manifest
```

### Modification du niveau d’exécution

Pour modifier le niveau d’exécution, éditez l’attribut `level` dans l’élément `requestedExecutionLevel` :

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### Exemples

#### Application standard (par défaut)

La plupart des applications devraient utiliser le niveau `asInvoker` par défaut :

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### Utilitaire système

Applications nécessitant des privilèges élevés lorsqu’ils sont disponibles :

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### Outil d’administration

Applications nécessitant toujours des privilèges d’administrateur :

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## Accès à l’interface utilisateur

L’attribut `uiAccess` détermine si votre application peut interagir avec des éléments d’interface utilisateur disposant de privilèges plus élevés. Dans la plupart des cas, il devrait rester défini sur `false`.

Définissez-le sur `true` uniquement si votre application doit :

- Envoyer des entrées à d’autres applications
- Piloter l’interface utilisateur d’autres applications
- Accéder aux éléments d’interface utilisateur de processus disposant de privilèges plus élevés

@note{type="caution" title="Exigences relatives à l’accès à l’interface utilisateur"}
La définition de `uiAccess="true"` exige que votre application soit :

- Signée numériquement avec un certificat émis par une autorité de certification de confiance
- Installée dans un emplacement sécurisé (Program Files ou Windows\System32)

@end

## Compilation avec des paramètres UAC personnalisés

Après avoir modifié votre modèle de manifeste, compilez votre application normalement :

```bash
wails3 build
```

Le processus de compilation intégrera automatiquement votre configuration UAC personnalisée dans l’exécutable.

## Vérification de la configuration UAC

Vous pouvez vérifier que vos paramètres UAC sont correctement intégrés à l’aide de l’outil `go-winres` :

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

Examinez ensuite le fichier manifeste extrait pour confirmer que votre configuration UAC est présente.

@note{type="tip" title="Persistance du manifeste"}
Contrairement à certains autres frameworks, la configuration UAC de Wails v3 est directement intégrée à l’exécutable pendant la compilation, ce qui garantit sa persistance lorsque l’application est copiée sur d’autres machines.

@end

## Dépannage

### Les invites UAC ne s’affichent pas

Si vous définissez `requireAdministrator` mais qu’aucune invite UAC ne s’affiche :

- Vérifiez que le manifeste est correctement intégré à votre exécutable
- Vérifiez que vous ne lancez pas l’application depuis un processus disposant déjà de privilèges élevés
- Assurez-vous que la syntaxe du manifeste est conforme au format XML

### L’application ne démarre pas

Si votre application ne démarre plus après avoir modifié les paramètres UAC :

- Vérifiez que la syntaxe XML du manifeste ne contient aucune erreur
- Vérifiez que la valeur du niveau d’exécution est valide
- Essayez de revenir à `asInvoker` pour isoler le problème

### Comportement incohérent d’une machine à l’autre

Si le comportement de l’UAC diffère d’une machine à l’autre :

- Assurez-vous que le manifeste est incorporé à l’exécutable (et non fourni séparément)
- Vérifiez que l’exécutable n’a pas été modifié après la compilation
- Vérifiez que l’UAC de Windows est activé sur la machine cible

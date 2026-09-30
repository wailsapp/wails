---
title: "Flux — fonctionnement interne"
description: "Fonctionnement du transport des flux, raisons de sa conception, signification des constantes de tampon et éléments non terminés"
slug: "guides/advanced/streams-internals"
sourcePath: "guides/advanced/streams-internals.md"
---

Référence destinée à toute personne — humaine ou agent — qui modifie le transport des flux. L’API destinée aux utilisateurs se trouve dans [Flux](/guides/streams/) ; cette page décrit les mécanismes sous-jacents et les raisons qui les motivent, car plusieurs décisions paraissent arbitraires tant que vous ne savez pas quels problèmes elles évitent.

## Fichiers

| fichier | rôle |
| --- | --- |
| `v3/pkg/application/stream.go` | API publique, `StreamConn`, `streamSink`, le gestionnaire et son registre |
| `v3/pkg/application/stream_session.go` | un chargement de page dans une fenêtre : la file d’attente sortante, les types de trames et la table des connexions |
| `v3/pkg/application/stream_transport.go` | les deux points de terminaison HTTP, le cadrage binaire, le réassemblage des fragments et le préambule du runtime |
| `v3/pkg/application/stream_server.go` | `-tags server` uniquement : véritable récepteur WebSocket |
| `v3/pkg/application/stream_prelude_{server,desktop}.go` | sélectionne le transport client au moment de servir le bundle |
| `v3/internal/runtime/desktop/@wailsio/runtime/src/stream.ts` | le client ayant la forme de `WebSocket` |
| `v3/tests/stream-performance/` | banc de test de charge (`-upload`, `-reloads`, balayages de scénarios) |

## Structure générale

Go→JS et JS→Go utilisent des mécanismes différents, et cette asymétrie constitue toute la conception.

```
Go                                  webview
──                                  ───────
Send() ─► per-window queue ─────────► GET  /wails/stream/poll   (held open)
                                      └─ one held request per window,
                                         carrying frames for every connection

Receive() ◄─ per-conn inbox ◄──────── POST /wails/stream/send   (one or more frames)
```

**Go→JS utilise une requête en attente.** La requête reste suspendue jusqu’à ce qu’un élément soit disponible pour livraison. Il n’existe délibérément aucun intervalle d’interrogation ni aucun mécanisme adaptatif : le serveur attend qu’une trame existe, si bien que la latence de livraison est déjà d’environ 0, et tout intervalle côté client ne pourrait que l’augmenter. Les trames qui arrivent pendant qu’une réponse est en transit s’accumulent et sont acheminées avec la suivante. L’aller-retour lui-même constitue donc la fenêtre de regroupement : celle-ci s’élargit à mesure que la charge augmente, sans qu’aucun mécanisme ne la mesure. Mesures : 1.0 trames par réponse à 100/s, encore 1.0 à 5000/s, 3.4 à 20000/s, et la latence p99 *diminue* lorsque le débit augmente.

**JS→Go utilise des requêtes POST ordinaires.** Les envois sont sérialisés par connexion au moyen d’une chaîne de promesses, car des appels `fetch` simultanés ne préservent pas l’ordre, alors que Go suppose que l’ordre d’envoi correspond à celui qu’il observe. Les trames accumulées derrière une requête en transit sont regroupées dans la requête POST suivante. Go ajoute la trame acceptée ou le préfixe accepté du lot à la boîte de réception de la connexion *avant* de répondre ; le client ne peut donc pas progresser au-delà des octets que Go n’a pas placés en file d’attente.

**Une seule requête d’interrogation en transit par fenêtre, multiplexant toutes les connexions.** Cette propriété garantit l’ordre par construction : une seule file d’attente, un seul mécanisme de vidage et aucun second chemin de livraison susceptible de dépasser le premier. Elle contourne également la limite HTTP/1.1 de six connexions par hôte sous Windows, où il s’agit de véritables requêtes réseau Chromium adressées à `http://wails.localhost`.

## Pourquoi ces décisions précises

Chacune de ces décisions porte la trace du travail sur le transport des événements. En supprimer une réintroduit un bug mesuré.

**Rien dans le chemin Go→JS ne sollicite le thread principal.** `Send` effectue l’ajout sous mutex, puis rend la main. Auparavant, les événements exécutaient leur évaluation directement dans le contexte d’appel lorsqu’ils étaient émis depuis le thread principal, alors qu’une émission antérieure provenant d’une goroutine était encore en file d’attente : 4.4 % des événements étaient inversés, sur les trois plateformes. Une seule file d’attente dotée d’un seul mécanisme de vidage ne peut pas produire ce résultat.

**Rien ne sollicite `evaluateJavaScript`, quelle que soit la taille.** L’insertion de la charge utile dans le code source à évaluer provoque une rétention de mémoire côté hôte au-delà d’un seuil propre à chaque plateforme : 11.6 Go sous macOS et 6.2 Go sous WebKitGTK à 100 × 1 Mo/s. Les flux ne s’en approchent jamais, ce qui explique pourquoi le balayage à débit d’octets constant reste stable pour toutes les tailles de trame.

**Les données de contrôle transitent dans les en-têtes, jamais dans le corps ni dans la chaîne de requête.** Avec les schémas d’URI personnalisés, la version 6.0 de WebKitGTK peut transmettre le corps des requêtes POST sous forme de paramètres de requête (`transport_http.go` comporte précisément une solution de repli pour ce cas), tandis que WebView2 limite la transmission du corps à environ 2 Mo.

**La réponse à la requête d’interrogation est binaire, et non JSON.** Les trames sont `[]byte` ; les encoder en base64 dans une enveloppe JSON ajouterait un surcoût de 33 % à chaque trame, ainsi qu’une analyse syntaxique sur le thread de l’interface utilisateur.

```
magic "WS1\0" | flags u8 | count u32 | count × ( connID u32 | kind u8 | len u32 | payload )
```

`kind` correspond à data / open / close / error. Il n’y a ni numéro de séquence ni accusé de réception : un WebSocket ne rejoue pas les données, et une connexion interrompue perd ce qui était en transit. Reproduire ce comportement est plus simple et plus honnête qu’utiliser un curseur que le tampon borné ne pourrait pas toujours satisfaire.

**Maintenir une requête en attente est sûr**, car chaque requête de la webview dispose déjà de sa propre goroutine. `dispatchWorkers` dans `assetserver_webview.go` est fixé à 0, avec un commentaire qui mentionne précisément ce cas ; l’activation de ce pool nécessiterait d’abord de borner la durée de vie des requêtes.

## Constantes de tampon

Toutes se trouvent dans `stream.go`. **Ce sont des constantes de compilation, et non des options** : il n’existe ni `Options.Streams` ni paramètre propre à chaque flux. Pour les modifier, vous devez éditer le fichier.

| constante | valeur | élément borné |
| --- | ---: | --- |
| `streamOutQueueBytes` | 8 Mo | octets mis en tampon par fenêtre en attente de collecte |
| `streamOutQueueDepth` | 256 | trames mises en tampon par fenêtre |
| `streamOutQueueBytesGlobal` / `streamOutQueueDepthGlobal` | 256 Mo / 8192 | données sortantes mises en tampon dans l’ensemble de l’application |
| `streamInQueueBytesGlobal` / `streamInQueueDepthGlobal` | 256 Mo / 8192 | données entrantes en attente de `Receive` dans l’ensemble de l’application |
| `streamMaxConnections` | 256 | connexions actives et fermetures en file d’attente dans une session |
| `streamMaxConnectionsGlobal` | 4096 | connexions actives dans l’ensemble de l’application |
| `streamOutCloseDepthGlobal` | 4096 | notifications de fermeture non distribuées dans l’ensemble de l’application |
| `streamMaxSessionsPerWindow` | 16 | sessions qu’une fenêtre peut conserver avant qu’une génération plus récente ne doive remplacer une génération antérieure |
| `streamMaxSessions` | 1024 | sessions dans l’ensemble de l’application |
| `streamOutControlDepth` / `streamOutControlDepthGlobal` | 256 / 4096 | trames de contrôle autres que de fermeture mises en file d’attente, par session et dans l’ensemble de l’application |
| `streamMaxChunkSets` / `streamMaxChunkTotal` | 256 / 4096 | téléversements incomplets par session et parties dans un même téléversement |
| `streamMaxChunkBytesGlobal` / `streamMaxChunkPartsGlobal` | 128 Mo / 4096 | charge utile des fragments et métadonnées des parties dans l’ensemble de l’application |
| `streamMaxChunkIDLen` | 64 octets | un identifiant d’ensemble de fragments fourni par le client |
| `streamMaxResponseBytes` | 1 Mo | une réponse d’interrogation |
| `streamHoldTimeout` | 20 s | durée pendant laquelle une interrogation vide reste en attente |
| `streamSessionTTL` | 60 s | aucune interrogation pendant cette durée et aucune connexion active ⇒ la session est morte |
| `streamSessionGrace` | 10 min | aucune interrogation pendant cette durée malgré des connexions actives ⇒ la session est morte |
| `streamSessionSweep` | 20 s | fréquence à laquelle le nettoyeur recherche les sessions mortes |
| `streamMaxFrameBytes` | 64 Mo | une trame dans l’un ou l’autre sens |
| `streamMaxNameLen` | 256 octets | un nom de flux enregistré ou demandé |
| `streamInQueueDepth` / `streamInQueueBytes` | 256 / 8 Mo | trames reçues et pas encore récupérées par `Receive` |

### Comment les choisir

**`streamOutQueueDepth` n’est délibérément pas `eventQueueCapacity` (64).** Cette constante a été mesurée pour une file vidée à raison d’une évaluation à la fois, où une profondeur supérieure n’apportait rien d’autre qu’une latence accrue pour les requêtes les plus lentes. Une interrogation vide les éléments par lots ; la profondeur doit donc couvrir ici la production d’un aller-retour, soit environ 25 trames à 5000 trames/s avec un aller-retour de 5 ms. 256 laisse de la marge pour une rafale sans bloquer le producteur.

**`streamOutQueueBytes` est la limite qui compte réellement**, car 256 trames de 1 Mo représentent 256 Mo. C’est le dernier rempart protégeant la mémoire de l’hôte lorsqu’un frontend cesse de récupérer les données.

Deux règles interagissent ici, et la seconde peut facilement être enfreinte par inadvertance :

- Les limites de profondeur et d’octets bornent l’*accumulation*.
- **Une file vide accepte toujours une trame, quelle que soit sa taille.** L’application inconditionnelle de la limite d’octets empêchait totalement l’envoi d’une trame dépassant cette limite : la condition d’attente ne pouvait jamais devenir vraie, si bien que `Send` restait bloqué indéfiniment et que `TrySend` signalait en permanence que la file était pleine. La taille de la trame ne dépend pas toujours de l’appelant ; une structure comportant un champ `[]byte` est sérialisée avec la taille qui en résulte.

**`streamMaxResponseBytes` existe à cause de Windows.** Le générateur de réponses WebView2 accumule l’intégralité du corps en mémoire et ne le transmet que dans `Finish` ; une réponse sans limite entraîne donc une allocation sans limite sur cette plateforme. Augmenter cette valeur n’améliore *pas* le débit sous Windows : les mesures montrent que le goulot d’étranglement sous Windows dépend du nombre d’octets et non du nombre de réponses. Le nombre de réponses/s varie d’un facteur 4 sur toute la plage de tailles de trame, tandis que le débit en Mo/s reste stable à environ 90.

**C’est la limite entrante qui met le frontend en attente.** Sur le poste de travail, `deliver` signale que la file est pleine, le point de terminaison répond `429` et le client retente la même trame ou le suffixe non accepté du lot avec une temporisation bornée entre les tentatives. Cela évite d’occuper un emplacement de requête de la webview pendant que le gestionnaire rattrape son retard. En mode serveur, la boucle de lecture du socket attend et laisse TCP appliquer la contre-pression. Sans cette limite, un gestionnaire qui tarde à appeler `Receive` pourrait faire croître sans limite la mémoire de l’hôte.

**Les trames de contrôle contournent les limites de données, mais disposent de limites de cycle de vie indépendantes.** La perte d’une trame de données sous l’effet de la contre-pression provoque un ralentissement ; la perte d’un accusé de réception d’ouverture laisse le frontend indéfiniment dans l’état `CONNECTING`, tandis que la perte d’une fermeture lui fait croire qu’une connexion morte est active. Les contrôles autres que les fermetures disposent donc de leur propre file bornée, distincte de celle dans laquelle les fermetures puisent, afin qu’une rafale d’ouvertures refusées ne puisse pas consommer la capacité dont une connexion acceptée a besoin pour signaler sa fin. Chaque session réserve également un emplacement de fermeture par connexion acceptée. Lorsque cette capacité est occupée, une nouvelle ouverture reçoit une contre-pression autorisant une nouvelle tentative avant son enregistrement.

**Les limites par session ont également leurs équivalents à l’échelle de l’application.** Sans elles, chaque session ou connexion admise pourrait conserver simultanément la totalité de son quota local. Les données sortantes et entrantes partagent donc des budgets distincts de 256 Mio / 8192 trames entre les transports de bureau et serveur. Les connexions actives disposent d’un budget de 4096 entrées, et les notifications de fermeture non remises d’un second budget de même taille. Atteindre un quota partagé applique le même comportement bloquant de `Send` ou non bloquant de `TrySend` que lorsqu’un quota local est atteint ; chaque vidage, réception, fermeture, échec d’écriture et chemin d’arrêt restitue sa réservation.

Ces deux budgets sont délibérément distincts, au lieu de former un quota unique qu’une connexion transmettrait à sa propre trame de fermeture. Chaque réservation est libérée par un seul et unique propriétaire : l’emplacement d’une connexion par `shutdown`, qui ne s’exécute qu’une fois, et celui d’une trame de fermeture par l’opération qui élimine cette trame, qu’il s’agisse d’un vidage ou de la destruction de sa session. La propriété qui passe d’une partie à une autre doit être transférée de manière atomique. Dans une révision antérieure, où une fermeture héritait de l’emplacement de la connexion, un emplacement était définitivement perdu chaque fois que la destruction intervenait entre une tentative de fermeture et l’échec de cette tentative.

**Les trames Go transfèrent la propriété ; les trames JavaScript sont copiées instantanément.** La méthode Go `Send` conserve la tranche de l’appelant jusqu’à son écriture par le transport ; après un appel réussi, les appelants ne doivent donc ni modifier ni réutiliser cet espace de stockage. La méthode JavaScript `send()` copie les entrées binaires modifiables avant de rendre la main, conformément à la sémantique de propriété des WebSocket natives. Cette règle asymétrique évite une deuxième copie de la trame entière dans Go, tout en maintenant un comportement intuitif de l’API côté navigateur.

**L’envoi en JavaScript respecte le contrat de mise en tampon des WebSocket.** `send()` ne peut pas bloquer ; une application peut donc mettre des données en file d’attente plus vite que le canal de requêtes du poste de travail ne les accepte, tout comme elle peut dépasser la capacité d’une WebSocket native. `bufferedAmount` inclut chaque octet conservé par ce socket et constitue le signal de contre-pression destiné à l’appelant ; les files côté hôte restent indépendamment bornées par les limites ci-dessus. Un échec définitif, la fermeture par le pair ou un appel local à `close()` libère les charges utiles conservées. La fermeture locale annule également une requête d’ouverture ou de données en attente sur `429` avant d’envoyer le contrôle de fermeture réservé ; ainsi, la contre-pression à l’admission ou côté récepteur ne peut pas laisser le socket bloqué dans l’état `CLOSING`.

**Le réassemblage des fragments dispose d’une allocation partagée de mémoire côté hôte.** Chaque session peut assembler une seule trame d’une taille maximale de 64 Mio, mais cette allocation ne doit pas être multipliée par le nombre de sessions admises. Les ensembles de fragments incomplets et pouvant faire l’objet d’une nouvelle tentative partagent donc un budget de charge utile admise de 128 Mio. L’achèvement d’un ensemble conserve brièvement à la fois ses parties et sa trame assemblée contiguë ; même en doublant cette allocation logique, l’utilisation reste donc dans la limite effective de mémoire de 256 Mio. Les parties conservées partagent également une allocation de métadonnées de 4096 entrées, afin que des fragments minuscules ou vides ne puissent pas faire grossir les tables de correspondance et les données de gestion des tranches sans approcher la limite en octets. Une requête qui dépasserait l’une ou l’autre allocation reçoit un signal de contre-pression autorisant une nouvelle tentative ; la livraison, le rejet, l’expiration ou l’arrêt de la session restitue au budget partagé à la fois les octets et les entrées des parties.

**L’interrogation ne réessaie qu’après les échecs récupérables.** Les erreurs réseau, les réponses d’expiration de requête (`408`), les réponses de données anticipées (`425`), la contre-pression (`429`) et les erreurs serveur (`5xx`) appliquent un délai exponentiel allant de 250 ms à 5 secondes. Les autres réponses `4xx` signalent des échecs de protocole ou de propriété et ferment immédiatement les Streams de la page ; `410` est le signal terminal normal d’une session retirée. La fermeture de la dernière connexion annule une interrogation en cours ou son temporisateur de délai, tandis qu’une connexion ouverte pendant cette phase d’arrêt lance une nouvelle boucle d’interrogation de remplacement.

**`streamSessionTTL` doit rester nettement supérieur à `streamHoldTimeout`**, sinon une session serait supprimée alors que sa propre interrogation est légitimement en attente.

Si vous optimisez pour une charge composée de nombreux petits messages, la limite de profondeur est atteinte en premier ; pour les charges utiles volumineuses, c’est la limite en octets. Il n’est nécessaire de modifier aucune des deux pour une application classique : sous macOS, les valeurs par défaut prennent en charge 634000 trames/s et 2100 Mo/s.

## Cycle de vie des connexions et des sessions

Une **session** correspond au chargement d’une page dans une fenêtre et est identifiée par un identifiant généré par le client, tel que le `clientId` du runtime. Les sessions sont créées de façon différée par la première requête reçue. Lorsqu’une plateforme ne peut pas identifier la fenêtre à l’origine de la requête (`windowID == 0`), le nombre d’identifiants de session reste limité globalement, mais leurs générations ne sont volontairement pas comparées : elles peuvent appartenir à des clients de navigateur indépendants dont les compteurs de génération sont sans rapport. Ces sessions expirent lors de leur fermeture ou à l’échéance de leur TTL, au lieu de se remplacer mutuellement.

Trois mécanismes assurent la fermeture, classés selon leur rapidité à détecter la situation qui la nécessite :

1. **Une interrogation provenant d’une session plus récente d’une fenêtre remplace les générations antérieures.** Lors d’un rechargement, la page reçoit un nouvel identifiant de session et incrémente une génération stockée dans le `sessionStorage` de cette fenêtre. La même valeur est répliquée dans `window.name`, qui persiste après les rechargements lorsque le stockage est désactivé, et s’appuie sur `performance.timeOrigin` (ou `Date.now()` sur les moteurs plus anciens), afin que l’effacement des deux espaces de stockage ne réinitialise pas le compteur à un. Chaque requête transporte l’identifiant et la génération de la session. Une interrogation ne retire que les générations de page inférieures, de sorte que l’ordonnancement du serveur ne puisse pas faire paraître une requête retardée de la page précédente plus récente que celle qui l’a remplacée. Si une stratégie bloque à la fois le stockage et `window.name`, l’ordre repose alors sur l’horloge de la page et dépend donc de l’attribution d’une origine temporelle ultérieure aux pages chargées plus tard. Le gestionnaire conserve, pour chaque fenêtre, un seuil de génération retirée afin qu’une requête déjà en cours ne puisse pas recréer l’ancienne page, sans devoir conserver tous les identifiants de session historiques. Les connexions de la session précédente se ferment immédiatement.
2. **La destruction d’une fenêtre** supprime toutes les sessions de cette fenêtre, comme `eventPayloadStore.dropWindow`.
3. **La suppression à l’échéance du TTL** récupère tout le reste, par exemple un moteur de rendu qui a planté ou une machine en veille.

Seuls les deux premiers mécanismes retirent la génération de la page. Le nettoyage par TTL supprime la session inactive sans avancer le seuil de génération retirée : une page cesse d’interroger dès que sa dernière connexion se ferme, mais cette même page, toujours chargée, doit pouvoir ouvrir un autre flux ultérieurement. Une génération réellement remplacée reste bloquée, car l’interrogation de la nouvelle page avance le seuil avant la suppression de l’ancienne session.

**Les WebViews d’Apple signalent les requêtes annulées.** Sous macOS et iOS, le rappel `stopURLSchemeTask` de WebKit annule le contexte de requête correspondant, de sorte que l’interrogation appartenant à une page quittée par navigation se débloque immédiatement. Le registre utilise comme clé l’identité conservée de la tâche native et supprime les entrées lorsque le traitement des requêtes les ferme. Sous Linux et Windows, le pont actuel n’expose toujours aucun rappel équivalent permettant une annulation anticipée ; une requête en attente y reste bloquée jusqu’à l’expiration du délai de maintien. Grâce à la règle 1, la *connexion* se ferme néanmoins rapidement. Dans les autres cas, l’annulation ne devient visible qu’au niveau de `EPIPE` sous Linux et de `Finish` sous Windows.

## Sélection du transport

`Stream(name)` consulte `window._wails.streamFactory`. Les builds serveur en installent un qui renvoie un véritable `WebSocket` ; dans les builds pour webview, il reste non défini et le client d’interrogation est utilisé.

La fabrique **doit** être installée avant l’exécution du corps de tout module, car les liaisons générées créent des flux au niveau du module. `custom.js` ne peut pas s’en charger : `loadOptionalScript` effectue une requête HEAD, puis ajoute une balise `<script>` ; il intervient donc beaucoup trop tard. La fabrique est plutôt ajoutée au début du bundle du runtime lorsque celui-ci est servi (`stream_prelude_server.go`), ce qui est synchrone par construction : les dépendances des modules ES sont évaluées avant les modules qui les importent.

Si vous ajoutez un troisième transport, placez-le lui aussi dans le préambule. Ne revenez pas à `custom.js`.

## Ce qui n’est pas terminé

|  | état |
| --- | --- |
| Annulation des requêtes depuis la couche de plateforme | **terminée pour Apple ; en attente pour Linux/Windows** — voir ci-dessus |
| Constantes de mémoire tampon sous forme d’options | non terminé ; disponible uniquement à la compilation |
| Flux typés | délibérément non terminé — par décision, les trames sont des `[]byte` |
| Traitement en pipeline (une deuxième interrogation en cours) | non terminé ; nécessiterait un réassemblage ordonné en JS |
| Équité entre les connexions | non terminé — les connexions d’une fenêtre partagent une même file d’attente ; une connexion qui l’inonde ralentit donc ses voisines |
| Regroupement des trames JS→Go | **terminé** — les trames accumulées derrière une requête en cours sont envoyées par lots de taille limitée ; les connexions peu chargées continuent d’envoyer une trame par requête POST |
| Débit sous Windows | environ 100 Mo/s, limité par la sérialisation `WebResourceRequested`. Les mémoires tampons partagées (`PostSharedBufferToScript`) constituent la solution envisagée ; des liaisons existent sous `internal/webview2/pkg/webview2/`, mais ne sont pas intégrées à `pkg/edge` |
| `wails3 dev` / Vite | **fonctionne** — vérifié avec un projet `vanilla-js` généré : le serveur de développement Vite utilise un proxy à `/`, et `/wails/stream/*` est mis en correspondance par le middleware du serveur de ressources avant le proxy ; les flux ne sont donc pas affectés |
| Fenêtres multiples | non testé sous charge, bien que les sessions soient, par construction, limitées à la fenêtre |

## Le paquet frontend en mode développement

Un projet généré importe `@wailsio/runtime` depuis **npm**, et non depuis le `/wails/runtime.js` intégré servi par le serveur de ressources. Sous `wails3 dev`, Vite le résout depuis `node_modules` ; une application construite avec un runtime publié ne verra donc pas les ajouts côté client effectués dans une branche.

Tant que les flux ne sont pas publiés, faites utiliser à une application de test le paquet de cette copie de travail :

```bash
task v3:install-runtime -- ./path/to/your-app/frontend
```

Cette commande reconstruit d’abord `dist/`, afin de toujours installer les sources actuelles. Pour annuler cette modification, exécutez  
`npm install @wailsio/runtime@latest` dans le même répertoire.

Notez qu’il existe **deux** sorties de compilation côté client et qu’il est facile d’en reconstruire une sans reconstruire  
l’autre : `task v3:runtime:build:package` produit le répertoire `dist/` du paquet npm (celui que le frontend d’une application  
importe), tandis que `task v3:runtime:build:assets` produit  
`bundledassets/runtime.js` (celui que la webview charge depuis le serveur de ressources). Toute modification de  
`stream.ts` nécessite de reconstruire les deux.

## Tests

```bash
go test ./pkg/application/ -run TestStream -race        # protocol, ordering, backpressure
go test -tags server ./pkg/application/ -run TestServerMode
pnpm --dir v3/internal/runtime/desktop/@wailsio/runtime test
```

Le test d’ordre est celui qui compte : huit goroutines effectuent des envois simultanés selon un  
compteur attribué sous le verrou de la file d’attente, et l’ordre de vidage doit correspondre exactement à  
l’ordre d’acceptation. Si ce test échoue, même une seule fois, l’invariant imposant un seul consommateur a été rompu.

Banc de test de charge :

```bash
go run ./tests/stream-performance -duration 20s              # full sweep
go run ./tests/stream-performance -upload -duration 10s      # JS→Go matrix
go run ./tests/stream-performance -reloads 6                 # connection lifecycle
```

Sous Windows, il doit s’exécuter dans la session de console interactive : une simple invocation par SSH s’interrompt dans  
la session 0 sans produire aucune sortie. Le binaire doit en outre être placé à un emplacement lisible à la fois par le  
compte SSH et le compte de la console, car les ACL de `C:\Users\<user>` en limitent l’accès à son propriétaire.

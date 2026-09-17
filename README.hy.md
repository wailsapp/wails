<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Սեղանի հավելվածների ստեղծում Go-ի և վեբ տեխնոլոգիաների միջոցով։
</p>

## Բովանդակություն

- [Ներածություն](#ներածություն)
- [Հնարավորություններ](#հնարավորություններ)
  - [Զարգացման ճանապարհային քարտեզ](#զարգացման-ճանապարհային-քարտեզ)
- [Սկսել աշխատանքը](#սկսել-աշխատանքը)
- [Հովանավորներ](#հովանավորներ)
- [Հաճախ տրվող հարցեր](#հաճախ-տրվող-հարցեր)
- [GitHub-ի աստղերը ժամանակի ընթացքում](#github-ի-աստղերը-ժամանակի-ընթացքում)
- [Մասնակիցներ](#մասնակիցներ)
- [Լիցենզիա](#լիցենզիա)
- [Ոգեշնչման աղբյուրներ](#ոգեշնչման-աղբյուրներ)

## Ներածություն

Go ծրագրերին վեբ ինտերֆեյս տրամադրելու ավանդական եղանակը ներկառուցված վեբ սերվերի օգտագործումն է։ **Wails**-ը առաջարկում է այլ մոտեցում․ այն հնարավորություն է տալիս Go կոդը և վեբ Frontend-ը փաթեթավորել **մեկ միասնական գործարկվող ֆայլի** մեջ։

Wails-ը տրամադրում է գործիքներ, որոնք հեշտացնում են նախագծի ստեղծումը, կոմպիլյացիան և փաթեթավորումը։ Մնում է միայն օգտագործել ձեր ստեղծագործական հնարավորությունները։

## Հնարավորություններ

- Backend-ի համար օգտագործել ստանդարտ Go
- Օգտագործել ցանկացած Frontend տեխնոլոգիա, որին արդեն ծանոթ եք
- Նախապես պատրաստված ձևանմուշների միջոցով արագ ստեղծել հարուստ ինտերֆեյսներ Go ծրագրերի համար
- JavaScript-ից հեշտությամբ կանչել Go-ի մեթոդները
- Go-ի Struct-ների և մեթոդների համար TypeScript-ի սահմանումների ավտոմատ ստեղծում
- Native երկխոսությունների և մենյուների աջակցություն
- Native մութ և լուսավոր ռեժիմների աջակցություն
- Ժամանակակից թափանցիկության և «Frosted Glass» պատուհանային էֆեկտների աջակցություն
- Go-ի և JavaScript-ի միջև իրադարձությունների միասնական համակարգ
- Հզոր CLI գործիք՝ նախագծերի արագ ստեղծման և Build-ի համար
- Բազմահարթակ աջակցություն
- Native ռենդերինգային շարժիչների օգտագործում՝ **առանց ներկառուցված բրաուզերի**

### Զարգացման ճանապարհային քարտեզ

Նախագծի ճանապարհային քարտեզը հասանելի է [այստեղ](https://github.com/wailsapp/wails/discussions/1484)։

Նոր հնարավորությունները և հանրային վարքագծի փոփոխությունները առաջարկվում են **WEP (Wails Enhancement Proposal)** մեխանիզմի միջոցով։ Դրանք ներկայացվում են որպես Draft Pull Request, այլ ոչ թե Feature Request Issue։

## Սկսել աշխատանքը

Wails-ն այժմ ունի երկու ակտիվ տարբերակ.

| Տարբերակ | Կարգավիճակ | Տեղադրում | Փաստաթղթեր |
|---|---|---|---|
| **v2** | Կայուն | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` | [wails.io](https://wails.io/) |
| **v3** | Բետա | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` | [v3.wails.io](https://v3.wails.io/) |

Տեղադրման ամբողջական հրահանգները հասանելի են [v2](https://wails.io/docs/gettingstarted/installation) և [v3](https://v3.wails.io) տարբերակների համար։

## Հովանավորներ

Նախագիծը աջակցվում է այն մարդկանց և ընկերությունների կողմից, որոնք նպաստում են դրա զարգացմանը։

<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

## Powered By

[![JetBrains logo.](https://resources.jetbrains.com/storage/products/company/brand/logos/jetbrains.svg)](https://jb.gg/OpenSource)

## Հաճախ տրվող հարցեր

### Արդյո՞ք Wails-ը Electron-ի այլընտրանք է։

Կախված է ձեր պահանջներից։ Wails-ը նախատեսված է Go ծրագրավորողների համար, ովքեր ցանկանում են հեշտությամբ ստեղծել թեթև Desktop հավելվածներ կամ իրենց գոյություն ունեցող Go ծրագրերին ավելացնել Frontend։

Wails-ը տրամադրում է նաև Native տարրեր, ինչպիսիք են մենյուները և երկխոսությունները, ուստի այն կարելի է դիտարկել որպես **Electron-ի թեթև այլընտրանք**։

### Ո՞ւմ համար է նախատեսված այս նախագիծը։

Go ծրագրավորողների համար, ովքեր ցանկանում են HTML/JS/CSS Frontend-ը փաթեթավորել իրենց հավելվածի հետ՝ առանց սերվեր ստեղծելու և ինտերֆեյսը ցուցադրելու համար առանձին բրաուզեր բացելու։

### Որտեղի՞ց է առաջացել Wails անվանումը։

Երբ նախագծի հեղինակը ծանոթացավ WebView-ի գաղափարին, նա մտածեց, որ ցանկանում է ստեղծել WebView հավելվածներ կառուցելու գործիքակազմ՝ մոտավորապես այնպես, ինչպես Rails-ը Ruby-ի համար։

Սկզբնական շրջանում անվանումը բառախաղ էր՝ **Webview on Rails**։ Պատահականորեն այն նաև հնչյունականորեն նման էր այն երկրի անգլերեն անվանը՝ [Wales](https://en.wikipedia.org/wiki/Wales), որտեղից նախագիծը ստեղծողն է։ Այդպես էլ անվանումը պահպանվեց։

## GitHub-ի աստղերը ժամանակի ընթացքում

<a href="https://github.com/wailsapp/wails/stargazers">
  <img alt="Wails-ի աստղերի պատմության գրաֆիկ" src="website/static/img/star-history.svg" width="800" />
</a>

## Մասնակիցներ

Մասնակիցների ցանկն այնքան է մեծացել, որ այն այլևս հնարավոր չէ ամբողջությամբ տեղադրել README-ում։

Բոլոր այն մարդիկ, ովքեր ներդրում են ունեցել նախագծի զարգացման գործում, ներկայացված են [մասնակիցների էջում](https://wails.io/credits#contributors)։

## Լիցենզիա

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Ոգեշնչման աղբյուրներ

Այս նախագիծը հիմնականում մշակվել է հետևյալ ալբոմների ոգեշնչմամբ.

- [Manic Street Preachers - Resistance Is Futile](https://open.spotify.com/album/1R2rsEUqXjIvAbzM0yHrxA)
- [Manic Street Preachers - This Is My Truth, Tell Me Yours](https://open.spotify.com/album/4VzCL9kjhgGQeKCiojK1YN)
- [The Midnight - Endless Summer](https://open.spotify.com/album/4Krg8zvprquh7TVn9OxZn)
- [Gary Numan - Savage (Songs from a Broken World)](https://open.spotify.com/album/3kMfsD07Q32HRWKRrpcexr)
- [Steve Vai - Passion & Warfare](https://open.spotify.com/album/0oL0OhrE2rYVns4IGj8h2m)
- [Ben Howard - Every Kingdom](https://open.spotify.com/album/1nJsbWm3Yy2DW1KIc1OKle)
- [Ben Howard - Noonday Dream](https://open.spotify.com/album/6astw05cTiXEc2OvyByaPs)
- [Adwaith - Melyn](https://open.spotify.com/album/2vBE40Rp60tl7rNqIZjaXM)
- [Gwidaith Hen Fran - Cedors Hen Wrach](https://open.spotify.com/album/3v2hrfNGINPLuDP0YDTOjm)
- [Metallica - Metallica](https://open.spotify.com/album/2Kh43m04B1UkVcpcRa1Zug)
- [Bloc Party - Silent Alarm](https://open.spotify.com/album/6SsIdN05HQg2GwYLfXuzLB)
- [Maxthor - Another World](https://open.spotify.com/album/3tklE2Fgw1hCIUstIwPBJF)
- [Alun Tan Lan - Y Distawrwydd](https://open.spotify.com/album/0c32OywcLpdJCWWMC6vB8v)
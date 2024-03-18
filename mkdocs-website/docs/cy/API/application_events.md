### Ar

API:
`Ar(eventType digwyddiadau.DdigwyddiadweithgangeningApplicationEventType, atebydd func(digwyddiad *Digwyddiad)) func()`

Mae `Ar()` yn cofrestru gwrandäwr digwyddiad ar gyfer digwyddiadau cymhwysiad penodol. Bydd y swyddogaeth atebydd a ddarperir yn cael ei sbarduno pan fydd y digwyddiad cysylltiedig yn digwydd. Mae'r swyddogaeth yn dychwelyd swyddogaeth y gellir ei galw i dynnu'r gwrandäwr.

### CofrestruArgraffwyr

API: 
`CofrestruArgraffwyr(eventType digwyddiadau.DdigwyddiadweithgangeningApplicationEventType, atebydd func(digwyddiad *Digwyddiad)) func()`

Mae `CofrestruArgraffwyr()` yn cofrestru atebydd i'w redeg fel crocen yn ystod digwyddiadau penodol. Caiff y crocenau hyn eu rhedeg cyn gwrandawyr sy'n gysylltiedig ag `Ar()`. Mae'r swyddogaeth yn dychwelyd swyddogaeth y gellir ei galw i dynnu'r bâs.
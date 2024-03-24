## Llusgo a Gollwng

Gellir galluogi llusgo a gollwng brodorol fesul ffenest. Yn syml, gosodwch yr
opsiwn cyflunio ffenest `EnableDragAndDrop` i `true` a bydd y ffenest yn
caniatáu i ffeiliau gael eu llusgoi arni. Pan fydd hyn yn digwydd, bydd y
digwyddiad `events.FilesDropped` yn cael ei yrru. Gellir yna nôl yr enwau
ffeil o `WindowEvent.Context()` gan ddefnyddio'r dull `DroppedFiles()`. Mae
hwn yn dychwelyd sleis o linynnau yn cynnwys yr enwau ffeil.
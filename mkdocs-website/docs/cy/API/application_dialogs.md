### Dangos Deialog Ynghylch

API: `ShowAboutDialog()`

Mae `ShowAboutDialog()` yn dangos blwch deialog "Ynghylch". Gall ddangos enw'r
cymhwysiad, disgrifiad ac eicon.

```go
    // Dangos y deialog ynghylch
    app.ShowAboutDialog()
```

### Gwybodaeth

API: `InfoDialog()`

Mae `InfoDialog()` yn creu ac yn dychwelyd esiampl newydd o `MessageDialog` gyda
`InfoDialogType`. Defnyddir y deialog hon fel arfer i ddangos negeseuon
gwybodaeth i'r defnyddiwr.

### Cwestiwn

API: `QuestionDialog()`

Mae `QuestionDialog()` yn creu ac yn dychwelyd esiampl newydd o `MessageDialog`
gyda `QuestionDialogType`. Defnyddir y deialog hon yn aml i ofyn cwestiwn i'r
defnyddiwr a disgwyl ymateb.

### Rhybudd

API: `WarningDialog()`

Mae `WarningDialog()` yn creu ac yn dychwelyd esiampl newydd o `MessageDialog`
gyda `WarningDialogType`. Fel y mae'r enw yn awgrymu, defnyddir y deialog hon yn
bennaf i ddangos negeseuon rhybudd i'r defnyddiwr.

### Gwall

API: `ErrorDialog()`

Mae `ErrorDialog()` yn creu ac yn dychwelyd esiampl newydd o `MessageDialog` gyda
`ErrorDialogType`. Cynlluniwyd y deialog hon i'w defnyddio pan fydd angen
dangos neges gwall i'r defnyddiwr.

### Agor Ffeil

API: `OpenFileDialog()`

Mae `OpenFileDialog()` yn creu ac yn dychwelyd esiampl newydd o
`OpenFileDialogStruct`. Mae'r deialog hon yn annog y defnyddiwr i ddewis un neu
ragor o ffeiliau o'u system ffeiliau.

### Cadw Ffeil

API: `SaveFileDialog()`

Mae `SaveFileDialog()` yn creu ac yn dychwelyd esiampl newydd o
`SaveFileDialogStruct`. Mae'r deialog hon yn annog y defnyddiwr i ddewis lleoliad
yn eu system ffeiliau lle y dylid cadw ffeil.

### Agor Cyfeiriadur

API: `OpenDirectoryDialog()`

Mae `OpenDirectoryDialog()` yn creu ac yn dychwelyd esiampl newydd o
`MessageDialog` gyda `OpenDirectoryDialogType`. Mae'r deialog hon yn galluogi'r
defnyddiwr i ddewis cyfeiriadur o'u system ffeiliau.
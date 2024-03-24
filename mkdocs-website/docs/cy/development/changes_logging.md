### Cofnodi

Roedd cofnodi yn v2 yn ddryslyd gan fod cofnodion cymhwysiad a chofnodion system (mewnol) yn defnyddio'r un cofnodwr. Rydym wedi ei symleiddio fel a ganlyn:

- Mae cofnodion mewnol yn cael eu trin nawr gan ddefnyddio'r cofnodwr `slog` safonol Go. Caiff hwn ei ffurfweddu gan ddefnyddio'r opsiwn `logger` yn yr opsiynau cymhwysiad. Yn ddiofyn, mae hwn yn defnyddio'r cofnodwr [tint](https://github.com/lmittmann/tint).
- Gellir cyflawni cofnodion cymhwysiad nawr trwy'r ciplug `log` newydd sy'n defnyddio `slog` o dan y rhyngwyneb. Mae'r ciplug hwn yn darparu API syml ar gyfer cofnodi i'r consol. Mae ar gael yn y naill iaith Go a JS.
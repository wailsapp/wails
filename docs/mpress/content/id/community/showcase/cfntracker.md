---
title: "CFN Tracker"
description: "Aplikasi desktop yang dibuat dengan Wails"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker) - Lacak pertandingan langsung dari profil CFN Street Fighter 6 atau V mana pun. Kunjungi [situs web](https://cfn.williamsjokvist.se/) untuk memulai.

## Fitur

- Pelacakan pertandingan secara real-time
- Penyimpanan log dan statistik pertandingan
- Dukungan untuk menampilkan statistik langsung di OBS melalui Browser Source
- Dukungan untuk SF6 dan SFV
- Kemampuan bagi pengguna untuk membuat tema OBS Browser sendiri dengan CSS

### Teknologi utama yang digunakan bersama Wails

- [Task](https://github.com/go-task/task) - membungkus CLI Wails agar perintah umum mudah digunakan
- [React](https://github.com/facebook/react) - dipilih karena ekosistemnya yang kaya (radix, framer-motion)
- [Bun](https://github.com/oven-sh/bun) - digunakan karena resolusi dependensinya yang cepat dan waktu build yang singkat
- [Rod](https://github.com/go-rod/rod) - otomatisasi browser headless untuk autentikasi dan polling perubahan
- [SQLite](https://github.com/mattn/go-sqlite3) - digunakan untuk menyimpan pertandingan, sesi, dan profil
- [Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) - aliran HTTP untuk mengirim pembaruan pelacakan ke browser source OBS
- [i18next](https://github.com/i18next/) - dengan konektor backend untuk menyajikan objek pelokalan dari lapisan Go
- [xstate](https://github.com/statelyai/xstate) - state machine untuk proses autentikasi dan pelacakan

---
title: "Utilitas Ekspor FileHound"
description: "Aplikasi desktop yang dibuat dengan Wails"
slug: "community/showcase/filehound"
sourcePath: "community/showcase/filehound.md"
---

![Utilitas Ekspor FileHound](/assets/showcase-images/filehound.webp)

[Utilitas Ekspor FileHound](https://www.filehound.co.uk/) FileHound adalah platform manajemen dokumen berbasis cloud yang dirancang untuk penyimpanan file yang aman, otomatisasi proses bisnis, dan kemampuan SmartCapture.

Utilitas Ekspor FileHound memungkinkan Administrator FileHound menjalankan tugas ekstraksi dokumen dan data secara aman untuk keperluan pencadangan dan pemulihan alternatif. Aplikasi ini akan mengunduh semua dokumen dan/atau metadata yang disimpan di FileHound berdasarkan filter yang Anda pilih. Metadata tersebut akan diekspor dalam format JSON dan XML.

Backend dibuat dengan:

- Go 1.15
- Wails 1.11.0
- go-sqlite3 1.14.6
- go-linq 3.2

Frontend menggunakan:

- Vue 2.6.11
- Vuex 3.4.0
- TypeScript
- Tailwind 1.9.6

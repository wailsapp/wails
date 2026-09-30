---
title: "Umpan Balik"
description: "Cara memberikan umpan balik dan melaporkan masalah untuk Wails v3"
slug: "feedback"
sourcePath: "feedback.md"
---

Kami menyambut (dan mendorong) umpan balik Anda! Harap cari issue atau diskusi yang sudah ada sebelum membuat yang baru. Berikut berbagai cara untuk berkontribusi:

@tabs
[Bug]
Jika Anda menemukan bug, silakan [buka issue](https://github.com/wailsapp/wails/issues/new/choose) di GitHub menggunakan templat laporan bug.

- Jelaskan bug dengan jelas menggunakan contoh sederhana yang dapat direproduksi. Jika dokumentasi tidak menjelaskan dengan jelas apa yang *seharusnya* terjadi, sertakan hal tersebut dalam laporan.
- Sertakan keluaran dari `wails3 doctor` dalam laporan Anda.
- Jika bug tersebut berupa perilaku yang tidak sesuai dengan dokumentasi saat ini, lakukan juga hal berikut:
  - Perbarui contoh yang sudah ada di direktori `v3/examples`, atau buat contoh baru yang menunjukkan masalah tersebut dengan jelas.
  - Buka [PR](https://github.com/wailsapp/wails/pulls) yang merujuk pada issue tersebut.


@note{type="caution"}
*Ingat*, perilaku yang tidak terduga belum tentu merupakan bug—mungkin perilakunya hanya tidak sesuai dengan yang Anda harapkan. Gunakan `Suggestions` untuk hal tersebut.

@end

Anda juga dipersilakan mendiskusikan bug di kanal [#v3](https://discord.gg/bdj28QNHmT) di Discord.

[Perbaikan]
Jika Anda memiliki perbaikan untuk bug atau peningkatan dokumentasi, silakan:

- Buka pull request di [repositori Wails](https://github.com/wailsapp/wails) dengan mengikuti [panduan kontribusi](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md).
- Cantumkan referensi ke setiap issue terkait dalam deskripsi PR.

[Peningkatan]
Fungsionalitas baru dan perubahan pada perilaku publik diusulkan melalui draft pull request **WEP (Wails Enhancement Proposal)**, bukan melalui issue permintaan fitur.

- Baca [proses WEP](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).
- Salin templat, lalu buka draft PR berjudul `[WEP] <title>` yang hanya berisi WEP dan materi pendukung.
- Anda dapat terlebih dahulu mendiskusikan ide secara informal di [GitHub Discussions](https://github.com/wailsapp/wails/discussions) atau kanal Discord [#v3](https://discord.gg/bdj28QNHmT), tetapi PR WEP wajib diajukan agar maintainer dapat mengambil keputusan.

[Memberikan Suara Dukungan]
- Tunjukkan dukungan untuk bug, WEP, dan diskusi menggunakan reaksi :thumbsup: di GitHub.
- Harap *jangan* sekadar menambahkan komentar seperti "+1" atau "saya juga".
- Tambahkan komentar jika Anda memiliki kontribusi yang bermakna, seperti "bug ini juga memengaruhi build ARM" atau "Pendekatan lain yang dapat digunakan adalah...".

@end

Issue yang telah diketahui dan pekerjaan yang sedang berlangsung dapat ditemukan [di sini](https://github.com/orgs/wailsapp/projects/6).

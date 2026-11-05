---
title: "Fitur Eksperimental"
description: "Tempat untuk eksperimen yang sedang dikembangkan di Wails v3 — apa saja eksperimen tersebut, mengapa dibuat, dan cara memberikan masukan."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Di sinilah eksperimen berada"}
Semua yang ada di bagian ini, sesuai definisinya, merupakan eksperimen. Fitur-fitur di sini harus diaktifkan secara opsional, dinonaktifkan secara default, dan bentuknya dapat berubah, namanya dapat diganti, atau dapat dihapus sepenuhnya antar-rilis. Jangan mengandalkannya untuk bagian penting apa pun kecuali Anda siap menghadapi perubahan yang tidak menentu.

@end

## Arti "Eksperimental"

Wails menyediakan eksperimen secara terbuka. Eksperimen adalah gagasan yang menurut kami cukup menjanjikan untuk diberikan kepada Anda lebih awal — tetapi belum sepenuhnya kami tetapkan. Kami memublikasikannya *karena* kami ingin belajar dari penggunaan nyata sebelum memutuskan apakah eksperimen tersebut akan menjadi bagian permanen Wails yang didukung.

Artinya, beberapa hal berikut berlaku untuk semua yang ada di bagian ini:

- **Fitur ini bersifat opsional.** Eksperimen tidak pernah mengubah perilaku default `wails3`. Anda mengaktifkannya dengan sengaja (biasanya melalui variabel lingkungan atau flag build), dan saat eksperimen dinonaktifkan, alur kerja Anda yang sudah ada tidak berubah sama sekali.
- **Fitur ini mungkin tidak dipertahankan.** Sebagian eksperimen berkembang menjadi fitur stabil. Eksperimen lainnya dirombak hingga tidak lagi menyerupai bentuk awalnya, atau dihentikan. Kami lebih memilih mencoba berbagai hal secara terbuka dan belajar dengan cepat daripada hanya merilis hal yang sudah benar-benar kami yakini.
- **API-nya belum ditetapkan.** Nama, flag, nilai default, dan perilaku dapat berubah antar-rilis selama suatu eksperimen masih dalam proses menemukan bentuk akhirnya. Catatan rilis akan menjelaskan perubahan tersebut, tetapi jangan mengharapkan jaminan stabilitas yang dimiliki fitur stabil.

## Kami menginginkan masukan Anda

Inilah bagian yang penting. Kelanjutan eksperimen bergantung pada masukan dari orang-orang yang benar-benar menggunakannya. Jika Anda mencoba salah satunya, kami benar-benar ingin mengetahui:

- Apakah eksperimen tersebut berfungsi untuk proyek Anda? Di bagian mana eksperimen itu gagal?
- Apakah eksperimen tersebut lebih cepat / lebih jelas / lebih nyaman — atau tidak sepadan dengan upaya untuk beralih?
- Apa yang perlu dipenuhi agar Anda menggunakannya secara default?

Masukan yang paling berguna bersifat konkret: apa yang Anda jalankan, apa yang Anda harapkan, dan apa yang sebenarnya terjadi. Setiap eksperimen memiliki utas tersendiri dalam kategori **Eksperimen** di GitHub Discussions:

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Diskusi Eksperimen" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Temukan utas untuk eksperimen yang Anda gunakan, lalu beri tahu kami bagaimana hasilnya, apa yang rusak, atau apa yang masih kurang."}
@end

## Eksperimen saat ini

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Build runner alternatif yang memahami Wails untuk Taskfile Anda yang sudah ada. Build inkremental yang lebih cepat, output terstruktur, dan berjalan paralel secara default."}
@linkcard{title="Kontrol LLM (MCP)" href="/guides/mcp-service/" description="Server Model Context Protocol bawaan yang memungkinkan agen LLM memeriksa, menguji, dan mengendalikan aplikasi Wails yang sedang berjalan — tanpa memerlukan kode pengguna dan diaktifkan dengan build tag."}
@end

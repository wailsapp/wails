---
title: "Snippet Expander"
description: "Aplikasi desktop yang dibuat dengan Wails"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Tangkapan Layar Snippet Expander](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Tangkapan layar jendela Select Snippet di Snippet Expander

![Tangkapan Layar Snippet Expander](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Tangkapan layar halaman Add Snippet di Snippet Expander

![Tangkapan Layar Snippet Expander](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Tangkapan layar jendela Search & Paste di Snippet Expander

[Snippet Expander](https://snippetexpander.org) adalah "Asisten kecil untuk cuplikan teks yang dapat diperluas" bagi Linux.

Snippet Expander terdiri atas aplikasi GUI yang dibuat dengan Wails untuk mengelola cuplikan dan pengaturan, serta mode jendela Search & Paste untuk memilih dan menempelkan cuplikan dengan cepat.

GUI berbasis Wails, CLI go-lang, dan daemon perluasan otomatis vala-lang semuanya berkomunikasi dengan daemon go-lang melalui D-Bus. Daemon tersebut menangani sebagian besar pekerjaan, yaitu mengelola basis data cuplikan dan pengaturan umum, serta menyediakan layanan untuk memperluas dan menempelkan cuplikan, dan sebagainya.

Lihat [kode sumber](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38) untuk mengetahui cara aplikasi Wails mengirim pesan dari UI ke backend yang kemudian diteruskan ke daemon, serta berlangganan peristiwa D-Bus untuk memantau perubahan cuplikan yang dilakukan melalui instans lain aplikasi atau CLI dan langsung menampilkannya di UI melalui peristiwa Wails.

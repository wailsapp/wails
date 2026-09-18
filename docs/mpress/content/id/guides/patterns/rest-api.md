---
title: "REST API"
description: "Mengekspos RESTful API dari aplikasi Wails v3 Anda"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
Halaman ini adalah penampung sementara. Konten lengkap tentang REST API akan segera tersedia.

@end

Anda dapat mengekspos RESTful API dari aplikasi Wails v3 menggunakan handler Go `net/http` standar atau framework seperti [Gin](/guides/gin-routing/). Karena Wails menyajikan frontend melalui handler aset, handler HTTP apa pun dapat dipasang dan diakses dari UI Anda atau klien lain.

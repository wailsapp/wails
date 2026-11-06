---
title: "Redis Viewer"
description: "GUI Redis desktop yang dibuat dengan Wails"
slug: "community/showcase/redisviewer"
sourcePath: "community/showcase/redisviewer.md"
---

![Tangkapan layar RedisViewer](/assets/showcase-images/redisviewer-overview1.webp)

![Tangkapan layar RedisViewer](/assets/showcase-images/redisviewer-overview2.webp)

[RedisViewer](https://redisviewer.com/) adalah GUI Redis desktop modern yang dibuat dengan Wails. Gunakan aplikasi ini untuk memeriksa nilai yang kompleks, menjalankan perintah, dan menganalisis performa Redis tanpa mengorbankan kualitas interaksi.

Dirancang berdasarkan arsitektur WebView + Go milik Wails, aplikasi ini menangani keyspace besar dan payload berat di backend Go, alih-alih mengirimkan semuanya ke frontend—dengan memanfaatkan model memori Go serta menghindari tekanan heap dan risiko kebocoran memori yang umum terjadi pada klien yang sangat bergantung pada JS. Antarmuka pengguna disesuaikan agar hanya merender konten yang tampil di layar, sehingga penelusuran set data berukuran sangat besar tetap responsif dan terasa menyerupai aplikasi desktop native.

[Kunjungi Situs Web Proyek](https://redisviewer.com/)

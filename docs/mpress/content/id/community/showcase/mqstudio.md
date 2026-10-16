---
title: "MQ Studio"
description: "Klien desktop yang mengutamakan penyimpanan lokal untuk RocketMQ, RabbitMQ, Kafka, dan lainnya"
slug: "community/showcase/mqstudio"
sourcePath: "community/showcase/mqstudio.md"
---

![Dialog koneksi baru MQ Studio yang menampilkan driver RocketMQ, Kafka, dan RabbitMQ](/assets/showcase-images/mqstudio.webp)

**[MQ Studio](https://mq-studio.amigoer.com)** adalah **klien desktop untuk antrean pesan yang mengutamakan penyimpanan lokal**, dibuat dengan **Go, Wails, dan React**. Setiap broker memiliki konsol sendiri: RocketMQ punya satu, Kafka punya yang lain, dan RabbitMQ menyediakan plugin pengelolaan. Antarmuka dan istilahnya berbeda, dan masing-masing merupakan layanan yang harus diterapkan serta dipelihara. MQ Studio menggantikan semuanya dengan satu aplikasi: setiap broker diakses melalui driver di balik antarmuka yang sama, sehingga halaman dan alur kerja tetap sama apa pun sistem yang terhubung.

## Fitur utama

- **Satu antarmuka untuk semua broker** — saat ini mendukung RocketMQ, RabbitMQ, dan Kafka; Pulsar, NATS, MQTT, dan SQS masuk dalam rencana pengembangan.
- **Topik, antrean, dan pesan** — periksa topik, antrean, exchange, dan binding; cari serta lacak pesan, pantau log, buat pesan dengan kunci dan header, kirim ulang, dan tangani pesan yang gagal terkirim.
- **Konsumen dan ketertinggalan** — grup, klien, langganan, dan ketertinggalan per partisi, dengan pengaturan ulang offset serta penanganan percobaan ulang dan antrean pesan gagal (DLQ).
- **Klaster dan peringatan** — kesehatan broker, metrik waktu proses, laju pemrosesan, penggunaan disk, dan notifikasi desktop bawaan.
- **Kemampuan koneksi yang transparan** — setiap driver menyatakan kemampuan endpoint yang sebenarnya, dan antarmuka hanya menawarkan operasi yang didukung broker.
- **Privasi secara bawaan** — konfigurasi tetap berada di perangkat Anda dan kredensial dienkripsi saat disimpan.

Tidak ada komponen server yang harus diterapkan, konsol web yang harus dipelihara, maupun telemetri. Wails memungkinkan hal ini: klien administrasi broker tersebut berupa pustaka Go. Lapisan driver mengaksesnya langsung dalam proses yang sama, sementara antarmukanya tetap berupa aplikasi React biasa — satu berkas biner per platform, bukan layanan yang harus dioperasikan seseorang.

Tersedia untuk macOS, Windows, dan Linux, dalam bahasa Inggris dan Mandarin.

[Situs web](https://mq-studio.amigoer.com) |
[GitHub](https://github.com/amigoer/mq-studio) |
[Unduh](https://github.com/amigoer/mq-studio/releases/latest)

---
title: "Wails v2 Dirilis"
description: "Catatan rilis dan pengumuman untuk Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![tangkapan layar montase](/assets/blog-images/montage.png)

## Akhirnya hadir!

Hari ini menandai perilisan [Wails](https://wails.io) v2. Sudah sekitar 18 bulan sejak versi alfa pertama v2 dan sekitar satu tahun sejak rilis beta pertama. Saya sungguh berterima kasih kepada semua orang yang terlibat dalam perkembangan proyek ini.

Salah satu alasan prosesnya memakan waktu selama itu adalah keinginan untuk mencapai suatu tingkat kelengkapan sebelum secara resmi menyebutnya v2. Kenyataannya, tidak pernah ada waktu yang sempurna untuk menandai sebuah rilis—selalu ada masalah yang belum terselesaikan atau fitur "satu lagi saja" yang ingin disertakan. Namun, menandai rilis mayor yang belum sempurna memberikan sedikit stabilitas bagi pengguna proyek sekaligus kesempatan bagi para pengembang untuk memulai dari titik yang baru.

Rilis ini jauh melampaui segala yang pernah saya harapkan. Saya berharap Anda menikmatinya sebesar kenikmatan yang kami rasakan saat mengembangkannya.

## Apa *itu* Wails?

Jika Anda belum mengenal Wails, Wails adalah proyek yang memungkinkan pemrogram Go menyediakan frontend yang kaya untuk program Go mereka dengan menggunakan teknologi web yang sudah dikenal. Wails merupakan alternatif Electron yang ringan dan berbasis Go. Informasi lebih lengkap tersedia di [situs resmi](https://v2.wails.io/docs/introduction).

## Apa yang baru?

Rilis v2 merupakan lompatan besar bagi proyek ini dan mengatasi banyak kendala pada v1. Jika Anda belum membaca satu pun postingan blog tentang rilis Beta untuk [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/), atau [Linux](/blog/wails-v2-beta-for-linux/), saya menyarankan Anda membacanya karena semua perubahan utama dibahas secara lebih mendetail. Ringkasnya:

- Komponen Webview2 untuk Windows yang mendukung standar web modern dan kemampuan debugging.
- [Tema gelap/terang](https://v2.wails.io/docs/reference/options#theme) + [tema khusus](https://v2.wails.io/docs/reference/options#customtheme) di Windows.
- Windows kini tidak lagi memerlukan CGO.
- Dukungan siap pakai untuk templat proyek Svelte, Vue, React, Preact, Lit, dan Vanilla.
- Integrasi [Vite](https://vitejs.dev/) yang menyediakan lingkungan pengembangan dengan hot reload untuk aplikasi Anda.
- [Menu](https://v2.wails.io/docs/guides/application-development#application-menu) dan [dialog](https://v2.wails.io/docs/reference/runtime/dialog) aplikasi native.
- Efek transparansi jendela native untuk [Windows](https://v2.wails.io/docs/reference/options#windowistranslucent) dan [macOS](https://v2.wails.io/docs/reference/options#windowistranslucent-1). Mendukung latar belakang Mica dan Acrylic.
- Buat [penginstal NSIS](https://v2.wails.io/docs/guides/windows-installer) dengan mudah untuk deployment di Windows.
- [Pustaka runtime](https://v2.wails.io/docs/reference/runtime/intro) yang kaya dan menyediakan metode utilitas untuk memanipulasi jendela, menangani event, dialog, menu, dan logging.
- Dukungan untuk [mengobfusikasi](https://v2.wails.io/docs/guides/obfuscated) aplikasi Anda menggunakan [garble](https://github.com/burrowers/garble).
- Dukungan untuk mengompresi aplikasi Anda menggunakan [UPX](https://upx.github.io/).
- Pembuatan TypeScript secara otomatis dari struct Go. Informasi selengkapnya tersedia [di sini](https://v2.wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Tidak ada pustaka atau DLL tambahan yang perlu didistribusikan bersama aplikasi Anda. Di platform mana pun.
- Tidak perlu membundel aset frontend. Cukup kembangkan aplikasi Anda seperti aplikasi web lainnya.

## Penghargaan & Terima Kasih

Mencapai v2 membutuhkan upaya yang sangat besar. Terdapat sekitar 2200 commit dari 89 kontributor antara versi alfa awal dan rilis hari ini, serta jauh lebih banyak orang yang menyediakan terjemahan, pengujian, umpan balik, dan bantuan di forum diskusi maupun pelacak masalah. Saya sungguh luar biasa berterima kasih kepada Anda semua. Saya juga ingin menyampaikan terima kasih yang sangat istimewa kepada semua sponsor proyek yang telah memberikan arahan, saran, dan umpan balik. Segala hal yang Anda lakukan sangat kami hargai.

Ada beberapa orang yang ingin saya beri penghargaan khusus:

Pertama-tama, terima kasih yang **sebesar-besarnya** kepada [@stffabi](https://github.com/stffabi) yang telah memberikan begitu banyak kontribusi yang bermanfaat bagi kita semua, sekaligus banyak membantu menangani berbagai masalah. Ia menghadirkan sejumlah fitur penting, seperti dukungan server pengembangan eksternal, yang mengubah penawaran mode pengembangan kami dengan memungkinkan kami memanfaatkan kemampuan luar biasa [Vite](https://vitejs.dev/). Tidak berlebihan jika dikatakan bahwa Wails v2 akan menjadi rilis yang jauh kurang menarik tanpa [kontribusinya yang luar biasa](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04). Terima kasih banyak, @stffabi!

Saya juga ingin menyampaikan apresiasi yang sebesar-besarnya kepada [@misitebao](https://github.com/misitebao) yang tanpa lelah memelihara situs web, menyediakan terjemahan bahasa Tionghoa, mengelola Crowdin, dan membantu penerjemah baru agar cepat beradaptasi. Ini merupakan tugas yang sangat penting, dan saya amat berterima kasih atas seluruh waktu dan upaya yang telah dicurahkan! Anda luar biasa!

Terakhir, tetapi tidak kalah penting, terima kasih yang sebesar-besarnya kepada Mat Ryer yang telah memberikan saran dan dukungan selama pengembangan v2. Menulis xBar bersama-sama menggunakan versi Alfa awal v2 turut membantu membentuk arah v2, sekaligus memberi saya pemahaman tentang sejumlah kelemahan desain dalam rilis-rilis awal. Dengan senang hati saya mengumumkan bahwa mulai hari ini, kami akan mulai mem-porting xBar ke Wails v2, dan xBar akan menjadi aplikasi unggulan proyek ini. Terima kasih, Mat!

## Pelajaran yang Dipetik

Ada sejumlah pelajaran yang kami petik dalam perjalanan menuju v2 dan pelajaran tersebut akan membentuk arah pengembangan selanjutnya.

## Rilis yang Lebih Kecil, Cepat, dan Terfokus

Selama pengembangan v2, banyak fitur dan perbaikan bug dikembangkan secara ad hoc. Hal ini menyebabkan siklus rilis menjadi lebih panjang dan proses debugging menjadi lebih sulit. Ke depannya, kami akan lebih sering membuat rilis dengan jumlah fitur yang lebih sedikit. Setiap rilis akan mencakup pembaruan dokumentasi serta pengujian menyeluruh. Kami berharap rilis yang lebih kecil, cepat, dan terfokus ini akan menghasilkan lebih sedikit regresi dan dokumentasi yang lebih berkualitas.

## Mendorong Keterlibatan

Saat memulai proyek ini, saya ingin segera membantu setiap orang yang mengalami masalah. Setiap masalah terasa "pribadi" dan saya ingin menyelesaikannya secepat mungkin. Pendekatan ini tidak berkelanjutan dan pada akhirnya justru merugikan keberlangsungan proyek. Ke depannya, saya akan memberi lebih banyak ruang bagi orang lain untuk ikut menjawab pertanyaan dan melakukan triase masalah. Akan sangat membantu jika tersedia perkakas untuk mendukung hal ini. Jadi, jika Anda memiliki saran, silakan bergabung dalam diskusi [di sini](https://github.com/wailsapp/wails/discussions/1855).

## Belajar Mengatakan Tidak

Semakin banyak orang yang terlibat dalam sebuah proyek Open Source, semakin banyak pula permintaan fitur tambahan yang mungkin bermanfaat atau mungkin tidak bermanfaat bagi sebagian besar orang. Fitur-fitur tersebut memerlukan waktu awal untuk dikembangkan dan di-debug, lalu menimbulkan biaya pemeliharaan berkelanjutan sejak saat itu. Saya sendiri paling bersalah dalam hal ini karena sering ingin "merebus lautan" alih-alih menyediakan fitur minimum yang layak. Ke depannya, kami perlu lebih sering mengatakan "Tidak" untuk penambahan fitur inti dan memusatkan energi pada cara memberdayakan pengembang agar dapat menyediakan fungsionalitas tersebut sendiri. Kami sedang mempertimbangkan plugin secara serius untuk skenario ini. Dengan demikian, siapa pun dapat memperluas proyek sesuai kebutuhan sekaligus memperoleh cara yang mudah untuk berkontribusi pada proyek.

## Menatap Masa Depan

Kami sudah mempertimbangkan begitu banyak fitur inti untuk ditambahkan ke Wails dalam siklus pengembangan mayor berikutnya. [Roadmap](https://github.com/wailsapp/wails/discussions/1484) dipenuhi berbagai gagasan menarik, dan saya tidak sabar untuk mulai mengerjakannya. Salah satu permintaan utama adalah dukungan beberapa jendela. Ini cukup rumit untuk diterapkan dengan benar, dan kami mungkin perlu mempertimbangkan API alternatif karena API saat ini tidak dirancang dengan kebutuhan tersebut. Berdasarkan sejumlah gagasan awal dan umpan balik, saya rasa Anda akan menyukai arah yang sedang kami pertimbangkan.

Secara pribadi, saya sangat antusias dengan prospek menjalankan aplikasi Wails di perangkat seluler. Kami sudah memiliki proyek demo yang menunjukkan bahwa aplikasi Wails dapat dijalankan di Android, jadi saya sangat ingin menjajaki sejauh mana hal ini dapat dikembangkan!

Poin terakhir yang ingin saya sampaikan adalah mengenai kesetaraan fitur. Sejak lama, salah satu prinsip utama kami adalah tidak menambahkan apa pun ke proyek ini tanpa dukungan lintas platform secara penuh. Meskipun sejauh ini hal tersebut (sebagian besar) terbukti dapat dicapai, prinsip ini benar-benar menghambat proyek dalam merilis fitur baru. Ke depannya, kami akan menerapkan pendekatan yang sedikit berbeda: setiap fitur baru yang tidak dapat segera dirilis untuk semua platform akan dirilis melalui konfigurasi atau API eksperimental. Dengan demikian, pengguna awal pada platform tertentu dapat mencoba fitur tersebut dan memberikan umpan balik yang akan digunakan untuk menyempurnakan desain akhirnya. Tentu saja, ini berarti stabilitas API tidak dijamin hingga fitur tersebut didukung sepenuhnya oleh semua platform yang dapat mendukungnya, tetapi setidaknya pengembangan tidak lagi terhambat.

## Kata Penutup

Saya sangat bangga dengan pencapaian kami dalam rilis V2. Sungguh mengagumkan melihat apa yang telah berhasil dibuat orang-orang menggunakan rilis beta sejauh ini, termasuk aplikasi berkualitas seperti [Varly](https://varly.app/), [Surge](https://getsurge.io/), dan [October](https://october.utf9k.net/). Saya mendorong Anda untuk melihat aplikasi-aplikasi tersebut.

Rilis ini terwujud berkat kerja keras banyak kontributor. Meskipun dapat diunduh dan digunakan secara gratis, rilis ini bukanlah hasil dari upaya tanpa biaya. Jangan salah, proyek ini membutuhkan pengorbanan yang besar. Pengorbanan itu bukan hanya waktu saya serta waktu setiap kontributor, tetapi juga waktu yang harus mereka lewatkan bersama teman dan keluarga. Karena itulah saya sangat berterima kasih atas setiap detik yang telah dicurahkan untuk mewujudkan proyek ini. Semakin banyak kontributor yang bergabung, semakin merata upaya ini dapat dibagi dan semakin banyak pula yang dapat kita capai bersama. Saya ingin mendorong Anda semua untuk memilih satu hal yang dapat Anda kontribusikan, baik itu mengonfirmasi bug yang dilaporkan seseorang, menyarankan perbaikan, memperbarui dokumentasi, maupun membantu orang yang membutuhkannya. Semua hal kecil ini memberikan dampak yang begitu besar! Akan sangat luar biasa jika Anda juga menjadi bagian dari perjalanan menuju v3.

Selamat menikmati!

&dash; Lea

PS: Jika Anda atau perusahaan Anda merasa Wails bermanfaat, harap pertimbangkan untuk [mensponsori proyek ini](https://github.com/sponsors/leaanthony). Terima kasih!

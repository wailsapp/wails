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

Hari ini menandai dirilisnya [Wails](https://wails.io) v2. Sudah sekitar 18 bulan sejak versi alfa v2 pertama dan sekitar satu tahun sejak rilis beta pertama. Saya sungguh berterima kasih kepada semua orang yang terlibat dalam perkembangan proyek ini.

Salah satu alasan prosesnya memakan waktu selama itu adalah keinginan untuk mencapai suatu tingkat kelengkapan sebelum secara resmi menyebutnya v2. Kenyataannya, tidak pernah ada waktu yang sempurna untuk menandai sebuah rilis—selalu ada masalah yang belum terselesaikan atau fitur "satu lagi saja" yang ingin disertakan. Namun, menandai rilis mayor yang belum sempurna dapat memberikan sedikit stabilitas bagi pengguna proyek, sekaligus memberi kesempatan kepada para pengembang untuk memulai kembali dengan lebih segar.

Rilis ini melampaui semua yang pernah saya harapkan. Saya berharap rilis ini memberi Anda kesenangan sebesar yang kami rasakan saat mengembangkannya.

## *Apa itu* Wails?

Jika Anda belum mengenal Wails, ini adalah proyek yang memungkinkan pemrogram Go menyediakan frontend kaya fitur untuk program Go mereka dengan menggunakan teknologi web yang sudah dikenal. Wails merupakan alternatif Electron yang ringan dan berbasis Go. Informasi lebih lengkap dapat ditemukan di [situs resmi](https://wails.io/docs/introduction).

## Apa yang baru?

Rilis v2 merupakan lompatan besar bagi proyek ini dan mengatasi banyak kendala yang ada di v1. Jika Anda belum membaca postingan blog tentang rilis Beta untuk [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/), atau [Linux](/blog/wails-v2-beta-for-linux/), saya menyarankan Anda membacanya karena semua perubahan utama dibahas secara lebih mendetail. Ringkasnya:

- Komponen Webview2 untuk Windows yang mendukung standar web modern dan kemampuan debugging.
- [Tema gelap/terang](https://wails.io/docs/reference/options#theme) + [tema khusus](https://wails.io/docs/reference/options#customtheme) di Windows.
- Windows kini tidak memerlukan CGO.
- Dukungan siap pakai untuk templat proyek Svelte, Vue, React, Preact, Lit, dan Vanilla.
- Integrasi [Vite](https://vitejs.dev/) yang menyediakan lingkungan pengembangan dengan hot reload untuk aplikasi Anda.
- [Menu](https://wails.io/docs/guides/application-development#application-menu) dan [dialog](https://wails.io/docs/reference/runtime/dialog) aplikasi native.
- Efek tembus cahaya sebagian pada jendela native untuk [Windows](https://wails.io/docs/reference/options#windowistranslucent) dan [macOS](https://wails.io/docs/reference/options#windowistranslucent-1). Mendukung latar Mica dan Acrylic.
- Buat [penginstal NSIS](https://wails.io/docs/guides/windows-installer) dengan mudah untuk deployment di Windows.
- [Pustaka runtime](https://wails.io/docs/reference/runtime/intro) yang kaya fitur dan menyediakan berbagai metode utilitas untuk memanipulasi jendela, menangani peristiwa, dialog, menu, serta logging.
- Dukungan untuk [mengobfuscasi](https://wails.io/docs/guides/obfuscated) aplikasi Anda menggunakan [garble](https://github.com/burrowers/garble).
- Dukungan untuk mengompresi aplikasi Anda menggunakan [UPX](https://upx.github.io/).
- Pembuatan TypeScript secara otomatis dari struct Go. Informasi selengkapnya tersedia [di sini](https://wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Tidak ada pustaka atau DLL tambahan yang perlu didistribusikan bersama aplikasi Anda, di platform mana pun.
- Tidak perlu membundel aset frontend. Cukup kembangkan aplikasi Anda seperti aplikasi web lainnya.

## Penghargaan & Terima Kasih

Mencapai v2 membutuhkan upaya yang sangat besar. Ada ~2200 commit dari 89 kontributor sejak versi alfa awal hingga rilis hari ini, serta jauh lebih banyak lagi pihak yang telah menyediakan terjemahan, melakukan pengujian, memberikan masukan, dan membantu di forum diskusi maupun pelacak isu. Saya sungguh luar biasa berterima kasih kepada Anda semua. Saya juga ingin menyampaikan terima kasih yang teramat istimewa kepada semua sponsor proyek yang telah memberikan arahan, saran, dan masukan. Segala hal yang Anda lakukan sangat kami hargai.

Ada beberapa orang yang ingin saya sebutkan secara khusus:

Pertama-tama, terima kasih yang **sebesar-besarnya** kepada [@stffabi](https://github.com/stffabi) yang telah memberikan begitu banyak kontribusi yang bermanfaat bagi kita semua, sekaligus menyediakan banyak dukungan untuk berbagai masalah. Ia telah menghadirkan beberapa fitur penting, seperti dukungan server pengembangan eksternal yang mengubah penawaran mode pengembangan kami dengan memungkinkan kami memanfaatkan kemampuan luar biasa [Vite](https://vitejs.dev/). Tidak berlebihan jika dikatakan bahwa Wails v2 akan menjadi rilis yang jauh kurang menarik tanpa [kontribusinya yang luar biasa](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04). Terima kasih banyak, @stffabi!

Saya juga ingin memberikan apresiasi yang sebesar-besarnya kepada [@misitebao](https://github.com/misitebao) yang tanpa lelah memelihara situs web, menyediakan terjemahan bahasa Mandarin, mengelola Crowdin, dan membantu penerjemah baru agar dapat segera berkontribusi. Ini merupakan tugas yang sangat penting, dan saya benar-benar berterima kasih atas seluruh waktu dan upaya yang dicurahkan untuknya! Anda luar biasa!

Terakhir, tetapi tidak kalah penting, terima kasih yang sebesar-besarnya kepada Mat Ryer yang telah memberikan saran dan dukungan selama pengembangan v2. Menulis xBar bersama-sama menggunakan versi Alfa awal v2 membantu membentuk arah v2, sekaligus memberi saya pemahaman mengenai beberapa kekurangan desain dalam rilis-rilis awal. Dengan gembira saya mengumumkan bahwa mulai hari ini, kami akan mulai mem-porting xBar ke Wails v2, dan aplikasi tersebut akan menjadi aplikasi unggulan proyek ini. Terima kasih, Mat!

## Pelajaran yang Dipetik

Ada sejumlah pelajaran yang kami petik dalam perjalanan menuju v2, yang akan membentuk arah pengembangan selanjutnya.

## Rilis yang Lebih Kecil, Lebih Cepat, dan Terfokus

Selama pengembangan v2, banyak fitur dan perbaikan bug dikembangkan secara ad hoc. Hal ini menyebabkan siklus rilis yang lebih panjang dan membuatnya lebih sulit untuk di-debug. Ke depannya, kami akan lebih sering membuat rilis dengan jumlah fitur yang lebih sedikit. Setiap rilis akan mencakup pembaruan dokumentasi serta pengujian menyeluruh. Semoga rilis yang lebih kecil, lebih cepat, dan terfokus ini menghasilkan lebih sedikit regresi serta dokumentasi yang lebih berkualitas.

## Mendorong Keterlibatan

Saat memulai proyek ini, saya ingin segera membantu setiap orang yang mengalami masalah. Setiap masalah terasa "personal" dan saya ingin menyelesaikannya secepat mungkin. Pendekatan ini tidak berkelanjutan dan pada akhirnya justru menghambat kelangsungan jangka panjang proyek. Ke depannya, saya akan memberikan lebih banyak ruang bagi orang lain untuk turut menjawab pertanyaan dan melakukan triase masalah. Akan sangat membantu jika tersedia peralatan untuk mendukung hal ini. Jadi, jika Anda memiliki saran, silakan bergabung dalam diskusi [di sini](https://github.com/wailsapp/wails/discussions/1855).

## Belajar Mengatakan Tidak

Semakin banyak orang yang terlibat dalam proyek Sumber Terbuka, semakin banyak pula permintaan fitur tambahan yang mungkin bermanfaat atau mungkin juga tidak bagi sebagian besar orang. Fitur-fitur ini akan membutuhkan waktu pada awalnya untuk dikembangkan dan di-debug, kemudian menimbulkan biaya pemeliharaan berkelanjutan. Saya sendiri adalah orang yang paling sering melakukan hal ini, karena kerap ingin mengerjakan terlalu banyak hal sekaligus alih-alih menyediakan fitur minimum yang layak. Ke depannya, kita perlu lebih sering mengatakan "Tidak" terhadap penambahan fitur inti dan memusatkan upaya pada cara memberdayakan pengembang agar dapat menyediakan fungsionalitas tersebut sendiri. Kami sedang mempertimbangkan plugin secara serius untuk skenario ini. Dengan demikian, siapa pun dapat memperluas proyek sesuai kebutuhan mereka, sekaligus memperoleh cara yang mudah untuk berkontribusi pada proyek.

## Menatap Masa Depan

Sudah ada begitu banyak fitur inti yang ingin kami tambahkan ke Wails dalam siklus pengembangan besar berikutnya. [Peta jalan](https://github.com/wailsapp/wails/discussions/1484) dipenuhi berbagai gagasan menarik, dan saya tidak sabar untuk mulai mengerjakannya. Salah satu permintaan utama adalah dukungan untuk banyak jendela. Fitur ini cukup rumit, dan agar dapat diterapkan dengan benar, kami mungkin perlu mempertimbangkan API alternatif karena API saat ini tidak dirancang untuk kebutuhan tersebut. Berdasarkan beberapa gagasan awal dan masukan yang kami terima, saya rasa Anda akan menyukai arah pengembangan yang sedang kami pertimbangkan.

Secara pribadi, saya sangat antusias dengan kemungkinan menjalankan aplikasi Wails di perangkat seluler. Kami sudah memiliki proyek demo yang menunjukkan bahwa aplikasi Wails dapat dijalankan di Android, jadi saya benar-benar tidak sabar untuk menjajaki sejauh mana kami dapat mengembangkannya!

Hal terakhir yang ingin saya bahas adalah kesetaraan fitur. Sejak lama, salah satu prinsip inti kami adalah tidak menambahkan apa pun ke proyek tanpa dukungan lintas platform sepenuhnya. Meskipun sejauh ini hal tersebut (sebagian besar) terbukti dapat dicapai, prinsip ini sangat menghambat proyek dalam merilis fitur baru. Ke depannya, kami akan menerapkan pendekatan yang sedikit berbeda: setiap fitur baru yang belum dapat segera dirilis untuk semua platform akan dirilis melalui konfigurasi atau API eksperimental. Dengan demikian, pengguna awal di platform tertentu dapat mencoba fitur tersebut dan memberikan masukan yang akan digunakan dalam desain akhir fitur itu. Tentu saja, ini berarti stabilitas API tidak dijamin hingga fitur tersebut didukung sepenuhnya oleh semua platform yang dapat mendukungnya, tetapi setidaknya pengembangan tidak lagi terhambat.

## Kata Penutup

Saya sangat bangga atas pencapaian kami dalam rilis V2. Sungguh menakjubkan melihat berbagai hal yang telah berhasil dibuat orang-orang menggunakan rilis beta sejauh ini. Aplikasi berkualitas seperti [Varly](https://varly.app/), [Surge](https://getsurge.io/), dan [October](https://october.utf9k.net/). Saya mendorong Anda untuk mencobanya.

Rilis ini terwujud berkat kerja keras banyak kontributor. Meskipun dapat diunduh dan digunakan secara gratis, rilis ini bukanlah hasil yang dicapai tanpa biaya. Jangan salah: proyek ini membutuhkan pengorbanan yang besar. Bukan hanya waktu saya dan waktu setiap kontributor yang tercurah, tetapi juga waktu kebersamaan mereka dengan teman dan keluarga yang harus dikorbankan. Karena itu, saya sangat berterima kasih atas setiap detik yang telah didedikasikan untuk mewujudkan proyek ini. Semakin banyak kontributor yang bergabung, semakin luas upaya ini dapat dibagi dan semakin banyak yang dapat kita capai bersama. Saya ingin mendorong Anda semua untuk memilih satu hal yang dapat Anda kontribusikan, entah itu mengonfirmasi bug yang dilaporkan seseorang, menyarankan perbaikan, memperbarui dokumentasi, atau membantu seseorang yang membutuhkannya. Semua hal kecil ini memberikan dampak yang begitu besar! Akan sangat luar biasa jika Anda juga menjadi bagian dari perjalanan menuju v3.

Selamat menikmati!

&dash; Lea

PS: Jika Anda atau perusahaan Anda merasakan manfaat Wails, harap pertimbangkan untuk [mensponsori proyek ini](https://github.com/sponsors/leaanthony). Terima kasih!

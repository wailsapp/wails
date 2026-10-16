---
title: "Bangun Aplikasi Desktop dengan Go"
description: "Aplikasi desktop native menggunakan Go dan teknologi web"
banner: {"content":"Wails v3 saat ini masih dalam versi beta. \u003ca href=\"https://v2.wails.io\"\u003eMencari dokumentasi v2?\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"Mulai","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"Lihat Tutorial","variant":"secondary"}],"image":{"alt":"Logo Wails","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"Bangun aplikasi yang indah dan berperforma tinggi menggunakan Go serta teknologi web modern. Satu basis kode. Tiga platform. Tanpa browser.*"}
sourcePath: "index.md"
---

<style>
  /* Hero background — the neon "digital Wales" mountain, full-bleed and fixed. */
  body::after {
    content: '';
    position: fixed;
    inset: 0;
    z-index: -1;
    background: url('/digital_wales_master.webp') center center / cover no-repeat;
    opacity: 0.35;
    pointer-events: none;
  }
  
  /* Gradient overlay - highest z-index */
  html::before {
    content: '';
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: 
      linear-gradient(to bottom, var(--sl-color-bg) 0%, transparent 15%, transparent 85%, var(--sl-color-bg) 100%),
      linear-gradient(to right, var(--sl-color-bg) 0%, transparent 10%, transparent 90%, var(--sl-color-bg) 100%);
    z-index: 0;
    pointer-events: none;
  }
  
  
  /* Hero padding */
  .hero {
    padding-top: 3rem !important;
    padding-bottom: 2rem !important;
  }
  
  .small-buttons {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-top: 1rem;
  }
  
  .small-buttons a {
    display: inline-block;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    border-radius: 6px;
    text-decoration: none;
    background: var(--sl-color-gray-6);
    color: var(--sl-color-white);
    border: 1px solid var(--sl-color-gray-5);
    transition: all 0.2s;
  }
  
  .small-buttons a:hover {
    background: var(--sl-color-gray-5);
    border-color: var(--sl-color-accent);
  }
  
  /* Round card corners and add translucent blur effect */
  .mpress-card {
    border-radius: 12px;
    background: rgba(var(--sl-color-gray-6-rgb), 0.6) !important;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
  }

</style>

<span class="browser-footnote" data-browser-footnote>* kecuali jika Anda memang menginginkannya</span>
<script>
(() => {
  const words = ['Desktop', 'Server', 'Mobile'];
  let index = 0;
  const init = () => {
    const tagline = document.querySelector('.mpress-frontmatter-hero-tagline');
    const note = document.querySelector('[data-browser-footnote]');
    if (tagline && note && !tagline.contains(note)) tagline.appendChild(note);
    const title = document.querySelector('.mpress-frontmatter-hero-title');
    if (!title || title.querySelector('.morph-word-title')) return;
    if (title.textContent.trim() !== 'Build Desktop Apps with Go') return;
    const titleLine = document.createElement('span');
    titleLine.className = 'morph-title-line';
    const prefix = document.createElement('span');
    prefix.textContent = 'Build';
    const wordWrap = document.createElement('span');
    wordWrap.className = 'morph-word-title';
    const activeWord = document.createElement('span');
    activeWord.textContent = words[0];
    wordWrap.append(activeWord);
    titleLine.append(prefix, wordWrap);
    title.replaceChildren(titleLine, document.createElement('br'), document.createTextNode('Apps with Go'));
    const rotate = () => {
      const word = title.querySelector('.morph-word-title span');
      if (!word) return;
      word.classList.add('morph-word-out');
      setTimeout(() => {
        index = (index + 1) % words.length;
        word.textContent = words[index];
        word.classList.remove('morph-word-out');
        setTimeout(rotate, 3600);
      }, 650);
    };
    setTimeout(rotate, 3600);
  };
  document.readyState === 'loading' ? document.addEventListener('DOMContentLoaded', init) : init();
})();
</script>

## Mulai Cepat

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

Setelah itu, Anda siap membuat proyek pertama:

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**Aplikasi Anda kini berjalan** dengan hot reload dan binding Go-ke-JS yang aman secara tipe.

Mengalami masalah dengan `wails3 setup`? Lihat [panduan instalasi manual](/quick-start/installation/).

## Aplikasi desktop Anda sudah menjadi aplikasi seluler

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Tidak perlu mengubah kode Go Anda.** `main.go` yang sama, layanan yang sama, frontend yang sama — Wails secara otomatis mengompilasinya untuk iOS dan Android. Fitur khusus platform (haptik, dialog native, area aman) tersedia saat Anda menginginkannya, tetapi Anda tidak memerlukannya untuk merilis aplikasi.

[Dokumentasi seluler →](/guides/mobile/)

---

## Mengapa Wails?

@cards{cols="2"}
🚀 Performa yang Dirasakan Pengguna
- Biner ~15MB dibandingkan dengan 150MB pada Electron
- Memori dasar ~10MB dibandingkan dengan 100MB+
- Waktu mulai &lt;0.5 detik dibandingkan dengan 2-3 detik
- Rendering native menggunakan WebView OS
- Tanpa overhead browser yang dibundel

---
⚙ Pengalaman Developer
- Satu basis kode Go untuk semua platform
- Framework web apa pun - React, Vue, Svelte
- Hot reload selama pengembangan
- Binding yang dibuat otomatis untuk memanggil Go dengan mudah dari Javascript
- IPC dalam memori. Tanpa port jaringan

---
✓ Siap untuk Desktop
- Beberapa jendela dengan siklus hidup masing-masing
- Menu native dan baki sistem
- Dialog file native platform
- Integrasi sistem dan pintasan
- Alat untuk penandatanganan kode dan pengemasan

---
▣ Desktop & Seluler
- Windows, macOS, Linux, iOS, Android
- Basis kode yang sama, tanpa penulisan ulang
- WebView native di setiap platform
- Tanpa port terbuka, tanpa server localhost
- Fitur platform tersedia saat diperlukan

@end

## Langkah Berikutnya

Berikutnya: [bangun aplikasi lengkap](/tutorials/03-notes-vanilla/), telusuri [contoh](https://github.com/wailsapp/wails/tree/master/v3/examples), atau lihat [referensi API](/reference/application/). Bermigrasi dari v2? Lihat [panduan peningkatan versi](/migration/v2-to-v3/).

@note{type="info" title="Wails v3 Beta"}
Wails v3 adalah perangkat lunak beta dengan API desktop yang stabil. Sejumlah tim telah menggunakannya di lingkungan produksi, tetapi sebaiknya melakukan pengujian menyeluruh sebelum penerapan selagi kami menyelesaikan penyempurnaan akhir untuk 3.0.

@end

@cards{cols="1"}
♥ Dukung Pengembangan Wails
Wails gratis dan bersumber terbuka, dibuat oleh developer untuk developer. Jika Wails membantu Anda membangun aplikasi yang luar biasa, pertimbangkan untuk mendukung pengembangannya secara berkelanjutan.

Sponsor Anda membantu memelihara proyek, meningkatkan dokumentasi, dan mengembangkan fitur baru yang bermanfaat bagi seluruh komunitas.

[Jadilah Sponsor →](https://github.com/sponsors/leaanthony)

@end

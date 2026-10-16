---
title: "修正文件"
description: "使用 M-Press 提交 Wails v3 文件修正 PR。"
sourcePath: "contributing/documentation.md"
---

歡迎提交修正 PR。您可以修正錯字、失效的連結、過時的範例、不清楚的說明或翻譯。若修正僅涉及文件，您不需要先建立議題，也不需要提供失敗的程式碼測試。

## 在本機預覽

建立 [wailsapp/wails](https://github.com/wailsapp/wails/fork) 的分支版本（fork），複製您的分支版本至本機，然後從 `master` 建立分支。

安裝鎖定版本的文件產生器：

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

您也可以從[M-Press v1.0.17發行版本](https://github.com/leaanthony/mpress/releases/tag/v1.0.17)下載經過驗證的二進位檔。

在 Wails 儲存庫的根目錄中：

```sh
mpress version
mpress dev
```

編輯 `docs/mpress/content/` 中的 `.md` 原始檔。英文是預設語言，其檔案直接位於該目錄中。現有翻譯位於 `fr/` 和 `id/` 等語言資料夾中。預覽會在您儲存時重新建置。

請保留每個頁面頂端的中繼資料區塊，以及成對的 `@...`／`@end` 元件。一般段落、標題、清單和圍欄程式碼區塊均可作為文字編輯。請勿編輯 `docs/mpress/site/` 中產生的檔案。

## 檢查修正內容

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

在瀏覽器中檢查變更後的頁面，並執行您修改過的所有程式碼範例。修正翻譯時，請將變更後的完整段落與英文原文比較。請保持命令、API 名稱、連結、程式碼範例和圖表連線不變。使用自然的技術用語，保留要求與注意事項；除了內文，也應翻譯可見的圖表標籤、導覽標籤和圖片說明。

每種已發布的語言都必須完整翻譯每個英文頁面。請勿使用英文預留位置或備援頁面。若只修正其中一種翻譯，可以僅變更該語言。如果您變更英文內容的含義，請更新其他已發布語言中對應的頁面；不需要重新產生無關的頁面。

## 提交拉取請求

針對 `master` 建立 PR。說明原本的問題、解釋您的修正，並列出您執行過的檢查。若有可見的版面配置變更，請附上螢幕擷取畫面；若變更了程式碼範例，請提供平台和版本詳細資訊。

不需要 Cloudflare 憑證或私人服務。公開 PR 的檢查會在不使用部署憑證的情況下建置並驗證靜態網站。

如需瞭解程式碼變更和功能提案，請參閱[參與 Wails 開發](/contributing/)。如需瞭解內部運作，請參閱[技術概觀](/contributing/overview/)。

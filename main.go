// x-browser: X (旧Twitter) 専用の超軽量デスクトップブラウザ。
// OS標準のWebViewエンジンをそのまま利用するため(macOS: WebKit /
// Windows: WebView2 / Linux: WebKitGTK)、Chromiumを同梱するElectron系より
// はるかに軽量・省メモリで動作する。
package main

import (
	"log"

	webview "github.com/webview/webview_go"
)

const (
	homeURL   = "https://x.com/home"
	winTitle  = "X"
	winWidth  = 1200
	winHeight = 840
)

// ページ読み込みのたびに注入するJS。
// キーボードショートカットだけで最低限のブラウザ操作(戻る/進む/更新/ホーム)を
// 実現し、ネイティブ側のバインドを増やさず軽量さを保つ。
const injectedJS = `
(function () {
  document.addEventListener('keydown', function (e) {
    var mod = e.metaKey || e.ctrlKey;

    // 更新: Cmd/Ctrl+R または F5
    if ((mod && e.key.toLowerCase() === 'r') || e.key === 'F5') {
      e.preventDefault();
      location.reload();
      return;
    }
    // 戻る: Alt+← または Cmd+[
    if ((e.altKey && e.key === 'ArrowLeft') || (mod && e.key === '[')) {
      e.preventDefault();
      history.back();
      return;
    }
    // 進む: Alt+→ または Cmd+]
    if ((e.altKey && e.key === 'ArrowRight') || (mod && e.key === ']')) {
      e.preventDefault();
      history.forward();
      return;
    }
    // ホームへ: Cmd/Ctrl+H
    if (mod && e.key.toLowerCase() === 'h') {
      e.preventDefault();
      location.href = 'https://x.com/home';
      return;
    }
  }, true);
})();
`

func main() {
	// 第1引数 true でデバッグ(開発者ツール)を有効化。
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle(winTitle)
	w.SetSize(winWidth, winHeight, webview.HintNone)

	// 新しいページが読み込まれるたびにショートカット用JSを再注入。
	w.Init(injectedJS)

	w.Navigate(homeURL)

	log.Println("x-browser 起動: ", homeURL)
	w.Run()
}

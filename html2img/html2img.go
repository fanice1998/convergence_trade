package html2img

import (
	"context"
	"os"

	"github.com/chromedp/chromedp"
)

// 透過 chromedp 將 html 渲染後得結果儲存成圖片
func SaveImage(fileURL string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte

	err := chromedp.Run(ctx, fullScreenshot(fileURL, &buf))
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("./examples/html/kline.png", buf, 0644); err != nil {
		panic(err)
	}
}

// 取得完整畫面圖片
func fullScreenshot(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.FullScreenshot(res, 90),
	}
}

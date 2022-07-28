package pdf

import (
	"context"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ServiceInterface interface {
	HTMLToPdf(html string) []byte
}

//TODO

type PDFService struct {
	chromeWebSocketURL string
}

func NewPDFService(chromeWebSocketURL string) *PDFService {
	return &PDFService{
		chromeWebSocketURL: chromeWebSocketURL,
	}
}

func (c *PDFService) HTMLToPdf(html string) []byte {
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), c.chromeWebSocketURL)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	// do this first so that its page.EventLoadEventFired event won't be caught
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
	); err != nil {
		panic(err)
	}

	var buf []byte
	var wg sync.WaitGroup
	wg.Add(1)
	// waiting to page load event
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev.(type) {
		case *page.EventLoadEventFired:
			go func() {
				if err := chromedp.Run(ctx, printToPDF(&buf)); err != nil {
					panic(err)
				}
				wg.Done()
			}()
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.PollFunction(
			"(html) => {document.open();document.write(html);document.close();return true;}",
			nil,
			chromedp.WithPollingArgs(html),
		),
	); err != nil {
		panic(err)
	}

	wg.Wait()
	page.Close()

	return buf
}

func printToPDF(res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf

			return nil
		}),
	}
}

package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"
    "time"
    "runtime"

    cdp "github.com/knq/chromedp"
    cdptypes "github.com/knq/chromedp/cdp"
    runner "github.com/knq/chromedp/runner"
)

func main() {
    var err error

    // create context
    ctxt, cancel := context.WithCancel(context.Background())
    defer cancel()

    start_load := time.Now()
    path := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    if runtime.GOOS != "windows" {
      path = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    }   
  
    c, err := cdp.New(ctxt, cdp.WithRunnerOptions(
      runner.Headless(path, 9222),
      runner.Flag("headless", true),
      runner.Flag("disable-gpu", true),
      runner.Flag("no-first-run", true),
      runner.Flag("no-default-browser-check", true),
      runner.Flag("window-size", "800,420"),
      runner.Flag("hide-scrollbars", "true"),
    ))
    //
    if err != nil {
      log.Fatal(err)
    }

    elapsed_load := time.Since(start_load)
    log.Printf("Chrome Loading Time: %s", elapsed_load)

    start := time.Now()

    // run task list
    // var site, res string

    for i := 0; i < 100; i++ {
        start = time.Now()
        log.Printf("Start Screenshot")
        err = c.Run(ctxt, screenshot())
        if err != nil {
            log.Fatal(err)
        }
        
        log.Printf("Screenshot took %s", time.Since(start))
    }


    // shutdown chrome
    err = c.Shutdown(ctxt)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Shutdown")
    }

    
    // wait for chrome to finish
    // err = c.Wait()
    // if err != nil {
    //     log.Fatal(err)
    // }

}

func screenshot() cdp.Tasks {
    
    var buf []byte
    output := cdp.Tasks{
        // cdp.Navigate("file:///Users/user/src/chromeless_lambda/local/test.html"),
        cdp.Navigate("http://google.com"),
        //cdp.Sleep(1 * time.Second),
        //cdp.WaitVisible(`#background-image`),
        cdp.CaptureScreenshot(&buf),
        cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
            return ioutil.WriteFile("/tmp/screenshot.png", buf, 0644)
            }),
    }
    
    return output
}

func googleSearch(q, text string, site, res *string) cdp.Tasks {
    var buf []byte
    sel := fmt.Sprintf(`//a[text()[contains(., '%s')]]`, text)
    return cdp.Tasks{
        cdp.Navigate(`https://www.google.com`),
        cdp.Sleep(2 * time.Second),
        cdp.WaitVisible(`#hplogo`, cdp.ByID),
        cdp.SendKeys(`#lst-ib`, q+"\n", cdp.ByID),
        cdp.WaitVisible(`#res`, cdp.ByID),
        cdp.Text(sel, res),
        cdp.Click(sel),
        cdp.Sleep(2 * time.Second),
        cdp.WaitVisible(`#footer`, cdp.ByQuery),
        cdp.WaitNotVisible(`div.v-middle > div.la-ball-clip-rotate`, cdp.ByQuery),
        cdp.Location(site),
        cdp.Screenshot(`#testimonials`, &buf, cdp.ByID),
        cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
            return ioutil.WriteFile("screenshot.png", buf, 0644)
        }),
    }
}
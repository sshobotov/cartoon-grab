package main

import (
    "flag"
    //"os"
    "image"
    "fmt"
    //"io/ioutil"
    "net/http"
    "launchpad.net/xmlpath"
	"code.google.com/p/gofpdf"
)

var sourceUrl       = flag.String("u", "", "Source URL")
var imageXPath      = flag.String("i", "", "Image xPath on the page")
var nextPageXPath   = flag.String("l", "", "Next page link xPath on the page")
const (
    userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36"
)

func main() {
    flag.Parse()

    if *initialUrl == "" {
        fmt.Println("Use -u to setup source url")
        return
    }
    if *imageXPath == "" {
        fmt.Println("Use -i to setup image xPath")
        return
    }
    if *nextPageXPath == "" {
        fmt.Println("Use -l to setup next page link xPath")
        return
    }
    pdf := gofpdf.New("P", "mm", "A4", "")
    client := &http.Client{}

    makePage(client, pdf, *initialUrl);
}

func makePage(client *http.Client, pdf *gofpdf.Fpdf, pageUrl string) (success boolean) {
    success = true

    if imgUrl, nextUrl, ok := collect(client, pageUrl); ok {
        add(pdf, imgUrl)
        if nextUrl != "" {
            makePage(client, pdf, nextUrl)
        }
    } else {
        success = false
    }
    return
}

func collect(client *http.Client, url string) (imgUrl string, nextUrl string, success boolean) {
    success = false

    req, err := http.NewRequest("GET", *sourceUrl, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    req.Header.Set("User-Agent", userAgent)

    response, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer response.Body.Close()

    if response.StatusCode >= 400 {
        fmt.Println("Unexpected source URL response status: %d", response.StatusCode)
        return
    }
    root, err := xmlpath.ParseHTML(response.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    path := xmlpath.MustCompile(*imageXPath + "/@src")
    if value, ok := path.String(root); ok {
        imgUrl = value
    }
    if (!ok) {
        success = true
        return
    }
    path := xmlpath.MustCompile(*nextPageXPath + "/@href")
    if value, ok := path.String(root); ok {
        nextUrl = value
    }
    return imgUrl, nextUrl, true
}

func add(client *http.Client, pdf *gofpdf.Fpdf, imgUrl string) (success string) {
    success = false

    req, err := http.NewRequest("GET", *sourceUrl, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    req.Header.Set("User-Agent", userAgent)

    response, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer response.Body.Close()

    m, _, err := image.Decode(resp.Body)
    if err != nil {
        fmt.Println(err)
    }
    g := m.Bounds()

    height := g.Dy()
    width := g.Dx()

    if (height < width) {
        pdf.AddPageFormat("L", nil)
    } else {
        pdf.AddPage()
    }

    
}
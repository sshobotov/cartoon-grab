package main

import (
    "flag"
    "fmt"
    "strings"
    "net/http"
    "launchpad.net/xmlpath"
	"code.google.com/p/gofpdf"
    "code.google.com/p/go-uuid/uuid"
)

var (
    initialUrl      = flag.String("u", "", "Source URL")
    imageXPath      = flag.String("i", "", "Image xPath on the page")
    nextPageXPath   = flag.String("l", "", "Next page link xPath on the page")
    baseUrl string
)
const (
    userAgent       = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36"
)

// cartoon-grab -u "http://www.mangareader.net/toukyou-kushu/1" -i "//table[@class='episode-table']//*[@id='imgholder']/a/img/@src" -l "//table[@class='episode-table']//*[@id='imgholder']/a/@href"

func main() {
    flag.Parse()

    if *initialUrl == "" {
        fmt.Println("Use -u to setup source url")
        return
    }
    urlParts := strings.SplitN(*initialUrl, "/", 4)
    baseUrl = urlParts[0] + "//" + urlParts[2]

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

    if ok := makePage(client, pdf, *initialUrl); ok {
        pdf.OutputFileAndClose(uuid.New() + ".pdf")
        fmt.Println("Done")
    }
}

func makePage(client *http.Client, pdf *gofpdf.Fpdf, pageUrl string) (success bool) {
    success = true

    if imgUrl, nextUrl, ok := collect(client, pageUrl); ok {
        if imgUrl == "" {
            return
        }
        added := add(client, pdf, imgUrl)
        if !added {
            success = false
            return
        }
        if nextUrl != "" {
            if (strings.HasPrefix(nextUrl, "/")) {
                nextUrl = baseUrl + nextUrl
            }
            fmt.Println("Next:", nextUrl)
            success = makePage(client, pdf, nextUrl)
        }
    } else {
        success = false
    }
    return
}

func collect(client *http.Client, url string) (imgUrl string, nextUrl string, success bool) {
    success = false

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Unable to create page request:", err)
        return
    }
    req.Close = true
    req.Header.Set("User-Agent", userAgent)

    response, err := client.Do(req)
    if err != nil {
        fmt.Println("Unable to complete page request:", err)
        return
    }
    defer response.Body.Close()

    if response.StatusCode >= 400 {
        fmt.Println("Unexpected source URL response status:", response.StatusCode)
        return
    }
    root, err := xmlpath.ParseHTML(response.Body)
    if err != nil {
        fmt.Println("Unable to parse HTML:", err)
        return
    }

    srcPath := xmlpath.MustCompile(*imageXPath)
    if value, ok := srcPath.String(root); ok {
        imgUrl = value
    } else {
        success = true
        return
    }
    hrefPath := xmlpath.MustCompile(*nextPageXPath)
    if value, ok := hrefPath.String(root); ok {
        nextUrl = value
    }
    return imgUrl, nextUrl, true
}

func add(client *http.Client, pdf *gofpdf.Fpdf, imgUrl string) (success bool) {
    success = false

    req, err := http.NewRequest("GET", imgUrl, nil)
    if err != nil {
        fmt.Println("Unable to create image request:", err)
        return
    }
    req.Close = true
    req.Header.Set("User-Agent", userAgent)

    response, err := client.Do(req)
    if err != nil {
        fmt.Println("Unable to complete image request:", err)
        return
    }
    defer response.Body.Close()

    tp := pdf.ImageTypeFromMime(response.Header["Content-Type"][0])
    infoPtr := pdf.RegisterImageReader(imgUrl, tp, response.Body)
    if !pdf.Ok() {
        return
    }
    width, height := infoPtr.Extent()

    if (height < width) {
        pdf.AddPageFormat("L", gofpdf.SizeType{ Wd: height, Ht: width })
    } else {
        pdf.AddPageFormat("P", gofpdf.SizeType{ Wd: width, Ht: height })
    }
    pdf.Image(imgUrl, 0, 0, width, height, false, tp, 0, "")
    
    success = pdf.Ok();
    return
}
package main

import (
    "flag"
    //"os"
    "fmt"
    "bytes"
    //"io/ioutil"
    "net/http"
    "launchpad.net/xmlpath"
	//"code.google.com/p/gofpdf"
)

var sourceUrl       = flag.String("u", "", "Source URL")
var imageXPath      = flag.String("i", "", "Image xPath on the page")
var nextPageXPath   = flag.String("l", "", "Next page link xPath on the page")

func main() {
    flag.Parse()

    if *sourceUrl == "" {
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
    client := &http.Client{}
    req, err := http.NewRequest("GET", *sourceUrl, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36")

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
    buf := new(bytes.Buffer)
    buf.ReadFrom(response.Body)
    s := buf.String()
    fmt.Println(s)
    path := xmlpath.MustCompile(*imageXPath + "/@src")
    if value, ok := path.String(root); ok {
        fmt.Println("Found:", value)
    }
}
package main

import (
    "encoding/json"
    "fmt"
    // "io/ioutil"
    "log"
    "net/http"
    // "net/http/httputil"
    "sync"
    "time"

    "github.com/julienschmidt/httprouter"
)

type Result struct {
    url string
    status int
}

func ClientGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    res := p.ByName("param")

    // ヘッダーセット
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // レスポンスに書き込む
    fmt.Fprintf(w, res)
}

func ClientPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    r.ParseForm()
    urls := r.Form["urls[]"]

    result := []Result{}

    GetResult(urls, &result, false)

    log.Println(result)
  
    res, err := json.Marshal(result)
    if err != nil {
        // 変換失敗時、500エラー
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ヘッダーセット
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // レスポンスに書き込む
    fmt.Fprintf(w, string(res))
}

func GetRequest(client *http.Client, url string, ch <-chan int, wg *sync.WaitGroup, result *[]Result) {
    
    defer func() {
        <- ch
        wg.Done()
    }()

    req, _ := http.NewRequest("GET", url, nil)
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    *result = append(*result, Result{ url: url, status: resp.StatusCode })
}

func PostRequest() {

}

func GetResult(urls []string, result *[]Result, isPost bool) {
    ch := make(chan int, 5)
    wg := &sync.WaitGroup{}

    log.Println(time.Now())

    client := &http.Client{}

    for _, url := range urls {
        if url == "" {
            continue
        }
        ch <- 1
        wg.Add(1)
        if isPost {
            break
        } else {
            go GetRequest(client, url, ch, wg, result)
        }
        
    }
    wg.Wait()

    log.Println(time.Now())
}

func main() {
    // HTTPルーターを初期化
    router := httprouter.New()

    router.GET("/Index/:param", ClientGet)
    router.POST("/Post", ClientPost)

    // Webサーバーを8080ポートで立ち上げる
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatal(err)
    }
}
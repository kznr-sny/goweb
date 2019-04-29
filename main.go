package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    // "net/http/httputil"
    "sync"
    "time"

    "github.com/julienschmidt/httprouter"
)

func ClientGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    res := p.ByName("param")
    // ヘッダーセット
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // レスポンスに書き込む
    fmt.Fprintf(w, res)
}

func ClientPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // 終了前に閉じる
    defer r.Body.Close()

    // ヘッダーセット
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        // 読み取り失敗時、400エラー
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // レスポンス用の構造体
    type ResponseParam struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
    var param ResponseParam

    // JSON => STRUCT
    err = json.Unmarshal(bodyBytes, &param)
    if err != nil {
        // 変換失敗時、400エラー
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // STRUCT => JSON
    res, _ := json.Marshal(param)
    // レスポンスに書き込む
    fmt.Fprintf(w, string(res))
}

func GetRequest(client *http.Client, url string)(int) {
    req, _ := http.NewRequest("GET", url, nil)

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return 0
    }
    defer resp.Body.Close()
    return resp.StatusCode

    // dumpResp, _ := httputil.DumpResponse(resp, true)
    // byteArray, _ := ioutil.ReadAll(resp.Body)
    // fmt.Println(string(byteArray))
}

func PostRequest() {

}

func GetFunc(client *http.Client, url string, ch <-chan int, wg *sync.WaitGroup) {
    log.Println(GetRequest(client, url), url)
    <- ch
    wg.Done()
}

func main() {
    // HTTPルーターを初期化
    router := httprouter.New()

    router.GET("/Index/:param", ClientGet)
    router.POST("/Post", ClientPost)

    urls := []string {
        "https://stackoverflow.com/",
        "http://yahoo.co.jp",
        "https://stackoverflow.com/",
        "http://yahoo.co.jp",
        "https://stackoverflow.com/",
        "http://yahoo.co.jp",
        "https://stackoverflow.com/",
        "http://yahoo.co.jp",
        "https://stackoverflow.com/",
        "http://yahoo.co.jp"}

    ch := make(chan int, 5)
    wg := &sync.WaitGroup{}

    log.Println(time.Now())
    client := &http.Client{}

    for _, url := range urls {
        ch <- 1
        wg.Add(1)
        go GetFunc(client, url, ch, wg)
    }
    wg.Wait()

    log.Println(time.Now())



    // Webサーバーを8080ポートで立ち上げる
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatal(err)
    }
}
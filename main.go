package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "sync"
    "strings"
    "time"

    "github.com/julienschmidt/httprouter"
    "github.com/kznr-sny/goweb/model"
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
    defer r.Body.Close()

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        // 読み取り失敗時、400エラー
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var rd model.RequestData
    var result []map[string]interface{}

    json.Unmarshal(bodyBytes, &rd)

    GetResult(rd, &result)
  
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

func GetResult(rd model.RequestData, result *[]map[string]interface{}) {
    ch := make(chan int, 5)
    wg := &sync.WaitGroup{}
    dict := make(map[int][]model.ParamData)
    client := &http.Client{}

    if len(rd.Params) > 0 {
        for _, v := range rd.Params {
            if _, ok := dict[v.ID]; ok {
                dict[v.ID] = append(dict[v.ID], model.ParamData{Key: v.Key, Value: v.Value})
            } else {
                var data []model.ParamData
                dict[v.ID] = append(data, model.ParamData{Key: v.Key, Value: v.Value})
            }
        }
    }
    
    log.Println(time.Now())

    for _, uri := range rd.Uris {
        if uri.Uri == "" {
            continue
        }
        ch <- 1
        wg.Add(1)
        if rd.IsPost {
            params := dict[uri.ID]
            go PostRequest(client, uri.Uri, params, ch, wg, result)
        } else {
            go GetRequest(client, uri.Uri, ch, wg, result)
        }
        
    }
    wg.Wait()

    log.Println(time.Now())
}

func GetRequest(client *http.Client, uri string, ch <-chan int, wg *sync.WaitGroup, result *[]map[string]interface{}) {
    
    defer func() {
        <- ch
        wg.Done()
    }()

    req, _ := http.NewRequest("GET", uri, nil)
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    *result = append(*result, map[string]interface{}{"uri": uri, "status": resp.StatusCode})
}

func PostRequest(client *http.Client, uri string, params []model.ParamData, ch <-chan int, wg *sync.WaitGroup, result *[]map[string]interface{}) {

    defer func() {
        <- ch
        wg.Done()
    }()

    postParams := url.Values{}
    for _, v := range params {
        if (v.Key == "") || (v.Value == "") {
            continue
        }
        postParams.Add(v.Key, v.Value)
    }
    
    req, _ := http.NewRequest("POST", uri, strings.NewReader(postParams.Encode()))
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    *result = append(*result, map[string]interface{}{"uri": uri, "status": resp.StatusCode})    
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
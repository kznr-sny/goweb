package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    res := p.ByName("param")
    // ヘッダーセット
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // レスポンスに書き込む
    fmt.Fprintf(w, res)
}

func GetJson(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func main() {
    // HTTPルーターを初期化
    router := httprouter.New()

    router.GET("/Index/:param", Index)
    router.POST("/GetJson", GetJson)

    // Webサーバーを8080ポートで立ち上げる
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatal(err)
    }
}
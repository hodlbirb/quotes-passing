package main;

import (
    "fmt"
    "strconv"
    "log"
    "os"
    "net/http"
    "encoding/json"
    "github.com/streadway/amqp"
)
import _ "net/http/pprof"

func main() {
    // Configure AMQP Server
    conn, conn_err := amqp.Dial(os.Getenv("AMQP_URL"))
    DieOnErr(conn_err, "Failed to connect to RabbitMQ")
    ch, ch_err := conn.Channel()
    DieOnErr(ch_err, "Failed to open a channel")
    queue, queue_err := ch.QueueDeclare(
        "sum",
        false,
        false,
        false,
        false,
        nil,
    )
    DieOnErr(queue_err, "Failed to declare a queue")
    // Configure HTTP Server
    http.HandleFunc("/sum", func(w http.ResponseWriter, req *http.Request) {
        // Set headers
        w.Header().Set("Content-Type", "application/json")
        // Decode url query params
        params := req.URL.Query()
        a, err := strconv.ParseFloat(params.Get("a"), 64)
        b, err := strconv.ParseFloat(params.Get("b"), 64)
        if err != nil {
            http.Error(w, "Failed to parse parameters", http.StatusBadRequest)
            return
        }
        // Encode result and write to HTTP response and publish to `sum` channel
        result := ResultResponse{Result: fmt.Sprintf("%.2f", (a+b))}
        json_result, _err := json.Marshal(result)
        if _err != nil {
            http.Error(w,"Server Error", http.StatusInternalServerError)
            return
        }
        pub_err := ch.Publish(
            "",
            queue.Name,
            false,
            false,
            amqp.Publishing {
                ContentType: "application/json",
                Body: json_result,
            })
        DieOnErr(pub_err, "Failed to publish a message")
        w.Write(json_result)
    })
    log.Println("[*] RPC server is up and running on port", os.Getenv("RPC_PORT"))
    log.Println(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("RPC_PORT")), nil))
}

func DieOnErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

type ResultResponse struct {
    Result string
}

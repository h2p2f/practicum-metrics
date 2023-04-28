# cmd/server

This code implements a server that listens on port 8080 (by default) and waits for a client to connect. Once connected, it stores the client's memory metrics into both memory and a file.

On the main page, users can see a list of all metrics that were stored in memory.

Upon start-up, the server loads all metrics from the file into memory.

There are two ways to send metrics to the server:
- by using name and value in the URL
- by using a JSON body

The client sends metrics in the following format:
```
ID    string   `json:"id"`
MType string   `json:"type"`
Delta *int64   `json:"delta,omitempty"`
Value *float64 `json:"value,omitempty"`
```
Where:
- ID is the unique metric name
- MType is the metric type (counter or gauge)
- Delta is the counter increment value
- Value is the gauge value

The server handles the following requests:
- POST "/update/{metric}/{key}/{value}" - updates metric with the given key and value
- GET "/value/{metric}/{key}" - returns the value of the given metric and key
- GET "/" - returns current values of all metrics stored in memory
- POST "/update/" - updates metric with the given JSON body
- GET "/value/" - returns current values of the given metric in JSON format.
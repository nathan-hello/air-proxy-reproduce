This repo is for an issue with [air](https://github.com/cosmtrek/air/) where it does not proxy websocket connections.

The site
1. Serves `index.html` which has an htmx websocket form
2. HTMX attempts to connect to /ws-endpoint via websocket
3. Upon seeing this request, /ws-endpoint upgrades the request from HTTP to websocket using [gorilla/websocket](https://github.com/gorilla/websocket)
4. When a websocket message is sent from the client (from the htmx form), /ws-endpoint sees it and sends it to all of the other clients connected
5. (Irrelevant to the issue) HTMX sees the html response and puts it in the bottom of the div of id "response"

The issue is that when Air tries to proxy the websocket connection, it will fail. In this sample repo, it will fail
silently and I won't get any logging in the ChatSocket for loop (where the websocket logic is).

Sometimes I get the error:
`proxy failed to forward the response body, err: http: request method or response status code does not allow body`

Which I assume is because the proxy is trying to parse HTTP requests, but this isn't HTTP, it's websocket.

To reproduce:
1. Clone repo `git clone https://github.com/nathan-hello/air-proxy-reproduce.git`
2. `cd air-proxy-reproduce`
3. Run `air`. The version of air that I am currently running is v1.52.0
4. Open `localhost:8090` in two browser tabs. Neither one can send or receive messages

Expected behavior:
1. Instead of running the program with `air`, use `go run main.go`.
2. Now in the two browser tabs go to `localhost:8080`
3. The two tabs should be able to send/receive messages to each other

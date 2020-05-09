# teleskope

<img src="teleskope-gif.gif"/>

```go
	r.HandleFunc("/list/ns", func(w http.ResponseWriter, r *http.Request) {
		controller.GetNamespaces(h, w, r)
	})
	r.HandleFunc("/list/dep/{ns}", func(w http.ResponseWriter, r *http.Request) {
		controller.GetDeployments(h, w, r)
	})
	r.HandleFunc("/dep/{ns}/{dep}", func(w http.ResponseWriter, r *http.Request) {
		controller.GetDeployment(h, w, r)
	})
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		controller.StreamUpdateds(h, w, r)
	})
```

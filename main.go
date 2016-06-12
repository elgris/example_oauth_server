package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
)

func main() {
	port := flag.Int("port", 8000, "Default port of the service")

	// TODO: implement real storage (memcached or Redis or whatever else DB-based)
	storage := NewTestStorage()

	// TODO: for now only 1 "application ID" is supported
	// need to make it configurable (with admin panel maybe?)
	storage.SetClient("test_client", &osin.DefaultClient{
		Id:          "test_client",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:1323/auth",
	})
	server := osin.NewServer(osin.NewServerConfig(), storage)

	// Authorization code endpoint
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()
		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			// TODO:
			// 1. check if request already comes from authentication website
			// 1.1 if not - redirect to website
			// 1.2. otherwise - authorize
			//
			// for now left deadly simple - each request to client 'test_client' succeeds

			ar.Authorized = true

			server.FinishAuthorizeRequest(resp, r, ar)
		}

		fmt.Println("ERR", resp.InternalError)
		osin.OutputJSON(resp, w, r)
	})

	// Access token endpoint
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			server.FinishAccessRequest(resp, r, ar)
		}

		fmt.Println("ERR", resp.InternalError)

		osin.OutputJSON(resp, w, r)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

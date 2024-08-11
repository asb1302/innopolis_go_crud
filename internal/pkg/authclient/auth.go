package authclient

import (
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

var client *fasthttp.HostClient

func Init(host string, isTLS bool) {
	client = &fasthttp.HostClient{
		Addr:  host,
		IsTLS: isTLS,
	}
}

func ValidateToken(token string) bool {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(buildURI("/get_user_info"))
	req.Header.Set(fasthttp.HeaderAuthorization, token)
	req.Header.SetHost(client.Addr)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return false
	}

	log.Println(resp)
	return resp.StatusCode() == http.StatusOK
}

func buildURI(path string) string {
	protocol := "http://"
	if client.IsTLS {
		protocol = "https://"
	}
	return protocol + client.Addr + path
}

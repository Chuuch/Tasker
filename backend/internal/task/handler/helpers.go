package handler

import "net/http"

const maxRequestBodySize = 1 << 20

func limitRequestBody(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxRequestBodySize,
	)
}

package qfusegate

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"

	"qiniupkg.com/x/rpc.v7"

	. "golang.org/x/net/context"
)

// ---------------------------------------------------------------------------

func fromReader(p unsafe.Pointer, n uintptr, r io.Reader) (err error) {

	b := ((*[1<<30]byte)(p))[:n]
	_, err = io.ReadFull(r, b)
	return
}

func gobResponseError(resp *http.Response) (err error) {

	e := &rpc.ErrorInfo{
		Err:   resp.Header.Get("X-Err"),
		Reqid: resp.Header.Get("X-Reqid"),
		Code:  resp.StatusCode,
	}
	if resp.StatusCode > 299 {
		if resp.ContentLength != 0 {
			ct, ok := resp.Header["Content-Type"]
			if ok && strings.HasPrefix(ct[0], "application/gob") {
				var errno int32
				err2 := fromReader(unsafe.Pointer(&errno), 4, resp.Body)
				if err2 != nil {
					e.Err = err2.Error()
				} else {
					e.Errno = int(errno)
				}
			}
		}
	}
	return e
}

func gobCallRet(ctx Context, ret interface{}, resp *http.Response) (err error) {

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode/100 == 2 {
		if ret != nil && resp.ContentLength != 0 {
			err = gob.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return gobResponseError(resp)
}

// ---------------------------------------------------------------------------

type gobClient struct {
	rpc.Client
}

func (r gobClient) Call(
	ctx Context, ret interface{}, method, url1 string) (err error) {

	resp, err := r.DoRequestWith(ctx, method, url1, "application/gob", nil, 0)
	if err != nil {
		return err
	}
	return gobCallRet(ctx, ret, resp)
}

func (r gobClient) CallWithGob(
	ctx Context, ret interface{}, method, url1 string, params interface{}) (err error) {

	var b bytes.Buffer
	err = gob.NewEncoder(&b).Encode(params)
	if err != nil {
		return err
	}

	resp, err := r.DoRequestWith(ctx, method, url1, "application/gob", &b, b.Len())
	if err != nil {
		return err
	}
	return gobCallRet(ctx, ret, resp)
}

// ---------------------------------------------------------------------------


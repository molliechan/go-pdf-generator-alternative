package gotenberg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/molliechan/go-pdf-generator-alternative/internal/helper"
)

const (
	paperWidth        string = "paperWidth"
	paperHeight       string = "paperHeight"
	marginTop         string = "marginTop"
	marginBottom      string = "marginBottom"
	marginLeft        string = "marginLeft"
	marginRight       string = "marginRight"
	preferCssPageSize string = "preferCssPageSize"
	printBackground   string = "printBackground"
	landscape         string = "landscape"
	scale             string = "scale"
	nativePageRanges  string = "nativePageRanges"
)

type Request struct {
	URL         string
	HttpHeaders map[string]string
	FormValues  map[string]string
	FormFiles   map[string][]byte
}

func NewRequest(url string) *Request {
	return &Request{
		URL:         url,
		HttpHeaders: make(map[string]string),
		FormValues:  make(map[string]string),
		FormFiles:   make(map[string][]byte),
	}
}

func (req *Request) SetFormValue(key string, value string) {
	req.FormValues[key] = value
}

func (req *Request) SetFormIndexFile(data []byte) {
	req.FormFiles["index.html"] = data
}

func (req *Request) SetMargins(margins [4]float64) {
	req.SetFormValue(marginTop, fmt.Sprintf("%f", margins[0]))
	req.SetFormValue(marginBottom, fmt.Sprintf("%f", margins[1]))
	req.SetFormValue(marginLeft, fmt.Sprintf("%f", margins[2]))
	req.SetFormValue(marginRight, fmt.Sprintf("%f", margins[3]))
}

type Client struct {
	Hostname   string
	HTTPClient *http.Client
}

// Store creates the resulting PDF to given destination.
func (c *Client) Store(req *Request, dest string) error {
	return c.storeContext(context.Background(), req, dest)
}

// StoreContext creates the resulting PDF to given destination.
// The created HTTP request can be canceled by the passed context.
func (c *Client) storeContext(ctx context.Context, req *Request, dest string) error {

	resp, err := c.postContext(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))

		return errors.New("failed to generate the result PDF")
	}
	return helper.WriteNewFile(dest, resp.Body)
}

// PostContext sends a request to the Gotenberg API
// and returns the response.
// The created HTTP request can be canceled by the passed context.
func (c *Client) postContext(ctx context.Context, req *Request) (*http.Response, error) {
	body, contentType, err := multipartForm(req)
	if err != nil {
		return nil, err
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	URL := fmt.Sprintf("%s%s", c.Hostname, req.URL)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", contentType)
	for key, value := range req.HttpHeaders {
		httpReq.Header.Set(key, value)
	}
	resp, err := c.HTTPClient.Do(httpReq) /* #nosec */
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func multipartForm(req *Request) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	// create form files
	for filename, data := range req.FormFiles {
		reader := ioutil.NopCloser(bytes.NewReader(data))

		part, err := writer.CreateFormFile("files", filename)
		if err != nil {
			return nil, "", fmt.Errorf("%s: creating form file: %v", filename, err)
		}

		_, err = io.Copy(part, reader)
		if err != nil {
			return nil, "", fmt.Errorf("%s: copying data: %v", filename, err)
		}
	}

	// create form values
	for key, value := range req.FormValues {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", fmt.Errorf("%s: writing form field: %v", key, err)
		}
	}

	return body, writer.FormDataContentType(), nil
}

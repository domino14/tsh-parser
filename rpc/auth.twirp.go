// Code generated by protoc-gen-twirp v8.1.2, DO NOT EDIT.
// source: auth.proto

package proto

import context "context"
import fmt "fmt"
import http "net/http"
import ioutil "io/ioutil"
import json "encoding/json"
import strconv "strconv"
import strings "strings"

import protojson "google.golang.org/protobuf/encoding/protojson"
import proto "google.golang.org/protobuf/proto"
import twirp "github.com/twitchtv/twirp"
import ctxsetters "github.com/twitchtv/twirp/ctxsetters"

// Version compatibility assertion.
// If the constant is not defined in the package, that likely means
// the package needs to be updated to work with this generated code.
// See https://twitchtv.github.io/twirp/docs/version_matrix.html
const _ = twirp.TwirpPackageMinVersion_8_1_0

// ===============================
// AuthenticationService Interface
// ===============================

type AuthenticationService interface {
	RegisterUser(context.Context, *NewUserRequest) (*NewUserResponse, error)

	GetJWT(context.Context, *JWTRequest) (*JWTResponse, error)

	WhoAmI(context.Context, *WhoAmIRequest) (*User, error)
}

// =====================================
// AuthenticationService Protobuf Client
// =====================================

type authenticationServiceProtobufClient struct {
	client      HTTPClient
	urls        [3]string
	interceptor twirp.Interceptor
	opts        twirp.ClientOptions
}

// NewAuthenticationServiceProtobufClient creates a Protobuf client that implements the AuthenticationService interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewAuthenticationServiceProtobufClient(baseURL string, client HTTPClient, opts ...twirp.ClientOption) AuthenticationService {
	if c, ok := client.(*http.Client); ok {
		client = withoutRedirects(c)
	}

	clientOpts := twirp.ClientOptions{}
	for _, o := range opts {
		o(&clientOpts)
	}

	// Using ReadOpt allows backwards and forwads compatibility with new options in the future
	literalURLs := false
	_ = clientOpts.ReadOpt("literalURLs", &literalURLs)
	var pathPrefix string
	if ok := clientOpts.ReadOpt("pathPrefix", &pathPrefix); !ok {
		pathPrefix = "/twirp" // default prefix
	}

	// Build method URLs: <baseURL>[<prefix>]/<package>.<Service>/<Method>
	serviceURL := sanitizeBaseURL(baseURL)
	serviceURL += baseServicePath(pathPrefix, "tshparser", "AuthenticationService")
	urls := [3]string{
		serviceURL + "RegisterUser",
		serviceURL + "GetJWT",
		serviceURL + "WhoAmI",
	}

	return &authenticationServiceProtobufClient{
		client:      client,
		urls:        urls,
		interceptor: twirp.ChainInterceptors(clientOpts.Interceptors...),
		opts:        clientOpts,
	}
}

func (c *authenticationServiceProtobufClient) RegisterUser(ctx context.Context, in *NewUserRequest) (*NewUserResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "tshparser")
	ctx = ctxsetters.WithServiceName(ctx, "AuthenticationService")
	ctx = ctxsetters.WithMethodName(ctx, "RegisterUser")
	caller := c.callRegisterUser
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *NewUserRequest) (*NewUserResponse, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*NewUserRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*NewUserRequest) when calling interceptor")
					}
					return c.callRegisterUser(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*NewUserResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*NewUserResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *authenticationServiceProtobufClient) callRegisterUser(ctx context.Context, in *NewUserRequest) (*NewUserResponse, error) {
	out := new(NewUserResponse)
	ctx, err := doProtobufRequest(ctx, c.client, c.opts.Hooks, c.urls[0], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

func (c *authenticationServiceProtobufClient) GetJWT(ctx context.Context, in *JWTRequest) (*JWTResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "tshparser")
	ctx = ctxsetters.WithServiceName(ctx, "AuthenticationService")
	ctx = ctxsetters.WithMethodName(ctx, "GetJWT")
	caller := c.callGetJWT
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *JWTRequest) (*JWTResponse, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*JWTRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*JWTRequest) when calling interceptor")
					}
					return c.callGetJWT(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*JWTResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*JWTResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *authenticationServiceProtobufClient) callGetJWT(ctx context.Context, in *JWTRequest) (*JWTResponse, error) {
	out := new(JWTResponse)
	ctx, err := doProtobufRequest(ctx, c.client, c.opts.Hooks, c.urls[1], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

func (c *authenticationServiceProtobufClient) WhoAmI(ctx context.Context, in *WhoAmIRequest) (*User, error) {
	ctx = ctxsetters.WithPackageName(ctx, "tshparser")
	ctx = ctxsetters.WithServiceName(ctx, "AuthenticationService")
	ctx = ctxsetters.WithMethodName(ctx, "WhoAmI")
	caller := c.callWhoAmI
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *WhoAmIRequest) (*User, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*WhoAmIRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*WhoAmIRequest) when calling interceptor")
					}
					return c.callWhoAmI(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*User)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*User) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *authenticationServiceProtobufClient) callWhoAmI(ctx context.Context, in *WhoAmIRequest) (*User, error) {
	out := new(User)
	ctx, err := doProtobufRequest(ctx, c.client, c.opts.Hooks, c.urls[2], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

// =================================
// AuthenticationService JSON Client
// =================================

type authenticationServiceJSONClient struct {
	client      HTTPClient
	urls        [3]string
	interceptor twirp.Interceptor
	opts        twirp.ClientOptions
}

// NewAuthenticationServiceJSONClient creates a JSON client that implements the AuthenticationService interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewAuthenticationServiceJSONClient(baseURL string, client HTTPClient, opts ...twirp.ClientOption) AuthenticationService {
	if c, ok := client.(*http.Client); ok {
		client = withoutRedirects(c)
	}

	clientOpts := twirp.ClientOptions{}
	for _, o := range opts {
		o(&clientOpts)
	}

	// Using ReadOpt allows backwards and forwads compatibility with new options in the future
	literalURLs := false
	_ = clientOpts.ReadOpt("literalURLs", &literalURLs)
	var pathPrefix string
	if ok := clientOpts.ReadOpt("pathPrefix", &pathPrefix); !ok {
		pathPrefix = "/twirp" // default prefix
	}

	// Build method URLs: <baseURL>[<prefix>]/<package>.<Service>/<Method>
	serviceURL := sanitizeBaseURL(baseURL)
	serviceURL += baseServicePath(pathPrefix, "tshparser", "AuthenticationService")
	urls := [3]string{
		serviceURL + "RegisterUser",
		serviceURL + "GetJWT",
		serviceURL + "WhoAmI",
	}

	return &authenticationServiceJSONClient{
		client:      client,
		urls:        urls,
		interceptor: twirp.ChainInterceptors(clientOpts.Interceptors...),
		opts:        clientOpts,
	}
}

func (c *authenticationServiceJSONClient) RegisterUser(ctx context.Context, in *NewUserRequest) (*NewUserResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "tshparser")
	ctx = ctxsetters.WithServiceName(ctx, "AuthenticationService")
	ctx = ctxsetters.WithMethodName(ctx, "RegisterUser")
	caller := c.callRegisterUser
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *NewUserRequest) (*NewUserResponse, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*NewUserRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*NewUserRequest) when calling interceptor")
					}
					return c.callRegisterUser(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*NewUserResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*NewUserResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *authenticationServiceJSONClient) callRegisterUser(ctx context.Context, in *NewUserRequest) (*NewUserResponse, error) {
	out := new(NewUserResponse)
	ctx, err := doJSONRequest(ctx, c.client, c.opts.Hooks, c.urls[0], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

func (c *authenticationServiceJSONClient) GetJWT(ctx context.Context, in *JWTRequest) (*JWTResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "tshparser")
	ctx = ctxsetters.WithServiceName(ctx, "AuthenticationService")
	ctx = ctxsetters.WithMethodName(ctx, "GetJWT")
	caller := c.callGetJWT
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *JWTRequest) (*JWTResponse, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*JWTRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*JWTRequest) when calling interceptor")
					}
					return c.callGetJWT(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*JWTResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*JWTResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *authenticationServiceJSONClient) callGetJWT(ctx context.Context, in *JWTRequest) (*JWTResponse, error) {
	out := new(JWTResponse)
	ctx, err := doJSONRequest(ctx, c.client, c.opts.Hooks, c.urls[1], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

func (c *authenticationServiceJSONClient) WhoAmI(ctx context.Context, in *WhoAmIRequest) (*User, error) {
	ctx = ctxsetters.WithPackageName(ctx, "tshparser")
	ctx = ctxsetters.WithServiceName(ctx, "AuthenticationService")
	ctx = ctxsetters.WithMethodName(ctx, "WhoAmI")
	caller := c.callWhoAmI
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *WhoAmIRequest) (*User, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*WhoAmIRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*WhoAmIRequest) when calling interceptor")
					}
					return c.callWhoAmI(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*User)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*User) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *authenticationServiceJSONClient) callWhoAmI(ctx context.Context, in *WhoAmIRequest) (*User, error) {
	out := new(User)
	ctx, err := doJSONRequest(ctx, c.client, c.opts.Hooks, c.urls[2], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

// ====================================
// AuthenticationService Server Handler
// ====================================

type authenticationServiceServer struct {
	AuthenticationService
	interceptor      twirp.Interceptor
	hooks            *twirp.ServerHooks
	pathPrefix       string // prefix for routing
	jsonSkipDefaults bool   // do not include unpopulated fields (default values) in the response
	jsonCamelCase    bool   // JSON fields are serialized as lowerCamelCase rather than keeping the original proto names
}

// NewAuthenticationServiceServer builds a TwirpServer that can be used as an http.Handler to handle
// HTTP requests that are routed to the right method in the provided svc implementation.
// The opts are twirp.ServerOption modifiers, for example twirp.WithServerHooks(hooks).
func NewAuthenticationServiceServer(svc AuthenticationService, opts ...interface{}) TwirpServer {
	serverOpts := newServerOpts(opts)

	// Using ReadOpt allows backwards and forwads compatibility with new options in the future
	jsonSkipDefaults := false
	_ = serverOpts.ReadOpt("jsonSkipDefaults", &jsonSkipDefaults)
	jsonCamelCase := false
	_ = serverOpts.ReadOpt("jsonCamelCase", &jsonCamelCase)
	var pathPrefix string
	if ok := serverOpts.ReadOpt("pathPrefix", &pathPrefix); !ok {
		pathPrefix = "/twirp" // default prefix
	}

	return &authenticationServiceServer{
		AuthenticationService: svc,
		hooks:                 serverOpts.Hooks,
		interceptor:           twirp.ChainInterceptors(serverOpts.Interceptors...),
		pathPrefix:            pathPrefix,
		jsonSkipDefaults:      jsonSkipDefaults,
		jsonCamelCase:         jsonCamelCase,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *authenticationServiceServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// handleRequestBodyError is used to handle error when the twirp server cannot read request
func (s *authenticationServiceServer) handleRequestBodyError(ctx context.Context, resp http.ResponseWriter, msg string, err error) {
	if context.Canceled == ctx.Err() {
		s.writeError(ctx, resp, twirp.NewError(twirp.Canceled, "failed to read request: context canceled"))
		return
	}
	if context.DeadlineExceeded == ctx.Err() {
		s.writeError(ctx, resp, twirp.NewError(twirp.DeadlineExceeded, "failed to read request: deadline exceeded"))
		return
	}
	s.writeError(ctx, resp, twirp.WrapError(malformedRequestError(msg), err))
}

// AuthenticationServicePathPrefix is a convenience constant that may identify URL paths.
// Should be used with caution, it only matches routes generated by Twirp Go clients,
// with the default "/twirp" prefix and default CamelCase service and method names.
// More info: https://twitchtv.github.io/twirp/docs/routing.html
const AuthenticationServicePathPrefix = "/twirp/tshparser.AuthenticationService/"

func (s *authenticationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "tshparser")
	ctx = ctxsetters.WithServiceName(ctx, "AuthenticationService")
	ctx = ctxsetters.WithResponseWriter(ctx, resp)

	var err error
	ctx, err = callRequestReceived(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	if req.Method != "POST" {
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}

	// Verify path format: [<prefix>]/<package>.<Service>/<Method>
	prefix, pkgService, method := parseTwirpPath(req.URL.Path)
	if pkgService != "tshparser.AuthenticationService" {
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}
	if prefix != s.pathPrefix {
		msg := fmt.Sprintf("invalid path prefix %q, expected %q, on path %q", prefix, s.pathPrefix, req.URL.Path)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}

	switch method {
	case "RegisterUser":
		s.serveRegisterUser(ctx, resp, req)
		return
	case "GetJWT":
		s.serveGetJWT(ctx, resp, req)
		return
	case "WhoAmI":
		s.serveWhoAmI(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}
}

func (s *authenticationServiceServer) serveRegisterUser(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveRegisterUserJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveRegisterUserProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *authenticationServiceServer) serveRegisterUserJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "RegisterUser")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	d := json.NewDecoder(req.Body)
	rawReqBody := json.RawMessage{}
	if err := d.Decode(&rawReqBody); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}
	reqContent := new(NewUserRequest)
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err = unmarshaler.Unmarshal(rawReqBody, reqContent); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}

	handler := s.AuthenticationService.RegisterUser
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *NewUserRequest) (*NewUserResponse, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*NewUserRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*NewUserRequest) when calling interceptor")
					}
					return s.AuthenticationService.RegisterUser(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*NewUserResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*NewUserResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *NewUserResponse
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *NewUserResponse and nil error while calling RegisterUser. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	marshaler := &protojson.MarshalOptions{UseProtoNames: !s.jsonCamelCase, EmitUnpopulated: !s.jsonSkipDefaults}
	respBytes, err := marshaler.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal json response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)

	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *authenticationServiceServer) serveRegisterUserProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "RegisterUser")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.handleRequestBodyError(ctx, resp, "failed to read request body", err)
		return
	}
	reqContent := new(NewUserRequest)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		s.writeError(ctx, resp, malformedRequestError("the protobuf request could not be decoded"))
		return
	}

	handler := s.AuthenticationService.RegisterUser
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *NewUserRequest) (*NewUserResponse, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*NewUserRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*NewUserRequest) when calling interceptor")
					}
					return s.AuthenticationService.RegisterUser(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*NewUserResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*NewUserResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *NewUserResponse
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *NewUserResponse and nil error while calling RegisterUser. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal proto response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *authenticationServiceServer) serveGetJWT(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveGetJWTJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveGetJWTProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *authenticationServiceServer) serveGetJWTJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetJWT")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	d := json.NewDecoder(req.Body)
	rawReqBody := json.RawMessage{}
	if err := d.Decode(&rawReqBody); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}
	reqContent := new(JWTRequest)
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err = unmarshaler.Unmarshal(rawReqBody, reqContent); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}

	handler := s.AuthenticationService.GetJWT
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *JWTRequest) (*JWTResponse, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*JWTRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*JWTRequest) when calling interceptor")
					}
					return s.AuthenticationService.GetJWT(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*JWTResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*JWTResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *JWTResponse
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *JWTResponse and nil error while calling GetJWT. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	marshaler := &protojson.MarshalOptions{UseProtoNames: !s.jsonCamelCase, EmitUnpopulated: !s.jsonSkipDefaults}
	respBytes, err := marshaler.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal json response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)

	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *authenticationServiceServer) serveGetJWTProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetJWT")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.handleRequestBodyError(ctx, resp, "failed to read request body", err)
		return
	}
	reqContent := new(JWTRequest)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		s.writeError(ctx, resp, malformedRequestError("the protobuf request could not be decoded"))
		return
	}

	handler := s.AuthenticationService.GetJWT
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *JWTRequest) (*JWTResponse, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*JWTRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*JWTRequest) when calling interceptor")
					}
					return s.AuthenticationService.GetJWT(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*JWTResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*JWTResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *JWTResponse
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *JWTResponse and nil error while calling GetJWT. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal proto response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *authenticationServiceServer) serveWhoAmI(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveWhoAmIJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveWhoAmIProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *authenticationServiceServer) serveWhoAmIJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "WhoAmI")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	d := json.NewDecoder(req.Body)
	rawReqBody := json.RawMessage{}
	if err := d.Decode(&rawReqBody); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}
	reqContent := new(WhoAmIRequest)
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err = unmarshaler.Unmarshal(rawReqBody, reqContent); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}

	handler := s.AuthenticationService.WhoAmI
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *WhoAmIRequest) (*User, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*WhoAmIRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*WhoAmIRequest) when calling interceptor")
					}
					return s.AuthenticationService.WhoAmI(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*User)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*User) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *User
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *User and nil error while calling WhoAmI. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	marshaler := &protojson.MarshalOptions{UseProtoNames: !s.jsonCamelCase, EmitUnpopulated: !s.jsonSkipDefaults}
	respBytes, err := marshaler.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal json response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)

	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *authenticationServiceServer) serveWhoAmIProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "WhoAmI")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.handleRequestBodyError(ctx, resp, "failed to read request body", err)
		return
	}
	reqContent := new(WhoAmIRequest)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		s.writeError(ctx, resp, malformedRequestError("the protobuf request could not be decoded"))
		return
	}

	handler := s.AuthenticationService.WhoAmI
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *WhoAmIRequest) (*User, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*WhoAmIRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*WhoAmIRequest) when calling interceptor")
					}
					return s.AuthenticationService.WhoAmI(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*User)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*User) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *User
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *User and nil error while calling WhoAmI. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal proto response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *authenticationServiceServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor1, 0
}

func (s *authenticationServiceServer) ProtocGenTwirpVersion() string {
	return "v8.1.2"
}

// PathPrefix returns the base service path, in the form: "/<prefix>/<package>.<Service>/"
// that is everything in a Twirp route except for the <Method>. This can be used for routing,
// for example to identify the requests that are targeted to this service in a mux.
func (s *authenticationServiceServer) PathPrefix() string {
	return baseServicePath(s.pathPrefix, "tshparser", "AuthenticationService")
}

var twirpFileDescriptor1 = []byte{
	// 326 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x4f, 0x4f, 0xf2, 0x40,
	0x10, 0xc6, 0xd3, 0xf7, 0x55, 0x84, 0x11, 0x24, 0x6e, 0xc4, 0x94, 0x9e, 0x48, 0x39, 0x88, 0x97,
	0x36, 0x8a, 0x89, 0x37, 0x13, 0x48, 0x8c, 0xca, 0xc1, 0x43, 0xc5, 0x90, 0x18, 0x13, 0xb2, 0x94,
	0x09, 0xbb, 0xd1, 0x76, 0xeb, 0xce, 0x56, 0xbe, 0xa4, 0x1f, 0xca, 0xb4, 0x94, 0x0a, 0x46, 0x2f,
	0x9e, 0x36, 0xcf, 0xfc, 0xf9, 0xcd, 0xcc, 0x93, 0x05, 0xe0, 0xa9, 0x11, 0x5e, 0xa2, 0x95, 0x51,
	0xac, 0x66, 0x48, 0x24, 0x5c, 0x13, 0x6a, 0x77, 0x08, 0x07, 0xf7, 0xb8, 0x7c, 0x24, 0xd4, 0x01,
	0xbe, 0xa5, 0x48, 0x86, 0x1d, 0xc1, 0x2e, 0x46, 0x5c, 0xbe, 0xda, 0x56, 0xc7, 0xea, 0xd5, 0x82,
	0x95, 0x60, 0x0e, 0x54, 0x13, 0x4e, 0xb4, 0x54, 0x7a, 0x6e, 0xff, 0xcb, 0x13, 0xa5, 0x76, 0x0f,
	0xa1, 0x59, 0x32, 0x28, 0x51, 0x31, 0xa1, 0x7b, 0x05, 0x30, 0x9a, 0x8c, 0xff, 0x8e, 0xec, 0xc2,
	0x7e, 0xde, 0xbf, 0xc2, 0x65, 0x00, 0xa3, 0x5e, 0x30, 0x5e, 0x03, 0x72, 0xe1, 0x3e, 0xc3, 0x4e,
	0x36, 0xf4, 0x17, 0x7c, 0x17, 0x1a, 0x6b, 0xdc, 0x54, 0x70, 0x12, 0xc5, 0x8c, 0xfa, 0x3a, 0x78,
	0xcb, 0x49, 0xb0, 0x36, 0x54, 0x25, 0x4d, 0xf9, 0x3c, 0x92, 0xb1, 0xfd, 0xbf, 0x63, 0xf5, 0xaa,
	0xc1, 0x9e, 0xa4, 0x41, 0x26, 0xdd, 0x26, 0x34, 0x26, 0x42, 0x0d, 0xa2, 0xbb, 0xe2, 0x8a, 0xf3,
	0x0f, 0x0b, 0x5a, 0x83, 0xd4, 0x08, 0x8c, 0x8d, 0x0c, 0xb9, 0x91, 0x2a, 0x7e, 0x40, 0xfd, 0x2e,
	0x43, 0x64, 0xd7, 0x50, 0x0f, 0x70, 0x21, 0xc9, 0xa0, 0xce, 0x17, 0x6a, 0x7b, 0xa5, 0xc1, 0xde,
	0xb6, 0xbb, 0x8e, 0xf3, 0x53, 0xaa, 0xb8, 0xf2, 0x12, 0x2a, 0x37, 0x68, 0x46, 0x93, 0x31, 0x6b,
	0x6d, 0x54, 0x7d, 0xf9, 0xe8, 0x1c, 0x7f, 0x0f, 0x17, 0x8d, 0x7d, 0xa8, 0xac, 0x56, 0x65, 0xf6,
	0x46, 0xc5, 0xd6, 0xf6, 0x4e, 0x73, 0x23, 0x93, 0x4d, 0x1d, 0x9e, 0x3e, 0x9d, 0x2c, 0xa4, 0x11,
	0xe9, 0xcc, 0x0b, 0x55, 0xe4, 0xcf, 0x55, 0x24, 0x63, 0x75, 0x76, 0xe1, 0x97, 0x55, 0xbe, 0x4e,
	0x42, 0x3f, 0xff, 0x2f, 0xb3, 0x4a, 0xfe, 0xf4, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x98, 0xff,
	0xbe, 0xbf, 0x44, 0x02, 0x00, 0x00,
}

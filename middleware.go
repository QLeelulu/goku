package goku

// middlewarer interface
// execute order: OnBeginRequest -> OnBeginMvcHandle -> {controller} -> OnEndMvcHandle -> OnEndRequest
// notice:
//		OnBeginRequest & OnEndRequest: All requests will be through these
//		OnBeginMvcHandle & OnEndMvcHandle: not matched route & static file are not through these
type Middlewarer interface {
	OnBeginRequest(ctx *HttpContext) (ActionResulter, error)
	OnBeginMvcHandle(ctx *HttpContext) (ActionResulter, error)
	OnEndMvcHandle(ctx *HttpContext) (ActionResulter, error)
	OnEndRequest(ctx *HttpContext) (ActionResulter, error)
}

// middleware handler, handle the middleware how tu execute
type MiddlewareHandler interface {
	BeginRequest(ctx *HttpContext) (ar ActionResulter, err error)
	BeginMvcHandle(ctx *HttpContext) (ar ActionResulter, err error)
	EndMvcHandle(ctx *HttpContext) (ar ActionResulter, err error)
	EndRequest(ctx *HttpContext) (ar ActionResulter, err error)
}

// the defaultmiddleware handler
type DefaultMiddlewareHandle struct {
	Middlewares []Middlewarer
}

func (mh *DefaultMiddlewareHandle) AddMiddleware(mw Middlewarer) {
	mh.Middlewares = append(mh.Middlewares, mw)
}

func (mh *DefaultMiddlewareHandle) BeginRequest(ctx *HttpContext) (ar ActionResulter, err error) {
	for _, mw := range mh.Middlewares {
		ar, err = mw.OnBeginRequest(ctx)
		if err != nil || ar != nil {
			return
		}
	}
	return
}

func (mh *DefaultMiddlewareHandle) EndRequest(ctx *HttpContext) (ar ActionResulter, err error) {
	for _, mw := range mh.Middlewares {
		ar, err = mw.OnEndRequest(ctx)
		if err != nil || ar != nil {
			return
		}
	}
	return
}

func (mh *DefaultMiddlewareHandle) BeginMvcHandle(ctx *HttpContext) (ar ActionResulter, err error) {
	for _, mw := range mh.Middlewares {
		ar, err = mw.OnBeginMvcHandle(ctx)
		if err != nil || ar != nil {
			return
		}
	}
	return
}

func (mh *DefaultMiddlewareHandle) EndMvcHandle(ctx *HttpContext) (ar ActionResulter, err error) {
	for _, mw := range mh.Middlewares {
		ar, err = mw.OnEndMvcHandle(ctx)
		if err != nil || ar != nil {
			return
		}
	}
	return
}

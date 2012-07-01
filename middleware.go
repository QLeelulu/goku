package goku

type Middlewarer interface {
	OnBeginRequest(ctx *HttpContext) (ActionResulter, error)
	OnBeginMvcHandle(ctx *HttpContext) (ActionResulter, error)
	//OnEndMvcHandle(ctx *HttpContext) (ActionResulter, error)
	OnEndRequest(ctx *HttpContext) (ActionResulter, error)
}

type MiddlewareHandler interface {
	BeginRequest(ctx *HttpContext) (ar ActionResulter, err error)
	BeginMvcHandle(ctx *HttpContext) (ar ActionResulter, err error)
	//EndMvcHandle(ctx *HttpContext) (ar ActionResulter, err error)
	EndRequest(ctx *HttpContext) (ar ActionResulter, err error)
}

type DefaultMiddlewareHandle struct {
	Middlewares []Middlewarer

	// beginRequests   []func(http.ResponseWriter, *http.Request) error
	// beginMvcHandles []func(http.ResponseWriter, *http.Request) error
	// endMvcHandles   []func(http.ResponseWriter, *http.Request) error
	// endRequests     []func(http.ResponseWriter, *http.Request) error

	// inited bool
}

// func (mh DefaultMiddlewareHandle) init() {
// 	if mh.inited {
// 		return
// 	}
// 	mh.beginRequests = new([]func(http.ResponseWriter, *http.Request) error)
// 	mh.beginMvcHandles = new([]func(http.ResponseWriter, *http.Request) error)
// 	mh.endMvcHandles = new([]func(http.ResponseWriter, *http.Request) error)
// 	mh.endRequests = new([]func(http.ResponseWriter, *http.Request) error)

// 	for _, mw := range mh.Middlewares {
// 		mh.beginRequests = append(mh.beginRequests, mw.OnBeginRequest)
// 		mh.beginMvcHandles = append(mh.beginMvcHandles, mw.OnBeginMvcHandle)
// 		mh.endMvcHandles = append(mh.endMvcHandles, mw.OnEndMvcHandle)
// 		mh.endRequests = append(mh.endRequests, mw.OnEndRequest)
// 	}
// 	mh.inited = true
// }

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

// func (mh *DefaultMiddlewareHandle) EndMvcHandle(ctx *HttpContext) (ar ActionResulter, err error) {
// 	for _, mw := range mh.Middlewares {
// 		ar, err = mw.OnEndMvcHandle(ctx)
// 		if err != nil || ar != nil {
// 			return
// 		}
// 	}
// 	return
// }

package goku

// Order of the filters execution is: 
//      1. OnActionExecuting
//      2. -> Execute Action -> return ActionResulter
//      3. OnActionExecuted
//      4. OnResultExecuting
//      5. -> ActionResulter.ExecuteResult
//      6. OnResultExecuted
type Filter interface {
	OnActionExecuting(ctx *HttpContext) (ActionResulter, error)
	OnActionExecuted(ctx *HttpContext) (ActionResulter, error)
	OnResultExecuting(ctx *HttpContext) (ActionResulter, error)
	OnResultExecuted(ctx *HttpContext) (ActionResulter, error)
}

func runFilterActionExecuting(ctx *HttpContext, filters []Filter) (ar ActionResulter, err error) {
	for _, f := range filters {
		ar, err = f.OnActionExecuting(ctx)
		if ctx.Canceled || ar != nil || err != nil {
			return
		}
	}
	return
}

func runFilterActionExecuted(ctx *HttpContext, filters []Filter) (ar ActionResulter, err error) {
	for _, f := range filters {
		ar, err = f.OnActionExecuted(ctx)
		if ctx.Canceled || ar != nil || err != nil {
			return
		}
	}
	return
}

func runFilterResultExecuting(ctx *HttpContext, filters []Filter) (ar ActionResulter, err error) {
	for _, f := range filters {
		ar, err = f.OnResultExecuting(ctx)
		if ctx.Canceled || ar != nil || err != nil {
			return
		}
	}
	return
}

func runFilterResultExecuted(ctx *HttpContext, filters []Filter) (ar ActionResulter, err error) {
	for _, f := range filters {
		ar, err = f.OnResultExecuted(ctx)
		if ctx.Canceled || ar != nil || err != nil {
			return
		}
	}
	return
}

package router

type middlewareChain struct {
	handlers []HandlerFunc
}

func (c *middlewareChain) execute(ctx *Context) error {
	// Execute middleware chain without recursion
	for _, h := range c.handlers {
		if err := h(ctx); err != nil {
			return err
		}
	}
	return nil
}

package gin

type H map[string]interface{}

type Context struct {
}

func (c *Context) BindJSON(i interface{}) error {
	return nil
}

func (c *Context) JSON(s int, i interface{}) {

}

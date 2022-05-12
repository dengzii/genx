package main

const (
	ginPkg = "github.com/gin-gonic/gin"
)

type ApiHandler struct {
	receiver      *GoField
	context       *GoField
	requestParam  *GoField
	responseParam *GoField
	err           *GoField

	*GoFunc
}

func NewApiHandler(fn *GoFunc) (*ApiHandler, error) {
	ret := &ApiHandler{
		GoFunc: fn,
	}

	for _, goField := range fn.getParamList() {
		if goField.typeName == "Context" && goField.pkgPath == ginPkg {
			ret.context = goField
			continue
		} else {
			ret.requestParam = goField
		}
	}

	for _, field := range fn.getResultList() {
		if field.typeName == "error" {
			ret.err = field
		} else {
			ret.responseParam = field
		}
	}
	return ret, nil
}

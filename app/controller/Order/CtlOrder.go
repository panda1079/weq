package controller

import (
	"core/library"
	"fmt"
)

type CtlOrder struct {
	out map[string]string
}

func (r *CtlOrder) OrderList(CH library.HttpInfo) {

	postData := map[string]string{
		"a": CH.R("a", "a"),
		"b": CH.R("b", "b"),
		"c": CH.R("c", "c"),
		"d": CH.R("d", "d"),
		"e": CH.R("e", "e"),
		"f": CH.R("f", "f"),
	}

	fmt.Print(postData)

	CH.OutJson(postData)

}

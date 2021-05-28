package main

import (
	"fmt"
	"reflect"
	"syscall/js"
)

type Element struct {
	elementType string
	attributes  map[string]string
	children    []Element
}

func createElement(elementType string, attributes map[string]string, children []interface{}) Element {
	// We initiate the return element
	result := Element{elementType: elementType, attributes: attributes}
	// We build the children element list recursively
	childrenList := make([]Element, 0)
	for _, child := range children {
		if reflect.TypeOf(child).String() == "string" {
			attributes := make(map[string]string)
			attributes["nodeValue"] = fmt.Sprintf("%v", child)
			childrenList = append(childrenList, Element{"TEXT_ELEMENT", attributes, make([]Element, 0)})
		} else {
			childElement, _ := child.(Element)
			childrenList = append(childrenList, childElement)
		}
	}
	result.children = childrenList
	return result
}

func printArgs() js.Func {
	printArgsFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println(this)
		objectPrototype := js.Global().Get("Object")
		keys := objectPrototype.Call("keys", args[0])
		for i := 0; i < keys.Length(); i++ {
			fmt.Println(keys.Index(i))
		}
		fmt.Println(args[0])
		return 0
	})
	return printArgsFunc
}
func getJSObjectKeys(object js.Value) []string {
	objectPrototype := js.Global().Get("Object")
	keys := objectPrototype.Call("keys", object)
	keyList := make([]string, 0)
	for i := 0; i < keys.Length(); i++ {
		keyList = append(keyList, keys.Index(i))
	}
	return keyList
}
func getMapStringStringFromJSObject(object js.Value) map[string]string {
	keys := getJSObjectKeys(object)
	mapStrStr := make(map[string]string)
	for _, key := range keys {
		mapStrStr[key] = object.Get(key)
	}
	return mapStrStr
}
func JSCreateElement() js.Func {
	createElementFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		elementType := args[0]
		attributes := getMapStringStringFromJSObject(args[1])
		createElement(elementType, attributes, make([]interface{}, 0))
		return 0
	})
	return createElementFunc
}
func main() {
	fmt.Println("Go Web Assembly Running")
	js.Global().Set("printArgs", printArgs())
	<-make(chan bool)
}

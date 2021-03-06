package controllers

type HelloController struct {
}

func (h *HelloController) Get() string {
	return "this is HelloController !"
}

//生成的请求地址是sayhello/{param:string}
func (h *HelloController) GetSayhelloBy(name string) string {
	//name := cxt.FormValue("name")
	return "hello " + name + " !"
}

//生成url的规则是，根据驼峰命名法将方法名先分成一个数组，判断第一个值是否合规（Get、Post、Any、All...)，
//判断是否使用了By，by后边会直接跟着参数
func (h *HelloController) PostSayhelloBy(name string) string {
	return "hello " + name + " !"
}

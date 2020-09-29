package controllors

import (
	"errors"
	"github.com/kataras/iris/v12"
	"overview/datamodels"
	"overview/services"
)

type MovieControllor struct {
	Service services.MovieService
}

func (c *MovieControllor) Get()(result []datamodels.Movie){
	return c.Service.GetAll()
}
// GetBy 返回一个 movie
// 演示:
// curl -i http://localhost:8080/movies/1
func (c *MovieControllor) GetBy(id int64) (datamodels.Movie, bool){
	return c.Service.GetById(id)
}
// PutBy 更新一个movie
// 演示:
// curl -i -X PUT -F "genre=Thriller" -F "poster=@/Users/kataras/Downloads/out.gif" http://localhost:8080/movies/1
func (c *MovieControllor) PutBy(cxt iris.Context, id int64)(datamodels.Movie,error){
	file, info, err := cxt.FormFile("poster")
	if err != nil{
		return datamodels.Movie{},errors.New("failed due form file 'poster' missing")
	}
	file.Close()

	poster := info.Filename
	genre := cxt.FormValue("genre")

	return c.Service.UpdatePosterAndGenreById(id,poster,genre)
}
// DeleteBy 删除一个 movie
// 演示:
// curl -i -X DELETE -u admin:password http://localhost:8080/movies/1
func (c *MovieControllor) DeleteBy(id int64) interface{}{
	wasDel := c.Service.DeleteById(id)
	if wasDel {
		return iris.Map{"deleteId":id}
	}
	//在这里，我们可以看到一个方法函数可以返回两种类型中的任何一种（map 或者 int）,
	// 我们不用指定特定的返回类型。
	return iris.StatusBadRequest
}
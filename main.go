package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"overview/datasource"
	"overview/repositories"
	"overview/services"
	"overview/web/controllors"
	"overview/web/middleware"
)

func main()  {
	app :=iris.New()
	app.Logger().SetLevel("debug")
	//app.Use(recover.New())
	//app.Use(logger.New())

	// 加载视图模板地址
	app.RegisterView(iris.HTML("./web/views",".html"))
	// 注册控制器
	// mvc.New(app.Party("/movies")).Handle(new(controllers.MovieController))
	//你也可以使用  `mvc.Configure` 方法拆分编写 MVC 应用程序的配置。
	// 如下所示：
	mvc.Configure(app.Party("/movies"),movies)

	app.Run(
			// Start the web server at localhost:8080
			iris.Addr(":8080"),
			// skip err server closed when CTRL/CMD+C pressed:
			iris.WithoutServerError(iris.ErrServerClosed),
			// enables faster json serialization and more:
			iris.WithOptimizations,
		)
}

// 注册控制器
// mvc.New(app.Party("/movies")).Handle(new(controllers.MovieController))
//你也可以使用  `mvc.Configure` 方法拆分编写 MVC 应用程序的配置。
// 如下所示：
func movies(app *mvc.Application){
	// Add the basic authentication(admin:password) middleware
	// for the /movies based requests.
	app.Router.Use(middleware.BasicAuth)

	// 使用数据源中的一些（内存）数据创建 movie 的数据库。
	repo := repositories.NewMovieRepository(datasource.Movies)
	// 创建 movie 的服务，我们将它绑定到 movie 应用程序。
	movieService := services.NewMovieService(repo)
	app.Register(movieService)
	//初始化控制器
	// 注意，你可以初始化多个控制器
	// 你也可以 使用 `movies.Party(relativePath)` 或者 `movies.Clone(app.Party(...))` 创建子应用。
	app.Handle(new(controllors.MovieControllor))
}
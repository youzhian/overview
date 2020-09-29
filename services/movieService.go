package services

import (
	"overview/datamodels"
	"overview/repositories"
)

// `MovieService` 会处理一些 `movie` 数据模型层的 CRUID 操作
// 这取决于 `movie` 存储库 的一些行为.
//这里将数据源和高级组件进行解耦
// 所以，我们可以在不做任何修改的情况下，轻松的切换使用不同的储库类型
// 这个是一个通用的接口
//因为我们可能需要在不的地方修改和尝试不同的逻辑
type MovieService interface {
	GetAll() []datamodels.Movie
	GetById(id int64) (datamodels.Movie, bool)
	DeleteById(id int64) bool
	UpdatePosterAndGenreById(id int64, poster string, genre string) (datamodels.Movie,error)
}
type movieService struct {
	repo repositories.MovieRepository
}
// NewMovieService 返回默认的movie服务层
func NewMovieService(repo repositories.MovieRepository) MovieService{
	return &movieService{repo: repo}
}

func (s *movieService) GetAll() []datamodels.Movie{
	return s.repo.SelectMany(func(_ datamodels.Movie) bool {
		return true
	},-1)
}

func (s *movieService)GetById(id int64) (datamodels.Movie, bool){
	return s.repo.Select(func(m datamodels.Movie) bool {
		return m.ID == id
	})
}
func (s *movieService) UpdatePosterAndGenreById(id int64,poster string,genre string)(datamodels.Movie,error){
	return s.repo.InsertOrUpdate(datamodels.Movie{ID: id,Poster: poster,Genre: genre})
}

func (s *movieService) DeleteById(id int64) bool{
	return s.repo.Delete(func(m datamodels.Movie) bool {
		return m.ID == id
	},1)
}
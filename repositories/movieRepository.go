package repositories

import (
	"errors"
	"overview/datamodels"
	"sync"
)

//query代表一种访客和它的查询动作
type Query func(movie datamodels.Movie) bool

// MovieRepository 会处理一些关于movie实例的基本操作
// 这是一个以测试为目的的接口，即是一个内存中的movie库
// 或是一个连接到数据库的实例
type MovieRepository interface {

	Exec(query Query,action Query,limit int,mode int)(ok bool)
	Select(query Query)(movie datamodels.Movie,found bool)
	SelectMany(query Query,limit int)(results []datamodels.Movie)
	InsertOrUpdate(movie datamodels.Movie)(updateMovie datamodels.Movie,err error)
	Delete(query Query,limit int)(deleted bool)
}
//movieMemoryRepository就是一个"MovieRepository"
//它负责存储内存中的数据（map）
type movieMemoryRepository struct {
	source map[int64]datamodels.Movie
	mu sync.RWMutex
}

// NewMovieRepository 返回一个新的基于内存的movie库
func NewMovieRepository(source map[int64]datamodels.Movie) MovieRepository{
	return &movieMemoryRepository{source: source}
}

const(
	// ReadOnlyMode will RLock(read) the data.
	ReadOnlyMode = iota
	// ReadWriteMode will Lock(write/read) the data.
	ReadWriteMode
)

func (r *movieMemoryRepository)Exec(query Query,action Query,actionLimit int,mode int)(ok bool){
	loops :=0

	if mode == ReadOnlyMode{
		r.mu.RLock()
		//运行完成后执行
		defer r.mu.RUnlock()
	}else {
		r.mu.Lock()
		//运行完成后执行
		defer r.mu.Unlock()
	}

	for _,movie:= range r.source{
		ok = query(movie)
		if ok{
			if action(movie){
				loops++
				if actionLimit >= loops{
					break
				}
			}
		}
	}

	return
}

// Select方法会收到一个查询方法
// 这个方法给出一个单独的movie实例
// 直到这个功能返回为true时停止迭代。
//
// 它返回最后一次查询成功所找到的结果的值
// 和最后的movie模型
// 以减少caller之间的通信
//
// 这是一个很简单但很聪明的雏形方法
// 我基本在所有会用到的地方使用自从我想到了它
// 也希望你们觉得好用
func (r *movieMemoryRepository)Select(query Query)(movie datamodels.Movie,found bool)  {
	found = r.Exec(query, func(m datamodels.Movie) bool {
		movie = m
		return true
	},1,ReadOnlyMode)
	if !found{
		movie = datamodels.Movie{}
	}
	return
}
// SelectMany作用相同于Select但是它返回一个切片
// 切片包含一个或多个实例
// 如果传入的参数limit<=0则返回所有
func (r *movieMemoryRepository)SelectMany(query Query,limit int)(results []datamodels.Movie){
	r.Exec(query, func(m datamodels.Movie) bool {
		results = append(results,m)
		return true
	},limit,ReadOnlyMode)
	return
}

func (r *movieMemoryRepository)InsertOrUpdate(movie datamodels.Movie)(updateMovie datamodels.Movie,err error){
	id := movie.ID

	if id == 0{//创建操作
		var lastId int64
		// 找到最大的ID，避免重复。
		// 在实际使用时您可以使用第三方库去生成
		// 一个string类型的UUID
		r.mu.RLock()
		for _,item :=range r.source{
			if item.ID > lastId{
				lastId = item.ID
			}
		}
		r.mu.RUnlock()
		id = lastId + 1
		movie.ID = id

		r.mu.Lock()
		r.source[id] = movie
		r.mu.Unlock()

		return movie,nil
	}else{
		// 更新操作是基于movie.ID的，
		// 在例子中我们允许了对poster和genre的更新（如果它们非空）。
		// 当然我们可以只是做单纯的数据替换操作:
		// r.source[id] = movie
		// 并注释掉下面的代码;
		current, exists := r.Select(func(m datamodels.Movie) bool {
			return m.ID == id
		})
		//数据不存在
		if !exists {
			//当根据ID查询不到数据时，抛出一个异常
			return datamodels.Movie{},errors.New("failed to update nonexistent movie")
		}
		if movie.Poster != ""{
			current.Poster = movie.Poster
		}
		if movie.Genre != ""{
			current.Genre = movie.Genre
		}
		r.mu.Lock()
		r.source[id] = current
		r.mu.Unlock()

		return movie,nil
	}
}

func (r *movieMemoryRepository)Delete(query Query,limit int)(deleted bool){
	return r.Exec(query, func(m datamodels.Movie) bool {
		delete(r.source,m.ID)
		return true
	},limit,ReadWriteMode)
}
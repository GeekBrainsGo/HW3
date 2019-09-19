package usecases

import (
	"fmt"

	"github.com/art-frela/HW3/karpov/domain"
)

// PostController - main controller for Posts
type PostController struct {
	//Posts []domain.PostInBlog
	PostRepo domain.PostRepository
	Log      Logger
}

// [Case - look posts]

// GetPosts get Posts from PostRepository and fill Posts field
func (pc *PostController) GetPosts(page, pagesize int) ([]domain.PostInBlog, error) {
	posts, err := pc.PostRepo.Find(page, pagesize)
	if err != nil {
		err = fmt.Errorf("get posts Find(limit=%d, offset=%d) error, %v", page, pagesize, err)
		return posts, err
	}
	//pc.Posts = posts
	return posts, nil
}

// GetSinglePost get Posts from PostRepository and fill Posts field
func (pc *PostController) GetSinglePost(id string) (domain.PostInBlog, error) {
	post, err := pc.PostRepo.FindByID(id)
	if err != nil {
		err = fmt.Errorf("get post by ID=%s, error, %v", id, err)
		return post, err
	}
	return post, nil
}

// SaveNewPost saves new post in the repository and returns new Post ID
func (pc *PostController) SaveNewPost(newpost domain.PostInBlog) (string, error) {
	postID, err := pc.PostRepo.Save(newpost)
	if err != nil {
		err = fmt.Errorf("save new post, error, %v", err)
		return postID, err
	}
	return postID, nil
}

// UpdPost updates exists post in the storage
func (pc *PostController) UpdPost(post domain.PostInBlog) error {
	err := pc.PostRepo.Update(post)
	if err != nil {
		err = fmt.Errorf("update post, error, %v", err)
		return err
	}
	return nil
}

//
//func (pc *PostController) prepareNewPostForm() error

// Logger - behavior of logging system
type Logger interface {
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Panic(args ...interface{})
}

// Paginator - using of pagination
type Paginator struct {
	Count       int
	Page        int
	PageCount   int
	PageSize    int
	HasNextPage bool
}

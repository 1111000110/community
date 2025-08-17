package model

import (
	"community.com/service/post/rpc/postservice"
)

func RpcPostToModelPost(post *postservice.Post) *Post {
	return &Post{
		PostId:   post.PostId,
		UserId:   post.UserId,
		Title:    post.Title,
		Content:  post.Content,
		Images:   post.Images,
		Theme:    post.Theme,
		Tags:     post.Tags,
		Status:   post.Status,
		CreateAt: post.CreateTime,
		UpdateAt: post.UpdateTime,
	}
}

func ModelPostToRpcPost(post *Post) *postservice.Post {
	return &postservice.Post{
		PostId:     post.PostId,
		UserId:     post.UserId,
		Title:      post.Title,
		Content:    post.Content,
		Images:     post.Images,
		Theme:      post.Theme,
		Tags:       post.Tags,
		Status:     post.Status,
		CreateTime: post.CreateAt,
		UpdateTime: post.UpdateAt,
	}
}
func ModelPostsToRpcPosts(posts []*Post) []*postservice.Post {
	result := make([]*postservice.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, ModelPostToRpcPost(post))
	}
	return result
}

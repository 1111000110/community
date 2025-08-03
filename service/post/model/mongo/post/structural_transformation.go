package model

import "community.com/service/post/rpc/postservice"

func PostToModelPost(post *postservice.Post) *Post {
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

func ModelPostToPost(post *Post) *postservice.Post {
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

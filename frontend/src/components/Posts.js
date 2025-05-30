import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { fetchPosts } from '../actions/postsActions';
import { fetchCategories } from '../actions/categoriesActions';

const Posts = () => {
  const dispatch = useDispatch();
  const posts = useSelector((state) => state.posts.posts);
  const categories = useSelector((state) => state.categories.categories);
  const isAuthenticated = useSelector((state) => state.auth.isAuthenticated);
  const [categoryName, setCategoryName] = useState('');
  const [categoryId, setCategoryId] = useState('');

  useEffect(() => {
    dispatch(fetchPosts());
    dispatch(fetchCategories());
  }, [dispatch]);

  useEffect(() => {
    if (categories.length > 0) {
      setCategoryName(categories[0].name);
      setCategoryId(categories[0].id);
    }
  }, [categories]);

  return (
    <div className="posts-container">
      <h2>Посты в категории: {categoryName}</h2>
      {posts.length === 0 ? (
        <p>В этой категории пока нет постов</p>
      ) : (
        <div className="posts-list">
          {posts.map((post) => (
            <div key={post.id} className="post-card">
              <h3>{post.title}</h3>
              <p>{post.content}</p>
              <div className="post-footer">
                <span>Автор: {post.author_username}</span>
                <span>Дата: {new Date(post.created_at).toLocaleString()}</span>
                <Link to={`/post/${post.id}`}>Читать далее</Link>
              </div>
            </div>
          ))}
        </div>
      )}
      {isAuthenticated && (
        <Link to={`/create-post/${categoryId}`} className="create-post-button">
          Создать новый пост
        </Link>
      )}
    </div>
  );
};

export default Posts; 
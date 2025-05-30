  <div className="post-container">
    <div className="post-header">
      <h2>{post.title}</h2>
      <div className="post-meta">
        <span>Автор: {post.author_username}</span>
        <span>Дата: {new Date(post.created_at).toLocaleString()}</span>
      </div>
    </div>
    <div className="post-content">
      <p>{post.content}</p>
    </div>
    <div className="post-actions">
      {post.author_id === currentUserId && (
        <>
          <button onClick={handleEdit}>Редактировать</button>
          <button onClick={handleDelete}>Удалить</button>
        </>
      )}
    </div>
    <Comments postId={post.id} />
  </div> 
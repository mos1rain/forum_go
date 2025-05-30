// ... existing code ...
      <div className="comments-section">
        <h3>Комментарии</h3>
        <form onSubmit={handleSubmit} className="comment-form">
          <textarea
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            placeholder="Напишите ваш комментарий..."
            required
          />
          <button type="submit">Оставить комментарий</button>
        </form>
        <div className="comments-list">
          {comments.map((comment) => (
            <div key={comment.id} className="comment">
              <div className="comment-header">
                <span className="username">{comment.username}</span>
                <span className="date">{new Date(comment.created_at).toLocaleString()}</span>
              </div>
              <p className="content">{comment.content}</p>
            </div>
          ))}
        </div>
      </div>
// ... existing code ...
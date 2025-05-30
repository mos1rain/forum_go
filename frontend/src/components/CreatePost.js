  <form onSubmit={handleSubmit}>
    <h2>Создать новый пост</h2>
    <div className="form-group">
      <label htmlFor="title">Заголовок:</label>
      <input
        type="text"
        id="title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
        placeholder="Введите заголовок поста"
      />
    </div>
    <div className="form-group">
      <label htmlFor="content">Содержание:</label>
      <textarea
        id="content"
        value={content}
        onChange={(e) => setContent(e.target.value)}
        required
        placeholder="Напишите содержание вашего поста"
      />
    </div>
    <button type="submit">Опубликовать пост</button>
  </form> 
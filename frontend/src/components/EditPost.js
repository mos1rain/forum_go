  <form onSubmit={handleSubmit} className="edit-post-form">
    <h2>Редактировать пост</h2>
    <div className="form-group">
      <label htmlFor="title">Заголовок:</label>
      <input
        type="text"
        id="title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
        placeholder="Введите новый заголовок"
      />
    </div>
    <div className="form-group">
      <label htmlFor="content">Содержание:</label>
      <textarea
        id="content"
        value={content}
        onChange={(e) => setContent(e.target.value)}
        required
        placeholder="Введите новое содержание"
      />
    </div>
    <div className="form-actions">
      <button type="submit">Сохранить изменения</button>
      <button type="button" onClick={handleCancel}>Отмена</button>
    </div>
  </form> 
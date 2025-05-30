  <div className="profile-container">
    <h2>Профиль пользователя</h2>
    <div className="profile-info">
      <div className="info-group">
        <label>Имя пользователя:</label>
        <span>{user.username}</span>
      </div>
      <div className="info-group">
        <label>Email:</label>
        <span>{user.email}</span>
      </div>
      <div className="info-group">
        <label>Дата регистрации:</label>
        <span>{new Date(user.created_at).toLocaleString()}</span>
      </div>
    </div>
    <div className="profile-stats">
      <h3>Статистика</h3>
      <div className="stats-grid">
        <div className="stat-item">
          <span className="stat-label">Созданные категории:</span>
          <span className="stat-value">{user.categories_count}</span>
        </div>
        <div className="stat-item">
          <span className="stat-label">Написанные посты:</span>
          <span className="stat-value">{user.posts_count}</span>
        </div>
        <div className="stat-item">
          <span className="stat-label">Комментарии:</span>
          <span className="stat-value">{user.comments_count}</span>
        </div>
      </div>
    </div>
  </div> 
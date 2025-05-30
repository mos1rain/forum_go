  <div className="categories-container">
    <h2>Категории форума</h2>
    {categories.length === 0 ? (
      <p>Пока нет созданных категорий</p>
    ) : (
      <div className="categories-list">
        {categories.map((category) => (
          <div key={category.id} className="category-card">
            <h3>{category.name}</h3>
            <p>{category.description}</p>
            <div className="category-footer">
              <span>Создатель: {category.creator_username}</span>
              {category.creator_id === currentUserId && (
                <button onClick={() => handleDeleteCategory(category.id)}>
                  Удалить категорию
                </button>
              )}
            </div>
          </div>
        ))}
      </div>
    )}
  </div> 
  <form onSubmit={handleSubmit}>
    <h2>Создать новую категорию</h2>
    <div className="form-group">
      <label htmlFor="name">Название категории:</label>
      <input
        type="text"
        id="name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
        placeholder="Введите название категории"
      />
    </div>
    <div className="form-group">
      <label htmlFor="description">Описание:</label>
      <textarea
        id="description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        required
        placeholder="Опишите, о чем будет эта категория"
      />
    </div>
    <button type="submit">Создать категорию</button>
  </form> 
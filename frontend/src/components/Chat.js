  <div className="chat-container">
    <div className="chat-messages">
      {messages.map((msg, index) => (
        <div key={index} className={`message ${msg.username === currentUser ? 'own-message' : ''}`}>
          <span className="username">{msg.username}:</span>
          <p>{msg.content}</p>
        </div>
      ))}
    </div>
    <form onSubmit={handleSubmit} className="chat-input">
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Введите ваше сообщение..."
      />
      <button type="submit">Отправить</button>
    </form>
  </div> 
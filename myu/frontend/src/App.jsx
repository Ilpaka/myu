// MessengerApp.jsx
import React, { useState, useEffect } from 'react';

const API_BASE = 'http://138.124.14.1:8080'; // Замени на IP твоего VPS

export default function MessengerApp() {
  const [users, setUsers] = useState([]);
  const [messages, setMessages] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);
  const [newUserName, setNewUserName] = useState('');
  const [messageText, setMessageText] = useState('');
  const [currentUserId, setCurrentUserId] = useState(0);

  const loadUsers = async () => {
    try {
      const response = await fetch(`${API_BASE}/users`);
      const data = await response.json();
      setUsers(data || []);
      if (data?.length && !currentUserId) {
        setCurrentUserId(data[0].id);
      }
    } catch (error) {
      console.error('Error loading users:', error);
    }
  };

  const loadMessages = async () => {
    if (!currentUserId) return;
    try {
      const response = await fetch(`${API_BASE}/messages/user/${currentUserId}`);
      const data = await response.json();
      setMessages(data || []);
    } catch (error) {
      console.error('Error loading messages:', error);
    }
  };

  const createUser = async () => {
    if (!newUserName.trim()) return;
    try {
      const response = await fetch(`${API_BASE}/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newUserName })
      });
      if (response.ok) {
        setNewUserName('');
        await loadUsers();
      }
    } catch (error) {
      console.error('Error creating user:', error);
    }
  };

  const sendMessage = async () => {
    if (!selectedUser || !messageText.trim() || !currentUserId) return;
    try {
      const response = await fetch(`${API_BASE}/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          text: messageText,
          from_id: currentUserId,
          to_id: selectedUser.id
        })
      });
      if (response.ok) {
        setMessageText('');
        await loadMessages();
      }
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  const markAsRead = async (id) => {
    try {
      await fetch(`${API_BASE}/messages/${id}/read`, { method: 'PATCH' });
      await loadMessages();
    } catch (error) {
      console.error('Error marking as read:', error);
    }
  };

  useEffect(() => {
    loadUsers();
  }, []);

  useEffect(() => {
    loadMessages();
    const interval = setInterval(loadMessages, 5000);
    return () => clearInterval(interval);
  }, [currentUserId]);

  return (
    <div style={{ padding: 16, maxWidth: 900, margin: '0 auto', fontFamily: 'Arial, sans-serif' }}>
      <h2>🚀 Remote Messenger</h2>
      <p style={{ color: '#666', fontSize: '14px', marginBottom: 20 }}>
        API: {API_BASE}
      </p>

      {/* Создание пользователя */}
      <section style={{ 
        marginBottom: 20, 
        padding: 15, 
        border: '1px solid #ddd', 
        borderRadius: 8, 
        backgroundColor: '#f9f9f9' 
      }}>
        <h4 style={{ margin: '0 0 10px 0' }}>👤 Создать пользователя</h4>
        <input
          type="text"
          value={newUserName}
          onChange={(e) => setNewUserName(e.target.value)}
          placeholder="Имя пользователя"
          style={{ 
            marginRight: 10, 
            padding: 8, 
            border: '1px solid #ccc', 
            borderRadius: 4,
            width: 200
          }}
          onKeyPress={(e) => e.key === 'Enter' && createUser()}
        />
        <button 
          onClick={createUser} 
          style={{ 
            padding: '8px 16px', 
            backgroundColor: '#007bff', 
            color: 'white', 
            border: 'none', 
            borderRadius: 4,
            cursor: 'pointer'
          }}
        >
          Создать
        </button>
      </section>

      {/* Выбор текущего пользователя */}
      <section style={{ marginBottom: 20 }}>
        <h4>🎭 Войти как:</h4>
        <select
          value={currentUserId || ''}
          onChange={(e) => setCurrentUserId(parseInt(e.target.value))}
          style={{ 
            padding: 8, 
            border: '1px solid #ccc', 
            borderRadius: 4,
            width: 250,
            fontSize: 14
          }}
        >
          <option value="" disabled>— выберите пользователя —</option>
          {users.map(u => (
            <option key={u.id} value={u.id}>{u.name} (ID: {u.id})</option>
          ))}
        </select>
      </section>

      {/* Отправка сообщения */}
      <section style={{ 
        marginBottom: 20, 
        padding: 15, 
        border: '1px solid #ddd', 
        borderRadius: 8,
        backgroundColor: '#f0f8ff'
      }}>
        <h4 style={{ margin: '0 0 10px 0' }}>✉️ Отправить сообщение</h4>
        
        <div style={{ marginBottom: 10 }}>
          <label style={{ display: 'block', marginBottom: 5, fontWeight: 'bold' }}>
            Получатель:
          </label>
          <select
            value={selectedUser?.id || ''}
            onChange={(e) => setSelectedUser(users.find(u => u.id === parseInt(e.target.value)))}
            style={{ 
              padding: 8, 
              border: '1px solid #ccc', 
              borderRadius: 4,
              width: '100%',
              maxWidth: 300
            }}
          >
            <option value="">— выберите получателя —</option>
            {users.filter(u => u.id !== currentUserId).map(u => (
              <option key={u.id} value={u.id}>{u.name} (ID: {u.id})</option>
            ))}
          </select>
        </div>

        <div style={{ marginBottom: 10 }}>
          <label style={{ display: 'block', marginBottom: 5, fontWeight: 'bold' }}>
            Сообщение:
          </label>
          <textarea
            rows={3}
            value={messageText}
            onChange={(e) => setMessageText(e.target.value)}
            placeholder="Введите сообщение..."
            style={{ 
              width: '100%', 
              padding: 8, 
              border: '1px solid #ccc', 
              borderRadius: 4,
              resize: 'vertical',
              fontFamily: 'Arial, sans-serif'
            }}
            onKeyPress={(e) => {
              if (e.key === 'Enter' && e.ctrlKey) {
                sendMessage();
              }
            }}
          />
          <small style={{ color: '#666' }}>Ctrl+Enter для отправки</small>
        </div>

        <button 
          onClick={sendMessage} 
          disabled={!selectedUser || !messageText.trim() || !currentUserId}
          style={{ 
            padding: '10px 20px', 
            backgroundColor: selectedUser && messageText.trim() && currentUserId ? '#28a745' : '#ccc',
            color: 'white', 
            border: 'none', 
            borderRadius: 4,
            cursor: selectedUser && messageText.trim() && currentUserId ? 'pointer' : 'not-allowed',
            fontSize: 16
          }}
        >
          Отправить 📤
        </button>
      </section>

      {/* Входящие сообщения */}
      <section>
        <h4>📬 Входящие сообщения {currentUserId ? `для ${users.find(u => u.id === currentUserId)?.name || 'пользователя'}` : ''}</h4>
        
        {!currentUserId ? (
          <p style={{ color: '#666', fontStyle: 'italic' }}>
            Выберите пользователя для просмотра сообщений
          </p>
        ) : messages.length === 0 ? (
          <div style={{ 
            padding: 20, 
            textAlign: 'center', 
            border: '1px dashed #ccc', 
            borderRadius: 8,
            color: '#666'
          }}>
            📭 Нет сообщений
          </div>
        ) : (
          <div>
            {messages.map(msg => (
              <div
                key={msg.id}
                style={{
                  border: '1px solid #ccc',
                  padding: 15,
                  marginBottom: 10,
                  borderRadius: 8,
                  backgroundColor: msg.is_read ? '#f8f9fa' : '#e3f2fd',
                  borderLeft: msg.is_read ? '4px solid #28a745' : '4px solid #007bff'
                }}
              >
                <div style={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center',
                  marginBottom: 8
                }}>
                  <div style={{ fontWeight: 'bold', color: '#333' }}>
                    От: {msg.from_name}
                    {!msg.is_read && <span style={{ color: 'red', marginLeft: 8 }}>🆕 НОВОЕ</span>}
                  </div>
                  <div style={{ fontSize: '12px', color: '#666' }}>
                    {msg.timestamp}
                  </div>
                </div>
                
                <div style={{ 
                  marginBottom: 10, 
                  fontSize: 16, 
                  lineHeight: 1.4,
                  padding: 8,
                  backgroundColor: 'white',
                  borderRadius: 4,
                  border: '1px solid #eee'
                }}>
                  {msg.text}
                </div>
                
                {!msg.is_read && (
                  <button
                    onClick={() => markAsRead(msg.id)}
                    style={{
                      padding: '6px 12px',
                      fontSize: '12px',
                      backgroundColor: '#17a2b8',
                      color: 'white',
                      border: 'none',
                      borderRadius: 4,
                      cursor: 'pointer'
                    }}
                  >
                    Прочитать ✅
                  </button>
                )}
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Статистика */}
      <div style={{ 
        marginTop: 30, 
        padding: 15, 
        backgroundColor: '#f8f9fa', 
        borderRadius: 8,
        textAlign: 'center'
      }}>
        <small style={{ color: '#666' }}>
          👥 Пользователей: {users.length} | 
          📨 Сообщений: {messages.length} | 
          🔄 Обновление каждые 5 секунд
        </small>
      </div>
    </div>
  );
}

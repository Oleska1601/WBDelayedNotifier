 const API_BASE_URL = 'http://localhost:8081/api/v1'; // Измените на ваш URL API
        
        // Создание уведомления
        document.getElementById('createForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = {
                channel: document.getElementById('channel').value,
                recipient: document.getElementById('recipient').value,
                message: document.getElementById('message').value,
                scheduled_at: new Date(document.getElementById('scheduledAt').value).toISOString()
            };
            
            try {
                const response = await fetch(`${API_BASE_URL}/notify`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });
                
                const result = await response.json();
                const resultDiv = document.getElementById('createResult');
                
                if (response.ok) {
                    resultDiv.innerHTML = `<span class="success">Успешно создано! ID: ${result.id}</span>`;
                } else {
                    resultDiv.innerHTML = `<span class="error">Ошибка: ${JSON.stringify(result, null, 2)}</span>`;
                }
            } catch (error) {
                document.getElementById('createResult').innerHTML = `<span class="error">Ошибка: ${error.message}</span>`;
            }
        });
        
        // Получение статуса уведомления
        document.getElementById('getForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const notificationId = document.getElementById('notificationId').value;
            
            try {
                const response = await fetch(`${API_BASE_URL}/notify/${notificationId}`);
                const result = await response.json();
                const resultDiv = document.getElementById('getResult');
                
                if (response.ok) {
                    resultDiv.innerHTML = `<span class="success">Статус: ${JSON.stringify(result, null, 2)}</span>`;
                } else {
                    resultDiv.innerHTML = `<span class="error">Ошибка: ${JSON.stringify(result, null, 2)}</span>`;
                }
            } catch (error) {
                document.getElementById('getResult').innerHTML = `<span class="error">Ошибка: ${error.message}</span>`;
            }
        });
        
        // Удаление уведомления
        document.getElementById('deleteForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const notificationId = document.getElementById('deleteNotificationId').value;
            
            try {
                const response = await fetch(`${API_BASE_URL}/notify/${notificationId}`, {
                    method: 'DELETE'
                });
                
                const resultDiv = document.getElementById('deleteResult');
                
                if (response.ok) {
                    resultDiv.innerHTML = `<span class="success">Уведомление ${notificationId} успешно удалено</span>`;
                } else {
                    const result = await response.json();
                    resultDiv.innerHTML = `<span class="error">Ошибка: ${JSON.stringify(result, null, 2)}</span>`;
                }
            } catch (error) {
                document.getElementById('deleteResult').innerHTML = `<span class="error">Ошибка: ${error.message}</span>`;
            }
        });
// Базовый класс для работы с API
class NotificationAPI {
    static baseURL = '/api';

    // Создание уведомления
  
    static async createNotification(notificationData) {
        try {
            const response = await fetch(`${this.baseURL}/notify`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(notificationData)
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Ошибка при создании уведомления:', error);
            throw error;
        }
    }

    // Получение уведомления по ID
    static async getNotification(id) {
        try {
            const response = await fetch(`${this.baseURL}/notify/${id}`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Ошибка при получении уведомления:', error);
            throw error;
        }
    }

    // Отмена уведомления
    static async cancelNotification(id) {
        try {
            const response = await fetch(`${this.baseURL}/notify/${id}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Ошибка при отмене уведомления:', error);
            throw error;
        }
    }
/*
    // Получение всех уведомлений
    static async getAllNotifications() {
        try {
            // В реальном API может быть endpoint для получения списка
            // Пока эмулируем получение через несколько запросов или бэкенд добавит /notifications
            const response = await fetch(`${this.baseURL}/notifications`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Ошибка при получении списка уведомлений:', error);
            throw error;
        }
    }
        */
}
// Базовый класс для работы с API
class NotificationAPI {
    static baseURL = '';

    // Создание уведомления
    static async createNotification(notificationData) {
        const response = await fetch(`${this.baseURL}/notify`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(notificationData)
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }

        return await response.json();
    }

    // Получение уведомления по ID
    static async getNotification(id) {
        const response = await fetch(`${this.baseURL}/notify/${id}`);
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }

        return await response.json();
    }

    // Отмена уведомления
    static async cancelNotification(id) {
        const response = await fetch(`${this.baseURL}/notify/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }

        return await response.json();
    }

    // Получение всех уведомлений
    static async getAllNotifications() {
        const response = await fetch(`${this.baseURL}/notifications`);
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }

        return await response.json();
    }
}
// Класс для управления уведомлениями
class NotificationManager {
    static notifications = [];
    static currentFilter = 'all';

    // Загрузить все уведомления
    static async loadAll() {
        try {
            this.notifications = await NotificationAPI.getAllNotifications();
            this.renderList();
            UI.updateStats(this.notifications);
        } catch (error) {
            UI.showError(error.message || 'Не удалось загрузить уведомления');
        }
    }

    // Создать новое уведомление
    static async createNotification(formData) {
        try {
            const newNotification = await NotificationAPI.createNotification(formData);
            this.notifications.unshift(newNotification); // Добавляем в начало
            this.renderList();
            UI.updateStats(this.notifications);
            UI.showSuccess('Уведомление успешно создано!');
        } catch (error) {
            UI.showError(error.message || 'Ошибка при создании уведомления');
        }
    }

    // Отменить уведомление
    static async cancelNotification(id) {
        try {
            await NotificationAPI.cancelNotification(id);
            
            // Обновляем статус локально
            const notification = this.notifications.find(n => n.id === id);
            if (notification) {
                notification.status = 'cancelled';
            }
            
            this.renderList();
            UI.updateStats(this.notifications);
            UI.showSuccess('Уведомление отменено');
        } catch (error) {
            UI.showError(error.message || 'Ошибка при отмене уведомления');
        }
    }

    // Отрендерить список уведомлений
    static renderList() {
        const tbody = document.getElementById('notifications-tbody');
        const filteredNotifications = this.getFilteredNotifications();

        if (filteredNotifications.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" style="text-align: center;">Нет уведомлений</td></tr>';
            return;
        }

        tbody.innerHTML = filteredNotifications.map(notification => `
            <tr>
                <td>${notification.id}</td>
                <td>${this.truncateText(notification.message, 50)}</td>
                <td>${this.formatDateTime(notification.scheduled_at)}</td>
                <td>${this.getChannelIcon(notification.channel)} ${notification.channel}</td>
                <td class="status-${notification.status}">${this.getStatusText(notification.status)}</td>
                <td>
                    ${notification.status === 'scheduled' ? 
                        `<button class="btn-danger" onclick="NotificationManager.cancelNotification('${notification.id}')">
                            Отменить
                         </button>` : 
                        '-'
                    }
                </td>
            </tr>
        `).join('');
    }

    // Применить фильтр
    static applyFilter(filter) {
        this.currentFilter = filter;
        this.renderList();
    }

    // Получить отфильтрованные уведомления
    static getFilteredNotifications() {
        if (this.currentFilter === 'all') {
            return this.notifications;
        }
        return this.notifications.filter(n => n.status === this.currentFilter);
    }

    // Вспомогательные методы
    static truncateText(text, maxLength) {
        return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
    }

    static formatDateTime(dateTimeString) {
        const date = new Date(dateTimeString);
        return date.toLocaleString('ru-RU');
    }

    static getChannelIcon(channel) {
        const icons = {
            email: '📧',
            telegram: '📱'
        };
        return icons[channel] || '📨';
    }

    static getStatusText(status) {
        const statuses = {
            scheduled: '⏳ Ожидает',
            sent: '✅ Отправлено',
            failed: '❌ Ошибка',
            cancelled: '🚫 Отменено'
        };
        return statuses[status] || status;
    }
}
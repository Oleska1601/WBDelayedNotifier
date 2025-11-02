// Класс для управления пользовательским интерфейсом
class UI {
    // Показать модальное окно
    static showModal(title, message) {
        const modal = document.getElementById('modal');
        const modalTitle = document.getElementById('modal-title');
        const modalMessage = document.getElementById('modal-message');

        modalTitle.textContent = title;
        modalMessage.textContent = message;
        modal.classList.remove('hidden');
    }

    // Скрыть модальное окно
    static hideModal() {
        const modal = document.getElementById('modal');
        modal.classList.add('hidden');
    }

    // Показать уведомление об успехе
    static showSuccess(message) {
        this.showModal('Успех', message);
    }

    // Показать ошибку
    static showError(message) {
        this.showModal('Ошибка', message);
    }

    // Обновить статистику
    // Обновить статистику
    static updateStats(notifications) {
        const totalCount = document.getElementById('total-count');
        const scheduledCount = document.getElementById('scheduled-count');
        const sentCount = document.getElementById('sent-count');
        const failedCount = document.getElementById('failed-count');
        const cancelledCount = document.getElementById('cancelled-count');

        const total = notifications.length;
        const scheduled = notifications.filter(n => n.status === 'scheduled').length;
        const sent = notifications.filter(n => n.status === 'sent').length;
        const failed = notifications.filter(n => n.status === 'failed').length;
        const cancelled = notifications.filter(n => n.status === 'cancelled').length;

        totalCount.textContent = total;
        scheduledCount.textContent = scheduled;
        sentCount.textContent = sent;
        failedCount.textContent = failed;
        cancelledCount.textContent = cancelled;
    }

    // Установить обработчики событий
    static setupEventListeners() {
        // Закрытие модального окна
        document.querySelector('.close').addEventListener('click', this.hideModal);
        document.getElementById('modal').addEventListener('click', (e) => {
            if (e.target.id === 'modal') {
                this.hideModal();
            }
        });

        // Обработка фильтров
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                // Убрать активный класс у всех кнопок
                document.querySelectorAll('.filter-btn').forEach(b => {
                    b.classList.remove('active');
                });
                // Добавить активный класс текущей кнопке
                e.target.classList.add('active');
                
                // Применить фильтр
                const filter = e.target.dataset.filter;
                NotificationManager.applyFilter(filter);
            });
        });
    }

    // Обновление полей получателя при смене канала
    static updateRecipientFields(channel) {
        const container = document.getElementById('recipient-fields');
        
        const fields = {
            email: `
                <div class="form-group">
                    <label for="recipient">Получатель:</label>
                    <input type="text" id="recipient" placeholder="Введите email">
                </div>
            `,
            telegram: `
                <div class="form-group">
                    <label for="recipient">Получатель:</label>
                    <input type="text" id="recipient" placeholder="Введите Telegram ID">
                </div>
            `
        };

        container.innerHTML = fields[channel] || fields.email;
    }
}
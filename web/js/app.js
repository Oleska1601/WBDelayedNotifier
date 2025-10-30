// Главный класс приложения
class App {
    constructor() {
        this.init();
    }

    async init() {
        // Инициализация UI
        UI.setupEventListeners();
        
        // Настройка обработчиков формы
        this.setupFormHandlers();
        
        // Загрузка начальных данных
        await NotificationManager.loadAll();
        
        // Запуск автообновления
        this.startAutoRefresh();
        
        console.log('Приложение инициализировано');
    }

    // Настройка обработчиков формы
    setupFormHandlers() {
        const form = document.getElementById('create-notification-form');
        
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = {
                message: document.getElementById('message').value,
                scheduled_at: document.getElementById('scheduled-time').value,
                channel: document.getElementById('channel').value
            };

            // Валидация
            if (!this.validateForm(formData)) {
                return;
            }

            // Создание уведомления
            await NotificationManager.createNotification(formData);
            
            // Очистка формы
            form.reset();
        });
    }

    // Валидация формы
    validateForm(formData) {
        if (!formData.message.trim()) {
            UI.showError('Введите текст сообщения');
            return false;
        }

        if (!formData.scheduled_at) {
            UI.showError('Выберите время отправки');
            return false;
        }

        const scheduledTime = new Date(formData.scheduled_at);
        const now = new Date();

        if (scheduledTime <= now) {
            UI.showError('Время отправки должно быть в будущем');
            return false;
        }

        return true;
    }

    // Автообновление списка каждые 30 секунд
    startAutoRefresh() {
        setInterval(async () => {
            await NotificationManager.loadAll();
        }, 30000); // 30 секунд
    }
}

// Инициализация приложения при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    new App();
});
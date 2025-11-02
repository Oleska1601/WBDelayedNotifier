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
                channel: document.getElementById('channel').value,
                recipient: document.getElementById('recipient').value
            };

            // Создание уведомления
            await NotificationManager.createNotification(formData);
            
            // Очистка формы
            form.reset();
        });

        // Обработчик смены канала
        document.getElementById('channel').addEventListener('change', (e) => {
            UI.updateRecipientFields(e.target.value);
        });
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
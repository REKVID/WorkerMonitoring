


# TgTaskBot - Telegram Bot для управления фотоотчетами

## 📝 Описание
TgTaskBot - это Telegram бот, разработанный для автоматизации процесса 
сбора и управления фотоотчетами.
Бот позволяет работникам отправлять фотографии с объектов, 
а администраторам - управлять и просматривать полученные отчеты.

## 🚀 Основные возможности

### 👷 Для работников:
- Выбор объекта работы
- Указание должности
- Отправка фотоотчетов
- Автоматическое сохранение фото с метаданными

### 👨‍💼 Для администраторов:
- Просмотр всех загруженных фотографий
- Фильтрация фото по дате/объекту/должности
- Управление списком объектов и должностей
- Просмотр файловой системы
- Доступ к логам системы

## 🛠 Технические требования

- Windows/Linux/MacOS
- Go 1.16 или выше
- Доступ к Telegram Bot API
- Права на создание файлов и папок

## 📂 Структура проекта

``` bash
TgTaskBot/
├── main.go                 # Точка входа в приложение
├── configuration.json      # Файл конфигурации
├── Config/                 # Пакет конфигурации
│   └── config.go
├── internal/              # Внутренняя логика
│   ├── admin/            # Функционал администратора
│   ├── bot/              # Основная логика бота
│   └── worker/           # Функционал работника
├── pkg/                  # Общие пакеты
│   ├── handlers/         # Обработчики
│   └── logger/           # Логирование
```

## ⚙️ Установка и настройка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/your-username/TgTaskBot.git
cd TgTaskBot
```

2. Создайте файл конфигурации `Config/configuration.json`:
```json
{
    "bot_token": "YOUR_BOT_TOKEN",
    "admin_key": "YOUR_ADMIN_KEY",
    "positions": ["Охранник", "Уборщик"],
    "objects": ["Объект 1", "Объект 2"],
    "base_dir": "./photos",
    "log_path": "./app.log"
}
```

3. Соберите проект:
```bash
# Для Windows
go build -o TgTaskBot.exe main.go

# Для Linux/MacOS
GOOS=windows GOARCH=amd64 go build -o TgTaskBot.exe main.go
```

## 🚀 Запуск

1. Через batch файл:
   - Запустите `start.bat`

2. Напрямую:
   - Запустите исполняемый файл `TgTaskBot.exe`

## 📱 Использование бота

### Начало работы:
1. Отправьте команду `/start`
2. Выберите роль (Работник/Администратор)

### Для работников:
1. Выберите роль "Работник"
2. Выберите объект из списка
3. Укажите свою должность
4. Отправьте фотографию

### Для администраторов:
1. Выберите роль "Администратор"
2. Введите ключ администратора
3. Используйте доступные опции:
   - Просмотр фотографий
   - Просмотр файловой системы
   - Просмотр логов
   - Добавление объектов/должностей

## 📁 Структура хранения фотографий

```
photos/
├── [Должность]/
│   ├── [Объект]/
│   │   └── [Дата]/
│   │       └── [Username].jpg
```

## 🔒 Безопасность
- Доступ к админ-панели защищен паролем
- Все действия логируются
- Фотографии хранятся в организованной структуре

## 📝 Логирование
- Все действия записываются в `app.log`
- Логи включают:
  - Действия пользователей
  - Ошибки системы
  - Загрузку фотографий
  - Административные действия

## ⚠️ Обработка ошибок
- Автоматическое создание необходимых директорий
- Проверка прав доступа
- Валидация входных данных
- Обработка сетевых ошибок

## 🔧 Устранение неполадок

### Распространенные проблемы:

1. **Ошибка конфигурации**
   - Проверьте наличие `configuration.json`
   - Убедитесь в правильности токена бота

2. **Ошибки доступа**
   - Проверьте права на создание файлов/папок
   - Убедитесь, что все необходимые директории существуют

3. **Проблемы с отправкой фото**
   - Проверьте подключение к интернету
   - Убедитесь, что размер фото не превышает лимиты Telegram

## Поддержка
При возникновении проблем:
1. Проверьте лог-файл `app.log`
2. Убедитесь в правильности конфигурации
3. Проверьте права доступа к файловой системе

## 🔄 Обновления
- Регулярно проверяйте репозиторий на наличие обновлений
- Сохраняйте резервную копию `configuration.json` перед обновлением

## 📜 Лицензия
MIT License - свободное использование и модификация

## Автор
https://github.com/REKVID
```

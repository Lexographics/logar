import en from '../en'
import type { BaseTranslation, Translation } from '../i18n-types'

const ru = {
	...en as Translation,
	
	login: {
		title: 'Вход',
		username: 'Имя пользователя',
		password: 'Пароль',
		login: 'Вход',
	},
	topbar: {
		my_account: 'Мой аккаунт',
		sign_out: 'Выйти',
	},
	dashboard: {
		title: 'Панель управления',
	},
	logs: {
		title: 'Журналы',
		live_stream_connected: 'Прямая трансляция подключена',
		live_stream_disconnected: 'Прямая трансляция отключена',
		reconnect: 'Переподключить',
		reconnect_tip: 'Попробовать переподключить прямую трансляцию',
		fields: {
			id: 'ID',
			timestamp: 'Отметка времени',
			category: 'Категория',
			message: 'Сообщение',
			severity: 'Уровень серьезности',
		},
		operators: {
			contains: 'Содержит',
			not_contains: 'Не содержит',
			starts_with: 'Начинается с',
			ends_with: 'Заканчивается на',
			between: 'Между',
			not_between: 'Не между',
			in: 'В списке',
			not_in: 'Не в списке',
		},
		value: 'Значение',
		values_comma_separated: 'Значения (через запятую)',
		now: 'Сейчас',
		from: 'От',
		to: 'До',
		add_filter: 'Добавить фильтр',
		time_format_tip: 'Совет: используйте формат DD-MM-YYYY HH:mm:ss.SSS',
		severity_tip: 'Уровни серьезности: 1 (Трассировка), 2 (Лог), 3 (Инфо), 4 (Предупреждение), 5 (Ошибка), 6 (Критическая ошибка)'
	},
	analytics: {
		title: 'Аналитика',
		header: 'Панель аналитики',
	},
	remote_actions: {
		title: 'Удаленные действия',
		header: 'Удаленные действия',
		loading: 'Загрузка...',
		no_actions: 'Действия не определены.',
		error: 'Ошибка',
		no_args: '(Нет аргументов)',
		invoke: 'Вызвать',
		units: {
			nsec: 'Наносекунда',
			usec: 'Микросекунда',
			msec: 'Миллисекунда',
			sec: 'Секунда',
			min: 'Минута',
			hour: 'Час',
			day: 'День',
			week: 'Неделя',
			month: 'Месяц',
			year: 'Год',
		},
	},
	user_sessions: {
		title: 'Сессии',
		header: 'Сессии',
		my_active_sessions: {
			title: 'Мои активные сессии',
			device: 'Устройство',
			last_activity: 'Последняя активность',
			created_at: 'Создано',
			you: '(Вы)',
			revoke: 'Отменить',
		},
		all_users: {
			title: 'Все пользователи',
			online: 'В сети',
			last_seen: 'Последний визит',
			admin: 'Администратор',
		}
	},
	settings: {
		title: 'Настройки',
		header: 'Настройки',
		theme: {
			title: 'Тема',
			select_theme: 'Выбрать тему',
			options: {
				system: 'Системная',
				light: 'Светлая',
				dark: 'Темная',
				dark_blue: 'Темно-синяя',
				light_green: 'Светло-зеленая',
				dark_green: 'Темно-зеленая',
				light_purple: 'Светло-фиолетовая',
				dark_purple: 'Темно-фиолетовая',
				light_pink: 'Светло-розовая',
				dark_pink: 'Темно-розовая',
			},
			current_theme: 'Текущая тема: {0}',
			theme_description: 'Системная тема будет соответствовать настройкам вашего устройства.',
		},
		language: {
			title: 'Язык и регион',
			language: {
				title: 'Язык',
			}
		},
		profile: {
			title: 'Профиль',
			username: 'Имя пользователя',
			display_name: 'Отображаемое имя',
			save: 'Сохранить',
		},
		account: {
			title: 'Управление аккаунтом',
			change_password: 'Изменить пароль',
			delete_account: 'Удалить аккаунт',
		},
	},
	help: {
		title: 'Помощь',
		header: 'Руководство Logar',
		about_logar: {
			title: 'О Logar',
			content: 'Logar — это легковесная и гибкая библиотека логирования для приложений на Go, предоставляющая комплексное решение для логирования с веб-интерфейсом для мониторинга и анализа.'
		},
		integration_guide: {
			title: 'Руководство по интеграции',
			content: 'Вот справочник по пакету Logar для Go'
		},
		troubleshooting: {
			title: 'Устранение неполадок',
			pre_twitter: 'Если у вас возникли проблемы, свяжитесь со мной в Twitter:',
			pre_github: 'или откройте проблему в',
			github_text: 'репозитории GitHub',
		},
		additional_resources: {
			title: 'Дополнительные ресурсы',
			github_repository: 'Репозиторий GitHub',
			api_docs: 'Документация API',
		}
	},
} satisfies BaseTranslation

export default ru
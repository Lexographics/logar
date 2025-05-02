import en from '../en'
import type { BaseTranslation, Translation } from '../i18n-types'

const kk = {
	...en as Translation,
	
	login: {
		title: 'Кіру',
		username: 'Пайдаланушы аты',
		password: 'Құпия сөз',
		login: 'Кіру',
	},
	topbar: {
		my_account: 'Аккаунтым',
		sign_out: 'Шығу',
	},
	dashboard: {
		title: 'Басқару панелі',
	},
	logs: {
		title: 'Журналдар',
		live_stream_connected: 'Тікелей трансляция қосылды',
		live_stream_disconnected: 'Тікелей трансляция ажыратылды',
		reconnect: 'Қайта қосу',
		reconnect_tip: 'Тікелей трансляцияны қайта қосуға әрекет жасау',
		fields: {
			id: 'ID',
			timestamp: 'Уақыт белгісі',
			category: 'Санат',
			message: 'Хабарлама',
			severity: 'Маңыздылық',
		},
		operators: {
			contains: 'Қамтиды',
			not_contains: 'Қамтымайды',
			starts_with: 'Мынадан басталады',
			ends_with: 'Мынамен аяқталады',
			between: 'Арасында',
			not_between: 'Арасында емес',
			in: 'Ішінде',
			not_in: 'Ішінде емес',
		},
		value: 'Мән',
		values_comma_separated: 'Мәндер (үтірмен бөлінген)',
		now: 'Қазір',
		from: 'Бастап',
		to: 'Дейін',
		add_filter: 'Сүзгі қосу',
		time_format_tip: 'Кеңес: DD-MM-YYYY HH:mm:ss.SSS форматты қолданыңыз',
		severity_tip: 'Маңыздылықтар: 1 (Із), 2 (Журнал), 3 (Ақпарат), 4 (Ескерту), 5 (Қате), 6 (Фатал)'
	},
	analytics: {
		title: 'Аналитика',
		header: 'Аналитика панелі',
	},
	remote_actions: {
		title: 'Қашықтан әрекеттер',
		header: 'Қашықтан әрекеттер',
		loading: 'Жүктелуде...',
		no_actions: 'Әрекеттер анықталмаған.',
		error: 'Қате',
		no_args: '(Аргументтер жоқ)',
		invoke: 'Шақыру',
		units: {
			nsec: "Наносекунд",
			usec: "Микросекунд",
			msec: "Миллисекунд",
			sec: "Секунд",
			min: "Минут",
			hour: "Сағат",
			day: "Күн",
			week: "Апта",
			month: "Ай",
			year: "Жыл",
		},
	},
	user_sessions: {
		title: 'Сессиялар',
		header: 'Сессиялар',
		my_active_sessions: {
			title: 'Белсенді сессияларым',
			device: 'Құрылғы',
			last_activity: 'Соңғы әрекет',
			created_at: 'Құрылған уақыты',
			you: '(Сіз)',
			revoke: 'Тоқтату', // Аяқтау ?
		},
		all_users: {
			title: 'Барлық пайдаланушылар',
			online: 'Желіде',
			last_seen: 'Соңғы көрінген уақыты',
			admin: 'Әкімші',
		}
	},
	settings: {
		title: 'Баптаулар',
		header: 'Баптаулар',
		theme: {
			title: 'Тақырып',
			select_theme: 'Тақырыпты таңдау',
			options: {
				light: 'Жарық',
				dark: 'Қараңғы',
				system: 'Жүйелік таңдау',
			},
			current_theme: 'Ағымдағы тақырып: {0}',
			theme_description: 'Жүйелік баптау құрылғыңыздың баптауларын пайдаланады.',
		},
		language: {
			title: 'Тіл және аймақ',
			language: {
				title: 'Тіл',
			}
		},
		profile: {
			title: 'Профиль',
			username: 'Пайдаланушы аты',
			display_name: 'Көрсетілетін аты',
			save: 'Сақтау',
		},
		account: {
			title: 'Аккаунтты басқару',
			change_password: 'Құпиясөзді өзгерту',
			delete_account: 'Аккаунтты жою',
		},
	},
	help: {
		title: 'Көмек',
		header: 'Logar көмек нұсқаулығы',
		about_logar: {
			title: 'Logar туралы',
			content: 'Logar — Go қосымшалары үшін жеңіл, икемді журналдау кітапханасы, ол веб-интерфейс арқылы мониторинг және талдау мүмкіндіктерін ұсынады.'
		},
		integration_guide: {
			title: 'Интеграция нұсқаулығы',
			content: 'Міне, Logar пакеті үшін Go анықтамасы'
		},
		troubleshooting: {
			title: 'Ақауларды жою',
			pre_twitter: 'Егер қандай да бір мәселелер туындаса, менімен Twitter-де хабарласыңыз:',
			pre_github: 'немесе мына жерде мәселе ашыңыз',
			github_text: 'GitHub репозиторийі',
		},
		additional_resources: {
			title: 'Қосымша ресурстар',
			github_repository: 'GitHub репозиторийі',
			api_docs: 'API құжаттамасы',
		}
	},
} satisfies BaseTranslation

export default kk
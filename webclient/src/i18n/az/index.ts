import en from '../en'
import type { BaseTranslation, Translation } from '../i18n-types'

const az = {
	...en as Translation,

	login: {
		title: 'Giriş',
		username: 'İstifadəçi adı',
		password: 'Şifre',
		login: 'Giriş Yap',
	},
	topbar: {
		my_account: 'Hesabım',
		sign_out: 'Çıxış',
	},
	dashboard: {
		title: 'İdarə paneli',
	},
	logs: {
		title: 'Jurnallar',
		live_stream_connected: 'Canlı yayım qoşuldu',
		live_stream_disconnected: 'Canlı yayım ayrıldı',
		reconnect: 'Yenidən qoşul',
		reconnect_tip: 'Canlı yayımı yenidən qoşmağa cəhd et',
		fields: {
			id: 'ID',
			timestamp: 'Zaman möhürü',
			category: 'Kateqoriya',
			message: 'Mesaj',
			severity: 'Ciddilik',
		},
		operators: {
			contains: 'Əhatə edir',
			not_contains: 'Əhatə etmir',
			starts_with: '... ilə başlayır',
			ends_with: '... ilə bitir',
			between: 'Arasında',
			not_between: 'Arasında deyil',
			in: 'İçində',
			not_in: 'İçində deyil',
		},
		value: 'Dəyər',
		values_comma_separated: 'Dəyərlər (vergüllə ayrılmış)',
		now: 'İndi',
		from: 'Kimdən',
		to: 'Kimə',
		add_filter: 'Filtr əlavə et',
		time_format_tip: 'İpucu: DD-MM-YYYY HH:mm:ss.SSS formatından istifadə edin',
		severity_tip: 'Ciddiliklər: 1 (İz), 2 (Jurnal), 3 (Məlumat), 4 (Xəbərdarlıq), 5 (Xəta), 6 (Kritik)'
	},
	analytics: {
		title: 'Analitika',
		header: 'Analitika Paneli',
		metrics: {
			title: 'Metriklər',
			total_visits: 'Ümumi Ziyarətlər',
			unique_visitors: 'Unikal Ziyarətçilər',
			active_visitors: 'Aktiv Ziyarətçilər',
			error_rate: 'Səhv Faizi',
			average_latency: 'Orta Gecikmə',
			latency_95th_percentile: '95-ci Faiz Gecikməsi',
			latency_99th_percentile: '99-cu Faiz Gecikməsi',
			outgoing_traffic: 'Gedən Trafik',
			incoming_traffic: 'Gələn Trafik',
		},
		top_pages: {
			title: 'Ən Populyar Səhifələr',
			page_text: '{0} baxış (%{1})',
		},
		browser_distribution: {
			title: 'Brauzer Dağılımı',
			header: 'Brauzer',
		},
		device_distribution: {
			title: 'Cihaz Dağılımı',
			header: 'Cihaz',
		},
		referer_distribution: {
			title: 'İstinad Dağılımı',
			header: 'İstinad',
		},
		instance_distribution: {
			title: 'Server Dağılımı',
			header: 'Server',
		},
	},
	remote_actions: {
		title: 'Uzaqdan Əmrlər',
		header: 'Uzaqdan Əmrlər',
		loading: 'Yüklənir...',
		no_actions: 'Əmrlər təyin edilməyib.',
		error: 'Xəta',
		no_args: '(Argument yoxdur)',
		invoke: 'İcra et',
		units: {
			nsec: "Nanosanə",
			usec: "Mikrosaniyə",
			msec: "Millisaniyə",
			sec: "Saniyə",
			min: "Dəqiqə",
			hour: "Saat",
			day: "Gün",
			week: "Həftə",
			month: "Ay",
			year: "İl",
		},
	},
	user_sessions: {
		title: 'Sessiyalar',
		header: 'Sessiyalar',
		my_active_sessions: {
			title: 'Aktiv Sessiyalarım',
			device: 'Cihaz',
			last_activity: 'Son Aktivlik',
			created_at: 'Yaradılma Tarixi',
			you: '(Siz)',
			revoke: 'Ləğv et',
		},
		all_users: {
			title: 'Bütün İstifadəçilər',
			online: 'Aktiv',
			last_seen: 'Son görülüş',
			admin: 'Admin',
		}
	},
	settings: {
		title: 'Ayarlar',
		header: 'Ayarlar',
		theme: {
			title: 'Mövzu',
			select_theme: 'Mövzunu seçin',
			options: {
				system: 'Sistem Seçimi',
				light: 'İşıqlı',
				dark: 'Qaranlıq',
				dark_blue: 'Qara Mavi',
				light_green: 'Açık Yeşil',
				dark_green: 'Qara Yeşil',
				light_purple: 'Açık Mor',
				dark_purple: 'Qara Mor',
				light_pink: 'Açık Pənklə',
				dark_pink: 'Qara Pənklə',
			},
			current_theme: 'Cari Mövzu: {0}',
			theme_description: 'Sistem Seçimi cihazınızın tənzimləmələrini izləyəcək.',
		},
		language: {
			title: 'Dil və Region',
			language: {
				title: 'Dil',
			}
		},
		profile: {
			title: 'Profil',
			username: 'İstifadəçi adı',
			display_name: 'Görüntü adı',
			save: 'Yadda saxla',
		},
		account: {
			title: 'Hesab İdarəetməsi',
			change_password: 'Şifrəni dəyiş',
			delete_account: 'Hesabı sil',
		},
	},
	help: {
		title: 'Kömək',
		header: 'Logar Kömək Bələdçisi',
		about_logar: {
			title: 'Logar Haqqında',
			content: 'Logar — Go tətbiqləri üçün yüngül və çevik jurnal kitabxanasıdır. Monitorinq və analiz üçün veb əsaslı interfeys təqdim edir.'
		},
		integration_guide: {
			title: 'İnteqrasiya Bələdçisi',
			content: 'Budur Logar paketinin Go istinad məlumatları'
		},
		troubleshooting: {
			title: 'Problemlərin Həlli',
			pre_twitter: 'Əgər hər hansı problem varsa, mənimlə Twitter-də əlaqə saxlayın:',
			pre_github: 'və ya GitHub-da məsələ açın',
			github_text: 'GitHub anbarı',
		},
		additional_resources: {
			title: 'Əlavə Məlumatlar',
			github_repository: 'GitHub Anbarı',
			api_docs: 'API Sənədləri',
		}
	},
} satisfies BaseTranslation

export default az
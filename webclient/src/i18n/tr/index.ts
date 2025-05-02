import en from '../en'
import type { BaseTranslation, Translation } from '../i18n-types'

const tr = {
	...en as Translation,
	
	login: {
		title: 'Giriş',
		username: 'Kullanıcı Adı',
		password: 'Şifre',
		login: 'Giriş Yap',
	},
	topbar: {
		my_account: 'Hesabım',
		sign_out: 'Çıkış Yap',
	},
	dashboard: {
		title: 'Gösterge Paneli',
	},
	logs: {
		title: 'Kayıtlar',
		live_stream_connected: 'Canlı Yayın Bağlandı',
		live_stream_disconnected: 'Canlı Yayın Kesildi',
		reconnect: 'Tekrar Bağlan',
		reconnect_tip: 'Canlı yayını tekrar bağlamaya çalış',
		fields: {
			id: 'ID',
			timestamp: 'Zaman Damgası',
			category: 'Kategori',
			message: 'Mesaj',
			severity: 'Önemlilik',
		},
		operators: {
			contains: 'İçeriyor',
			not_contains: 'İçermiyor',
			starts_with: '... ile başlayan',
			ends_with: '... ile biten',
			between: 'Arasında',
			not_between: 'Arasında değil',
			in: 'İçinde',
			not_in: 'İçinde değil',
		},
		value: 'Değer',
		values_comma_separated: 'Değerler (virgülle ayrılmış)',
		now: 'Şimdi',
		from: 'Başlangıç',
		to: 'Bitiş',
		add_filter: 'Filtre Ekle',
		time_format_tip: 'İpucu: DD-MM-YYYY HH:mm:ss.SSS formatını kullanın',
		severity_tip: 'Ciddiyetler: 1 (Trace), 2 (Log), 3 (Info), 4 (Warn), 5 (Error), 6 (Fatal)'
	},
	analytics: {
		title: 'Analitik',
		header: 'Analitik Paneli',
	},
	remote_actions: {
		title: 'Uzak Eylemler',
		header: 'Uzak Eylemler',
		loading: 'Yükleniyor...',
		no_actions: 'Tanımlı eylem yok',
		error: 'Hata',
		no_args: '(Argüman yok)',
		invoke: 'Çağır',
		units: {
			nsec: "Nanosaniye",
			usec: "Mikrosaniye",
			msec: "Milisaniye",
			sec: "Saniye",
			min: "Dakika",
			hour: "Saat",
			day: "Gün",
			week: "Hafta",
			month: "Ay",
			year: "Yıl",
		},
	},
	user_sessions: {
		title: 'Oturumlar',
		header: 'Oturumlar',
		my_active_sessions: {
			title: 'Aktif Oturumlarım',
			device: 'Cihaz',
			last_activity: 'Son Aktiflik',
			created_at: 'Oluşturulma Zamanı',
			you: '(Siz)',
			revoke: 'İptal Et',
		},
		all_users: {
			title: 'Tüm Kullanıcılar',
			online: 'Aktif',
			last_seen: 'Son görülme',
			admin: 'Yönetici',
		}
	},
	settings: {
		title: 'Ayarlar',
		header: 'Ayarlar',
		theme: {
			title: 'Tema',
			select_theme: 'Tema Seçin',
			options: {
				light: 'Açık',
				dark: 'Koyu',
				system: 'Sistem Tercihi',
			},
			current_theme: 'Geçerli Tema: {0}',
			theme_description: 'Sistem tercihi, cihaz ayarlarınıza göre ayarlanır.',
		},
		language: {
			title: 'Dil ve Bölge',
			language: {
				title: 'Dil',
			}
		},
		profile: {
			title: 'Profil',
			username: 'Kullanıcı Adı',
			display_name: 'Görünen Ad',
			save: 'Kaydet',
		},
		account: {
			title: 'Hesap Yönetimi',
			change_password: 'Şifreyi Değiştir',
			delete_account: 'Hesabı Sil',
		},
	},
	help: {
		title: 'Yardım',
		header: 'Logar Yardım Kılavuzu',
		about_logar: {
			title: 'Logar Hakkında',
			content: 'Logar, Go uygulamaları için hafif, esnek bir günlükleme kütüphanesidir. Web tabanlı bir arayüzle izleme ve analiz için kapsamlı bir çözüm sunar.'
		},
		integration_guide: {
			title: 'Entegrasyon Kılavuzu',
			content: 'Logar paketinin Go referansı'
		},
		troubleshooting: {
			title: 'Sorun Giderme',
			pre_twitter: 'Bir sorunla karşılaşırsanız, lütfen Twitter üzerinden benimle iletişime geçin:',
			pre_github: 'veya GitHub\'da bir issue oluşturun:',
			github_text: 'GitHub repo',
		},
		additional_resources: {
			title: 'Ek Kaynaklar',
			github_repository: 'GitHub Deposu',
			api_docs: 'API Dokümantasyonu',
		}
	},
} satisfies BaseTranslation

export default tr

import type { BaseTranslation } from '../i18n-types'

const en = {
	login: {
		title: 'Login',
		username: 'Username',
		password: 'Password',
		login: 'Login',
	},
	topbar: {
		my_account: 'My Account',
		sign_out: 'Sign Out',
	},
	dashboard: {
		title: 'Dashboard',
	},
	logs: {
		title: 'Logs',
		live_stream_connected: 'Live Stream Connected',
		live_stream_disconnected: 'Live Stream Disconnected',
		reconnect: 'Reconnect',
		reconnect_tip: 'Attempt to reconnect the live stream',
		fields: {
			id: 'ID',
			timestamp: 'Timestamp',
			category: 'Category',
			message: 'Message',
			severity: 'Severity',
		},
		operators: {
			contains: 'Contains',
			not_contains: 'Not Contains',
			starts_with: 'Starts With',
			ends_with: 'Ends With',
			between: 'Between',
			not_between: 'Not Between',
			in: 'In',
			not_in: 'Not In',
		},
		value: 'Value',
		values_comma_separated: 'Values (comma-separated)',
		now: 'Now',
		from: 'From',
		to: 'To',
		add_filter: 'Add Filter',
		time_format_tip: 'Tip: Use format DD-MM-YYYY HH:mm:ss.SSS',
		severity_tip: 'Severities: 1 (Trace), 2 (Log), 3 (Info), 4 (Warn), 5 (Error), 6 (Fatal)'
	},
	analytics: {
		title: 'Analytics',
		header: 'Analytics Dashboard',
	},
	remote_actions: {
		title: 'Remote Actions',
		header: 'Remote Actions',
		loading: 'Loading...',
		no_actions: 'No actions defined.',
		error: 'Error',
		no_args: '(No arguments)',
		invoke: 'Invoke',
		units: {
			nsec: "Nanosecond",
			usec: "Microsecond",
			msec: "Millisecond",
			sec: "Second",
			min: "Minute",
			hour: "Hour",
			day: "Day",
			week: "Week",
			month: "Month",
			year: "Year",
		},
	},
	user_sessions: {
		title: 'Sessions',
		header: 'Sessions',
		my_active_sessions: {
			title: 'My Active Sessions',
			device: 'Device',
			last_activity: 'Last Activity',
			created_at: 'Created At',
			you: '(You)',
			revoke: 'Revoke',
		},
		all_users: {
			title: 'All Users',
			online: 'Online',
			last_seen: 'Last seen',
			admin: 'Admin',
		}
	},
	settings: {
		title: 'Settings',
		header: 'Settings',
		theme: {
			title: 'Theme',
			select_theme: 'Select Theme',
			options: {
				light: 'Light',
				dark: 'Dark',
				system: 'System Preference',
			},
			current_theme: 'Current Theme: {0}',
			theme_description: 'System preference will follow your device settings.',
		},
		language: {
			title: 'Language & Region',
			language: {
				title: 'Language',
			}
		},
		profile: {
			title: 'Profile',
			username: 'Username',
			display_name: 'Display Name',
			save: 'Save',
		},
		account: {
			title: 'Account Management',
			change_password: 'Change Password',
			delete_account: 'Delete Account',
		},
	},
	help: {
		title: 'Help',
		header: 'Logar Help Guide',
		about_logar: {
			title: 'About Logar',
			content: 'Logar is a lightweight, flexible logging library for Go applications that provides a comprehensive solution for application logging with a web-based interface for monitoring and analysis.'
		},
		integration_guide: {
			title: 'Integration Guide',
			content: 'Here is the Go reference of the Logar package'
		},
		troubleshooting: {
			title: 'Troubleshooting',
			pre_twitter: 'If you have any issues, please contact me on twitter:',
			pre_github: 'or open an issue on the',
			github_text: 'GitHub repository',
		},
		additional_resources: {
			title: 'Additional Resources',
			github_repository: 'GitHub Repository',
			api_docs: 'API Documentation',
		}
	},
} satisfies BaseTranslation

export default en

import en from '../en'
import type { BaseTranslation, Translation } from '../i18n-types'

const zh = {
	...en as Translation,
	
	login: {
		title: '登录',
		username: '用户名',
		password: '密码',
		login: '登录',
	},
	topbar: {
		my_account: '我的账户',
		sign_out: '退出',
	},
	dashboard: {
		title: '仪表板',
	},
	logs: {
		title: '日志',
		live_stream_connected: '直播已连接',
		live_stream_disconnected: '直播已断开',
		reconnect: '重新连接',
		reconnect_tip: '尝试重新连接直播流',
		fields: {
			id: 'ID',
			timestamp: '时间戳',
			category: '类别',
			message: '消息',
			severity: '严重程度',
		},
		operators: {
			contains: '包含',
			not_contains: '不包含',
			starts_with: '以...开始',
			ends_with: '以...结束',
			between: '介于之间',
			not_between: '不在之间',
			in: '在列表中',
			not_in: '不在列表中',
		},
		value: '值',
		values_comma_separated: '值（用逗号分隔）',
		now: '现在',
		from: '从',
		to: '到',
		add_filter: '添加过滤器',
		filters: '过滤器',
		load_more: '加载更多',
		loading: '加载中...',
		no_more_logs: '没有更多日志可显示',
		time_format_tip: '提示：使用格式 DD-MM-YYYY HH:mm:ss.SSS',
		severity_tip: '严重等级：1（追踪）、2（日志）、3（信息）、4（警告）、5（错误）、6（致命）'
	},
	analytics: {
		title: '分析',
		header: '分析仪表盘',
		metrics: {
			title: '指标',
			total_visits: '总访问量',
			unique_visitors: '独立访客数',
			active_visitors: '活跃访客数',
			error_rate: '错误率',
			average_latency: '平均延迟',
			latency_95th_percentile: '95百分位延迟',
			latency_99th_percentile: '99百分位延迟',
			outgoing_traffic: '出站流量',
			incoming_traffic: '入站流量',
		},
		top_pages: {
			title: '热门页面',
			page_text: '{0} 次浏览（{1}%）',
		},
		browser_distribution: {
			title: '浏览器分布',
			header: '浏览器',
		},
		device_distribution: {
			title: '设备分布',
			header: '设备',
		},
		referer_distribution: {
			title: '来源分布',
			header: '来源',
		},
		instance_distribution: {
			title: '服务器分布',
			header: '服务器',
		},
	},
	remote_actions: {
		title: '远程操作',
		header: '远程操作',
		loading: '加载中...',
		no_actions: '未定义操作。',
		error: '错误',
		no_args: '（没有参数）',
		invoke: '调用',
		units: {
			nsec: '纳秒',
			usec: '微秒',
			msec: '毫秒',
			sec: '秒',
			min: '分钟',
			hour: '小时',
			day: '天',
			week: '周',
			month: '月',
			year: '年',
		},
	},
	user_sessions: {
		title: '会话',
		header: '会话',
		my_active_sessions: {
			title: '我的活跃会话',
			device: '设备',
			last_activity: '上次活动时间',
			created_at: '创建时间',
			you: '(你)',
			revoke: '撤销',
		},
		all_users: {
			title: '所有用户',
			online: '在线',
			last_seen: '上次在线',
			admin: '管理员',
		}
	},
	settings: {
		title: '设置',
		header: '设置',
		theme: {
			title: '主题',
			select_theme: '选择主题',
			options: {
				system: '系统偏好',
				light: '浅色',
				dark: '深色',
				dark_blue: '深蓝色',
				light_green: '浅绿色',
				dark_green: '深绿色',
				light_purple: '浅紫色',
				dark_purple: '深紫色',
				light_pink: '浅粉色',
				dark_pink: '深粉色',
			},
			current_theme: '当前主题：{0}',
			theme_description: '系统偏好将遵循您的设备设置',
		},
		language: {
			title: '语言和地区',
			language: {
				title: '语言',
			}
		},
		profile: {
			title: '个人资料',
			username: '用户名',
			display_name: '显示名称',
			save: '保存',
		},
		account: {
			title: '账户管理',
			change_password: '更改密码',
			delete_account: '删除账户',
		},
	},
	help: {
		title: '帮助',
		header: 'Logar 使用指南',
		about_logar: {
			title: '关于 Logar',
			content: 'Logar 是一个轻量且灵活的 Go 应用日志库，提供基于 Web 的界面用于日志监控与分析的全面解决方案。'
		},
		integration_guide: {
			title: '集成指南',
			content: '以下是 Logar 包的 Go 语言参考'
		},
		troubleshooting: {
			title: '故障排查',
			pre_twitter: '如果遇到问题，请通过 Twitter 联系我：',
			pre_github: '或在此提交问题：',
			github_text: 'GitHub 仓库',
		},
		additional_resources: {
			title: '其他资源',
			github_repository: 'GitHub 仓库',
			api_docs: 'API 文档',
		}
	},
} satisfies BaseTranslation

export default zh
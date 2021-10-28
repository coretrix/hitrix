module.exports = {
    base: '/hitrix/',
    title: 'Hitrix - golang framework',
    description: 'Golang Framework for high traffic applications. Designed for speed up development time.',
    head: [
        ['link', {rel: "shortcut icon", href: "logo-favicon.png"}],
        ['meta', {name: 'theme-color', content: '#D7A318'}],
        ['meta', {name: 'apple-mobile-web-app-capable', content: 'yes'}],
        ['meta', {name: 'apple-mobile-web-app-status-bar-style', content: 'black'}]
    ],
    themeConfig: {
        repo: 'https://github.com/coretrix/hitrix',
        docsRepo: 'https://github.com/coretrix/hitrix',
        logo: '/logo-favicon-90x90.png',
        editLinks: true,
        docsDir: 'docs/docs',
        editLinkText: '',
        lastUpdated: true,
        smoothScroll: true,
        navbar: [
            {
                text: 'Guide',
                link: '/',
            },
            {
                text: 'Roadmap',
                link: '/roadmap/'
            },
            {
                text: 'CoreTrix Rules',
                link: '/rules/'
            },
        ],
        sidebar: [
            {
                text: 'Guide',
                children: [
                    {
                        text: 'Introduction',
                        collapsible: false,
                        link: '/',
                    },
                    {
                        text: 'Services',
                        children: [
                            {
                                text: 'App',
                                link: '/guide/services/app',
                            },
                            {
                                text: 'Config',
                                link: '/guide/services/config',
                            },
                            {
                                text: 'ORM Engine',
                                link: '/guide/services/orm_engine',
                            },
                            {
                                text: 'ORM Engine Context',
                                link: '/guide/services/orm_engine_context',
                            },
                            {
                                text: 'Amazon S3',
                                link: '/guide/services/amazon_s3',
                            },
                            {
                                text: 'OSS - Google',
                                link: '/guide/services/oss_google',
                            },
                            {
                                text: 'API Logger',
                                link: '/guide/services/api_logger',
                            },
                            {
                                text: 'Authentication',
                                link: '/guide/services/authentication',
                            },
                            {
                                text: 'Clock',
                                link: '/guide/services/clock',
                            },
                            {
                                text: 'Checkout',
                                link: '/guide/services/checkout',
                            },
                            {
                                text: 'CRUD',
                                link: '/guide/services/crud',
                            },
                            {
                                text: 'DDOS',
                                link: '/guide/services/ddos',
                            },
                            {
                                text: 'Dynamic link',
                                link: '/guide/services/dynamic_link',
                            },
                            {
                                text: 'Error logger',
                                link: '/guide/services/error_logger',
                            },
                            {
                                text: 'Firebase cloud messaging',
                                link: '/guide/services/fcm',
                            },
                            {
                                text: 'File extractor',
                                link: '/guide/services/file_extractor',
                            },
                            {
                                text: 'Localizer',
                                link: '/guide/services/localizer',
                            },
                            {
                                text: 'Mandrill - Transactional emails',
                                link: '/guide/services/mandrill',
                            },
                            {
                                text: 'JWT',
                                link: '/guide/services/jwt',
                            },
                            {
                                text: 'Password',
                                link: '/guide/services/password',
                            },
                            {
                                text: 'PDF',
                                link: '/guide/services/pdf',
                            },
                            {
                                text: 'Slack',
                                link: '/guide/services/slack',
                            },
                            {
                                text: 'SMS',
                                link: '/guide/services/sms',
                            },
                            {
                                text: 'Setting',
                                link: '/guide/services/setting',
                            }, {
                                text: 'Stripe',
                                link: '/guide/services/stripe',
                            },
                            {
                                text: 'Uploader',
                                link: '/guide/services/uploader',
                            },
                            {
                                text: 'WebSocket',
                                link: '/guide/services/websocket',
                            },
                            {
                                text: 'Exporter',
                                link: '/guide/services/exporter',
                            },
                        ],
                    },
                    {
                        text: 'Features',
                        children: [
                            {
                                text: 'Flags',
                                link: '/guide/features/flags',
                            },
                            {
                                text: 'Background scripts',
                                link: '/guide/features/script',
                            },
                            {
                                text: 'Seeder',
                                link: '/guide/features/seeder',
                            },
                            {
                                text: 'Validator',
                                link: '/guide/features/validator',
                            },
                            {
                                text: 'Pagination',
                                link: '/guide/features/pagination',
                            },
                            {
                                text: 'Integration test',
                                link: '/guide/features/test',
                            },
                            {
                                text: 'Helper',
                                link: '/guide/features/helper',
                            },
                            {
                                text: 'Goroutine',
                                link: '/guide/features/goroutine',
                            },
                        ],
                    },
                ],
            },
        ],

        plugins: [
            '@vuepress/plugin-search',
            {
                searchMaxSuggestions: 10
            },
            '@vuepress/plugin-back-to-top',
            '@vuepress/plugin-medium-zoom',
            [
                'vuepress-plugin-sitemap',
                {hostname: 'https://coretrix.github.io/hitrix'}
            ],
            // [
            //     '@vuepress/google-analytics',
            //     {
            //         'ga': ''
            //     }
            // ]
        ]
    }
}

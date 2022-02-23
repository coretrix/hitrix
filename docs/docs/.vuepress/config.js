const config = require('./config/index')

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
        algolia: config.Algolia,
        navbar: config.Navigation,
        sidebar: config.Sidebar,
    },
    plugins: [
        ['@vuepress/plugin-search', config.Search],
        ['@vuepress/plugin-back-to-top', true],
        ['@vuepress/plugin-medium-zoom', true],
        ['vuepress-plugin-sitemap', { hostname: 'https://coretrix.github.io/hitrix' }],
        // ['@vuepress/google-analytics', { 'ga': ''}]
    ]
}

module.exports = {
    locales: {
        '/': {
            placeholder: 'Search...'
        }
    },
    maxSuggestions: 10,
    getExtraFields: (page) => page.frontmatter.tags ?? []
}

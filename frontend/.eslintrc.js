module.exports = {
    root: true,

    env: {
        node: true
    },

    extends: [
        // https://vuejs.github.io/eslint-plugin-vue/rules/#priority-c-recommended-minimizing-arbitrary-choices-and-cognitive-overhead
        'plugin:vue/recommended',
        // https://github.com/vuejs/eslint-config-standard
        '@vue/standard'
    ],

    rules: {
        'no-console': 'off',
        'no-debugger': 'off',
        'vue/max-attributes-per-line': ['error', {
            'singleline': 10,
            'multiline': {
                'max': 3,
                'allowFirstLine': true
            }
        }],
        'camelcase': 'off',
        'quotes': ['error', 'single'],
        'semi': ['error', 'never'],
        'max-len': ['error', { 'code': 150, 'ignoreStrings': true }],
        "indent": "off",
        "no-tabs": "off",
    },

    parserOptions: {
        parser: 'babel-eslint',
    }
}

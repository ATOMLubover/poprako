import js from '@eslint/js';
import tsPlugin from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import vuePlugin from 'eslint-plugin-vue';
import vueParser from 'vue-eslint-parser';

export default [{
        ignores: ['dist/**', 'node_modules/**', 'src/**/*.d.ts', 'src/**/*.js'],
    },
    {
        files: ['src/**/*.{ts,vue}'],
        languageOptions: {
            parser: vueParser,
            parserOptions: {
                parser: tsParser,
                ecmaVersion: 'latest',
                sourceType: 'module',
                extraFileExtensions: ['.vue'],
            },
            globals: {
                window: 'readonly',
                document: 'readonly',
                localStorage: 'readonly',
                location: 'readonly',
                URLSearchParams: 'readonly',
            },
        },
        plugins: {
            '@typescript-eslint': tsPlugin,
            vue: vuePlugin,
        },
        rules: {
            ...js.configs.recommended.rules,
            ...tsPlugin.configs.recommended.rules,
            ...vuePlugin.configs['vue3-recommended'].rules,
            'vue/multi-word-component-names': 'off',
            '@typescript-eslint/no-explicit-any': 'off',
        },
    },
];
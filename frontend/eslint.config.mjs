import react from 'eslint-plugin-react'
import prettier from 'eslint-plugin-prettier'
import globals from 'globals'
import tsParser from '@typescript-eslint/parser'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import js from '@eslint/js'
import { FlatCompat } from '@eslint/eslintrc'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
})

const config = [
  {
    ignores: [
      'node_modules',
      '.pnp',
      '**/.pnp.js',
      '.yarn/install-state.gz',
      'coverage',
      '.next/',
      'out/',
      'build',
      '**/.DS_Store',
      '**/*.pem',
      '**/npm-debug.log*',
      '**/yarn-debug.log*',
      '**/yarn-error.log*',
      '**/.env*.local',
      '**/.vercel',
      '**/*.tsbuildinfo',
      '**/next-env.d.ts',
      '**/.env*',
      '!**/.env.sample',
      '**/eslintrc.config.mjs',
      '**/pnpm-lock.yaml',
    ],
  },
  ...compat.extends(
    'plugin:react/recommended',
    'next/core-web-vitals',
    'prettier'
  ),
  {
    plugins: { react, prettier },
    languageOptions: {
      globals: { ...globals.browser },
      parser: tsParser,
      ecmaVersion: 'latest',
      sourceType: 'module',
      parserOptions: { project: './tsconfig.json' },
    },
    rules: {
      '@typescript-eslint/consistent-type-definitions': 'off',
      '@typescript-eslint/consistent-type-imports': 'off',
      '@typescript-eslint/explicit-function-return-type': 'off',
      '@typescript-eslint/no-misused-promises': 'off',
      '@typescript-eslint/strict-boolean-expressions': 'off',
      '@typescript-eslint/prefer-nullish-coalescing': 'off',
      '@typescript-eslint/no-confusing-void-expression': 'off',
      'object-shorthand': 'off',
    },
  },
]

export default config

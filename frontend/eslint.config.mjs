import js from '@eslint/js'
import { FlatCompat } from '@eslint/eslintrc'
import { includeIgnoreFile } from '@eslint/compat'
import { defineConfig, globalIgnores } from 'eslint/config'
import { fileURLToPath } from 'node:url'
import stylistic from '@stylistic/eslint-plugin'

const gitignorePath = fileURLToPath(new URL('.gitignore', import.meta.url))

export default defineConfig([
  js.configs.recommended,
  stylistic.configs.recommended,
  new FlatCompat().extends('next/core-web-vitals', 'next/typescript'),

  includeIgnoreFile(gitignorePath, 'Imported .gitignore patterns'),
  globalIgnores(['*.config.*']),
])

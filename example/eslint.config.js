import js from '@eslint/js'
import { FlatCompat } from '@eslint/eslintrc'
import { includeIgnoreFile } from "@eslint/compat";
import { defineConfig, globalIgnores } from 'eslint/config';
import { fileURLToPath } from 'node:url';
import stylistic from '@stylistic/eslint-plugin'

const gitignorePath = fileURLToPath(new URL('.gitignore', import.meta.url));

export default defineConfig([
  new FlatCompat().extends('next/core-web-vitals', 'next/typescript'),
  js.configs.recommended,
  stylistic.configs.recommended,

  includeIgnoreFile(gitignorePath, 'Imported .gitignore patterns'),
  globalIgnores([
    '*.config.*',
  ]),
])

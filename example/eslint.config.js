import js from '@eslint/js'
import { FlatCompat } from '@eslint/eslintrc'
import stylistic from '@stylistic/eslint-plugin'

const compat = new FlatCompat({})
const config = [
  ...compat.extends('next/core-web-vitals', 'next/typescript'),
  stylistic.configs.recommended,
  js.configs.recommended,
]
export default config

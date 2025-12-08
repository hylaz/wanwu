// Documentation for this file: https://prettier.io/docs/en/options.html
module.exports = {
  printWidth: 80,
  endOfLine: 'auto',
  semi: true,
  tabWidth: 2,
  useTabs: false,
  singleQuote: true,
  trailingComma: 'all',
  bracketSpacing: true,
  jsxBracketSameLine: false,
  arrowParens: 'avoid',
  overrides: [
    {
      files: '.prettierrc',
      options: { parser: 'json', trailingComma: 'none' },
    },
    {
      files: '.babelrc',
      options: { parser: 'json', trailingComma: 'none' },
    },
    {
      files: '*.json',
      options: { trailingComma: 'none' },
    },
  ],
  htmlWhitespaceSensitivity: 'ignore',
};

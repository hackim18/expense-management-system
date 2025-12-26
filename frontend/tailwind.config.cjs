const daisyui = require('daisyui')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./app/**/*.{vue,js,ts}', './components/**/*.{vue,js,ts}'],
  theme: {
    extend: {
      boxShadow: {
        soft: '0 24px 60px -40px rgba(31, 26, 22, 0.6)',
        glow: '0 14px 50px -20px rgba(15, 118, 110, 0.5)'
      }
    }
  },
  plugins: [daisyui],
  daisyui: {
    themes: [
      {
        rupate: {
          primary: '#d97706',
          secondary: '#0f766e',
          accent: '#1f1a16',
          neutral: '#1f1a16',
          'base-100': '#fffaf2',
          'base-200': '#f6efe4',
          'base-300': '#e8dfd2',
          info: '#38bdf8',
          success: '#10b981',
          warning: '#f59e0b',
          error: '#f43f5e'
        }
      }
    ]
  }
}

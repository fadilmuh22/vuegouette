import colors from 'tailwindcss/colors'

/** @type {import('tailwindcss').Config} */

const colorsMap = {
  default: {
    // darkish
    primary: '#111827',
    secondary: '#1f2937',
    background: '#131618',
    onbackground: '#F1F3F5',
    'neutral-light': '#F1F3F5',
    'neutral-dark': '#34495E',
  },
  senti: {
    // greenish
    primary: '#064e3b',
    secondary: '#065f46',
    background: '#F1F3F5',
    onbackground: '#34495E',
    'neutral-light': '#F1F3F5',
    'neutral-dark': '#34495E',
  },
  cv: {
    // blueish
    primary: '#2a0666',
    secondary: '#4b4453',
    background: '#F1F3F5',
    onbackground: '#34495E',
    'neutral-light': '#F1F3F5',
    'neutral-dark': '#34495E',
  },
}

export default {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  content: [],
  darkMode: 'selector',
  theme: {
    extend: {
      screens: {
        xs: '320px',
        sm: '640px',
        md: '768px',
        lg: '1024px',
        xl: '1280px',
        '2xl': '1536px',
      },
    },
    colors: {
      ...colorsMap['default'],
      transparent: 'transparent',
      current: 'currentColor',
      black: colors.black,
      white: colors.white,
      gray: colors.gray,
      emerald: colors.emerald,
      indigo: colors.indigo,
      yellow: colors.yellow,
    },
  },
  plugins: [],
}

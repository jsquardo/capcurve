/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        surface: {
          DEFAULT: '#0f1117',
          raised: '#161b27',
          overlay: '#1e2535',
        },
        brand: {
          DEFAULT: '#3b82f6',
          dim: '#1d4ed8',
        },
        peak: '#f59e0b',
        value: '#10b981',
        overpaid: '#ef4444',
        neutral: '#6b7280',
      },
      fontFamily: {
        display: ['var(--font-display)', 'sans-serif'],
        body: ['var(--font-body)', 'sans-serif'],
        mono: ['var(--font-mono)', 'monospace'],
      },
    },
  },
  plugins: [],
}

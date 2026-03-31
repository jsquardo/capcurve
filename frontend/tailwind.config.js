/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        app: 'rgb(var(--color-bg) / <alpha-value>)',
        elevated: 'rgb(var(--color-bg-elevated) / <alpha-value>)',
        panel: 'rgb(var(--color-bg-panel) / <alpha-value>)',
        overlay: 'rgb(var(--color-bg-overlay) / <alpha-value>)',
        text: {
          DEFAULT: 'rgb(var(--color-text) / <alpha-value>)',
          muted: 'rgb(var(--color-text-muted) / <alpha-value>)',
          subtle: 'rgb(var(--color-text-subtle) / <alpha-value>)',
        },
        border: {
          DEFAULT: 'rgb(var(--color-border) / <alpha-value>)',
          strong: 'rgb(var(--color-border-strong) / <alpha-value>)',
        },
        accent: {
          DEFAULT: 'rgb(var(--color-accent) / <alpha-value>)',
          strong: 'rgb(var(--color-accent-strong) / <alpha-value>)',
        },
        link: 'rgb(var(--color-link) / <alpha-value>)',
        success: 'rgb(var(--color-success) / <alpha-value>)',
        danger: 'rgb(var(--color-danger) / <alpha-value>)',
        projection: 'rgb(var(--color-projection) / <alpha-value>)',
        surface: {
          DEFAULT: 'rgb(var(--color-bg) / <alpha-value>)',
          raised: 'rgb(var(--color-bg-elevated) / <alpha-value>)',
          overlay: 'rgb(var(--color-bg-overlay) / <alpha-value>)',
        },
        brand: {
          DEFAULT: 'rgb(var(--color-link) / <alpha-value>)',
          dim: 'rgb(var(--color-link) / 0.72)',
        },
        peak: 'rgb(var(--color-accent) / <alpha-value>)',
        value: 'rgb(var(--color-success) / <alpha-value>)',
        overpaid: 'rgb(var(--color-danger) / <alpha-value>)',
        neutral: 'rgb(var(--color-text-muted) / <alpha-value>)',
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

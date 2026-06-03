/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        pertamina: '#003d7a',
        shell: '#e81c1c',
        vivo: '#3b82f6',
        bp: '#00a651',
        primary: '#667eea',
        secondary: '#764ba2',
      },
    },
  },
  plugins: [],
}

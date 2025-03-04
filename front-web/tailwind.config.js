/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
     "./src/App.tsx" 
  ],
  theme: {
    extend: {},
  },
  plugins: [],

  mode: 'jit', 
  purge: false,
}
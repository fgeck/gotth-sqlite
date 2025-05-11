/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [ "/**/*.html", "./**/*.templ", "./**/*.go", ],
    safelist: [],
    theme: {
        extend: {
          colors: {
            gray: {
              700: '#374151', // Dark mode sidebar background
              800: '#1f2937', // Dark mode body background
            },
            blue: {
              400: '#60a5fa', // Dark mode accent
              600: '#2563eb',  // Light mode accent
            }
          }
        }
    }
}

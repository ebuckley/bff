/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
        "../pkg/server/templates.go",
    ],
    theme: {
        extend: {},
    },
    plugins: [],
}

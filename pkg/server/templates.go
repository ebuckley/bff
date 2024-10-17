package server

import "html/template"

var index = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
 <title>{{.Heading}}</title>
 <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-100 p-4">
 <h1 class="text-3xl font-bold mb-4">{{.Heading}}</h1>
 <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
   {{range .Actions}}
     <a href="/a/{{.Slug}}" class="bg-blue-500 text-white text-lg font-bold py-4 px-6 rounded-lg hover:bg-blue-700">
       {{.Name}}
     </a>
   {{end}}
 </div>
</body>
</html>
`))

var action = template.Must(template.New("action").Parse(`
<!DOCTYPE html>
<html>
<head>
 <title>{{.Heading}}</title>
 <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-100 p-4">
 <h1 class="text-3xl font-bold mb-4">{{.Heading}}</h1>
 <div id="app"></div>
</body>
</html>
`))

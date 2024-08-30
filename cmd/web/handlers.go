package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Server", "Go")
	/*
		Добавляет заголовок к карте заголовков

		1. Заголовки должны быть добавлены строго до вызова w.WriteHeader() или w.Write().
	*/

	var files = []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/home.gohtml",
	}
	/*
		Инициализируется фрагмент, содержащий пути к двум файлам.

		1. Файл, содержащий базовый шаблон, должен быть *первым* файлом во фрагменте.
	*/

	ts, err := template.ParseFiles(files...)
	/*
		Считывает файлы шаблона в набор шаблонов

		1. Путь к файлу должен быть либо относительным к текущему рабочему каталогу, либо абсолютным путем.
	*/
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(writer, "base", nil) // записывает содержимое шаблона в качестве текста ответа.
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func snippetView(writer http.ResponseWriter, request *http.Request) {
	var id, err = strconv.Atoi(request.PathValue("id")) // извлечение значения из сегмента с подстановочными знаками
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d...\n", id)
	/*
		Поскольку http.ResponseWriter реализует интерфейс io.Writer, вместо
		var msg = fmt.Sprintf("Display a specific snippet with ID %d...\n", id)
		writer.Write([]byte(msg))

		можно использовать любой метод, аргументом которого является io.Writer:
		io.WriteString(w, msg)
		fmt.Fprint(w, msg)
	*/
}

func snippetCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Display a form for creating a new snippet...\n"))
}

func snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusCreated)
	/*
		Устанавливает код состояния

		1. Вызов w.WriteHeader() возможен только один раз для каждого ответа.
		2. Если не вызывать w.WriteHeader() явно,то при первом вызове w.Write()
			автоматически будет отправлен 200 код состояния.
		3. Для отправки кода состояния отличного от 200, необходимо вызвать w.WriteHeader()
			перед вызовом w.Write().
	*/

	writer.Write([]byte("Save a new snippet...\n"))
}

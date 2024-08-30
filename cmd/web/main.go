package main

import (
	"log"
	"net/http"
)

/*
Основные компоненты web-приложения:
1. Обработчики: заполняют заголовки и тела HTTP-ответов
2. Маршрутизатор: сопоставляет шаблоны маршрутизации URL веб-приложения и соответствующих обработчиков
3. Web-сервер: прослушивает порт

Движение HTTP-запроса
Web-сервер -> Маршрутизатор -> Обработчик
*/

func main() {
	var mux = http.NewServeMux() // маршрутизатор
	/*
		Пакет http предоставляет функции http.Handle() и http.HandleFunc(),
		которые позволяют регистрировать маршруты без явного указания servemux.
		Это достигается с использованием переменной http.DefaultServeMux, которая инициализируется go автоматически.
		Использование http.DefaultServeMux, как и любой глобальной переменной не рекомендуется.

		При вызове http.ListenAndServe(":4000", nil), сервер будет использовать http.DefaultServeMux для маршрутизации.

		1. Приоритет при обработке шаблонов предоставляется более "конкретному" шаблону.
		2. Обработка шаблонов не зависит от их порядка объявления в коде.
		3. Когда два шаблона перекрываются (одинаково "конкретны") возникает panic.
			Например "/post/new/{id}" и "/post/{author}/latest" перекрываются,
			потому что они оба соответствуют пути запроса /post/new/latest,
			но неясно, какой из них должен иметь приоритет.
	*/

	var fileServer = http.FileServer(http.Dir("./ui/static/"))
	// Создает файловый сервер, который раздает файлы из директории ./static

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	// Оборачивает файловый сервер в StripPrefix, который удаляет /static из начала URL перед передачей его файловому серверу

	mux.HandleFunc("/{$}", home) // обработчик для '/' URL шаблона
	/*
		Если шаблон оканчивается на '/' (subtree path pattern),
		он будет найден всякий раз, когда начало пути URL запроса совпадает с путем поддерева.
		Этого эффекта можно избежать, используя '{$}' после последнего '/', например mux.HandleFunc("/{$}", home)

		Пути URL запросов автоматически очищаются.
	*/

	mux.HandleFunc("/snippet/view/{id}", snippetView)
	/*
		Сегменты с подстановочными знаками в шаблоне маршрута обозначаются подстановочным идентификатором внутри {}.
		Путь запроса может содержать любое непустое значение для сегментов с подстановочными знаками.

		Каждый сегмент пути может содержать только один подстановочный знак, заполняющий весь сегмент пути
	*/

	mux.HandleFunc("GET /snippet/create", snippetCreate)
	/*
		Если шаблон не оканчивается на '/', он будет найден только в том случае,
		если URL-адрес запроса полностью совпадает с шаблоном.

		'GET' определяет конкретный метод. Спецификатор метода должен быть написан в верхнем регистре,
		за ним следует хотя бы один пробел (допустимы как пробелы, так и табуляция).
		Спецификатор GET соответствует методам GET и HEAD.
	*/

	mux.HandleFunc("POST /snippet/create", snippetCreatePost)
	/*
		Использование curl
		GET:
		>>> curl -i localhost:4000/snippet/create
		HTTP/1.1 200 OK
		Date: Sat, 24 Aug 2024 16:27:54 GMT
		Content-Length: 44
		Content-Type: text/plain; charset=utf-8

		Display a form for creating a new snippet...

		POST:
		>>> curl -i -d "" localhost:4000/snippet/create
		HTTP/1.1 200 OK
		Date: Sat, 24 Aug 2024 16:29:05 GMT
		Content-Length: 21
		Content-Type: text/plain; charset=utf-8

		Save a new snippet...

		DELETE: (не реализован)
		>>> curl -i -X DELETE localhost:4000/snippet/create
		HTTP/1.1 405 Method Not Allowed
		Allow: GET, HEAD, POST
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff
		Date: Sat, 24 Aug 2024 16:32:33 GMT
		Content-Length: 19

		Method Not Allowed
	*/

	log.Println("Listening on :4000")

	var err = http.ListenAndServe(":4000", mux) // запуск web-сервера
	/*
		ошибка, которую вернет всегда http.ListenAndServe != nil

		1. addr должен быть в формате "host:port". При отсутствии хоста, сервер будет прослушивать все доступные сетевые интерфейсы компьютера.
		2. Допускается использование именованных портов, таких как ":http" или ":http-alt" вместо числа.
			В это случае функция `http.ListenAndServe()` попытается найти соответствующий номер порта в файле /etc/services при запуске сервера,
			возвращая ошибку, если совпадение не может быть найдено.
	*/

	log.Fatal(err)
}

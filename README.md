# grader
Сервис Web Client
 * html templates для пользователей + админка
 * rate limiter на загрузку файлов 
 * swagger для документации хэндлеров
 * JWT токен в куке
 * загрузка файла main.go для запуска review - хранение файла в s3 minio
 * паттерн transactional outbox, создание события в таблице outbox с дальнейшей передачей в очередь rabbit
 * обработка callback вызова от сервиса grader, авторизация через header Authorization - отдельный JWT

Сервис Queue Processor
* переиспользован с прошлой задачи, достаточно было поменять конфиг, переключить на rabbit с указанием названия target очереди и url для внешнего вызова

Сервис Grader
* handler для обработки запроса от queue processor, payload в теле запроса
* загрузка файла из s3 minio, локальное хранение с дальнейшей прокидкой в docker container 
* worker pool для запуска docker и запуска тестов
* callback вызов

Docker image 
* сборка с копией проекта
![img.png](img.png)
* при запуске прокидывается в нужную директорию по номеру HW загруженный файл из review и запускается тест
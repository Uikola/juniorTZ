Тестовое задание на Junior Golang Developer

Запуск:

1) Разместить директорию envs в директории config и добавить в ней dev.env файл со
следующим содержимым: </br>
PORT=:your_port </br>
CONN_STRING=host=db port=5432 user=postgres password=password dbname=juniortzdb sslmode=disable </br>
DRIVER_NAME=postgres </br>
TIMEOUT=your_timeout </br>
IDLE_TIMEOUT=your_idle_timeout </br>
</br>
2) Запустить: docker build -t juniortz -f Dockerfile .
3) Запустить: docker-compose up

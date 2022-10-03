<?php

//отправка запроса

$data = file_get_contents("https://yandex.ru");

file_put_contents('yandex.html', $data);

//считываем файл

$fileData = file_get_contents("./Docker/docker-compose.yml");

echo $fileData;
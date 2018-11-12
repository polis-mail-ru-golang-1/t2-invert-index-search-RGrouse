# Поиск релевантного файла

Утилита выполняющая поиск указанной фразы по указанным файлам. На вход принимает директорию с файлами и адрес интерфейса, который нужно слушать. В начале работы производит индексацию файлов, строит по каждому обратный индекс. Затем сравнивает файлы по наилучшему совпадению содержимого с поисковой фразой. Выводит на экран файлы, где были найдены токены из поисковой фразы в порядке наилучшего соответствия. Наилучшее соответствие - все токены из поисковой фразы встретились в файле наибольшее количество раз. Если ни одного токена из фразы не найдено в файле, файл не выводится.

Директория с файлами, интерфейс и уровень логирования задаются переменными среды (SDIR, LISTEN и LOG_LEVEL соответственно). Например: `SDIR="./search" LISTEN="127.0.0.1:8080" LOG_LEVEL="debug"`. 
Поисковая фраза задается в параметре 'q' GET-запроса по пути '/search'. Например: `127.0.0.1:8080/search?q=your%20query`.

Результат выводится в HTTP-ответе в виде:
```
- file1.txt; совпадений - 2
- file2.txt; совпадений - 1
```

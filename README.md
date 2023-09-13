# Todo.txt agenda view

This is a simple script to generate an agenda view from todo.txt files.
It will read all files listed in the file passed into the first argument `TODO_FILES`.
It will look for due: tags and sort them by date.

The todos will then be presented in the following format: 

```text
Weekly Agenda

[PAST DUE]

Wednesday, September 13, 2023

Thursday, September 14, 2023

Friday, September 15, 2023

Saturday, September 16, 2023

Sunday, September 17, 2023

Monday, September 18, 2023

Tuesday, September 19, 2023
```

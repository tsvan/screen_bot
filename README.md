# screen bot


## Requirements
### OpenCv
Общая инструкция по установке

https://gocv.io/getting-started/windows/

1)MinGW-W64 ставить `posix-seh` Лучше последнюю, в инструкции была указана 7.3.0, с ней была ошибка компиляции
с 8.3.0 всё собралось
2) Перед `chdir %GOPATH%\src\gocv.io\x\gocv` нужно скопировать папку из
`pkg\mod\gocv.io\x` поместить её в `%GOPATH%\src\gocv.io\x` и переименовать в `gocv`

`go env GOPATH` что бы узнать GOPATH

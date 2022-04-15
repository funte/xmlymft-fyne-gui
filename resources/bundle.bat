@REM 打包资源

rm bundled.go
fyne bundle -package resources msyh.ttc >> bundled.go
fyne bundle -package resources -append icon.png >> bundled.go

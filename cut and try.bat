rem -- C:\A\repo\LinkedIn_Location_Data_Cleaning\cut and try.bat

cd C:\A\repo\LinkedIn_Location_Data_Cleaning

"C:\Program Files (x86)\Notepad++\Notepad++.exe" C:\A\repo\LinkedIn_Location_Data_Cleaning\LinkedIn_Location_Data_Cleaning.go

cls

gofmt -w LinkedIn_Location_Data_Cleaning.go

go run LinkedIn_Location_Data_Cleaning.go < C:\A\repo\LinkedIn_Location_Data_Cleaning\stdin.txt > C:\A\repo\LinkedIn_Location_Data_Cleaning\stdout.txt

"C:\Program Files (x86)\Notepad++\Notepad++.exe" C:\A\repo\LinkedIn_Location_Data_Cleaning\stdout.txt

pause

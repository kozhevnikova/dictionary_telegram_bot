# itmeansbot: a Telegram instant messenger inline bot 
### Objective: _to provide users with definitions of English words._ 

# Installation
Use *git clone* to get the code from the remote repository to your laptop.    
Before starting the project, the developer must fill the main.go file with the necessary credentials.    

## API settings
Oxford Dictionary is used for the project.  
To start the code, the developer must get an application ID and application key from the official website and enter them in const() section in main.go file in lines number 20 and 21. 

## Bot settings
Get your token from BotFather and use it in main.go file in line number 36.   
` err := os.Setenv("token","your_token_here")
`

# Additional information 
https://medium.com/@jkozhevnikova/inline-dictionary-in-your-telegram-part-i-introduction-32990d8664b4
https://medium.com/@jkozhevnikova/inline-dictionary-in-your-telegram-part-ii-explanation-6935494acf9e 

# Author
Jane Kozhevnikova - jane.kozhevnikova@gmail.com

# Status
Inactive

## License
MIT


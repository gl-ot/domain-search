# Domain search
Searches for available `.com` domains with godaddy api.

### How it works
You wish to have domain e.g. `football`. Of course football.com already taken.

The script puts every word of the word list in the beginning and the end of your domain football
 and checks if domain is available and then saves result in file.

### How to use
* Put a list of emphasizing words into `words.txt` splitting them with a new line. (repo already has a default one)
* Download search-$os according to your OS and run `./search-$os -domain=football`
* After the program has finished `available.txt` will appear with result.

For more information run `./search-$os --help`
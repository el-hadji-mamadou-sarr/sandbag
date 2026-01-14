1. some packets might be lost
2. it does not
3. pas de réponse pour le serveur
4. dans le kernel

the UDP is very fast but not reliable, sometimes a receivefrom can hang until the server give error


- UDP is fast
- no connexion state like in tcp and fail silently,
- set timeout 

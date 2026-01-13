1. donc le round robin permet de distribuer le trafic sur les serveurs de maniere séquentiel. donc si les tous les serveurs sont identiques, donc chaque serveur aura le méme trafic.
2. round robin fails si les servuers sont différents;
3. la latence de queue permet c'est le 1% de requétes qui sont plus lentes que la moyenne. Lorsque le load balancing est mal fait, un des serveurs pourrais recevoir plus de trafic que les autres ce qui provoquerait un goulots dans le réseau qui entrainerait la latence de queue.
4. 
5. health est utilisé pour checker si le serveur est toujours en cours de fonctionnement
6. Si les retries ne sont pas controlé, cela pourrais entrainé la non disponibilité du service
7. c'est des erreurs qui déclenchent d'autres erreurs sur toute la chaine.
8. pour renvoyer au client une réponse rapide.
9. c'est une session dont la state est enregistrée sur la mémoire du serveur. c'est dangereux parceque si l'algorithme du load balancer change, le client ne vas plus garder sa session.
10. 
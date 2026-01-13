1. pour des requétes http, les connexions ne sont pas maintenu longtemps, la lenteur du serveur n'est pas proportionnelle à la durée de conexion
2. 
3. c'est le 1% des requétes qui ont une latence élevée
4. EWMA utlise la moyenne qui a été fait en t-1 pour calculer la moyenne à l'instant t. 
5. parce que les 2 serveurs dans P2C sont choisis au hasard
6. least connexions
7. Pour éviter le sticky sessions car la state de la session est enregistrée sur le serveur, et si ce serveur n'est plus healthy, le client ne va pas étre redirrigé sur un autre.
8. si tous les serveurs sont identiques, healthy et stable et pas de latence
9. 
10. le load balancer doit fail fast.


mini-project:

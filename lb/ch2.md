1. un health check 200 OK peux cacher des problémes soujacentes comme la par exemple la latence du service.
2. un backend lent vient du fait que les requétes prennent du temps pour répondre et un backend mort c'est lorsque que le backend n'arrive plus à sortir une réponse
3. 
4. par exemple un RR avec avec des temps de réponses différentes sor tous les serveurs stables signifie que un des backend est dégradé.
5. les requêtes qui étaient en cours d'exécution ne seront pas aboutit donc le client attend indéfiniement pour une réponse
6. le draining permet de stopper les nouvelles requétes et de permettre les requétes en cours d'écution sur le serveur d'aboutir avant de kill le serveur. 
7.
8. lorsque le backend devient degraded, le load balancing ne va pas fonctionner correctement.
9. parce que le load balancing doit fail fast pour protéger le system et donner une réponse à l'utilisateur pour rééssayer.
10. si le mauvais health check vient du fait que un des services qui dépend de ce service n'est pas ready pour circuler correctement du trafic

--- 

Mini-Project

correction:
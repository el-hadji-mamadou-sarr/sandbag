1. plus le tick rate est élevé, plus le traitement des 10 événements est rapide. cela vient du fait que il y'a plus l'effet bloquant et le serveur doit envoyer des événements dans un max 1/tick_rate
2. on va perdre des packets puisque dans l'interval du tick_rate certain packets ne seront pas transmis.
3. Comme il y'a plus cet effet bloquant, c'est le selectors sur le serveur qui le gére
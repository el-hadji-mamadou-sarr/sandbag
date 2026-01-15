1. c'est parce que il y'a pas la logique de séquence des packets intégré dans le protocol
2. un packet perdu n'est jamais recouvrable, alors qu'un packet en retard peux étre réagencé.
3. le timestamp c'est pour calculer la latence alors que la seq c'est pour l'ordre des packets
4. les serveurs de jeux ont bbesoin d'étre performant en terme de transmission de packet et UDP offre cette possibilité. Le TCP a beaucoup de fonctionnalité qu'un serveur de jeux n'a besoin par exemple l'ACK pour le mécanisme d'accusé de réception qui est une étape pas importante dans un serveur de jeux.
5. byte order, sequence, timestamp, payload.



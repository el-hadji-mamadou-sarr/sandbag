1. ack permet juste de savoir si tous les packets ont été reçu ça ne justifie pas que le systéme est fiable c'est juste dun moyen de vérification
2. parce dans le cas de l'ack cumulatif, l'ack n'est pas envoyé constamment ce qui permettrait de ne pas avoir les goulots d'étranglement, moins de bande passante
3. packé en vol: packé sans le header l'ack et le packé acké avec ce header pour savoir la séquence max qui n'est pas perdu.
4. Une fenêtre trop grande impliquerait que le client peut envoyer dans cette fenêtre un débit important de packets
5. Pour un jeu, on a plutôt besoin du dernier packet pour la position du joueur, ses états etc..
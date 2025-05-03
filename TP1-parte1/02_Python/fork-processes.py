import random
import time 
import sys
import os

PLAYER = 5
THROWS = 10

def player(id):
    sys.stdout.write(f"Player {id} enters the game.\n")
    points = 0
    for i in range(THROWS):
        dice = random.randint(1, 6)
        points += dice
        sys.stdout.write(f"Player {id} - Throw {i + 1}: {dice}\n")
        time.sleep(random.uniform(0.1, 0.3))
    sys.stdout.write(f"Player {id} finished with {points} points.\n")
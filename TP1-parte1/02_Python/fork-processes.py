import random
import time 
import sys
import os

PLAYER = 5
THROWS = 10
DICE_MIN = 1
DICE_MAX = 6
SLEEP_MIN = 0.1
SLEEP_MAX = 0.3
WAIT_BLOCKING = 0
THROW_INDEX_OFFSET = 1

def player(id):
    sys.stdout.write(f"Player {id} enters the game.\n")
    points = 0
    for i in range(THROWS):
        dice = random.randint(DICE_MIN, DICE_MAX)
        points += dice
        sys.stdout.write(f"Player {id} - Throw {i + THROW_INDEX_OFFSET}: {dice}\n")
        time.sleep(random.uniform(SLEEP_MIN, SLEEP_MAX))
    sys.stdout.write(f"Player {id} finished with {points} points.\n")


def main():
    processes = []

    for player_id in range(PLAYER):
        pid = os.fork()

        if pid < 0:
            sys.exit(f"Error while creating process n° {player_id}")
        if pid:
            processes.append(pid)
        else:
            player(player_id)
            os._exit(os.EX_OK)

    for pid in processes:
        os.waitpid(pid, WAIT_BLOCKING)

    print("All players have finished")


if __name__ == "__main__":
    main()
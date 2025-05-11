import random
import time
import numpy as np
from concurrent.futures import ThreadPoolExecutor

TEAM_NAMES = [
  "Boca Juniors", "River Plate", "San Lorenzo", "Ferro", "Huracan", "Velez",
  "Estudiantes (LP)", "Belgrano", "Lanus", "Talleres (C)", "Dep Espaniol",
  "San Martin (T)", "Dep Mandiyu (C)", "Rosario Central", "Independiente",
  "Racing Club", "Gimnasia (LP)", "Platense", "Argentinos", "Newells"
]
POINTS_WIN = 2
POINTS_DRAW = 1
GOAL_LAMBDA = 1.0
SIM_SLEEP_MIN = 0.1
SIM_SLEEP_MAX = 0.15
MAX_WORKERS = 10


def generate_round_robin(teams):
  teams = list(teams)
  if len(teams) % 2:
    teams.append(None)
  rounds = []
  n = len(teams)
  for _ in range(n - 1):
    pairs = []
    for i in range(n // 2):
      home, away = teams[i], teams[-1 - i]
      if home and away:
        pairs.append((home, away))
    rounds.append(pairs)
    teams = [teams[0], teams[-1]] + teams[1:-1]
  return rounds


class StandingsTable:
  def __init__(self, teams):
    template = {
      'pts': 0, 'played': 0, 'wins': 0,
      'draws': 0, 'losses': 0,
      'gf': 0, 'ga': 0, 'gd': 0
    }
    self.stats = {team: template.copy() for team in teams}

  def update(self, home, away, hg, ag):
    self._record_match(home, away, hg, ag)
    self._assign_points(home, away, hg, ag)

  def _record_match(self, home, away, hg, ag):
    for team, goals_for, goals_against in ((home, hg, ag), (away, ag, hg)):
      s = self.stats[team]
      s['played'] += 1
      s['gf'] += goals_for
      s['ga'] += goals_against
      s['gd'] = s['gf'] - s['ga']

  def _assign_points(self, home, away, hg, ag):
    if hg > ag:
      self._win(home, away)
    elif hg < ag:
      self._win(away, home)
    else:
      self._draw(home, away)

  def _win(self, winner, loser):
    self.stats[winner]['wins'] += 1
    self.stats[loser]['losses'] += 1
    self.stats[winner]['pts'] += POINTS_WIN

  def _draw(self, team1, team2):
    for team in (team1, team2):
      self.stats[team]['draws'] += 1
      self.stats[team]['pts'] += POINTS_DRAW

  def get_rankings(self):
    return sorted(
      self.stats.items(),
      key=lambda x: (x[1]['pts'], x[1]['gd'], x[1]['gf']),
      reverse=True
    )


def simulate_match(match):
  home, away = match
  time.sleep(random.uniform(SIM_SLEEP_MIN, SIM_SLEEP_MAX))
  home_goals = np.random.poisson(GOAL_LAMBDA)
  away_goals = np.random.poisson(GOAL_LAMBDA)
  return home, away, int(home_goals), int(away_goals)


def run_sequential(rounds):
  table = StandingsTable(TEAM_NAMES)
  start = time.perf_counter()
  for round_matches in rounds:
    for match in round_matches:
      result = simulate_match(match)
      table.update(*result)
  duration = time.perf_counter() - start
  print("-- Sequential Simulation --")
  print_standings(table)
  print(f"Total time: {duration:.3f} seconds")


def run_concurrent(rounds, workers=None):
  table = StandingsTable(TEAM_NAMES)
  start = time.perf_counter()
  for round_matches in rounds:
    with ThreadPoolExecutor(max_workers=workers) as executor:
      futures = [executor.submit(simulate_match, m) for m in round_matches]
      for future in futures:
        table.update(*future.result())
  duration = time.perf_counter() - start
  print("-- Concurrent Simulation --")
  print_standings(table)
  print(f"Total time: {duration:.3f} seconds")


def print_standings(table):
  header = f"{'Equipo':<20} {'Pts':>4} {'PJ':>3} {'PG':>3} {'PE':>3} {'PP':>3} {'GF':>3} {'GC':>3} {'DIF':>4}"
  print(header)
  print("-" * len(header))
  for team, stats in table.get_rankings():
    print(f"{team:<20} {stats['pts']:>4} {stats['played']:>3} {stats['wins']:>3} {stats['draws']:>3} {stats['losses']:>3} {stats['gf']:>3} {stats['ga']:>3} {stats['gd']:>4}")
  print("-" * len(header))


if __name__ == '__main__':
  schedule = generate_round_robin(TEAM_NAMES)
  run_sequential(schedule)
  print()
  run_concurrent(schedule, workers=MAX_WORKERS)
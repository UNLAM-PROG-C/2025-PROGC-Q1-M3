import threading
import time
import random
import sys

# Constants
PARKING_LOT_CAPACITY = 6
PARKING_LOT_ENTRY_CONTROL = 5
PIT_CAPACITY = 3
SERVICE_CAPACITY = 2
NUM_ASSISTANTS = 2

# Semaphores & Locks
parking_lot = threading.BoundedSemaphore(PARKING_LOT_CAPACITY)
parking_lot_entry = threading.BoundedSemaphore(PARKING_LOT_ENTRY_CONTROL)
main_lane = threading.Semaphore(1)
pit_lane = threading.Semaphore(1)
pit_area = threading.BoundedSemaphore(PIT_CAPACITY)
service_area = threading.BoundedSemaphore(SERVICE_CAPACITY)
assistants = threading.Semaphore(NUM_ASSISTANTS)
aaron_lock = threading.Lock()
charles_lock = threading.Lock()

class Car(threading.Thread):

  def __init__(self, car_id):
    super().__init__()
    self.car_id = car_id

  def run(self):
    try:
      self.arrive_and_park()
      self.inspect()
      self.transfer_to_pit()
      self.repair()
      self.transfer_to_service()
      self.oil_change()
      self.transfer_to_wash()
      self.wash()
      self.exit_shop()
    except Exception as e:
      print(f"Car {self.car_id} error: {e}")

  def arrive_and_park(self):
    print(f"Car {self.car_id}: Arrived, waiting for parking...")
    parking_lot_entry.acquire()
    parking_lot.acquire()
    main_lane.acquire()
    print(f"Car {self.car_id}: Entering lot.")
    time.sleep(random.uniform(0.1, 0.3))
    main_lane.release()

  def inspect(self):
    with aaron_lock:
      print(f"Car {self.car_id}: Inspecting by Aaron.")
      time.sleep(random.uniform(0.5, 1.0))

  def transfer_to_pit(self):
    pit_lane.acquire()
    pit_area.acquire()
    assistants.acquire()
    print(f"Car {self.car_id}: Moved to pit.")
    parking_lot.release()
    parking_lot_entry.release()
    assistants.release()
    time.sleep(random.uniform(0.1, 0.3))
    pit_lane.release()

  def repair(self):
    with charles_lock:
      print(f"Car {self.car_id}: Repair by Charles.")
      time.sleep(random.uniform(1.0, 2.0))

  def transfer_to_service(self):
    assistants.acquire()
    assistants.acquire()
    service_area.acquire()
    print(f"Car {self.car_id}: Moved to service.")
    pit_area.release()
    time.sleep(random.uniform(0.1, 0.3))

  def oil_change(self):
    print(f"Car {self.car_id}: Oil change.")
    time.sleep(random.uniform(0.5, 1.0))

  def transfer_to_wash(self):
    parking_lot.acquire()
    print(f"Car {self.car_id}: Heading to wash.")
    assistants.release()
    assistants.release()
    time.sleep(random.uniform(0.1, 0.3))

  def wash(self):
    print(f"Car {self.car_id}: Washing.")
    time.sleep(random.uniform(0.5, 1.5))
    service_area.release()

  def exit_shop(self):
    main_lane.acquire()
    print(f"Car {self.car_id}: Exiting shop.")
    time.sleep(random.uniform(0.1, 0.3))
    main_lane.release()
    parking_lot.release()
    print(f"Car {self.car_id}: Picked up.")

def simulate(n):
  cars = []
  for i in range(n):
    car = Car(i+1)
    cars.append(car)
    car.start()
  for c in cars:
    c.join()

if __name__ == '__main__':
  if len(sys.argv) != 2:
    print("Usage: python gas_monkey.py <num_cars>")
    sys.exit(1)
  simulate(int(sys.argv[1]))
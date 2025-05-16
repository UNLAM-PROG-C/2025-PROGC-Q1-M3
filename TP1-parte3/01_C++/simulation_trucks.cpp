#include <iostream>
#include <barrier>
#include <thread>
#include <string>
#include <utility>
#include <cassert>
#include <semaphore>
#include <memory>
#include <mutex> 
#include <atomic>
#include <functional>

const int SIMULATION_LIMIT = 1e9+1;

std::pair<int,int> process_args(int argc, char* argv[])
{
  try
  {
    if (argc != 3)
    {
      throw std::invalid_argument("Wrong number of arguments");
    }

    int num_trucks = std::stoi(std::string(argv[1]));
    int num_travels =std::stoi(std::string(argv[2]));

    if (num_trucks <= 0 || num_travels <= 0 || num_trucks > SIMULATION_LIMIT || num_travels > SIMULATION_LIMIT) 
    {
      throw std::invalid_argument("The number of trucks and travels must be positive integers up to " + std::to_string(SIMULATION_LIMIT));
    }

    return std::make_pair(num_trucks, num_travels);
  } catch (const std::exception& e)
  {
    std::cerr << "Error: " << e.what() << std::endl;
    std::cerr << "Usage: " << argv[0] << " <num_trucks> <num_travels>" << std::endl;
    exit(1);
  }
  assert(false);
}

const int MIN_TRAVEL_TIME = 18, MAX_TRAVEL_TIME = 24;
const int GAS_LOAD_TIME = 1, FERNANDEZ_GAS_STATIONS = 2;
const int LOAD_TIME = 2, UNLOAD_TIME = 2;

std::counting_semaphore<SIMULATION_LIMIT> tapiales_travels(0), fernandez_travels(0);
std::binary_semaphore tapiales_load(1), fernandez_load(1), tapiales_unload(1), fernandez_unload(1);
std::counting_semaphore<FERNANDEZ_GAS_STATIONS> fernandez_gas_station(FERNANDEZ_GAS_STATIONS);
std::mutex cout_mutex;

std::atomic<int> hours_passed(0);

void advance_hour()
{
  hours_passed++;
  std::lock_guard<std::mutex> cout_lock(cout_mutex);
  std::cout << std::string(50,'=');
  std::cout << "[HOUR " << hours_passed << "]" ;
  std::cout << std::string(50,'=') << std::endl;
}

using CompletitionFunction = std::function<void()>;

std::shared_ptr<std::barrier<CompletitionFunction>> clock_barrier = nullptr;

void simulate_time_passage(int hours_to_pass, std::string message = "")
{
  for(int i = 0; i < hours_to_pass; i++)
  {

    cout_mutex.lock();
    // std::cout<< "[HOUR "<<hours_passed<<" ] \t" ;
    std::cout<< message << std::endl;
    cout_mutex.unlock();

    clock_barrier->arrive_and_wait();
  }
}


void load_in_tapiales(int truck_id)
{
  while(!tapiales_load.try_acquire())
  {
    simulate_time_passage(1, "Truck " + std::to_string(truck_id) + " is waiting to load in Tapiales");
  }
  simulate_time_passage( LOAD_TIME, "Truck " + std::to_string(truck_id) + " is loading in Tapiales");
  tapiales_load.release();
}

void travel_from_tapiales_to_fernandez(int truck_id)
{
  int travel_time = rand() % (MAX_TRAVEL_TIME - MIN_TRAVEL_TIME + 1) + MIN_TRAVEL_TIME;
  simulate_time_passage( travel_time, "Truck " + std::to_string(truck_id) + " is traveling from Tapiales to Fernandez");
}   

void unload_in_fernandez(int truck_id)
{
  while(!fernandez_unload.try_acquire())
  {
    simulate_time_passage( 1, "Truck " + std::to_string(truck_id) + " is waiting to unload in Fernandez");
  }
  simulate_time_passage( UNLOAD_TIME, "Truck " + std::to_string(truck_id) + " is unloading in Fernandez");
  fernandez_unload.release();
}

void tapiales_to_fernandez(int truck_id)
{
  load_in_tapiales(truck_id);
  travel_from_tapiales_to_fernandez(truck_id);
  unload_in_fernandez(truck_id);
}

void load_in_fernandez(int truck_id)
{
  while(!fernandez_load.try_acquire()){
    simulate_time_passage(1, "Truck " + std::to_string(truck_id) + " is waiting to load in Fernandez");
  }
  simulate_time_passage( LOAD_TIME, "Truck " + std::to_string(truck_id) + " is loading in Fernandez");
  fernandez_load.release();
}

void load_gas_in_fernandez(int truck_id)
{
  while(!fernandez_gas_station.try_acquire())
  {
    simulate_time_passage( 1, "Truck " + std::to_string(truck_id) + " is waiting to load gas in Fernandez");
  }
  simulate_time_passage( GAS_LOAD_TIME, "Truck " + std::to_string(truck_id) + " is loading gas in Fernandez");
  fernandez_gas_station.release();
}

void travel_from_fernandez_to_tapiales(int truck_id)
{
  int travel_time = rand() % (MAX_TRAVEL_TIME - MIN_TRAVEL_TIME + 1) + MIN_TRAVEL_TIME;
  simulate_time_passage( travel_time, "Truck " + std::to_string(truck_id) + " is traveling from Fernandez to Tapiales");
}

void unload_in_tapiales(int truck_id)
{
  while(!tapiales_unload.try_acquire())
  {
    simulate_time_passage(1, "Truck " + std::to_string(truck_id) + " is waiting to unload in Tapiales");
  }
  simulate_time_passage( UNLOAD_TIME, "Truck " + std::to_string(truck_id) + " is unloading in Tapiales");
  tapiales_unload.release();
}

void fernandez_to_tapiales(int truck_id)
{
  load_in_fernandez(truck_id);
  load_gas_in_fernandez(truck_id);
  travel_from_fernandez_to_tapiales(truck_id);
  unload_in_tapiales(truck_id);
}

void truck_simulation_core(int truck_id)
{
  while(true)
  {
    if(! tapiales_travels.try_acquire())
    {
      break;
    }
    tapiales_to_fernandez(truck_id);
        
    if(! fernandez_travels.try_acquire())
    {
      break;
    }
        
    fernandez_to_tapiales(truck_id);
  }
}

int truck_simulation(int truck_id)
{
  cout_mutex.lock();
  std::cout << "Truck " << truck_id << " is starting its travels" << std::endl;
  cout_mutex.unlock();
  
  clock_barrier->arrive_and_wait();
  
  truck_simulation_core(truck_id);
  
  cout_mutex.lock();
  std::cout << "Truck " << truck_id << " has finished its travels" << std::endl;
  cout_mutex.unlock();
  
  clock_barrier->arrive_and_drop();
  
  return hours_passed;
}

void run_main_simulation(int num_trucks){
    std::thread trucks[num_trucks];
    for(int truck_id = 1; truck_id <= num_trucks; truck_id++){
        trucks[truck_id-1] = std::thread(truck_simulation, truck_id);
    }

    for(int i = 0; i < num_trucks; i++){
        trucks[i].join();
    }

    std::cout << "Simulation finished" << std::endl;
}

int main_simulation(int num_trucks, int num_travels){
    std::cout<< "Starting simulation with " << num_trucks << " trucks and " << num_travels << " travels" << std::endl;
    tapiales_travels.release(num_travels);
    fernandez_travels.release(num_travels);
    
    clock_barrier = std::make_shared<std::barrier<CompletitionFunction>>(num_trucks, advance_hour);
    
    run_main_simulation(num_trucks);
    
    std::cout << "All trucks have finished their travels in " << hours_passed << " hours "<< std::endl;
    const int HOURS_IN_A_DAY = 24;
    std::cout << "This is "<<hours_passed/HOURS_IN_A_DAY << " days and "<< hours_passed % HOURS_IN_A_DAY << " hours " << std::endl;
    return 0;
}

int main(int argc, char* argv[]){
    auto [num_trucks, num_travels] = process_args(argc, argv);
    main_simulation(num_trucks, num_travels);
    return 0;
}


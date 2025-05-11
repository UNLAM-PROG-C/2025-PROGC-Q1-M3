use std::thread; 
use rand::{self,Rng};

const MIN_CHAKRA:u32 = 5;
const MAX_CHAKRA:u32 = 10;
const MIN_DURATION_S:f32 = 0.1;
const MAX_DURATION_S:f32 = 0.2;
const SUCCESS_RATE:f32 = 0.5;

fn clone() -> u32 {
    let chakra = rand::rng().random_range(MIN_CHAKRA..=MAX_CHAKRA);
    let mut succeses = 0u32;

    for _ in 0..chakra{
        let duration = rand::rng().random_range(MIN_DURATION_S..=MAX_DURATION_S);
        thread::sleep(std::time::Duration::from_secs_f32(duration));
        if rand::random::<f32>() < SUCCESS_RATE {
            succeses += 1;
        }
    }
    succeses
}

fn thread_work(clones: u32) -> u32 {
    let mut successes = 0;
    for _ in 0..clones {
        successes += clone();
    }
    successes
}

fn get_enviorment() -> (usize, usize){
    let num_clones = std::env::var("NUM_CLONES")
    .expect("NUM_CLONES not set")
    .parse::<usize>()
    .expect("NUM_CLONES must be a positive integer");

    let num_threads = std::env::var("NUM_THREADS").
    expect("NUM_THREADS not set")
    .parse::<usize>()
    .expect("NUM_THREADS must be a positive integer");

    (num_clones, num_threads)
}


fn thread_master() -> u32{
    let (num_clones, num_threads) = get_enviorment();

    (0..num_clones).collect::<Vec<_>>()
    .chunks(num_threads)
    .into_iter()
    .map(|chunk|{
        let len = chunk.len();
        thread::spawn(move || thread_work(len as u32) )           
    }).map(|thread| thread.join().unwrap())
    .reduce(|acc, sum| acc + sum)
    .unwrap()
}

fn main(){
    let time_0 = std::time::Instant::now();
    let level = thread_master();
    let time_1 = std::time::Instant::now();
    let elapsed = time_1.duration_since(time_0);
    println!("Level: {}", level);
    println!("Time elapsed in seconds: {:?}", elapsed.as_secs_f32());
}
/*
    Lo óptimo es que la cantidad de clones sea un múltiplo de la cantidad de hilos, para que cada hilo tenga la misma cantidad de clones.
    Si no es así, el último hilo tendrá menos clones que los demás y por tanto no se aprovechará completamente la capacidad de procesamiento.
    Como tanto el tiempo consumido como el nivel alcanzado son lineales respecto a la cantidad de clones, ninguna cantidad de clones es mejor que otra,
    salvo por lo comentado sobre la relación con la cantidad de hilos.

    Esto asumiendo que los hilos corresponde a computo realmente paralelo. En caso contrario, no habría una diferencia de eficiencia entre
    distintas cantidades de clones.
*/
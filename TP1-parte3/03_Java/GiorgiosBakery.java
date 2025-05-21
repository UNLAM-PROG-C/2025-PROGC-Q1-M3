import java.util.concurrent.*;
import java.util.*;

public class GiorgiosBakery 
{

  private static final int TABLE_CAPACITY = 20;
  private static final int BASKET_CAPACITY = 50;
  private static final int COUNTER_CAPACITY = 30;

  private static final BlockingQueue<String> doughTable = new ArrayBlockingQueue<>(TABLE_CAPACITY);
  private static final BlockingQueue<String> breadBasket = new ArrayBlockingQueue<>(BASKET_CAPACITY);
  private static final BlockingQueue<String> counter = new ArrayBlockingQueue<>(COUNTER_CAPACITY);

  private static final Semaphore scale = new Semaphore(1);
  private static final Semaphore labelMachines = new Semaphore(2);

  private static final Random random = new Random();
  private static int totalSales = 0;
  private static final Object salesLock = new Object();

  private static int totalClients;
  private static final Object clientLock = new Object();
  private static int clientsServed = 0;

  private static final Object doughTableLock = new Object();

  private static final String[] names = {"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank"};

  public static void main(String[] args) 
  {
    if (args.length != 1) 
    {
      System.out.println("Usage: java GiorgiosBakery <number_of_clients>");
      return;
    }

    totalClients = Integer.parseInt(args[0]);

    Thread master = new Thread(new Kneader("Master", 2));
    Thread assistant = new Thread(new Kneader("Assistant", 1));
    Thread baker = new Thread(new Baker());
    Thread packer1 = new Thread(new Packer("Packer 1"));
    Thread packer2 = new Thread(new Packer("Packer 2"));

    master.start();
    assistant.start();
    baker.start();
    packer1.start();
    packer2.start();

    List<Thread> clientThreads = new ArrayList<>();
    for (int i = 1; i <= totalClients; i++) 
    {
      Thread client = new Thread(new Client(i));
      client.start();
      clientThreads.add(client);
    }

    for (Thread client : clientThreads) 
    {
      try 
      {
        client.join();
      } catch (InterruptedException e) 
      {
        Thread.currentThread().interrupt();
      }
    }

    System.out.println("All clients served. Bakery is closing.");
    System.exit(0);
  }

  static class Kneader implements Runnable 
  {
    private final String name;
    private final int doughBallsPerTurn;

    public Kneader(String name, int doughBallsPerTurn) 
    {
      this.name = name;
      this.doughBallsPerTurn = doughBallsPerTurn;
    }

    @Override
    public void run() 
    {
      try 
      {
        while (true)
        {
          synchronized (clientLock)
          {
            if (clientsServed >= totalClients) break;
          }
          System.out.println(name + " is kneading: " + doughBallsPerTurn + " dough ball(s)");
          Thread.sleep(5000);

          synchronized (doughTableLock) 
          {
            while (doughTable.size() > TABLE_CAPACITY - this.doughBallsPerTurn) {
              doughTableLock.wait();
            }

            for (int i = 0; i < doughBallsPerTurn; i++) 
            {
              doughTable.put("Dough Ball");
              doughTableLock.notifyAll();
            }
            System.out.println(name + " put " + doughBallsPerTurn + " dough ball(s). Balls in Table: " + doughTable.size());
          }
        }
      } catch (InterruptedException e) 
      {
        Thread.currentThread().interrupt();
      }
    }
  }

  static class Baker implements Runnable 
  {
    @Override
    public void run() 
    {
      try 
      {
        while (true)
        {
          synchronized (clientLock)
          {
            if (clientsServed >= totalClients && doughTable.size() < 5) break;
          }
          synchronized (doughTableLock) 
          {
            while (doughTable.size() < 5) 
            {
              doughTableLock.wait();
            }
            List<String> batch = new ArrayList<>();
            for (int i = 0; i < 5; i++) {
              batch.add(doughTable.take());
            }
            System.out.println("Baker took 5 dough balls. Remaining in Table: " + doughTable.size());
            Thread.sleep(1000);
            doughTableLock.notifyAll();
          }
          System.out.println("Baker baking batch of 5 dough balls...");
          Thread.sleep(10000);

          synchronized (breadBasket)
          {
            while (breadBasket.remainingCapacity() < 5)
            {
              System.out.println("Oven blocked. Waiting for basket space...");
              breadBasket.wait();
            }
            for (int i = 0; i < 5; i++)
            {
              breadBasket.put("Bread");
            }
            System.out.println("Batch baked. Basket: " + breadBasket.size());
            breadBasket.notifyAll();
          }
        }
      } catch (InterruptedException e) 
      {
        Thread.currentThread().interrupt();
      }
    }
  }

  static class Packer implements Runnable 
  {
    private final String name;

    public Packer(String name) 
    {
      this.name = name;
    }

    @Override
    public void run() 
    {
      try 
      {
        while (true)
        {
          synchronized (clientLock)
          {
            if (clientsServed >= totalClients && breadBasket.size() < 3) break;
          }

          List<String> packageItems = new ArrayList<>();
          for (int i = 0; i < 3; i++)
          {
            packageItems.add(breadBasket.take());
          }
          System.out.println(name + " took 3 breads. Remaining in Basket: " + breadBasket.size());

          synchronized (breadBasket)
          {
            breadBasket.notifyAll();
          }

          scale.acquire();
          System.out.println(name + " weighing package...");
          Thread.sleep(1000);
          scale.release();

          labelMachines.acquire();
          System.out.println(name + " labeling package...");
          Thread.sleep(1000);
          labelMachines.release();

          counter.put("Package");
          System.out.println(name + " placed package. Counter: " + counter.size());
        }
      } catch (InterruptedException e) {
        Thread.currentThread().interrupt();
      }
    }
  }

  static class Client implements Runnable 
  {
    private final int id;
    private final String name;
    
    public Client(int id) 
    {
      this.id = id;
      this.name = names[random.nextInt(names.length)];
    }

    @Override
    public void run() 
    {
      try {
        int packages = 1 + random.nextInt(3);
        for (int i = 0; i < packages; i++)
        {
          counter.take();
          System.out.println("Client " + id + " (" + name + ") took 1 package. Remaining in Counter: " + counter.size());
        }
        synchronized (salesLock)
        {
          totalSales += packages;
        }
        synchronized (clientLock)
        {
          clientsServed++;
        }
        System.out.println("Client " + id + " (" + name + ") bought " + packages + " package(s). Total sales: " + totalSales);
      } catch (InterruptedException e) {
        Thread.currentThread().interrupt();
      }
    }
  }
}

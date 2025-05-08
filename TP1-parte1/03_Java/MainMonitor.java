import java.io.*;
import java.util.*;

public class MainMonitor 
{
  private static final String[] ZONES = 
  {
    "Basement", "Attic", "Kitchen", 
    "Bedroom", "Garden", "Mausoleum"
  };

  public static void main(String[] args) throws IOException, InterruptedException 
  {
    if (args.length != 2) 
    {
      System.out.println("Usage: java MainMonitor <duration> <interval>");
      return;
    }

    int duration = Integer.parseInt(args[0]);
    int interval = Integer.parseInt(args[1]);

    List<Process> processes = new ArrayList<>();
    List<BufferedReader> readers = new ArrayList<>();

    for (String zone : ZONES) 
    {
      ProcessBuilder builder = new ProcessBuilder
      (
        "java", "CameraProcess", zone, 
        String.valueOf(duration), 
        String.valueOf(interval)
      );
      builder.redirectErrorStream(true);
      Process process = builder.start();
      processes.add(process);
      readers.add(new BufferedReader(new InputStreamReader(process.getInputStream())));
    }

    boolean finishedCameras = false;
    while (!finishedCameras) 
    {
      finishedCameras = true;
      for (int i = 0; i < readers.size(); i++) 
      {
        BufferedReader reader = readers.get(i);
        if (reader.ready()) 
        {
          String line = reader.readLine();
          if (line != null) 
          {
            System.out.println(line);
            finishedCameras = false;
          }
        } else 
        {
          if (processes.get(i).isAlive()) 
          {
            finishedCameras = false;
          }
        }
      }
    }

    for (Process process : processes) 
    {
      process.waitFor();
    }

    System.out.println("All camera processes have finished.");
  }
}

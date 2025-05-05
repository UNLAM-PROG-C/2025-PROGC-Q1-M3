import java.util.Random;

public class CameraProcess {

  public static final int NO_ACTIVITY = 0;
  public static final int MOVEMENT_DETECTED = 1;
  public static final int THERMAL_ANOMALY = 2;
  public static final int STRANGE_SHADOW = 3;
  public static final int NOISE_DETECTED = 4;

  public static final String[] EVENT_NAMES = {
    "No Activity",
    "Movement Detected",
    "Thermal Anomaly",
    "Strange Shadow",
    "Noise Detected"
  };

  public static void main(String[] args) {
    if (args.length != 3) {
      System.out.println("Usage: java CameraProcess <zone> <duration> <interval>");
      return;
    }

    String zone = args[0];
    int durationSeconds = Integer.parseInt(args[1]);
    int intervalSeconds = Integer.parseInt(args[2]);

    int eventCount = 0;
    Random randomEvent = new Random();

    long endTime = System.currentTimeMillis() + (durationSeconds * 1000);

    while (System.currentTimeMillis() < endTime) {
      int eventId = randomEvent.nextInt(5);
      System.out.printf("[CAMERA-%s] Zone: %s | Event: %s%n", 
        ProcessHandle.current().pid(), 
        zone, 
        EVENT_NAMES[eventId]);
      if (eventId != NO_ACTIVITY) {
        eventCount++;
      }
      try {
        Thread.sleep(intervalSeconds * 1000);
      } catch (InterruptedException e) {
        e.printStackTrace();
      }
    }

    System.out.printf("[CAMERA-%s] Zone: %s | Paranormal events: %d%n", 
      ProcessHandle.current().pid(), 
      zone, 
      eventCount);
  }
}

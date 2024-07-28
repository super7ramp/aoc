import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Arrays;
import java.util.List;
import java.util.stream.IntStream;

void main() throws IOException {
    final List<String> lines = Files.readAllLines(Path.of("input.txt"));
    final long[] times = parseLine(lines.getFirst());
    final long[] distances = parseLine(lines.getLast());
    System.out.println("Times: " + Arrays.toString(times));
    System.out.println("Distances: " + Arrays.toString(distances));

    final long numberOfWays = numberOfWays(times, distances);
    System.out.println("Number of ways to beat the record (part 1): " + numberOfWays);

    final long numberOfWaysPart2 = numberOfWays(new long[]{58819676L}, new long[]{434104122191218L});
    System.out.println("Number of ways to beat the record (part 2): " + numberOfWaysPart2);
}

private static long[] parseLine(String lines) {
    return Arrays.stream(lines.split(" +"))
            .map(String::trim)
            .skip(1)
            .mapToLong(Long::parseLong)
            .toArray();
}

private long numberOfWays(long[] times, long[] distances) {
    return IntStream.range(0, times.length)
            .mapToLong(i -> solutionCount(times[i], distances[i]))
            .reduce((a, b) -> a * b)
            .orElseThrow();
}

private long solutionCount(long raceTimeInMs, long minimumDistanceInMm) {
    long solutionCount = 0;
    for (int holdDurationInMs = 1; holdDurationInMs < raceTimeInMs; holdDurationInMs++) {
        final long distanceCoveredInMm = (raceTimeInMs - holdDurationInMs) * holdDurationInMs;
        if (distanceCoveredInMm > minimumDistanceInMm) {
            solutionCount++;
        }
    }
    return solutionCount;
}
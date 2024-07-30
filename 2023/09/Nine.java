import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Arrays;
import java.util.List;

void main() throws IOException {
    final List<String> lines = Files.readAllLines(Path.of("input.txt"));

    final int nextValueSum = lines.stream()
            .map(this::parseSeries)
            .mapToInt(this::extrapolateNextValue)
            .sum();
    System.out.println("Sum of next values: " + nextValueSum);

    final int previousValueSum = lines.stream()
            .map(this::parseSeries)
            .mapToInt(this::extrapolatePreviousValue)
            .sum();
    System.out.println("Sum of previous values: " + previousValueSum);
}

private int[] parseSeries(final String s) {
    final String[] entries = s.split(" +");
    return Arrays.stream(entries)
            .mapToInt(Integer::parseInt)
            .toArray();
}

private int extrapolateNextValue(final int[] values) {
    final var differences = new int[values.length - 1];
    for (int i = 0; i < values.length - 1; i++) {
        differences[i] = values[i + 1] - values[i];
    }
    final int nextDifference;
    if (Arrays.stream(differences).anyMatch(i -> i != 0)) {
        nextDifference = extrapolateNextValue(differences);
    } else {
        nextDifference = 0;
    }
    return values[values.length - 1] + nextDifference;
}

private int extrapolatePreviousValue(int[] values) {
    final var differences = new int[values.length - 1];
    for (int i = 0; i < values.length - 1; i++) {
        differences[i] = values[i + 1] - values[i];
    }
    final int previousDifference;
    if (Arrays.stream(differences).anyMatch(i -> i != 0)) {
        previousDifference = extrapolatePreviousValue(differences);
    } else {
        previousDifference = 0;
    }
    return values[0] - previousDifference;
}
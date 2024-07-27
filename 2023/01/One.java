import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.stream.Stream;

void main() throws IOException {
    try (final Stream<String> lines = Files.lines(Path.of("input.txt"))) {
        final int sum = lines.map(this::prepareLine)
                .mapToInt(this::findCalibrationValue)
                .sum();
        System.out.println(sum);
    }
}

private String prepareLine(final String line) {
    return line
            .replaceAll("one", "o1e")
            .replaceAll("two", "t2o")
            .replaceAll("three", "t3e")
            .replaceAll("four", "4")
            .replaceAll("five", "5e")
            .replaceAll("six", "6")
            .replaceAll("seven", "7n")
            .replaceAll("eight", "8t")
            .replaceAll("nine", "n9e");
}

private int findCalibrationValue(final String line) {
    final int[] digits = line.chars()
            .filter(Character::isDigit)
            .map(Character::getNumericValue)
            .toArray();
    final int calibrationValue = digits[0] * 10 + digits[digits.length - 1];
    System.out.println(line + " -> " + calibrationValue);
    return calibrationValue;
}

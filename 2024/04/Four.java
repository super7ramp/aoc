import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.regex.Pattern;

static final Pattern XMAS = Pattern.compile("XMAS");
static final Pattern SAMX = Pattern.compile("SAMX");

void main() throws IOException {
    final String input = input();
    final long rowMatchCount = countXmases(input);
    System.out.println("Row match count: " + rowMatchCount);

    final String transposed = transpose(input);
    final long columnMatchCount = countXmases(transposed);
    System.out.println("Tranposed:\n" + transposed);
    System.out.println("Column match count: " + columnMatchCount);

    final String rotatedPlus45 = rotatePlus45(input);
    final long descendingDiagonalMatchCount = countXmases(rotatedPlus45);
    System.out.println("Rotated + 45°:\n" + rotatedPlus45);
    System.out.println("Descending diagonal match count: " + descendingDiagonalMatchCount);

    final String rotatedLess45 = rotateLess45(input);
    final long ascendingDiagonalMatchCount = countXmases(rotatedLess45);
    System.out.println("Rotated - 45°:\n" + rotatedLess45);
    System.out.println("Ascending diagonal match count: " + ascendingDiagonalMatchCount);

    final long total = rowMatchCount + columnMatchCount + descendingDiagonalMatchCount + ascendingDiagonalMatchCount;
    System.out.println("Total xmas count: " + total);
}

String input() throws IOException {
    final var inputPath = Path.of("input.txt");
    return Files.readString(inputPath);
}

/**
 * Transposes the input string, i.e. swaps rows and columns, i.e. rotates by 90° clockwise.
 * <pre>
 *     ABC         ADG
 *     DEF   -->   BEH
 *     GHI         CFI
 * </pre>
 *
 * @param input the input string
 * @return the transposed string
 */
String transpose(final String input) {
    final var transposed = new StringBuilder();
    final String[] lines = input.split("\n");
    for (int i = 0; i < lines[0].length(); i++) {
        for (final String line : lines) {
            transposed.append(line.charAt(i));
        }
        transposed.append("\n");
    }
    return transposed.toString();
}

/**
 * Rotates the input string by 45° clockwise.
 * <pre>
 *     ABC         A
 *     DEF   -->   DB
 *     GHI         GEC
 *                 HF
 *                 I
 * </pre>
 *
 * @param input the input string
 * @return the rotated string
 */
String rotatePlus45(final String input) {
    final String[] lines = input.split("\n");
    final var rotated = new StringBuilder();
    // top left corner + center
    // 0,0
    // 0,1 1,0
    // 0,2 1,1 2,0
    for (int i = 0; i < lines.length; i++) {
        for (int j = 0; j <= i; j++) {
            rotated.append(lines[j].charAt(i - j));
        }
        rotated.append("\n");
    }
    // bottom right corner
    // 1,2 2,1
    // 2,2
    for (int i = 1; i < lines.length; i++) {
        for (int j = i; j < lines.length; j++) {
            rotated.append(lines[j].charAt(lines.length - 1 - j + i));
        }
        rotated.append("\n");
    }
    return rotated.toString();
}

/**
 * Rotates the input string by 45° counter-clockwise.
 * <pre>
 *     ABC         C
 *     DEF   -->   BF
 *     GHI         AEI
 *                 DH
 *                 G
 * </pre>
 *
 * @param input the input string
 * @return the rotated string
 */
String rotateLess45(final String input) {
    final var rotated = new StringBuilder();
    final String[] lines = input.split("\n");
    // top right corner + center
    // 2,0
    // 1,0 2,1
    // 0,0 1,1 2,2
    for (int i = 0; i < lines.length; i++) {
        for (int j = 0; j <= i; j++) {
            rotated.append(lines[i - j].charAt(lines.length - 1 - j));
        }
        rotated.append("\n");
    }
    // bottom left corner
    // 0,1 1,2
    // 0,2
    for (int i = 1; i < lines.length; i++) {
        for (int j = i; j < lines.length; j++) {
            rotated.append(lines[j].charAt(j - i));
        }
        rotated.append("\n");
    }
    return rotated.toString();
}

long countXmases(final String input) {
    // Not using the same matcher with XMAS|SAMX because it would miss overlapping matches
    return XMAS.matcher(input).results().count() + SAMX.matcher(input).results().count();
}
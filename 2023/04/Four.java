import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

private static final Pattern CARD_LINE = Pattern.compile("^Card +\\d+:(?<winning>( +\\d{1,2})+) \\|(?<actual>( +\\d{1,2})+)$");
private static final Pattern NUMBER = Pattern.compile("\\d{1,2}");

void main() throws IOException {
    final List<String> lines = Files.readAllLines(Paths.get("input.txt"));
    final long[] matchCounts = lines.stream()
            .mapToLong(this::numberMatchCount)
            .toArray();
    final long totalWorth = Arrays.stream(matchCounts)
            .map(count -> count == 0 ? 0 : 1L << (count - 1))
            .sum();
    System.out.println("Total worth (part 1): " + totalWorth);

    final long[] cardCounts = new long[matchCounts.length];
    Arrays.fill(cardCounts, 1);
    for (int card = 0; card < matchCounts.length; card++) {
        final long count = matchCounts[card];
        for (int wonCard = card + 1; wonCard < card + 1 + count; wonCard++) {
            cardCounts[wonCard] += cardCounts[card];
        }
    }
    final long cardCountsSum = Arrays.stream(cardCounts).sum();
    System.out.println("Total scratchcards (part 2): " + cardCountsSum);
}

private long numberMatchCount(final String line) {
    final Matcher matcher = CARD_LINE.matcher(line);
    if (!matcher.matches()) {
        throw new IllegalArgumentException("Invalid line: " + line);
    }

    final String actual = matcher.group("actual");
    final List<Integer> actualNumbers = NUMBER.matcher(actual)
            .results()
            .map(result -> Integer.valueOf(result.group()))
            .toList();

    final String winning = matcher.group("winning");
    return NUMBER.matcher(winning)
            .results()
            .map(result -> Integer.valueOf(result.group()))
            .filter(actualNumbers::contains)
            .count();
}

